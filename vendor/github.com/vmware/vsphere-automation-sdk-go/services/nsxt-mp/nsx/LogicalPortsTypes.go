// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: LogicalPorts.
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

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_VIF = "VIF"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_LOGICALROUTER = "LOGICALROUTER"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_BRIDGEENDPOINT = "BRIDGEENDPOINT"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_DHCP_SERVICE = "DHCP_SERVICE"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_METADATA_PROXY = "METADATA_PROXY"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_L2VPN_SESSION = "L2VPN_SESSION"

// Possible value for ``attachmentType`` of method LogicalPorts#list.
const LogicalPorts_LIST_ATTACHMENT_TYPE_NONE = "NONE"

func logicalPortsCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
	fieldNameMap["logical_port"] = "LogicalPort"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalPortsCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
}

func logicalPortsCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
	fieldNameMap["logical_port"] = "LogicalPort"
	paramsTypeMap["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
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
		"logical_port",
		"POST",
		"/api/v1/logical-ports",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalPortsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lport_id"] = vapiBindings_.NewStringType()
	fields["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["lport_id"] = "LportId"
	fieldNameMap["detach"] = "Detach"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalPortsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func logicalPortsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lport_id"] = vapiBindings_.NewStringType()
	fields["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["lport_id"] = "LportId"
	fieldNameMap["detach"] = "Detach"
	paramsTypeMap["lport_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["lportId"] = vapiBindings_.NewStringType()
	pathParams["lport_id"] = "lportId"
	queryParams["detach"] = "detach"
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
		"/api/v1/logical-ports/{lportId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalPortsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lport_id"] = vapiBindings_.NewStringType()
	fieldNameMap["lport_id"] = "LportId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalPortsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
}

func logicalPortsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lport_id"] = vapiBindings_.NewStringType()
	fieldNameMap["lport_id"] = "LportId"
	paramsTypeMap["lport_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["lportId"] = vapiBindings_.NewStringType()
	pathParams["lport_id"] = "lportId"
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
		"/api/v1/logical-ports/{lportId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalPortsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["attachment_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["attachment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["container_ports_only"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["parent_vif_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["attachment_id"] = "AttachmentId"
	fieldNameMap["attachment_type"] = "AttachmentType"
	fieldNameMap["bridge_cluster_id"] = "BridgeClusterId"
	fieldNameMap["container_ports_only"] = "ContainerPortsOnly"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["diagnostic"] = "Diagnostic"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["parent_vif_id"] = "ParentVifId"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["switching_profile_id"] = "SwitchingProfileId"
	fieldNameMap["transport_node_id"] = "TransportNodeId"
	fieldNameMap["transport_zone_id"] = "TransportZoneId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalPortsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalPortListResultBindingType)
}

func logicalPortsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["attachment_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["attachment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["container_ports_only"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["parent_vif_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["attachment_id"] = "AttachmentId"
	fieldNameMap["attachment_type"] = "AttachmentType"
	fieldNameMap["bridge_cluster_id"] = "BridgeClusterId"
	fieldNameMap["container_ports_only"] = "ContainerPortsOnly"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["diagnostic"] = "Diagnostic"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["parent_vif_id"] = "ParentVifId"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["switching_profile_id"] = "SwitchingProfileId"
	fieldNameMap["transport_node_id"] = "TransportNodeId"
	fieldNameMap["transport_zone_id"] = "TransportZoneId"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["transport_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["attachment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["attachment_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["container_ports_only"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["parent_vif_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["switching_profile_id"] = "switching_profile_id"
	queryParams["sort_by"] = "sort_by"
	queryParams["transport_node_id"] = "transport_node_id"
	queryParams["bridge_cluster_id"] = "bridge_cluster_id"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["attachment_type"] = "attachment_type"
	queryParams["included_fields"] = "included_fields"
	queryParams["attachment_id"] = "attachment_id"
	queryParams["transport_zone_id"] = "transport_zone_id"
	queryParams["container_ports_only"] = "container_ports_only"
	queryParams["diagnostic"] = "diagnostic"
	queryParams["parent_vif_id"] = "parent_vif_id"
	queryParams["logical_switch_id"] = "logical_switch_id"
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
		"/api/v1/logical-ports",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalPortsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lport_id"] = vapiBindings_.NewStringType()
	fields["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
	fieldNameMap["lport_id"] = "LportId"
	fieldNameMap["logical_port"] = "LogicalPort"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalPortsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
}

func logicalPortsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lport_id"] = vapiBindings_.NewStringType()
	fields["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
	fieldNameMap["lport_id"] = "LportId"
	fieldNameMap["logical_port"] = "LogicalPort"
	paramsTypeMap["lport_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logical_port"] = vapiBindings_.NewReferenceType(nsxModel.LogicalPortBindingType)
	paramsTypeMap["lportId"] = vapiBindings_.NewStringType()
	pathParams["lport_id"] = "lportId"
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
		"logical_port",
		"PUT",
		"/api/v1/logical-ports/{lportId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
