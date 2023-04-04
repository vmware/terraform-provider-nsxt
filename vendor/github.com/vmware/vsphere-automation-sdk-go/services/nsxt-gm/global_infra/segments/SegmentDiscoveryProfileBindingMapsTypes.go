// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: SegmentDiscoveryProfileBindingMaps.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package segments

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_global_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"reflect"
)

func segmentDiscoveryProfileBindingMapsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func SegmentDiscoveryProfileBindingMapsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func segmentDiscoveryProfileBindingMapsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infra_segment_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infraSegmentId"] = vapiBindings_.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = vapiBindings_.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/global-manager/api/v1/global-infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func SegmentDiscoveryProfileBindingMapsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
}

func segmentDiscoveryProfileBindingMapsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infra_segment_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infraSegmentId"] = vapiBindings_.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = vapiBindings_.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/global-manager/api/v1/global-infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func SegmentDiscoveryProfileBindingMapsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapListResultBindingType)
}

func segmentDiscoveryProfileBindingMapsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["infra_segment_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["infraSegmentId"] = vapiBindings_.NewStringType()
	pathParams["infra_segment_id"] = "infraSegmentId"
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/global-manager/api/v1/global-infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func SegmentDiscoveryProfileBindingMapsPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func segmentDiscoveryProfileBindingMapsPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infra_segment_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	paramsTypeMap["infraSegmentId"] = vapiBindings_.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = vapiBindings_.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/global-manager/api/v1/global-infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentDiscoveryProfileBindingMapsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func SegmentDiscoveryProfileBindingMapsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
}

func segmentDiscoveryProfileBindingMapsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra_segment_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	fields["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	fieldNameMap["infra_segment_id"] = "InfraSegmentId"
	fieldNameMap["segment_discovery_profile_binding_map_id"] = "SegmentDiscoveryProfileBindingMapId"
	fieldNameMap["segment_discovery_profile_binding_map"] = "SegmentDiscoveryProfileBindingMap"
	paramsTypeMap["segment_discovery_profile_binding_map_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["infra_segment_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["segment_discovery_profile_binding_map"] = vapiBindings_.NewReferenceType(nsx_global_policyModel.SegmentDiscoveryProfileBindingMapBindingType)
	paramsTypeMap["infraSegmentId"] = vapiBindings_.NewStringType()
	paramsTypeMap["segmentDiscoveryProfileBindingMapId"] = vapiBindings_.NewStringType()
	pathParams["segment_discovery_profile_binding_map_id"] = "segmentDiscoveryProfileBindingMapId"
	pathParams["infra_segment_id"] = "infraSegmentId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/global-manager/api/v1/global-infra/segments/{infraSegmentId}/segment-discovery-profile-binding-maps/{segmentDiscoveryProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
