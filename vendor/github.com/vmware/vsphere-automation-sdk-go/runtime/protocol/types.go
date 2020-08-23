/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package protocol

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
)

type OperationMetadata struct {
	methodDefinition *core.MethodDefinition
	inputType        bindings.StructType
	outputType       bindings.BindingType
	errorBindingMap  map[string]bindings.BindingType
	restMetadata     OperationRestMetadata
}

func NewOperationMetadata(methodDefinition *core.MethodDefinition,
	inputType bindings.StructType,
	outputType bindings.BindingType,
	errorBindingMap map[string]bindings.BindingType,
	restMetadata OperationRestMetadata) OperationMetadata {
	return OperationMetadata{methodDefinition: methodDefinition,
		inputType:       inputType,
		outputType:      outputType,
		errorBindingMap: errorBindingMap,
		restMetadata:    restMetadata}
}

func (meta OperationMetadata) MethodDefinition() *core.MethodDefinition {
	return meta.methodDefinition
}

func (meta OperationMetadata) InputType() bindings.StructType {
	return meta.inputType
}

func (meta OperationMetadata) OutputType() bindings.BindingType {
	return meta.outputType
}

func (meta OperationMetadata) ErrorBindingMap() map[string]bindings.BindingType {
	return meta.errorBindingMap
}

func (meta OperationMetadata) RestMetadata() OperationRestMetadata {
	return meta.restMetadata
}

// fields and fieldsNameMap defines the bindingtype and name of field respectively of @BodyField annotation
// Rest metadata for name and types of query, header and
// body parameters of an operation. Example:
//	meta.ParamsTypeMap["input.nested.bparam"] = bindings.NewListType(bindings.NewStringType(), reflect.TypeOf([]string{}))
//	meta.ParamsTypeMap["input.nested.hparam"] = bindings.NewStringType()
//	meta.ParamsTypeMap["input.nested.qparam"] = bindings.NewStringType()
//	meta.QueryParams["qparam"] = "input.nested.qparam"
//	meta.HeaderParams["Hparam"] = "input.nested.hparam"
//	meta.BodyParam = "input.nested.bparam"
//	httpMethod = "GET|POST|UPDATE|PATCH|DELETE"
//	urlTemplate = "/newannotations/properties/{id}"
type OperationRestMetadata struct {
	// Name of all the field name wrappers that should be present in Data Input Value
	fields       map[string]bindings.BindingType
	fieldNameMap map[string]string
	// Flattened types of all parameters. Key is fully qualified field name
	paramsTypeMap map[string]bindings.BindingType
	//Names of rest parameter to fully qualified canonical name of the field
	pathParamsNameMap    map[string]string
	queryParamsNameMap   map[string]string
	headerParamsNameMap  map[string]string
	dispatchHeaderParams map[string]string
	bodyFieldsMap        map[string]string
	//Encoded dispatch parameters
	dispatchParam string
	//Fully qualified field name canonical name of body param
	bodyParamActualName string
	//HTTP action for the operation
	httpMethod string
	//HTTP URL for the operation
	urlTemplate string
	// Content-Type that operation consumes
	operationConsumes string
	// HTTP response success code
	successCode int
	// Field name of response body
	responseBodyName string
	// vAPI error name to HTTP response error code mapping
	errorCodeMap map[string]int
	// Map from result field name to http header name
	resultHeadersNameMap map[string]string
	// Map from error field name to http header name
	errorHeadersNameMap map[string]map[string]string
}

func NewOperationRestMetadata(
	fields map[string]bindings.BindingType,
	fieldNameMap map[string]string,
	paramsTypeMap map[string]bindings.BindingType,
	pathParamsNameMap map[string]string,
	queryParamsNameMap map[string]string,
	headerParamsNameMap map[string]string,
	dispatchHeaderParams map[string]string,
	bodyFieldsMap map[string]string,
	dispatchParam string,
	bodyParamActualName string,
	httpMethod string,
	urlTemplate string,
	operationConsumes string,
	resultHeadersNameMap map[string]string,
	successCode int,
	responseBodyName string,
	errorHeadersNameMap map[string]map[string]string,
	errorCodeMap map[string]int) OperationRestMetadata {

	return OperationRestMetadata{
		fields:               fields,
		fieldNameMap:         fieldNameMap,
		paramsTypeMap:        paramsTypeMap,
		pathParamsNameMap:    pathParamsNameMap,
		queryParamsNameMap:   queryParamsNameMap,
		headerParamsNameMap:  headerParamsNameMap,
		dispatchHeaderParams: dispatchHeaderParams,
		bodyFieldsMap:        bodyFieldsMap,
		dispatchParam:        dispatchParam,
		bodyParamActualName:  bodyParamActualName,
		httpMethod:           httpMethod,
		urlTemplate:          urlTemplate,
		operationConsumes:    operationConsumes,
		successCode:          successCode,
		responseBodyName:     responseBodyName,
		errorCodeMap:         errorCodeMap,
		resultHeadersNameMap: resultHeadersNameMap,
		errorHeadersNameMap:  errorHeadersNameMap}
}

func (meta OperationRestMetadata) Fields() map[string]bindings.BindingType {
	return meta.fields
}

func (meta OperationRestMetadata) FieldNameMap() map[string]string {
	return meta.fieldNameMap
}

func (meta OperationRestMetadata) ParamsTypeMap() map[string]bindings.BindingType {
	return meta.paramsTypeMap
}

func (meta OperationRestMetadata) PathParamsNameMap() map[string]string {
	return meta.pathParamsNameMap
}

func (meta OperationRestMetadata) QueryParamsNameMap() map[string]string {
	return meta.queryParamsNameMap
}

func (meta OperationRestMetadata) HeaderParamsNameMap() map[string]string {
	return meta.headerParamsNameMap
}

func (meta OperationRestMetadata) DispatchHeaderParams() map[string]string {
	return meta.dispatchHeaderParams
}

func (meta OperationRestMetadata) BodyFieldsMap() map[string]string {
	return meta.bodyFieldsMap
}

func (meta OperationRestMetadata) DispatchParam() string {
	return meta.dispatchParam
}

func (meta OperationRestMetadata) BodyParamActualName() string {
	return meta.bodyParamActualName
}

func (meta OperationRestMetadata) HttpMethod() string {
	return meta.httpMethod
}

func (meta OperationRestMetadata) UrlTemplate() string {
	return meta.urlTemplate
}

func (meta OperationRestMetadata) OperationConsumes() string {
	return meta.operationConsumes
}

func (meta OperationRestMetadata) SuccessCode() int {
	return meta.successCode
}

func (meta OperationRestMetadata) ResponseBodyName() string {
	return meta.responseBodyName
}

func (meta OperationRestMetadata) ErrorCodeMap() map[string]int {
	return meta.errorCodeMap
}

func (meta OperationRestMetadata) ResultHeadersNameMap() map[string]string {
	return meta.resultHeadersNameMap
}
func (meta OperationRestMetadata) ErrorHeadersNameMap() map[string]map[string]string {
	return meta.errorHeadersNameMap
}

func (meta OperationRestMetadata) GetUrlPath(
	pathVariableFields map[string][]string, queryParamFields map[string][]string,
	dispatchParam string) string {
	urlPath := meta.urlTemplate
	// Substitute path variables with values in the template
	for fieldName, fields := range pathVariableFields {
		val := url.QueryEscape(strings.Join(fields, ""))
		urlPath = strings.Replace(urlPath, fmt.Sprintf("{%s}", fieldName), val, 1)
	}

	// Construct the query params portion of the url
	queryPrams := []string{}

	// Add dispatch parameter first if it presents
	if dispatchParam != "" {
		queryPrams = append(queryPrams, dispatchParam)
	}

	// Add other operation query parameters
	for fieldName, fieldStr := range queryParamFields {
		for _, e := range fieldStr {
			qpKey := url.QueryEscape(fieldName)
			qpVal := url.QueryEscape(e)
			qparam := fmt.Sprintf("%s=%s", qpKey, qpVal)
			queryPrams = append(queryPrams, qparam)
		}
	}
	queryParamStr := strings.Join(queryPrams, "&")

	if queryParamStr != "" {
		// Append the query params portion if it exists
		var connector string
		if strings.ContainsAny(urlPath, "?") {
			connector = "&"
		} else {
			connector = "?"
		}
		urlPath = strings.Join([]string{urlPath, queryParamStr}, connector)
	}
	return urlPath
}

func (meta OperationRestMetadata) PathVariableFieldNames() []string {
	return getListOfMapValues(meta.pathParamsNameMap)
}

func (meta OperationRestMetadata) QueryParamFieldNames() []string {
	return getListOfMapValues(meta.queryParamsNameMap)
}

func (meta OperationRestMetadata) HeaderFieldNames() []string {
	return getListOfMapValues(meta.headerParamsNameMap)
}

func (meta OperationRestMetadata) DispatchHeaderNames() []string {
	return getListOfMapValues(meta.dispatchHeaderParams)
}

func (meta OperationRestMetadata) BodyFieldsNames() []string {
	return getListOfMapValues(meta.bodyFieldsMap)
}

func getListOfMapValues(mapValue map[string]string) []string {
	fields := []string{}
	for k := range mapValue {
		fields = append(fields, k)
	}
	return fields
}
