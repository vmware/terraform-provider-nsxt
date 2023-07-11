// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: DiscoveredNodes
// Used by client-side stubs.

package fabric

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type DiscoveredNodesClient interface {

	// NSX components are installaed on host and transport node is created with given configurations.
	//
	// @param nodeExtIdParam (required)
	// @param transportNodeParam (required)
	// @param overrideNsxOwnershipParam Override NSX Ownership (optional, default to false)
	// @return com.vmware.nsx.model.TransportNode
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Createtransportnode(nodeExtIdParam string, transportNodeParam nsxModel.TransportNode, overrideNsxOwnershipParam *bool) (nsxModel.TransportNode, error)

	// Returns information about a specific discovered node.
	//
	// @param nodeExtIdParam (required)
	// @return com.vmware.nsx.model.DiscoveredNode
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(nodeExtIdParam string) (nsxModel.DiscoveredNode, error)

	// Returns information about all discovered nodes.
	//
	// @param cmLocalIdParam Local Id of the discovered node in the Compute Manager (optional)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param displayNameParam Display name of discovered node (optional)
	// @param externalIdParam External id of the discovered node, ex. a mo-ref from VC (optional)
	// @param hasParentParam Discovered node has a parent compute collection or is a standalone host (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param ipAddressParam IP address of the discovered node (optional)
	// @param nodeIdParam Id of the fabric node created from the discovered node (optional)
	// @param nodeTypeParam Discovered Node type like HostNode (optional)
	// @param originIdParam Id of the compute manager from where this node was discovered (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param parentComputeCollectionParam External id of the compute collection to which this node belongs (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.DiscoveredNodeListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cmLocalIdParam *string, cursorParam *string, displayNameParam *string, externalIdParam *string, hasParentParam *string, includedFieldsParam *string, ipAddressParam *string, nodeIdParam *string, nodeTypeParam *string, originIdParam *string, pageSizeParam *int64, parentComputeCollectionParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.DiscoveredNodeListResult, error)

	// When transport node profile (TNP) is applied to a cluster, if any validation fails (e.g. VMs running on host) then transport node (TN) is not created. In that case after the required action is taken (e.g. VMs powered off), you can call this API to try to create TN for that discovered node. Do not call this API if Transport Node already exists for the discovered node. In that case use API on transport node. /transport-nodes/<transport-node-id>?action=restore_cluster_config
	//
	// @param nodeExtIdParam (required)
	// @param overrideNsxOwnershipParam Override NSX Ownership (optional, default to false)
	// @return com.vmware.nsx.model.TransportNode
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Reapplyclusterconfig(nodeExtIdParam string, overrideNsxOwnershipParam *bool) (nsxModel.TransportNode, error)
}

type discoveredNodesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewDiscoveredNodesClient(connector vapiProtocolClient_.Connector) *discoveredNodesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.fabric.discovered_nodes")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"createtransportnode":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "createtransportnode"),
		"get":                  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":                 vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"reapplyclusterconfig": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "reapplyclusterconfig"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	dIface := discoveredNodesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &dIface
}

func (dIface *discoveredNodesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := dIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (dIface *discoveredNodesClient) Createtransportnode(nodeExtIdParam string, transportNodeParam nsxModel.TransportNode, overrideNsxOwnershipParam *bool) (nsxModel.TransportNode, error) {
	typeConverter := dIface.connector.TypeConverter()
	executionContext := dIface.connector.NewExecutionContext()
	operationRestMetaData := discoveredNodesCreatetransportnodeRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(discoveredNodesCreatetransportnodeInputType(), typeConverter)
	sv.AddStructField("NodeExtId", nodeExtIdParam)
	sv.AddStructField("TransportNode", transportNodeParam)
	sv.AddStructField("OverrideNsxOwnership", overrideNsxOwnershipParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.TransportNode
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := dIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.discovered_nodes", "createtransportnode", inputDataValue, executionContext)
	var emptyOutput nsxModel.TransportNode
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), DiscoveredNodesCreatetransportnodeOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.TransportNode), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), dIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (dIface *discoveredNodesClient) Get(nodeExtIdParam string) (nsxModel.DiscoveredNode, error) {
	typeConverter := dIface.connector.TypeConverter()
	executionContext := dIface.connector.NewExecutionContext()
	operationRestMetaData := discoveredNodesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(discoveredNodesGetInputType(), typeConverter)
	sv.AddStructField("NodeExtId", nodeExtIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.DiscoveredNode
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := dIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.discovered_nodes", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.DiscoveredNode
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), DiscoveredNodesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.DiscoveredNode), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), dIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (dIface *discoveredNodesClient) List(cmLocalIdParam *string, cursorParam *string, displayNameParam *string, externalIdParam *string, hasParentParam *string, includedFieldsParam *string, ipAddressParam *string, nodeIdParam *string, nodeTypeParam *string, originIdParam *string, pageSizeParam *int64, parentComputeCollectionParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.DiscoveredNodeListResult, error) {
	typeConverter := dIface.connector.TypeConverter()
	executionContext := dIface.connector.NewExecutionContext()
	operationRestMetaData := discoveredNodesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(discoveredNodesListInputType(), typeConverter)
	sv.AddStructField("CmLocalId", cmLocalIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DisplayName", displayNameParam)
	sv.AddStructField("ExternalId", externalIdParam)
	sv.AddStructField("HasParent", hasParentParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("IpAddress", ipAddressParam)
	sv.AddStructField("NodeId", nodeIdParam)
	sv.AddStructField("NodeType", nodeTypeParam)
	sv.AddStructField("OriginId", originIdParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("ParentComputeCollection", parentComputeCollectionParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.DiscoveredNodeListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := dIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.discovered_nodes", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.DiscoveredNodeListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), DiscoveredNodesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.DiscoveredNodeListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), dIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (dIface *discoveredNodesClient) Reapplyclusterconfig(nodeExtIdParam string, overrideNsxOwnershipParam *bool) (nsxModel.TransportNode, error) {
	typeConverter := dIface.connector.TypeConverter()
	executionContext := dIface.connector.NewExecutionContext()
	operationRestMetaData := discoveredNodesReapplyclusterconfigRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(discoveredNodesReapplyclusterconfigInputType(), typeConverter)
	sv.AddStructField("NodeExtId", nodeExtIdParam)
	sv.AddStructField("OverrideNsxOwnership", overrideNsxOwnershipParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.TransportNode
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := dIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.discovered_nodes", "reapplyclusterconfig", inputDataValue, executionContext)
	var emptyOutput nsxModel.TransportNode
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), DiscoveredNodesReapplyclusterconfigOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.TransportNode), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), dIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
