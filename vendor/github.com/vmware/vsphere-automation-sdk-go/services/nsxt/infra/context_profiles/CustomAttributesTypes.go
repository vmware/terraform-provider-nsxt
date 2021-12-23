// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: CustomAttributes.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package context_profiles

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

// Possible value for ``action`` of method CustomAttributes#create.
const CustomAttributes_CREATE_ACTION_ADD = "add"

// Possible value for ``action`` of method CustomAttributes#create.
const CustomAttributes_CREATE_ACTION_REMOVE = "remove"

func customAttributesCreateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
	fields["action"] = bindings.NewStringType()
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	fieldNameMap["action"] = "Action"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func customAttributesCreateOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func customAttributesCreateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
	fields["action"] = bindings.NewStringType()
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	fieldNameMap["action"] = "Action"
	paramsTypeMap["action"] = bindings.NewStringType()
	paramsTypeMap["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
	queryParams["action"] = "action"
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

func customAttributesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func customAttributesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func customAttributesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
	fieldNameMap["policy_attributes"] = "PolicyAttributes"
	paramsTypeMap["policy_attributes"] = bindings.NewReferenceType(model.PolicyAttributesBindingType)
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
