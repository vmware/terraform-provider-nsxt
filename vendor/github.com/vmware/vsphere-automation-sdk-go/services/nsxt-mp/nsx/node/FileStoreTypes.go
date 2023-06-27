// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: FileStore.
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

func fileStoreCopyfromremotefileInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["file_name"] = vapiBindings_.NewStringType()
	fields["copy_from_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyFromRemoteFilePropertiesBindingType)
	fieldNameMap["file_name"] = "FileName"
	fieldNameMap["copy_from_remote_file_properties"] = "CopyFromRemoteFileProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreCopyfromremotefileOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.FilePropertiesBindingType)
}

func fileStoreCopyfromremotefileRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["file_name"] = vapiBindings_.NewStringType()
	fields["copy_from_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyFromRemoteFilePropertiesBindingType)
	fieldNameMap["file_name"] = "FileName"
	fieldNameMap["copy_from_remote_file_properties"] = "CopyFromRemoteFileProperties"
	paramsTypeMap["copy_from_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyFromRemoteFilePropertiesBindingType)
	paramsTypeMap["file_name"] = vapiBindings_.NewStringType()
	paramsTypeMap["fileName"] = vapiBindings_.NewStringType()
	pathParams["file_name"] = "fileName"
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
		"action=copy_from_remote_file",
		"copy_from_remote_file_properties",
		"POST",
		"/api/v1/node/file-store/{fileName}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.concurrent_change": 409, "com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreCopytoremotefileInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["file_name"] = vapiBindings_.NewStringType()
	fields["copy_to_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyToRemoteFilePropertiesBindingType)
	fieldNameMap["file_name"] = "FileName"
	fieldNameMap["copy_to_remote_file_properties"] = "CopyToRemoteFileProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreCopytoremotefileOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func fileStoreCopytoremotefileRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["file_name"] = vapiBindings_.NewStringType()
	fields["copy_to_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyToRemoteFilePropertiesBindingType)
	fieldNameMap["file_name"] = "FileName"
	fieldNameMap["copy_to_remote_file_properties"] = "CopyToRemoteFileProperties"
	paramsTypeMap["file_name"] = vapiBindings_.NewStringType()
	paramsTypeMap["copy_to_remote_file_properties"] = vapiBindings_.NewReferenceType(nsxModel.CopyToRemoteFilePropertiesBindingType)
	paramsTypeMap["fileName"] = vapiBindings_.NewStringType()
	pathParams["file_name"] = "fileName"
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
		"action=copy_to_remote_file",
		"copy_to_remote_file_properties",
		"POST",
		"/api/v1/node/file-store/{fileName}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreCreateremotedirectoryInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["create_remote_directory_properties"] = vapiBindings_.NewReferenceType(nsxModel.CreateRemoteDirectoryPropertiesBindingType)
	fieldNameMap["create_remote_directory_properties"] = "CreateRemoteDirectoryProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreCreateremotedirectoryOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func fileStoreCreateremotedirectoryRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["create_remote_directory_properties"] = vapiBindings_.NewReferenceType(nsxModel.CreateRemoteDirectoryPropertiesBindingType)
	fieldNameMap["create_remote_directory_properties"] = "CreateRemoteDirectoryProperties"
	paramsTypeMap["create_remote_directory_properties"] = vapiBindings_.NewReferenceType(nsxModel.CreateRemoteDirectoryPropertiesBindingType)
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
		"action=create_remote_directory",
		"create_remote_directory_properties",
		"POST",
		"/api/v1/node/file-store",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.timed_out": 500, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreDeleteInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["file_name"] = vapiBindings_.NewStringType()
	fieldNameMap["file_name"] = "FileName"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreDeleteOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func fileStoreDeleteRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["file_name"] = vapiBindings_.NewStringType()
	fieldNameMap["file_name"] = "FileName"
	paramsTypeMap["file_name"] = vapiBindings_.NewStringType()
	paramsTypeMap["fileName"] = vapiBindings_.NewStringType()
	pathParams["file_name"] = "fileName"
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
		"/api/v1/node/file-store/{fileName}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["file_name"] = vapiBindings_.NewStringType()
	fieldNameMap["file_name"] = "FileName"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.FilePropertiesBindingType)
}

func fileStoreGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["file_name"] = vapiBindings_.NewStringType()
	fieldNameMap["file_name"] = "FileName"
	paramsTypeMap["file_name"] = vapiBindings_.NewStringType()
	paramsTypeMap["fileName"] = vapiBindings_.NewStringType()
	pathParams["file_name"] = "fileName"
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
		"/api/v1/node/file-store/{fileName}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.FilePropertiesListResultBindingType)
}

func fileStoreListRestMetadata() vapiProtocol_.OperationRestMetadata {
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
		"/api/v1/node/file-store",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func fileStoreRetrievesshfingerprintInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ssh_fingerprint_properties"] = vapiBindings_.NewReferenceType(nsxModel.SshFingerprintPropertiesBindingType)
	fieldNameMap["ssh_fingerprint_properties"] = "SshFingerprintProperties"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FileStoreRetrievesshfingerprintOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.SshFingerprintPropertiesBindingType)
}

func fileStoreRetrievesshfingerprintRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ssh_fingerprint_properties"] = vapiBindings_.NewReferenceType(nsxModel.SshFingerprintPropertiesBindingType)
	fieldNameMap["ssh_fingerprint_properties"] = "SshFingerprintProperties"
	paramsTypeMap["ssh_fingerprint_properties"] = vapiBindings_.NewReferenceType(nsxModel.SshFingerprintPropertiesBindingType)
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
		"action=retrieve_ssh_fingerprint",
		"ssh_fingerprint_properties",
		"POST",
		"/api/v1/node/file-store",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
