// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: SegmentSecurityProfileBindingMaps.
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

func segmentSecurityProfileBindingMapsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfileBindingMapsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentSecurityProfileBindingMapsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["segment_security_profile_binding_map_id"] = "segmentSecurityProfileBindingMapId"
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/segment-security-profile-binding-maps/{segmentSecurityProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentSecurityProfileBindingMapsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfileBindingMapsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
}

func segmentSecurityProfileBindingMapsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["segment_security_profile_binding_map_id"] = "segmentSecurityProfileBindingMapId"
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/segment-security-profile-binding-maps/{segmentSecurityProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentSecurityProfileBindingMapsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfileBindingMapsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapListResultBindingType)
}

func segmentSecurityProfileBindingMapsListRestMetadata() protocol.OperationRestMetadata {
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
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
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
		"/global-manager/api/v1/global-infra/segments/{segmentId}/segment-security-profile-binding-maps",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentSecurityProfileBindingMapsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	fieldNameMap["segment_security_profile_binding_map"] = "SegmentSecurityProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfileBindingMapsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentSecurityProfileBindingMapsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	fieldNameMap["segment_security_profile_binding_map"] = "SegmentSecurityProfileBindingMap"
	paramsTypeMap["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["segment_security_profile_binding_map_id"] = "segmentSecurityProfileBindingMapId"
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
		"segment_security_profile_binding_map",
		"PATCH",
		"/global-manager/api/v1/global-infra/segments/{segmentId}/segment-security-profile-binding-maps/{segmentSecurityProfileBindingMapId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func segmentSecurityProfileBindingMapsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	fieldNameMap["segment_security_profile_binding_map"] = "SegmentSecurityProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfileBindingMapsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
}

func segmentSecurityProfileBindingMapsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["segment_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	fields["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	fieldNameMap["segment_id"] = "SegmentId"
	fieldNameMap["segment_security_profile_binding_map_id"] = "SegmentSecurityProfileBindingMapId"
	fieldNameMap["segment_security_profile_binding_map"] = "SegmentSecurityProfileBindingMap"
	paramsTypeMap["segment_security_profile_binding_map"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingMapBindingType)
	paramsTypeMap["segment_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile_binding_map_id"] = bindings.NewStringType()
	paramsTypeMap["segmentId"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileBindingMapId"] = bindings.NewStringType()
	pathParams["segment_id"] = "segmentId"
	pathParams["segment_security_profile_binding_map_id"] = "segmentSecurityProfileBindingMapId"
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
		"segment_security_profile_binding_map",
		"PUT",
		"/global-manager/api/v1/global-infra/segments/{segmentId}/segment-security-profile-binding-maps/{segmentSecurityProfileBindingMapId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
