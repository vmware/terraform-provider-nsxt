// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: ContainerClusters.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package fabric

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"reflect"
)

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_PAS = "PAS"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_PKS = "PKS"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_KUBERNETES = "Kubernetes"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_OPENSHIFT = "Openshift"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_WCP = "WCP"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_WCP_GUEST = "WCP_Guest"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_AKS = "AKS"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_EKS = "EKS"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_TKGM = "TKGm"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_TKGI = "TKGi"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_GKE = "GKE"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_GARDENER = "Gardener"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_RANCHER = "Rancher"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_TAS = "TAS"

// Possible value for ``clusterType`` of method ContainerClusters#list.
const ContainerClusters_LIST_CLUSTER_TYPE_OTHER = "Other"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_VSPHERE = "vSphere"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_AWS = "AWS"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_AZURE = "Azure"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_GOOGLE = "Google"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_VMC = "VMC"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_KVM = "KVM"

// Possible value for ``infraType`` of method ContainerClusters#list.
const ContainerClusters_LIST_INFRA_TYPE_BAREMETAL = "Baremetal"

func containerClustersGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["container_cluster_id"] = vapiBindings_.NewStringType()
	fieldNameMap["container_cluster_id"] = "ContainerClusterId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ContainerClustersGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ContainerClusterBindingType)
}

func containerClustersGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["container_cluster_id"] = vapiBindings_.NewStringType()
	fieldNameMap["container_cluster_id"] = "ContainerClusterId"
	paramsTypeMap["container_cluster_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["containerClusterId"] = vapiBindings_.NewStringType()
	pathParams["container_cluster_id"] = "containerClusterId"
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
		"/api/v1/fabric/container-clusters/{containerClusterId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func containerClustersListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cluster_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["infra_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["scope_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cluster_type"] = "ClusterType"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["infra_type"] = "InfraType"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["scope_id"] = "ScopeId"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ContainerClustersListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ContainerClusterListResultBindingType)
}

func containerClustersListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cluster_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["infra_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["scope_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cluster_type"] = "ClusterType"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["infra_type"] = "InfraType"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["scope_id"] = "ScopeId"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["infra_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["scope_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["cluster_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["infra_type"] = "infra_type"
	queryParams["scope_id"] = "scope_id"
	queryParams["cluster_type"] = "cluster_type"
	queryParams["sort_by"] = "sort_by"
	queryParams["page_size"] = "page_size"
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
		"/api/v1/fabric/container-clusters",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
