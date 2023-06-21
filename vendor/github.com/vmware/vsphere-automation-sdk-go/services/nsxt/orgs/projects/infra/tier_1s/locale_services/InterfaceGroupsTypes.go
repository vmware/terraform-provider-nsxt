// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: InterfaceGroups.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package locale_services

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func interfaceGroupsDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func InterfaceGroupsDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func interfaceGroupsDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	paramsTypeMap["locale_service_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["project_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["group_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["org_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["orgId"] = vapiBindings_.NewStringType()
	paramsTypeMap["projectId"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1Id"] = vapiBindings_.NewStringType()
	paramsTypeMap["localeServiceId"] = vapiBindings_.NewStringType()
	paramsTypeMap["groupId"] = vapiBindings_.NewStringType()
	pathParams["group_id"] = "groupId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	pathParams["project_id"] = "projectId"
	pathParams["org_id"] = "orgId"
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
		"/policy/api/v1/orgs/{orgId}/projects/{projectId}/infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/interface-groups/{groupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfaceGroupsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func InterfaceGroupsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
}

func interfaceGroupsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	paramsTypeMap["locale_service_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["project_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["group_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["org_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["orgId"] = vapiBindings_.NewStringType()
	paramsTypeMap["projectId"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1Id"] = vapiBindings_.NewStringType()
	paramsTypeMap["localeServiceId"] = vapiBindings_.NewStringType()
	paramsTypeMap["groupId"] = vapiBindings_.NewStringType()
	pathParams["group_id"] = "groupId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	pathParams["project_id"] = "projectId"
	pathParams["org_id"] = "orgId"
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
		"/policy/api/v1/orgs/{orgId}/projects/{projectId}/infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/interface-groups/{groupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfaceGroupsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func InterfaceGroupsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupListResultBindingType)
}

func interfaceGroupsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["locale_service_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["project_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["org_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["include_mark_for_delete_objects"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["orgId"] = vapiBindings_.NewStringType()
	paramsTypeMap["projectId"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1Id"] = vapiBindings_.NewStringType()
	paramsTypeMap["localeServiceId"] = vapiBindings_.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	pathParams["project_id"] = "projectId"
	pathParams["org_id"] = "orgId"
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
		"/policy/api/v1/orgs/{orgId}/projects/{projectId}/infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/interface-groups",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfaceGroupsPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fields["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	fieldNameMap["tier1_interface_group"] = "Tier1InterfaceGroup"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func InterfaceGroupsPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func interfaceGroupsPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fields["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	fieldNameMap["tier1_interface_group"] = "Tier1InterfaceGroup"
	paramsTypeMap["locale_service_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["project_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["group_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["org_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	paramsTypeMap["orgId"] = vapiBindings_.NewStringType()
	paramsTypeMap["projectId"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1Id"] = vapiBindings_.NewStringType()
	paramsTypeMap["localeServiceId"] = vapiBindings_.NewStringType()
	paramsTypeMap["groupId"] = vapiBindings_.NewStringType()
	pathParams["group_id"] = "groupId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	pathParams["project_id"] = "projectId"
	pathParams["org_id"] = "orgId"
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
		"tier1_interface_group",
		"PATCH",
		"/policy/api/v1/orgs/{orgId}/projects/{projectId}/infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/interface-groups/{groupId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func interfaceGroupsUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fields["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	fieldNameMap["tier1_interface_group"] = "Tier1InterfaceGroup"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func InterfaceGroupsUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
}

func interfaceGroupsUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_id"] = vapiBindings_.NewStringType()
	fields["project_id"] = vapiBindings_.NewStringType()
	fields["tier1_id"] = vapiBindings_.NewStringType()
	fields["locale_service_id"] = vapiBindings_.NewStringType()
	fields["group_id"] = vapiBindings_.NewStringType()
	fields["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	fieldNameMap["org_id"] = "OrgId"
	fieldNameMap["project_id"] = "ProjectId"
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["group_id"] = "GroupId"
	fieldNameMap["tier1_interface_group"] = "Tier1InterfaceGroup"
	paramsTypeMap["locale_service_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["project_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["group_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["org_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1_interface_group"] = vapiBindings_.NewReferenceType(nsx_policyModel.Tier1InterfaceGroupBindingType)
	paramsTypeMap["orgId"] = vapiBindings_.NewStringType()
	paramsTypeMap["projectId"] = vapiBindings_.NewStringType()
	paramsTypeMap["tier1Id"] = vapiBindings_.NewStringType()
	paramsTypeMap["localeServiceId"] = vapiBindings_.NewStringType()
	paramsTypeMap["groupId"] = vapiBindings_.NewStringType()
	pathParams["group_id"] = "groupId"
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	pathParams["project_id"] = "projectId"
	pathParams["org_id"] = "orgId"
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
		"tier1_interface_group",
		"PUT",
		"/policy/api/v1/orgs/{orgId}/projects/{projectId}/infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/interface-groups/{groupId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
