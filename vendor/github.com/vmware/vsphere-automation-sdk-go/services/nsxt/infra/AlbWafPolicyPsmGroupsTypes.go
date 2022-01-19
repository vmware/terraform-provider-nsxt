// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: AlbWafPolicyPsmGroups.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func albWafPolicyPsmGroupsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["force"] = "Force"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albWafPolicyPsmGroupsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func albWafPolicyPsmGroupsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["force"] = "Force"
	paramsTypeMap["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	paramsTypeMap["albWafpolicypsmgroupId"] = bindings.NewStringType()
	pathParams["alb_wafpolicypsmgroup_id"] = "albWafpolicypsmgroupId"
	queryParams["force"] = "force"
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
		"/policy/api/v1/infra/alb-waf-policy-psm-groups/{albWafpolicypsmgroupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albWafPolicyPsmGroupsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albWafPolicyPsmGroupsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
}

func albWafPolicyPsmGroupsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	paramsTypeMap["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	paramsTypeMap["albWafpolicypsmgroupId"] = bindings.NewStringType()
	pathParams["alb_wafpolicypsmgroup_id"] = "albWafpolicypsmgroupId"
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
		"/policy/api/v1/infra/alb-waf-policy-psm-groups/{albWafpolicypsmgroupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albWafPolicyPsmGroupsListInputType() bindings.StructType {
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

func albWafPolicyPsmGroupsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBWafPolicyPSMGroupApiResponseBindingType)
}

func albWafPolicyPsmGroupsListRestMetadata() protocol.OperationRestMetadata {
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
		"/policy/api/v1/infra/alb-waf-policy-psm-groups",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albWafPolicyPsmGroupsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["a_LB_waf_policy_PSM_group"] = "ALBWafPolicyPSMGroup"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albWafPolicyPsmGroupsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func albWafPolicyPsmGroupsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["a_LB_waf_policy_PSM_group"] = "ALBWafPolicyPSMGroup"
	paramsTypeMap["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	paramsTypeMap["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	paramsTypeMap["albWafpolicypsmgroupId"] = bindings.NewStringType()
	pathParams["alb_wafpolicypsmgroup_id"] = "albWafpolicypsmgroupId"
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
		"a_LB_waf_policy_PSM_group",
		"PATCH",
		"/policy/api/v1/infra/alb-waf-policy-psm-groups/{albWafpolicypsmgroupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albWafPolicyPsmGroupsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["a_LB_waf_policy_PSM_group"] = "ALBWafPolicyPSMGroup"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albWafPolicyPsmGroupsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
}

func albWafPolicyPsmGroupsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	fields["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	fieldNameMap["alb_wafpolicypsmgroup_id"] = "AlbWafpolicypsmgroupId"
	fieldNameMap["a_LB_waf_policy_PSM_group"] = "ALBWafPolicyPSMGroup"
	paramsTypeMap["alb_wafpolicypsmgroup_id"] = bindings.NewStringType()
	paramsTypeMap["a_LB_waf_policy_PSM_group"] = bindings.NewReferenceType(model.ALBWafPolicyPSMGroupBindingType)
	paramsTypeMap["albWafpolicypsmgroupId"] = bindings.NewStringType()
	pathParams["alb_wafpolicypsmgroup_id"] = "albWafpolicypsmgroupId"
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
		"a_LB_waf_policy_PSM_group",
		"PUT",
		"/policy/api/v1/infra/alb-waf-policy-psm-groups/{albWafpolicypsmgroupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
