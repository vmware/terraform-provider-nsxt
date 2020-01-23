/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: LbNodeUsageSummary.
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





func lbNodeUsageSummaryGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_usages"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["enforcement_point_path"] = "EnforcementPointPath"
	fieldNameMap["include_usages"] = "IncludeUsages"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func lbNodeUsageSummaryGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.AggregateLBNodeUsageSummaryBindingType)
}

func lbNodeUsageSummaryGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_usages"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["enforcement_point_path"] = "EnforcementPointPath"
	fieldNameMap["include_usages"] = "IncludeUsages"
	paramsTypeMap["include_usages"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["enforcement_point_path"] = bindings.NewOptionalType(bindings.NewStringType())
	queryParams["include_usages"] = "include_usages"
	queryParams["enforcement_point_path"] = "enforcement_point_path"
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
		"/policy/api/v1/infra/lb-node-usage-summary",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


