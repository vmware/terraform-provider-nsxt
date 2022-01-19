// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: UiViews.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"reflect"
)

func uiViewsCreateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["view"] = bindings.NewReferenceType(model.ViewBindingType)
	fieldNameMap["view"] = "View"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func uiViewsCreateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ViewBindingType)
}

func uiViewsCreateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["view"] = bindings.NewReferenceType(model.ViewBindingType)
	fieldNameMap["view"] = "View"
	paramsTypeMap["view"] = bindings.NewReferenceType(model.ViewBindingType)
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"view",
		"POST",
		"/global-manager/api/v1/ui-views",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func uiViewsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["view_id"] = bindings.NewStringType()
	fieldNameMap["view_id"] = "ViewId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func uiViewsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func uiViewsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["view_id"] = bindings.NewStringType()
	fieldNameMap["view_id"] = "ViewId"
	paramsTypeMap["view_id"] = bindings.NewStringType()
	paramsTypeMap["viewId"] = bindings.NewStringType()
	pathParams["view_id"] = "viewId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"",
		"DELETE",
		"/global-manager/api/v1/ui-views/{viewId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func uiViewsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tag"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["view_ids"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["widget_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tag"] = "Tag"
	fieldNameMap["view_ids"] = "ViewIds"
	fieldNameMap["widget_id"] = "WidgetId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func uiViewsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ViewListBindingType)
}

func uiViewsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tag"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["view_ids"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["widget_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tag"] = "Tag"
	fieldNameMap["view_ids"] = "ViewIds"
	fieldNameMap["widget_id"] = "WidgetId"
	paramsTypeMap["widget_id"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["tag"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["view_ids"] = bindings.NewOptionalType(bindings.NewStringType())
	queryParams["view_ids"] = "view_ids"
	queryParams["widget_id"] = "widget_id"
	queryParams["tag"] = "tag"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"",
		"GET",
		"/global-manager/api/v1/ui-views",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func uiViewsGet0InputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["view_id"] = bindings.NewStringType()
	fieldNameMap["view_id"] = "ViewId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func uiViewsGet0OutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ViewBindingType)
}

func uiViewsGet0RestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["view_id"] = bindings.NewStringType()
	fieldNameMap["view_id"] = "ViewId"
	paramsTypeMap["view_id"] = bindings.NewStringType()
	paramsTypeMap["viewId"] = bindings.NewStringType()
	pathParams["view_id"] = "viewId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"",
		"GET",
		"/global-manager/api/v1/ui-views/{viewId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func uiViewsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["view_id"] = bindings.NewStringType()
	fields["view"] = bindings.NewReferenceType(model.ViewBindingType)
	fieldNameMap["view_id"] = "ViewId"
	fieldNameMap["view"] = "View"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func uiViewsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ViewBindingType)
}

func uiViewsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["view_id"] = bindings.NewStringType()
	fields["view"] = bindings.NewReferenceType(model.ViewBindingType)
	fieldNameMap["view_id"] = "ViewId"
	fieldNameMap["view"] = "View"
	paramsTypeMap["view_id"] = bindings.NewStringType()
	paramsTypeMap["view"] = bindings.NewReferenceType(model.ViewBindingType)
	paramsTypeMap["viewId"] = bindings.NewStringType()
	pathParams["view_id"] = "viewId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"view",
		"PUT",
		"/global-manager/api/v1/ui-views/{viewId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
