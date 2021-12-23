// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: IntrusionServicePolicies.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

// Possible value for ``operation`` of method IntrusionServicePolicies#revise.
const IntrusionServicePolicies_REVISE_OPERATION_TOP = "insert_top"

// Possible value for ``operation`` of method IntrusionServicePolicies#revise.
const IntrusionServicePolicies_REVISE_OPERATION_BOTTOM = "insert_bottom"

// Possible value for ``operation`` of method IntrusionServicePolicies#revise.
const IntrusionServicePolicies_REVISE_OPERATION_AFTER = "insert_after"

// Possible value for ``operation`` of method IntrusionServicePolicies#revise.
const IntrusionServicePolicies_REVISE_OPERATION_BEFORE = "insert_before"

func intrusionServicePoliciesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func intrusionServicePoliciesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func intrusionServicePoliciesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["policy_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["policyId"] = bindings.NewStringType()
	pathParams["policy_id"] = "policyId"
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
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies/{policyId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func intrusionServicePoliciesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func intrusionServicePoliciesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
}

func intrusionServicePoliciesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["policy_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["policyId"] = bindings.NewStringType()
	pathParams["policy_id"] = "policyId"
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
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies/{policyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func intrusionServicePoliciesListInputType() bindings.StructType {
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

func intrusionServicePoliciesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IdsSecurityPolicyListResultBindingType)
}

func intrusionServicePoliciesListRestMetadata() protocol.OperationRestMetadata {
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
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func intrusionServicePoliciesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func intrusionServicePoliciesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func intrusionServicePoliciesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["policy_id"] = bindings.NewStringType()
	paramsTypeMap["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["policyId"] = bindings.NewStringType()
	pathParams["policy_id"] = "policyId"
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
		"ids_security_policy",
		"PATCH",
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies/{policyId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func intrusionServicePoliciesReviseInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fields["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	fieldNameMap["anchor_path"] = "AnchorPath"
	fieldNameMap["operation"] = "Operation"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func intrusionServicePoliciesReviseOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
}

func intrusionServicePoliciesReviseRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fields["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	fieldNameMap["anchor_path"] = "AnchorPath"
	fieldNameMap["operation"] = "Operation"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["policy_id"] = bindings.NewStringType()
	paramsTypeMap["anchor_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["operation"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["policyId"] = bindings.NewStringType()
	pathParams["policy_id"] = "policyId"
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
		"ids_security_policy",
		"POST",
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies/{policyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func intrusionServicePoliciesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func intrusionServicePoliciesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
}

func intrusionServicePoliciesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["policy_id"] = bindings.NewStringType()
	fields["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["policy_id"] = "PolicyId"
	fieldNameMap["ids_security_policy"] = "IdsSecurityPolicy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["policy_id"] = bindings.NewStringType()
	paramsTypeMap["ids_security_policy"] = bindings.NewReferenceType(model.IdsSecurityPolicyBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["policyId"] = bindings.NewStringType()
	pathParams["policy_id"] = "policyId"
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
		"ids_security_policy",
		"PUT",
		"/policy/api/v1/infra/domains/{domainId}/intrusion-service-policies/{policyId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
