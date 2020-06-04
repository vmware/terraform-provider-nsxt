
/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Client stubs for service: IpfixCollectorProfiles
 * Functions that implement the generated IpfixCollectorProfilesClient interface
 */


package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

type DefaultIpfixCollectorProfilesClient struct {
	interfaceName       string
	interfaceDefinition core.InterfaceDefinition
	methodIdentifiers   []core.MethodIdentifier
	methodNameToDefMap  map[string]*core.MethodDefinition
	errorBindingMap     map[string]bindings.BindingType
	interfaceIdentifier core.InterfaceIdentifier
	connector           client.Connector
}

func NewDefaultIpfixCollectorProfilesClient(connector client.Connector) *DefaultIpfixCollectorProfilesClient {
	interfaceName := "com.vmware.nsx_policy.infra.ipfix_collector_profiles"
	interfaceIdentifier := core.NewInterfaceIdentifier(interfaceName)
	methodIdentifiers := []core.MethodIdentifier{
		core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		core.NewMethodIdentifier(interfaceIdentifier, "get"),
		core.NewMethodIdentifier(interfaceIdentifier, "list"),
		core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorBindingMap := make(map[string]bindings.BindingType)
	errorBindingMap[errors.AlreadyExists{}.Error()] = errors.AlreadyExistsBindingType()
	errorBindingMap[errors.AlreadyInDesiredState{}.Error()] = errors.AlreadyInDesiredStateBindingType()
	errorBindingMap[errors.Canceled{}.Error()] = errors.CanceledBindingType()
	errorBindingMap[errors.ConcurrentChange{}.Error()] = errors.ConcurrentChangeBindingType()
	errorBindingMap[errors.Error{}.Error()] = errors.ErrorBindingType()
	errorBindingMap[errors.FeatureInUse{}.Error()] = errors.FeatureInUseBindingType()
	errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errorBindingMap[errors.InvalidArgument{}.Error()] = errors.InvalidArgumentBindingType()
	errorBindingMap[errors.InvalidElementConfiguration{}.Error()] = errors.InvalidElementConfigurationBindingType()
	errorBindingMap[errors.InvalidElementType{}.Error()] = errors.InvalidElementTypeBindingType()
	errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errorBindingMap[errors.NotAllowedInCurrentState{}.Error()] = errors.NotAllowedInCurrentStateBindingType()
	errorBindingMap[errors.OperationNotFound{}.Error()] = errors.OperationNotFoundBindingType()
	errorBindingMap[errors.ResourceBusy{}.Error()] = errors.ResourceBusyBindingType()
	errorBindingMap[errors.ResourceInUse{}.Error()] = errors.ResourceInUseBindingType()
	errorBindingMap[errors.ResourceInaccessible{}.Error()] = errors.ResourceInaccessibleBindingType()
	errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errorBindingMap[errors.TimedOut{}.Error()] = errors.TimedOutBindingType()
	errorBindingMap[errors.UnableToAllocateResource{}.Error()] = errors.UnableToAllocateResourceBindingType()
	errorBindingMap[errors.Unauthenticated{}.Error()] = errors.UnauthenticatedBindingType()
	errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errorBindingMap[errors.UnexpectedInput{}.Error()] = errors.UnexpectedInputBindingType()
	errorBindingMap[errors.Unsupported{}.Error()] = errors.UnsupportedBindingType()
	errorBindingMap[errors.UnverifiedPeer{}.Error()] = errors.UnverifiedPeerBindingType()


	iIface := DefaultIpfixCollectorProfilesClient{interfaceName: interfaceName, methodIdentifiers: methodIdentifiers, interfaceDefinition: interfaceDefinition, errorBindingMap: errorBindingMap, interfaceIdentifier: interfaceIdentifier, connector: connector}
	iIface.methodNameToDefMap = make(map[string]*core.MethodDefinition)
	iIface.methodNameToDefMap["delete"] = iIface.deleteMethodDefinition()
	iIface.methodNameToDefMap["get"] = iIface.getMethodDefinition()
	iIface.methodNameToDefMap["list"] = iIface.listMethodDefinition()
	iIface.methodNameToDefMap["patch"] = iIface.patchMethodDefinition()
	iIface.methodNameToDefMap["update"] = iIface.updateMethodDefinition()
	return &iIface
}

func (iIface *DefaultIpfixCollectorProfilesClient) Delete(ipfixCollectorProfileIdParam string) error {
	typeConverter := iIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(iIface.interfaceIdentifier, "delete")
	sv := bindings.NewStructValueBuilder(ipfixCollectorProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("IpfixCollectorProfileId", ipfixCollectorProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixCollectorProfilesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := iIface.connector.NewExecutionContext()
	methodResult := iIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (iIface *DefaultIpfixCollectorProfilesClient) Get(ipfixCollectorProfileIdParam string) (model.IPFIXCollectorProfile, error) {
	typeConverter := iIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(iIface.interfaceIdentifier, "get")
	sv := bindings.NewStructValueBuilder(ipfixCollectorProfilesGetInputType(), typeConverter)
	sv.AddStructField("IpfixCollectorProfileId", ipfixCollectorProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXCollectorProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixCollectorProfilesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := iIface.connector.NewExecutionContext()
	methodResult := iIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.IPFIXCollectorProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixCollectorProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXCollectorProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *DefaultIpfixCollectorProfilesClient) List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXCollectorProfileListResult, error) {
	typeConverter := iIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(iIface.interfaceIdentifier, "list")
	sv := bindings.NewStructValueBuilder(ipfixCollectorProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXCollectorProfileListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixCollectorProfilesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := iIface.connector.NewExecutionContext()
	methodResult := iIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.IPFIXCollectorProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixCollectorProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXCollectorProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *DefaultIpfixCollectorProfilesClient) Patch(ipfixCollectorProfileIdParam string, iPFIXCollectorProfileParam model.IPFIXCollectorProfile) error {
	typeConverter := iIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(iIface.interfaceIdentifier, "patch")
	sv := bindings.NewStructValueBuilder(ipfixCollectorProfilesPatchInputType(), typeConverter)
	sv.AddStructField("IpfixCollectorProfileId", ipfixCollectorProfileIdParam)
	sv.AddStructField("IPFIXCollectorProfile", iPFIXCollectorProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixCollectorProfilesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := iIface.connector.NewExecutionContext()
	methodResult := iIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (iIface *DefaultIpfixCollectorProfilesClient) Update(ipfixCollectorProfileIdParam string, iPFIXCollectorProfileParam model.IPFIXCollectorProfile) (model.IPFIXCollectorProfile, error) {
	typeConverter := iIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(iIface.interfaceIdentifier, "update")
	sv := bindings.NewStructValueBuilder(ipfixCollectorProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("IpfixCollectorProfileId", ipfixCollectorProfileIdParam)
	sv.AddStructField("IPFIXCollectorProfile", iPFIXCollectorProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPFIXCollectorProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := ipfixCollectorProfilesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	iIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := iIface.connector.NewExecutionContext()
	methodResult := iIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.IPFIXCollectorProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ipfixCollectorProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPFIXCollectorProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}


func (iIface *DefaultIpfixCollectorProfilesClient) Invoke(ctx *core.ExecutionContext, methodId core.MethodIdentifier, inputDataValue data.DataValue) core.MethodResult {
	methodResult := iIface.connector.GetApiProvider().Invoke(iIface.interfaceName, methodId.Name(), inputDataValue, ctx)
	return methodResult
}


func (iIface *DefaultIpfixCollectorProfilesClient) deleteMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(iIface.interfaceName)
	typeConverter := iIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesDeleteInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesDeleteOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "delete")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	iIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.delete method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (iIface *DefaultIpfixCollectorProfilesClient) getMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(iIface.interfaceName)
	typeConverter := iIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesGetInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesGetOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "get")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	iIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.get method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (iIface *DefaultIpfixCollectorProfilesClient) listMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(iIface.interfaceName)
	typeConverter := iIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesListInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesListOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "list")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	iIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.list method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (iIface *DefaultIpfixCollectorProfilesClient) patchMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(iIface.interfaceName)
	typeConverter := iIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesPatchInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesPatchOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "patch")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	iIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.patch method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (iIface *DefaultIpfixCollectorProfilesClient) updateMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(iIface.interfaceName)
	typeConverter := iIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesUpdateInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(ipfixCollectorProfilesUpdateOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "update")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	iIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	iIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultIpfixCollectorProfilesClient.update method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}
