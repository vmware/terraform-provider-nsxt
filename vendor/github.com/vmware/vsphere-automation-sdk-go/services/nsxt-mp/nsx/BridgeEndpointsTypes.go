// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: BridgeEndpoints.
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

func bridgeEndpointsCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
	fieldNameMap["bridge_endpoint"] = "BridgeEndpoint"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func BridgeEndpointsCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
}

func bridgeEndpointsCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
	fieldNameMap["bridge_endpoint"] = "BridgeEndpoint"
	paramsTypeMap["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
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
		"bridge_endpoint",
		"POST",
		"/api/v1/bridge-endpoints",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func bridgeEndpointsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func BridgeEndpointsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func bridgeEndpointsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	paramsTypeMap["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["bridgeendpointId"] = vapiBindings_.NewStringType()
	pathParams["bridgeendpoint_id"] = "bridgeendpointId"
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
		"/api/v1/bridge-endpoints/{bridgeendpointId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func bridgeEndpointsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func BridgeEndpointsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
}

func bridgeEndpointsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	paramsTypeMap["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["bridgeendpointId"] = vapiBindings_.NewStringType()
	pathParams["bridgeendpoint_id"] = "bridgeendpointId"
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
		"/api/v1/bridge-endpoints/{bridgeendpointId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func bridgeEndpointsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["bridge_endpoint_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vlan_transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["bridge_cluster_id"] = "BridgeClusterId"
	fieldNameMap["bridge_endpoint_profile_id"] = "BridgeEndpointProfileId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["vlan_transport_zone_id"] = "VlanTransportZoneId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func BridgeEndpointsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointListResultBindingType)
}

func bridgeEndpointsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["bridge_endpoint_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vlan_transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["bridge_cluster_id"] = "BridgeClusterId"
	fieldNameMap["bridge_endpoint_profile_id"] = "BridgeEndpointProfileId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["logical_switch_id"] = "LogicalSwitchId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["vlan_transport_zone_id"] = "VlanTransportZoneId"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["bridge_endpoint_profile_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vlan_transport_zone_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["bridge_cluster_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["logical_switch_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["bridge_endpoint_profile_id"] = "bridge_endpoint_profile_id"
	queryParams["sort_by"] = "sort_by"
	queryParams["vlan_transport_zone_id"] = "vlan_transport_zone_id"
	queryParams["bridge_cluster_id"] = "bridge_cluster_id"
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
		"/api/v1/bridge-endpoints",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func bridgeEndpointsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fields["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	fieldNameMap["bridge_endpoint"] = "BridgeEndpoint"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func BridgeEndpointsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
}

func bridgeEndpointsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	fields["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
	fieldNameMap["bridgeendpoint_id"] = "BridgeendpointId"
	fieldNameMap["bridge_endpoint"] = "BridgeEndpoint"
	paramsTypeMap["bridge_endpoint"] = vapiBindings_.NewReferenceType(nsxModel.BridgeEndpointBindingType)
	paramsTypeMap["bridgeendpoint_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["bridgeendpointId"] = vapiBindings_.NewStringType()
	pathParams["bridgeendpoint_id"] = "bridgeendpointId"
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
		"bridge_endpoint",
		"PUT",
		"/api/v1/bridge-endpoints/{bridgeendpointId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
