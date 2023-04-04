// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: FullSyncAction.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package infra

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"reflect"
)

// Possible value for ``action`` of method FullSyncAction#create.
const FullSyncAction_CREATE_ACTION_REQUEST_FULL_SYNC = "request_full_sync"

// Possible value for ``action`` of method FullSyncAction#create.
const FullSyncAction_CREATE_ACTION_REQUEST_NOTIFICATIONS_FULL_SYNC = "request_notifications_full_sync"

// Possible value for ``action`` of method FullSyncAction#create.
const FullSyncAction_CREATE_ACTION_ABORT_CURRENT_SYNC = "abort_current_sync"

// Possible value for ``action`` of method FullSyncAction#create.
const FullSyncAction_CREATE_ACTION_PURGE_HISTORY = "purge_history"

// Possible value for ``syncType`` of method FullSyncAction#create.
const FullSyncAction_CREATE_SYNC_TYPE_SYNC = "gm_to_lm_full_sync"

func fullSyncActionCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["action"] = vapiBindings_.NewStringType()
	fields["sync_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["action"] = "Action"
	fieldNameMap["sync_type"] = "SyncType"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func FullSyncActionCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func fullSyncActionCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["action"] = vapiBindings_.NewStringType()
	fields["sync_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["action"] = "Action"
	fieldNameMap["sync_type"] = "SyncType"
	paramsTypeMap["sync_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["action"] = vapiBindings_.NewStringType()
	queryParams["sync_type"] = "sync_type"
	queryParams["action"] = "action"
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
		"POST",
		"/policy/api/v1/infra/full-sync-action",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
