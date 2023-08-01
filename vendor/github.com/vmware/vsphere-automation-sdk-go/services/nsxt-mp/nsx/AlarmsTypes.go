// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Alarms.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package nsx

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"reflect"
)

// Possible value for ``newStatus`` of method Alarms#setstatus.
const Alarms_SETSTATUS_NEW_STATUS_OPEN = "OPEN"

// Possible value for ``newStatus`` of method Alarms#setstatus.
const Alarms_SETSTATUS_NEW_STATUS_ACKNOWLEDGED = "ACKNOWLEDGED"

// Possible value for ``newStatus`` of method Alarms#setstatus.
const Alarms_SETSTATUS_NEW_STATUS_SUPPRESSED = "SUPPRESSED"

// Possible value for ``newStatus`` of method Alarms#setstatus.
const Alarms_SETSTATUS_NEW_STATUS_RESOLVED = "RESOLVED"

// Possible value for ``newStatus`` of method Alarms#setstatus0.
const Alarms_SETSTATUS_0_NEW_STATUS_OPEN = "OPEN"

// Possible value for ``newStatus`` of method Alarms#setstatus0.
const Alarms_SETSTATUS_0_NEW_STATUS_ACKNOWLEDGED = "ACKNOWLEDGED"

// Possible value for ``newStatus`` of method Alarms#setstatus0.
const Alarms_SETSTATUS_0_NEW_STATUS_SUPPRESSED = "SUPPRESSED"

// Possible value for ``newStatus`` of method Alarms#setstatus0.
const Alarms_SETSTATUS_0_NEW_STATUS_RESOLVED = "RESOLVED"

func alarmsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alarm_id"] = vapiBindings_.NewStringType()
	fieldNameMap["alarm_id"] = "AlarmId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlarmsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.AlarmBindingType)
}

func alarmsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alarm_id"] = vapiBindings_.NewStringType()
	fieldNameMap["alarm_id"] = "AlarmId"
	paramsTypeMap["alarm_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["alarmId"] = vapiBindings_.NewStringType()
	pathParams["alarm_id"] = "alarmId"
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
		"/api/v1/alarms/{alarmId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func alarmsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["after"] = "After"
	fieldNameMap["before"] = "Before"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["event_tag"] = "EventTag"
	fieldNameMap["event_type"] = "EventType"
	fieldNameMap["feature_name"] = "FeatureName"
	fieldNameMap["id"] = "Id"
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_resource_type"] = "NodeResourceType"
	fieldNameMap["org"] = "Org"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["project"] = "Project"
	fieldNameMap["severity"] = "Severity"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["status"] = "Status"
	fieldNameMap["vpc"] = "Vpc"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlarmsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.AlarmsListResultBindingType)
}

func alarmsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["after"] = "After"
	fieldNameMap["before"] = "Before"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["event_tag"] = "EventTag"
	fieldNameMap["event_type"] = "EventType"
	fieldNameMap["feature_name"] = "FeatureName"
	fieldNameMap["id"] = "Id"
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_resource_type"] = "NodeResourceType"
	fieldNameMap["org"] = "Org"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["project"] = "Project"
	fieldNameMap["severity"] = "Severity"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["status"] = "Status"
	fieldNameMap["vpc"] = "Vpc"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["cursor"] = "cursor"
	queryParams["severity"] = "severity"
	queryParams["feature_name"] = "feature_name"
	queryParams["before"] = "before"
	queryParams["org"] = "org"
	queryParams["vpc"] = "vpc"
	queryParams["project"] = "project"
	queryParams["intent_path"] = "intent_path"
	queryParams["sort_by"] = "sort_by"
	queryParams["event_tag"] = "event_tag"
	queryParams["event_type"] = "event_type"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["after"] = "after"
	queryParams["id"] = "id"
	queryParams["node_resource_type"] = "node_resource_type"
	queryParams["node_id"] = "node_id"
	queryParams["page_size"] = "page_size"
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
		"/api/v1/alarms",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func alarmsSetstatusInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["alarm_id"] = vapiBindings_.NewStringType()
	fields["new_status"] = vapiBindings_.NewStringType()
	fields["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fieldNameMap["alarm_id"] = "AlarmId"
	fieldNameMap["new_status"] = "NewStatus"
	fieldNameMap["suppress_duration"] = "SuppressDuration"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlarmsSetstatusOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.AlarmBindingType)
}

func alarmsSetstatusRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["alarm_id"] = vapiBindings_.NewStringType()
	fields["new_status"] = vapiBindings_.NewStringType()
	fields["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fieldNameMap["alarm_id"] = "AlarmId"
	fieldNameMap["new_status"] = "NewStatus"
	fieldNameMap["suppress_duration"] = "SuppressDuration"
	paramsTypeMap["new_status"] = vapiBindings_.NewStringType()
	paramsTypeMap["alarm_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["alarmId"] = vapiBindings_.NewStringType()
	pathParams["alarm_id"] = "alarmId"
	queryParams["new_status"] = "new_status"
	queryParams["suppress_duration"] = "suppress_duration"
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
		"action=set_status",
		"",
		"POST",
		"/api/v1/alarms/{alarmId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func alarmsSetstatus0InputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["new_status"] = vapiBindings_.NewStringType()
	fields["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["new_status"] = "NewStatus"
	fieldNameMap["after"] = "After"
	fieldNameMap["before"] = "Before"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["event_tag"] = "EventTag"
	fieldNameMap["event_type"] = "EventType"
	fieldNameMap["feature_name"] = "FeatureName"
	fieldNameMap["id"] = "Id"
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_resource_type"] = "NodeResourceType"
	fieldNameMap["org"] = "Org"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["project"] = "Project"
	fieldNameMap["severity"] = "Severity"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["status"] = "Status"
	fieldNameMap["suppress_duration"] = "SuppressDuration"
	fieldNameMap["vpc"] = "Vpc"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AlarmsSetstatus0OutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func alarmsSetstatus0RestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["new_status"] = vapiBindings_.NewStringType()
	fields["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["new_status"] = "NewStatus"
	fieldNameMap["after"] = "After"
	fieldNameMap["before"] = "Before"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["event_tag"] = "EventTag"
	fieldNameMap["event_type"] = "EventType"
	fieldNameMap["feature_name"] = "FeatureName"
	fieldNameMap["id"] = "Id"
	fieldNameMap["intent_path"] = "IntentPath"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["node_resource_type"] = "NodeResourceType"
	fieldNameMap["org"] = "Org"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["project"] = "Project"
	fieldNameMap["severity"] = "Severity"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	fieldNameMap["status"] = "Status"
	fieldNameMap["suppress_duration"] = "SuppressDuration"
	fieldNameMap["vpc"] = "Vpc"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["severity"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["feature_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["before"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["org"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["vpc"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["new_status"] = vapiBindings_.NewStringType()
	paramsTypeMap["project"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["intent_path"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["suppress_duration"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["event_tag"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["event_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["after"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_resource_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	paramsTypeMap["status"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	queryParams["cursor"] = "cursor"
	queryParams["severity"] = "severity"
	queryParams["feature_name"] = "feature_name"
	queryParams["before"] = "before"
	queryParams["org"] = "org"
	queryParams["vpc"] = "vpc"
	queryParams["new_status"] = "new_status"
	queryParams["project"] = "project"
	queryParams["intent_path"] = "intent_path"
	queryParams["sort_by"] = "sort_by"
	queryParams["suppress_duration"] = "suppress_duration"
	queryParams["event_tag"] = "event_tag"
	queryParams["event_type"] = "event_type"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["after"] = "after"
	queryParams["id"] = "id"
	queryParams["node_resource_type"] = "node_resource_type"
	queryParams["node_id"] = "node_id"
	queryParams["page_size"] = "page_size"
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
		"action=set_status",
		"",
		"POST",
		"/api/v1/alarms",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
