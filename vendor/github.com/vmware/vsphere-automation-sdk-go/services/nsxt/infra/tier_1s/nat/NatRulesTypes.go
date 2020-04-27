/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: NatRules.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package nat

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func natRulesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func natRulesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func natRulesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["nat_rule_id"] = bindings.NewStringType()
	paramsTypeMap["nat_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["natId"] = bindings.NewStringType()
	paramsTypeMap["natRuleId"] = bindings.NewStringType()
	pathParams["nat_rule_id"] = "natRuleId"
	pathParams["nat_id"] = "natId"
	pathParams["tier1_id"] = "tier1Id"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/nat/{natId}/nat-rules/{natRuleId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func natRulesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func natRulesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyNatRuleBindingType)
}

func natRulesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["nat_rule_id"] = bindings.NewStringType()
	paramsTypeMap["nat_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["natId"] = bindings.NewStringType()
	paramsTypeMap["natRuleId"] = bindings.NewStringType()
	pathParams["nat_rule_id"] = "natRuleId"
	pathParams["nat_id"] = "natId"
	pathParams["tier1_id"] = "tier1Id"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/nat/{natId}/nat-rules/{natRuleId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func natRulesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func natRulesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyNatRuleListResultBindingType)
}

func natRulesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["nat_id"] = bindings.NewStringType()
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["natId"] = bindings.NewStringType()
	pathParams["nat_id"] = "natId"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/nat/{natId}/nat-rules",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func natRulesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fields["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	fieldNameMap["policy_nat_rule"] = "PolicyNatRule"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func natRulesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func natRulesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fields["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	fieldNameMap["policy_nat_rule"] = "PolicyNatRule"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	paramsTypeMap["nat_rule_id"] = bindings.NewStringType()
	paramsTypeMap["nat_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["natId"] = bindings.NewStringType()
	paramsTypeMap["natRuleId"] = bindings.NewStringType()
	pathParams["nat_rule_id"] = "natRuleId"
	pathParams["nat_id"] = "natId"
	pathParams["tier1_id"] = "tier1Id"
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
		"policy_nat_rule",
		"PATCH",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/nat/{natId}/nat-rules/{natRuleId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func natRulesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fields["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	fieldNameMap["policy_nat_rule"] = "PolicyNatRule"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func natRulesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyNatRuleBindingType)
}

func natRulesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["nat_id"] = bindings.NewStringType()
	fields["nat_rule_id"] = bindings.NewStringType()
	fields["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["nat_id"] = "NatId"
	fieldNameMap["nat_rule_id"] = "NatRuleId"
	fieldNameMap["policy_nat_rule"] = "PolicyNatRule"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["policy_nat_rule"] = bindings.NewReferenceType(model.PolicyNatRuleBindingType)
	paramsTypeMap["nat_rule_id"] = bindings.NewStringType()
	paramsTypeMap["nat_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["natId"] = bindings.NewStringType()
	paramsTypeMap["natRuleId"] = bindings.NewStringType()
	pathParams["nat_rule_id"] = "natRuleId"
	pathParams["nat_id"] = "natId"
	pathParams["tier1_id"] = "tier1Id"
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
		"policy_nat_rule",
		"PUT",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/nat/{natId}/nat-rules/{natRuleId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


