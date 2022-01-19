// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: SegmentDiscoveryProfileBindingMaps
// Used by client-side stubs.

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

const _ = core.SupportedByRuntimeVersion1

type SegmentDiscoveryProfileBindingMapsClient interface {

	// API will delete Segment Discovery Profile Binding Profile
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) error

	// API will get Infra Segment Discovery Profile Binding Map
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
	// @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) (model.SegmentDiscoveryProfileBindingMap, error)

	// API will list all Infra Segment Discovery Profile Binding Maps in current segment id.
	//
	// @param infraSegmentIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMapListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentDiscoveryProfileBindingMapListResult, error)

	// API will create Infra Segment Discovery Profile Binding Map. For objects with no binding maps, default profile is applied.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
	// @param segmentDiscoveryProfileBindingMapParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) error

	// API will update Infra Segment Discovery Profile Binding Map. For objects with no binding maps, default profile is applied.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
	// @param segmentDiscoveryProfileBindingMapParam (required)
	// @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) (model.SegmentDiscoveryProfileBindingMap, error)
}

type segmentDiscoveryProfileBindingMapsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewSegmentDiscoveryProfileBindingMapsClient(connector client.Connector) *segmentDiscoveryProfileBindingMapsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	sIface := segmentDiscoveryProfileBindingMapsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &sIface
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := sIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) Delete(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentDiscoveryProfileBindingMapsDeleteInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMapId", segmentDiscoveryProfileBindingMapIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentDiscoveryProfileBindingMapsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) Get(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) (model.SegmentDiscoveryProfileBindingMap, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentDiscoveryProfileBindingMapsGetInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMapId", segmentDiscoveryProfileBindingMapIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentDiscoveryProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentDiscoveryProfileBindingMapsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps", "get", inputDataValue, executionContext)
	var emptyOutput model.SegmentDiscoveryProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentDiscoveryProfileBindingMapsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentDiscoveryProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentDiscoveryProfileBindingMapListResult, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentDiscoveryProfileBindingMapsListInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentDiscoveryProfileBindingMapListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentDiscoveryProfileBindingMapsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps", "list", inputDataValue, executionContext)
	var emptyOutput model.SegmentDiscoveryProfileBindingMapListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentDiscoveryProfileBindingMapsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentDiscoveryProfileBindingMapListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) Patch(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentDiscoveryProfileBindingMapsPatchInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMapId", segmentDiscoveryProfileBindingMapIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMap", segmentDiscoveryProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentDiscoveryProfileBindingMapsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (sIface *segmentDiscoveryProfileBindingMapsClient) Update(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) (model.SegmentDiscoveryProfileBindingMap, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentDiscoveryProfileBindingMapsUpdateInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMapId", segmentDiscoveryProfileBindingMapIdParam)
	sv.AddStructField("SegmentDiscoveryProfileBindingMap", segmentDiscoveryProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentDiscoveryProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentDiscoveryProfileBindingMapsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.segments.segment_discovery_profile_binding_maps", "update", inputDataValue, executionContext)
	var emptyOutput model.SegmentDiscoveryProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentDiscoveryProfileBindingMapsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentDiscoveryProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
