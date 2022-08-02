// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: GlobalInfra.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package nsx_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func globalInfraGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["base_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["filter"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["type_filter"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["base_path"] = "BasePath"
	fieldNameMap["filter"] = "Filter"
	fieldNameMap["type_filter"] = "TypeFilter"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func globalInfraGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.InfraBindingType)
}

func globalInfraGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["base_path"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["filter"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["type_filter"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["base_path"] = "BasePath"
	fieldNameMap["filter"] = "Filter"
	fieldNameMap["type_filter"] = "TypeFilter"
	paramsTypeMap["base_path"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["filter"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["type_filter"] = bindings.NewOptionalType(bindings.NewStringType())
	queryParams["filter"] = "filter"
	queryParams["base_path"] = "base_path"
	queryParams["type_filter"] = "type_filter"
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
		"/policy/api/v1/global-infra",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func globalInfraPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["infra"] = bindings.NewReferenceType(model.InfraBindingType)
	fields["enforce_revision_check"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["infra"] = "Infra"
	fieldNameMap["enforce_revision_check"] = "EnforceRevisionCheck"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func globalInfraPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func globalInfraPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["infra"] = bindings.NewReferenceType(model.InfraBindingType)
	fields["enforce_revision_check"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["infra"] = "Infra"
	fieldNameMap["enforce_revision_check"] = "EnforceRevisionCheck"
	paramsTypeMap["enforce_revision_check"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["infra"] = bindings.NewReferenceType(model.InfraBindingType)
	queryParams["enforce_revision_check"] = "enforce_revision_check"
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
		"infra",
		"PATCH",
		"/policy/api/v1/global-infra",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
