// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: TransportNodeProfiles.
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

func transportNodeProfilesCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_profile"] = "TransportNodeProfile"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeProfilesCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
}

func transportNodeProfilesCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["transport_node_profile"] = "TransportNodeProfile"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	paramsTypeMap["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	queryParams["override_nsx_ownership"] = "override_nsx_ownership"
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
		"transport_node_profile",
		"POST",
		"/api/v1/transport-node-profiles",
		"",
		resultHeaders,
		201,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeProfilesDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeProfilesDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func transportNodeProfilesDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	paramsTypeMap["transport_node_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeProfileId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_profile_id"] = "transportNodeProfileId"
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
		"/api/v1/transport-node-profiles/{transportNodeProfileId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeProfilesGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeProfilesGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
}

func transportNodeProfilesGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	paramsTypeMap["transport_node_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeProfileId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_profile_id"] = "transportNodeProfileId"
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
		"/api/v1/transport-node-profiles/{transportNodeProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeProfilesListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeProfilesListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileListResultBindingType)
}

func transportNodeProfilesListRestMetadata() vapiProtocol_.OperationRestMetadata {
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
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
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
		"/api/v1/transport-node-profiles",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func transportNodeProfilesUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fields["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	fields["esx_mgmt_if_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["if_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["ping_ip"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["skip_validation"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["vnic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vnic_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	fieldNameMap["transport_node_profile"] = "TransportNodeProfile"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TransportNodeProfilesUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
}

func transportNodeProfilesUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["transport_node_profile_id"] = vapiBindings_.NewStringType()
	fields["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	fields["esx_mgmt_if_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["if_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["ping_ip"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["skip_validation"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["vnic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vnic_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["transport_node_profile_id"] = "TransportNodeProfileId"
	fieldNameMap["transport_node_profile"] = "TransportNodeProfile"
	fieldNameMap["esx_mgmt_if_migration_dest"] = "EsxMgmtIfMigrationDest"
	fieldNameMap["if_id"] = "IfId"
	fieldNameMap["override_nsx_ownership"] = "OverrideNsxOwnership"
	fieldNameMap["ping_ip"] = "PingIp"
	fieldNameMap["skip_validation"] = "SkipValidation"
	fieldNameMap["vnic"] = "Vnic"
	fieldNameMap["vnic_migration_dest"] = "VnicMigrationDest"
	paramsTypeMap["ping_ip"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["transport_node_profile"] = vapiBindings_.NewReferenceType(nsxModel.TransportNodeProfileBindingType)
	paramsTypeMap["vnic"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["skip_validation"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["esx_mgmt_if_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["if_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vnic_migration_dest"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["override_nsx_ownership"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["transport_node_profile_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["transportNodeProfileId"] = vapiBindings_.NewStringType()
	pathParams["transport_node_profile_id"] = "transportNodeProfileId"
	queryParams["ping_ip"] = "ping_ip"
	queryParams["vnic"] = "vnic"
	queryParams["skip_validation"] = "skip_validation"
	queryParams["esx_mgmt_if_migration_dest"] = "esx_mgmt_if_migration_dest"
	queryParams["if_id"] = "if_id"
	queryParams["vnic_migration_dest"] = "vnic_migration_dest"
	queryParams["override_nsx_ownership"] = "override_nsx_ownership"
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
		"transport_node_profile",
		"PUT",
		"/api/v1/transport-node-profiles/{transportNodeProfileId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
