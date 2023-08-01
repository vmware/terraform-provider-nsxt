// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: BridgeEndpointProfiles
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

type BridgeEndpointProfilesClient interface {

	// Creates a Bridge Endpoint Profile. Profile contains edge cluster id, indexes of the member nodes, fialover mode and high availability mode for a Bridge EndPoint
	//
	// @param bridgeEndpointProfileParam (required)
	// @return com.vmware.nsx.model.BridgeEndpointProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(bridgeEndpointProfileParam nsxModel.BridgeEndpointProfile) (nsxModel.BridgeEndpointProfile, error)

	// Deletes the specified Bridge Endpoint Profile.
	//
	// @param bridgeendpointprofileIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(bridgeendpointprofileIdParam string) error

	// Returns information about a specified bridge endpoint profile.
	//
	// @param bridgeendpointprofileIdParam (required)
	// @return com.vmware.nsx.model.BridgeEndpointProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(bridgeendpointprofileIdParam string) (nsxModel.BridgeEndpointProfile, error)

	// Returns information about all configured bridge endoint profiles
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param edgeClusterIdParam Edge Cluster Identifier (optional)
	// @param failoverModeParam (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.BridgeEndpointProfileListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, edgeClusterIdParam *string, failoverModeParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.BridgeEndpointProfileListResult, error)

	// Modifies a existing bridge endpoint profile.
	//
	// @param bridgeendpointprofileIdParam (required)
	// @param bridgeEndpointProfileParam (required)
	// @return com.vmware.nsx.model.BridgeEndpointProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(bridgeendpointprofileIdParam string, bridgeEndpointProfileParam nsxModel.BridgeEndpointProfile) (nsxModel.BridgeEndpointProfile, error)
}

type bridgeEndpointProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewBridgeEndpointProfilesClient(connector vapiProtocolClient_.Connector) *bridgeEndpointProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.bridge_endpoint_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	bIface := bridgeEndpointProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &bIface
}

func (bIface *bridgeEndpointProfilesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := bIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (bIface *bridgeEndpointProfilesClient) Create(bridgeEndpointProfileParam nsxModel.BridgeEndpointProfile) (nsxModel.BridgeEndpointProfile, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointProfilesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointProfilesCreateInputType(), typeConverter)
	sv.AddStructField("BridgeEndpointProfile", bridgeEndpointProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpointProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoint_profiles", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpointProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointProfilesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpointProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointProfilesClient) Delete(bridgeendpointprofileIdParam string) error {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointProfilesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("BridgeendpointprofileId", bridgeendpointprofileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoint_profiles", "delete", inputDataValue, executionContext)
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

func (bIface *bridgeEndpointProfilesClient) Get(bridgeendpointprofileIdParam string) (nsxModel.BridgeEndpointProfile, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointProfilesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointProfilesGetInputType(), typeConverter)
	sv.AddStructField("BridgeendpointprofileId", bridgeendpointprofileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpointProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoint_profiles", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpointProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpointProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointProfilesClient) List(cursorParam *string, edgeClusterIdParam *string, failoverModeParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.BridgeEndpointProfileListResult, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("FailoverMode", failoverModeParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpointProfileListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoint_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpointProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpointProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *bridgeEndpointProfilesClient) Update(bridgeendpointprofileIdParam string, bridgeEndpointProfileParam nsxModel.BridgeEndpointProfile) (nsxModel.BridgeEndpointProfile, error) {
	typeConverter := bIface.connector.TypeConverter()
	executionContext := bIface.connector.NewExecutionContext()
	operationRestMetaData := bridgeEndpointProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(bridgeEndpointProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("BridgeendpointprofileId", bridgeendpointprofileIdParam)
	sv.AddStructField("BridgeEndpointProfile", bridgeEndpointProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.BridgeEndpointProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := bIface.connector.GetApiProvider().Invoke("com.vmware.nsx.bridge_endpoint_profiles", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.BridgeEndpointProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), BridgeEndpointProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.BridgeEndpointProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
