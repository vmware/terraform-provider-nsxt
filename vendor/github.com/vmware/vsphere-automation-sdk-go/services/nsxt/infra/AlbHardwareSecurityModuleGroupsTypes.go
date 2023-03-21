// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: AlbHardwareSecurityModuleGroups.
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

func albHardwareSecurityModuleGroupsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["force"] = "Force"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlbHardwareSecurityModuleGroupsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func albHardwareSecurityModuleGroupsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["force"] = "Force"
	paramsTypeMap["force"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["albHardwaresecuritymodulegroupId"] = vapiBindings_.NewStringType()
	pathParams["alb_hardwaresecuritymodulegroup_id"] = "albHardwaresecuritymodulegroupId"
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
		"/policy/api/v1/infra/alb-hardware-security-module-groups/{albHardwaresecuritymodulegroupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albHardwareSecurityModuleGroupsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlbHardwareSecurityModuleGroupsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
}

func albHardwareSecurityModuleGroupsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	paramsTypeMap["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["albHardwaresecuritymodulegroupId"] = vapiBindings_.NewStringType()
	pathParams["alb_hardwaresecuritymodulegroup_id"] = "albHardwaresecuritymodulegroupId"
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
		"/policy/api/v1/infra/alb-hardware-security-module-groups/{albHardwaresecuritymodulegroupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albHardwareSecurityModuleGroupsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlbHardwareSecurityModuleGroupsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupApiResponseBindingType)
}

func albHardwareSecurityModuleGroupsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
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
		"/policy/api/v1/infra/alb-hardware-security-module-groups",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albHardwareSecurityModuleGroupsPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["a_LB_hardware_security_module_group"] = "ALBHardwareSecurityModuleGroup"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlbHardwareSecurityModuleGroupsPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func albHardwareSecurityModuleGroupsPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["a_LB_hardware_security_module_group"] = "ALBHardwareSecurityModuleGroup"
	paramsTypeMap["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	paramsTypeMap["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["albHardwaresecuritymodulegroupId"] = vapiBindings_.NewStringType()
	pathParams["alb_hardwaresecuritymodulegroup_id"] = "albHardwaresecuritymodulegroupId"
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
		"a_LB_hardware_security_module_group",
		"PATCH",
		"/policy/api/v1/infra/alb-hardware-security-module-groups/{albHardwaresecuritymodulegroupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albHardwareSecurityModuleGroupsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["a_LB_hardware_security_module_group"] = "ALBHardwareSecurityModuleGroup"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlbHardwareSecurityModuleGroupsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
}

func albHardwareSecurityModuleGroupsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	fields["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	fieldNameMap["alb_hardwaresecuritymodulegroup_id"] = "AlbHardwaresecuritymodulegroupId"
	fieldNameMap["a_LB_hardware_security_module_group"] = "ALBHardwareSecurityModuleGroup"
	paramsTypeMap["a_LB_hardware_security_module_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.ALBHardwareSecurityModuleGroupBindingType)
	paramsTypeMap["alb_hardwaresecuritymodulegroup_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["albHardwaresecuritymodulegroupId"] = vapiBindings_.NewStringType()
	pathParams["alb_hardwaresecuritymodulegroup_id"] = "albHardwaresecuritymodulegroupId"
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
		"a_LB_hardware_security_module_group",
		"PUT",
		"/policy/api/v1/infra/alb-hardware-security-module-groups/{albHardwaresecuritymodulegroupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
