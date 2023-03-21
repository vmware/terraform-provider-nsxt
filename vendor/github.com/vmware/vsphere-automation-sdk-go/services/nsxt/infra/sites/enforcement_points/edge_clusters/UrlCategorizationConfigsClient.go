// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: UrlCategorizationConfigs
// Used by client-side stubs.

package edge_clusters

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type UrlCategorizationConfigsClient interface {

	// Delete PolicyUrlCategorizationConfig. If deleted, the URL categorization will be disabled for that edge cluster.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	//
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
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) (nsx_policyModel.PolicyUrlCategorizationConfig, error)

	// Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @param policyUrlCategorizationConfigParam (required)
	// @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam nsx_policyModel.PolicyUrlCategorizationConfig) (nsx_policyModel.PolicyUrlCategorizationConfig, error)

	// Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
	//
	// @param siteIdParam (required)
	// @param enforcementPointIdParam (required)
	// @param edgeClusterIdParam (required)
	// @param urlCategorizationConfigIdParam (required)
	// @param policyUrlCategorizationConfigParam (required)
	// @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam nsx_policyModel.PolicyUrlCategorizationConfig) (nsx_policyModel.PolicyUrlCategorizationConfig, error)
}

type urlCategorizationConfigsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewUrlCategorizationConfigsClient(connector vapiProtocolClient_.Connector) *urlCategorizationConfigsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"patch":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	uIface := urlCategorizationConfigsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &uIface
}

func (uIface *urlCategorizationConfigsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := uIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (uIface *urlCategorizationConfigsClient) Delete(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) error {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := urlCategorizationConfigsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(urlCategorizationConfigsDeleteInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Get(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) (nsx_policyModel.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := urlCategorizationConfigsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(urlCategorizationConfigsGetInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "get", inputDataValue, executionContext)
	var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UrlCategorizationConfigsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsx_policyModel.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Patch(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam nsx_policyModel.PolicyUrlCategorizationConfig) (nsx_policyModel.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := urlCategorizationConfigsPatchRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(urlCategorizationConfigsPatchInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	sv.AddStructField("PolicyUrlCategorizationConfig", policyUrlCategorizationConfigParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "patch", inputDataValue, executionContext)
	var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UrlCategorizationConfigsPatchOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsx_policyModel.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *urlCategorizationConfigsClient) Update(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam nsx_policyModel.PolicyUrlCategorizationConfig) (nsx_policyModel.PolicyUrlCategorizationConfig, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := urlCategorizationConfigsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(urlCategorizationConfigsUpdateInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementPointId", enforcementPointIdParam)
	sv.AddStructField("EdgeClusterId", edgeClusterIdParam)
	sv.AddStructField("UrlCategorizationConfigId", urlCategorizationConfigIdParam)
	sv.AddStructField("PolicyUrlCategorizationConfig", policyUrlCategorizationConfigParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.edge_clusters.url_categorization_configs", "update", inputDataValue, executionContext)
	var emptyOutput nsx_policyModel.PolicyUrlCategorizationConfig
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UrlCategorizationConfigsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsx_policyModel.PolicyUrlCategorizationConfig), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
