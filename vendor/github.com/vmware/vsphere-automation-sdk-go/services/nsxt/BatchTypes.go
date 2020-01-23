/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Batch.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package nsx_policy

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func batchCreateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["batch_request"] = bindings.NewReferenceType(model.BatchRequestBindingType)
	fields["atomic"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["batch_request"] = "BatchRequest"
	fieldNameMap["atomic"] = "Atomic"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func batchCreateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.BatchResponseBindingType)
}

func batchCreateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["batch_request"] = bindings.NewReferenceType(model.BatchRequestBindingType)
	fields["atomic"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["batch_request"] = "BatchRequest"
	fieldNameMap["atomic"] = "Atomic"
	paramsTypeMap["batch_request"] = bindings.NewReferenceType(model.BatchRequestBindingType)
	paramsTypeMap["atomic"] = bindings.NewOptionalType(bindings.NewBooleanType())
	queryParams["atomic"] = "atomic"
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
		"batch_request",
		"POST",
		"/policy/api/v1/batch",
		resultHeaders,
		201,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


