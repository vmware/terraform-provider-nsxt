/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Multicast.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package locale_services

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func multicastGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func multicastGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
}

func multicastGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/multicast",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func multicastPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_multicast_config"] = "PolicyMulticastConfig"
	fieldNameMap["override"] = "Override"
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
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_multicast_config"] = "PolicyMulticastConfig"
	fieldNameMap["override"] = "Override"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	queryParams["override"] = "override"
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
		"policy_multicast_config",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/multicast",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func multicastUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_multicast_config"] = "PolicyMulticastConfig"
	fieldNameMap["override"] = "Override"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func multicastUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
}

func multicastUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_services_id"] = bindings.NewStringType()
	fields["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	fields["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_services_id"] = "LocaleServicesId"
	fieldNameMap["policy_multicast_config"] = "PolicyMulticastConfig"
	fieldNameMap["override"] = "Override"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["override"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["policy_multicast_config"] = bindings.NewReferenceType(model.PolicyMulticastConfigBindingType)
	paramsTypeMap["locale_services_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServicesId"] = bindings.NewStringType()
	pathParams["locale_services_id"] = "localeServicesId"
	pathParams["tier0_id"] = "tier0Id"
	queryParams["override"] = "override"
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
		"policy_multicast_config",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServicesId}/multicast",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


