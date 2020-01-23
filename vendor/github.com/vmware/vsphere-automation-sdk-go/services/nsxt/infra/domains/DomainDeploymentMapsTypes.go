/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: DomainDeploymentMaps.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package domains

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func domainDeploymentMapsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func domainDeploymentMapsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func domainDeploymentMapsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["domainDeploymentMapId"] = bindings.NewStringType()
	pathParams["domain_deployment_map_id"] = "domainDeploymentMapId"
	pathParams["domain_id"] = "domainId"
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
		"DELETE",
		"/policy/api/v1/infra/domains/{domainId}/domain-deployment-maps/{domainDeploymentMapId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func domainDeploymentMapsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func domainDeploymentMapsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
}

func domainDeploymentMapsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["domainDeploymentMapId"] = bindings.NewStringType()
	pathParams["domain_deployment_map_id"] = "domainDeploymentMapId"
	pathParams["domain_id"] = "domainId"
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
		"/policy/api/v1/infra/domains/{domainId}/domain-deployment-maps/{domainDeploymentMapId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func domainDeploymentMapsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func domainDeploymentMapsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.DomainDeploymentMapListResultBindingType)
}

func domainDeploymentMapsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["domainId"] = bindings.NewStringType()
	pathParams["domain_id"] = "domainId"
	queryParams["cursor"] = "cursor"
	queryParams["sort_ascending"] = "sort_ascending"
	queryParams["included_fields"] = "included_fields"
	queryParams["sort_by"] = "sort_by"
	queryParams["include_mark_for_delete_objects"] = "include_mark_for_delete_objects"
	queryParams["page_size"] = "page_size"
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
		"/policy/api/v1/infra/domains/{domainId}/domain-deployment-maps",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func domainDeploymentMapsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fields["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	fieldNameMap["domain_deployment_map"] = "DomainDeploymentMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func domainDeploymentMapsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func domainDeploymentMapsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fields["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	fieldNameMap["domain_deployment_map"] = "DomainDeploymentMap"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["domainDeploymentMapId"] = bindings.NewStringType()
	pathParams["domain_deployment_map_id"] = "domainDeploymentMapId"
	pathParams["domain_id"] = "domainId"
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
		"domain_deployment_map",
		"PATCH",
		"/policy/api/v1/infra/domains/{domainId}/domain-deployment-maps/{domainDeploymentMapId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func domainDeploymentMapsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fields["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	fieldNameMap["domain_deployment_map"] = "DomainDeploymentMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func domainDeploymentMapsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
}

func domainDeploymentMapsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["domain_id"] = bindings.NewStringType()
	fields["domain_deployment_map_id"] = bindings.NewStringType()
	fields["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	fieldNameMap["domain_id"] = "DomainId"
	fieldNameMap["domain_deployment_map_id"] = "DomainDeploymentMapId"
	fieldNameMap["domain_deployment_map"] = "DomainDeploymentMap"
	paramsTypeMap["domain_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["domain_deployment_map"] = bindings.NewReferenceType(model.DomainDeploymentMapBindingType)
	paramsTypeMap["domainId"] = bindings.NewStringType()
	paramsTypeMap["domainDeploymentMapId"] = bindings.NewStringType()
	pathParams["domain_deployment_map_id"] = "domainDeploymentMapId"
	pathParams["domain_id"] = "domainId"
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
		"domain_deployment_map",
		"PUT",
		"/policy/api/v1/infra/domains/{domainId}/domain-deployment-maps/{domainDeploymentMapId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


