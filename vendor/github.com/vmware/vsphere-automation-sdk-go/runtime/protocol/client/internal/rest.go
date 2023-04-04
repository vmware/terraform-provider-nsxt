/* Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package internal

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/rest"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

// RESTHttpProtocol provides functionality for making and handling http REST calls to a server
type RESTHttpProtocol struct {
	Connector
	*HttpTransport
	Options protocol.RestClientOptions
}

// NewRESTHttpProtocol instantiates instance of RESTHttpProtocol
func NewRESTHttpProtocol(c Connector, restOptions protocol.RestClientOptions) *RESTHttpProtocol {
	return &RESTHttpProtocol{
		Connector:     c,
		HttpTransport: NewHttpTransport(),
		Options:       restOptions,
	}
}

// SecurityContextSerializerMap gets map of all security serializers
func (r *RESTHttpProtocol) SecurityContextSerializerMap() map[string]protocol.SecurityContextSerializer {
	return r.Options.SecurityContextSerializers()
}

// Invoke from given parameters makes http request to remote server
func (r *RESTHttpProtocol) Invoke(
	serviceId string,
	operationId string,
	inputValue data.DataValue,
	ctx *core.ExecutionContext) core.MethodResult {

	request, err := r.buildRequest(serviceId, operationId, inputValue, ctx)
	if err != nil {
		vapiErr := getRequestError(err.Error())
		return getErrorMethodResult(vapiErr)
	}

	response, err := r.executeRequest(ctx, request)
	if err != nil {
		vapiErr := getVAPIError(err)
		return getErrorMethodResult(vapiErr)
	}

	return r.handleResponse(ctx, response)
}

func (r *RESTHttpProtocol) buildRequest(
	_ string,
	operationId string,
	inputValue data.DataValue,
	ctx *core.ExecutionContext) (*http.Request, error) {

	if ctx == nil {
		panic("execution context can't be nil")
	}

	// Get operation metadata from connector
	restMetadata, err := r.retrieveOperationRestMetadata(ctx)
	if err != nil {
		return nil, err
	}

	inputStructValue, ok := inputValue.(*data.StructValue)
	if !ok {
		return nil, l10n.NewRuntimeErrorNoParam("vapi.protocol.server.rest.response.not_structure")
	}

	secCtxSerializer, err := r.getSecurityContextSerializer(ctx.SecurityContext())
	if err != nil {
		return nil, err
	}

	// Serialize urlPath, inputHeaders and requestBody
	serializedRequest, err := rest.SerializeRequestsWithSecCtxSerializers(
		inputStructValue, ctx, restMetadata, secCtxSerializer)
	if err != nil {
		return nil, err
	}

	httpRequest, err := r.buildHTTPRequest(ctx, serializedRequest, restMetadata)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(operationId, lib.TaskInvocationString) {
		queryParams := httpRequest.URL.Query()
		queryParams.Add(lib.TaskRESTQueryKey, "true")
		httpRequest.URL.RawQuery = queryParams.Encode()
	}

	// fill in request with go context from bindings
	if ctx.Context() != nil {
		httpRequest = httpRequest.WithContext(ctx.Context())
	}

	return httpRequest, nil
}

func (r *RESTHttpProtocol) retrieveOperationRestMetadata(ctx *core.ExecutionContext) (*protocol.OperationRestMetadata, error) {
	contextValue, err := ctx.ConnectionMetadata(core.RESTMetadataKey)
	if err != nil {
		return nil, l10n.NewRuntimeErrorNoParam("vapi.bindings.stub.rest_metadata.unavailable")
	}
	restMetadata, ok := contextValue.(protocol.OperationRestMetadata)
	if !ok {
		err := l10n.NewRuntimeErrorNoParam("vapi.bindings.stub.rest_metadata.type.mismatch")
		return nil, err
	}
	return &restMetadata, nil
}

func (r *RESTHttpProtocol) getSecurityContextSerializer(securityContext core.SecurityContext) (protocol.SecurityContextSerializer, error) {
	if securityContext != nil {
		// Get schemeID of the securityContext
		schemeID, err := rest.GetSecurityCtxStrValue(securityContext, security.AUTHENTICATION_SCHEME_ID)
		if err != nil {
			return nil, err
		}

		if schemeID == nil {
			return nil, nil
		}

		log.Debugf("SecurityContext schemeID is: ", schemeID)
		// Find the appropriate SecurityContextSerializer based on the schemeID
		secCtxSerializer, ok := r.Options.SecurityContextSerializers()[*schemeID]
		if !ok {
			log.Debug("No appropriate SecurityContextSerializer for schemeID %s. "+
				"Security related HTTP headers will not be added to request", schemeID)
			return nil, nil
		}
		return secCtxSerializer, nil
	}
	return nil, nil
}

func (r *RESTHttpProtocol) buildHTTPRequest(
	ctx *core.ExecutionContext,
	serializedRequest *rest.Request,
	restMetadata *protocol.OperationRestMetadata) (*http.Request, error) {

	body := strings.NewReader(serializedRequest.RequestBody())
	url := r.Connector.Address() + serializedRequest.URLPath()
	method := restMetadata.HttpMethod()
	log.Debugf("Invoking action: %q and url: %q", method, url)

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for headerKey, headerValues := range serializedRequest.InputHeaders() {
		for _, headerValue := range headerValues {
			request.Header.Set(headerKey, headerValue)
		}
	}

	// if not set up to what InputHeaders point to use enabledDefaultContentType flag
	if _, ok := request.Header[lib.HTTP_CONTENT_TYPE_HEADER]; !ok && r.Options.EnableDefaultContentType() {
		request.Header.Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	}

	request.Header.Set(lib.HTTP_USER_AGENT_HEADER, GetRuntimeUserAgentHeader())

	copyContextsToHeaders(ctx, request.Header)

	return request, nil
}

func (r *RESTHttpProtocol) handleResponse(ctx *core.ExecutionContext, response *http.Response) core.MethodResult {
	defer func() {
		if response != nil && response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				panic(err)
			}
		}
	}()

	restMetadata, err := r.retrieveOperationRestMetadata(ctx)
	if err != nil {
		responseError := getResponseError(err.Error())
		return getErrorMethodResult(responseError)
	}

	return r.PrepareMethodResult(response, restMetadata)
}

func (r *RESTHttpProtocol) PrepareMethodResult(response *http.Response, restMetadata *protocol.OperationRestMetadata) core.MethodResult {
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		responseError := getResponseError(err.Error())
		return getErrorMethodResult(responseError)
	}

	responseBody := string(resp)
	methodResult, err := rest.DeserializeResponse(response.StatusCode, response.Header, responseBody, restMetadata)
	if err != nil {
		responseError := getResponseError(err.Error())
		return getErrorMethodResult(responseError)
	}
	return methodResult
}
