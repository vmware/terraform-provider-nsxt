/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"reflect"
	"strconv"
	"strings"
)

var structValuePtr = reflect.TypeOf((*data.StructValue)(nil))
var stringValuePtr = reflect.TypeOf((*data.StringValue)(nil))
var integerValuePtr = reflect.TypeOf((*data.IntegerValue)(nil))
var doubleValuePtr = reflect.TypeOf((*data.DoubleValue)(nil))
var listValuePtr = reflect.TypeOf((*data.ListValue)(nil))
var optionalValuePtr = reflect.TypeOf((*data.OptionalValue)(nil))
var errorValuePtr = reflect.TypeOf((*data.ErrorValue)(nil))
var voidValuePtr = reflect.TypeOf((*data.VoidValue)(nil))
var boolValuePtr = reflect.TypeOf((*data.BooleanValue)(nil))
var blobValuePtr = reflect.TypeOf(((*data.BlobValue)(nil)))
var secretValuePtr = reflect.TypeOf((*data.SecretValue)(nil))
var methodResultType = reflect.TypeOf(core.MethodResult{})
var executionContextPtr = reflect.TypeOf((*core.ExecutionContext)(nil))
var jsonRpcRequestType = reflect.TypeOf(JsonRpc20Request{})
var jsonRpcResponseType = reflect.TypeOf(JsonRpc20Response{})
var jsonRpcRequestErrorType = reflect.TypeOf((*JsonRpcRequestError)(nil))

// Encodes vapi structs to json.
type JsonRpcEncoder struct {
	// Setting this field to true will replace value of the secret fields to "redacted".
	// This should be used for logging purposes only.
	redactSecretFields bool
}

func NewJsonRpcEncoder() *JsonRpcEncoder {
	return &JsonRpcEncoder{}
}

func (j *JsonRpcEncoder) SetRedactSecretFields(redact bool) {
	j.redactSecretFields = redact
}

func (j *JsonRpcEncoder) Encode(value interface{}) ([]byte, error) {
	return dispatcher(value, j.redactSecretFields)
}

func getCustomSerializer(val interface{}, redactSecret bool) interface{} {
	switch reflect.TypeOf(val) {
	case structValuePtr:
		var structValueSerializer = NewStructValueSerializer(val.(*data.StructValue), redactSecret)
		return structValueSerializer
	case stringValuePtr:
		var stringValueSerializer = NewStringValueSerializer(val.(*data.StringValue))
		return stringValueSerializer
	case integerValuePtr:
		var integerValueSerializer = NewIntegerValueSerializer(val.(*data.IntegerValue))
		return integerValueSerializer
	case doubleValuePtr:
		var doubleValueSerializer = NewDoubleValueSerializer(val.(*data.DoubleValue))
		return doubleValueSerializer
	case listValuePtr:
		var listValueSerializer = NewListValueSerializer(val.(*data.ListValue), redactSecret)
		return listValueSerializer
	case optionalValuePtr:
		var optionalValueSerializer = NewOptionalValueSerializer(val.(*data.OptionalValue), redactSecret)
		return optionalValueSerializer
	case errorValuePtr:
		return NewErrorValueSerializer(val.(*data.ErrorValue), redactSecret)
	case voidValuePtr:
		return NewVoidValueSerializer()
	case boolValuePtr:
		return NewBooleanValueSerializer(val.(*data.BooleanValue))
	case methodResultType:
		return NewMethodResultSerializer(val.(core.MethodResult), redactSecret)
	case blobValuePtr:
		return NewBlobValueSerializer(val.(*data.BlobValue))
	case secretValuePtr:
		if redactSecret {
			return NewRedactedSecretValueSerializer()
		} else {
			return NewSecretValueSerializer(val.(*data.SecretValue))
		}

	default:
		log.Infof("could not find serializer for %s", reflect.TypeOf(val))
		return nil
	}
}

type RedactedSecretValueSerializer struct {
}

func NewRedactedSecretValueSerializer() *RedactedSecretValueSerializer {
	return &RedactedSecretValueSerializer{}
}

func (rss *RedactedSecretValueSerializer) MarshalJSON() ([]byte, error) {
	var result = map[string]interface{}{}
	result["SECRET"] = "*redacted*"
	return json.Marshal(result)
}

type SecretValueSerializer struct {
	secretValue *data.SecretValue
}

func (svs *SecretValueSerializer) MarshalJSON() ([]byte, error) {
	var result = map[string]interface{}{}
	result["SECRET"] = svs.secretValue.Value()
	return json.Marshal(result)
}

func NewSecretValueSerializer(value *data.SecretValue) *SecretValueSerializer {
	return &SecretValueSerializer{secretValue: value}
}

type BlobValueSerializer struct {
	blobValue *data.BlobValue
}

func (bvs *BlobValueSerializer) MarshalJSON() ([]byte, error) {
	var result = map[string]interface{}{}
	base64EncodedString := base64.StdEncoding.EncodeToString(bvs.blobValue.Value())
	result["BINARY"] = base64EncodedString
	return json.Marshal(result)
}

func NewBlobValueSerializer(value *data.BlobValue) *BlobValueSerializer {
	return &BlobValueSerializer{blobValue: value}
}

type JsonRpcRequestErrorSerializer struct {
	jsonRpcRequestError *JsonRpcRequestError
}

func NewJsonRpcRequestErrorSerializer(requestError *JsonRpcRequestError) *JsonRpcRequestErrorSerializer {
	return &JsonRpcRequestErrorSerializer{jsonRpcRequestError: requestError}
}

func (j *JsonRpcRequestErrorSerializer) MarshalJSON() ([]byte, error) {

	var result = make(map[string]interface{})
	jsonRpcRequest := j.jsonRpcRequestError.jsonRpc20Request
	if jsonRpcRequest != nil {
		result[lib.JSONRPC] = jsonRpcRequest.version
		if !j.jsonRpcRequestError.jsonRpc20Request.notification {
			result[lib.JSONRPC_ID] = j.jsonRpcRequestError.jsonRpc20Request.id
		}
	}
	jsonRpc20Error := j.jsonRpcRequestError.jsonRpc20Error
	if jsonRpc20Error != nil {
		var error = make(map[string]interface{})
		error[lib.ERROR_CODE] = jsonRpc20Error.Code()
		error[lib.ERROR_MESSAGE] = jsonRpc20Error.Message()
		if jsonRpc20Error.data != nil {
			error[lib.ERROR_DATA] = jsonRpc20Error.data
		}
		result[lib.METHOD_RESULT_ERROR] = error
	}

	return json.Marshal(result)
}

type StructValueSerializer struct {
	structValue  *data.StructValue
	redactSecret bool
}

func NewStructValueSerializer(value *data.StructValue, redactSecret bool) *StructValueSerializer {
	return &StructValueSerializer{structValue: value, redactSecret: redactSecret}
}

func (svs *StructValueSerializer) MarshalJSON() ([]byte, error) {
	var itemMapS = make(map[string]interface{})
	for key, val := range svs.structValue.Fields() {
		itemMapS[key] = getCustomSerializer(val, svs.redactSecret)
	}
	var resultMap = make(map[string]interface{})
	resultMap[svs.structValue.Name()] = itemMapS
	var wrapperMap = make(map[string]interface{})
	wrapperMap[svs.structValue.Type().String()] = resultMap
	return json.Marshal(wrapperMap)
}

type StringValueSerializer struct {
	stringValue *data.StringValue
}

func NewStringValueSerializer(value *data.StringValue) *StringValueSerializer {
	return &StringValueSerializer{stringValue: value}
}

func (svs *StringValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(svs.stringValue.Value())
}

type IntegerValueSerializer struct {
	integerValue *data.IntegerValue
}

func (svs *IntegerValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(svs.integerValue.Value())
}

func NewIntegerValueSerializer(value *data.IntegerValue) *IntegerValueSerializer {

	return &IntegerValueSerializer{integerValue: value}
}

type DoubleValueSerializer struct {
	doubleValue *data.DoubleValue
}

func (svs *DoubleValueSerializer) MarshalJSON() ([]byte, error) {
	var buf []byte
	splitResult := strings.Split(svs.doubleValue.String(), ".")
	var prec int
	if len(splitResult) > 1 {
		prec = len(splitResult[1])
	} else {
		//setting minimum precision to 1
		prec = 1
	}
	buf = strconv.AppendFloat(buf, svs.doubleValue.Value(), 'f', prec, 64)
	return buf, nil
}

func NewDoubleValueSerializer(value *data.DoubleValue) *DoubleValueSerializer {
	return &DoubleValueSerializer{doubleValue: value}
}

type ListValueSerializer struct {
	listValue    *data.ListValue
	redactSecret bool
}

func (lvs *ListValueSerializer) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, 0)
	for _, element := range lvs.listValue.List() {
		result = append(result, getCustomSerializer(element, lvs.redactSecret))
	}
	return json.Marshal(result)
}
func NewListValueSerializer(value *data.ListValue, redactSecret bool) *ListValueSerializer {
	return &ListValueSerializer{listValue: value, redactSecret: redactSecret}
}

type OptionalValueSerializer struct {
	optionalValue *data.OptionalValue
	redactSecret  bool
}

func (ovs *OptionalValueSerializer) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	if ovs.optionalValue.IsSet() {
		result[ovs.optionalValue.Type().String()] = getCustomSerializer(ovs.optionalValue.Value(), ovs.redactSecret)
	} else {
		result[ovs.optionalValue.Type().String()] = nil
	}

	return json.Marshal(result)
}

func NewOptionalValueSerializer(value *data.OptionalValue, redactSecret bool) *OptionalValueSerializer {
	return &OptionalValueSerializer{optionalValue: value, redactSecret: redactSecret}
}

type ErrorValueSerializer struct {
	errorValue   *data.ErrorValue
	redactSecret bool
}

func (evs *ErrorValueSerializer) MarshalJSON() ([]byte, error) {
	var itemMapS = make(map[string]interface{})
	for key, val := range evs.errorValue.Fields() {
		itemMapS[key] = getCustomSerializer(val, evs.redactSecret)
	}
	var resultMap = make(map[string]interface{})
	resultMap[evs.errorValue.Name()] = itemMapS
	var wrapperMap = make(map[string]interface{})
	wrapperMap[evs.errorValue.Type().String()] = resultMap
	return json.Marshal(wrapperMap)
}

func NewErrorValueSerializer(value *data.ErrorValue, redactSecret bool) *ErrorValueSerializer {
	return &ErrorValueSerializer{errorValue: value, redactSecret: redactSecret}
}

type VoidValueSerializer struct {
}

func (vvs *VoidValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
func NewVoidValueSerializer() *VoidValueSerializer {
	return &VoidValueSerializer{}
}

type BooleanValueSerializer struct {
	booleanValue *data.BooleanValue
}

func (bvs *BooleanValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(bvs.booleanValue.Value())
}

func NewBooleanValueSerializer(value *data.BooleanValue) *BooleanValueSerializer {
	return &BooleanValueSerializer{booleanValue: value}
}

type MethodResultSerializer struct {
	methodResult core.MethodResult
	redactSecret bool
}

func NewMethodResultSerializer(methodResult core.MethodResult, redactSecret bool) *MethodResultSerializer {
	return &MethodResultSerializer{methodResult: methodResult, redactSecret: redactSecret}
}

func (methodResultSerializer *MethodResultSerializer) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	var methodResult = methodResultSerializer.methodResult
	if methodResult.IsSuccess() {
		result[lib.METHOD_RESULT_OUTPUT] = getCustomSerializer(methodResult.Output(), methodResultSerializer.redactSecret)
	} else {
		result[lib.METHOD_RESULT_ERROR] = getCustomSerializer(methodResult.Error(), methodResultSerializer.redactSecret)
	}
	return json.Marshal(result)
}

type ApplicationContextSerializer struct {
	appContext *core.ApplicationContext
}

func NewApplicationContextSerializer(appContext *core.ApplicationContext) *ApplicationContextSerializer {
	return &ApplicationContextSerializer{appContext: appContext}
}

func (acs *ApplicationContextSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(acs.appContext)
}

type SecurityContextSerializer struct {
	securityContext core.SecurityContext
}

func NewSecurityContextSerializer(context core.SecurityContext) *SecurityContextSerializer {
	return &SecurityContextSerializer{securityContext: context}
}

func (scs *SecurityContextSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(scs.securityContext)
}

type ExecutionContextSerializer struct {
	executionContext *core.ExecutionContext
}

func NewExecutionContextSerializer(executionContext *core.ExecutionContext) *ExecutionContextSerializer {
	return &ExecutionContextSerializer{executionContext: executionContext}
}

func (ecs *ExecutionContextSerializer) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	if ecs.executionContext.SecurityContext() != nil {
		result[lib.SECURITY_CONTEXT] = NewSecurityContextSerializer(ecs.executionContext.SecurityContext())
	}
	result[lib.APPLICATION_CONTEXT] = NewApplicationContextSerializer(ecs.executionContext.ApplicationContext())
	return json.Marshal(result)
}

func dispatcher(val interface{}, redactSecret bool) ([]byte, error) {

	switch reflect.TypeOf(val) {
	case structValuePtr:
		var structValueSerializer = NewStructValueSerializer(val.(*data.StructValue), redactSecret)
		return json.Marshal(structValueSerializer)
	case listValuePtr:
		var listValueSerializer = NewListValueSerializer(val.(*data.ListValue), redactSecret)
		return json.Marshal(listValueSerializer)
	case optionalValuePtr:
		var optionalValueSerializer = NewOptionalValueSerializer(val.(*data.OptionalValue), redactSecret)
		return json.Marshal(optionalValueSerializer)
	case errorValuePtr:
		var errorValueSerializer = NewErrorValueSerializer(val.(*data.ErrorValue), redactSecret)
		return json.Marshal(errorValueSerializer)
	case methodResultType:
		return json.Marshal(NewMethodResultSerializer(val.(core.MethodResult), redactSecret))
	case executionContextPtr:
		return json.Marshal(NewExecutionContextSerializer(val.(*core.ExecutionContext)))
	case jsonRpcRequestType:
		return json.Marshal(NewJsonRpc20RequestSerializer(val.(JsonRpc20Request)))
	case jsonRpcResponseType:
		return json.Marshal(NewJsonRpc20ResponseSerializer(val.(JsonRpc20Response), redactSecret))
	case jsonRpcRequestErrorType:
		return json.Marshal(NewJsonRpcRequestErrorSerializer(val.(*JsonRpcRequestError)))
	default:
		fmt.Println(reflect.TypeOf(val))
	}
	return nil, errors.New("Could not Serialize")
}

type JsonRpc20RequestSerializer struct {
	jsonRequest JsonRpc20Request
}

func NewJsonRpc20RequestSerializer(jsonRequest JsonRpc20Request) JsonRpc20RequestSerializer {
	return JsonRpc20RequestSerializer{jsonRequest: jsonRequest}
}

func (jsonreq JsonRpc20RequestSerializer) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	result[lib.JSONRPC] = jsonreq.jsonRequest.version
	result[lib.JSONRPC_METHOD] = jsonreq.jsonRequest.method
	result[lib.JSONRPC_PARAMS] = jsonreq.jsonRequest.params
	if !jsonreq.jsonRequest.notification {
		result[lib.JSONRPC_ID] = jsonreq.jsonRequest.id
	}
	return json.Marshal(result)
}

type JsonRpc20ResponseSerializer struct {
	jsonResponse JsonRpc20Response
	redactSecret bool
}

func NewJsonRpc20ResponseSerializer(response JsonRpc20Response, redactSecret bool) JsonRpc20ResponseSerializer {
	return JsonRpc20ResponseSerializer{response, redactSecret}
}

func (jsonResp JsonRpc20ResponseSerializer) MarshalJSON() ([]byte, error) {
	var result = make(map[string]interface{})
	result[lib.JSONRPC] = jsonResp.jsonResponse.version
	result[lib.JSONRPC_ID] = jsonResp.jsonResponse.id
	if jsonResp.jsonResponse.result != nil {
		result[lib.METHOD_RESULT] = getCustomSerializer(jsonResp.jsonResponse.result, jsonResp.redactSecret)
	} else {
		result[lib.METHOD_RESULT_ERROR] = jsonResp.jsonResponse.error
		if jsonResp.jsonResponse.error[lib.ERROR_CODE] == string(JSONRPC_PARSE_ERROR) ||
			jsonResp.jsonResponse.error[lib.ERROR_CODE] == JSONRPC_INVALID_REQUEST {
			result[lib.JSONRPC_ID] = nil
		}
	}
	return json.Marshal(result)

}
