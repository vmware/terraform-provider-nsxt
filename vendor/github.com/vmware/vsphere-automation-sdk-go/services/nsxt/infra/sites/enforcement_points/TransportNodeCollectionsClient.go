// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: TransportNodeCollections
// Used by client-side stubs.

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type TransportNodeCollectionsClient interface {

	// By deleting transport node collection, we are detaching the transport node profile(TNP) from the compute collection. It has no effect on existing transport nodes. However, new hosts added to the compute collection will no longer be automatically converted to NSX transport node. Detaching TNP from compute collection does not delete TNP.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error

	// Returns transport node collection by id
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @return com.vmware.nsx_policy.model.HostTransportNodeCollection
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) (model.HostTransportNodeCollection, error)

	// This API configures a compute collection for security. In the request body, specify a Transport Node Collection with only the ID of the target compute collection meant for security. Specifically, a Transport Node Profile ID should not be specified. This API will define a system-generated security Transport Node Profile and apply it on the compute collection to create the Transport Node Collection.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @param hostTransportNodeCollectionParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Installformicroseg(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection) error

	// Returns all Transport Node collections
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param clusterMoidParam Managed object ID of cluster in VC (optional)
	// @param computeCollectionIdParam Compute collection id (optional)
	// @param vcInstanceUuidParam UUID for VC deployment (optional)
	// @return com.vmware.nsx_policy.model.HostTransportNodeCollectionListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(siteIdParam string, enforcementpointIdParam string, clusterMoidParam *string, computeCollectionIdParam *string, vcInstanceUuidParam *string) (model.HostTransportNodeCollectionListResult, error)

	// Attach different transport node profile to compute collection by updating transport node collection.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @param hostTransportNodeCollectionParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection) error

	// This API uninstalls NSX applied to the Transport Node Collection with the ID corresponding to the one specified in the request.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Removensx(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error

	// This API is relevant for compute collection on which vLCM is enabled. This API should be invoked to retry the realization of transport node profile on the compute collection. This is useful when profile realization had failed because of error in vLCM. This API has no effect if vLCM is not enabled on the computer collection.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Retryprofilerealization(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error

	// When transport node collection is created the hosts which are part of compute collection will be prepared automatically i.e. NSX Manager attempts to install the NSX components on hosts. Transport nodes for these hosts are created using the configuration specified in transport node profile. Pass apply_profile to false, if you do not want to apply transport node profile on the existing transport node with overridden host flag set and ignore overridden hosts flag is set to true on the transport node profile.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param transportNodeCollectionsIdParam (required)
	// @param hostTransportNodeCollectionParam (required)
	// @param applyProfileParam Indicates if the Transport Node Profile (TNP) configuration should be applied during creation (optional, default to true)
	// @return com.vmware.nsx_policy.model.HostTransportNodeCollection
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionsIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection, applyProfileParam *bool) (model.HostTransportNodeCollection, error)
}

type transportNodeCollectionsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewTransportNodeCollectionsClient(connector client.Connector) *transportNodeCollectionsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete":                  core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                     core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"installformicroseg":      core.NewMethodIdentifier(interfaceIdentifier, "installformicroseg"),
		"list":                    core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":                   core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"removensx":               core.NewMethodIdentifier(interfaceIdentifier, "removensx"),
		"retryprofilerealization": core.NewMethodIdentifier(interfaceIdentifier, "retryprofilerealization"),
		"update":                  core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	tIface := transportNodeCollectionsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &tIface
}

func (tIface *transportNodeCollectionsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := tIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (tIface *transportNodeCollectionsClient) Delete(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsDeleteInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Get(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) (model.HostTransportNodeCollection, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsGetInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNodeCollection
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "get", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNodeCollection
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), transportNodeCollectionsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNodeCollection), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Installformicroseg(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsInstallformicrosegInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	sv.AddStructField("HostTransportNodeCollection", hostTransportNodeCollectionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsInstallformicrosegRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "installformicroseg", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) List(siteIdParam string, enforcementpointIdParam string, clusterMoidParam *string, computeCollectionIdParam *string, vcInstanceUuidParam *string) (model.HostTransportNodeCollectionListResult, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsListInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("ClusterMoid", clusterMoidParam)
	sv.AddStructField("ComputeCollectionId", computeCollectionIdParam)
	sv.AddStructField("VcInstanceUuid", vcInstanceUuidParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNodeCollectionListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "list", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNodeCollectionListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), transportNodeCollectionsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNodeCollectionListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Patch(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsPatchInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	sv.AddStructField("HostTransportNodeCollection", hostTransportNodeCollectionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Removensx(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsRemovensxInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsRemovensxRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "removensx", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Retryprofilerealization(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionIdParam string) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsRetryprofilerealizationInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionId", transportNodeCollectionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsRetryprofilerealizationRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "retryprofilerealization", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportNodeCollectionsClient) Update(siteIdParam string, enforcementpointIdParam string, transportNodeCollectionsIdParam string, hostTransportNodeCollectionParam model.HostTransportNodeCollection, applyProfileParam *bool) (model.HostTransportNodeCollection, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(transportNodeCollectionsUpdateInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("TransportNodeCollectionsId", transportNodeCollectionsIdParam)
	sv.AddStructField("HostTransportNodeCollection", hostTransportNodeCollectionParam)
	sv.AddStructField("ApplyProfile", applyProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNodeCollection
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := transportNodeCollectionsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.transport_node_collections", "update", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNodeCollection
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), transportNodeCollectionsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNodeCollection), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
