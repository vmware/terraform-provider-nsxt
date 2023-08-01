// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: BridgeEndpoints
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

type BridgeEndpointsClient interface {

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param bridgeEndpointParam (required)
	// @return com.vmware.nsx.model.BridgeEndpoint
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(bridgeEndpointParam nsxModel.BridgeEndpoint) (nsxModel.BridgeEndpoint, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param bridgeendpointIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(bridgeendpointIdParam string) error

	// Returns information about a specified bridge endpoint.
	//
	// @param bridgeendpointIdParam (required)
	// @return com.vmware.nsx.model.BridgeEndpoint
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(bridgeendpointIdParam string) (nsxModel.BridgeEndpoint, error)

	// Returns information about all configured bridge endoints
	//
	// @param bridgeClusterIdParam Bridge Cluster Identifier (optional)
	// @param bridgeEndpointProfileIdParam Bridge endpoint profile used by the edge cluster (optional)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param logicalSwitchIdParam Logical Switch Identifier (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param vlanTransportZoneIdParam VLAN transport zone id used by the edge cluster (optional)
	// @return com.vmware.nsx.model.BridgeEndpointListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(bridgeClusterIdParam *string, bridgeEndpointProfileIdParam *string, cursorParam *string, includedFieldsParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, vlanTransportZoneIdParam *string) (nsxModel.BridgeEndpointListResult, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param bridgeendpointIdParam (required)
	// @param bridgeEndpointParam (required)
	// @return com.vmware.nsx.model.BridgeEndpoint
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(bridgeendpointIdParam string, bridgeEndpointParam nsxModel.BridgeEndpoint) (nsxModel.BridgeEndpoint, error)
}

type bridgeEndpointsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewBridgeEndpointsClient(connector vapiProtocolClient_.Connector) *bridgeEndpointsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.bridge_endpoints")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	bIface := bridgeEndpointsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &bIface
}

func (bIface *bridgeEndpointsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := bIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (bIface *bridgeEndpointsClient) Create(bridgeEndpointParam nsxModel.BridgeEndpoint) (nsxModel.BridgeEndpoint, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointsCreateInputType(), typeConverter)
	sv.AddStructField("BridgeEndpoint", bridgeEndpointParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpoint
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoints", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpoint
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpoint), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointsClient) Delete(bridgeendpointIdParam string) error {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointsDeleteInputType(), typeConverter)
	sv.AddStructField("BridgeendpointId", bridgeendpointIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoints", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (bIface *bridgeEndpointsClient) Get(bridgeendpointIdParam string) (nsxModel.BridgeEndpoint, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointsGetInputType(), typeConverter)
	sv.AddStructField("BridgeendpointId", bridgeendpointIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpoint
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoints", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpoint
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpoint), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointsClient) List(bridgeClusterIdParam *string, bridgeEndpointProfileIdParam *string, cursorParam *string, includedFieldsParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, vlanTransportZoneIdParam *string) (nsxModel.BridgeEndpointListResult, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointsListInputType(), typeConverter)
	sv.AddStructField("BridgeClusterId", bridgeClusterIdParam)
	sv.AddStructField("BridgeEndpointProfileId", bridgeEndpointProfileIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("LogicalSwitchId", logicalSwitchIdParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("VlanTransportZoneId", vlanTransportZoneIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpointListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoints", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpointListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpointListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointsClient) Update(bridgeendpointIdParam string, bridgeEndpointParam nsxModel.BridgeEndpoint) (nsxModel.BridgeEndpoint, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointsUpdateInputType(), typeConverter)
	sv.AddStructField("BridgeendpointId", bridgeendpointIdParam)
	sv.AddStructField("BridgeEndpoint", bridgeEndpointParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpoint
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoints", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpoint
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpoint), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
