// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: OrgRoot.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package nsx_policy

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func orgRootGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["base_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["type_filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["base_path"] = "BasePath"
	fieldNameMap["filter"] = "Filter"
	fieldNameMap["type_filter"] = "TypeFilter"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func OrgRootGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsx_policyModel.OrgRootBindingType)
}

func orgRootGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["base_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["type_filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["base_path"] = "BasePath"
	fieldNameMap["filter"] = "Filter"
	fieldNameMap["type_filter"] = "TypeFilter"
	paramsTypeMap["filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["base_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["type_filter"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["filter"] = "filter"
	queryParams["base_path"] = "base_path"
	queryParams["type_filter"] = "type_filter"
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
		"/policy/api/v1/org-root",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func orgRootPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org_root"] = vapiBindings_.NewReferenceType(nsx_policyModel.OrgRootBindingType)
	fields["enforce_revision_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["org_root"] = "OrgRoot"
	fieldNameMap["enforce_revision_check"] = "EnforceRevisionCheck"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func OrgRootPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func orgRootPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["org_root"] = vapiBindings_.NewReferenceType(nsx_policyModel.OrgRootBindingType)
	fields["enforce_revision_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["org_root"] = "OrgRoot"
	fieldNameMap["enforce_revision_check"] = "EnforceRevisionCheck"
	paramsTypeMap["org_root"] = vapiBindings_.NewReferenceType(nsx_policyModel.OrgRootBindingType)
	paramsTypeMap["enforce_revision_check"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	queryParams["enforce_revision_check"] = "enforce_revision_check"
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
		"org_root",
		"PATCH",
		"/policy/api/v1/org-root",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
