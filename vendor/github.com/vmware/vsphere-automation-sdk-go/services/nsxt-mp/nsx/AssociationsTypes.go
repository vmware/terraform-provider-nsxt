// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Data type definitions file for service: Associations.
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

// Possible value for ``associatedResourceType`` of method Associations#list.
const Associations_LIST_ASSOCIATED_RESOURCE_TYPE_NSGROUP = "NSGroup"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_NSGROUP = "NSGroup"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_IPSET = "IPSet"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_MACSET = "MACSet"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_LOGICALSWITCH = "LogicalSwitch"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_LOGICALPORT = "LogicalPort"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_VIRTUALMACHINE = "VirtualMachine"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_DIRECTORYGROUP = "DirectoryGroup"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_VIRTUALNETWORKINTERFACE = "VirtualNetworkInterface"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_TRANSPORTNODE = "TransportNode"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_IPADDRESS = "IPAddress"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_CLOUDNATIVESERVICEINSTANCE = "CloudNativeServiceInstance"

// Possible value for ``resourceType`` of method Associations#list.
const Associations_LIST_RESOURCE_TYPE_PHYSICALSERVER = "PhysicalServer"

func associationsListInputType() vapiBindings_.StructType {
	fields := make(map[string]vapiBindings_.BindingType)
	fieldNameMap := make(map[string]string)
	fields["associated_resource_type"] = vapiBindings_.NewStringType()
	fields["resource_id"] = vapiBindings_.NewStringType()
	fields["resource_type"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["fetch_ancestors"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["associated_resource_type"] = "AssociatedResourceType"
	fieldNameMap["resource_id"] = "ResourceId"
	fieldNameMap["resource_type"] = "ResourceType"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["fetch_ancestors"] = "FetchAncestors"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []vapiBindings_.Validator{}
	return vapiBindings_.NewStructType("operation-input", fields, reflect.TypeOf(vapiData_.StructValue{}), fieldNameMap, validators)
}

func AssociationsListOutputType() vapiBindings_.BindingType {
	return vapiBindings_.NewReferenceType(nsxModel.AssociationListResultBindingType)
}

func associationsListRestMetadata() vapiProtocol_.OperationRestMetadata {
	fields := map[string]vapiBindings_.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]vapiBindings_.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["associated_resource_type"] = vapiBindings_.NewStringType()
	fields["resource_id"] = vapiBindings_.NewStringType()
	fields["resource_type"] = vapiBindings_.NewStringType()
	fields["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["fetch_ancestors"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fields["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	fields["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	fields["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	fieldNameMap["associated_resource_type"] = "AssociatedResourceType"
	fieldNameMap["resource_id"] = "ResourceId"
	fieldNameMap["resource_type"] = "ResourceType"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["fetch_ancestors"] = "FetchAncestors"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["cursor"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["sort_ascending"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["included_fields"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["resource_type"] = vapiBindings_.NewStringType()
	paramsTypeMap["resource_id"] = vapiBindings_.NewStringType()
	paramsTypeMap["sort_by"] = vapiBindings_.NewOptionalType(vapiBindings_.NewStringType())
	paramsTypeMap["associated_resource_type"] = vapiBindings_.NewStringType()
	paramsTypeMap["fetch_ancestors"] = vapiBindings_.NewOptionalType(vapiBindings_.NewBooleanType())
	paramsTypeMap["page_size"] = vapiBindings_.NewOptionalType(vapiBindings_.NewIntegerType())
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["resource_type"] = "resource_type"
	queryParams["resource_id"] = "resource_id"
	queryParams["sort_by"] = "sort_by"
	queryParams["associated_resource_type"] = "associated_resource_type"
	queryParams["fetch_ancestors"] = "fetch_ancestors"
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
		"/api/v1/associations",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400, "com.vmware.vapi.std.errors.unauthorized": 403, "com.vmware.vapi.std.errors.service_unavailable": 503, "com.vmware.vapi.std.errors.internal_server_error": 500, "com.vmware.vapi.std.errors.not_found": 404})
}
