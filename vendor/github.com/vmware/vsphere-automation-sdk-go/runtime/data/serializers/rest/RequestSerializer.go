/* Copyright © 2019, 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/cleanjson"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

type Request struct {
	urlPath      string
	inputHeaders map[string][]string
	requestBody  string
}

func NewRequest(urlPath string, inputHeaders map[string][]string, requestBody string) *Request {
	return &Request{urlPath: urlPath, inputHeaders: inputHeaders, requestBody: requestBody}
}

func (result *Request) URLPath() string {
	return result.urlPath
}

func (result *Request) SetURLPath(path string) {
	result.urlPath = path
}

func (result *Request) InputHeaders() map[string][]string {
	return result.inputHeaders
}

func (result *Request) setInputHeaders(headers map[string]interface{}) error {
	for k, v := range headers {
		if strV, ok := v.(string); ok {
			result.inputHeaders[k] = []string{strV}
		} else {
			return fmt.Errorf("Expected a string value, but received %s", reflect.TypeOf(v))
		}
	}
	return nil
}

func (result *Request) RequestBody() string {
	return result.requestBody
}

func (result *Request) SetRequestBody(body string) {
	result.requestBody = body
}

// SerializeRequests Deprecated: Use SerializeRequestsWithSecCtxSerializers() instead.
// SerializeRequests serializes a request as a REST request
// Return Request with urlPath, inputHeaders and requestBody
func SerializeRequests(inputValue *data.StructValue, ctx *core.ExecutionContext,
	metadata *protocol.OperationRestMetadata) (*Request, error) {
	result, err := SerializeInput(inputValue, metadata)
	if err != nil {
		log.Error(err)
		return result, err
	}
	var secCtx core.SecurityContext
	if ctx != nil {
		secCtx = ctx.SecurityContext()
	}

	authzHeaders, err := getAuthorizationHeaders(secCtx)
	if err != nil {
		return nil, err
	}

	if err = result.setInputHeaders(authzHeaders); err != nil {
		return nil, err
	}

	return result, nil
}

// SerializeRequestsWithSecCtxSerializers serializes a request into a REST request
// Return Request with urlPath, inputHeaders and requestBody
func SerializeRequestsWithSecCtxSerializers(inputValue *data.StructValue, execCtx *core.ExecutionContext,
	metadata *protocol.OperationRestMetadata, secCtxSerializer SecurityContextSerializer) (*Request, error) {
	result, err := SerializeInput(inputValue, metadata)
	if err != nil {
		log.Error(err)
		return result, err
	}

	var secCtx core.SecurityContext
	if execCtx != nil {
		secCtx = execCtx.SecurityContext()
	}

	var authzHeaders map[string]interface{}
	if secCtx != nil && secCtxSerializer != nil {
		authzHeaders, err = secCtxSerializer.Serialize(secCtx)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if err = result.setInputHeaders(authzHeaders); err != nil {
		return nil, err
	}
	return result, nil
}

func createParamsFieldMap(inputValue *data.StructValue, fieldsMapBindingtype map[string]bindings.BindingType, fieldNamesMap map[string]string) (map[string][]string, error) {
	fields := map[string][]string{}
	dataValueToJSONEncoder := cleanjson.NewDataValueToJsonEncoder()
	for field, fieldName := range fieldNamesMap {
		val, err := getNestedParams(field, inputValue, fieldsMapBindingtype)
		if err != nil {
			log.Debugf("Error in request serialization: %s ", err.Error())
			return nil, err
		}
		if reflect.TypeOf(val) == data.ListValuePtr {
			temp := []string{}
			for _, e := range val.(*data.ListValue).List() {
				elem, err := unquote(dataValueToJSONEncoder.Encode(e))
				if err != nil {
					return nil, err
				}
				temp = append(temp, elem)
			}
			fields[fieldName] = temp
		} else {
			fields, err = replaceFieldValue(fields, fieldName, val)
			if err != nil {
				return nil, err
			}
		}
	}
	return fields, nil
}

// SerializeInput serializes the input value
// Return Request with urlPath, inputHeaders and requestBody
func SerializeInput(inputValue *data.StructValue, metadata *protocol.OperationRestMetadata) (*Request, error) {
	if metadata == nil {
		return nil, l10n.NewRuntimeErrorNoParam("vapi.data.serializers.rest.metadata.value.nil")
	}

	var requestBodyInputValue data.DataValue
	var err error

	fieldsMapBindingtype := metadata.Fields()

	// Add @Body Field
	requestBodyParam := metadata.BodyParamActualName()
	if requestBodyParam != "" {
		requestBodyInputValue, err = getNestedParams(requestBodyParam, inputValue, fieldsMapBindingtype)
		if err != nil {
			log.Errorf("Error serializing Body Param, : %s ", err)
			return nil, err
		}
	} else {
		requestBodyInputValue = data.NewStructValue(inputValue.Name(), map[string]data.DataValue{})
	}

	// Add @Path Fields
	pathFields, err := createParamsFieldMap(inputValue, fieldsMapBindingtype, metadata.PathParamsNameMap())
	if err != nil {
		return nil, err
	}

	// Add @Header Fields
	headerFields, err := createParamsFieldMap(inputValue, fieldsMapBindingtype, metadata.HeaderParamsNameMap())
	if err != nil {
		return nil, err
	}

	// Add dispatch Header Fields
	dispatchHeaders := metadata.DispatchHeaderParams()
	for k, value := range dispatchHeaders {
		if _, ok := headerFields[k]; !ok || value != "" {
			headerFields[k] = []string{value}
		}
	}

	// Update Content-Type if consumes is provided
	consumes := metadata.OperationConsumes()
	if consumes != "" {
		headerFields["Content-Type"] = []string{consumes}
	}

	// Add @Query Fields
	queryFields, err := createParamsFieldMap(inputValue, fieldsMapBindingtype, metadata.QueryParamsNameMap())
	if err != nil {
		return nil, err
	}

	// Add @BodyField Fields
	bodyFieldsMap := metadata.BodyFieldsMap()
	for fieldName, paramName := range bodyFieldsMap {
		if val, err := inputValue.Field(paramName); err == nil {
			// BodyField wrapper will also be a StructValue type
			requestBodyInputValue.(*data.StructValue).SetField(fieldName, val)
		} else {
			log.Error("Input datavalue doesn't contain BodyField %s which is required ", paramName)
		}
	}

	dataValueToJSONEncoder := cleanjson.NewDataValueToJsonEncoder()

	var requestBodyStr string
	if consumes == lib.FORM_URL_CONTENT_TYPE {
		formURLBody, err := encodeDataValueToFormURL(requestBodyInputValue, dataValueToJSONEncoder)
		if err != nil {
			return nil, err
		} else if formURLBody != nil {
			requestBodyStr = *formURLBody
		}
	} else {
		requestBodyStr, err = dataValueToJSONEncoder.Encode(requestBodyInputValue)
		if err != nil {
			return nil, err
		}
	}
	/**
	 * ref: https://golang.org/pkg/encoding/json/#Marshal
	 *  A nil pointer encodes as the null JSON value
	 */
	if requestBodyStr == "null" {
		requestBodyStr = ""
	}

	// Process dispatch query param
	dispatchParams := queryEscapeDispatchParams(metadata.DispatchParam(), queryFields)

	urlPath := metadata.GetUrlPath(pathFields, queryFields, dispatchParams)
	return NewRequest(urlPath, headerFields, requestBodyStr), nil
}

func queryEscapeDispatchParams(dispatchParamsStr string, queryFields map[string][]string) string {
	dispatchParams := []string{}
	for _, p := range strings.Split(dispatchParamsStr, "&") {
		dispatchList := strings.Split(p, "=")
		temp := []string{}
		k := dispatchList[0]

		/**
		* Added to handle allowed characters
		* example: where `=` can be part of the value
		* https://opengrok.eng.vmware.com/source/xref/vapi-main.perforce.1666/vapi-core/idl-toolkit/vmodl-models/src/dist/models/test/vmodl/rest-native/VerbParams.vmodl#109
		 */
		for _, v := range dispatchList {
			temp = append(temp, url.QueryEscape(v))
		}

		// Add it to url only if key doesnt exist in query field
		if _, ok := queryFields[k]; !ok {
			dispatchParams = append(dispatchParams, strings.Join(temp, "="))
		}
	}
	return strings.Join(dispatchParams, "&")
}

func replaceFieldValue(fields map[string][]string, fieldName string, fieldValue data.DataValue) (map[string][]string, error) {
	fieldStr, err := convertDataValueToString(fieldValue)
	if err != nil {
		return nil, err
	}
	if fieldStr != "" {
		fields[fieldName] = []string{fieldStr}
	}
	return fields, nil
}

// Serialize data value to string for use in Request URL
func convertDataValueToString(dataValue data.DataValue) (string, error) {
	dataValueToJSONEncoder := cleanjson.NewDataValueToJsonEncoder()
	switch reflect.TypeOf(dataValue) {
	case data.OptionalValuePtr:
		optionalVal := dataValue.(*data.OptionalValue)
		if optionalVal.IsSet() {
			return convertDataValueToString(optionalVal.Value())
		}
		return "", nil
	case data.BoolValuePtr:
		return dataValueToJSONEncoder.Encode(dataValue)
	case data.StringValuePtr:
		return unquote(dataValueToJSONEncoder.Encode(dataValue))
	case data.IntegerValuePtr:
		return dataValueToJSONEncoder.Encode(dataValue)
	case data.DoubleValuePtr:
		return dataValueToJSONEncoder.Encode(dataValue)
	case data.BlobValuePtr:
		return unquote(dataValueToJSONEncoder.Encode(dataValue))
	case data.SecretValuePtr:
		return unquote(dataValueToJSONEncoder.Encode(dataValue))
	default:
		return "", l10n.NewRuntimeError(
			"vapi.data.serializers.rest.unsupported_data_value",
			map[string]string{"type": dataValue.Type().String()})
	}
}

func unquote(encodedStr string, err error) (string, error) {
	if err != nil {
		log.Error(err)
		return encodedStr, err
	}
	if strings.Contains(encodedStr, "\"") {
		encodedStr, err = strconv.Unquote(encodedStr)
		if err != nil {
			log.Debugf("Error in while unquoting : %s ", err.Error())
			return encodedStr, err
		}
	}
	return encodedStr, nil
}

// Form URL Encode Support
func encodeDataValueToFormURL(dataValue data.DataValue, dataValueToJSONEncoder *cleanjson.DataValueToJsonEncoder) (*string, error) {
	switch reflect.TypeOf(dataValue) {
	case data.OptionalValuePtr:
		optionalVal := dataValue.(*data.OptionalValue)
		if optionalVal.IsSet() {
			return encodeDataValueToFormURL(optionalVal.Value(), dataValueToJSONEncoder)
		}
		return nil, nil
	case data.BoolValuePtr:
		str, err := dataValueToJSONEncoder.Encode(dataValue)
		return &str, err
	case data.StringValuePtr:
		str, err := unquote(dataValueToJSONEncoder.Encode(dataValue))
		return &str, err
	case data.IntegerValuePtr:
		str, err := dataValueToJSONEncoder.Encode(dataValue)
		return &str, err
	case data.DoubleValuePtr:
		str, err := dataValueToJSONEncoder.Encode(dataValue)
		return &str, err
	case data.BlobValuePtr:
		str, err := unquote(dataValueToJSONEncoder.Encode(dataValue))
		return &str, err
	case data.SecretValuePtr:
		str, err := unquote(dataValueToJSONEncoder.Encode(dataValue))
		return &str, err
	case data.StructValuePtr:
		return structValueToFormURLEncode(dataValue.(*data.StructValue), dataValueToJSONEncoder)
	default:
		return nil, l10n.NewRuntimeError("vapi.data.serializers.rest.unhandled.datavalue.formurl",
			map[string]string{"datavalue": dataValue.Type().String()})
	}
}

func structValueToFormURLEncode(formURLStructVal *data.StructValue, dataValueToJSONEncoder *cleanjson.DataValueToJsonEncoder) (*string, error) {
	params := url.Values{}
	for k, v := range formURLStructVal.Fields() {
		val, err := encodeDataValueToFormURL(v, dataValueToJSONEncoder)
		if val != nil {
			params.Add(k, *val)
		} else if err != nil {
			return nil, l10n.NewRuntimeError("vapi.data.serializers.rest.formurl.field.error",
				map[string]string{"field": k, "msg": err.Error()})
		}
	}
	res := params.Encode()
	return &res, nil
}

// Deprecated: use Serialize(ctx core.SecurityContext) in SecurityContextSerializer instead
// Get the authorization headers for the corresponding security context
func getAuthorizationHeaders(securityContext core.SecurityContext) (map[string]interface{}, error) {
	if securityContext == nil {
		log.Info("securityContext is nil")
		return map[string]interface{}{}, nil
	}

	switch securityContext.Property(security.AUTHENTICATION_SCHEME_ID) {
	case security.USER_PASSWORD_SCHEME_ID:
		return getUsernameCtxHeaders(securityContext)
	case security.SESSION_SCHEME_ID:
		return getSessionCtxHeaders(securityContext)
	case security.OAUTH_SCHEME_ID:
		return getOauthCtxHeaders(securityContext)
	default:
		return nil, fmt.Errorf("Invalid authentication scheme ID, expected one of [%s, %s, %s], actual value %s",
			security.USER_PASSWORD_SCHEME_ID, security.SESSION_SCHEME_ID, security.OAUTH_SCHEME_ID,
			securityContext.Property(security.AUTHENTICATION_SCHEME_ID))
	}
}

// Deprecated: use UserPwdSecContextSerializer instead
// Get the authorization headers for username security context
func getUsernameCtxHeaders(securityContext core.SecurityContext) (map[string]interface{}, error) {
	username, err := GetSecurityCtxStrValue(securityContext, security.USER_KEY)
	if err != nil {
		return nil, err
	}

	password, err := GetSecurityCtxStrValue(securityContext, security.PASSWORD_KEY)
	if err != nil {
		return nil, err
	}

	credentialString := fmt.Sprintf("%s:%s", *username, *password)
	base64EncodedVal := base64.StdEncoding.EncodeToString([]byte(credentialString))
	return map[string]interface{}{"Authorization": "Basic " + base64EncodedVal}, nil
}

// Deprecated: use SessionSecContextSerializer instead
// Get the authorization headers for session security context
func getSessionCtxHeaders(securityContext core.SecurityContext) (map[string]interface{}, error) {
	sessionID, err := GetSecurityCtxStrValue(securityContext, security.SESSION_ID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{security.SESSION_ID_KEY: sessionID}, nil
}

// Deprecated: use OauthSecContextSerializer instead
//  Get the authorization headers for oauth security context.
func getOauthCtxHeaders(securityContext core.SecurityContext) (map[string]interface{}, error) {
	oauthToken, err := GetSecurityCtxStrValue(securityContext, security.ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{security.CSP_AUTH_TOKEN_KEY: oauthToken}, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetSecurityCtxStrValue returns value of the given key in *string.
// Error will be raised if securityContext is nil or value is not string type
func GetSecurityCtxStrValue(securityContext core.SecurityContext, propKey string) (*string, error) {
	if securityContext == nil {
		err := errors.New("securityContext can't be nil")
		return nil, err
	}

	securityContextMap := securityContext.GetAllProperties()
	var propVal interface{}
	var ok bool
	if propVal, ok = securityContextMap[propKey]; !ok {
		log.Debugf("%s is not present in the security context", propKey)
		return nil, nil
	}

	if propVal == nil {
		return nil, nil
	}

	if propValueStr, ok := propVal.(string); ok {
		return &propValueStr, nil
	}

	err := fmt.Errorf("Invalid type for %s, expected type string, actual type %s",
		propKey, reflect.TypeOf(propVal).String())
	log.Error(err)
	return nil, err
}

func getNestedParams(field string, inputValue *data.StructValue, fields map[string]bindings.BindingType) (data.DataValue, error) {
	tokens := strings.Split(field, ".")
	if len(tokens) > 1 {
		currentField := tokens[0]
		if fieldBinding, ok := fields[currentField]; ok {
			return getNestedParam(tokens, inputValue, fieldBinding)
		}
		return nil, l10n.NewRuntimeError("vapi.data.serializers.rest.nested.invalid.args",
			map[string]string{"param": strings.Join(tokens, ".")})
	}
	return inputValue.Field(field)
}

func getNestedParam(tokens []string, inputValue *data.StructValue, fieldBinding bindings.BindingType) (data.DataValue, error) {
	currentField := tokens[0]
	if len(tokens) == 1 {
		return inputValue.Field(currentField)
	}

	newInputValue, err := inputValue.Field(currentField)
	if err != nil {
		return nil, err
	}
	fieldBindingType := reflect.TypeOf(fieldBinding)
	currentField = tokens[1]

	if temp, ok := newInputValue.(*data.StructValue); ok {
		if fieldBindingType == bindings.ReferenceBindingType {
			return getNestedParam(tokens[1:], temp, fieldBinding.(bindings.ReferenceType).Resolve().(bindings.StructType).Field(currentField))
		} else if fieldBindingType == bindings.StructBindingType {
			return getNestedParam(tokens[1:], temp, fieldBinding.(bindings.StructType).Field(currentField))
		}
	} else if optVal, ok := newInputValue.(*data.OptionalValue); ok && fieldBindingType == bindings.OptionalBindingType {
		if !optVal.IsSet() || len(tokens) == 1 {
			return optVal, nil
		}

		optElementBinding := fieldBinding.(bindings.OptionalType).ElementType()
		if temp, ok := optVal.Value().(*data.StructValue); ok && reflect.TypeOf(optElementBinding) == bindings.ReferenceBindingType {
			structBindingType := optElementBinding.(bindings.ReferenceType).Resolve().(bindings.StructType)
			return getNestedParam(tokens[1:], temp, structBindingType.Field(currentField))
		}
	} else if fieldBindingType == bindings.OptionalBindingType {
		return data.NewOptionalValue(nil), nil
	}
	return nil, l10n.NewRuntimeError("vapi.data.serializers.rest.nested.invalid.args",
		map[string]string{"param": strings.Join(tokens, ".")})
}
