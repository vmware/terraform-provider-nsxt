/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"reflect"
	"strings"
)

//Decodes json into vapi structs
type JsonRpcDecoder struct{}

func NewJsonRpcDecoder() *JsonRpcDecoder {
	return &JsonRpcDecoder{}
}

func (j *JsonRpcDecoder) Decode(data []byte, v interface{}) error {
	return nil
}

func (j *JsonRpcDecoder) DecodeDataValue(data string) (data.DataValue, error) {
	result, err := deserializeJsonString(data)
	if err != nil {
		return nil, err
	}
	return j.GetDataValue(result)
}

func (j *JsonRpcDecoder) GetDataValue(jsonDataValue interface{}) (data.DataValue, error) {

	switch result := jsonDataValue.(type) {
	case map[string]interface{}:
		inputData := jsonDataValue.(map[string]interface{})
		for k, val := range inputData {
			switch k {
			case "STRUCTURE":
				return j.DeSerializeStructValue(val.(map[string]interface{}))
			case "OPTIONAL":
				return j.DeSerializeOptionalValue(val)
			case "ERROR":
				return j.DeSerializeErrorValue(val.(map[string]interface{}))
			case "BINARY":
				return j.DeSerializeBinaryValue(val)
			case "SECRET":
				return j.DeSerializeSecretValue(val)
			}
		}

	case []interface{}:
		return j.DeSerializeListValue(jsonDataValue.([]interface{}))
	case string:
		return data.NewStringValue(jsonDataValue.(string)), nil
	case bool:
		return data.NewBooleanValue(jsonDataValue.(bool)), nil
	case json.Number:
		// By default, json marshaller converts json number to float64.
		// Use the below strategy to distinguish between floating point numbers and int64
		intValue, err := result.Int64()
		if err == nil {
			return data.NewIntegerValue(intValue), nil
		}
		floatValue, err := result.Float64()
		if err == nil {
			return data.NewDoubleValue(floatValue), nil
		}
	case nil:
		return nil, nil
	case interface{}:
		log.Debugf("Got type of json object as interface{}. Determine its equivalent DataValue")
	default:
		log.Debugf("Could not determine appropriate DataValue for %s", result)
	}

	return nil, nil
}
func (j *JsonRpcDecoder) DeSerializeSecretValue(jsonSecretValue interface{}) (*data.SecretValue, error) {
	return data.NewSecretValue(jsonSecretValue.(string)), nil
}
func (j *JsonRpcDecoder) DeSerializeListValue(jsonListValue []interface{}) (*data.ListValue, error) {
	var listValue = data.NewListValue()
	for _, element := range jsonListValue {
		listElementDataValue, dvError := j.GetDataValue(element)
		if dvError != nil {
			return nil, dvError
		}
		listValue.Add(listElementDataValue)
	}
	return listValue, nil
}
func (j *JsonRpcDecoder) DeSerializeStructValue(jsonStructValue map[string]interface{}) (*data.StructValue, error) {
	var structName string
	var fields map[string]interface{}
	for key, val := range jsonStructValue {
		structName = key
		fields = val.(map[string]interface{})
	}

	var structVal = data.NewStructValue(structName, nil)
	for fieldName, fieldJsonValue := range fields {
		var fieldDataValue, dvError = j.GetDataValue(fieldJsonValue)
		if dvError != nil {
			return nil, dvError
		}
		structVal.SetField(fieldName, fieldDataValue)
	}
	return structVal, nil

}

func (j *JsonRpcDecoder) DeSerializeErrorValue(jsonStructValue map[string]interface{}) (*data.ErrorValue, error) {
	var errorName string
	var fields map[string]interface{}
	for key, val := range jsonStructValue {
		errorName = key
		fields = val.(map[string]interface{})
	}
	var structVal = data.NewErrorValue(errorName, nil)
	for fieldName, fieldJsonValue := range fields {
		fieldDataValue, dvError := j.GetDataValue(fieldJsonValue)
		if dvError != nil {
			return nil, dvError
		}
		structVal.SetField(fieldName, fieldDataValue)
	}
	return structVal, nil

}
func (decoder *JsonRpcDecoder) DeSerializeOptionalValue(i interface{}) (*data.OptionalValue, error) {
	dataValue, dvError := decoder.GetDataValue(i)
	if dvError != nil {
		return nil, dvError
	}
	return data.NewOptionalValue(dataValue), nil
}

/**
Deserialize Methodresult
*/
func (decoder *JsonRpcDecoder) DeSerializeMethodResult(methodResultInput map[string]interface{}) (core.MethodResult, error) {
	if val, ok := methodResultInput[lib.METHOD_RESULT_OUTPUT]; ok {
		var output, err = decoder.GetDataValue(val)
		if err != nil {
			return core.MethodResult{}, err
		}
		return core.NewMethodResult(output, nil), nil
	} else if val, ok := methodResultInput[lib.METHOD_RESULT_ERROR]; ok {
		var methodResultError, err = decoder.GetDataValue(val)
		if err != nil {
			return core.MethodResult{}, err
		}
		return core.NewMethodResult(nil, methodResultError.(*data.ErrorValue)), nil
	}

	return core.MethodResult{}, errors.New("error de-serializing methodresult")
}

func (decoder *JsonRpcDecoder) DeSerializeApplicationContext(appCtxData interface{}) (*core.ApplicationContext, error) {
	if reflect.TypeOf(appCtxData) == nil {
		return core.NewApplicationContext(nil), nil
	}
	var appCtxMap = make(map[string]*string)
	var appCtxMapData = appCtxData.(map[string]interface{})
	for key, val := range appCtxMapData {
		if valStr, ok := val.(string); ok {
			appCtxMap[key] = &valStr
		} else {
			var stringPtr *string = nil
			appCtxMap[key] = stringPtr
		}
	}
	return core.NewApplicationContext(appCtxMap), nil
}

func (decoder *JsonRpcDecoder) DeSerializeSecurityContext(secCtxData map[string]interface{}) (core.SecurityContext, error) {
	if reflect.TypeOf(secCtxData) == nil {
		return core.NewSecurityContextImpl(), nil
	} else {
		secCtx := core.NewSecurityContextImpl()
		secCtx.SetContextMap(secCtxData)
		return secCtx, nil
	}
}

/**
Deserialize ExecutionContext
*/
func (decoder *JsonRpcDecoder) DeSerializeExecutionContext(executionContext interface{}) (*core.ExecutionContext, error) {
	var executionContextMap map[string]interface{}
	if executionContextMapVal, isMap := executionContext.(map[string]interface{}); !isMap {
		return core.NewExecutionContext(nil, nil), nil
	} else {
		executionContextMap = executionContextMapVal
	}
	var appCtx *core.ApplicationContext = nil
	var err error
	if val, ok := executionContextMap[lib.APPLICATION_CONTEXT]; ok {
		appCtx, err = decoder.DeSerializeApplicationContext(val)
		if err != nil {
			return nil, err
		}
	}
	var secCtx core.SecurityContext
	if val, ok := executionContextMap[lib.SECURITY_CONTEXT]; ok {
		secCtx, err = decoder.DeSerializeSecurityContext(val.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
	} else {
		log.Debugf("SecurityContext not passed in the request. Creating an empty security context")
		secCtx = core.NewSecurityContextImpl()
	}
	return core.NewExecutionContext(appCtx, secCtx), nil

}

// Get id value from request.
// id can be string, int or nil.
// it is important find out the type of id and return it.
// when id is a number, its type is json.Number but its default representation is string.
// In this case, type has to be represented as int otherwise request and response id will be of different type.
func (decoder *JsonRpcDecoder) getIdValue(id interface{}) (interface{}, error) {
	var idString, isString = id.(string)
	if !isString {
		//check if id is an int.
		var jsonNumber, numError = id.(json.Number)
		if !numError {
			//id is not string or json.Number. check if its nil
			if reflect.TypeOf(id) != nil {
				log.Errorf("JSON RPC request id must be string or int or nil")
				return nil, errors.New("JSON RPC request id must be string or int or nil")
			} else {
				return nil, nil
			}
		} else {
			//check if number is int64
			var idInt, intErr = jsonNumber.Int64()
			if intErr != nil {
				log.Errorf("JSON RPC request id must be string or int or nil")
				return nil, errors.New("JSON RPC request id must be string or int or nil")
			} else {
				return idInt, nil
			}
		}

	} else {
		return idString, nil
	}
}

func (decoder *JsonRpcDecoder) getJsonRpc20Request(request map[string]interface{}) (JsonRpc20Request, *JsonRpc20Error) {
	var version = request[lib.JSONRPC].(string)
	var notification = true
	var id interface{}
	var err error
	if idValue, ok := request[lib.JSONRPC_ID]; ok {
		notification = false
		id, err = decoder.getIdValue(idValue)
		if err != nil {
			return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(err.Error())
		}
	}
	var method, ok = request[lib.JSONRPC_METHOD].(string)
	if !ok {
		log.Errorf("JSON RPC request method must be string ")
		return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(nil)
	}

	if paramValue, ok := request[lib.JSONRPC_PARAMS]; ok {
		paramMap, isObject := paramValue.(map[string]interface{})
		if !isObject && reflect.TypeOf(paramValue) != nil {
			return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(nil)
		}
		return NewJsonRpc20Request(version, method, paramMap, id, notification), nil
	} else {
		return NewJsonRpc20Request(version, method, nil, id, notification), nil
	}
}

func (decoder *JsonRpcDecoder) GetJsonRpc20Response(response map[string]interface{}) (JsonRpc20Response, *JsonRpc20Error) {
	var version string
	if versionValue, ok := response[lib.JSONRPC]; ok {
		version = versionValue.(string)
	} else {
		return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("jsonrpc version not present")
	}
	var id interface{}
	if responseId, ok := response[lib.JSONRPC_ID]; ok {
		idVal, err := decoder.getIdValue(responseId)
		if err != nil {
			return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("Invalid Id")
		} else {
			id = idVal
		}

	}
	if result, ok := response[lib.METHOD_RESULT]; ok {
		return NewJsonRpc20Response(version, id, result, nil), nil
	} else {
		var err = response[lib.METHOD_RESULT_ERROR].(map[string]interface{})
		return NewJsonRpc20Response(version, id, nil, err), nil
	}

}

func deserializeJsonString(inputString string) (map[string]interface{}, error) {
	var requestObject = make(map[string]interface{})
	d := json.NewDecoder(strings.NewReader(inputString))
	d.UseNumber()
	if err := d.Decode(&requestObject); err != nil {
		return nil, err
	} else {
		return requestObject, nil
	}
}

/**
 * Deserialize incoming json which could be in byte or string format.
 */
func DeSerializeJson(request interface{}) (map[string]interface{}, error) {
	if requestString, ok := request.(string); ok {
		return deserializeJsonString(requestString)
	} else if requestBytes, ok := request.([]byte); ok {
		requestString := string(requestBytes)
		return deserializeJsonString(requestString)
	}
	return nil, errors.New("Error Deserializing json")
}

/**
 * Deserialize Json Rpc Request
 *
 * request can be string or bytes.
 */
func (decoder *JsonRpcDecoder) DeSerializeRequest(request interface{}) (JsonRpc20Request, *JsonRpc20Error) {
	requestObject, err := DeSerializeJson(request)
	if err != nil {
		return JsonRpc20Request{}, NewJsonRpcErrorParseError(err)
	}
	return decoder.getJsonRpc20Request(requestObject)

}

func (decoder *JsonRpcDecoder) DeSerializeResponse(response interface{}) (JsonRpc20Response, *JsonRpc20Error) {
	responseObj, err := DeSerializeJson(response)
	if err != nil {
		return JsonRpc20Response{}, NewJsonRpcErrorParseError(err)
	}
	return decoder.GetJsonRpc20Response(responseObj)
}
func (decoder *JsonRpcDecoder) DeSerializeBinaryValue(val interface{}) (data.DataValue, error) {
	//binary value is base64 encoded on the wire. decode it first.
	base64DecodedString, err := base64.StdEncoding.DecodeString(val.(string))
	if err != nil {
		return nil, err
	}
	byt := base64DecodedString
	return data.NewBlobValue(byt), nil
}
