/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: IpsecVpnTunnelProfiles.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func ipsecVpnTunnelProfilesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipsecVpnTunnelProfilesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipsecVpnTunnelProfilesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	paramsTypeMap["tunnel_profile_id"] = bindings.NewStringType()
	paramsTypeMap["tunnelProfileId"] = bindings.NewStringType()
	pathParams["tunnel_profile_id"] = "tunnelProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"DELETE",
		"/policy/api/v1/infra/ipsec-vpn-tunnel-profiles/{tunnelProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipsecVpnTunnelProfilesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipsecVpnTunnelProfilesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
}

func ipsecVpnTunnelProfilesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	paramsTypeMap["tunnel_profile_id"] = bindings.NewStringType()
	paramsTypeMap["tunnelProfileId"] = bindings.NewStringType()
	pathParams["tunnel_profile_id"] = "tunnelProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/ipsec-vpn-tunnel-profiles/{tunnelProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipsecVpnTunnelProfilesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipsecVpnTunnelProfilesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPSecVpnTunnelProfileListResultBindingType)
}

func ipsecVpnTunnelProfilesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/ipsec-vpn-tunnel-profiles",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipsecVpnTunnelProfilesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fields["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	fieldNameMap["ip_sec_vpn_tunnel_profile"] = "IpSecVpnTunnelProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipsecVpnTunnelProfilesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipsecVpnTunnelProfilesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fields["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	fieldNameMap["ip_sec_vpn_tunnel_profile"] = "IpSecVpnTunnelProfile"
	paramsTypeMap["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	paramsTypeMap["tunnel_profile_id"] = bindings.NewStringType()
	paramsTypeMap["tunnelProfileId"] = bindings.NewStringType()
	pathParams["tunnel_profile_id"] = "tunnelProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"ip_sec_vpn_tunnel_profile",
		"PATCH",
		"/policy/api/v1/infra/ipsec-vpn-tunnel-profiles/{tunnelProfileId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func ipsecVpnTunnelProfilesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fields["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	fieldNameMap["ip_sec_vpn_tunnel_profile"] = "IpSecVpnTunnelProfile"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipsecVpnTunnelProfilesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
}

func ipsecVpnTunnelProfilesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tunnel_profile_id"] = bindings.NewStringType()
	fields["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	fieldNameMap["tunnel_profile_id"] = "TunnelProfileId"
	fieldNameMap["ip_sec_vpn_tunnel_profile"] = "IpSecVpnTunnelProfile"
	paramsTypeMap["ip_sec_vpn_tunnel_profile"] = bindings.NewReferenceType(model.IPSecVpnTunnelProfileBindingType)
	paramsTypeMap["tunnel_profile_id"] = bindings.NewStringType()
	paramsTypeMap["tunnelProfileId"] = bindings.NewStringType()
	pathParams["tunnel_profile_id"] = "tunnelProfileId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"ip_sec_vpn_tunnel_profile",
		"PUT",
		"/policy/api/v1/infra/ipsec-vpn-tunnel-profiles/{tunnelProfileId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


