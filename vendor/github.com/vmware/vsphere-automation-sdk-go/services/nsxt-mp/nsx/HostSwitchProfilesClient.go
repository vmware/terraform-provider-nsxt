// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: HostSwitchProfiles
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

type HostSwitchProfilesClient interface {

	// Creates a hostswitch profile. The resource_type is required. For uplink profiles, the teaming and policy parameters are required. By default, the mtu is 1600 and the transport_vlan is 0. The supported MTU range is 1280 through (uplink_mtu_threshold). (uplink_mtu_threshold) is 9000 by default. Range can be extended by modifying (uplink_mtu_threshold) in SwitchingGlobalConfig to the required upper threshold.
	//  This api is now deprecated. Please use new api - PUT policy/api/v1/infra/host-switch-profiles/uplinkProfile1
	//
	// Deprecated: This API element is deprecated.
	//
	// @param baseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.BaseHostSwitchProfile.
	// @return com.vmware.nsx.model.BaseHostSwitchProfile
	// The return value will contain all the properties defined in nsxModel.BaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(baseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Deletes a specified hostswitch profile.
	//  This api is now deprecated. Please use new api - DELETE policy/api/v1/infra/host-switch-profiles/uplinkProfile1
	//
	// Deprecated: This API element is deprecated.
	//
	// @param hostSwitchProfileIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(hostSwitchProfileIdParam string) error

	// Returns information about a specified hostswitch profile.
	//  This api is now deprecated. Please use new api - policy/api/v1/infra/host-switch-profiles/uplinkProfile1
	//
	// Deprecated: This API element is deprecated.
	//
	// @param hostSwitchProfileIdParam (required)
	// @return com.vmware.nsx.model.BaseHostSwitchProfile
	// The return value will contain all the properties defined in nsxModel.BaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(hostSwitchProfileIdParam string) (*vapiData_.StructValue, error)

	// Returns information about the configured hostswitch profiles. Hostswitch profiles define networking policies for hostswitches (sometimes referred to as bridges in OVS). Currently, only uplink teaming is supported. Uplink teaming allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order.
	//  This api is now deprecated. Please use new api - policy/api/v1/infra/host-switch-profiles
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param deploymentTypeParam Supported edge deployment type. (optional)
	// @param hostswitchProfileTypeParam Supported HostSwitch profiles. (optional)
	// @param includeSystemOwnedParam Whether the list result contains system resources (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param nodeTypeParam Fabric node type for which uplink profiles are to be listed (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param uplinkTeamingPolicyNameParam The host switch profile's uplink teaming policy name (optional)
	// @return com.vmware.nsx.model.HostSwitchProfilesListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeSystemOwnedParam *bool, includedFieldsParam *string, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (nsxModel.HostSwitchProfilesListResult, error)

	// Modifies a specified hostswitch profile. The body of the PUT request must include the resource_type. For uplink profiles, the put request must also include teaming parameters. Modifiable attributes include display_name, mtu, and transport_vlan. For uplink teaming policies, uplink_name and policy are also modifiable.
	//  This api is now deprecated. Please use new api - PATCH policy/api/v1/infra/host-switch-profiles/uplinkProfile1
	//
	// Deprecated: This API element is deprecated.
	//
	// @param hostSwitchProfileIdParam (required)
	// @param baseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.BaseHostSwitchProfile.
	// @return com.vmware.nsx.model.BaseHostSwitchProfile
	// The return value will contain all the properties defined in nsxModel.BaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(hostSwitchProfileIdParam string, baseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type hostSwitchProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewHostSwitchProfilesClient(connector vapiProtocolClient_.Connector) *hostSwitchProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.host_switch_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	hIface := hostSwitchProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &hIface
}

func (hIface *hostSwitchProfilesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := hIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (hIface *hostSwitchProfilesClient) Create(baseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesCreateInputType(), typeConverter)
	sv.AddStructField("BaseHostSwitchProfile", baseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx.host_switch_profiles", "create", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Delete(hostSwitchProfileIdParam string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx.host_switch_profiles", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Get(hostSwitchProfileIdParam string) (*vapiData_.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesGetInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx.host_switch_profiles", "get", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeSystemOwnedParam *bool, includedFieldsParam *string, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (nsxModel.HostSwitchProfilesListResult, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DeploymentType", deploymentTypeParam)
	sv.AddStructField("HostswitchProfileType", hostswitchProfileTypeParam)
	sv.AddStructField("IncludeSystemOwned", includeSystemOwnedParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("NodeType", nodeTypeParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("UplinkTeamingPolicyName", uplinkTeamingPolicyNameParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.HostSwitchProfilesListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx.host_switch_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.HostSwitchProfilesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.HostSwitchProfilesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Update(hostSwitchProfileIdParam string, baseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	sv.AddStructField("BaseHostSwitchProfile", baseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx.host_switch_profiles", "update", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
