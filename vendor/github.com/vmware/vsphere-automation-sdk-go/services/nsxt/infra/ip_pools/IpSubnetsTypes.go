/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: IpSubnets.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package ip_pools

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func ipSubnetsDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipSubnetsDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipSubnetsDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	paramsTypeMap["ip_pool_id"] = bindings.NewStringType()
	paramsTypeMap["ip_subnet_id"] = bindings.NewStringType()
	paramsTypeMap["ipPoolId"] = bindings.NewStringType()
	paramsTypeMap["ipSubnetId"] = bindings.NewStringType()
	pathParams["ip_pool_id"] = "ipPoolId"
	pathParams["ip_subnet_id"] = "ipSubnetId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"/policy/api/v1/infra/ip-pools/{ipPoolId}/ip-subnets/{ipSubnetId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipSubnetsGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipSubnetsGetOutputType() bindings.BindingType {
	return bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
}

func ipSubnetsGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	paramsTypeMap["ip_pool_id"] = bindings.NewStringType()
	paramsTypeMap["ip_subnet_id"] = bindings.NewStringType()
	paramsTypeMap["ipPoolId"] = bindings.NewStringType()
	paramsTypeMap["ipSubnetId"] = bindings.NewStringType()
	pathParams["ip_pool_id"] = "ipPoolId"
	pathParams["ip_subnet_id"] = "ipSubnetId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"/policy/api/v1/infra/ip-pools/{ipPoolId}/ip-subnets/{ipSubnetId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipSubnetsListInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipSubnetsListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.IpAddressPoolSubnetListResultBindingType)
}

func ipSubnetsListRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fields["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	fields["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["cursor"] = "Cursor"
	fieldNameMap["include_mark_for_delete_objects"] = "IncludeMarkForDeleteObjects"
	fieldNameMap["included_fields"] = "IncludedFields"
	fieldNameMap["page_size"] = "PageSize"
	fieldNameMap["sort_ascending"] = "SortAscending"
	fieldNameMap["sort_by"] = "SortBy"
	paramsTypeMap["ip_pool_id"] = bindings.NewStringType()
	paramsTypeMap["included_fields"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["page_size"] = bindings.NewOptionalType(bindings.NewIntegerType())
	paramsTypeMap["include_mark_for_delete_objects"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["cursor"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_by"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["sort_ascending"] = bindings.NewOptionalType(bindings.NewBooleanType())
	paramsTypeMap["ipPoolId"] = bindings.NewStringType()
	pathParams["ip_pool_id"] = "ipPoolId"
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
		dispatchHeaderParams,
		bodyFieldsMap,
		"",
		"",
		"GET",
		"/policy/api/v1/infra/ip-pools/{ipPoolId}/ip-subnets",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipSubnetsPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fields["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	fieldNameMap["ip_address_pool_subnet"] = "IpAddressPoolSubnet"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipSubnetsPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func ipSubnetsPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fields["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	fieldNameMap["ip_address_pool_subnet"] = "IpAddressPoolSubnet"
	paramsTypeMap["ip_pool_id"] = bindings.NewStringType()
	paramsTypeMap["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	paramsTypeMap["ip_subnet_id"] = bindings.NewStringType()
	paramsTypeMap["ipPoolId"] = bindings.NewStringType()
	paramsTypeMap["ipSubnetId"] = bindings.NewStringType()
	pathParams["ip_pool_id"] = "ipPoolId"
	pathParams["ip_subnet_id"] = "ipSubnetId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"ip_address_pool_subnet",
		"PATCH",
		"/policy/api/v1/infra/ip-pools/{ipPoolId}/ip-subnets/{ipSubnetId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func ipSubnetsUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fields["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	fieldNameMap["ip_address_pool_subnet"] = "IpAddressPoolSubnet"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func ipSubnetsUpdateOutputType() bindings.BindingType {
	return bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
}

func ipSubnetsUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["ip_pool_id"] = bindings.NewStringType()
	fields["ip_subnet_id"] = bindings.NewStringType()
	fields["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	fieldNameMap["ip_pool_id"] = "IpPoolId"
	fieldNameMap["ip_subnet_id"] = "IpSubnetId"
	fieldNameMap["ip_address_pool_subnet"] = "IpAddressPoolSubnet"
	paramsTypeMap["ip_pool_id"] = bindings.NewStringType()
	paramsTypeMap["ip_address_pool_subnet"] = bindings.NewDynamicStructType([]bindings.ReferenceType{bindings.NewReferenceType(model.IpAddressPoolSubnetBindingType),}, bindings.REST)
	paramsTypeMap["ip_subnet_id"] = bindings.NewStringType()
	paramsTypeMap["ipPoolId"] = bindings.NewStringType()
	paramsTypeMap["ipSubnetId"] = bindings.NewStringType()
	pathParams["ip_pool_id"] = "ipPoolId"
	pathParams["ip_subnet_id"] = "ipSubnetId"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
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
		"ip_address_pool_subnet",
		"PUT",
		"/policy/api/v1/infra/ip-pools/{ipPoolId}/ip-subnets/{ipSubnetId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


