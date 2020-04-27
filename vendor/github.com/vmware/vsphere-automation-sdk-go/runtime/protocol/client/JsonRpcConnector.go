/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package client

import (
	"bytes"
	"encoding/json"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/common"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/server/rpc/msg"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type JsonRpcConnector struct {
	url                string
	httpClient         http.Client
	securityContext    core.SecurityContext
	appContext         *core.ApplicationContext
	typeConverter      *bindings.TypeConverter
	connectionMetadata map[string]interface{}
}

func NewJsonRpcConnector(url string, client http.Client) *JsonRpcConnector {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	return &JsonRpcConnector{url: url, httpClient: client, typeConverter: typeConverter}
}

func (j *JsonRpcConnector) ApplicationContext() *core.ApplicationContext {
	return j.appContext
}

func (j *JsonRpcConnector) SetApplicationContext(ctx *core.ApplicationContext) {
	j.appContext = ctx
}

func (j *JsonRpcConnector) SecurityContext() core.SecurityContext {
	return j.securityContext
}

func (j *JsonRpcConnector) SetSecurityContext(ctx core.SecurityContext) {
	j.securityContext = ctx
}

func (j *JsonRpcConnector) NewExecutionContext() *core.ExecutionContext {
	if j.appContext == nil {
		j.appContext = common.NewDefaultApplicationContext()
	} else {
		common.InsertOperationId(j.appContext)
	}
	return core.NewExecutionContext(j.appContext, j.securityContext)
}

func (j *JsonRpcConnector) GetApiProvider() core.APIProvider {
	return j
}

func (j *JsonRpcConnector) TypeConverter() *bindings.TypeConverter {
	return j.typeConverter
}

func (j *JsonRpcConnector) SetConnectionMetadata(connectionMetadata map[string]interface{}) {
	j.connectionMetadata = connectionMetadata
}

func (j *JsonRpcConnector) ConnectionMetadata() map[string]interface{} {
	return j.connectionMetadata
}

func (j *JsonRpcConnector) Invoke(serviceId string, operationId string, inputValue data.DataValue, ctx *core.ExecutionContext) core.MethodResult {
	if ctx == nil {
		ctx = j.NewExecutionContext()
	}
	if !ctx.ApplicationContext().HasProperty(lib.OPID) {
		common.InsertOperationId(ctx.ApplicationContext())
	}
	opId := ctx.ApplicationContext().GetProperty(lib.OPID)

	var params = make(map[string]interface{})
	var jsonRpcEncoder = msg.NewJsonRpcEncoder()
	var encodedInput, encodingError = jsonRpcEncoder.Encode(inputValue)
	if encodingError != nil {
		log.Error("Error encoding input")
		log.Error(encodingError)
		err := l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": encodingError.Error()})
		err_val := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, err_val)
	}
	encodedInputMap := map[string]interface{}{}
	var unmarshalError = json.Unmarshal(encodedInput, &encodedInputMap)
	if unmarshalError != nil {
		log.Error(unmarshalError)
		err := l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": unmarshalError.Error()})
		err_val := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, err_val)
	}
	log.Debugf("Invoking with input %+v", string(encodedInput))
	params[lib.REQUEST_INPUT] = encodedInputMap
	params[lib.EXECUTION_CONTEXT] = msg.NewExecutionContextSerializer(ctx)
	params[lib.REQUEST_OPERATION_ID] = operationId
	params[lib.REQUEST_SERVICE_ID] = serviceId
	var jsonRpcRequest = msg.NewJsonRpc20Request(lib.JSONRPC_VERSION, lib.JSONRPC_INVOKE, params, opId, false)
	var jsonRpcRequestSerializer = msg.NewJsonRpc20RequestSerializer(jsonRpcRequest)
	var requestBytes, marshallError = jsonRpcRequestSerializer.MarshalJSON()
	if marshallError != nil {
		log.Error(unmarshalError)
		err := l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": marshallError.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	req, nrequestErr := http.NewRequest(http.MethodPost, j.url, bytes.NewBuffer(requestBytes))
	if nrequestErr != nil {
		log.Error(nrequestErr)
		err := l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": nrequestErr.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return core.NewMethodResult(nil, errVal)
	}

	req.Header.Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	req.Header.Set(lib.VAPI_SERVICE_HEADER, serviceId)
	req.Header.Set(lib.VAPI_OPERATION_HEADER, operationId)
	req.Header.Set(lib.HTTP_USER_AGENT_HEADER, GetRuntimeUserAgentHeader())
	if j.connectionMetadata["isStreamingResponse"] == true {
		req.Header.Set(lib.HTTP_ACCEPT, lib.VAPI_STREAMING_HEADER_VALUE)
	}
	if ctx.Context() != nil {
		req = req.WithContext(ctx.Context())
	}
	CopyContextsToHeaders(ctx, req)

	if j.connectionMetadata["isStreamingResponse"] == true {
		methodResult := core.NewMethodResult(nil, nil)
		methodResultChan := make(chan core.MethodResult)
		methodResult.SetResponseStream(methodResultChan)
		go func() {
			response, requestErr := j.httpClient.Do(req)
			if requestErr != nil {
				if strings.HasSuffix(requestErr.Error(), "connection refused") {
					err := l10n.NewRuntimeErrorNoParam("vapi.server.unavailable")
					errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
					methodResultChan <- core.NewMethodResult(nil, errVal)
					return
				} else if strings.HasSuffix(requestErr.Error(), "i/o timeout") {
					err := l10n.NewRuntimeErrorNoParam("vapi.server.timedout")
					errVal := bindings.CreateErrorValueFromMessages(bindings.TIMEDOUT_ERROR_DEF, []error{err})
					methodResultChan <- core.NewMethodResult(nil, errVal)
					return
				}
			}
			defer response.Body.Close()
			dec := json.NewDecoder(httputil.NewChunkedReader(response.Body))
			dec.UseNumber()
			jsonRpcDecoder := msg.NewJsonRpcDecoder()
			for {
				var m map[string]interface{}
				if err := dec.Decode(&m); err == io.EOF {
					break
				} else if err != nil {
					// err is not nil when ctx is cancelled.
					log.Debug(err)
					break
				}
				jsonRpcResponse, responseDeserializationErr := jsonRpcDecoder.GetJsonRpc20Response(m)
				if responseDeserializationErr != nil {
					log.Error(responseDeserializationErr)
					err := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": responseDeserializationErr.Message()})
					errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{err})
					methodResultChan <- core.NewMethodResult(nil, errVal)
					break
				}
				if jsonRpcResponseMap, ok := jsonRpcResponse.Result().(map[string]interface{}); ok {
					// check for terminal frame
					if len(jsonRpcResponseMap) == 0 {
						log.Debug("Recieved terminal frame")
						break
					}
					methodResult, err := jsonRpcDecoder.DeSerializeMethodResult(jsonRpcResponseMap)
					if err != nil {
						log.Error(err)
						runtimeError := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": err.Error()})
						errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
						methodResultChan <- core.NewMethodResult(nil, errVal)
						break
					}
					methodResultChan <- methodResult
					if methodResult.Error() != nil {
						log.Debug("Recieved error frame")
						break
					}
				} else {
					runtimeError := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": ""})
					errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
					methodResultChan <- core.NewMethodResult(nil, errVal)
				}

			}
			close(methodResultChan)

		}()
		return methodResult
	}
	// non streaming responses.
	response, requestErr := j.httpClient.Do(req)
	if requestErr != nil {
		if strings.HasSuffix(requestErr.Error(), "connection refused") {
			err := l10n.NewRuntimeErrorNoParam("vapi.server.unavailable")
			errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		} else if strings.HasSuffix(requestErr.Error(), "i/o timeout") {
			err := l10n.NewRuntimeErrorNoParam("vapi.server.timedout")
			errVal := bindings.CreateErrorValueFromMessages(bindings.TIMEDOUT_ERROR_DEF, []error{err})
			return core.NewMethodResult(nil, errVal)
		}
	}
	defer response.Body.Close()
	resp, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Error(readErr)
		runtimeError := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": readErr.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
		return core.NewMethodResult(nil, errVal)
	}
	jsonRpcDecoder := msg.NewJsonRpcDecoder()
	jsonRpcResponse, deserializeError := jsonRpcDecoder.DeSerializeResponse(resp)
	//TODO
	//simplify DeserializeResponse API to return only deserialization error instead of jsonrpcerror
	if deserializeError != nil {
		log.Error(deserializeError)
		runtimeError := l10n.NewRuntimeErrorNoParam("vapi.protocol.client.response.error")
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
		return core.NewMethodResult(nil, errVal)
	}
	if jsonRpcResponseMap, ok := jsonRpcResponse.Result().(map[string]interface{}); ok {
		methodResult, err := jsonRpcDecoder.DeSerializeMethodResult(jsonRpcResponseMap)
		if err != nil {
			log.Error(err)
			runtimeError := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": err.Error()})
			errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
			return core.NewMethodResult(nil, errVal)
		}
		return methodResult
	} else {
		runtimeError := l10n.NewRuntimeError("vapi.server.response.error", map[string]string{"errMsg": ""})
		errVal := bindings.CreateErrorValueFromMessages(bindings.INTERNAL_SERVER_ERROR_DEF, []error{runtimeError})
		return core.NewMethodResult(nil, errVal)
	}
}
