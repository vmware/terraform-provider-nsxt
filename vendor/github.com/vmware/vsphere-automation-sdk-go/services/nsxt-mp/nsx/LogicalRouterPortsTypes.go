// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: LogicalRouterPorts.
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

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERUPLINKPORT = "LogicalRouterUpLinkPort"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERDOWNLINKPORT = "LogicalRouterDownLinkPort"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERLINKPORTONTIER0 = "LogicalRouterLinkPortOnTIER0"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERLINKPORTONTIER1 = "LogicalRouterLinkPortOnTIER1"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERLOOPBACKPORT = "LogicalRouterLoopbackPort"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERIPTUNNELPORT = "LogicalRouterIPTunnelPort"

// Possible value for ``resourceType`` of method LogicalRouterPorts#list.
const LogicalRouterPorts_LIST_RESOURCE_TYPE_LOGICALROUTERCENTRALIZEDSERVICEPORT = "LogicalRouterCentralizedServicePort"

func logicalRouterPortsCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
	fieldNameMap["logical_router_port"] = "LogicalRouterPort"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRouterPortsCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
}

func logicalRouterPortsCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
	fieldNameMap["logical_router_port"] = "LogicalRouterPort"
	paramsTypeMap["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
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
		"logical_router_port",
		"POST",
		"/api/v1/logical-router-ports",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRouterPortsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fields["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	fieldNameMap["cascade_delete_linked_ports"] = "CascadeDeleteLinkedPorts"
	fieldNameMap["force"] = "Force"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRouterPortsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func logicalRouterPortsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fields["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	fieldNameMap["cascade_delete_linked_ports"] = "CascadeDeleteLinkedPorts"
	fieldNameMap["force"] = "Force"
	paramsTypeMap["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["logical_router_port_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logicalRouterPortId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_port_id"] = "logicalRouterPortId"
	queryParams["cascade_delete_linked_ports"] = "cascade_delete_linked_ports"
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
		"",
		"",
		"DELETE",
		"/api/v1/logical-router-ports/{logicalRouterPortId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRouterPortsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRouterPortsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
}

func logicalRouterPortsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	paramsTypeMap["logical_router_port_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logicalRouterPortId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_port_id"] = "logicalRouterPortId"
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
		"/api/v1/logical-router-ports/{logicalRouterPortId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRouterPortsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["resource_type"] = "ResourceType"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRouterPortsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortListResultBindingType)
}

func logicalRouterPortsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["resource_type"] = "ResourceType"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["logical_router_id"] = "logical_router_id"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["resource_type"] = "resource_type"
	queryParams["sort_by"] = "sort_by"
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
		"/api/v1/logical-router-ports",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRouterPortsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fields["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	fieldNameMap["logical_router_port"] = "LogicalRouterPort"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRouterPortsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
}

func logicalRouterPortsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_port_id"] = vapiBindings_.NewStringType()
	fields["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
	fieldNameMap["logical_router_port_id"] = "LogicalRouterPortId"
	fieldNameMap["logical_router_port"] = "LogicalRouterPort"
	paramsTypeMap["logical_router_port"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsxModel.LogicalRouterPortBindingType)})
	paramsTypeMap["logical_router_port_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logicalRouterPortId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_port_id"] = "logicalRouterPortId"
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
		"logical_router_port",
		"PUT",
		"/api/v1/logical-router-ports/{logicalRouterPortId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
