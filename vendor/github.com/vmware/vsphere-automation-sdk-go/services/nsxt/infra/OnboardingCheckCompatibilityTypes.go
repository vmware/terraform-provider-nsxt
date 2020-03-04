/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: OnboardingCheckCompatibility.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package infra

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func onboardingCheckCompatibilityCreateInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["site_node_connection_info"] = bindings.NewReferenceType(model.SiteNodeConnectionInfoBindingType)
	fieldNameMap["site_node_connection_info"] = "SiteNodeConnectionInfo"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func onboardingCheckCompatibilityCreateOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.CompatibilityCheckResultBindingType)
}

func onboardingCheckCompatibilityCreateRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["site_node_connection_info"] = bindings.NewReferenceType(model.SiteNodeConnectionInfoBindingType)
	fieldNameMap["site_node_connection_info"] = "SiteNodeConnectionInfo"
	paramsTypeMap["site_node_connection_info"] = bindings.NewReferenceType(model.SiteNodeConnectionInfoBindingType)
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
		"site_node_connection_info",
		"POST",
		"/policy/api/v1/infra/onboarding-check-compatibility",
		resultHeaders,
		201,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


