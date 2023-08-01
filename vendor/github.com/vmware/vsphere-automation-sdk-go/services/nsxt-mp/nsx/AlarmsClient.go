// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Alarms
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

type AlarmsClient interface {

	// Returns alarm associated with alarm-id. If HTTP status 404 is returned, this means the specified alarm-id is invalid or the alarm with alarm-id has been deleted. An alarm is deleted by the system if it is RESOLVED and older than eight days. The system can also delete the remaining RESOLVED alarms sooner to free system resources when too many alarms are being generated. When this happens the oldest day's RESOLVED alarms are deleted first.
	//
	// @param alarmIdParam (required)
	// @return com.vmware.nsx.model.Alarm
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(alarmIdParam string) (nsxModel.Alarm, error)

	// Returns a list of all Alarms currently known to the system.
	//
	// @param afterParam Timestamp in milliseconds since epoch (optional)
	// @param beforeParam Timestamp in milliseconds since epoch (optional)
	// @param cursorParam Cursor for pagination (optional)
	// @param eventTagParam Event tag (optional)
	// @param eventTypeParam Event Type Filter (optional)
	// @param featureNameParam Feature Name (optional)
	// @param idParam Alarm ID (optional)
	// @param intentPathParam Intent Path for entity ID (optional)
	// @param nodeIdParam Node ID (optional)
	// @param nodeResourceTypeParam Node Resource Type (optional)
	// @param orgParam Org ID (optional)
	// @param pageSizeParam Page Size for pagination (optional)
	// @param projectParam Project ID (optional)
	// @param severityParam Severity (optional)
	// @param sortAscendingParam Represents order of sorting the values (optional, default to true)
	// @param sortByParam Key for sorting on this column (optional)
	// @param statusParam Status (optional)
	// @param vpcParam VPC ID (optional)
	// @return com.vmware.nsx.model.AlarmsListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(afterParam *int64, beforeParam *int64, cursorParam *string, eventTagParam *string, eventTypeParam *string, featureNameParam *string, idParam *string, intentPathParam *string, nodeIdParam *string, nodeResourceTypeParam *string, orgParam *string, pageSizeParam *int64, projectParam *string, severityParam *string, sortAscendingParam *bool, sortByParam *string, statusParam *string, vpcParam *string) (nsxModel.AlarmsListResult, error)

	// Update status of an Alarm. The new_status value can be OPEN, ACKNOWLEDGED, SUPPRESSED, or RESOLVED. If new_status is SUPPRESSED, the suppress_duration query parameter must also be specified.
	//
	// @param alarmIdParam (required)
	// @param newStatusParam Status (required)
	// @param suppressDurationParam Duration in hours for which Alarm should be suppressed (optional)
	// @return com.vmware.nsx.model.Alarm
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Setstatus(alarmIdParam string, newStatusParam string, suppressDurationParam *int64) (nsxModel.Alarm, error)

	// Bulk update the status of zero or more Alarms that match the specified filters. The new_status value can be OPEN, ACKNOWLEDGED, SUPPRESSED, or RESOLVED. If new_status is SUPPRESSED, the suppress_duration query parameter must also be specified.
	//
	// @param newStatusParam Status (required)
	// @param afterParam Timestamp in milliseconds since epoch (optional)
	// @param beforeParam Timestamp in milliseconds since epoch (optional)
	// @param cursorParam Cursor for pagination (optional)
	// @param eventTagParam Event tag (optional)
	// @param eventTypeParam Event Type Filter (optional)
	// @param featureNameParam Feature Name (optional)
	// @param idParam Alarm ID (optional)
	// @param intentPathParam Intent Path for entity ID (optional)
	// @param nodeIdParam Node ID (optional)
	// @param nodeResourceTypeParam Node Resource Type (optional)
	// @param orgParam Org ID (optional)
	// @param pageSizeParam Page Size for pagination (optional)
	// @param projectParam Project ID (optional)
	// @param severityParam Severity (optional)
	// @param sortAscendingParam Represents order of sorting the values (optional, default to true)
	// @param sortByParam Key for sorting on this column (optional)
	// @param statusParam Status (optional)
	// @param suppressDurationParam Duration in hours for which Alarm should be suppressed (optional)
	// @param vpcParam VPC ID (optional)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Setstatus0(newStatusParam string, afterParam *int64, beforeParam *int64, cursorParam *string, eventTagParam *string, eventTypeParam *string, featureNameParam *string, idParam *string, intentPathParam *string, nodeIdParam *string, nodeResourceTypeParam *string, orgParam *string, pageSizeParam *int64, projectParam *string, severityParam *string, sortAscendingParam *bool, sortByParam *string, statusParam *string, suppressDurationParam *int64, vpcParam *string) error
}

type alarmsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewAlarmsClient(connector vapiProtocolClient_.Connector) *alarmsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.alarms")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"get":         vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":        vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"setstatus":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "setstatus"),
		"setstatus_0": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "setstatus_0"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	aIface := alarmsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &aIface
}

func (aIface *alarmsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := aIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (aIface *alarmsClient) Get(alarmIdParam string) (nsxModel.Alarm, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	operationRestMetaData := alarmsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(alarmsGetInputType(), typeConverter)
	sv.AddStructField("AlarmId", alarmIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.Alarm
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx.alarms", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.Alarm
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), AlarmsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.Alarm), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *alarmsClient) List(afterParam *int64, beforeParam *int64, cursorParam *string, eventTagParam *string, eventTypeParam *string, featureNameParam *string, idParam *string, intentPathParam *string, nodeIdParam *string, nodeResourceTypeParam *string, orgParam *string, pageSizeParam *int64, projectParam *string, severityParam *string, sortAscendingParam *bool, sortByParam *string, statusParam *string, vpcParam *string) (nsxModel.AlarmsListResult, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	operationRestMetaData := alarmsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(alarmsListInputType(), typeConverter)
	sv.AddStructField("After", afterParam)
	sv.AddStructField("Before", beforeParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("EventTag", eventTagParam)
	sv.AddStructField("EventType", eventTypeParam)
	sv.AddStructField("FeatureName", featureNameParam)
	sv.AddStructField("Id", idParam)
	sv.AddStructField("IntentPath", intentPathParam)
	sv.AddStructField("NodeId", nodeIdParam)
	sv.AddStructField("NodeResourceType", nodeResourceTypeParam)
	sv.AddStructField("Org", orgParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("Project", projectParam)
	sv.AddStructField("Severity", severityParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("Status", statusParam)
	sv.AddStructField("Vpc", vpcParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.AlarmsListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx.alarms", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.AlarmsListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), AlarmsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.AlarmsListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *alarmsClient) Setstatus(alarmIdParam string, newStatusParam string, suppressDurationParam *int64) (nsxModel.Alarm, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	operationRestMetaData := alarmsSetstatusRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(alarmsSetstatusInputType(), typeConverter)
	sv.AddStructField("AlarmId", alarmIdParam)
	sv.AddStructField("NewStatus", newStatusParam)
	sv.AddStructField("SuppressDuration", suppressDurationParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.Alarm
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx.alarms", "setstatus", inputDataValue, executionContext)
	var emptyOutput nsxModel.Alarm
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), AlarmsSetstatusOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.Alarm), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *alarmsClient) Setstatus0(newStatusParam string, afterParam *int64, beforeParam *int64, cursorParam *string, eventTagParam *string, eventTypeParam *string, featureNameParam *string, idParam *string, intentPathParam *string, nodeIdParam *string, nodeResourceTypeParam *string, orgParam *string, pageSizeParam *int64, projectParam *string, severityParam *string, sortAscendingParam *bool, sortByParam *string, statusParam *string, suppressDurationParam *int64, vpcParam *string) error {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	operationRestMetaData := alarmsSetstatus0RestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(alarmsSetstatus0InputType(), typeConverter)
	sv.AddStructField("NewStatus", newStatusParam)
	sv.AddStructField("After", afterParam)
	sv.AddStructField("Before", beforeParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("EventTag", eventTagParam)
	sv.AddStructField("EventType", eventTypeParam)
	sv.AddStructField("FeatureName", featureNameParam)
	sv.AddStructField("Id", idParam)
	sv.AddStructField("IntentPath", intentPathParam)
	sv.AddStructField("NodeId", nodeIdParam)
	sv.AddStructField("NodeResourceType", nodeResourceTypeParam)
	sv.AddStructField("Org", orgParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("Project", projectParam)
	sv.AddStructField("Severity", severityParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("Status", statusParam)
	sv.AddStructField("SuppressDuration", suppressDurationParam)
	sv.AddStructField("Vpc", vpcParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx.alarms", "setstatus_0", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}
