/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Evpn.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package tier_0s

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func evpnGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.EvpnConfigBindingType)
}

func evpnGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/evpn",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func evpnPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["evpn_config"] = "EvpnConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func evpnPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["evpn_config"] = "EvpnConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"evpn_config",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/evpn",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func evpnUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["evpn_config"] = "EvpnConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func evpnUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.EvpnConfigBindingType)
}

func evpnUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["evpn_config"] = "EvpnConfig"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["evpn_config"] = bindings.NewReferenceType(model.EvpnConfigBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"evpn_config",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/evpn",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


