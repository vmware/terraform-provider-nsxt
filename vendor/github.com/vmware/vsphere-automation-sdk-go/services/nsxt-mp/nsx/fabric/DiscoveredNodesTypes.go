// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: DiscoveredNodes.
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

// Possible value for ``hasParent`` of method DiscoveredNodes#list.
const DiscoveredNodes_LIST_HAS_PARENT_TRUE = "true"

// Possible value for ``hasParent`` of method DiscoveredNodes#list.
const DiscoveredNodes_LIST_HAS_PARENT_FALSE = "false"

func discoveredNodesCreatetransportnodeInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fields["transport_node"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["node_ext_id"] = "NodeExtId"
	fieldNameMap["transport_node"] = "TransportNode"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func DiscoveredNodesCreatetransportnodeOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeBindingType)
}

func discoveredNodesCreatetransportnodeRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fields["transport_node"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["node_ext_id"] = "NodeExtId"
	fieldNameMap["transport_node"] = "TransportNode"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	paramsTypeMap["transport_node"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeBindingType)
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["node_ext_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["nodeExtId"] = vapiBindings_.NewStringType()
	pathParams["node_ext_id"] = "nodeExtId"
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
		"action=create_transport_node",
		"transport_node",
		"POST",
		"/api/v1/fabric/discovered-nodes/{nodeExtId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func discoveredNodesGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fieldNameMap["node_ext_id"] = "NodeExtId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func DiscoveredNodesGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.DiscoveredNodeBindingType)
}

func discoveredNodesGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fieldNameMap["node_ext_id"] = "NodeExtId"
	paramsTypeMap["node_ext_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["nodeExtId"] = vapiBindings_.NewStringType()
	pathParams["node_ext_id"] = "nodeExtId"
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
		"/api/v1/fabric/discovered-nodes/{nodeExtId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func discoveredNodesListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["has_parent"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["ip_address"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["parent_compute_collection"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cm_local_id"] = "CmLocalId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["display_name"] = "DisplayName"
	fieldNameMap["external_id"] = "ExternalId"
	fieldNameMap["has_parent"] = "HasParent"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["ip_address"] = "IpAddress"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_type"] = "NodeType"
	fieldNameMap["origin_id"] = "OriginId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["parent_compute_collection"] = "ParentComputeCollection"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func DiscoveredNodesListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.DiscoveredNodeListResultBindingType)
}

func discoveredNodesListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["has_parent"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["ip_address"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["parent_compute_collection"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cm_local_id"] = "CmLocalId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["display_name"] = "DisplayName"
	fieldNameMap["external_id"] = "ExternalId"
	fieldNameMap["has_parent"] = "HasParent"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["ip_address"] = "IpAddress"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_type"] = "NodeType"
	fieldNameMap["origin_id"] = "OriginId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["parent_compute_collection"] = "ParentComputeCollection"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["has_parent"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["ip_address"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["parent_compute_collection"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["cm_local_id"] = "cm_local_id"
	queryParams["has_parent"] = "has_parent"
	queryParams["external_id"] = "external_id"
	queryParams["origin_id"] = "origin_id"
	queryParams["ip_address"] = "ip_address"
	queryParams["sort_by"] = "sort_by"
	queryParams["display_name"] = "display_name"
	queryParams["node_type"] = "node_type"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["parent_compute_collection"] = "parent_compute_collection"
	queryParams["node_id"] = "node_id"
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
		"/api/v1/fabric/discovered-nodes",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func discoveredNodesReapplyclusterconfigInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["node_ext_id"] = "NodeExtId"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func DiscoveredNodesReapplyclusterconfigOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeBindingType)
}

func discoveredNodesReapplyclusterconfigRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_ext_id"] = vapiBindings_.NewStringType()
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["node_ext_id"] = "NodeExtId"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["node_ext_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["nodeExtId"] = vapiBindings_.NewStringType()
	pathParams["node_ext_id"] = "nodeExtId"
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
		"action=reapply_cluster_config",
		"",
		"POST",
		"/api/v1/fabric/discovered-nodes/{nodeExtId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
