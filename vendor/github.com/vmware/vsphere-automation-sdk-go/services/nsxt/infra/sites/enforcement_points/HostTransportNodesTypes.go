// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: HostTransportNodes.
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

// Possible value for ``action`` of method HostTransportNodes#updatemaintenancemode.
const HostTransportNodes_UPDATEMAINTENANCEMODE_ACTION_ENTER_MAINTENANCE_MODE = "enter_maintenance_mode"

// Possible value for ``action`` of method HostTransportNodes#updatemaintenancemode.
const HostTransportNodes_UPDATEMAINTENANCEMODE_ACTION_FORCED_ENTER_MAINTENANCE_MODE = "forced_enter_maintenance_mode"

// Possible value for ``action`` of method HostTransportNodes#updatemaintenancemode.
const HostTransportNodes_UPDATEMAINTENANCEMODE_ACTION_EXIT_MAINTENANCE_MODE = "exit_maintenance_mode"

func hostTransportNodesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["unprepare_host"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["force"] = "Force"
	fieldNameMap["unprepare_host"] = "UnprepareHost"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesDeleteRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["unprepare_host"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["force"] = "Force"
	fieldNameMap["unprepare_host"] = "UnprepareHost"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["unprepare_host"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
	queryParams["unprepare_host"] = "unprepare_host"
	queryParams["force"] = "force"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeBindingType)
}

func hostTransportNodesGetRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["discovered_node_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["in_maintenance_mode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["node_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["node_types"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["transport_zone_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["discovered_node_id"] = "DiscoveredNodeId"
	fieldNameMap["in_maintenance_mode"] = "InMaintenanceMode"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["node_ip"] = "NodeIp"
	fieldNameMap["node_types"] = "NodeTypes"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["transport_zone_path"] = "TransportZonePath"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeListResultBindingType)
}

func hostTransportNodesListRestMetadata() protocol.OperationRestMetadata {
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
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["discovered_node_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["in_maintenance_mode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["node_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["node_types"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["transport_zone_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["discovered_node_id"] = "DiscoveredNodeId"
	fieldNameMap["in_maintenance_mode"] = "InMaintenanceMode"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["node_ip"] = "NodeIp"
	fieldNameMap["node_types"] = "NodeTypes"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["transport_zone_path"] = "TransportZonePath"
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["discovered_node_id"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["transport_zone_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["in_maintenance_mode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["node_types"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["node_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	queryParams["cursor"] = "cursor"
	queryParams["node_types"] = "node_types"
	queryParams["discovered_node_id"] = "discovered_node_id"
	queryParams["node_ip"] = "node_ip"
	queryParams["transport_zone_path"] = "transport_zone_path"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["in_maintenance_mode"] = "in_maintenance_mode"
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
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesMigratetovdsInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["skip_maintmode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["skip_maintmode"] = "SkipMaintmode"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesMigratetovdsOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesMigratetovdsRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["skip_maintmode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["skip_maintmode"] = "SkipMaintmode"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["skip_maintmode"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
	queryParams["skip_maintmode"] = "skip_maintmode"
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
		"action=migrate_to_vds",
		"",
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	fields["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["host_transport_node"] = "HostTransportNode"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesPatchRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	fields["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["host_transport_node"] = "HostTransportNode"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	paramsTypeMap["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	paramsTypeMap["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
	queryParams["ping_ip"] = "ping_ip"
	queryParams["vnic"] = "vnic"
	queryParams["skip_validation"] = "skip_validation"
	queryParams["esx_mgmt_if_migration_dest"] = "esx_mgmt_if_migration_dest"
	queryParams["if_id"] = "if_id"
	queryParams["vnic_migration_dest"] = "vnic_migration_dest"
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
		"host_transport_node",
		"PATCH",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesRestoreclusterconfigInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesRestoreclusterconfigOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesRestoreclusterconfigRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
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
		"action=restore_cluster_config",
		"",
		"PUT",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesResynchostconfigInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesResynchostconfigOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesResynchostconfigRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
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
		"action=resync_host_config",
		"",
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	fields["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["host_transport_node"] = "HostTransportNode"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.HostTransportNodeBindingType)
}

func hostTransportNodesUpdateRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	fields["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["host_transport_node"] = "HostTransportNode"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	paramsTypeMap["vnic"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["skip_validation"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["host_transport_node"] = bindings.NewReferenceType(model.HostTransportNodeBindingType)
	paramsTypeMap["if_id"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["esx_mgmt_if_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["ping_ip"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["vnic_migration_dest"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
	queryParams["ping_ip"] = "ping_ip"
	queryParams["vnic"] = "vnic"
	queryParams["skip_validation"] = "skip_validation"
	queryParams["esx_mgmt_if_migration_dest"] = "esx_mgmt_if_migration_dest"
	queryParams["if_id"] = "if_id"
	queryParams["vnic_migration_dest"] = "vnic_migration_dest"
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
		"host_transport_node",
		"PUT",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostTransportNodesUpdatemaintenancemodeInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_id"] = bindings.NewStringType()
	fields["enforcementpoint_id"] = bindings.NewStringType()
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["action"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["action"] = "Action"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func hostTransportNodesUpdatemaintenancemodeOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func hostTransportNodesUpdatemaintenancemodeRestMetadata() protocol.OperationRestMetadata {
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
	fields["host_transport_node_id"] = bindings.NewStringType()
	fields["action"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["site_id"] = "SiteId"
	fieldNameMap["enforcementpoint_id"] = "EnforcementpointId"
	fieldNameMap["host_transport_node_id"] = "HostTransportNodeId"
	fieldNameMap["action"] = "Action"
	paramsTypeMap["host_transport_node_id"] = bindings.NewStringType()
	paramsTypeMap["enforcementpoint_id"] = bindings.NewStringType()
	paramsTypeMap["action"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewStringType()
	paramsTypeMap["siteId"] = bindings.NewStringType()
	paramsTypeMap["enforcementpointId"] = bindings.NewStringType()
	paramsTypeMap["hostTransportNodeId"] = bindings.NewStringType()
	pathParams["enforcementpoint_id"] = "enforcementpointId"
	pathParams["site_id"] = "siteId"
	pathParams["host_transport_node_id"] = "hostTransportNodeId"
	queryParams["action"] = "action"
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
		"POST",
		"/policy/api/v1/infra/sites/{siteId}/enforcement-points/{enforcementpointId}/host-transport-nodes/{hostTransportNodeId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
