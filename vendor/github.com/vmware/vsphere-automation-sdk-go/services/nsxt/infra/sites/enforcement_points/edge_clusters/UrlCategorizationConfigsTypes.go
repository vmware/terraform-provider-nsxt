// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: UrlCategorizationConfigs.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package edge_clusters

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func urlCategorizationConfigsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func urlCategorizationConfigsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func urlCategorizationConfigsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func urlCategorizationConfigsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func urlCategorizationConfigsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/url-categorization-configs/{urlCategorizationConfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func urlCategorizationConfigsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fields["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func urlCategorizationConfigsPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fields["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = bindings.NewStringType()
	paramsTypeMap["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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

func urlCategorizationConfigsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fields["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func urlCategorizationConfigsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
}

func urlCategorizationConfigsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["url_categorization_config_id"] = bindings.NewStringType()
	fields["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["url_categorization_config_id"] = "UrlCategorizationConfigId"
	fieldNameMap["policy_url_categorization_config"] = "PolicyUrlCategorizationConfig"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["url_categorization_config_id"] = bindings.NewStringType()
	paramsTypeMap["policy_url_categorization_config"] = bindings.NewReferenceType(model.PolicyUrlCategorizationConfigBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	paramsTypeMap["urlCategorizationConfigId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
	pathParams["url_categorization_config_id"] = "urlCategorizationConfigId"
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
