// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: LogicalSwitches
// Used by client-side stubs.

package nsx

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type LogicalSwitchesClient interface {

	// Creates a new logical switch. The request must include the transport_zone_id, display_name, and admin_state (UP or DOWN). The replication_mode (MTEP or SOURCE) is required for overlay logical switches, but not for VLAN-based logical switches. A vlan needs to be provided for VLAN-based logical switches.
	//  This api is now deprecated. Please use new api -/infra/segments/<segment-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalSwitchParam (required)
	// @return com.vmware.nsx.model.LogicalSwitch
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(logicalSwitchParam nsxModel.LogicalSwitch) (nsxModel.LogicalSwitch, error)

	// Removes a logical switch from the associated overlay or VLAN transport zone. By default, a logical switch cannot be deleted if there are logical ports on the switch, or it is added to a NSGroup. Cascade option can be used to delete all ports and the logical switch. Detach option can be used to delete the logical switch forcibly.
	//  This api is now deprecated. Please use new api - /infra/segments/<segment-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lswitchIdParam (required)
	// @param cascadeParam Delete a Logical Switch and all the logical ports in it, if none of the logical ports have any attachment. (optional, default to false)
	// @param detachParam Force delete a logical switch (optional, default to false)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(lswitchIdParam string, cascadeParam *bool, detachParam *bool) error

	// Returns information about the specified logical switch Id.
	//  This api is now deprecated. Please use new api - /infra/segments/<segment-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lswitchIdParam (required)
	// @return com.vmware.nsx.model.LogicalSwitch
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(lswitchIdParam string) (nsxModel.LogicalSwitch, error)

	// Returns information about all configured logical switches.
	//  This api is now deprecated. Please use new api - /infra/segments
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param diagnosticParam Flag to enable showing of transit logical switch. (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param switchTypeParam Logical Switch type (optional)
	// @param switchingProfileIdParam Switching Profile identifier (optional)
	// @param transportTypeParam Mode of transport supported in the transport zone for this logical switch (optional)
	// @param transportZoneIdParam Transport zone identifier (optional)
	// @param uplinkTeamingPolicyNameParam The logical switch's uplink teaming policy name (optional)
	// @param vlanParam Virtual Local Area Network Identifier (optional)
	// @param vniParam VNI of the OVERLAY LogicalSwitch(es) to return. (optional)
	// @return com.vmware.nsx.model.LogicalSwitchListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, diagnosticParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, switchTypeParam *string, switchingProfileIdParam *string, transportTypeParam *string, transportZoneIdParam *string, uplinkTeamingPolicyNameParam *string, vlanParam *int64, vniParam *int64) (nsxModel.LogicalSwitchListResult, error)

	// Modifies attributes of an existing logical switch. Modifiable attributes include admin_state, replication_mode, switching_profile_ids and VLAN spec. You cannot modify the original transport_zone_id.
	//  This api is now deprecated. Please use new api - PATCH /infra/segments/<segment-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lswitchIdParam (required)
	// @param logicalSwitchParam (required)
	// @return com.vmware.nsx.model.LogicalSwitch
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(lswitchIdParam string, logicalSwitchParam nsxModel.LogicalSwitch) (nsxModel.LogicalSwitch, error)
}

type logicalSwitchesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewLogicalSwitchesClient(connector vapiProtocolClient_.Connector) *logicalSwitchesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.logical_switches")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	lIface := logicalSwitchesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *logicalSwitchesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *logicalSwitchesClient) Create(logicalSwitchParam nsxModel.LogicalSwitch) (nsxModel.LogicalSwitch, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalSwitchesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalSwitchesCreateInputType(), typeConverter)
	sv.AddStructField("LogicalSwitch", logicalSwitchParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalSwitch
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_switches", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalSwitch
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalSwitchesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalSwitch), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalSwitchesClient) Delete(lswitchIdParam string, cascadeParam *bool, detachParam *bool) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalSwitchesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalSwitchesDeleteInputType(), typeConverter)
	sv.AddStructField("LswitchId", lswitchIdParam)
	sv.AddStructField("Cascade", cascadeParam)
	sv.AddStructField("Detach", detachParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_switches", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (lIface *logicalSwitchesClient) Get(lswitchIdParam string) (nsxModel.LogicalSwitch, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalSwitchesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalSwitchesGetInputType(), typeConverter)
	sv.AddStructField("LswitchId", lswitchIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalSwitch
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_switches", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalSwitch
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalSwitchesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalSwitch), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalSwitchesClient) List(cursorParam *string, diagnosticParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, switchTypeParam *string, switchingProfileIdParam *string, transportTypeParam *string, transportZoneIdParam *string, uplinkTeamingPolicyNameParam *string, vlanParam *int64, vniParam *int64) (nsxModel.LogicalSwitchListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalSwitchesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalSwitchesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("Diagnostic", diagnosticParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("SwitchType", switchTypeParam)
	sv.AddStructField("SwitchingProfileId", switchingProfileIdParam)
	sv.AddStructField("TransportType", transportTypeParam)
	sv.AddStructField("TransportZoneId", transportZoneIdParam)
	sv.AddStructField("UplinkTeamingPolicyName", uplinkTeamingPolicyNameParam)
	sv.AddStructField("Vlan", vlanParam)
	sv.AddStructField("Vni", vniParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalSwitchListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_switches", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalSwitchListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalSwitchesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalSwitchListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalSwitchesClient) Update(lswitchIdParam string, logicalSwitchParam nsxModel.LogicalSwitch) (nsxModel.LogicalSwitch, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalSwitchesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalSwitchesUpdateInputType(), typeConverter)
	sv.AddStructField("LswitchId", lswitchIdParam)
	sv.AddStructField("LogicalSwitch", logicalSwitchParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalSwitch
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_switches", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalSwitch
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalSwitchesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalSwitch), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
