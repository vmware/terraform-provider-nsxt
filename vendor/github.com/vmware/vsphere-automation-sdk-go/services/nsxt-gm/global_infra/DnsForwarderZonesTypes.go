/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: DnsForwarderZones.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package global_infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func dnsForwarderZonesDeleteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dnsForwarderZonesDeleteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func dnsForwarderZonesDeleteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	paramsTypeMap["dns_forwarder_zone_id"] = bindings.NewStringType()
	paramsTypeMap["dnsForwarderZoneId"] = bindings.NewStringType()
	pathParams["dns_forwarder_zone_id"] = "dnsForwarderZoneId"
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
		"/global-manager/api/v1/global-infra/dns-forwarder-zones/{dnsForwarderZoneId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func dnsForwarderZonesGetInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dnsForwarderZonesGetOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
}

func dnsForwarderZonesGetRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	paramsTypeMap["dns_forwarder_zone_id"] = bindings.NewStringType()
	paramsTypeMap["dnsForwarderZoneId"] = bindings.NewStringType()
	pathParams["dns_forwarder_zone_id"] = "dnsForwarderZoneId"
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
		"/global-manager/api/v1/global-infra/dns-forwarder-zones/{dnsForwarderZoneId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func dnsForwarderZonesListInputType() bindings.StructType {
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

func dnsForwarderZonesListOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyDnsForwarderZoneListResultBindingType)
}

func dnsForwarderZonesListRestMetadata() protocol.OperationRestMetadata {
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
		"/global-manager/api/v1/global-infra/dns-forwarder-zones",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func dnsForwarderZonesPatchInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fields["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	fieldNameMap["policy_dns_forwarder_zone"] = "PolicyDnsForwarderZone"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dnsForwarderZonesPatchOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func dnsForwarderZonesPatchRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fields["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	fieldNameMap["policy_dns_forwarder_zone"] = "PolicyDnsForwarderZone"
	paramsTypeMap["dns_forwarder_zone_id"] = bindings.NewStringType()
	paramsTypeMap["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	paramsTypeMap["dnsForwarderZoneId"] = bindings.NewStringType()
	pathParams["dns_forwarder_zone_id"] = "dnsForwarderZoneId"
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
		"policy_dns_forwarder_zone",
		"PATCH",
		"/global-manager/api/v1/global-infra/dns-forwarder-zones/{dnsForwarderZoneId}",
		"",
		resultHeaders,
		204,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}

func dnsForwarderZonesUpdateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fields["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	fieldNameMap["policy_dns_forwarder_zone"] = "PolicyDnsForwarderZone"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func dnsForwarderZonesUpdateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
}

func dnsForwarderZonesUpdateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	dispatchHeaderParams := map[string]string{}
	bodyFieldsMap := map[string]string{}
	fields["dns_forwarder_zone_id"] = bindings.NewStringType()
	fields["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	fieldNameMap["dns_forwarder_zone_id"] = "DnsForwarderZoneId"
	fieldNameMap["policy_dns_forwarder_zone"] = "PolicyDnsForwarderZone"
	paramsTypeMap["dns_forwarder_zone_id"] = bindings.NewStringType()
	paramsTypeMap["policy_dns_forwarder_zone"] = bindings.NewReferenceType(model.PolicyDnsForwarderZoneBindingType)
	paramsTypeMap["dnsForwarderZoneId"] = bindings.NewStringType()
	pathParams["dns_forwarder_zone_id"] = "dnsForwarderZoneId"
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
		"policy_dns_forwarder_zone",
		"PUT",
		"/global-manager/api/v1/global-infra/dns-forwarder-zones/{dnsForwarderZoneId}",
		"",
		resultHeaders,
		200,
		"",
		errorHeaders,
		map[string]int{"com.vmware.vapi.std.errors.invalid_request": 400,"com.vmware.vapi.std.errors.unauthorized": 403,"com.vmware.vapi.std.errors.service_unavailable": 503,"com.vmware.vapi.std.errors.internal_server_error": 500,"com.vmware.vapi.std.errors.not_found": 404})
}


