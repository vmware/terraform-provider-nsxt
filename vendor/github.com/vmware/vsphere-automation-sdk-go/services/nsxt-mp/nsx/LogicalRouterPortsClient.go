// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: LogicalRouterPorts
// Used by client-side stubs.

package nsx

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type LogicalRouterPortsClient interface {

	// Creates a logical router port. The required parameters include resource_type (LogicalRouterUpLinkPort, LogicalRouterDownLinkPort, LogicalRouterLinkPort, LogicalRouterLoopbackPort, LogicalRouterCentralizedServicePort); and logical_router_id (the router to which each logical router port is assigned). The service_bindings parameter is optional.
	//
	//  Please use below Policy APIs.
	//  PATCH /policy/api/v1/infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  PATCH /policy/api/v1/infra/tier-1s/<tier-1-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  PATCH /policy/api/v1/infra/tier-1s/<tier-1-id>/segments/<segment-id> for DOWNLINK
	//  PATCH /policy/api/v1/infra/segments/<segment-id> for DOWNLINK
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterPortParam (required)
	// The parameter must contain all the properties defined in nsxModel.LogicalRouterPort.
	// @return com.vmware.nsx.model.LogicalRouterPort
	// The return value will contain all the properties defined in nsxModel.LogicalRouterPort.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(logicalRouterPortParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Deletes the specified logical router port. You must delete logical router ports before you can delete the associated logical router. To Delete Tier0 router link port you must have to delete attached tier1 router link port, otherwise pass \"force=true\" as query param to force delete the Tier0 router link port.
	//
	//  Please use below Policy APIs.
	//  DELETE /policy/api/v1/infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  DELETE /policy/api/v1/infra/tier-1s/<tier-1-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  DELETE /policy/api/v1/infra/tier-1s/<tier-1-id>/segments/<segment-id> for DOWNLINK
	//  DELETE /policy/api/v1/infra/segments/<segment-id> for DOWNLINK
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterPortIdParam (required)
	// @param cascadeDeleteLinkedPortsParam Flag to specify whether to delete related logical switch ports (optional, default to false)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(logicalRouterPortIdParam string, cascadeDeleteLinkedPortsParam *bool, forceParam *bool) error

	// Returns information about the specified logical router port.
	//
	//  Please use below Policy APIs.
	//  GET /policy/api/v1/infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  GET /policy/api/v1/infra/tier-1s/<tier-1-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  GET /policy/api/v1/infra/tier-1s/<tier-1-id>/segments/<segment-id> for DOWNLINK
	//  GET /policy/api/v1/infra/segments/<segment-id> for DOWNLINK
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterPortIdParam (required)
	// @return com.vmware.nsx.model.LogicalRouterPort
	// The return value will contain all the properties defined in nsxModel.LogicalRouterPort.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(logicalRouterPortIdParam string) (*vapiData_.StructValue, error)

	// Returns information about all logical router ports. Information includes the resource_type (LogicalRouterUpLinkPort, LogicalRouterDownLinkPort, LogicalRouterLinkPort, LogicalRouterLoopbackPort, LogicalRouterCentralizedServicePort); logical_router_id (the router to which each logical router port is assigned); and any service_bindings (such as DHCP relay service). The GET request can include a query parameter (logical_router_id or logical_switch_id).
	//
	//  Please use below Policy APIs.
	//  GET /policy/api/v1/infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/interfaces
	//  GET /policy/api/v1/infra/tier-1s/<tier-1-id>/locale-services/<locale-service-id>/interfaces
	//  GET /policy/api/v1/infra/tier-1s/<tier-1-id>/segments for DOWNLINK
	//  GET /policy/api/v1/infra/segments for DOWNLINK
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param logicalRouterIdParam Logical Router identifier (optional)
	// @param logicalSwitchIdParam Logical Switch identifier (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param resourceTypeParam Resource types of logical router port (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.LogicalRouterPortListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, logicalRouterIdParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, resourceTypeParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.LogicalRouterPortListResult, error)

	// Modifies the specified logical router port. Required parameters include the resource_type and logical_router_id. Modifiable parameters include the resource_type (LogicalRouterUpLinkPort, LogicalRouterDownLinkPort, LogicalRouterLinkPort, LogicalRouterLoopbackPort, LogicalRouterCentralizedServicePort), logical_router_id (to reassign the port to a different router), and service_bindings.
	//
	//  Please use below Policy APIs.
	//  PUT /policy/api/v1/infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  PUT /policy/api/v1/infra/tier-1s/<tier-1-id>/locale-services/<locale-service-id>/interfaces/<interface-id>
	//  PUT /policy/api/v1/infra/tier-1s/<tier-1-id>/segments/<segment-id> for DOWNLINK
	//  PUT /policy/api/v1/infra/segments/<segment-id> for DOWNLINK
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterPortIdParam (required)
	// @param logicalRouterPortParam (required)
	// The parameter must contain all the properties defined in nsxModel.LogicalRouterPort.
	// @return com.vmware.nsx.model.LogicalRouterPort
	// The return value will contain all the properties defined in nsxModel.LogicalRouterPort.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(logicalRouterPortIdParam string, logicalRouterPortParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type logicalRouterPortsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewLogicalRouterPortsClient(connector vapiProtocolClient_.Connector) *logicalRouterPortsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.logical_router_ports")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	lIface := logicalRouterPortsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *logicalRouterPortsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *logicalRouterPortsClient) Create(logicalRouterPortParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRouterPortsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRouterPortsCreateInputType(), typeConverter)
	sv.AddStructField("LogicalRouterPort", logicalRouterPortParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_router_ports", "create", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRouterPortsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRouterPortsClient) Delete(logicalRouterPortIdParam string, cascadeDeleteLinkedPortsParam *bool, forceParam *bool) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRouterPortsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRouterPortsDeleteInputType(), typeConverter)
	sv.AddStructField("LogicalRouterPortId", logicalRouterPortIdParam)
	sv.AddStructField("CascadeDeleteLinkedPorts", cascadeDeleteLinkedPortsParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_router_ports", "delete", inputDataValue, executionContext)
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

func (lIface *logicalRouterPortsClient) Get(logicalRouterPortIdParam string) (*vapiData_.StructValue, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRouterPortsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRouterPortsGetInputType(), typeConverter)
	sv.AddStructField("LogicalRouterPortId", logicalRouterPortIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_router_ports", "get", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRouterPortsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRouterPortsClient) List(cursorParam *string, includedFieldsParam *string, logicalRouterIdParam *string, logicalSwitchIdParam *string, pageSizeParam *int64, resourceTypeParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.LogicalRouterPortListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRouterPortsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRouterPortsListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	sv.AddStructField("LogicalSwitchId", logicalSwitchIdParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("ResourceType", resourceTypeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouterPortListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_router_ports", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouterPortListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRouterPortsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouterPortListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRouterPortsClient) Update(logicalRouterPortIdParam string, logicalRouterPortParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRouterPortsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRouterPortsUpdateInputType(), typeConverter)
	sv.AddStructField("LogicalRouterPortId", logicalRouterPortIdParam)
	sv.AddStructField("LogicalRouterPort", logicalRouterPortParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_router_ports", "update", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRouterPortsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
