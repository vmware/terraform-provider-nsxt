// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Users.
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

func usersActivateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_password_property"] = "NodeUserPasswordProperty"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersActivateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersActivateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_password_property"] = "NodeUserPasswordProperty"
	paramsTypeMap["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"action=activate",
		"node_user_password_property",
		"POST",
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersCreateaudituserInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersCreateaudituserOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersCreateaudituserRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	paramsTypeMap["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
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
		"action=create_audit_user",
		"node_user_properties",
		"POST",
		"/api/v1/node/users",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersCreateuserInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersCreateuserOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersCreateuserRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	paramsTypeMap["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
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
		"action=create_user",
		"node_user_properties",
		"POST",
		"/api/v1/node/users",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersDeactivateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersDeactivateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersDeactivateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"action=deactivate",
		"",
		"POST",
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func usersDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fieldNameMap["userid"] = "Userid"
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesListResultBindingType)
}

func usersListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
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
		"/api/v1/node/users",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersResetownpasswordInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["reset_node_user_own_password_properties"] = vapiBindings_.NewReferenceType(nsxModel.ResetNodeUserOwnPasswordPropertiesBindingType)
	fieldNameMap["reset_node_user_own_password_properties"] = "ResetNodeUserOwnPasswordProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersResetownpasswordOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func usersResetownpasswordRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["reset_node_user_own_password_properties"] = vapiBindings_.NewReferenceType(nsxModel.ResetNodeUserOwnPasswordPropertiesBindingType)
	fieldNameMap["reset_node_user_own_password_properties"] = "ResetNodeUserOwnPasswordProperties"
	paramsTypeMap["reset_node_user_own_password_properties"] = vapiBindings_.NewReferenceType(nsxModel.ResetNodeUserOwnPasswordPropertiesBindingType)
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
		"action=reset_own_password",
		"reset_node_user_own_password_properties",
		"POST",
		"/api/v1/node/users",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersResetpasswordInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_password_property"] = "NodeUserPasswordProperty"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersResetpasswordOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func usersResetpasswordRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_password_property"] = "NodeUserPasswordProperty"
	paramsTypeMap["node_user_password_property"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPasswordPropertyBindingType)
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"action=reset_password",
		"node_user_password_property",
		"POST",
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func usersUpdateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func UsersUpdateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
}

func usersUpdateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["userid"] = vapiBindings_.NewStringType()
	fields["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	fieldNameMap["userid"] = "Userid"
	fieldNameMap["node_user_properties"] = "NodeUserProperties"
	paramsTypeMap["node_user_properties"] = vapiBindings_.NewReferenceType(nsxModel.NodeUserPropertiesBindingType)
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	paramsTypeMap["userid"] = vapiBindings_.NewStringType()
	pathParams["userid"] = "userid"
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
		"node_user_properties",
		"PUT",
		"/api/v1/node/users/{userid}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
