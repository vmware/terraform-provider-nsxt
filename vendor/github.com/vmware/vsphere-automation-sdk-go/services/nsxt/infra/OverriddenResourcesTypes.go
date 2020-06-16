/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: OverriddenResources.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func overriddenResourcesListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["intent_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["site_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["site_path"] = "SitePath"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func overriddenResourcesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.OverriddenResourceListResultBindingType)
}

func overriddenResourcesListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["intent_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["site_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["site_path"] = "SitePath"
	paramsTypeMap["site_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["intent_path"] = bindings.NewOptionalType(bindings.NewStringType())
	queryParams["site_path"] = "site_path"
	queryParams["intent_path"] = "intent_path"
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
		"/policy/api/v1/infra/overridden-resources",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


