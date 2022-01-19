// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: EvpnTunnelEndpoints.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func evpnTunnelEndpointsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnTunnelEndpointsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func evpnTunnelEndpointsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tunnel_endpoint_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["tunnelEndpointId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["tunnel_endpoint_id"] = "tunnelEndpointId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/evpn-tunnel-endpoints/{tunnelEndpointId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func evpnTunnelEndpointsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnTunnelEndpointsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
}

func evpnTunnelEndpointsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tunnel_endpoint_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["tunnelEndpointId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["tunnel_endpoint_id"] = "tunnelEndpointId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/evpn-tunnel-endpoints/{tunnelEndpointId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func evpnTunnelEndpointsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnTunnelEndpointsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.EvpnTunnelEndpointConfigListResultBindingType)
}

func evpnTunnelEndpointsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/evpn-tunnel-endpoints",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func evpnTunnelEndpointsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fields["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	fieldNameMap["evpn_tunnel_endpoint_config"] = "EvpnTunnelEndpointConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnTunnelEndpointsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func evpnTunnelEndpointsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fields["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	fieldNameMap["evpn_tunnel_endpoint_config"] = "EvpnTunnelEndpointConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	paramsTypeMap["tunnel_endpoint_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["tunnelEndpointId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["tunnel_endpoint_id"] = "tunnelEndpointId"
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
		"evpn_tunnel_endpoint_config",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/evpn-tunnel-endpoints/{tunnelEndpointId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func evpnTunnelEndpointsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fields["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	fieldNameMap["evpn_tunnel_endpoint_config"] = "EvpnTunnelEndpointConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnTunnelEndpointsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
}

func evpnTunnelEndpointsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["tunnel_endpoint_id"] = bindings.NewStringType()
	fields["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["tunnel_endpoint_id"] = "TunnelEndpointId"
	fieldNameMap["evpn_tunnel_endpoint_config"] = "EvpnTunnelEndpointConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["evpn_tunnel_endpoint_config"] = bindings.NewReferenceType(model.EvpnTunnelEndpointConfigBindingType)
	paramsTypeMap["tunnel_endpoint_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	paramsTypeMap["tunnelEndpointId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["tunnel_endpoint_id"] = "tunnelEndpointId"
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
		"evpn_tunnel_endpoint_config",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/evpn-tunnel-endpoints/{tunnelEndpointId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
