// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: AlbAutoScaleLaunchConfigs.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"reflect"
)

func albAutoScaleLaunchConfigsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["force"] = "Force"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albAutoScaleLaunchConfigsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func albAutoScaleLaunchConfigsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["force"] = "Force"
	paramsTypeMap["force"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	paramsTypeMap["albAutoscalelaunchconfigId"] = bindings.NewStringType()
	pathParams["alb_autoscalelaunchconfig_id"] = "albAutoscalelaunchconfigId"
	queryParams["force"] = "force"
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
		"DELETE",
		"/policy/api/v1/infra/alb-auto-scale-launch-configs/{albAutoscalelaunchconfigId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albAutoScaleLaunchConfigsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albAutoScaleLaunchConfigsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
}

func albAutoScaleLaunchConfigsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	paramsTypeMap["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	paramsTypeMap["albAutoscalelaunchconfigId"] = bindings.NewStringType()
	pathParams["alb_autoscalelaunchconfig_id"] = "albAutoscalelaunchconfigId"
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
		"/policy/api/v1/infra/alb-auto-scale-launch-configs/{albAutoscalelaunchconfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albAutoScaleLaunchConfigsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albAutoScaleLaunchConfigsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigApiResponseBindingType)
}

func albAutoScaleLaunchConfigsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
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
		"/policy/api/v1/infra/alb-auto-scale-launch-configs",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albAutoScaleLaunchConfigsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["a_LB_auto_scale_launch_config"] = "ALBAutoScaleLaunchConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albAutoScaleLaunchConfigsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func albAutoScaleLaunchConfigsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["a_LB_auto_scale_launch_config"] = "ALBAutoScaleLaunchConfig"
	paramsTypeMap["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	paramsTypeMap["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	paramsTypeMap["albAutoscalelaunchconfigId"] = bindings.NewStringType()
	pathParams["alb_autoscalelaunchconfig_id"] = "albAutoscalelaunchconfigId"
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
		"a_LB_auto_scale_launch_config",
		"PATCH",
		"/policy/api/v1/infra/alb-auto-scale-launch-configs/{albAutoscalelaunchconfigId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func albAutoScaleLaunchConfigsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["a_LB_auto_scale_launch_config"] = "ALBAutoScaleLaunchConfig"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func albAutoScaleLaunchConfigsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
}

func albAutoScaleLaunchConfigsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	fields["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	fieldNameMap["alb_autoscalelaunchconfig_id"] = "AlbAutoscalelaunchconfigId"
	fieldNameMap["a_LB_auto_scale_launch_config"] = "ALBAutoScaleLaunchConfig"
	paramsTypeMap["a_LB_auto_scale_launch_config"] = bindings.NewReferenceType(model.ALBAutoScaleLaunchConfigBindingType)
	paramsTypeMap["alb_autoscalelaunchconfig_id"] = bindings.NewStringType()
	paramsTypeMap["albAutoscalelaunchconfigId"] = bindings.NewStringType()
	pathParams["alb_autoscalelaunchconfig_id"] = "albAutoscalelaunchconfigId"
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
		"a_LB_auto_scale_launch_config",
		"PUT",
		"/policy/api/v1/infra/alb-auto-scale-launch-configs/{albAutoscalelaunchconfigId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
