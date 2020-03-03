/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: IpfixL2Profiles.
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





func ipfixL2ProfilesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["override"] = "Override"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixL2ProfilesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixL2ProfilesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["override"] = "Override"
	paramsTypeMap["ipfix_l2_profile_id"] = bindings.NewStringType()
	paramsTypeMap["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["ipfixL2ProfileId"] = bindings.NewStringType()
	pathParams["ipfix_l2_profile_id"] = "ipfixL2ProfileId"
	queryParams["override"] = "override"
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
		"/policy/api/v1/infra/ipfix-l2-profiles/{ipfixL2ProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipfixL2ProfilesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixL2ProfilesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
}

func ipfixL2ProfilesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	paramsTypeMap["ipfix_l2_profile_id"] = bindings.NewStringType()
	paramsTypeMap["ipfixL2ProfileId"] = bindings.NewStringType()
	pathParams["ipfix_l2_profile_id"] = "ipfixL2ProfileId"
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
		"/policy/api/v1/infra/ipfix-l2-profiles/{ipfixL2ProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipfixL2ProfilesListInputType() bindings.StructType {
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

func ipfixL2ProfilesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXL2ProfileListResultBindingType)
}

func ipfixL2ProfilesListRestMetadata() protocol.OperationRestMetadata {
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
		"/policy/api/v1/infra/ipfix-l2-profiles",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipfixL2ProfilesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["i_PFIX_l2_profile"] = "IPFIXL2Profile"
	fieldNameMap["override"] = "Override"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixL2ProfilesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipfixL2ProfilesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["i_PFIX_l2_profile"] = "IPFIXL2Profile"
	fieldNameMap["override"] = "Override"
	paramsTypeMap["ipfix_l2_profile_id"] = bindings.NewStringType()
	paramsTypeMap["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	paramsTypeMap["ipfixL2ProfileId"] = bindings.NewStringType()
	pathParams["ipfix_l2_profile_id"] = "ipfixL2ProfileId"
	queryParams["override"] = "override"
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
		"i_PFIX_l2_profile",
		"PATCH",
		"/policy/api/v1/infra/ipfix-l2-profiles/{ipfixL2ProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipfixL2ProfilesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["i_PFIX_l2_profile"] = "IPFIXL2Profile"
	fieldNameMap["override"] = "Override"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipfixL2ProfilesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
}

func ipfixL2ProfilesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["ipfix_l2_profile_id"] = bindings.NewStringType()
	fields["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["ipfix_l2_profile_id"] = "IpfixL2ProfileId"
	fieldNameMap["i_PFIX_l2_profile"] = "IPFIXL2Profile"
	fieldNameMap["override"] = "Override"
	paramsTypeMap["ipfix_l2_profile_id"] = bindings.NewStringType()
	paramsTypeMap["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["i_PFIX_l2_profile"] = bindings.NewReferenceType(model.IPFIXL2ProfileBindingType)
	paramsTypeMap["ipfixL2ProfileId"] = bindings.NewStringType()
	pathParams["ipfix_l2_profile_id"] = "ipfixL2ProfileId"
	queryParams["override"] = "override"
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
		"i_PFIX_l2_profile",
		"PUT",
		"/policy/api/v1/infra/ipfix-l2-profiles/{ipfixL2ProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


