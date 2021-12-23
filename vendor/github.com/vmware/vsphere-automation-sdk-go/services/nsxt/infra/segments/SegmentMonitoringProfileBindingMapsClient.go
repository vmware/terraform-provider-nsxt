// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: SegmentMonitoringProfileBindingMaps
// Used by client-side stubs.

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type SegmentMonitoringProfileBindingMapsClient interface {

	// API will delete Infra Segment Monitoring Profile Binding Profile.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) error

	// API will get Infra Segment Monitoring Profile Binding Map.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
	// @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) (model.SegmentMonitoringProfileBindingMap, error)

	// API will list all Infra Segment Monitoring Profile Binding Maps in current segment id.
	//
	// @param infraSegmentIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMapListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentMonitoringProfileBindingMapListResult, error)

	// API will create infra segment monitoring profile binding map.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
	// @param segmentMonitoringProfileBindingMapParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) error

	// API will update Infra Segment Monitoring Profile Binding Map.
	//
	// @param infraSegmentIdParam Infra Segment ID (required)
	// @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
	// @param segmentMonitoringProfileBindingMapParam (required)
	// @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) (model.SegmentMonitoringProfileBindingMap, error)
}

type segmentMonitoringProfileBindingMapsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewSegmentMonitoringProfileBindingMapsClient(connector client.Connector) *segmentMonitoringProfileBindingMapsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	sIface := segmentMonitoringProfileBindingMapsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &sIface
}

func (sIface *segmentMonitoringProfileBindingMapsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := sIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (sIface *segmentMonitoringProfileBindingMapsClient) Delete(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentMonitoringProfileBindingMapsDeleteInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMapId", segmentMonitoringProfileBindingMapIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentMonitoringProfileBindingMapsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps", "delete", inputDataValue, executionContext)
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

func (sIface *segmentMonitoringProfileBindingMapsClient) Get(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) (model.SegmentMonitoringProfileBindingMap, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentMonitoringProfileBindingMapsGetInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMapId", segmentMonitoringProfileBindingMapIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentMonitoringProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentMonitoringProfileBindingMapsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps", "get", inputDataValue, executionContext)
	var emptyOutput model.SegmentMonitoringProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentMonitoringProfileBindingMapsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentMonitoringProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *segmentMonitoringProfileBindingMapsClient) List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentMonitoringProfileBindingMapListResult, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentMonitoringProfileBindingMapsListInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentMonitoringProfileBindingMapListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentMonitoringProfileBindingMapsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps", "list", inputDataValue, executionContext)
	var emptyOutput model.SegmentMonitoringProfileBindingMapListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentMonitoringProfileBindingMapsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentMonitoringProfileBindingMapListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *segmentMonitoringProfileBindingMapsClient) Patch(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentMonitoringProfileBindingMapsPatchInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMapId", segmentMonitoringProfileBindingMapIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMap", segmentMonitoringProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentMonitoringProfileBindingMapsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps", "patch", inputDataValue, executionContext)
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

func (sIface *segmentMonitoringProfileBindingMapsClient) Update(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) (model.SegmentMonitoringProfileBindingMap, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(segmentMonitoringProfileBindingMapsUpdateInputType(), typeConverter)
	sv.AddStructField("InfraSegmentId", infraSegmentIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMapId", segmentMonitoringProfileBindingMapIdParam)
	sv.AddStructField("SegmentMonitoringProfileBindingMap", segmentMonitoringProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SegmentMonitoringProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := segmentMonitoringProfileBindingMapsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.segments.segment_monitoring_profile_binding_maps", "update", inputDataValue, executionContext)
	var emptyOutput model.SegmentMonitoringProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), segmentMonitoringProfileBindingMapsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SegmentMonitoringProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
