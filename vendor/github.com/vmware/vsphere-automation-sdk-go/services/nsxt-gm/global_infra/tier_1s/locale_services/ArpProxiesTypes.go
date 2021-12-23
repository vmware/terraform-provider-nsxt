// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: ArpProxies.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"reflect"
)

// Possible value for ``source`` of method ArpProxies#list.
const ArpProxies_LIST_SOURCE_REALTIME = "realtime"

// Possible value for ``source`` of method ArpProxies#list.
const ArpProxies_LIST_SOURCE_CACHED = "cached"

func arpProxiesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["source"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["enforcement_point_path"] = "EnforcementPointPath"
	fieldNameMap["source"] = "Source"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func arpProxiesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyArpProxyTableListResultBindingType)
}

func arpProxiesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["tier1_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["source"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier1_id"] = "Tier1Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["enforcement_point_path"] = "EnforcementPointPath"
	fieldNameMap["source"] = "Source"
	paramsTypeMap["source"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier1_id"] = bindings.NewStringType()
	paramsTypeMap["tier1Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier1_id"] = "tier1Id"
	pathParams["locale_service_id"] = "localeServiceId"
	queryParams["enforcement_point_path"] = "enforcement_point_path"
	queryParams["source"] = "source"
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
		"/global-manager/api/v1/global-infra/tier-1s/{tier1Id}/locale-services/{localeServiceId}/arp-proxies",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
