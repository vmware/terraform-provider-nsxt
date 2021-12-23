// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: DhcpStaticBindingConfigs.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"reflect"
)

func dhcpStaticBindingConfigsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dhcpStaticBindingConfigsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func dhcpStaticBindingConfigsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["binding_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["bindingId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["binding_id"] = "bindingId"
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/dhcp-static-binding-configs/{bindingId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func dhcpStaticBindingConfigsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dhcpStaticBindingConfigsGetOutputType() bindings.BindingType {
	return bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
}

func dhcpStaticBindingConfigsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["binding_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["bindingId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["binding_id"] = "bindingId"
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/dhcp-static-binding-configs/{bindingId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func dhcpStaticBindingConfigsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dhcpStaticBindingConfigsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.DhcpStaticBindingConfigListResultBindingType)
}

func dhcpStaticBindingConfigsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/dhcp-static-binding-configs",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func dhcpStaticBindingConfigsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fields["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	fieldNameMap["dhcp_static_binding_config"] = "DhcpStaticBindingConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dhcpStaticBindingConfigsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func dhcpStaticBindingConfigsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fields["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	fieldNameMap["dhcp_static_binding_config"] = "DhcpStaticBindingConfig"
	paramsTypeMap["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["binding_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["bindingId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["binding_id"] = "bindingId"
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
		"dhcp_static_binding_config",
		"PATCH",
		"/global-manager/api/v1/global-infra/segments/{segmentId}/dhcp-static-binding-configs/{bindingId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func dhcpStaticBindingConfigsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fields["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	fieldNameMap["dhcp_static_binding_config"] = "DhcpStaticBindingConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dhcpStaticBindingConfigsUpdateOutputType() bindings.BindingType {
	return bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
}

func dhcpStaticBindingConfigsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["binding_id"] = bindings.NewStringType()
	fields["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["binding_id"] = "BindingId"
	fieldNameMap["dhcp_static_binding_config"] = "DhcpStaticBindingConfig"
	paramsTypeMap["dhcp_static_binding_config"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.DhcpStaticBindingConfigBindingType)}, bindings.REST)
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["binding_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["bindingId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["binding_id"] = "bindingId"
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
		"dhcp_static_binding_config",
		"PUT",
		"/global-manager/api/v1/global-infra/segments/{segmentId}/dhcp-static-binding-configs/{bindingId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
