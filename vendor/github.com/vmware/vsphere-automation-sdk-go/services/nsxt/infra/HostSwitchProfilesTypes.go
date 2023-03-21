// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: HostSwitchProfiles.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package infra

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

// Possible value for ``deploymentType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_DEPLOYMENT_TYPE_VIRTUAL_MACHINE = "VIRTUAL_MACHINE"

// Possible value for ``deploymentType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_DEPLOYMENT_TYPE_PHYSICAL_MACHINE = "PHYSICAL_MACHINE"

// Possible value for ``deploymentType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_DEPLOYMENT_TYPE_UNKNOWN = "UNKNOWN"

// Possible value for ``hostswitchProfileType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYUPLINKHOSTSWITCHPROFILE = "PolicyUplinkHostSwitchProfile"

// Possible value for ``hostswitchProfileType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYLLDPHOSTSWITCHPROFILE = "PolicyLldpHostSwitchProfile"

// Possible value for ``hostswitchProfileType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYNIOCPROFILE = "PolicyNiocProfile"

// Possible value for ``hostswitchProfileType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYEXTRACONFIGHOSTSWITCHPROFILE = "PolicyExtraConfigHostSwitchProfile"

// Possible value for ``hostswitchProfileType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_HOSTSWITCH_PROFILE_TYPE_POLICYVTEPHAHOSTSWITCHPROFILE = "PolicyVtepHAHostSwitchProfile"

// Possible value for ``nodeType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_NODE_TYPE_EDGENODE = "EdgeNode"

// Possible value for ``nodeType`` of method HostSwitchProfiles#list.
const HostSwitchProfiles_LIST_NODE_TYPE_PUBLICCLOUDGATEWAYNODE = "PublicCloudGatewayNode"

func hostSwitchProfilesDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func HostSwitchProfilesDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func hostSwitchProfilesDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	paramsTypeMap["host_switch_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["hostSwitchProfileId"] = vapiBindings_.NewStringType()
	pathParams["host_switch_profile_id"] = "hostSwitchProfileId"
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
		"/policy/api/v1/infra/host-switch-profiles/{hostSwitchProfileId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostSwitchProfilesGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func HostSwitchProfilesGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
}

func hostSwitchProfilesGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	paramsTypeMap["host_switch_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["hostSwitchProfileId"] = vapiBindings_.NewStringType()
	pathParams["host_switch_profile_id"] = "hostSwitchProfileId"
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
		"/policy/api/v1/infra/host-switch-profiles/{hostSwitchProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostSwitchProfilesListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["deployment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["hostswitch_profile_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["include_system_owned"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["max_active_uplink_count"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["deployment_type"] = "DeploymentType"
	fieldNameMap["hostswitch_profile_type"] = "HostswitchProfileType"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["include_system_owned"] = "IncludeSystemOwned"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["max_active_uplink_count"] = "MaxActiveUplinkCount"
	fieldNameMap["node_type"] = "NodeType"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["uplink_teaming_policy_name"] = "UplinkTeamingPolicyName"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func HostSwitchProfilesListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.PolicyHostSwitchProfilesListResultBindingType)
}

func hostSwitchProfilesListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["deployment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["hostswitch_profile_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["include_system_owned"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["max_active_uplink_count"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["deployment_type"] = "DeploymentType"
	fieldNameMap["hostswitch_profile_type"] = "HostswitchProfileType"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["include_system_owned"] = "IncludeSystemOwned"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["max_active_uplink_count"] = "MaxActiveUplinkCount"
	fieldNameMap["node_type"] = "NodeType"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["uplink_teaming_policy_name"] = "UplinkTeamingPolicyName"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["hostswitch_profile_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["uplink_teaming_policy_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["deployment_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["include_system_owned"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["max_active_uplink_count"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["node_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["hostswitch_profile_type"] = "hostswitch_profile_type"
	queryParams["uplink_teaming_policy_name"] = "uplink_teaming_policy_name"
	queryParams["deployment_type"] = "deployment_type"
	queryParams["include_system_owned"] = "include_system_owned"
	queryParams["max_active_uplink_count"] = "max_active_uplink_count"
	queryParams["node_type"] = "node_type"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
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
		"/policy/api/v1/infra/host-switch-profiles",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostSwitchProfilesPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fields["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	fieldNameMap["policy_base_host_switch_profile"] = "PolicyBaseHostSwitchProfile"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func HostSwitchProfilesPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
}

func hostSwitchProfilesPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fields["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	fieldNameMap["policy_base_host_switch_profile"] = "PolicyBaseHostSwitchProfile"
	paramsTypeMap["host_switch_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	paramsTypeMap["hostSwitchProfileId"] = vapiBindings_.NewStringType()
	pathParams["host_switch_profile_id"] = "hostSwitchProfileId"
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
		"policy_base_host_switch_profile",
		"PATCH",
		"/policy/api/v1/infra/host-switch-profiles/{hostSwitchProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func hostSwitchProfilesUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fields["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	fieldNameMap["policy_base_host_switch_profile"] = "PolicyBaseHostSwitchProfile"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func HostSwitchProfilesUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
}

func hostSwitchProfilesUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["host_switch_profile_id"] = vapiBindings_.NewStringType()
	fields["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	fieldNameMap["host_switch_profile_id"] = "HostSwitchProfileId"
	fieldNameMap["policy_base_host_switch_profile"] = "PolicyBaseHostSwitchProfile"
	paramsTypeMap["host_switch_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["policy_base_host_switch_profile"] = vapiBindings_.NewDynamicStructType([]vapiBindings_.ReferenceType{vapiBindings_.NewReferenceType(nsx_policyModel.PolicyBaseHostSwitchProfileBindingType)})
	paramsTypeMap["hostSwitchProfileId"] = vapiBindings_.NewStringType()
	pathParams["host_switch_profile_id"] = "hostSwitchProfileId"
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
		"policy_base_host_switch_profile",
		"PUT",
		"/policy/api/v1/infra/host-switch-profiles/{hostSwitchProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
