//TODO refactor this class to share code with jsonrpc connector
/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

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
)

type RestConnector struct {
	url                      string
	httpClient               http.Client
	securityContext          core.SecurityContext
	appContext               *core.ApplicationContext
	connectionMetadata       map[string]interface{}
	enableDefaultContentType bool
}

func NewRestConnector(url string, client http.Client) *RestConnector {
	return &RestConnector{url: url, httpClient: client, enableDefaultContentType: true}
}

func (j *RestConnector) ApplicationContext() *core.ApplicationContext {
	return j.appContext
}

func (j *RestConnector) SetApplicationContext(ctx *core.ApplicationContext) {
	j.appContext = ctx
}

func (j *RestConnector) SecurityContext() core.SecurityContext {
	return j.securityContext
}

func (j *RestConnector) SetSecurityContext(ctx core.SecurityContext) {
	j.securityContext = ctx
}

func (j *RestConnector) SetConnectionMetadata(connectionMetadata map[string]interface{}) {
	j.connectionMetadata = connectionMetadata
}

func (j *RestConnector) ConnectionMetadata() map[string]interface{} {
	return j.connectionMetadata
}

// If enableDefaultContentType is True then Header[Content-Type] gets overwritten to value 'application/json'
func (j *RestConnector) SetEnableDefaultContentType(enableDefaultContentType bool) {
	j.enableDefaultContentType = enableDefaultContentType
}

func (j *RestConnector) NewExecutionContext() *core.ExecutionContext {
	if j.appContext == nil {
		j.appContext = common.NewDefaultApplicationContext()
	} else {
		common.InsertOperationId(j.appContext)
	}
	return core.NewExecutionContext(j.appContext, j.securityContext)
}

func (j *RestConnector) GetApiProvider() core.APIProvider {
	return j
}

func (j *RestConnector) TypeConverter() *bindings.TypeConverter {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.REST)
	return typeConverter
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
	for k, v := range serializedRequest.InputHeaders() {
		req.Header.Set(k, v)
	}
	if _, ok := req.Header[lib.HTTP_CONTENT_TYPE_HEADER]; !ok && j.enableDefaultContentType {
		req.Header.Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	}
	req.Header.Set(lib.HTTP_USER_AGENT_HEADER, GetRuntimeUserAgentHeader())
	CopyContextsToHeaders(ctx, req)
	return req, nil
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

func (j *RestConnector) Invoke(serviceID string, operationID string,
	inputValue data.DataValue, ctx *core.ExecutionContext) core.MethodResult {
	//TODO do we need serviceID and opID for rest connector?
	if ctx == nil {
		ctx = j.NewExecutionContext()
	}
	if !ctx.ApplicationContext().HasProperty(lib.OPID) {
		common.InsertOperationId(ctx.ApplicationContext())
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

	// Serialize urlPath, inputHeaders and requestBody
	serializedRequest, err := rest.SerializeRequests(inputStructValue, ctx, restMetadata)
	if err != nil {
		err := l10n.NewRuntimeError("vapi.data.serializers.json.marshall.error",
			map[string]string{"errorMessage": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	req, err := j.buildHTTPRequest(serializedRequest, ctx, restMetadata)
	response, err := j.httpClient.Do(req)
	if err != nil {
		errString := err.Error()
		if strings.HasSuffix(errString, "connection refused") {
			err := l10n.NewRuntimeErrorNoParam("vapi.server.unavailable")
			errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		} else if strings.HasSuffix(errString, "i/o timeout") {
			err := l10n.NewRuntimeErrorNoParam("vapi.server.timedout")
			errVal := bindings.CreateErrorValueFromMessages(bindings.TIMEDOUT_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		} else if strings.Contains(errString, "x509") {
			err := l10n.NewRuntimeErrorNoParam("vapi.security.authentication.certificate.invalid")
			errVal := bindings.CreateErrorValueFromMessages(bindings.UNAUTHENTICATED_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		} else {
			// TODO add more specific errors
			err := l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": errString})
			errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		}
	}
	defer response.Body.Close()
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.response.error")
		// TODO create an appropriate binding error for this
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}
	methodResult, err := rest.DeserializeResponse(response.StatusCode, string(resp))
	if err != nil {
		err := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.response.unmarshall.error")
		// TODO create an appropriate binding error for this
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}
	return methodResult
}
