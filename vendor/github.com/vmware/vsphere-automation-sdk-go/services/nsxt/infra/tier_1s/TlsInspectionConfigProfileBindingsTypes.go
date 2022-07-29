// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: TlsInspectionConfigProfileBindings.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func tlsInspectionConfigProfileBindingsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tlsInspectionConfigProfileBindingsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func tlsInspectionConfigProfileBindingsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["tlsInspectionConfigProfileBindingId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["tls_inspection_config_profile_binding_id"] = "tlsInspectionConfigProfileBindingId"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/tls-inspection-config-profile-bindings/{tlsInspectionConfigProfileBindingId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tlsInspectionConfigProfileBindingsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tlsInspectionConfigProfileBindingsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
}

func tlsInspectionConfigProfileBindingsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["tlsInspectionConfigProfileBindingId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["tls_inspection_config_profile_binding_id"] = "tlsInspectionConfigProfileBindingId"
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
		"/policy/api/v1/infra/tier-1s/{tier1Id}/tls-inspection-config-profile-bindings/{tlsInspectionConfigProfileBindingId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tlsInspectionConfigProfileBindingsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fields["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	fieldNameMap["tls_config_profile_binding_map"] = "TlsConfigProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tlsInspectionConfigProfileBindingsPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
}

func tlsInspectionConfigProfileBindingsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fields["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	fieldNameMap["tls_config_profile_binding_map"] = "TlsConfigProfileBindingMap"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	paramsTypeMap["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["tlsInspectionConfigProfileBindingId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["tls_inspection_config_profile_binding_id"] = "tlsInspectionConfigProfileBindingId"
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
		"tls_config_profile_binding_map",
		"PATCH",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/tls-inspection-config-profile-bindings/{tlsInspectionConfigProfileBindingId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tlsInspectionConfigProfileBindingsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fields["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	fieldNameMap["tls_config_profile_binding_map"] = "TlsConfigProfileBindingMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tlsInspectionConfigProfileBindingsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
}

func tlsInspectionConfigProfileBindingsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	fields["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["tls_inspection_config_profile_binding_id"] = "TlsInspectionConfigProfileBindingId"
	fieldNameMap["tls_config_profile_binding_map"] = "TlsConfigProfileBindingMap"
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tls_inspection_config_profile_binding_id"] = bindings.NewStringType()
	paramsTypeMap["tls_config_profile_binding_map"] = bindings.NewReferenceType(model.TlsConfigProfileBindingMapBindingType)
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["tlsInspectionConfigProfileBindingId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["tls_inspection_config_profile_binding_id"] = "tlsInspectionConfigProfileBindingId"
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
		"tls_config_profile_binding_map",
		"PUT",
		"/policy/api/v1/infra/tier-1s/{tier1Id}/tls-inspection-config-profile-bindings/{tlsInspectionConfigProfileBindingId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
