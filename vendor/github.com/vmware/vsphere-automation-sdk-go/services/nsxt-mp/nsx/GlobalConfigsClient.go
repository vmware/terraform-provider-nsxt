// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: GlobalConfigs
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

type GlobalConfigsClient interface {

	// Returns global configurations that belong to the config type. This rest routine is deprecated, and will be removed after a year.
	//
	//  Use the following Policy APIs for different global configs -
	//  OperationCollectorGlobalConfig GET /policy/api/v1/infra/ops-global-config
	//  RoutingGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  SecurityGlobalConfig GET /policy/api/v1/infra/security-global-config
	//  IdsGlobalConfig GET /policy/api/v1/infra/settings/firewall/security/intrusion-services
	//  FipsGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  SwitchingGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  FirewallGlobalConfig GET policy/api/v1/infra/settings/firewall/security
	//
	// @param configTypeParam (required)
	// @return com.vmware.nsx.model.GlobalConfigs
	// The return value will contain all the properties defined in nsxModel.GlobalConfigs.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(configTypeParam string) (*vapiData_.StructValue, error)

	// Returns global configurations of a NSX domain grouped by the config types. These global configurations are valid across NSX domain for their respective types unless they are overridden by a more granular configurations. This rest routine is deprecated, and will be removed after a year.
	//
	//  The list API is disaggregated to feature verticals.
	//  Use the following Policy APIs for different global configs -
	//  OperationCollectorGlobalConfig GET /policy/api/v1/infra/ops-global-config
	//  RoutingGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  SecurityGlobalConfig GET /policy/api/v1/infra/security-global-config
	//  IdsGlobalConfig GET /policy/api/v1/infra/settings/firewall/security/intrusion-services
	//  FipsGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  SwitchingGlobalConfig GET /policy/api/v1/infra/connectivity-global-config
	//  FirewallGlobalConfig GET policy/api/v1/infra/settings/firewall/security
	// @return com.vmware.nsx.model.GlobalConfigsListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List() (nsxModel.GlobalConfigsListResult, error)

	// It is similar to update global configurations but this request would trigger update even if the configs are unmodified. However, the realization of the new configurations is config-type specific. Refer to config-type specific documentation for details about the configuration push state. This rest routine is deprecated, and will be removed after a year.
	//
	//  Use the following Policy APIs for different global configs -
	//  OperationCollectorGlobalConfig PUT /policy/api/v1/infra/ops-global-config
	//  RoutingGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  SecurityGlobalConfig PUT /policy/api/v1/infra/security-global-config
	//  IdsGlobalConfig PUT /policy/api/v1/infra/settings/firewall/security/intrusion-services
	//  FipsGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  SwitchingGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  FirewallGlobalConfig PUT /policy/api/v1/infra/settings/firewall/security
	//
	// @param configTypeParam (required)
	// @param globalConfigsParam (required)
	// The parameter must contain all the properties defined in nsxModel.GlobalConfigs.
	// @return com.vmware.nsx.model.GlobalConfigs
	// The return value will contain all the properties defined in nsxModel.GlobalConfigs.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Resyncconfig(configTypeParam string, globalConfigsParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Updates global configurations that belong to a config type. The request must include the updated values along with the unmodified values. The values that are updated(different) would trigger update to config-type specific state. However, the realization of the new configurations is config-type specific. Refer to config-type specific documentation for details about the config- uration push state. This rest routine is deprecated, and will be removed after a year.
	//
	//  Use the following Policy APIs for different global configs -
	//  OperationCollectorGlobalConfig PUT /policy/api/v1/infra/ops-global-config
	//  RoutingGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  SecurityGlobalConfig PUT /policy/api/v1/infra/security-global-config
	//  IdsGlobalConfig PUT /policy/api/v1/infra/settings/firewall/security/intrusion-services
	//  FipsGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  SwitchingGlobalConfig PUT /policy/api/v1/infra/connectivity-global-config
	//  FirewallGlobalConfig PUT /policy/api/v1/infra/settings/firewall/security
	//
	// @param configTypeParam (required)
	// @param globalConfigsParam (required)
	// The parameter must contain all the properties defined in nsxModel.GlobalConfigs.
	// @return com.vmware.nsx.model.GlobalConfigs
	// The return value will contain all the properties defined in nsxModel.GlobalConfigs.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(configTypeParam string, globalConfigsParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type globalConfigsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewGlobalConfigsClient(connector vapiProtocolClient_.Connector) *globalConfigsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.global_configs")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"get":          vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":         vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"resyncconfig": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "resyncconfig"),
		"update":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	gIface := globalConfigsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &gIface
}

func (gIface *globalConfigsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := gIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (gIface *globalConfigsClient) Get(configTypeParam string) (*vapiData_.StructValue, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	operationRestMetaData := globalConfigsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(globalConfigsGetInputType(), typeConverter)
	sv.AddStructField("ConfigType", configTypeParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx.global_configs", "get", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), GlobalConfigsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *globalConfigsClient) List() (nsxModel.GlobalConfigsListResult, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	operationRestMetaData := globalConfigsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(globalConfigsListInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.GlobalConfigsListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx.global_configs", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.GlobalConfigsListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), GlobalConfigsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.GlobalConfigsListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *globalConfigsClient) Resyncconfig(configTypeParam string, globalConfigsParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	operationRestMetaData := globalConfigsResyncconfigRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(globalConfigsResyncconfigInputType(), typeConverter)
	sv.AddStructField("ConfigType", configTypeParam)
	sv.AddStructField("GlobalConfigs", globalConfigsParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx.global_configs", "resyncconfig", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), GlobalConfigsResyncconfigOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *globalConfigsClient) Update(configTypeParam string, globalConfigsParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	operationRestMetaData := globalConfigsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(globalConfigsUpdateInputType(), typeConverter)
	sv.AddStructField("ConfigType", configTypeParam)
	sv.AddStructField("GlobalConfigs", globalConfigsParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx.global_configs", "update", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), GlobalConfigsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
