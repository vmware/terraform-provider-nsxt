// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: UrlCategorizationConfigs.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package edge_clusters

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func urlCategorizationConfigsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UrlCategorizationConfigsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func urlCategorizationConfigsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	paramsTypeMap["edge_cluster_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["site_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcement_point_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["siteId"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcementPointId"] = vapiBindings_.NewStringType()
	paramsTypeMap["edgeClusterId"] = vapiBindings_.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = vapiBindings_.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func urlCategorizationConfigsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UrlCategorizationConfigsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	paramsTypeMap["edge_cluster_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["site_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcement_point_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["siteId"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcementPointId"] = vapiBindings_.NewStringType()
	paramsTypeMap["edgeClusterId"] = vapiBindings_.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = vapiBindings_.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func urlCategorizationConfigsPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fields["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UrlCategorizationConfigsPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fields["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	paramsTypeMap["edge_cluster_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["site_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcement_point_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	paramsTypeMap["siteId"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcementPointId"] = vapiBindings_.NewStringType()
	paramsTypeMap["edgeClusterId"] = vapiBindings_.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = vapiBindings_.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"policy_url_categorization_config",
		"PATCH",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func urlCategorizationConfigsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fields["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UrlCategorizationConfigsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = vapiBindings_.NewStringType()
	fields["enforcement_point_id"] = vapiBindings_.NewStringType()
	fields["edge_cluster_id"] = vapiBindings_.NewStringType()
	fields["url_categorization_config_id"] = vapiBindings_.NewStringType()
	fields["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	paramsTypeMap["edge_cluster_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["site_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcement_point_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["policy_url_categorization_config"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyUrlCategorizationConfigBindingType)
	paramsTypeMap["siteId"] = vapiBindings_.NewStringType()
	paramsTypeMap["enforcementPointId"] = vapiBindings_.NewStringType()
	paramsTypeMap["edgeClusterId"] = vapiBindings_.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = vapiBindings_.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"policy_url_categorization_config",
		"PUT",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
