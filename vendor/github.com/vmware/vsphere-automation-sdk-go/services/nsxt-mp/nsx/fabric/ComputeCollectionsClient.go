// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: ComputeCollections
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

type ComputeCollectionsClient interface {

	// Perform action specific to NSX on the compute-collection. cc-ext-id should be of type VC_Cluster.
	//
	// @param ccExtIdParam (required)
	// @param actionParam Supported actions on compute-collection (optional)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(ccExtIdParam string, actionParam *string) error

	// Returns information about a specific compute collection.
	//
	// @param ccExtIdParam (required)
	// @return com.vmware.nsx.model.ComputeCollection
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(ccExtIdParam string) (nsxModel.ComputeCollection, error)

	// Returns information about all compute collections.
	//
	// @param cmLocalIdParam Local Id of the compute collection in the Compute Manager (optional)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param discoveredNodeIdParam Id of the discovered node which belongs to this Compute Collection (optional)
	// @param displayNameParam Name of the ComputeCollection in source compute manager (optional)
	// @param externalIdParam External ID of the ComputeCollection in the source Compute manager, e.g. mo-ref in VC (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param nodeIdParam Id of the fabric node created from a discovered node belonging to this Compute Collection (optional)
	// @param originIdParam Id of the compute manager from where this Compute Collection was discovered (optional)
	// @param originTypeParam ComputeCollection type like VC_Cluster. Here the Compute Manager type prefix would help in differentiating similar named Compute Collection types from different Compute Managers (optional)
	// @param ownerIdParam Id of the owner of compute collection in the Compute Manager (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.ComputeCollectionListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cmLocalIdParam *string, cursorParam *string, discoveredNodeIdParam *string, displayNameParam *string, externalIdParam *string, includedFieldsParam *string, nodeIdParam *string, originIdParam *string, originTypeParam *string, ownerIdParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.ComputeCollectionListResult, error)
}

type computeCollectionsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewComputeCollectionsClient(connector vapiProtocolClient_.Connector) *computeCollectionsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.fabric.compute_collections")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	cIface := computeCollectionsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &cIface
}

func (cIface *computeCollectionsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := cIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (cIface *computeCollectionsClient) Create(ccExtIdParam string, actionParam *string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := computeCollectionsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(computeCollectionsCreateInputType(), typeConverter)
	sv.AddStructField("CcExtId", ccExtIdParam)
	sv.AddStructField("Action", actionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.compute_collections", "create", inputDataValue, executionContext)
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

func (cIface *computeCollectionsClient) Get(ccExtIdParam string) (nsxModel.ComputeCollection, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := computeCollectionsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(computeCollectionsGetInputType(), typeConverter)
	sv.AddStructField("CcExtId", ccExtIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ComputeCollection
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.compute_collections", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.ComputeCollection
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ComputeCollectionsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ComputeCollection), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *computeCollectionsClient) List(cmLocalIdParam *string, cursorParam *string, discoveredNodeIdParam *string, displayNameParam *string, externalIdParam *string, includedFieldsParam *string, nodeIdParam *string, originIdParam *string, originTypeParam *string, ownerIdParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.ComputeCollectionListResult, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	operationRestMetaData := computeCollectionsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(computeCollectionsListInputType(), typeConverter)
	sv.AddStructField("CmLocalId", cmLocalIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DiscoveredNodeId", discoveredNodeIdParam)
	sv.AddStructField("DisplayName", displayNameParam)
	sv.AddStructField("ExternalId", externalIdParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("NodeId", nodeIdParam)
	sv.AddStructField("OriginId", originIdParam)
	sv.AddStructField("OriginType", originTypeParam)
	sv.AddStructField("OwnerId", ownerIdParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.ComputeCollectionListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.compute_collections", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.ComputeCollectionListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), ComputeCollectionsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.ComputeCollectionListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
