/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: IpfixDfwCollectorProfiles.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package global_infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func ipfixDfwCollectorProfilesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixDfwCollectorProfilesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixDfwCollectorProfilesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	paramsTypeMap["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	paramsTypeMap["ipfixDfwCollectorProfileId"] = bindings.NewStringType()
	pathParams["ipfix_dfw_collector_profile_id"] = "ipfixDfwCollectorProfileId"
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
		"/global-manager/api/v1/global-infra/ipfix-dfw-collector-profiles/{ipfixDfwCollectorProfileId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixDfwCollectorProfilesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixDfwCollectorProfilesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
}

func ipfixDfwCollectorProfilesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	paramsTypeMap["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	paramsTypeMap["ipfixDfwCollectorProfileId"] = bindings.NewStringType()
	pathParams["ipfix_dfw_collector_profile_id"] = "ipfixDfwCollectorProfileId"
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
		"/global-manager/api/v1/global-infra/ipfix-dfw-collector-profiles/{ipfixDfwCollectorProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixDfwCollectorProfilesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixDfwCollectorProfilesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXDFWCollectorProfileListResultBindingType)
}

func ipfixDfwCollectorProfilesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
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
		"/global-manager/api/v1/global-infra/ipfix-dfw-collector-profiles",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixDfwCollectorProfilesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fields["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	fieldNameMap["i_PFIXDFW_collector_profile"] = "IPFIXDFWCollectorProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixDfwCollectorProfilesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixDfwCollectorProfilesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fields["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	fieldNameMap["i_PFIXDFW_collector_profile"] = "IPFIXDFWCollectorProfile"
	paramsTypeMap["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	paramsTypeMap["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	paramsTypeMap["ipfixDfwCollectorProfileId"] = bindings.NewStringType()
	pathParams["ipfix_dfw_collector_profile_id"] = "ipfixDfwCollectorProfileId"
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
		"i_PFIXDFW_collector_profile",
		"PATCH",
		"/global-manager/api/v1/global-infra/ipfix-dfw-collector-profiles/{ipfixDfwCollectorProfileId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixDfwCollectorProfilesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fields["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	fieldNameMap["i_PFIXDFW_collector_profile"] = "IPFIXDFWCollectorProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixDfwCollectorProfilesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
}

func ipfixDfwCollectorProfilesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	fields["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	fieldNameMap["ipfix_dfw_collector_profile_id"] = "IpfixDfwCollectorProfileId"
	fieldNameMap["i_PFIXDFW_collector_profile"] = "IPFIXDFWCollectorProfile"
	paramsTypeMap["ipfix_dfw_collector_profile_id"] = bindings.NewStringType()
	paramsTypeMap["i_PFIXDFW_collector_profile"] = bindings.NewReferenceType(model.IPFIXDFWCollectorProfileBindingType)
	paramsTypeMap["ipfixDfwCollectorProfileId"] = bindings.NewStringType()
	pathParams["ipfix_dfw_collector_profile_id"] = "ipfixDfwCollectorProfileId"
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
		"i_PFIXDFW_collector_profile",
		"PUT",
		"/global-manager/api/v1/global-infra/ipfix-dfw-collector-profiles/{ipfixDfwCollectorProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


