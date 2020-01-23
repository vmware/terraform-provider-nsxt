/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest



import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/cleanjson"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

type Request struct {
	urlPath      string
	inputHeaders map[string]string
	requestBody  string
}

func NewRequest(urlPath string, inputHeaders map[string]string, requestBody string) *Request {
	return &Request{urlPath: urlPath, inputHeaders: inputHeaders, requestBody: requestBody}
}

func (result *Request) URLPath() string {
	return result.urlPath
}

func (result *Request) SetURLPath(path string) {
	result.urlPath = path
}

func (result *Request) InputHeaders() map[string]string {
	return result.inputHeaders
}

func (result *Request) SetInputHeaders(headers map[string]string) {
	for k, v := range headers {
		result.inputHeaders[k] = v
	}
}

func (result *Request) RequestBody() string {
	return result.requestBody
}

func (result *Request) SetRequestBody(body string) {
	result.requestBody = body
}

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
	result.SetInputHeaders(authzHeaders)

	return result, nil
}

func createParamsFieldMap( inputValue *data.StructValue, fieldsMapBindingtype map[string]bindings.BindingType, fieldNamesMap map[string]string) (map[string]string, error) {
	fields := map[string]string{}
	for field, fieldName := range fieldNamesMap {
		if val, err := getNestedParams(field, inputValue, fieldsMapBindingtype); err == nil {
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
	var bodyfieldAnnotations bool = false

	fieldsMapBindingtype := metadata.Fields()

	// Add @Body Field
	requestBodyParam := metadata.BodyParamActualName()
	if requestBodyParam != "" {
		requestBodyInputValue, err = getNestedParams(requestBodyParam, inputValue, fieldsMapBindingtype)
		if err != nil {
			log.Debugf("Body Param was not found, Error: %s ", err)
			return nil, err
		}
	} else {
		requestBodyInputValue = data.NewStructValue(inputValue.Name(), map[string]data.DataValue{})
		bodyfieldAnnotations = true
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

	// Add @Query Fields
	queryFields, err := createParamsFieldMap(inputValue, fieldsMapBindingtype, metadata.QueryParamsNameMap())
	if err != nil {
		return nil, err
	}

	// Add @BodyField Fields
	if  bodyfieldAnnotations {
		requestFields := requestBodyInputValue.(*data.StructValue).FieldNames()
		for field, _ := range metadata.FieldNameMap() {
			if stringInSlice(field, requestFields) {
				continue
			} else if val, err := inputValue.Field(field); err == nil {
				requestBodyInputValue.(*data.StructValue).SetField(field, val)
			} else {
				log.Error("Input datavalue doesn't contain BodyField %s which is required ", field)
			}
		}
	}

	dataValueToJSONEncoder := cleanjson.NewDataValueToJsonEncoder()
	requestBodyStr, err := dataValueToJSONEncoder.Encode(requestBodyInputValue)
	if err != nil {
		return nil, err
	}

	urlPath := metadata.GetUrlPath(pathFields, queryFields, metadata.DispatchParam())
	return NewRequest(urlPath, headerFields, requestBodyStr), nil
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func replaceFieldValue(fields map[string]string, fieldName string, fieldValue data.DataValue) (map[string]string, error) {
	fieldStr, err := convertDataValueToString(fieldValue)
	if err != nil {
		return nil, err
	}
	if fieldStr != "" {
		fields[fieldName] = fieldStr
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
	case data.ListValuePtr:
		return dataValueToJSONEncoder.Encode(dataValue)
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
	encodedStr, err = strconv.Unquote(encodedStr)
	if err != nil {
		log.Error(err)
		return encodedStr, err
	}
	return encodedStr, nil
}

// Get the authorization headers for the corresponding security context
func getAuthorizationHeaders(securityContext core.SecurityContext) (map[string]string, error) {
	if securityContext == nil {
		log.Info("securityContext is nil")
		return map[string]string{}, nil
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

// Get the authorization headers for username security context
func getUsernameCtxHeaders(securityContext core.SecurityContext) (map[string]string, error) {
	username, err := getSecurityCtxStrValue(securityContext, security.USER_KEY)
	if err != nil {
		return nil, err
	}

	password, err := getSecurityCtxStrValue(securityContext, security.PASSWORD_KEY)
	if err != nil {
		return nil, err
	}

	credentialString := fmt.Sprintf("%s:%s", username, password)
	base64EncodedVal := base64.StdEncoding.EncodeToString([]byte(credentialString))
	return map[string]string{"Authorization": "Basic " + base64EncodedVal}, nil
}

// Get the authorization headers for session security context
func getSessionCtxHeaders(securityContext core.SecurityContext) (map[string]string, error) {
	sessionID, err := getSecurityCtxStrValue(securityContext, security.SESSION_ID)
	if err != nil {
		return nil, err
	}
	return map[string]string{security.SESSION_ID_KEY: sessionID}, nil
}

//  Get the authorization headers for oauth security context.
func getOauthCtxHeaders(securityContext core.SecurityContext) (map[string]string, error) {
	oauthToken, err := getSecurityCtxStrValue(securityContext, security.ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	return map[string]string{security.CSP_AUTH_TOKEN_KEY: oauthToken}, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getSecurityCtxStrValue(securityContext core.SecurityContext, propKey string) (string, error) {
	securityContextMap := securityContext.GetAllProperties()
	var propVal interface{}
	if propValue, ok := securityContextMap[propKey]; !ok {
		err := fmt.Errorf("%s is not present in the security context", propKey)
		log.Error(err)
		return "", err
	} else {
		propVal = propValue
	}

	if propVal == nil {
		log.Debugf("%s value is nil, use \"\" instead", propKey)
		return "", nil
	}

	if propValueStr, ok := propVal.(string); ok {
		return propValueStr, nil
	}

	err := fmt.Errorf("Invalid type for %s, expected type string, actual type %s",
		propKey, reflect.TypeOf(propVal).String())
	log.Error(err)
	return "", err
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
		if temp, ok := optVal.Value().(*data.StructValue); ok &&  reflect.TypeOf(optElementBinding) == bindings.ReferenceBindingType {
			return getNestedParam(tokens[1:], temp, optElementBinding.(bindings.ReferenceType).Resolve())
		}
	} else if fieldBindingType == bindings.OptionalBindingType {
		return data.NewOptionalValue(nil), nil
	}
	return nil, l10n.NewRuntimeError( "vapi.data.serializers.rest.nested.invalid.args",
		map[string]string{"param": strings.Join(tokens, ".")})
}