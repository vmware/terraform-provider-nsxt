/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/server"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const ENABLE_VAPI_PROVIDER_WIRE_LOGGING = "ENABLE_VAPI_PROVIDER_WIRE_LOGGING"

type JsonRpcHandler struct {
	apiProvider         core.APIProvider
	vapiAppCtxConstants map[string]string
	jsonRpcEncoder      *JsonRpcEncoder
	jsonRpcDecoder      *JsonRpcDecoder
	requestProcessors   []server.RequestPreProcessor
}

func NewJsonRpcHandler(apiProvider core.APIProvider) *JsonRpcHandler {
	var jsonRpcEncoder = NewJsonRpcEncoder()
	var jsonRpcDecoder = NewJsonRpcDecoder()
	vapiAppCtxConstants := map[string]string{
		"opid":                "opId",
		"actid":               "actId",
		"$showunreleasedapis": "$showUnreleasedAPIs",
		"$useragent":          "$userAgent",
		"$donotroute":         "$doNotRoute",
		"vmwaresessionid":     "vmwareSessionId",
		"activationid":        "ActivationId"}
	return &JsonRpcHandler{apiProvider: apiProvider,
		jsonRpcEncoder:      jsonRpcEncoder,
		jsonRpcDecoder:      jsonRpcDecoder,
		vapiAppCtxConstants: vapiAppCtxConstants}
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
	var request, requestDeserializationError = j.jsonRpcDecoder.getJsonRpc20Request(requestObj)
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

func (j *JsonRpcHandler) AddRequestPreProcessor(reqProcessor server.RequestPreProcessor) {
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

	j.copyHeadersToContexts(executionContext, r)
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
	var jsonRpc20Response = j.prepareResponseBody(request, methodResult, serviceId, operationId)
	j.sendResponse(jsonRpc20Response, rw, methodResult.Error())
}

func (j *JsonRpcHandler) copyHeadersToContexts(ctx *core.ExecutionContext, r *http.Request) {
	appCtx := ctx.ApplicationContext()
	secCtx := ctx.SecurityContext()
	for key, value := range r.Header {
		keyLowerCase := strings.ToLower(key)
		s := strings.Split(keyLowerCase, lib.VAPI_HEADER_PREFIX)
		if len(s) > 1 {
			// Override values in appCtx with headers with "vapi-ctx-" prefix
			// The values from HTTP headers override the body values.
			// if there are multiple values for the same header, the first entry will be chosen.
			if j.vapiAppCtxConstants[s[1]] != "" {
				appCtx.SetProperty(j.vapiAppCtxConstants[s[1]], &value[0])
			} else {
				appCtx.SetProperty(s[1], &value[0])
			}
		} else {
			switch keyLowerCase {
			case lib.HTTP_ACCEPT_LANGUAGE:
				appCtx.SetProperty(lib.HTTP_ACCEPT_LANGUAGE, &value[0])
			case lib.VAPI_SESSION_HEADER:
				// Override session value in security context with authorization header
				secCtx.SetProperty(security.SESSION_ID, value[0])
				secCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.SESSION_SCHEME_ID)
			}
		}
	}
	// When the request has $useragent header, it will override the custom one if present.
	if userAgentVal, ok := r.Header["User-Agent"]; ok {
		appCtx.SetProperty("$userAgent", &userAgentVal[0])
	}
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

func (j *JsonRpcHandler) prepareResponseBody(request JsonRpc20Request, methodResult core.MethodResult, serviceId string, operationId string) JsonRpc20Response {
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
