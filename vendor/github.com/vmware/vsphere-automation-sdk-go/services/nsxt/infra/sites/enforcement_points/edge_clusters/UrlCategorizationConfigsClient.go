// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: UrlCategorizationConfigs
// Used by client-side stubs.

package edge_clusters

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type UrlCategorizationConfigsClient interface {

	// Delete PolicyUrlCategorizationConfig. If deleted, the URL categorization will be disabled for that edge cluster.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) error

	// Gets a PolicyUrlCategorizationConfig. This returns the details of the config like whether the URL categorization is enabled or disabled, the id of the context profiles which are used to filter the categories, and the update frequency of the data from the cloud.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) (model.PolicyUrlCategorizationConfig, error)

	// Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @param policyUrlCategorizationConfigParam (required)
	// @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error)

	// Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @param policyUrlCategorizationConfigParam (required)
	// @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error)
}

type urlCategorizationConfigsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewUrlCategorizationConfigsClient(connector client.Connector) *urlCategorizationConfigsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	uIface := urlCategorizationConfigsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &uIface
}

func (uIface *urlCategorizationConfigsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := uIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (uIface *urlCategorizationConfigsClient) Delete(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) error {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(urlCategorizationConfigsDeleteInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := urlCategorizationConfigsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	uIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Get(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) (model.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(urlCategorizationConfigsGetInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.PolicyUrlCategorizationConfig
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := urlCategorizationConfigsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	uIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "get", inputDataValue, executionContext)
	var emptyOutput model.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), urlCategorizationConfigsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Patch(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(urlCategorizationConfigsPatchInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	sv.AddStructField("PolicyUrlCategorizationConfig", policyUrlCategorizationConfigParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.PolicyUrlCategorizationConfig
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := urlCategorizationConfigsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	uIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "patch", inputDataValue, executionContext)
	var emptyOutput model.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), urlCategorizationConfigsPatchOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Update(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(urlCategorizationConfigsUpdateInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	sv.AddStructField("PolicyUrlCategorizationConfig", policyUrlCategorizationConfigParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.PolicyUrlCategorizationConfig
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := urlCategorizationConfigsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	uIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "update", inputDataValue, executionContext)
	var emptyOutput model.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), urlCategorizationConfigsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
