/* Copyright © 2019-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/common"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/server"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/server/tracing"
)

const ENABLE_VAPI_PROVIDER_WIRE_LOGGING = "ENABLE_VAPI_PROVIDER_WIRE_LOGGING"

type JsonRpcHandler struct {
	apiProvider       core.APIProvider
	jsonRpcEncoder    *JsonRpcEncoder
	jsonRpcDecoder    *JsonRpcDecoder
	requestProcessors []core.JSONRPCRequestPreProcessor
	tracer            tracing.StartSpan
}

type JsonRpcHandlerOption func(handler *JsonRpcHandler)

func NewJsonRpcHandler(apiProvider core.APIProvider, opts ...JsonRpcHandlerOption) *JsonRpcHandler {
	var jsonRpcEncoder = NewJsonRpcEncoder()
	var jsonRpcDecoder = NewJsonRpcDecoder()
	res := &JsonRpcHandler{apiProvider: apiProvider,
		jsonRpcEncoder: jsonRpcEncoder,
		jsonRpcDecoder: jsonRpcDecoder,
		tracer:         tracing.NoopTracer}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// WithTracer option allows NewJsonRpcHandler to set API server tracer that will be used to extract http headers and
// create server spans for incoming calls.
func WithTracer(tracer tracing.StartSpan) JsonRpcHandlerOption {
	return func(handler *JsonRpcHandler) {
		handler.tracer = tracer
	}
}

func (j *JsonRpcHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(rw, "Invalid request method", http.StatusMethodNotAllowed)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error reading request body", http.StatusInternalServerError)
	}
	requestObj, err := DeSerializeJson(body)
	if err != nil {
		log.Error("Error deserializing jsonrpc request")
		log.Error(err)
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest("Error deserializing jsonrpc request"), nil)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}
	var request, requestDeserializationError = getJsonRpc20Request(requestObj)
	if requestDeserializationError != nil {
		log.Error("Error deserializing jsonrpc request")
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest("Error deserializing jsonrpc request"), nil)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}
	for _, reqProcessor := range j.requestProcessors {
		err := reqProcessor.Process(&requestObj)
		if err != nil {
			log.Error("Encountered error during preprocessing of json request")
			log.Error(err)
			var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest(err), &request)
			j.sendResponse(jsonRpcRequestError, rw, nil)
			return
		}
	}
	j.processJsonRpcRequest(rw, r, request)
}

func (j *JsonRpcHandler) AddRequestPreProcessor(reqProcessor core.JSONRPCRequestPreProcessor) {
	j.requestProcessors = append(j.requestProcessors, reqProcessor)
}

func (j *JsonRpcHandler) sendResponse(response interface{}, rw http.ResponseWriter, error *data.ErrorValue) {
	var result, encodeError = j.jsonRpcEncoder.Encode(response)
	if encodeError != nil {
		log.Error(encodeError)
	}
	rw.Header().Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	if error != nil {
		//Accessing directly to avoid canonicalizing header key.
		rw.Header()[lib.VAPI_ERROR] = []string{error.Name()}
	}
	_, writeErr := rw.Write(result)
	if writeErr != nil {
		log.Error(writeErr)
	}
}

func (j *JsonRpcHandler) processJsonRpcRequest(rw http.ResponseWriter, r *http.Request, request JsonRpc20Request) {

	var serviceId, serviceIdError = j.serviceId(request)
	if serviceIdError != nil {
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidParams(serviceIdError.Error()), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	var operationId, operationIdError = j.operationId(request)
	if operationIdError != nil {
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidParams(operationIdError.Error()), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	//check for json-rpc 1.1 header and body operationid and service id mismatch
	vapiServiceId := r.Header.Get("vapi-service")
	vapiOperationId := r.Header.Get("vapi-operation")
	if vapiServiceId != "" && vapiOperationId == "" {
		log.Error("operation identifier missing in HTTP header")
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcMismatchOperationIdError("operation identifier missing in HTTP header"), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}
	if vapiOperationId != "" && vapiServiceId == "" {
		log.Error("service identifier missing in HTTP header")
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcMismatchOperationIdError("service identifier missing in HTTP header"), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	if (vapiServiceId != "" && vapiServiceId != serviceId) || (vapiOperationId != "" && vapiOperationId != operationId) {
		// Throw JSON-RPC 2.0 error with code -31001 and message "Mismatching operation identifier in HTTP header and payload"
		// when headers do not match the body
		log.Error("Mismatching operation identifier in HTTP header and payload")
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcMismatchOperationIdError(nil), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	var executionContext, execContextError = j.executionContext(request)
	if execContextError != nil {
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest(execContextError.Error()), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	protocol := tracing.JsonRpc
	if vapiServiceId != "" {
		protocol = tracing.JsonRpc11
	}
	traceCtx, serverSpan := j.tracer(serviceId, operationId, protocol, r)
	defer serverSpan.Finish()

	cancelCtx, cancelFunc := context.WithCancel(traceCtx)
	defer cancelFunc()

	executionContext.WithContext(cancelCtx)

	// TODO add header parser as comma-separated headers are not distinct values in https://pkg.go.dev/net/http#Header.Values
	// TODO headers should be case insensitive - utilize https://pkg.go.dev/net/http#Header.Values (requires go v1.14)
	hasStreamAcceptType := common.StringSliceContains(r.Header[lib.HTTP_ACCEPT], lib.VAPI_STREAMING_CLEAN_JSON_CONTENT_TYPE) ||
		common.StringSliceContains(r.Header[lib.HTTP_ACCEPT], lib.VAPI_STREAMING_CONTENT_TYPE)
	hasMonoAcceptType := common.StringSliceContains(r.Header[lib.HTTP_ACCEPT], lib.JSON_CONTENT_TYPE)

	if hasMonoAcceptType && hasStreamAcceptType {
		executionContext.WithContext(context.WithValue(executionContext.Context(), core.ResponseTypeKey, core.MonoAndStreamRespose))
	} else if hasMonoAcceptType {
		executionContext.WithContext(context.WithValue(executionContext.Context(), core.ResponseTypeKey, core.OnlyMonoResponse))
	} else if hasStreamAcceptType {
		executionContext.WithContext(context.WithValue(executionContext.Context(), core.ResponseTypeKey, core.OnlyStreamResponse))
	}

	server.CopyHeadersToContexts(executionContext, r)
	if !executionContext.ApplicationContext().HasProperty(lib.OPID) {
		log.Debug("opId was not present for the request")
	} else {
		opId := executionContext.ApplicationContext().GetProperty(lib.OPID)
		log.Debugf("Processing operation with opId %s", *opId)
	}

	var inputValue, inputError = j.input(request)
	if inputError != nil {
		var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest(inputError.Error()), &request)
		j.sendResponse(jsonRpcRequestError, rw, nil)
		return
	}

	_, logWireValue := os.LookupEnv(ENABLE_VAPI_PROVIDER_WIRE_LOGGING)
	if logWireValue {
		jsonRpcEncoder := NewJsonRpcEncoder()
		jsonRpcEncoder.SetRedactSecretFields(true)
		encodedInput, encodeError := jsonRpcEncoder.Encode(inputValue)
		if encodeError != nil {
			var jsonRpcRequestError = NewJsonRpcRequestError(NewJsonRpcErrorInvalidRequest(encodeError.Error()), &request)
			j.sendResponse(jsonRpcRequestError, rw, nil)
			return
		}
		log.Debugf("Handling new request with input %+v", string(encodedInput))
	}

	var methodResult = j.apiProvider.Invoke(serviceId, operationId, inputValue, executionContext)
	if methodResult.IsResponseStream() {
		rw.Header().Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.VAPI_STREAMING_CONTENT_TYPE)
		flusher, _ := rw.(http.Flusher)

		for i := range methodResult.ResponseStream() {
			var jsonRpc20Response = j.prepareResponseBody(request, i)
			var err = j.SendResponseFrame(jsonRpc20Response, rw, i.Error(), flusher)

			// Sending the response failed. Notify the everyone listening on the cancel
			// context done channel. Then terminate the stream.
			if err != nil {
				log.Errorf("Sending stream response failed with: %v", err)
				serverSpan.LogError(err)
				return
			}

			// If method implementation returned an error, after sending it back in the
			// response, break the stream loop. The termination frame is sent after that.
			if i.Error() != nil {
				serverSpan.LogVapiError(i.Error())
				break
			}
		}

		// Put the stream termination frame at the end of the response
		// to signal the client that the stream is finished.
		j.sendTerminalFrame(&request, rw, flusher)
	} else {
		if !methodResult.IsSuccess() {
			serverSpan.LogVapiError(methodResult.Error())
		}
		var jsonRpc20Response = j.prepareResponseBody(request, methodResult)
		j.sendResponse(jsonRpc20Response, rw, methodResult.Error())
	}
}

func (j *JsonRpcHandler) SendResponseFrame(response interface{}, rw http.ResponseWriter, errorVal *data.ErrorValue, flusher http.Flusher) error {
	var frame, encodeError = j.jsonRpcEncoder.Encode(response)
	if encodeError != nil {
		return encodeError
	}
	if errorVal != nil {
		//Accessing directly to avoid canonicalizing header key.
		rw.Header()[lib.VAPI_ERROR] = []string{errorVal.Name()}
	}
	var err = j.writeFrame(frame, rw)
	if err != nil {
		return err
	}
	flusher.Flush()
	return nil
}

func (j *JsonRpcHandler) sendTerminalFrame(request *JsonRpc20Request, rw http.ResponseWriter, flusher http.Flusher) {
	log.Debug("Sending terminal frame")
	terminalFrame := make(map[string]interface{})
	terminalFrame[lib.JSONRPC] = request.version
	terminalFrame[lib.JSONRPC_ID] = request.id
	terminalFrame[lib.METHOD_RESULT] = map[string]interface{}{}
	frameBytes, _ := json.Marshal(terminalFrame)

	var err = j.writeFrame(frameBytes, rw)
	if err != nil {
		log.Error(err)
	}
	flusher.Flush()
}

func (j *JsonRpcHandler) writeFrame(value []byte, rw http.ResponseWriter) error {
	dataLengthInHex := fmt.Sprintf("%x", len(value))
	_, err := rw.Write([]byte(dataLengthInHex))
	if err != nil {
		return err
	}
	_, err = rw.Write(lib.CRLFBytes)
	if err != nil {
		return err
	}
	_, err = rw.Write(value)
	if err != nil {
		return err
	}
	_, err = rw.Write(lib.CRLFBytes)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonRpcHandler) operationId(request JsonRpc20Request) (string, error) {
	var params = request.Params()
	if val, ok := params[lib.REQUEST_OPERATION_ID]; ok {
		return val.(string), nil
	} else {
		return "", errors.New("params does not have key: operationId")
	}
}

func (j *JsonRpcHandler) serviceId(request JsonRpc20Request) (string, error) {
	var params = request.Params()
	if val, ok := params[lib.REQUEST_SERVICE_ID]; ok {
		return val.(string), nil
	} else {
		return "", errors.New("params does not have key: serviceId")
	}
}

func (j *JsonRpcHandler) input(request JsonRpc20Request) (data.DataValue, error) {
	var params = request.Params()
	if val, ok := params[lib.REQUEST_INPUT]; ok {
		var structValue, dvError = j.jsonRpcDecoder.GetDataValue(val)
		if dvError != nil {
			return nil, dvError
		}
		return structValue, nil
	} else {
		return nil, errors.New("params does not have key: input")
	}
}

func (j *JsonRpcHandler) executionContext(request JsonRpc20Request) (*core.ExecutionContext, error) {
	var params = request.Params()
	if val, ok := params[lib.EXECUTION_CONTEXT]; ok {
		var executionContext, dsError = j.jsonRpcDecoder.DeSerializeExecutionContext(val)
		if dsError != nil {
			return nil, errors.New("error de-serializing execution context")
		}
		return executionContext, nil
	} else {
		return nil, errors.New("params does not have key: ctx")
	}
}

func (j *JsonRpcHandler) prepareResponseBody(request JsonRpc20Request, methodResult core.MonoResult) JsonRpc20Response {
	_, logWireValue := os.LookupEnv(ENABLE_VAPI_PROVIDER_WIRE_LOGGING)
	if logWireValue {
		jsonRpcEncoder := NewJsonRpcEncoder()
		jsonRpcEncoder.SetRedactSecretFields(true)
		encodedMethodResult, _ := jsonRpcEncoder.Encode(methodResult)
		log.Debugf("Sending response with output %+v", string(encodedMethodResult))

	}
	var version = request.version
	var id = request.id
	var jsonRpc20Response = NewJsonRpc20Response(version, id, methodResult, nil)
	return jsonRpc20Response
}
