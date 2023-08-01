// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: EdgeClusters
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

type EdgeClustersClient interface {

	// Creates a new edge cluster. It only supports homogeneous members. The TransportNodes backed by EdgeNode are only allowed in cluster members. DeploymentType (VIRTUAL_MACHINE|PHYSICAL_MACHINE) of these EdgeNodes is recommended to be the same. EdgeCluster supports members of different deployment types.
	//
	// @param edgeClusterParam (required)
	// @return com.vmware.nsx.model.EdgeCluster
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(edgeClusterParam nsxModel.EdgeCluster) (nsxModel.EdgeCluster, error)

	// Deletes the specified edge cluster.
	//
	// @param edgeClusterIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(edgeClusterIdParam string) error

	// Returns information about the specified edge cluster.
	//
	// @param edgeClusterIdParam (required)
	// @return com.vmware.nsx.model.EdgeCluster
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(edgeClusterIdParam string) (nsxModel.EdgeCluster, error)

	// Returns information about the configured edge clusters, which enable you to group together transport nodes of the type EdgeNode and apply fabric profiles to all members of the edge cluster. Each edge node can participate in only one edge cluster.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.EdgeClusterListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.EdgeClusterListResult, error)

	// Relocate auto allocated service contexts from edge node at given index. For API to perform relocate and remove action the edge node at given index must only have auto allocated service contexts. If any manually allocated service context is present on the edge cluster member, then the task will not be performed. Also, it is recommended to move edge node for which relocate and remove action is being performed into maintenance mode, before executing the API. If edge is not not moved into maintenance mode, then API will move edge node into maintenance mode before performing the actual relocate and remove task.To maintain high availability, Edge cluster should have at least two healthy edge nodes for relocate and removal. Once relocate action is performed successfully, the edge node will be removed from the edge cluster.
	//
	// @param edgeClusterIdParam (required)
	// @param edgeClusterMemberIndexParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Relocateandremove(edgeClusterIdParam string, edgeClusterMemberIndexParam nsxModel.EdgeClusterMemberIndex) error

	// Replace the transport node in the specified member of the edge-cluster. This is a disruptive action. This will move all the LogicalRouterPorts(uplink and routerLink) host on the old transport_node to the new transport_node. The transportNode cannot be present in another member of any edgeClusters.
	//
	// @param edgeClusterIdParam (required)
	// @param edgeClusterMemberTransportNodeParam (required)
	// @return com.vmware.nsx.model.EdgeCluster
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Replacetransportnode(edgeClusterIdParam string, edgeClusterMemberTransportNodeParam nsxModel.EdgeClusterMemberTransportNode) (nsxModel.EdgeCluster, error)

	// Modifies the specified edge cluster. Modifiable parameters include the description, display_name, transport-node-id. If the optional fabric_profile_binding is included, resource_type and profile_id are required. User should do a GET on the edge-cluster and obtain the payload and retain the member_index of the existing members as returning in the GET output. For new member additions, the member_index cannot be defined by the user, user can read the system allocated index to the new member in the output of this API call or by doing a GET call. User cannot use this PUT api to replace the transport_node of an existing member because this is a disruption action, we have exposed a explicit API for doing so, refer to \"ReplaceEdgeClusterMemberTransportNode\" EdgeCluster only supports homogeneous members. The TransportNodes backed by EdgeNode are only allowed in cluster members. DeploymentType (VIRTUAL_MACHINE|PHYSICAL_MACHINE) of these EdgeNodes is recommended to be the same. EdgeCluster supports members of different deployment types.
	//
	// @param edgeClusterIdParam (required)
	// @param edgeClusterParam (required)
	// @return com.vmware.nsx.model.EdgeCluster
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(edgeClusterIdParam string, edgeClusterParam nsxModel.EdgeCluster) (nsxModel.EdgeCluster, error)
}

type edgeClustersClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewEdgeClustersClient(connector vapiProtocolClient_.Connector) *edgeClustersClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.edge_clusters")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create":               vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete":               vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":                 vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"relocateandremove":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "relocateandremove"),
		"replacetransportnode": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "replacetransportnode"),
		"update":               vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	eIface := edgeClustersClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &eIface
}

func (eIface *edgeClustersClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := eIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (eIface *edgeClustersClient) Create(edgeClusterParam nsxModel.EdgeCluster) (nsxModel.EdgeCluster, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersCreateInputType(), typeConverter)
	sv.AddStructField("EdgeCluster", edgeClusterParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EdgeCluster
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.EdgeCluster
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EdgeClustersCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EdgeCluster), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *edgeClustersClient) Delete(edgeClusterIdParam string) error {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersDeleteInputType(), typeConverter)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (eIface *edgeClustersClient) Get(edgeClusterIdParam string) (nsxModel.EdgeCluster, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersGetInputType(), typeConverter)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EdgeCluster
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.EdgeCluster
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EdgeClustersGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EdgeCluster), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *edgeClustersClient) List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.EdgeClusterListResult, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EdgeClusterListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.EdgeClusterListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EdgeClustersListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EdgeClusterListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *edgeClustersClient) Relocateandremove(edgeClusterIdParam string, edgeClusterMemberIndexParam nsxModel.EdgeClusterMemberIndex) error {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersRelocateandremoveRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersRelocateandremoveInputType(), typeConverter)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("EdgeClusterMemberIndex", edgeClusterMemberIndexParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "relocateandremove", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (eIface *edgeClustersClient) Replacetransportnode(edgeClusterIdParam string, edgeClusterMemberTransportNodeParam nsxModel.EdgeClusterMemberTransportNode) (nsxModel.EdgeCluster, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersReplacetransportnodeRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersReplacetransportnodeInputType(), typeConverter)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("EdgeClusterMemberTransportNode", edgeClusterMemberTransportNodeParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EdgeCluster
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "replacetransportnode", inputDataValue, executionContext)
	var emptyOutput nsxModel.EdgeCluster
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EdgeClustersReplacetransportnodeOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EdgeCluster), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (eIface *edgeClustersClient) Update(edgeClusterIdParam string, edgeClusterParam nsxModel.EdgeCluster) (nsxModel.EdgeCluster, error) {
	typeConverter := eIface.connector.TypeConverter()
	executionContext := eIface.connector.NewExecutionContext()
	operationRestMetaData := edgeClustersUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(edgeClustersUpdateInputType(), typeConverter)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("EdgeCluster", edgeClusterParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.EdgeCluster
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := eIface.connector.GetApiProvider().Invoke("com.vmware.nsx.edge_clusters", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.EdgeCluster
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), EdgeClustersUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.EdgeCluster), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), eIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
