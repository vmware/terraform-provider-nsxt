// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: LogicalSwitches.
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

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_DEFAULT = "DEFAULT"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_SERVICE_PLANE = "SERVICE_PLANE"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_DHCP_RELAY = "DHCP_RELAY"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_GLOBAL = "GLOBAL"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_INTER_ROUTER = "INTER_ROUTER"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_EVPN = "EVPN"

// Possible value for ``switchType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_SWITCH_TYPE_DVPG = "DVPG"

// Possible value for ``transportType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_TRANSPORT_TYPE_OVERLAY = "OVERLAY"

// Possible value for ``transportType`` of method LogicalSwitches#list.
const LogicalSwitches_LIST_TRANSPORT_TYPE_VLAN = "VLAN"

func logicalSwitchesCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
	fieldNameMap["logical_switch"] = "LogicalSwitch"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalSwitchesCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
}

func logicalSwitchesCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
	fieldNameMap["logical_switch"] = "LogicalSwitch"
	paramsTypeMap["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
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
		"logical_switch",
		"POST",
		"/api/v1/logical-switches",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalSwitchesDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fields["cascade"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["lswitch_id"] = "LswitchId"
	fieldNameMap["cascade"] = "Cascade"
	fieldNameMap["detach"] = "Detach"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalSwitchesDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func logicalSwitchesDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fields["cascade"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["lswitch_id"] = "LswitchId"
	fieldNameMap["cascade"] = "Cascade"
	fieldNameMap["detach"] = "Detach"
	paramsTypeMap["lswitch_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["cascade"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["detach"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["lswitchId"] = vapiBindings_.NewStringType()
	pathParams["lswitch_id"] = "lswitchId"
	queryParams["cascade"] = "cascade"
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
		"/api/v1/logical-switches/{lswitchId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalSwitchesGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fieldNameMap["lswitch_id"] = "LswitchId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalSwitchesGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
}

func logicalSwitchesGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fieldNameMap["lswitch_id"] = "LswitchId"
	paramsTypeMap["lswitch_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["lswitchId"] = vapiBindings_.NewStringType()
	pathParams["lswitch_id"] = "lswitchId"
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
		"/api/v1/logical-switches/{lswitchId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalSwitchesListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switch_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vlan"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["vni"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["diagnostic"] = "Diagnostic"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["switch_type"] = "SwitchType"
	fieldNameMap["switching_profile_id"] = "SwitchingProfileId"
	fieldNameMap["transport_type"] = "TransportType"
	fieldNameMap["transport_zone_id"] = "TransportZoneId"
	fieldNameMap["uplink_teaming_policy_name"] = "UplinkTeamingPolicyName"
	fieldNameMap["vlan"] = "Vlan"
	fieldNameMap["vni"] = "Vni"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalSwitchesListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchListResultBindingType)
}

func logicalSwitchesListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switch_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vlan"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["vni"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["diagnostic"] = "Diagnostic"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["switch_type"] = "SwitchType"
	fieldNameMap["switching_profile_id"] = "SwitchingProfileId"
	fieldNameMap["transport_type"] = "TransportType"
	fieldNameMap["transport_zone_id"] = "TransportZoneId"
	fieldNameMap["uplink_teaming_policy_name"] = "UplinkTeamingPolicyName"
	fieldNameMap["vlan"] = "Vlan"
	fieldNameMap["vni"] = "Vni"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["switching_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["switch_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vni"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["vlan"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["transport_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["diagnostic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["switching_profile_id"] = "switching_profile_id"
	queryParams["sort_by"] = "sort_by"
	queryParams["switch_type"] = "switch_type"
	queryParams["uplink_teaming_policy_name"] = "uplink_teaming_policy_name"
	queryParams["vni"] = "vni"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["vlan"] = "vlan"
	queryParams["included_fields"] = "included_fields"
	queryParams["transport_type"] = "transport_type"
	queryParams["transport_zone_id"] = "transport_zone_id"
	queryParams["diagnostic"] = "diagnostic"
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
		"/api/v1/logical-switches",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func logicalSwitchesUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fields["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
	fieldNameMap["lswitch_id"] = "LswitchId"
	fieldNameMap["logical_switch"] = "LogicalSwitch"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func LogicalSwitchesUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
}

func logicalSwitchesUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["lswitch_id"] = vapiBindings_.NewStringType()
	fields["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
	fieldNameMap["lswitch_id"] = "LswitchId"
	fieldNameMap["logical_switch"] = "LogicalSwitch"
	paramsTypeMap["lswitch_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["logical_switch"] = vapiBindings_.NewReferenceType(nsxModel.LogicalSwitchBindingType)
	paramsTypeMap["lswitchId"] = vapiBindings_.NewStringType()
	pathParams["lswitch_id"] = "lswitchId"
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
		"logical_switch",
		"PUT",
		"/api/v1/logical-switches/{lswitchId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
