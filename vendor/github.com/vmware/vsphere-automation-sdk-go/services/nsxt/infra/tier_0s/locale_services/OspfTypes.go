// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Ospf.
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

func ospfGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ospfGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
}

func ospfGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func ospfPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["ospf_routing_config"] = "OspfRoutingConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ospfPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
}

func ospfPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["ospf_routing_config"] = "OspfRoutingConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"ospf_routing_config",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func ospfUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["ospf_routing_config"] = "OspfRoutingConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ospfUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
}

func ospfUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["ospf_routing_config"] = "OspfRoutingConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["ospf_routing_config"] = bindings.NewReferenceType(model.OspfRoutingConfigBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"ospf_routing_config",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/ospf",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
