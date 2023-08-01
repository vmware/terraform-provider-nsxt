// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: LogicalRouters.
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

// Possible value for ``routerType`` of method LogicalRouters#list.
const LogicalRouters_LIST_ROUTER_TYPE_TIER0 = "TIER0"

// Possible value for ``routerType`` of method LogicalRouters#list.
const LogicalRouters_LIST_ROUTER_TYPE_TIER1 = "TIER1"

// Possible value for ``routerType`` of method LogicalRouters#list.
const LogicalRouters_LIST_ROUTER_TYPE_VRF = "VRF"

func logicalRoutersCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
	fieldNameMap["logical_router"] = "LogicalRouter"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
}

func logicalRoutersCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
	fieldNameMap["logical_router"] = "LogicalRouter"
	paramsTypeMap["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
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
		"logical_router",
		"POST",
		"/api/v1/logical-routers",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["cascade_delete_linked_ports"] = "CascadeDeleteLinkedPorts"
	fieldNameMap["force"] = "Force"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func logicalRoutersDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["cascade_delete_linked_ports"] = "CascadeDeleteLinkedPorts"
	fieldNameMap["force"] = "Force"
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["cascade_delete_linked_ports"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["logicalRouterId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_id"] = "logicalRouterId"
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
		"/api/v1/logical-routers/{logicalRouterId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
}

func logicalRoutersGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logicalRouterId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_id"] = "logicalRouterId"
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
		"/api/v1/logical-routers/{logicalRouterId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["router_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vrfs_on_logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["router_type"] = "RouterType"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["vrfs_on_logical_router_id"] = "VrfsOnLogicalRouterId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterListResultBindingType)
}

func logicalRoutersListRestMetadata() vapiProtocol_.OperationRestMetadata {
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
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["router_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vrfs_on_logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["router_type"] = "RouterType"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["vrfs_on_logical_router_id"] = "VrfsOnLogicalRouterId"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["router_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["vrfs_on_logical_router_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["router_type"] = "router_type"
	queryParams["page_size"] = "page_size"
	queryParams["vrfs_on_logical_router_id"] = "vrfs_on_logical_router_id"
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
		"/api/v1/logical-routers",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersReallocateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["service_router_allocation_config"] = vapiBindings_.NewReferenceType(nsxModel.ServiceRouterAllocationConfigBindingType)
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["service_router_allocation_config"] = "ServiceRouterAllocationConfig"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersReallocateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
}

func logicalRoutersReallocateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["service_router_allocation_config"] = vapiBindings_.NewReferenceType(nsxModel.ServiceRouterAllocationConfigBindingType)
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["service_router_allocation_config"] = "ServiceRouterAllocationConfig"
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["service_router_allocation_config"] = vapiBindings_.NewReferenceType(nsxModel.ServiceRouterAllocationConfigBindingType)
	paramsTypeMap["logicalRouterId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_id"] = "logicalRouterId"
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
		"action=reallocate",
		"service_router_allocation_config",
		"POST",
		"/api/v1/logical-routers/{logicalRouterId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersReprocessInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersReprocessOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func logicalRoutersReprocessRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logicalRouterId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_id"] = "logicalRouterId"
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
		"action=reprocess",
		"",
		"POST",
		"/api/v1/logical-routers/{logicalRouterId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalRoutersUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["logical_router"] = "LogicalRouter"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalRoutersUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
}

func logicalRoutersUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_router_id"] = vapiBindings_.NewStringType()
	fields["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
	fieldNameMap["logical_router_id"] = "LogicalRouterId"
	fieldNameMap["logical_router"] = "LogicalRouter"
	paramsTypeMap["logical_router_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logical_router"] = vapiBindings_.NewReferenceType(nsxModel.LogicalRouterBindingType)
	paramsTypeMap["logicalRouterId"] = vapiBindings_.NewStringType()
	pathParams["logical_router_id"] = "logicalRouterId"
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
		"logical_router",
		"PUT",
		"/api/v1/logical-routers/{logicalRouterId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
