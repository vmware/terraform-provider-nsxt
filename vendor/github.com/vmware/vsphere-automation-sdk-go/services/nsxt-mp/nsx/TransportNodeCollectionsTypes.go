// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: TransportNodeCollections.
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

func transportNodeCollectionsCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	fields["apply_profile"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_collection"] = "TransportNodeCollection"
	fieldNameMap["apply_profile"] = "ApplyProfile"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
}

func transportNodeCollectionsCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	fields["apply_profile"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_collection"] = "TransportNodeCollection"
	fieldNameMap["apply_profile"] = "ApplyProfile"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	paramsTypeMap["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	paramsTypeMap["apply_profile"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	queryParams["apply_profile"] = "apply_profile"
	queryParams["override_nsx_ownership"] = "override_nsx_ownership"
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
		"transport_node_collection",
		"POST",
		"/api/v1/transport-node-collections",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func transportNodeCollectionsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
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
		"/api/v1/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
}

func transportNodeCollectionsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
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
		"/api/v1/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cluster_moid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["compute_collection_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vc_instance_uuid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cluster_moid"] = "ClusterMoid"
	fieldNameMap["compute_collection_id"] = "ComputeCollectionId"
	fieldNameMap["vc_instance_uuid"] = "VcInstanceUuid"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionListResultBindingType)
}

func transportNodeCollectionsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cluster_moid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["compute_collection_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vc_instance_uuid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cluster_moid"] = "ClusterMoid"
	fieldNameMap["compute_collection_id"] = "ComputeCollectionId"
	fieldNameMap["vc_instance_uuid"] = "VcInstanceUuid"
	paramsTypeMap["cluster_moid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["compute_collection_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vc_instance_uuid"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["cluster_moid"] = "cluster_moid"
	queryParams["compute_collection_id"] = "compute_collection_id"
	queryParams["vc_instance_uuid"] = "vc_instance_uuid"
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
		"/api/v1/transport-node-collections",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsRetryprofilerealizationInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsRetryprofilerealizationOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func transportNodeCollectionsRetryprofilerealizationRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	paramsTypeMap["transport_node_collection_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeCollectionId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
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
		"action=retry_profile_realization",
		"",
		"POST",
		"/api/v1/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeCollectionsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fields["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["transport_node_collection"] = "TransportNodeCollection"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeCollectionsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
}

func transportNodeCollectionsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_collection_id"] = vapiBindings_.NewStringType()
	fields["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_collection_id"] = "TransportNodeCollectionId"
	fieldNameMap["transport_node_collection"] = "TransportNodeCollection"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	paramsTypeMap["transport_node_collection_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transport_node_collection"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeCollectionBindingType)
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["transportNodeCollectionId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_collection_id"] = "transportNodeCollectionId"
	queryParams["override_nsx_ownership"] = "override_nsx_ownership"
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
		"transport_node_collection",
		"PUT",
		"/api/v1/transport-node-collections/{transportNodeCollectionId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
