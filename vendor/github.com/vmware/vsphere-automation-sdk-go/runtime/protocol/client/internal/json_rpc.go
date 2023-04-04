/* Copyright © 2019, 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/common"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/cleanjson"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	vapiHttp "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/http"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/server/rpc/msg"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// JsonRpcHttpProtocol provides functionality for making and handling http JSON-RPC calls to a server
type JsonRpcHttpProtocol struct {
	Connector
	*HttpTransport
	requestPreProcessors  []core.JSONRPCRequestPreProcessor
	responseHandlers      []vapiHttp.ClientResponseHandler
	streamingAcceptHeader string
}

var _ StreamingTransport = NewJsonRpcHttpProtocol(nil)

// NewJsonRpcHttpProtocol instantiates instance of JsonRpcHttpProtocol
func NewJsonRpcHttpProtocol(c Connector) *JsonRpcHttpProtocol {
	return &JsonRpcHttpProtocol{
		Connector:     c,
		HttpTransport: NewHttpTransport(),
		requestPreProcessors: []core.JSONRPCRequestPreProcessor{
			security.NewJSONSsoSigner(),
		},
		responseHandlers: []vapiHttp.ClientResponseHandler{
			NewFramesResponseHandler(),
			NewCleanJsonFramesResponseHandler(),
			NewRegularResponseHandler(),
		},
		streamingAcceptHeader: lib.VAPI_STREAMING_CONTENT_TYPE,
	}
}

// SetStreamingProtocol sets streaming wire protocol
// Default option is lib.VAPI_STREAMING_CONTENT_TYPE
// lib.VAPI_STREAMING_CLEAN_JSON_CONTENT_TYPE can be used as well for shorter wire format.
func (j *JsonRpcHttpProtocol) SetStreamingProtocol(acceptHeader string) {
	j.streamingAcceptHeader = acceptHeader
}

func (j *JsonRpcHttpProtocol) SetClientFrameDeserializer(deserializer vapiHttp.ClientFrameDeserializer) {
	for _, handler := range j.responseHandlers {
		if h, ok := handler.(vapiHttp.ClientFramesResponseHandler); ok {
			h.SetClientFrameDeserializer(deserializer)
		}
	}
}

// Invoke from given parameters makes http request to remote server
func (j *JsonRpcHttpProtocol) Invoke(
	serviceId string,
	operationId string,
	inputValue data.DataValue,
	ctx *core.ExecutionContext) core.MethodResult {

	request, err := j.buildRequest(serviceId, operationId, inputValue, ctx)
	if err != nil {
		vapiErr := getRequestError(err.Error())
		return getErrorMethodResult(vapiErr)
	}

	response, err := j.executeRequest(ctx, request)
	if err != nil {
		vapiErr := getVAPIError(err)
		return getErrorMethodResult(vapiErr)
	}

	return j.handleResponse(ctx, response)
}

func (j *JsonRpcHttpProtocol) buildRequest(
	serviceId string,
	operationId string,
	inputValue data.DataValue,
	ctx *core.ExecutionContext) (*http.Request, error) {

	if ctx == nil {
		panic("execution context can't be nil")
	}
	if !ctx.ApplicationContext().HasProperty(lib.OPID) {
		common.InsertOperationId(ctx.ApplicationContext())
	}
	opId := ctx.ApplicationContext().GetProperty(lib.OPID)

	var jsonRpcEncoder = msg.NewJsonRpcEncoder()
	encodedInput, err := jsonRpcEncoder.Encode(inputValue)
	if err != nil {
		return nil, err
	}

	encodedInputMap := map[string]interface{}{}
	err = json.Unmarshal(encodedInput, &encodedInputMap)
	if err != nil {
		return nil, err
	}

	log.Debugf("Invoking with input %+v", string(encodedInput))

	executionCtxJson, err := ctx.JSON()
	if err != nil {
		return nil, err
	}

	var params = make(map[string]interface{})
	params[lib.REQUEST_INPUT] = encodedInputMap
	params[lib.EXECUTION_CONTEXT] = executionCtxJson
	params[lib.REQUEST_OPERATION_ID] = operationId
	params[lib.REQUEST_SERVICE_ID] = serviceId
	var jsonRpcRequest = msg.NewJsonRpc20Request(lib.JSONRPC_VERSION, lib.JSONRPC_INVOKE, params, opId, false)

	var requestJson = jsonRpcRequest.JSON()
	for _, preProcessor := range j.requestPreProcessors {
		err = preProcessor.Process(&requestJson)
		if err != nil {
			return nil, err
		}
	}

	var jsonRpcRequestSerializer = msg.NewJsonRpc20RequestSerializer(jsonRpcRequest)
	requestBytes, err := jsonRpcRequestSerializer.MarshalJSON()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, j.Connector.Address(), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set(lib.HTTP_CONTENT_TYPE_HEADER, lib.JSON_CONTENT_TYPE)
	request.Header.Set(lib.VAPI_SERVICE_HEADER, serviceId)
	request.Header.Set(lib.VAPI_OPERATION_HEADER, operationId)
	request.Header.Set(lib.HTTP_USER_AGENT_HEADER, GetRuntimeUserAgentHeader())

	request.Header.Set(lib.HTTP_ACCEPT, lib.JSON_CONTENT_TYPE)

	responseType := core.AcceptableResponseType(ctx.Context())

	if responseType.AcceptsStreamResponse() {
		request.Header.Add(lib.HTTP_ACCEPT, j.streamingAcceptHeader)
	}

	copyContextsToHeaders(ctx, request.Header)

	return request, nil
}

func (j *JsonRpcHttpProtocol) handleResponse(
	ctx *core.ExecutionContext,
	response *http.Response) core.MethodResult {
	contentType := response.Header[lib.HTTP_CONTENT_TYPE_HEADER]
	if contentType == nil {
		responseError := newResponseError(
			bindings.INTERNAL_SERVER_ERROR_DEF,
			"vapi.protocol.client.response.error.missingContentType",
			nil)
		return getErrorMethodResult(responseError)
	}

	for _, handler := range j.responseHandlers {
		result, err := handler.HandleResponse(ctx.Context(), response)
		if err != nil {
			responseError := getResponseError(err.Error())
			return getErrorMethodResult(responseError)
		}
		if result != nil {
			return result
		}
	}

	responseError := newResponseError(
		bindings.INTERNAL_SERVER_ERROR_DEF,
		"vapi.protocol.client.response.error.unknownContentType",
		map[string]string{"contentType": strings.Join(contentType, ", ")})
	return getErrorMethodResult(responseError)
}

type RegularResponseHandler struct {
	decoder *msg.JsonRpcDecoder
}

var _ vapiHttp.ClientResponseHandler = &RegularResponseHandler{}

func NewRegularResponseHandler() *RegularResponseHandler {
	return &RegularResponseHandler{decoder: msg.NewJsonRpcDecoder()}
}

func (r *RegularResponseHandler) HandleResponse(_ context.Context, response *http.Response) (core.MethodResult, error) {
	if !common.StringSliceContains(response.Header[lib.HTTP_CONTENT_TYPE_HEADER], lib.JSON_CONTENT_TYPE) {
		return nil, nil
	}

	defer closeResponse(response)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jsonRpcResponse, deserializeError := msg.DeSerializeResponse(responseBody)
	//TODO: simplify DeserializeResponse API to return only error instead of JsonRpc20Error
	if deserializeError != nil {
		responseError := getResponseError(deserializeError.Message())
		result := getErrorMethodResult(responseError)
		return result, nil
	}

	if jsonRpcResponseMap, ok := jsonRpcResponse.Result().(map[string]interface{}); ok {
		methodResult, err := r.decoder.DeserializeMethodResult(jsonRpcResponseMap)
		if err != nil {
			return nil, err
		}
		return methodResult, nil
	} else if jsonRpcErrorResponseMap := jsonRpcResponse.Error(); jsonRpcErrorResponseMap != nil {
		methodResult, err := r.decoder.DeserializeMethodResult(jsonRpcErrorResponseMap)
		if err != nil {
			return nil, err
		}
		return methodResult, nil
	}

	err = newResponseError(
		bindings.INTERNAL_SERVER_ERROR_DEF,
		"vapi.protocol.client.response.error.invalidResponse",
		nil)
	return nil, err
}

type FramesResponseHandler struct {
	deserializer vapiHttp.ClientFrameDeserializer
}

var _ vapiHttp.ClientFramesResponseHandler = &FramesResponseHandler{}

func NewFramesResponseHandler() *FramesResponseHandler {
	return &FramesResponseHandler{
		deserializer: NewJsonFrameDeserializer(msg.NewJsonRpcDecoder()),
	}
}

func (c *FramesResponseHandler) SetClientFrameDeserializer(deserializer vapiHttp.ClientFrameDeserializer) {
	c.deserializer = deserializer
}

func (c *FramesResponseHandler) HandleResponse(
	ctx context.Context,
	response *http.Response) (core.MethodResult, error) {

	contentType := response.Header[lib.HTTP_CONTENT_TYPE_HEADER]
	if !common.StringSliceContains(contentType, lib.VAPI_STREAMING_CONTENT_TYPE) {
		return nil, nil
	}

	return handleFramesResponse(ctx, c.deserializer, response)
}

type CleanJsonFramesResponseHandler struct {
	deserializer vapiHttp.ClientFrameDeserializer
}

var _ vapiHttp.ClientFramesResponseHandler = &CleanJsonFramesResponseHandler{}

func NewCleanJsonFramesResponseHandler() *CleanJsonFramesResponseHandler {
	return &CleanJsonFramesResponseHandler{
		deserializer: NewJsonFrameDeserializer(cleanjson.NewJsonToDataValueDecoder()),
	}
}

func (c *CleanJsonFramesResponseHandler) SetClientFrameDeserializer(deserializer vapiHttp.ClientFrameDeserializer) {
	c.deserializer = deserializer
}

func (c *CleanJsonFramesResponseHandler) HandleResponse(
	ctx context.Context,
	response *http.Response) (core.MethodResult, error) {

	contentType := response.Header[lib.HTTP_CONTENT_TYPE_HEADER]
	if !common.StringSliceContains(contentType, lib.VAPI_STREAMING_CLEAN_JSON_CONTENT_TYPE) {
		return nil, nil
	}

	return handleFramesResponse(ctx, c.deserializer, response)
}

func handleFramesResponse(
	ctx context.Context,
	deserializer vapiHttp.ClientFrameDeserializer,
	response *http.Response) (core.MethodResult, error) {

	frames := make(chan []byte)

	framesReadCtx := core.WithErrorContext(ctx)
	go getFrames(framesReadCtx, frames, response)

	methodResultFrames, err := deserializer.DeserializeFrames(framesReadCtx, frames)
	if err != nil {
		return nil, err
	}

	methodResult := core.NewStreamMethodResult(methodResultFrames, func() { framesReadCtx.Cancel(nil) })
	return methodResult, nil

}

func getFrames(ctx *core.ErrorContext, frames chan []byte, response *http.Response) {
	defer close(frames)
	defer closeResponse(response)

	vapiFrameReader := vapiHttp.NewVapiFrameReader(response.Body)

	for {
		// check if client canceled request before every frame read
		select {
		case <-ctx.Done():
			// client cancel
			log.Info("Client canceled request")
			return
		default:
			//continue execution
		}

		frameData, err := vapiFrameReader.ReadFrame()
		if err == io.EOF || err == context.Canceled {
			// end of stream or request canceled by client
			return
		} else if err != nil {
			ctx.Cancel(getResponseError(err.Error()))
			return
		}
		frames <- frameData
	}
}

type JsonFrameDeserializer struct {
	decoder serializers.MethodResultDeserializer
}

var _ vapiHttp.ClientFrameDeserializer = &JsonFrameDeserializer{}

func NewJsonFrameDeserializer(decoder serializers.MethodResultDeserializer) *JsonFrameDeserializer {
	return &JsonFrameDeserializer{decoder: decoder}
}

func (j *JsonFrameDeserializer) DeserializeFrames(ctx context.Context, frames chan []byte) (chan core.MonoResult, error) {
	result := make(chan core.MonoResult)

	go func() {
		defer close(result)

		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				// context canceled
				if err == context.Canceled {
					return
				}
				// another error occurred
				if err != nil {
					// error populated in getFrames
					responseError := getResponseError(err.Error())
					result <- getErrorMonoResult(responseError)
				}
				return
			case frame := <-frames:
				deserializedFrame, err := j.deserializeFrame(frame)
				if err != nil {
					responseError := getResponseError(err.Error())
					result <- getErrorMonoResult(responseError)
				}

				if deserializedFrame != nil {
					result <- deserializedFrame
				} else {
					return
				}
			}
		}
	}()

	return result, nil
}

func (j *JsonFrameDeserializer) deserializeFrame(frameData []byte) (core.MethodResult, error) {
	var jsonFrame map[string]interface{}
	dec := json.NewDecoder(bytes.NewBuffer(frameData))
	dec.UseNumber()
	err := dec.Decode(&jsonFrame)
	if err != nil {
		return nil, err
	}

	jsonRpcResponse, responseDeserializationErr := msg.GetJsonRpc20Response(jsonFrame)
	//TODO: simplify GetJsonRpc20Response API to return only error instead of JsonRpc20Error
	if responseDeserializationErr != nil {
		responseError := getResponseError(responseDeserializationErr.Message())
		return nil, responseError
	}

	return j.getMethodResultFrame(&jsonRpcResponse, j.decoder)
}

func (j *JsonFrameDeserializer) getMethodResultFrame(
	jsonRpcResponse *msg.JsonRpc20Response,
	decoder serializers.MethodResultDeserializer) (core.MethodResult, error) {
	if jsonRpcResponseMap, ok := jsonRpcResponse.Result().(map[string]interface{}); ok {
		// check for terminal frame
		if len(jsonRpcResponseMap) == 0 {
			log.Debug("Received terminal frame")
			return nil, nil
		}
		methodResult, err := decoder.DeserializeMethodResult(jsonRpcResponseMap)
		if err != nil {
			return nil, err
		}
		return methodResult, nil
	} else if jsonRpcErrorResponseMap := jsonRpcResponse.Error(); jsonRpcErrorResponseMap != nil {
		methodResult, err := decoder.DeserializeMethodResult(jsonRpcErrorResponseMap)
		if err != nil {
			return nil, err
		}
		return methodResult, nil
	}

	responseError := newResponseError(
		bindings.INTERNAL_SERVER_ERROR_DEF,
		"vapi.protocol.client.response.error.invalidFrame",
		nil)
	return nil, responseError
}

func closeResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		_ = r.Body.Close()
	}
}
