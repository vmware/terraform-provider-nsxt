/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Tier0DeploymentMaps.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package locale_services

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func tier0DeploymentMapsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tier0DeploymentMapsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func tier0DeploymentMapsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["tier0DeploymentMapId"] = bindings.NewStringType()
	pathParams["tier0_deployment_map_id"] = "tier0DeploymentMapId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/tier-0-deployment-maps/{tier0DeploymentMapId}",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func tier0DeploymentMapsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tier0DeploymentMapsGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
}

func tier0DeploymentMapsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["tier0DeploymentMapId"] = bindings.NewStringType()
	pathParams["tier0_deployment_map_id"] = "tier0DeploymentMapId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/tier-0-deployment-maps/{tier0DeploymentMapId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func tier0DeploymentMapsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tier0DeploymentMapsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier0DeploymentMapListResultBindingType)
}

func tier0DeploymentMapsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/tier-0-deployment-maps",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func tier0DeploymentMapsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fields["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	fieldNameMap["tier0_deployment_map"] = "Tier0DeploymentMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tier0DeploymentMapsPatchOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
}

func tier0DeploymentMapsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fields["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	fieldNameMap["tier0_deployment_map"] = "Tier0DeploymentMap"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["tier0DeploymentMapId"] = bindings.NewStringType()
	pathParams["tier0_deployment_map_id"] = "tier0DeploymentMapId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"tier0_deployment_map",
		"PATCH",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/tier-0-deployment-maps/{tier0DeploymentMapId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func tier0DeploymentMapsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fields["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	fieldNameMap["tier0_deployment_map"] = "Tier0DeploymentMap"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func tier0DeploymentMapsUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
}

func tier0DeploymentMapsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["tier0_id"] = bindings.NewStringType()
	fields["locale_service_id"] = bindings.NewStringType()
	fields["tier0_deployment_map_id"] = bindings.NewStringType()
	fields["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	fieldNameMap["tier0_id"] = "Tier0Id"
	fieldNameMap["locale_service_id"] = "LocaleServiceId"
	fieldNameMap["tier0_deployment_map_id"] = "Tier0DeploymentMapId"
	fieldNameMap["tier0_deployment_map"] = "Tier0DeploymentMap"
	paramsTypeMap["tier0_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map_id"] = bindings.NewStringType()
	paramsTypeMap["locale_service_id"] = bindings.NewStringType()
	paramsTypeMap["tier0_deployment_map"] = bindings.NewReferenceType(model.Tier0DeploymentMapBindingType)
	paramsTypeMap["tier0Id"] = bindings.NewStringType()
	paramsTypeMap["localeServiceId"] = bindings.NewStringType()
	paramsTypeMap["tier0DeploymentMapId"] = bindings.NewStringType()
	pathParams["tier0_deployment_map_id"] = "tier0DeploymentMapId"
	pathParams["tier0_id"] = "tier0Id"
	pathParams["locale_service_id"] = "localeServiceId"
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
		"tier0_deployment_map",
		"PUT",
		"/policy/api/v1/infra/tier-0s/{tier0Id}/locale-services/{localeServiceId}/tier-0-deployment-maps/{tier0DeploymentMapId}",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


