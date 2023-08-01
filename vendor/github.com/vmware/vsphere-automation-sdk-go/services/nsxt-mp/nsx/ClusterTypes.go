// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Cluster.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package nsx

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"reflect"
)

// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_GLOBAL_MANAGER = "GLOBAL_MANAGER"

// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_LOCAL_MANAGER = "LOCAL_MANAGER"

// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_LOCAL_LOCAL_MANAGER = "LOCAL_LOCAL_MANAGER"

// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_NSX_INTELLIGENCE = "NSX_INTELLIGENCE"

// Possible value for ``force`` of method Cluster#removenode.
const Cluster_REMOVENODE_FORCE_TRUE = "true"

// Possible value for ``force`` of method Cluster#removenode.
const Cluster_REMOVENODE_FORCE_FALSE = "false"

// Possible value for ``gracefulShutdown`` of method Cluster#removenode.
const Cluster_REMOVENODE_GRACEFUL_SHUTDOWN_TRUE = "true"

// Possible value for ``gracefulShutdown`` of method Cluster#removenode.
const Cluster_REMOVENODE_GRACEFUL_SHUTDOWN_FALSE = "false"

// Possible value for ``ignoreRepositoryIpCheck`` of method Cluster#removenode.
const Cluster_REMOVENODE_IGNORE_REPOSITORY_IP_CHECK_TRUE = "true"

// Possible value for ``ignoreRepositoryIpCheck`` of method Cluster#removenode.
const Cluster_REMOVENODE_IGNORE_REPOSITORY_IP_CHECK_FALSE = "false"

func clusterBackuptoremoteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["frame_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["site_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["frame_type"] = "FrameType"
	fieldNameMap["site_id"] = "SiteId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterBackuptoremoteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterBackuptoremoteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["frame_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["site_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["frame_type"] = "FrameType"
	fieldNameMap["site_id"] = "SiteId"
	paramsTypeMap["site_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["frame_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["site_id"] = "site_id"
	queryParams["frame_type"] = "frame_type"
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
		"action=backup_to_remote",
		"",
		"POST",
		"/api/v1/cluster",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	paramsTypeMap["target_uri"] = vapiBindings_.NewStringType()
	paramsTypeMap["target_node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetNodeId"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetUri"] = vapiBindings_.NewStringType()
	pathParams["target_uri"] = "targetUri"
	pathParams["target_node_id"] = "targetNodeId"
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
		"POST",
		"/api/v1/cluster/{targetNodeId}/{targetUri}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	paramsTypeMap["target_uri"] = vapiBindings_.NewStringType()
	paramsTypeMap["target_node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetNodeId"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetUri"] = vapiBindings_.NewStringType()
	pathParams["target_uri"] = "targetUri"
	pathParams["target_node_id"] = "targetNodeId"
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
		"/api/v1/cluster/{targetNodeId}/{targetUri}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ClusterConfigBindingType)
}

func clusterGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
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
		"/api/v1/cluster",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterGet0InputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_id"] = vapiBindings_.NewStringType()
	fieldNameMap["node_id"] = "NodeId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterGet0OutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ClusterNodeInfoBindingType)
}

func clusterGet0RestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_id"] = vapiBindings_.NewStringType()
	fieldNameMap["node_id"] = "NodeId"
	paramsTypeMap["node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["nodeId"] = vapiBindings_.NewStringType()
	pathParams["node_id"] = "nodeId"
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
		"/api/v1/cluster/{nodeId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterGet1InputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterGet1OutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterGet1RestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	paramsTypeMap["target_uri"] = vapiBindings_.NewStringType()
	paramsTypeMap["target_node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetNodeId"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetUri"] = vapiBindings_.NewStringType()
	pathParams["target_uri"] = "targetUri"
	pathParams["target_node_id"] = "targetNodeId"
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
		"/api/v1/cluster/{targetNodeId}/{targetUri}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterJoinclusterInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["join_cluster_parameters"] = vapiBindings_.NewReferenceType(nsxModel.JoinClusterParametersBindingType)
	fieldNameMap["join_cluster_parameters"] = "JoinClusterParameters"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterJoinclusterOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ClusterConfigurationBindingType)
}

func clusterJoinclusterRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["join_cluster_parameters"] = vapiBindings_.NewReferenceType(nsxModel.JoinClusterParametersBindingType)
	fieldNameMap["join_cluster_parameters"] = "JoinClusterParameters"
	paramsTypeMap["join_cluster_parameters"] = vapiBindings_.NewReferenceType(nsxModel.JoinClusterParametersBindingType)
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
		"action=join_cluster",
		"join_cluster_parameters",
		"POST",
		"/api/v1/cluster",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterRemovenodeInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_id"] = vapiBindings_.NewStringType()
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["graceful_shutdown"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["ignore_repository_ip_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["force"] = "Force"
	fieldNameMap["graceful_shutdown"] = "GracefulShutdown"
	fieldNameMap["ignore_repository_ip_check"] = "IgnoreRepositoryIpCheck"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterRemovenodeOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ClusterConfigurationBindingType)
}

func clusterRemovenodeRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_id"] = vapiBindings_.NewStringType()
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["graceful_shutdown"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["ignore_repository_ip_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["force"] = "Force"
	fieldNameMap["graceful_shutdown"] = "GracefulShutdown"
	fieldNameMap["ignore_repository_ip_check"] = "IgnoreRepositoryIpCheck"
	paramsTypeMap["ignore_repository_ip_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["graceful_shutdown"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["nodeId"] = vapiBindings_.NewStringType()
	pathParams["node_id"] = "nodeId"
	queryParams["graceful_shutdown"] = "graceful-shutdown"
	queryParams["ignore_repository_ip_check"] = "ignore-repository-ip-check"
	queryParams["force"] = "force"
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
		"action=remove_node",
		"",
		"POST",
		"/api/v1/cluster/{nodeId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterSummarizeinventorytoremoteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterSummarizeinventorytoremoteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterSummarizeinventorytoremoteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
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
		"action=summarize_inventory_to_remote",
		"",
		"POST",
		"/api/v1/cluster",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func clusterUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ClusterUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func clusterUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["target_node_id"] = vapiBindings_.NewStringType()
	fields["target_uri"] = vapiBindings_.NewStringType()
	fieldNameMap["target_node_id"] = "TargetNodeId"
	fieldNameMap["target_uri"] = "TargetUri"
	paramsTypeMap["target_uri"] = vapiBindings_.NewStringType()
	paramsTypeMap["target_node_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetNodeId"] = vapiBindings_.NewStringType()
	paramsTypeMap["targetUri"] = vapiBindings_.NewStringType()
	pathParams["target_uri"] = "targetUri"
	pathParams["target_node_id"] = "targetNodeId"
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
		"PUT",
		"/api/v1/cluster/{targetNodeId}/{targetUri}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
