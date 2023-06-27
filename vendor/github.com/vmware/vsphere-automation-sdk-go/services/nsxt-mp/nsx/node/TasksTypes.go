// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Tasks.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package node

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"reflect"
)

func tasksCancelInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["task_id"] = vapiBindings_.NewStringType()
	fieldNameMap["task_id"] = "TaskId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TasksCancelOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func tasksCancelRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["task_id"] = vapiBindings_.NewStringType()
	fieldNameMap["task_id"] = "TaskId"
	paramsTypeMap["task_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["taskId"] = vapiBindings_.NewStringType()
	pathParams["task_id"] = "taskId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		dispatchHeaderParams,
		bodyFieldsMap,
		"action=cancel",
		"",
		"POST",
		"/api/v1/node/tasks/{taskId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tasksDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["task_id"] = vapiBindings_.NewStringType()
	fieldNameMap["task_id"] = "TaskId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TasksDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func tasksDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["task_id"] = vapiBindings_.NewStringType()
	fieldNameMap["task_id"] = "TaskId"
	paramsTypeMap["task_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["taskId"] = vapiBindings_.NewStringType()
	pathParams["task_id"] = "taskId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/api/v1/node/tasks/{taskId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tasksGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["task_id"] = vapiBindings_.NewStringType()
	fields["suppress_redirect"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["task_id"] = "TaskId"
	fieldNameMap["suppress_redirect"] = "SuppressRedirect"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TasksGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ApplianceManagementTaskPropertiesBindingType)
}

func tasksGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["task_id"] = vapiBindings_.NewStringType()
	fields["suppress_redirect"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fieldNameMap["task_id"] = "TaskId"
	fieldNameMap["suppress_redirect"] = "SuppressRedirect"
	paramsTypeMap["task_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["suppress_redirect"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["taskId"] = vapiBindings_.NewStringType()
	pathParams["task_id"] = "taskId"
	queryParams["suppress_redirect"] = "suppress_redirect"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/api/v1/node/tasks/{taskId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func tasksListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_method"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_uri"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["user"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["fields"] = "Fields"
	fieldNameMap["request_method"] = "RequestMethod"
	fieldNameMap["request_path"] = "RequestPath"
	fieldNameMap["request_uri"] = "RequestUri"
	fieldNameMap["status"] = "Status"
	fieldNameMap["user"] = "User"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func TasksListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ApplianceManagementTaskListResultBindingType)
}

func tasksListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_method"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["request_uri"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["user"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["fields"] = "Fields"
	fieldNameMap["request_method"] = "RequestMethod"
	fieldNameMap["request_path"] = "RequestPath"
	fieldNameMap["request_uri"] = "RequestUri"
	fieldNameMap["status"] = "Status"
	fieldNameMap["user"] = "User"
	paramsTypeMap["request_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["request_method"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["user"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["request_uri"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["request_path"] = "request_path"
	queryParams["request_method"] = "request_method"
	queryParams["fields"] = "fields"
	queryParams["user"] = "user"
	queryParams["request_uri"] = "request_uri"
	queryParams["status"] = "status"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]map[string]string{}
	return vapiProtocol_.NewOperationRestMetadata(
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
		"/api/v1/node/tasks",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
