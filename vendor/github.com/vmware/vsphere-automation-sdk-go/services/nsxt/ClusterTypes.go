/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: Cluster.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package nsx_policy

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)

// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_GLOBAL_MANAGER = "GLOBAL_MANAGER"
// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_LOCAL_MANAGER = "LOCAL_MANAGER"
// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_LOCAL_LOCAL_MANAGER = "LOCAL_LOCAL_MANAGER"
// Possible value for ``frameType`` of method Cluster#backuptoremote.
const Cluster_BACKUPTOREMOTE_FRAME_TYPE_NSX_INTELLIGENCE = "NSX_INTELLIGENCE"




func clusterBackuptoremoteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["frame_type"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["site_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["frame_type"] = "FrameType"
	fieldNameMap["site_id"] = "SiteId"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func clusterBackuptoremoteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func clusterBackuptoremoteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["frame_type"] = bindings.NewOptionalType(bindings.NewStringType())
	fields["site_id"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["frame_type"] = "FrameType"
	fieldNameMap["site_id"] = "SiteId"
	paramsTypeMap["frame_type"] = bindings.NewOptionalType(bindings.NewStringType())
	paramsTypeMap["site_id"] = bindings.NewOptionalType(bindings.NewStringType())
	queryParams["site_id"] = "site_id"
	queryParams["frame_type"] = "frame_type"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"action=backup_to_remote",
		"",
		"POST",
		"/policy/api/v1/cluster",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}

func clusterSummarizeinventorytoremoteInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func clusterSummarizeinventorytoremoteOutputType() bindings.BindingType {
	return bindings.NewVoidType()
}

func clusterSummarizeinventorytoremoteRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"action=summarize_inventory_to_remote",
		"",
		"POST",
		"/policy/api/v1/cluster",
		resultHeaders,
		204,
		errorHeaders,
		map[string]int{"InvalidRequest": 400,"Unauthorized": 403,"ServiceUnavailable": 503,"InternalServerError": 500,"NotFound": 404})
}


