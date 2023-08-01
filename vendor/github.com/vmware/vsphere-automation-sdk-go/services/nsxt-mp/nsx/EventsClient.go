// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Events
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

type EventsClient interface {

	// Returns event associated with event-id.
	//
	// @param eventIdParam (required)
	// @return com.vmware.nsx.model.MonitoringEvent
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(eventIdParam string) (nsxModel.MonitoringEvent, error)

	// Returns a list of all Events defined in NSX.
	// @return com.vmware.nsx.model.EventListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List() (nsxModel.EventListResult, error)

	// Reset all user configurable values for event identified by event-id to factory defaults.
	//
	// @param eventIdParam (required)
	// @return com.vmware.nsx.model.MonitoringEvent
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Setdefault(eventIdParam string) (nsxModel.MonitoringEvent, error)

	// Update event identified by event-id.
	//
	// @param eventIdParam (required)
	// @param monitoringEventParam (required)
	// @return com.vmware.nsx.model.MonitoringEvent
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(eventIdParam string, monitoringEventParam nsxModel.MonitoringEvent) (nsxModel.MonitoringEvent, error)
}

type eventsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewEventsClient(connector vapiProtocolClient_.Connector) *eventsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.events")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"get":        vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"setdefault": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "setdefault"),
		"update":     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	eIface := eventsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &eIface
}

func (eIface *eventsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := eIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (eIface *eventsClient) Get(eventIdParam string) (nsxModel.MonitoringEvent, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := eventsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(eventsGetInputType(), typeConverter)
	sv.AddStructField("EventId", eventIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MonitoringEvent
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.events", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.MonitoringEvent
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EventsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MonitoringEvent), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *eventsClient) List() (nsxModel.EventListResult, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := eventsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(eventsListInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EventListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.events", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.EventListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EventsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EventListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *eventsClient) Setdefault(eventIdParam string) (nsxModel.MonitoringEvent, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := eventsSetdefaultRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(eventsSetdefaultInputType(), typeConverter)
	sv.AddStructField("EventId", eventIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MonitoringEvent
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.events", "setdefault", inputDataValue, executionContext)
	var emptyOutput nsxModel.MonitoringEvent
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EventsSetdefaultOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MonitoringEvent), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *eventsClient) Update(eventIdParam string, monitoringEventParam nsxModel.MonitoringEvent) (nsxModel.MonitoringEvent, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := eventsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(eventsUpdateInputType(), typeConverter)
	sv.AddStructField("EventId", eventIdParam)
	sv.AddStructField("MonitoringEvent", monitoringEventParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MonitoringEvent
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.events", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.MonitoringEvent
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EventsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MonitoringEvent), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
