/* Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

//TODO refactor this class to share code with jsonrpc connector

package client

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/common"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/rest"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

type RestConnector struct {
	url        string
	httpClient http.Client
	options    *connectorOptions
	provider   core.APIProvider

	statusCode int
}

type connectorOptions struct {
	decorators                   []core.APIProviderDecorator
	securityContext              core.SecurityContext
	appContext                   *core.ApplicationContext
	connectionMetadata           map[string]interface{}
	enableDefaultContentType     bool
	requestProcessors            []core.RequestProcessor
	responseAcceptors            []core.ResponseAcceptor
	typeConverter                *bindings.TypeConverter
	securityContextSerializerMap map[string]rest.SecurityContextSerializer
}

type ConnectorOption func(*connectorOptions)

func NewRestConnector(url string, client http.Client, options ...ConnectorOption) *RestConnector {
	connectorOptions := defaultRESTConnectorOptions()

	for _, o := range options {
		o(connectorOptions)
	}

	connector := &RestConnector{
		url:        url,
		httpClient: client,
		options:    connectorOptions,
	}
	connector.provider = connector

	chainDecorators(connector)

	return connector
}

func chainDecorators(connector *RestConnector) {
	for _, decorator := range connector.options.decorators {
		connector.provider = decorator(connector.provider)
	}
}

func defaultRESTConnectorOptions() *connectorOptions {
	return &connectorOptions{
		securityContextSerializerMap: map[string]rest.SecurityContextSerializer{
			security.USER_PASSWORD_SCHEME_ID: rest.NewUserPwdSecContextSerializer(),
			security.SESSION_SCHEME_ID:       rest.NewSessionSecContextSerializer(),
			security.OAUTH_SCHEME_ID:         rest.NewOauthSecContextSerializer(),
		},
		typeConverter:            bindings.NewTypeConverter(bindings.InRestMode()),
		enableDefaultContentType: true,
	}
}

// WithDecorators sets decorating APIProvider.
// APIProvider decorators wrap execution of Invoke method to extend runtime with additional functionality.
func WithDecorators(decorators ...core.APIProviderDecorator) ConnectorOption {
	return func(options *connectorOptions) {
		options.decorators = append(options.decorators, decorators...)
	}
}

// WithRequestProcessors defines request processors for RestConnector
func WithRequestProcessors(requestProcessors ...core.RequestProcessor) ConnectorOption {
	return func(options *connectorOptions) {
		options.requestProcessors = requestProcessors
	}
}

// WithResponseAcceptors defines response acceptors for RestConnector
func WithResponseAcceptors(responseAcceptors ...core.ResponseAcceptor) ConnectorOption {
	return func(options *connectorOptions) {
		options.responseAcceptors = responseAcceptors
	}
}

func (j *RestConnector) ApplicationContext() *core.ApplicationContext {
	return j.options.appContext
}

func (j *RestConnector) StatusCode() int {
	return j.statusCode
}

func (j *RestConnector) SetApplicationContext(ctx *core.ApplicationContext) {
	j.options.appContext = ctx
}

func (j *RestConnector) SecurityContext() core.SecurityContext {
	return j.options.securityContext
}

func (j *RestConnector) SetSecurityContext(ctx core.SecurityContext) {
	j.options.securityContext = ctx
}

func (j *RestConnector) SetConnectionMetadata(connectionMetadata map[string]interface{}) {
	j.options.connectionMetadata = connectionMetadata
}

func (j *RestConnector) ConnectionMetadata() map[string]interface{} {
	return j.options.connectionMetadata
}

// If enableDefaultContentType is True then Header[Content-Type] gets overwritten to value 'application/json'
func (j *RestConnector) SetEnableDefaultContentType(enableDefaultContentType bool) {
	j.options.enableDefaultContentType = enableDefaultContentType
}

func (j *RestConnector) SetSecCtxSerializer(schemeID string, serializer rest.SecurityContextSerializer) {
	j.options.securityContextSerializerMap[schemeID] = serializer
}

func (j *RestConnector) SecurityContextSerializerMap() map[string]rest.SecurityContextSerializer {
	return j.options.securityContextSerializerMap
}

// AddRequestProcessor adds request processor to connector.
// Request processors are executed right before request is made to the server
// Deprecated: use WithRequestProcessors instead
func (j *RestConnector) AddRequestProcessor(processor rest.RequestProcessor) {
	j.options.requestProcessors = append(j.options.requestProcessors, processor.Process)
}

// RequestProcessors gets list of all request processors defined in connector
func (j *RestConnector) RequestProcessors() []core.RequestProcessor {
	return j.options.requestProcessors
}

// ResponseAcceptors gets list of all response acceptors defined in connector
func (j *RestConnector) ResponseAcceptors() []core.ResponseAcceptor {
	return j.options.responseAcceptors
}

func (j *RestConnector) NewExecutionContext() *core.ExecutionContext {
	if j.options.appContext == nil {
		j.options.appContext = common.NewDefaultApplicationContext()
	}
	appCopy := j.options.appContext.Copy()
	common.InsertOperationId(appCopy)
	return core.NewExecutionContext(appCopy, j.options.securityContext)
}

func (j *RestConnector) GetApiProvider() core.APIProvider {
	return j.provider
}

func (j *RestConnector) TypeConverter() *bindings.TypeConverter {
	return j.options.typeConverter
}

func (j *RestConnector) buildHTTPRequest(serializedRequest *rest.Request,
	ctx *core.ExecutionContext, restMetadata *protocol.OperationRestMetadata) (*http.Request, error) {
	body := strings.NewReader(serializedRequest.RequestBody())
	url := j.url + serializedRequest.URLPath()
	method := restMetadata.HttpMethod()
	log.Debugf("Invoking action: %q and url: %q", method, url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, vlist := range serializedRequest.InputHeaders() {
		for _, v := range vlist {
			req.Header.Set(k, v)
		}
	}
	if _, ok := req.Header[lib.HTTP_CONTENT_TYPE_HEADER]; !ok && j.options.enableDefaultContentType {
		req.Header.Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	}
	req.Header.Set(lib.HTTP_USER_AGENT_HEADER, GetRuntimeUserAgentHeader())
	CopyContextsToHeaders(ctx, req)
	return req, nil
}

func (j *RestConnector) Invoke(serviceID string, operationID string,
	inputValue data.DataValue, ctx *core.ExecutionContext) core.MethodResult {

	//TODO do we need serviceID and opID for rest connector?
	if ctx == nil {
		ctx = j.NewExecutionContext()
	}

	// Get operation metadata from connector
	restMetadata, errVal := j.retrieveOperationRestMetadata()
	if errVal != nil {
		return core.NewMethodResult(nil, errVal)
	}

	var inputStructValue *data.StructValue
	if structValue, ok := inputValue.(*data.StructValue); ok {
		inputStructValue = structValue
	} else {
		err := l10n.NewRuntimeErrorNoParam("vapi.protocol.server.rest.response.not_structure")
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	securityCtx := ctx.SecurityContext()
	var secCtxSerializer rest.SecurityContextSerializer
	if securityCtx != nil {
		// Get schemeID of the securityContext
		schemeID, err := rest.GetSecurityCtxStrValue(securityCtx, security.AUTHENTICATION_SCHEME_ID)
		if err != nil {
			log.Error(err)
			err := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.request.error")
			errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		}

		if schemeID != nil {
			log.Debug("SecurityContext schemeID is ", schemeID)
			// Find the approprate SecurityContextSerializer based on the schemeID
			if serializer, ok := j.options.securityContextSerializerMap[*schemeID]; ok {
				secCtxSerializer = serializer
			} else {
				log.Debug("No appropriate SecurityContextSerializer for schemeID %s. HTTP headers will not be added to request", schemeID)
			}
		}
	}

	// Serialize urlPath, inputHeaders and requestBody
	serializedRequest, err := rest.SerializeRequestsWithSecCtxSerializers(
		inputStructValue, ctx, restMetadata, secCtxSerializer)
	if err != nil {
		err := l10n.NewRuntimeError("vapi.data.serializers.json.marshall.error",
			map[string]string{"errorMessage": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	req, err := j.buildHTTPRequest(serializedRequest, ctx, restMetadata)

	if ctx.Context() != nil {
		req = req.WithContext(ctx.Context())
	}

	// Allow client to access the req object before sending the http request
	allRequestProcessors := append(ctx.RuntimeData().GetRequestProcessors(), j.options.requestProcessors...)
	for _, preProcessor := range allRequestProcessors {
		err := preProcessor(req)
		if err != nil {
			log.Debug(err)
			err := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.request.error")
			errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		}
	}

	response, err := j.httpClient.Do(req)

	if err != nil {
		errVal := getVAPIError(err)
		return core.NewMethodResult(nil, errVal)
	}

	// Provides access to raw response coming from the server
	allResponseAcceptors := append(ctx.RuntimeData().GetResponseAcceptors(), j.options.responseAcceptors...)
	for _, responseAcceptor := range allResponseAcceptors {
		responseAcceptor(response)
	}

	return j.PrepareMethodResult(response, restMetadata)
}

func (j *RestConnector) retrieveOperationRestMetadata() (*protocol.OperationRestMetadata, *data.ErrorValue) {
	var connMetadata interface{}
	if connMeta, ok := j.ConnectionMetadata()[lib.REST_METADATA]; ok {
		connMetadata = connMeta
	} else {
		err := l10n.NewRuntimeErrorNoParam("vapi.bindings.stub.rest_metadata.unavailable")
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return nil, errVal
	}

	if restMeta, ok := connMetadata.(protocol.OperationRestMetadata); ok {
		return &restMeta, nil
	}

	err := l10n.NewRuntimeErrorNoParam("vapi.bindings.stub.rest_metadata.type.mismatch")
	errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
	return nil, errVal
}

func (j *RestConnector) PrepareMethodResult(response *http.Response, restMetadata *protocol.OperationRestMetadata) core.MethodResult {
	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
	}()
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.response.error")
		// TODO create an appropriate binding error for this
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	// assign status code
	j.statusCode = response.StatusCode

	respHeader := response.Header
	responseBody := string(resp)
	methodResult, err := rest.DeserializeResponse(response.StatusCode, respHeader, responseBody, restMetadata)
	if err != nil {
		err := l10n.NewRuntimeError("vapi.protocol.client.response.unmarshall.error", map[string]string{"responseBody": responseBody})
		// TODO create an appropriate binding error for this
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}
	return methodResult
}
