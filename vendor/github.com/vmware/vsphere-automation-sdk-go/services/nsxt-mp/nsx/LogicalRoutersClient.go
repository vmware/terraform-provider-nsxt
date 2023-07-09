// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: LogicalRouters
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

type LogicalRoutersClient interface {

	// Creates a logical router. The required parameters are router_type (TIER0 or TIER1) and edge_cluster_id (TIER0 only). Optional parameters include internal and external transit network addresses.
	//
	//  Please use below policy apis instead of this API.
	//  PATCH /infra/tier-0s/<id>
	//  PATCH /infra/tier-0s/<id>/locale-services/<id>
	//  PATCH /infra/tier-1s/<id>
	//  PATCH /infra/tier-1s/<id>/locale-services/<id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterParam (required)
	// @return com.vmware.nsx.model.LogicalRouter
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(logicalRouterParam nsxModel.LogicalRouter) (nsxModel.LogicalRouter, error)

	// Deletes the specified logical router. You must delete associated logical router ports before you can delete a logical router. Otherwise use force delete which will delete all related ports and other entities associated with that LR. To force delete logical router pass force=true in query param.
	//
	//  Please use below policy apis instead of this API.
	//  DELETE /infra/tier-0s/<id>
	//  DELETE /infra/tier-1s/<id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterIdParam (required)
	// @param cascadeDeleteLinkedPortsParam Flag to specify whether to delete related logical switch ports (optional, default to false)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(logicalRouterIdParam string, cascadeDeleteLinkedPortsParam *bool, forceParam *bool) error

	// Returns information about the specified logical router.
	//
	//  Please use below policy apis instead of this API.
	//  GET /infra/tier-0s/<id>
	//  GET /infra/tier-0s/<id>/locale-services/<id>
	//  GET /infra/tier-1s/<id>
	//  GET /infra/tier-1s/<id>/locale-services/<id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterIdParam (required)
	// @return com.vmware.nsx.model.LogicalRouter
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(logicalRouterIdParam string) (nsxModel.LogicalRouter, error)

	// Returns information about all logical routers, including the UUID, internal and external transit network addresses, and the router type (TIER0 or TIER1). You can get information for only TIER0 routers or only the TIER1 routers by including the router_type query parameter.
	//
	//  Please use below policy apis instead of this API.
	//  GET /infra/tier-0s/
	//  GET /infra/tier-1s/
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param routerTypeParam Type of Logical Router (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param vrfsOnLogicalRouterIdParam List all VRFs on the specified logical router. (optional)
	// @return com.vmware.nsx.model.LogicalRouterListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, routerTypeParam *string, sortAscendingParam *bool, sortByParam *string, vrfsOnLogicalRouterIdParam *string) (nsxModel.LogicalRouterListResult, error)

	// API to re allocate edge node placement for TIER1 logical router. You can re-allocate service routers of TIER1 in same edge cluster or different edge cluster. You can also place edge nodes manually and provide maximum two indices for HA mode ACTIVE_STANDBY. To re-allocate on new edge cluster you must have existing edge cluster for TIER1 logical router. This will be disruptive operation and all existing statistics of logical router will be remove.
	//
	//  In policy there will be no equivalent API but it will be achieved automatically when you will change edge cluster or will remove edge nodes from tier1 locale service
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterIdParam (required)
	// @param serviceRouterAllocationConfigParam (required)
	// @return com.vmware.nsx.model.LogicalRouter
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Reallocate(logicalRouterIdParam string, serviceRouterAllocationConfigParam nsxModel.ServiceRouterAllocationConfig) (nsxModel.LogicalRouter, error)

	// Reprocess logical router configuration and configuration of related entities like logical router ports, static routing, etc. Any missing Updates are published to controller.
	//
	//  Please use below policy apis instead of this API.
	//  POST /infra/tier-0s/<tier-0-id>?action=reprocess
	//  POST /infra/tier-1s/<tier-1-id>?action=reprocess
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Reprocess(logicalRouterIdParam string) error

	// Modifies the specified logical router. Modifiable attributes include the internal_transit_network, external_transit_networks, and edge_cluster_id (for TIER0 routers).
	//
	//  Please use below policy apis instead of this API.
	//  PUT /infra/tier-0s/<id>
	//  PUT /infra/tier-0s/<id>/locale-services/<id>
	//  PUT /infra/tier-1s/<id>
	//  PUT /infra/tier-1s/<id>/locale-services/<id>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param logicalRouterIdParam (required)
	// @param logicalRouterParam (required)
	// @return com.vmware.nsx.model.LogicalRouter
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(logicalRouterIdParam string, logicalRouterParam nsxModel.LogicalRouter) (nsxModel.LogicalRouter, error)
}

type logicalRoutersClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewLogicalRoutersClient(connector vapiProtocolClient_.Connector) *logicalRoutersClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.logical_routers")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create":     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete":     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":        vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"reallocate": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "reallocate"),
		"reprocess":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "reprocess"),
		"update":     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	lIface := logicalRoutersClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *logicalRoutersClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *logicalRoutersClient) Create(logicalRouterParam nsxModel.LogicalRouter) (nsxModel.LogicalRouter, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersCreateInputType(), typeConverter)
	sv.AddStructField("LogicalRouter", logicalRouterParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouter
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouter
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRoutersCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouter), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRoutersClient) Delete(logicalRouterIdParam string, cascadeDeleteLinkedPortsParam *bool, forceParam *bool) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersDeleteInputType(), typeConverter)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	sv.AddStructField("CascadeDeleteLinkedPorts", cascadeDeleteLinkedPortsParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "delete", inputDataValue, executionContext)
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

func (lIface *logicalRoutersClient) Get(logicalRouterIdParam string) (nsxModel.LogicalRouter, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersGetInputType(), typeConverter)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouter
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouter
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRoutersGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouter), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRoutersClient) List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, routerTypeParam *string, sortAscendingParam *bool, sortByParam *string, vrfsOnLogicalRouterIdParam *string) (nsxModel.LogicalRouterListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("RouterType", routerTypeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("VrfsOnLogicalRouterId", vrfsOnLogicalRouterIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouterListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouterListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRoutersListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouterListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRoutersClient) Reallocate(logicalRouterIdParam string, serviceRouterAllocationConfigParam nsxModel.ServiceRouterAllocationConfig) (nsxModel.LogicalRouter, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersReallocateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersReallocateInputType(), typeConverter)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	sv.AddStructField("ServiceRouterAllocationConfig", serviceRouterAllocationConfigParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouter
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "reallocate", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouter
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRoutersReallocateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouter), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *logicalRoutersClient) Reprocess(logicalRouterIdParam string) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersReprocessRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersReprocessInputType(), typeConverter)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "reprocess", inputDataValue, executionContext)
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

func (lIface *logicalRoutersClient) Update(logicalRouterIdParam string, logicalRouterParam nsxModel.LogicalRouter) (nsxModel.LogicalRouter, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := logicalRoutersUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(logicalRoutersUpdateInputType(), typeConverter)
	sv.AddStructField("LogicalRouterId", logicalRouterIdParam)
	sv.AddStructField("LogicalRouter", logicalRouterParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LogicalRouter
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.logical_routers", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.LogicalRouter
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LogicalRoutersUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LogicalRouter), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
