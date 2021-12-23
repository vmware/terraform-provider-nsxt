// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Areas.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package ospf

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func areasDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func areasDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func areasDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["area_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["areaId"] = bindings.NewStringType()
	pathParams["area_id"] = "areaId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf/areas/{areaId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func areasGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func areasGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfAreaConfigBindingType)
}

func areasGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["area_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["areaId"] = bindings.NewStringType()
	pathParams["area_id"] = "areaId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf/areas/{areaId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func areasListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func areasListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfAreaConfigListResultBindingType)
}

func areasListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf/areas",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func areasPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fields["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	fieldNameMap["ospf_area_config"] = "OspfAreaConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func areasPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfAreaConfigBindingType)
}

func areasPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fields["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	fieldNameMap["ospf_area_config"] = "OspfAreaConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["area_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["areaId"] = bindings.NewStringType()
	pathParams["area_id"] = "areaId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"ospf_area_config",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf/areas/{areaId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func areasUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fields["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	fieldNameMap["ospf_area_config"] = "OspfAreaConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func areasUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfAreaConfigBindingType)
}

func areasUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["area_id"] = bindings.NewStringType()
	fields["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["area_id"] = "AreaId"
	fieldNameMap["ospf_area_config"] = "OspfAreaConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["ospf_area_config"] = bindings.NewReferenceType(model.OspfAreaConfigBindingType)
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["area_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["areaId"] = bindings.NewStringType()
	pathParams["area_id"] = "areaId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"ospf_area_config",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf/areas/{areaId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
