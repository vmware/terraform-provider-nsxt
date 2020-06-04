
/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Client stubs for service: BfdProfiles
 * Functions that implement the generated BfdProfilesClient interface
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

type DefaultBfdProfilesClient struct {
	interfaceName       string
	interfaceDefinition core.InterfaceDefinition
	methodIdentifiers   []core.MethodIdentifier
	methodNameToDefMap  map[string]*core.MethodDefinition
	errorBindingMap     map[string]bindings.BindingType
	interfaceIdentifier core.InterfaceIdentifier
	connector           client.Connector
}

func NewDefaultBfdProfilesClient(connector client.Connector) *DefaultBfdProfilesClient {
	interfaceName := "com.vmware.nsx_policy.infra.bfd_profiles"
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


	bIface := DefaultBfdProfilesClient{interfaceName: interfaceName, methodIdentifiers: methodIdentifiers, interfaceDefinition: interfaceDefinition, errorBindingMap: errorBindingMap, interfaceIdentifier: interfaceIdentifier, connector: connector}
	bIface.methodNameToDefMap = make(map[string]*core.MethodDefinition)
	bIface.methodNameToDefMap["delete"] = bIface.deleteMethodDefinition()
	bIface.methodNameToDefMap["get"] = bIface.getMethodDefinition()
	bIface.methodNameToDefMap["list"] = bIface.listMethodDefinition()
	bIface.methodNameToDefMap["patch"] = bIface.patchMethodDefinition()
	bIface.methodNameToDefMap["update"] = bIface.updateMethodDefinition()
	return &bIface
}

func (bIface *DefaultBfdProfilesClient) Delete(bfdProfileIdParam string) error {
	typeConverter := bIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(bIface.interfaceIdentifier, "delete")
	sv := bindings.NewStructValueBuilder(bfdProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("BfdProfileId", bfdProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := bfdProfilesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	bIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := bIface.connector.NewExecutionContext()
	methodResult := bIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (bIface *DefaultBfdProfilesClient) Get(bfdProfileIdParam string) (model.BfdProfile, error) {
	typeConverter := bIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(bIface.interfaceIdentifier, "get")
	sv := bindings.NewStructValueBuilder(bfdProfilesGetInputType(), typeConverter)
	sv.AddStructField("BfdProfileId", bfdProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.BfdProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := bfdProfilesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	bIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := bIface.connector.NewExecutionContext()
	methodResult := bIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.BfdProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), bfdProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.BfdProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *DefaultBfdProfilesClient) List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.BfdProfileListResult, error) {
	typeConverter := bIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(bIface.interfaceIdentifier, "list")
	sv := bindings.NewStructValueBuilder(bfdProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.BfdProfileListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := bfdProfilesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	bIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := bIface.connector.NewExecutionContext()
	methodResult := bIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.BfdProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), bfdProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.BfdProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (bIface *DefaultBfdProfilesClient) Patch(bfdProfileIdParam string, bfdProfileParam model.BfdProfile) error {
	typeConverter := bIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(bIface.interfaceIdentifier, "patch")
	sv := bindings.NewStructValueBuilder(bfdProfilesPatchInputType(), typeConverter)
	sv.AddStructField("BfdProfileId", bfdProfileIdParam)
	sv.AddStructField("BfdProfile", bfdProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := bfdProfilesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	bIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := bIface.connector.NewExecutionContext()
	methodResult := bIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (bIface *DefaultBfdProfilesClient) Update(bfdProfileIdParam string, bfdProfileParam model.BfdProfile) (model.BfdProfile, error) {
	typeConverter := bIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(bIface.interfaceIdentifier, "update")
	sv := bindings.NewStructValueBuilder(bfdProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("BfdProfileId", bfdProfileIdParam)
	sv.AddStructField("BfdProfile", bfdProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.BfdProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := bfdProfilesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	bIface.connector.SetConnectionMetadata(connectionMetadata)
	executionContext := bIface.connector.NewExecutionContext()
	methodResult := bIface.Invoke(executionContext, methodIdentifier, inputDataValue)
	var emptyOutput model.BfdProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), bfdProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.BfdProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), bIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}


func (bIface *DefaultBfdProfilesClient) Invoke(ctx *core.ExecutionContext, methodId core.MethodIdentifier, inputDataValue data.DataValue) core.MethodResult {
	methodResult := bIface.connector.GetApiProvider().Invoke(bIface.interfaceName, methodId.Name(), inputDataValue, ctx)
	return methodResult
}


func (bIface *DefaultBfdProfilesClient) deleteMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(bIface.interfaceName)
	typeConverter := bIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(bfdProfilesDeleteInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(bfdProfilesDeleteOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "delete")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	bIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.delete method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (bIface *DefaultBfdProfilesClient) getMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(bIface.interfaceName)
	typeConverter := bIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(bfdProfilesGetInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(bfdProfilesGetOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "get")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	bIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.get method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (bIface *DefaultBfdProfilesClient) listMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(bIface.interfaceName)
	typeConverter := bIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(bfdProfilesListInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(bfdProfilesListOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "list")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	bIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.list method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (bIface *DefaultBfdProfilesClient) patchMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(bIface.interfaceName)
	typeConverter := bIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(bfdProfilesPatchInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(bfdProfilesPatchOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "patch")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	bIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.patch method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (bIface *DefaultBfdProfilesClient) updateMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(bIface.interfaceName)
	typeConverter := bIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(bfdProfilesUpdateInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(bfdProfilesUpdateOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "update")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	bIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	bIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultBfdProfilesClient.update method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}
