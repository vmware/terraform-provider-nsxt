
/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Client stubs for service: FloodProtectionProfileBindings
 * Functions that implement the generated FloodProtectionProfileBindingsClient interface
 */


package tier_1s

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

type DefaultFloodProtectionProfileBindingsClient struct {
	interfaceName       string
	interfaceDefinition core.InterfaceDefinition
	methodIdentifiers   []core.MethodIdentifier
	methodNameToDefMap  map[string]*core.MethodDefinition
	errorBindingMap     map[string]bindings.BindingType
	interfaceIdentifier core.InterfaceIdentifier
	connector           client.Connector
}

func NewDefaultFloodProtectionProfileBindingsClient(connector client.Connector) *DefaultFloodProtectionProfileBindingsClient {
	interfaceName := "com.vmware.nsx_policy.infra.tier_1s.flood_protection_profile_bindings"
	interfaceIdentifier := core.NewInterfaceIdentifier(interfaceName)
	methodIdentifiers := []core.MethodIdentifier{
		core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		core.NewMethodIdentifier(interfaceIdentifier, "get"),
		core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorBindingMap := make(map[string]bindings.BindingType)
	errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errorBindingMap[errors.InvalidArgument{}.Error()] = errors.InvalidArgumentBindingType()
	errorBindingMap[errors.OperationNotFound{}.Error()] = errors.OperationNotFoundBindingType()
	errorBindingMap[errors.UnexpectedInput{}.Error()] = errors.UnexpectedInputBindingType()
	errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errorBindingMap[errors.TimedOut{}.Error()] = errors.TimedOutBindingType()
	fIface := DefaultFloodProtectionProfileBindingsClient{interfaceName: interfaceName, methodIdentifiers: methodIdentifiers, interfaceDefinition: interfaceDefinition, errorBindingMap: errorBindingMap, interfaceIdentifier: interfaceIdentifier, connector: connector}
	fIface.methodNameToDefMap = make(map[string]*core.MethodDefinition)
	fIface.methodNameToDefMap["delete"] = fIface.deleteMethodDefinition()
	fIface.methodNameToDefMap["get"] = fIface.getMethodDefinition()
	fIface.methodNameToDefMap["patch"] = fIface.patchMethodDefinition()
	fIface.methodNameToDefMap["update"] = fIface.updateMethodDefinition()
	return &fIface
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) Delete(tier1IdParam string, floodProtectionProfileBindingIdParam string) error {
	typeConverter := fIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(fIface.interfaceIdentifier, "delete")
	sv := bindings.NewStructValueBuilder(floodProtectionProfileBindingsDeleteInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("FloodProtectionProfileBindingId", floodProtectionProfileBindingIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := floodProtectionProfileBindingsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	fIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := fIface.Invoke(fIface.connector.NewExecutionContext(), methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) Get(tier1IdParam string, floodProtectionProfileBindingIdParam string) (model.FloodProtectionProfileBindingMap, error) {
	typeConverter := fIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(fIface.interfaceIdentifier, "get")
	sv := bindings.NewStructValueBuilder(floodProtectionProfileBindingsGetInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("FloodProtectionProfileBindingId", floodProtectionProfileBindingIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.FloodProtectionProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := floodProtectionProfileBindingsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	fIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := fIface.Invoke(fIface.connector.NewExecutionContext(), methodIdentifier, inputDataValue)
	var emptyOutput model.FloodProtectionProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), floodProtectionProfileBindingsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.FloodProtectionProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) Patch(tier1IdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) error {
	typeConverter := fIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(fIface.interfaceIdentifier, "patch")
	sv := bindings.NewStructValueBuilder(floodProtectionProfileBindingsPatchInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("FloodProtectionProfileBindingId", floodProtectionProfileBindingIdParam)
	sv.AddStructField("FloodProtectionProfileBindingMap", floodProtectionProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := floodProtectionProfileBindingsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	fIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := fIface.Invoke(fIface.connector.NewExecutionContext(), methodIdentifier, inputDataValue)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) Update(tier1IdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) (model.FloodProtectionProfileBindingMap, error) {
	typeConverter := fIface.connector.TypeConverter()
	methodIdentifier := core.NewMethodIdentifier(fIface.interfaceIdentifier, "update")
	sv := bindings.NewStructValueBuilder(floodProtectionProfileBindingsUpdateInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("FloodProtectionProfileBindingId", floodProtectionProfileBindingIdParam)
	sv.AddStructField("FloodProtectionProfileBindingMap", floodProtectionProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.FloodProtectionProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := floodProtectionProfileBindingsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	fIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := fIface.Invoke(fIface.connector.NewExecutionContext(), methodIdentifier, inputDataValue)
	var emptyOutput model.FloodProtectionProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), floodProtectionProfileBindingsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.FloodProtectionProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.errorBindingMap[methodResult.Error().Name()])
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}


func (fIface *DefaultFloodProtectionProfileBindingsClient) Invoke(ctx *core.ExecutionContext, methodId core.MethodIdentifier, inputDataValue data.DataValue) core.MethodResult {
	methodResult := fIface.connector.GetApiProvider().Invoke(fIface.interfaceName, methodId.Name(), inputDataValue, ctx)
	return methodResult
}


func (fIface *DefaultFloodProtectionProfileBindingsClient) deleteMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(fIface.interfaceName)
	typeConverter := fIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsDeleteInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsDeleteOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "delete")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	fIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.delete method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) getMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(fIface.interfaceName)
	typeConverter := fIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsGetInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsGetOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "get")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	fIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.get method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) patchMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(fIface.interfaceName)
	typeConverter := fIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsPatchInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsPatchOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "patch")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	fIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.patch method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}

func (fIface *DefaultFloodProtectionProfileBindingsClient) updateMethodDefinition() *core.MethodDefinition {
	interfaceIdentifier := core.NewInterfaceIdentifier(fIface.interfaceName)
	typeConverter := fIface.connector.TypeConverter()

	input, inputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsUpdateInputType())
	output, outputError := typeConverter.ConvertToDataDefinition(floodProtectionProfileBindingsUpdateOutputType())
	if inputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's input - %s",
			bindings.VAPIerrorsToError(inputError).Error())
		return nil
	}
	if outputError != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's output - %s",
			bindings.VAPIerrorsToError(outputError).Error())
		return nil
	}
	methodIdentifier := core.NewMethodIdentifier(interfaceIdentifier, "update")
	errorDefinitions := make([]data.ErrorDefinition, 0)
	fIface.errorBindingMap[errors.InvalidRequest{}.Error()] = errors.InvalidRequestBindingType()
	errDef1, errError1 := typeConverter.ConvertToDataDefinition(errors.InvalidRequestBindingType())
	if errError1 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's errors.InvalidRequest error - %s",
			bindings.VAPIerrorsToError(errError1).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef1.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.Unauthorized{}.Error()] = errors.UnauthorizedBindingType()
	errDef2, errError2 := typeConverter.ConvertToDataDefinition(errors.UnauthorizedBindingType())
	if errError2 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's errors.Unauthorized error - %s",
			bindings.VAPIerrorsToError(errError2).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef2.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.ServiceUnavailable{}.Error()] = errors.ServiceUnavailableBindingType()
	errDef3, errError3 := typeConverter.ConvertToDataDefinition(errors.ServiceUnavailableBindingType())
	if errError3 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's errors.ServiceUnavailable error - %s",
			bindings.VAPIerrorsToError(errError3).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef3.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.InternalServerError{}.Error()] = errors.InternalServerErrorBindingType()
	errDef4, errError4 := typeConverter.ConvertToDataDefinition(errors.InternalServerErrorBindingType())
	if errError4 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's errors.InternalServerError error - %s",
			bindings.VAPIerrorsToError(errError4).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef4.(data.ErrorDefinition))
	fIface.errorBindingMap[errors.NotFound{}.Error()] = errors.NotFoundBindingType()
	errDef5, errError5 := typeConverter.ConvertToDataDefinition(errors.NotFoundBindingType())
	if errError5 != nil {
		log.Errorf("Error in ConvertToDataDefinition for DefaultFloodProtectionProfileBindingsClient.update method's errors.NotFound error - %s",
			bindings.VAPIerrorsToError(errError5).Error())
		return nil
	}
	errorDefinitions = append(errorDefinitions, errDef5.(data.ErrorDefinition))

	methodDefinition := core.NewMethodDefinition(methodIdentifier, input, output, errorDefinitions)
	return &methodDefinition
}
