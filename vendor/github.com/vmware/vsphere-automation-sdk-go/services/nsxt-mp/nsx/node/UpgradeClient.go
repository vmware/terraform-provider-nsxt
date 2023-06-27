// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Upgrade
// Used by client-side stubs.

package node

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type UpgradeClient interface {

	// Get upgrade task status for the given task of the given bundle. Both bundle_name and task_id must be provided, otherwise you will receive a 404 NOT FOUND response. This api is deprecated. Please use API GET /node/upgrade/status-summary to get upgrade status when upgrade is in progress.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param bundleNameParam Bundle Name (optional)
	// @param upgradeTaskIdParam Upgrade Task ID (optional)
	// @return com.vmware.nsx.model.UpgradeTaskProperties
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(bundleNameParam *string, upgradeTaskIdParam *string) (nsxModel.UpgradeTaskProperties, error)
}

type upgradeClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewUpgradeClient(connector vapiProtocolClient_.Connector) *upgradeClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.node.upgrade")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"get": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	uIface := upgradeClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &uIface
}

func (uIface *upgradeClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := uIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (uIface *upgradeClient) Get(bundleNameParam *string, upgradeTaskIdParam *string) (nsxModel.UpgradeTaskProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := upgradeGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(upgradeGetInputType(), typeConverter)
	sv.AddStructField("BundleName", bundleNameParam)
	sv.AddStructField("UpgradeTaskId", upgradeTaskIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.UpgradeTaskProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.upgrade", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.UpgradeTaskProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UpgradeGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.UpgradeTaskProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
