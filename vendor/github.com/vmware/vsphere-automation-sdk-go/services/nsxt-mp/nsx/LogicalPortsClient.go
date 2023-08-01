// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: LogicalPorts
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

type LogicalPortsClient interface {

	// Creates a new logical switch port. The required parameters are the associated logical_switch_id and admin_state (UP or DOWN). Optional parameters are the attachment and switching_profile_ids. If you don't specify switching_profile_ids, default switching profiles are assigned to the port. If you don't specify an attachment, the switch port remains empty. To configure an attachment, you must specify an id, and optionally you can specify an attachment_type (VIF or LOGICALROUTER). The attachment_type is VIF by default.
	//  This api is now deprecated. Please use new api PUT /infra/segments/<segment-id>/ports/<port-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalPortParam (required)
	// @return com.vmware.nsx.model.LogicalPort
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(logicalPortParam nsxModel.LogicalPort) (nsxModel.LogicalPort, error)

	// Deletes the specified logical switch port. By default, if logical port has attachments, or it is added to any NSGroup, the deletion will be failed. Option detach could be used for deleting logical port forcibly.
	//  This api is now deprecated. Please use new api - DELETE /infra/segments/<segment-id>/ports/<port-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lportIdParam (required)
	// @param detachParam force delete even if attached or referenced by a group (optional, default to false)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(lportIdParam string, detachParam *bool) error

	// Returns information about a specified logical port.
	//  Please use corresponding policy API /infra/segments/<segment-id>/ports/<lport-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lportIdParam (required)
	// @return com.vmware.nsx.model.LogicalPort
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(lportIdParam string) (nsxModel.LogicalPort, error)

	// Returns information about all configured logical switch ports. Logical switch ports connect to VM virtual network interface cards (NICs). Each logical port is associated with one logical switch.
	//  This api is now deprecated. Please use new api - /infra/segments/<segment-id>/ports
	//
	// Deprecated: This API element is deprecated.
	//
	// @param attachmentIdParam Logical Port attachment Id (optional)
	// @param attachmentTypeParam Type of attachment for logical port; for query only. (optional)
	// @param bridgeClusterIdParam Bridge Cluster identifier (optional)
	// @param containerPortsOnlyParam Only container VIF logical ports will be returned if true (optional, default to false)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param diagnosticParam Flag to enable showing of transit logical port. (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param logicalSwitchIdParam Logical Switch identifier (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param parentVifIdParam ID of the VIF of type PARENT (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param switchingProfileIdParam Network Profile identifier (optional)
	// @param transportNodeIdParam Transport node identifier (optional)
	// @param transportZoneIdParam Transport zone identifier (optional)
	// @return com.vmware.nsx.model.LogicalPortListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(attachmentIdParam *string, attachmentTypeParam *string, bridgeClusterIdParam *string, containerPortsOnlyParam *bool, cursorParam *string, diagnosticParam *bool, includedFieldsParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, parentVifIdParam *string, sortAscendingParam *bool, sortByParam *string, switchingProfileIdParam *string, transportNodeIdParam *string, transportZoneIdParam *string) (nsxModel.LogicalPortListResult, error)

	// Modifies an existing logical switch port. Parameters that can be modified include attachment_type (LOGICALROUTER, VIF), admin_state (UP or DOWN), attachment id and switching_profile_ids. You cannot modify the logical_switch_id. In other words, you cannot move an existing port from one switch to another switch.
	//  This api is now deprecated. Please use new api - /infra/segments/<segment-id>/ports/<port-id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param lportIdParam (required)
	// @param logicalPortParam (required)
	// @return com.vmware.nsx.model.LogicalPort
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(lportIdParam string, logicalPortParam nsxModel.LogicalPort) (nsxModel.LogicalPort, error)
}

type logicalPortsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewLogicalPortsClient(connector vapiProtocolClient_.Connector) *logicalPortsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.logical_ports")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	lIface := logicalPortsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *logicalPortsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *logicalPortsClient) Create(logicalPortParam nsxModel.LogicalPort) (nsxModel.LogicalPort, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalPortsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalPortsCreateInputType(), typeConverter)
	sv.AddStructField("LogicalPort", logicalPortParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalPort
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_ports", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalPort
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalPortsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalPort), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalPortsClient) Delete(lportIdParam string, detachParam *bool) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalPortsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalPortsDeleteInputType(), typeConverter)
	sv.AddStructField("LportId", lportIdParam)
	sv.AddStructField("Detach", detachParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_ports", "delete", inputDataValue, executionContext)
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

func (lIface *logicalPortsClient) Get(lportIdParam string) (nsxModel.LogicalPort, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalPortsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalPortsGetInputType(), typeConverter)
	sv.AddStructField("LportId", lportIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalPort
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_ports", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalPort
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalPortsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalPort), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalPortsClient) List(attachmentIdParam *string, attachmentTypeParam *string, bridgeClusterIdParam *string, containerPortsOnlyParam *bool, cursorParam *string, diagnosticParam *bool, includedFieldsParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, parentVifIdParam *string, sortAscendingParam *bool, sortByParam *string, switchingProfileIdParam *string, transportNodeIdParam *string, transportZoneIdParam *string) (nsxModel.LogicalPortListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalPortsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalPortsListInputType(), typeConverter)
	sv.AddStructField("AttachmentId", attachmentIdParam)
	sv.AddStructField("AttachmentType", attachmentTypeParam)
	sv.AddStructField("BridgeClusterId", bridgeClusterIdParam)
	sv.AddStructField("ContainerPortsOnly", containerPortsOnlyParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("Diagnostic", diagnosticParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("LogicalSwitchId", logicalSwitchIdParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("ParentVifId", parentVifIdParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("SwitchingProfileId", switchingProfileIdParam)
	sv.AddStructField("TransportNodeId", transportNodeIdParam)
	sv.AddStructField("TransportZoneId", transportZoneIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalPortListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_ports", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalPortListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalPortsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalPortListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalPortsClient) Update(lportIdParam string, logicalPortParam nsxModel.LogicalPort) (nsxModel.LogicalPort, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalPortsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalPortsUpdateInputType(), typeConverter)
	sv.AddStructField("LportId", lportIdParam)
	sv.AddStructField("LogicalPort", logicalPortParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalPort
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_ports", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalPort
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalPortsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalPort), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
