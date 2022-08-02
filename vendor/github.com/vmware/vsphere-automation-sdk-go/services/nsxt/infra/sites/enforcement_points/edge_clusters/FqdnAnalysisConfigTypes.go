// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: FqdnAnalysisConfig.
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

func fqdnAnalysisConfigDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func fqdnAnalysisConfigDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func fqdnAnalysisConfigDeleteRestMetadata() protocol.OperationRestMetadata {
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
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/fqdn-analysis-config",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fqdnAnalysisConfigGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func fqdnAnalysisConfigGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
}

func fqdnAnalysisConfigGetRestMetadata() protocol.OperationRestMetadata {
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
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/fqdn-analysis-config",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fqdnAnalysisConfigPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["fqdn_analysis_config"] = "FqdnAnalysisConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func fqdnAnalysisConfigPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
}

func fqdnAnalysisConfigPatchRestMetadata() protocol.OperationRestMetadata {
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
	fields["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["fqdn_analysis_config"] = "FqdnAnalysisConfig"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
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
		"fqdn_analysis_config",
		"PATCH",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/fqdn-analysis-config",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fqdnAnalysisConfigUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcement_point_id"] = bindings.NewStringType()
	fields["edge_cluster_id"] = bindings.NewStringType()
	fields["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["fqdn_analysis_config"] = "FqdnAnalysisConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func fqdnAnalysisConfigUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
}

func fqdnAnalysisConfigUpdateRestMetadata() protocol.OperationRestMetadata {
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
	fields["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcement_point_id"] = "EnforcementPointId"
	fieldNameMap["edge_cluster_id"] = "EdgeClusterId"
	fieldNameMap["fqdn_analysis_config"] = "FqdnAnalysisConfig"
	paramsTypeMap["edge_cluster_id"] = bindings.NewStringType()
	paramsTypeMap["fqdn_analysis_config"] = bindings.NewReferenceType(model.FqdnAnalysisConfigBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["enforcement_point_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementPointId"] = bindings.NewStringType()
	paramsTypeMap["edgeClusterId"] = bindings.NewStringType()
	pathParams["site_id"] = "siteId"
	pathParams["edge_cluster_id"] = "edgeClusterId"
	pathParams["enforcement_point_id"] = "enforcementPointId"
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
		"fqdn_analysis_config",
		"PUT",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementPointId}/edge-clusters/{edgeClusterId}/fqdn-analysis-config",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
