// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Cluster
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

type ClusterClient interface {

	// Request one-time backup. The backup will be uploaded using the same server configuration as for automatic backup.
	//
	// @param frameTypeParam Frame type (optional, default to LOCAL_LOCAL_MANAGER)
	// @param siteIdParam Site ID (optional, default to localhost)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Backuptoremote(frameTypeParam *string, siteIdParam *string) error

	// Invoke POST request on target cluster node
	//
	// @param targetNodeIdParam Target node UUID or keyword self (required)
	// @param targetUriParam URI of API to invoke on target node (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(targetNodeIdParam string, targetUriParam string) error

	// Invoke DELETE request on target cluster node
	//
	// @param targetNodeIdParam Target node UUID or keyword self (required)
	// @param targetUriParam URI of API to invoke on target node (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(targetNodeIdParam string, targetUriParam string) error

	// Returns information about the NSX cluster configuration. An NSX cluster has two functions or purposes, commonly referred to as \"roles.\" These two roles are control and management. Each NSX installation has a single cluster. Separate NSX clusters do not share data. In other words, a given data-plane node is attached to only one cluster, not to multiple clusters.
	// @return com.vmware.nsx.model.ClusterConfig
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get() (nsxModel.ClusterConfig, error)

	// Returns information about the specified NSX cluster node.
	//
	// @param nodeIdParam (required)
	// @return com.vmware.nsx.model.ClusterNodeInfo
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get0(nodeIdParam string) (nsxModel.ClusterNodeInfo, error)

	// Invoke GET request on target cluster node
	//
	// @param targetNodeIdParam Target node UUID or keyword self (required)
	// @param targetUriParam URI of API to invoke on target node (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get1(targetNodeIdParam string, targetUriParam string) error

	// Join this node to a NSX Cluster
	//
	// @param joinClusterParametersParam (required)
	// @return com.vmware.nsx.model.ClusterConfiguration
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Joincluster(joinClusterParametersParam nsxModel.JoinClusterParameters) (nsxModel.ClusterConfiguration, error)

	// Detach a node from the Cluster
	//
	// @param nodeIdParam UUID of the node (required)
	// @param forceParam (optional)
	// @param gracefulShutdownParam (optional, default to false)
	// @param ignoreRepositoryIpCheckParam (optional, default to false)
	// @return com.vmware.nsx.model.ClusterConfiguration
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Removenode(nodeIdParam string, forceParam *string, gracefulShutdownParam *string, ignoreRepositoryIpCheckParam *string) (nsxModel.ClusterConfiguration, error)

	// Request one-time inventory summary. The backup will be uploaded using the same server configuration as for an automatic backup.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Summarizeinventorytoremote() error

	// Invoke PUT request on target cluster node
	//
	// @param targetNodeIdParam Target node UUID or keyword self (required)
	// @param targetUriParam URI of API to invoke on target node (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(targetNodeIdParam string, targetUriParam string) error
}

type clusterClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewClusterClient(connector vapiProtocolClient_.Connector) *clusterClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.cluster")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"backuptoremote":             vapiCore_.NewMethodIdentifier(interfaceIdentifier, "backuptoremote"),
		"create":                     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete":                     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                        vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"get_0":                      vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get_0"),
		"get_1":                      vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get_1"),
		"joincluster":                vapiCore_.NewMethodIdentifier(interfaceIdentifier, "joincluster"),
		"removenode":                 vapiCore_.NewMethodIdentifier(interfaceIdentifier, "removenode"),
		"summarizeinventorytoremote": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "summarizeinventorytoremote"),
		"update":                     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	cIface := clusterClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &cIface
}

func (cIface *clusterClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := cIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (cIface *clusterClient) Backuptoremote(frameTypeParam *string, siteIdParam *string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterBackuptoremoteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterBackuptoremoteInputType(), typeConverter)
	sv.AddStructField("FrameType", frameTypeParam)
	sv.AddStructField("SiteId", siteIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "backuptoremote", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *clusterClient) Create(targetNodeIdParam string, targetUriParam string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterCreateInputType(), typeConverter)
	sv.AddStructField("TargetNodeId", targetNodeIdParam)
	sv.AddStructField("TargetUri", targetUriParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "create", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *clusterClient) Delete(targetNodeIdParam string, targetUriParam string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterDeleteInputType(), typeConverter)
	sv.AddStructField("TargetNodeId", targetNodeIdParam)
	sv.AddStructField("TargetUri", targetUriParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *clusterClient) Get() (nsxModel.ClusterConfig, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterGetInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ClusterConfig
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.ClusterConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ClusterGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ClusterConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *clusterClient) Get0(nodeIdParam string) (nsxModel.ClusterNodeInfo, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterGet0RestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterGet0InputType(), typeConverter)
	sv.AddStructField("NodeId", nodeIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ClusterNodeInfo
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "get_0", inputDataValue, executionContext)
	var emptyOutput nsxModel.ClusterNodeInfo
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ClusterGet0OutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ClusterNodeInfo), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *clusterClient) Get1(targetNodeIdParam string, targetUriParam string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterGet1RestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterGet1InputType(), typeConverter)
	sv.AddStructField("TargetNodeId", targetNodeIdParam)
	sv.AddStructField("TargetUri", targetUriParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "get_1", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *clusterClient) Joincluster(joinClusterParametersParam nsxModel.JoinClusterParameters) (nsxModel.ClusterConfiguration, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterJoinclusterRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterJoinclusterInputType(), typeConverter)
	sv.AddStructField("JoinClusterParameters", joinClusterParametersParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ClusterConfiguration
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "joincluster", inputDataValue, executionContext)
	var emptyOutput nsxModel.ClusterConfiguration
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ClusterJoinclusterOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ClusterConfiguration), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *clusterClient) Removenode(nodeIdParam string, forceParam *string, gracefulShutdownParam *string, ignoreRepositoryIpCheckParam *string) (nsxModel.ClusterConfiguration, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterRemovenodeRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterRemovenodeInputType(), typeConverter)
	sv.AddStructField("NodeId", nodeIdParam)
	sv.AddStructField("Force", forceParam)
	sv.AddStructField("GracefulShutdown", gracefulShutdownParam)
	sv.AddStructField("IgnoreRepositoryIpCheck", ignoreRepositoryIpCheckParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ClusterConfiguration
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "removenode", inputDataValue, executionContext)
	var emptyOutput nsxModel.ClusterConfiguration
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ClusterRemovenodeOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ClusterConfiguration), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *clusterClient) Summarizeinventorytoremote() error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterSummarizeinventorytoremoteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterSummarizeinventorytoremoteInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "summarizeinventorytoremote", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *clusterClient) Update(targetNodeIdParam string, targetUriParam string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := clusterUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(clusterUpdateInputType(), typeConverter)
	sv.AddStructField("TargetNodeId", targetNodeIdParam)
	sv.AddStructField("TargetUri", targetUriParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.cluster", "update", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}
