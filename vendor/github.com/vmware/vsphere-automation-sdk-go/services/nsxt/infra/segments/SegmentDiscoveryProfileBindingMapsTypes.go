// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: SegmentDiscoveryProfileBindingMaps.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func segmentDiscoveryProfileBindingMapsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentDiscoveryProfileBindingMapsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentDiscoveryProfileBindingMapsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["infra_segment_id"] = bindings.NewStringType()
	paramsTypeMap["infraSegmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
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
		"/policy/api/v1/infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentDiscoveryProfileBindingMapsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
}

func segmentDiscoveryProfileBindingMapsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["infra_segment_id"] = bindings.NewStringType()
	paramsTypeMap["infraSegmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
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
		"/policy/api/v1/infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentDiscoveryProfileBindingMapsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapListResultBindingType)
}

func segmentDiscoveryProfileBindingMapsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["infra_segment_id"] = bindings.NewStringType()
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["infraSegmentId"] = bindings.NewStringType()
	pathParams["infra_segment_id"] = "infraSegmentId"
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
		"/policy/api/v1/infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentDiscoveryProfileBindingMapsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentDiscoveryProfileBindingMapsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["infra_segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	paramsTypeMap["infraSegmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
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
		"segment_discovery_profile_binding_map",
		"PATCH",
		"/policy/api/v1/infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentDiscoveryProfileBindingMapsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
}

func segmentDiscoveryProfileBindingMapsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["infra_segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_discovery_profile_binding_map"] = bindings.NewReferenceType(model.SegmentDiscoveryProfileBindingMapBindingType)
	paramsTypeMap["infraSegmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
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
		"segment_discovery_profile_binding_map",
		"PUT",
		"/policy/api/v1/infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
