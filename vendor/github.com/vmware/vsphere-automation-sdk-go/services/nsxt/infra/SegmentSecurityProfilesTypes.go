/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: SegmentSecurityProfiles.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func segmentSecurityProfilesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfilesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentSecurityProfilesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	paramsTypeMap["segment_security_profile_id"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileId"] = bindings.NewStringType()
	pathParams["segment_security_profile_id"] = "segmentSecurityProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"DELETE",
		"/policy/api/v1/infra/segment-security-profiles/{segmentSecurityProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func segmentSecurityProfilesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfilesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
}

func segmentSecurityProfilesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	paramsTypeMap["segment_security_profile_id"] = bindings.NewStringType()
	paramsTypeMap["segmentSecurityProfileId"] = bindings.NewStringType()
	pathParams["segment_security_profile_id"] = "segmentSecurityProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/segment-security-profiles/{segmentSecurityProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func segmentSecurityProfilesListInputType() bindings.StructType {
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

func segmentSecurityProfilesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileListResultBindingType)
}

func segmentSecurityProfilesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
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
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/segment-security-profiles",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func segmentSecurityProfilesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fields["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	fieldNameMap["segment_security_profile"] = "SegmentSecurityProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfilesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func segmentSecurityProfilesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fields["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	fieldNameMap["segment_security_profile"] = "SegmentSecurityProfile"
	paramsTypeMap["segment_security_profile_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	paramsTypeMap["segmentSecurityProfileId"] = bindings.NewStringType()
	pathParams["segment_security_profile_id"] = "segmentSecurityProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"segment_security_profile",
		"PATCH",
		"/policy/api/v1/infra/segment-security-profiles/{segmentSecurityProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func segmentSecurityProfilesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fields["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	fieldNameMap["segment_security_profile"] = "SegmentSecurityProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func segmentSecurityProfilesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
}

func segmentSecurityProfilesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["segment_security_profile_id"] = bindings.NewStringType()
	fields["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	fieldNameMap["segment_security_profile_id"] = "SegmentSecurityProfileId"
	fieldNameMap["segment_security_profile"] = "SegmentSecurityProfile"
	paramsTypeMap["segment_security_profile_id"] = bindings.NewStringType()
	paramsTypeMap["segment_security_profile"] = bindings.NewReferenceType(model.SegmentSecurityProfileBindingType)
	paramsTypeMap["segmentSecurityProfileId"] = bindings.NewStringType()
	pathParams["segment_security_profile_id"] = "segmentSecurityProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"segment_security_profile",
		"PUT",
		"/policy/api/v1/infra/segment-security-profiles/{segmentSecurityProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


