// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: ComputeCollections.
// Includes binding types of a structures and enumerations defined in the service.
// Shared by client-side stubs and server-side skeletons to ensure type
// compatibility.

package fabric

import (
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocol_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"reflect"
)

// Possible value for ``action`` of method ComputeCollections#create.
const ComputeCollections_CREATE_ACTION_NSX = "remove_nsx"

func computeCollectionsCreateInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cc_ext_id"] = vapiBindings_.NewStringType()
	fields["action"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cc_ext_id"] = "CcExtId"
	fieldNameMap["action"] = "Action"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ComputeCollectionsCreateOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewVoidType()
}

func computeCollectionsCreateRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cc_ext_id"] = vapiBindings_.NewStringType()
	fields["action"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cc_ext_id"] = "CcExtId"
	fieldNameMap["action"] = "Action"
	paramsTypeMap["action"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["cc_ext_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["ccExtId"] = vapiBindings_.NewStringType()
	pathParams["cc_ext_id"] = "ccExtId"
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
		"/api/v1/fabric/compute-collections/{ccExtId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func computeCollectionsGetInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cc_ext_id"] = vapiBindings_.NewStringType()
	fieldNameMap["cc_ext_id"] = "CcExtId"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ComputeCollectionsGetOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ComputeCollectionBindingType)
}

func computeCollectionsGetRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cc_ext_id"] = vapiBindings_.NewStringType()
	fieldNameMap["cc_ext_id"] = "CcExtId"
	paramsTypeMap["cc_ext_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["ccExtId"] = vapiBindings_.NewStringType()
	pathParams["cc_ext_id"] = "ccExtId"
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
		"/api/v1/fabric/compute-collections/{ccExtId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}

func computeCollectionsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["discovered_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["owner_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cm_local_id"] = "CmLocalId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["discovered_node_id"] = "DiscoveredNodeId"
	fieldNameMap["display_name"] = "DisplayName"
	fieldNameMap["external_id"] = "ExternalId"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["origin_id"] = "OriginId"
	fieldNameMap["origin_type"] = "OriginType"
	fieldNameMap["owner_id"] = "OwnerId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func ComputeCollectionsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.ComputeCollectionListResultBindingType)
}

func computeCollectionsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["discovered_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["origin_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["owner_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["cm_local_id"] = "CmLocalId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["discovered_node_id"] = "DiscoveredNodeId"
	fieldNameMap["display_name"] = "DisplayName"
	fieldNameMap["external_id"] = "ExternalId"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["node_id"] = "NodeId"
	fieldNameMap["origin_id"] = "OriginId"
	fieldNameMap["origin_type"] = "OriginType"
	fieldNameMap["owner_id"] = "OwnerId"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["cm_local_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["owner_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["external_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["origin_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["display_name"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["origin_type"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["discovered_node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["node_id"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["cm_local_id"] = "cm_local_id"
	queryParams["owner_id"] = "owner_id"
	queryParams["external_id"] = "external_id"
	queryParams["origin_id"] = "origin_id"
	queryParams["sort_by"] = "sort_by"
	queryParams["display_name"] = "display_name"
	queryParams["origin_type"] = "origin_type"
	queryParams["discovered_node_id"] = "discovered_node_id"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["node_id"] = "node_id"
	queryParams["page_size"] = "page_size"
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
		"/api/v1/fabric/compute-collections",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
