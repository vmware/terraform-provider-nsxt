// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Multicast.
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

func multicastGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func multicastGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
}

func multicastGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/multicast",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func multicastPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_tier1_multicast_config"] = "PolicyTier1MulticastConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func multicastPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func multicastPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_tier1_multicast_config"] = "PolicyTier1MulticastConfig"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
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
		"policy_tier1_multicast_config",
		"PATCH",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/multicast",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func multicastUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_tier1_multicast_config"] = "PolicyTier1MulticastConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func multicastUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
}

func multicastUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_tier1_multicast_config"] = "PolicyTier1MulticastConfig"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["policy_tier1_multicast_config"] = bindings.NewReferenceType(model.PolicyTier1MulticastConfigBindingType)
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier1_id"] = "tier1Id"
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
		"policy_tier1_multicast_config",
		"PUT",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/locale-services/{localeServicesId}/multicast",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
