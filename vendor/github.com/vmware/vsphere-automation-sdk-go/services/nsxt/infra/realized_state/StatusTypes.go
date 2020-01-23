/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Status.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package realized_state

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func statusGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["intent_path"] = bindings.NewStringType()
	fields["include_enforced_status"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["include_enforced_status"] = "IncludeEnforcedStatus"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func statusGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ConsolidatedRealizedStatusBindingType)
}

func statusGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["intent_path"] = bindings.NewStringType()
	fields["include_enforced_status"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["include_enforced_status"] = "IncludeEnforcedStatus"
	paramsTypeMap["include_enforced_status"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["intent_path"] = bindings.NewStringType()
	queryParams["intent_path"] = "intent_path"
	queryParams["include_enforced_status"] = "include_enforced_status"
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
		"/policy/api/v1/infra/realized-state/status",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


