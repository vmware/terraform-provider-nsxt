// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: GatewayPolicies.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"reflect"
)

// Possible value for ``operation`` of method GatewayPolicies#revise.
const GatewayPolicies_REVISE_OPERATION_TOP = "insert_top"

// Possible value for ``operation`` of method GatewayPolicies#revise.
const GatewayPolicies_REVISE_OPERATION_BOTTOM = "insert_bottom"

// Possible value for ``operation`` of method GatewayPolicies#revise.
const GatewayPolicies_REVISE_OPERATION_AFTER = "insert_after"

// Possible value for ``operation`` of method GatewayPolicies#revise.
const GatewayPolicies_REVISE_OPERATION_BEFORE = "insert_before"

func gatewayPoliciesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func gatewayPoliciesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["gatewayPolicyId"] = bindings.NewStringType()
	pathParams["gateway_policy_id"] = "gatewayPolicyId"
	pathParams["domain_id"] = "domainId"
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
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies/{gatewayPolicyId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func gatewayPoliciesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.GatewayPolicyBindingType)
}

func gatewayPoliciesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["gatewayPolicyId"] = bindings.NewStringType()
	pathParams["gateway_policy_id"] = "gatewayPolicyId"
	pathParams["domain_id"] = "domainId"
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
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies/{gatewayPolicyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func gatewayPoliciesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["include_rule_count"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["include_rule_count"] = "IncludeRuleCount"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.GatewayPolicyListResultBindingType)
}

func gatewayPoliciesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["include_rule_count"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["include_rule_count"] = "IncludeRuleCount"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["include_rule_count"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["domainId"] = bindings.NewStringType()
	pathParams["domain_id"] = "domainId"
	queryParams["cursor"] = "cursor"
	queryParams["include_rule_count"] = "include_rule_count"
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
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func gatewayPoliciesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func gatewayPoliciesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["gatewayPolicyId"] = bindings.NewStringType()
	pathParams["gateway_policy_id"] = "gatewayPolicyId"
	pathParams["domain_id"] = "domainId"
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
		"gateway_policy",
		"PATCH",
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies/{gatewayPolicyId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func gatewayPoliciesReviseInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fields["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	fieldNameMap["anchor_path"] = "AnchorPath"
	fieldNameMap["operation"] = "Operation"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesReviseOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.GatewayPolicyBindingType)
}

func gatewayPoliciesReviseRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fields["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	fieldNameMap["anchor_path"] = "AnchorPath"
	fieldNameMap["operation"] = "Operation"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["gateway_policy_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["gatewayPolicyId"] = bindings.NewStringType()
	pathParams["gateway_policy_id"] = "gatewayPolicyId"
	pathParams["domain_id"] = "domainId"
	queryParams["anchor_path"] = "anchor_path"
	queryParams["operation"] = "operation"
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
		"action=revise",
		"gateway_policy",
		"POST",
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies/{gatewayPolicyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func gatewayPoliciesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func gatewayPoliciesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.GatewayPolicyBindingType)
}

func gatewayPoliciesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["gateway_policy_id"] = bindings.NewStringType()
	fields["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["gateway_policy_id"] = "GatewayPolicyId"
	fieldNameMap["gateway_policy"] = "GatewayPolicy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy_id"] = bindings.NewStringType()
	paramsTypeMap["gateway_policy"] = bindings.NewReferenceType(model.GatewayPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["gatewayPolicyId"] = bindings.NewStringType()
	pathParams["gateway_policy_id"] = "gatewayPolicyId"
	pathParams["domain_id"] = "domainId"
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
		"gateway_policy",
		"PUT",
		"/global-manager/api/v1/global-infra/domains/{domainId}/gateway-policies/{gatewayPolicyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
