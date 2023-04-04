// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: CustomAttributes.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package context_profiles

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

// Possible value for ``action`` of method CustomAttributes#create.
const CustomAttributes_CREATE_ACTION_ADD = "add"

// Possible value for ``action`` of method CustomAttributes#create.
const CustomAttributes_CREATE_ACTION_REMOVE = "remove"

func customAttributesCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
	fields["action"] = vapiBindings_.NewStringType()
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	fieldNameMap["action"] = "Action"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func CustomAttributesCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func customAttributesCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
	fields["action"] = vapiBindings_.NewStringType()
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	fieldNameMap["action"] = "Action"
	paramsTypeMap["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
	paramsTypeMap["action"] = vapiBindings_.NewStringType()
	queryParams["action"] = "action"
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
		"policy_attributes",
		"POST",
		"/policy/api/v1/infra/context-profiles/custom-attributes",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func customAttributesPatchInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func CustomAttributesPatchOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func customAttributesPatchRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	paramsTypeMap["policy_attributes"] = vapiBindings_.NewReferenceType(nsx_policyModel.PolicyAttributesBindingType)
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
		"policy_attributes",
		"PATCH",
		"/policy/api/v1/infra/context-profiles/custom-attributes",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
