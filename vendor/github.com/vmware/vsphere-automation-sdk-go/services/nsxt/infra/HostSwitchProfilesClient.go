// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: HostSwitchProfiles
// Used by client-side stubs.

package infra

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsx_policyModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type HostSwitchProfilesClient interface {

	// Deletes a specified hostswitch profile.
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
	//
	// @param hostSwitchProfileIdParam (required)
	// @return com.vmware.nsx_policy.model.PolicyBaseHostSwitchProfile
	// The return value will contain all the properties defined in nsx_policyModel.PolicyBaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(hostSwitchProfileIdParam string) (*vapiData_.StructValue, error)

	// Returns information about the configured hostswitch profiles. Hostswitch profiles define networking policies for hostswitches (sometimes referred to as bridges in OVS). Currently, following profiles are supported. UplinkHostSwitchProfile, LldpHostSwitchProfile, NiocProfile & ExtraConfigHostSwitchProfile. Uplink profile - teaming defined in this profile allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order. Lldp profile - Enable or disable sending LLDP packets NiocProfile - Network I/O Control settings: defines limits, shares and reservations for various host traffic types. ExtraConfig - Vendor specific configuration on HostSwitch, logical switch or logical port
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param deploymentTypeParam Supported edge deployment type. (optional)
	// @param hostswitchProfileTypeParam Supported HostSwitch profiles. (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includeSystemOwnedParam Whether the list result contains system resources (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param maxActiveUplinkCountParam Filter uplink profiles by number of active links in teaming policy. (optional)
	// @param nodeTypeParam Fabric node type for which uplink profiles are to be listed (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param uplinkTeamingPolicyNameParam The host switch profile's uplink teaming policy name (optional)
	// @return com.vmware.nsx_policy.model.PolicyHostSwitchProfilesListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeMarkForDeleteObjectsParam *bool, includeSystemOwnedParam *bool, includedFieldsParam *string, maxActiveUplinkCountParam *int64, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (nsx_policyModel.PolicyHostSwitchProfilesListResult, error)

	// Patch a hostswitch profile. The resource_type is required and needs to be one of the following, UplinkHostSwitchProfile, LldpHostSwitchProfile, NiocProfile & ExtraConfigHostSwitchProfile. Uplink profile - For uplink profiles, the teaming and policy parameters are required. By default, the mtu is 1600 and the transport_vlan is 0. The supported MTU range is 1280 through (uplink_mtu_threshold). uplink_mtu_threshold is 9000 by default. Range can be extended by modifying (uplink_mtu_threshold) in SwitchingGlobalConfig to the required upper threshold. Teaming defined in this profile allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order. Lldp profile - Enable or disable sending LLDP packets NiocProfile - Network I/O Control settings: defines limits, shares and reservations for various host traffic types. ExtraConfig - Vendor specific configuration on HostSwitch, logical switch or logical port
	//
	// @param hostSwitchProfileIdParam (required)
	// @param policyBaseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in nsx_policyModel.PolicyBaseHostSwitchProfile.
	// @return com.vmware.nsx_policy.model.PolicyBaseHostSwitchProfile
	// The return value will contain all the properties defined in nsx_policyModel.PolicyBaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Create or update a hostswitch profile. The resource_type is required and needs to be one of the following, UplinkHostSwitchProfile, LldpHostSwitchProfile, NiocProfile & ExtraConfigHostSwitchProfile. Uplink profile - For uplink profiles, the teaming and policy parameters are required. By default, the mtu is 1600 and the transport_vlan is 0. The supported MTU range is 1280 through (uplink_mtu_threshold). uplink_mtu_threshold is 9000 by default. Range can be extended by modifying (uplink_mtu_threshold) in SwitchingGlobalConfig to the required upper threshold. Teaming defined in this profile allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order. Lldp profile - Enable or disable sending LLDP packets NiocProfile - Network I/O Control settings: defines limits, shares and reservations for various host traffic types. ExtraConfig - Vendor specific configuration on HostSwitch, logical switch or logical port
	//
	// @param hostSwitchProfileIdParam (required)
	// @param policyBaseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in nsx_policyModel.PolicyBaseHostSwitchProfile.
	// @return com.vmware.nsx_policy.model.PolicyBaseHostSwitchProfile
	// The return value will contain all the properties defined in nsx_policyModel.PolicyBaseHostSwitchProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type hostSwitchProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewHostSwitchProfilesClient(connector vapiProtocolClient_.Connector) *hostSwitchProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.host_switch_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "patch"),
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

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "delete", inputDataValue, executionContext)
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

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "get", inputDataValue, executionContext)
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

func (hIface *hostSwitchProfilesClient) List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeMarkForDeleteObjectsParam *bool, includeSystemOwnedParam *bool, includedFieldsParam *string, maxActiveUplinkCountParam *int64, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (nsx_policyModel.PolicyHostSwitchProfilesListResult, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DeploymentType", deploymentTypeParam)
	sv.AddStructField("HostswitchProfileType", hostswitchProfileTypeParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludeSystemOwned", includeSystemOwnedParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("MaxActiveUplinkCount", maxActiveUplinkCountParam)
	sv.AddStructField("NodeType", nodeTypeParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("UplinkTeamingPolicyName", uplinkTeamingPolicyNameParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsx_policyModel.PolicyHostSwitchProfilesListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsx_policyModel.PolicyHostSwitchProfilesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsx_policyModel.PolicyHostSwitchProfilesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Patch(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesPatchRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesPatchInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	sv.AddStructField("PolicyBaseHostSwitchProfile", policyBaseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "patch", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), HostSwitchProfilesPatchOutputType())
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

func (hIface *hostSwitchProfilesClient) Update(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	operationRestMetaData := hostSwitchProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(hostSwitchProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	sv.AddStructField("PolicyBaseHostSwitchProfile", policyBaseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "update", inputDataValue, executionContext)
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
