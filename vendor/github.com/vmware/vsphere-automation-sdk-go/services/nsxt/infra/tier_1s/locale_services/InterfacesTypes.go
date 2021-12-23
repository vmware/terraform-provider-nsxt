// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Interfaces.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func interfacesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func interfacesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func interfacesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	paramsTypeMap["interface_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["interfaceId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["interface_id"] = "interfaceId"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/interfaces/{interfaceId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfacesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func interfacesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier1InterfaceBindingType)
}

func interfacesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	paramsTypeMap["interface_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["interfaceId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["interface_id"] = "interfaceId"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/interfaces/{interfaceId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfacesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func interfacesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier1InterfaceListResultBindingType)
}

func interfacesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/interfaces",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfacesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fields["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	fieldNameMap["tier1_interface"] = "Tier1Interface"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func interfacesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func interfacesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fields["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	fieldNameMap["tier1_interface"] = "Tier1Interface"
	paramsTypeMap["interface_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["interfaceId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["interface_id"] = "interfaceId"
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
		"tier1_interface",
		"PATCH",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/interfaces/{interfaceId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfacesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fields["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	fieldNameMap["tier1_interface"] = "Tier1Interface"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func interfacesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier1InterfaceBindingType)
}

func interfacesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["interface_id"] = bindings.NewStringType()
	fields["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["interface_id"] = "InterfaceId"
	fieldNameMap["tier1_interface"] = "Tier1Interface"
	paramsTypeMap["interface_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_interface"] = bindings.NewReferenceType(model.Tier1InterfaceBindingType)
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["interfaceId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["interface_id"] = "interfaceId"
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
		"tier1_interface",
		"PUT",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/interfaces/{interfaceId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
