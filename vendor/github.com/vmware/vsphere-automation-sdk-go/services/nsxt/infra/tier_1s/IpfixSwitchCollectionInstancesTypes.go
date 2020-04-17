/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: IpfixSwitchCollectionInstances.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package tier_1s

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func ipfixSwitchCollectionInstancesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixSwitchCollectionInstancesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixSwitchCollectionInstancesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	paramsTypeMap["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["ipfixSwitchCollectionInstanceId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["ipfix_switch_collection_instance_id"] = "ipfixSwitchCollectionInstanceId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/ipfix-switch-collection-instances/{ipfixSwitchCollectionInstanceId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixSwitchCollectionInstancesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixSwitchCollectionInstancesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
}

func ipfixSwitchCollectionInstancesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	paramsTypeMap["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["ipfixSwitchCollectionInstanceId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["ipfix_switch_collection_instance_id"] = "ipfixSwitchCollectionInstanceId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/ipfix-switch-collection-instances/{ipfixSwitchCollectionInstanceId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixSwitchCollectionInstancesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixSwitchCollectionInstancesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceListResultBindingType)
}

func ipfixSwitchCollectionInstancesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
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
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/ipfix-switch-collection-instances",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixSwitchCollectionInstancesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fields["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	fieldNameMap["i_PFIX_switch_collection_instance"] = "IPFIXSwitchCollectionInstance"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixSwitchCollectionInstancesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixSwitchCollectionInstancesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fields["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	fieldNameMap["i_PFIX_switch_collection_instance"] = "IPFIXSwitchCollectionInstance"
	paramsTypeMap["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["ipfixSwitchCollectionInstanceId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["ipfix_switch_collection_instance_id"] = "ipfixSwitchCollectionInstanceId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"i_PFIX_switch_collection_instance",
		"PATCH",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/ipfix-switch-collection-instances/{ipfixSwitchCollectionInstanceId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipfixSwitchCollectionInstancesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fields["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	fieldNameMap["i_PFIX_switch_collection_instance"] = "IPFIXSwitchCollectionInstance"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixSwitchCollectionInstancesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
}

func ipfixSwitchCollectionInstancesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	fields["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["ipfix_switch_collection_instance_id"] = "IpfixSwitchCollectionInstanceId"
	fieldNameMap["i_PFIX_switch_collection_instance"] = "IPFIXSwitchCollectionInstance"
	paramsTypeMap["ipfix_switch_collection_instance_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["i_PFIX_switch_collection_instance"] = bindings.NewReferenceType(model.IPFIXSwitchCollectionInstanceBindingType)
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["ipfixSwitchCollectionInstanceId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["ipfix_switch_collection_instance_id"] = "ipfixSwitchCollectionInstanceId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"i_PFIX_switch_collection_instance",
		"PUT",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/ipfix-switch-collection-instances/{ipfixSwitchCollectionInstanceId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


