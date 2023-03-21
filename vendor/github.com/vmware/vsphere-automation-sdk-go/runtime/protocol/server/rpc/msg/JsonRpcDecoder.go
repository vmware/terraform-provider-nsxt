/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"reflect"
	"strings"
)

//JsonRpcDecoder decodes unmarshalled json object into vapi data value format
type JsonRpcDecoder struct {
	serializers.MethodResultDeserializerBase
}

func NewJsonRpcDecoder() *JsonRpcDecoder {
	i := &JsonRpcDecoder{}
	i.Impl = i
	return i
}

func (decoder *JsonRpcDecoder) Decode(data []byte, v interface{}) error {
	return nil
}

func (decoder *JsonRpcDecoder) DecodeDataValue(data string) (data.DataValue, error) {
	result, err := deserializeJsonString(data)
	if err != nil {
		return nil, err
	}
	return decoder.GetDataValue(result)
}

func (decoder *JsonRpcDecoder) GetDataValue(jsonDataValue interface{}) (data.DataValue, error) {

	switch result := jsonDataValue.(type) {
	case map[string]interface{}:
		inputData := jsonDataValue.(map[string]interface{})
		for k, val := range inputData {
			switch k {
			case "STRUCTURE":
				return decoder.DeSerializeStructValue(val.(map[string]interface{}))
			case "OPTIONAL":
				return decoder.DeSerializeOptionalValue(val)
			case "ERROR":
				return decoder.DeSerializeErrorValue(val.(map[string]interface{}))
			case "BINARY":
				return decoder.DeSerializeBinaryValue(val)
			case "SECRET":
				return decoder.DeSerializeSecretValue(val)
			}
		}

	case []interface{}:
		return decoder.DeSerializeListValue(jsonDataValue.([]interface{}))
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
func (decoder *JsonRpcDecoder) DeSerializeSecretValue(jsonSecretValue interface{}) (*data.SecretValue, error) {
	return data.NewSecretValue(jsonSecretValue.(string)), nil
}
func (decoder *JsonRpcDecoder) DeSerializeListValue(jsonListValue []interface{}) (*data.ListValue, error) {
	var listValue = data.NewListValue()
	for _, element := range jsonListValue {
		listElementDataValue, dvError := decoder.GetDataValue(element)
		if dvError != nil {
			return nil, dvError
		}
		listValue.Add(listElementDataValue)
	}
	return listValue, nil
}
func (decoder *JsonRpcDecoder) DeSerializeStructValue(jsonStructValue map[string]interface{}) (*data.StructValue, error) {
	var structName string
	var fields map[string]interface{}
	for key, val := range jsonStructValue {
		structName = key
		fields = val.(map[string]interface{})
	}

	var structVal = data.NewStructValue(structName, nil)
	for fieldName, fieldJsonValue := range fields {
		var fieldDataValue, dvError = decoder.GetDataValue(fieldJsonValue)
		if dvError != nil {
			return nil, dvError
		}
		structVal.SetField(fieldName, fieldDataValue)
	}
	return structVal, nil

}

func (decoder *JsonRpcDecoder) DeSerializeErrorValue(jsonStructValue map[string]interface{}) (*data.ErrorValue, error) {
	var errorName string
	var fields map[string]interface{}
	for key, val := range jsonStructValue {
		errorName = key
		fields = val.(map[string]interface{})
	}
	var structVal = data.NewErrorValue(errorName, nil)
	for fieldName, fieldJsonValue := range fields {
		fieldDataValue, dvError := decoder.GetDataValue(fieldJsonValue)
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

func (decoder *JsonRpcDecoder) DeSerializeBinaryValue(val interface{}) (data.DataValue, error) {
	//binary value is base64 encoded on the wire. decode it first.
	base64DecodedString, err := base64.StdEncoding.DecodeString(val.(string))
	if err != nil {
		return nil, err
	}
	byt := base64DecodedString
	return data.NewBlobValue(byt), nil
}

func DeSerializeJson(request interface{}) (map[string]interface{}, error) {
	if requestString, ok := request.(string); ok {
		return deserializeJsonString(requestString)
	} else if requestBytes, ok := request.([]byte); ok {
		requestString := string(requestBytes)
		return deserializeJsonString(requestString)
	}
	return nil, errors.New("error deserializing json")
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
