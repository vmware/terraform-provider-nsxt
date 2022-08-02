// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: TransportNodeCollections.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func transportNodeCollectionsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func transportNodeCollectionsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
}

func transportNodeCollectionsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsInstallformicrosegInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsInstallformicrosegOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func transportNodeCollectionsInstallformicrosegRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"action=install_for_microseg",
		"host_transport_node_collection",
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["cluster_moid"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["compute_collection_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vc_instance_uuid"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["cluster_moid"] = "ClusterMoid"
	fieldNameMap["compute_collection_id"] = "ComputeCollectionId"
	fieldNameMap["vc_instance_uuid"] = "VcInstanceUuid"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeCollectionListResultBindingType)
}

func transportNodeCollectionsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["cluster_moid"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["compute_collection_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vc_instance_uuid"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["cluster_moid"] = "ClusterMoid"
	fieldNameMap["compute_collection_id"] = "ComputeCollectionId"
	fieldNameMap["vc_instance_uuid"] = "VcInstanceUuid"
	paramsTypeMap["compute_collection_id"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["vc_instance_uuid"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["cluster_moid"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	queryParams["cluster_moid"] = "cluster_moid"
	queryParams["compute_collection_id"] = "compute_collection_id"
	queryParams["vc_instance_uuid"] = "vc_instance_uuid"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func transportNodeCollectionsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"host_transport_node_collection",
		"PATCH",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsRemovensxInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsRemovensxOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func transportNodeCollectionsRemovensxRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"action=remove_nsx",
		"",
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsRetryprofilerealizationInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsRetryprofilerealizationOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func transportNodeCollectionsRetryprofilerealizationRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collection_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	pathParams["site_id"] = "siteId"
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
		"action=retry_profile_realization",
		"",
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collections_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fields["apply_profile"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collections_id"] = "TransportNodeCollectionsId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	fieldNameMap["apply_profile"] = "ApplyProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func transportNodeCollectionsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
}

func transportNodeCollectionsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["transport_node_collections_id"] = bindings.NewStringType()
	fields["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	fields["apply_profile"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["transport_node_collections_id"] = "TransportNodeCollectionsId"
	fieldNameMap["host_transport_node_collection"] = "HostTransportNodeCollection"
	fieldNameMap["apply_profile"] = "ApplyProfile"
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["apply_profile"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["host_transport_node_collection"] = bindings.NewReferenceType(model.HostTransportNodeCollectionBindingType)
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["transport_node_collections_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["transportNodeCollectionsId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["transport_node_collections_id"] = "transportNodeCollectionsId"
	queryParams["apply_profile"] = "apply_profile"
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
		"host_transport_node_collection",
		"PUT",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/transport-node-collections/{transportNodeCollectionsId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
