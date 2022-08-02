// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: HostSwitchProfiles
// Used by client-side stubs.

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type HostSwitchProfilesClient interface {

	// Deletes a specified hostswitch profile.
	//
	// @param hostSwitchProfileIdParam (required)
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
	// The return value will contain all the properties defined in model.PolicyBaseHostSwitchProfile.
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(hostSwitchProfileIdParam string) (*data.StructValue, error)

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
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeMarkForDeleteObjectsParam *bool, includeSystemOwnedParam *bool, includedFieldsParam *string, maxActiveUplinkCountParam *int64, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (model.PolicyHostSwitchProfilesListResult, error)

	// Patch a hostswitch profile. The resource_type is required and needs to be one of the following, UplinkHostSwitchProfile, LldpHostSwitchProfile, NiocProfile & ExtraConfigHostSwitchProfile. Uplink profile - For uplink profiles, the teaming and policy parameters are required. By default, the mtu is 1600 and the transport_vlan is 0. The supported MTU range is 1280 through (uplink_mtu_threshold). uplink_mtu_threshold is 9000 by default. Range can be extended by modifying (uplink_mtu_threshold) in SwitchingGlobalConfig to the required upper threshold. Teaming defined in this profile allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order. Lldp profile - Enable or disable sending LLDP packets NiocProfile - Network I/O Control settings: defines limits, shares and reservations for various host traffic types. ExtraConfig - Vendor specific configuration on HostSwitch, logical switch or logical port
	//
	// @param hostSwitchProfileIdParam (required)
	// @param policyBaseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in model.PolicyBaseHostSwitchProfile.
	// @return com.vmware.nsx_policy.model.PolicyBaseHostSwitchProfile
	// The return value will contain all the properties defined in model.PolicyBaseHostSwitchProfile.
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *data.StructValue) (*data.StructValue, error)

	// Create or update a hostswitch profile. The resource_type is required and needs to be one of the following, UplinkHostSwitchProfile, LldpHostSwitchProfile, NiocProfile & ExtraConfigHostSwitchProfile. Uplink profile - For uplink profiles, the teaming and policy parameters are required. By default, the mtu is 1600 and the transport_vlan is 0. The supported MTU range is 1280 through (uplink_mtu_threshold). uplink_mtu_threshold is 9000 by default. Range can be extended by modifying (uplink_mtu_threshold) in SwitchingGlobalConfig to the required upper threshold. Teaming defined in this profile allows NSX to load balance traffic across different physical NICs (PNICs) on the hypervisor hosts. Multiple teaming policies are supported, including LACP active, LACP passive, load balancing based on source ID, and failover order. Lldp profile - Enable or disable sending LLDP packets NiocProfile - Network I/O Control settings: defines limits, shares and reservations for various host traffic types. ExtraConfig - Vendor specific configuration on HostSwitch, logical switch or logical port
	//
	// @param hostSwitchProfileIdParam (required)
	// @param policyBaseHostSwitchProfileParam (required)
	// The parameter must contain all the properties defined in model.PolicyBaseHostSwitchProfile.
	// @return com.vmware.nsx_policy.model.PolicyBaseHostSwitchProfile
	// The return value will contain all the properties defined in model.PolicyBaseHostSwitchProfile.
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *data.StructValue) (*data.StructValue, error)
}

type hostSwitchProfilesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewHostSwitchProfilesClient(connector client.Connector) *hostSwitchProfilesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.host_switch_profiles")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	hIface := hostSwitchProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &hIface
}

func (hIface *hostSwitchProfilesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := hIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (hIface *hostSwitchProfilesClient) Delete(hostSwitchProfileIdParam string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostSwitchProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostSwitchProfilesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Get(hostSwitchProfileIdParam string) (*data.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostSwitchProfilesGetInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *data.StructValue
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostSwitchProfilesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "get", inputDataValue, executionContext)
	var emptyOutput *data.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostSwitchProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(*data.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) List(cursorParam *string, deploymentTypeParam *string, hostswitchProfileTypeParam *string, includeMarkForDeleteObjectsParam *bool, includeSystemOwnedParam *bool, includedFieldsParam *string, maxActiveUplinkCountParam *int64, nodeTypeParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, uplinkTeamingPolicyNameParam *string) (model.PolicyHostSwitchProfilesListResult, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostSwitchProfilesListInputType(), typeConverter)
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
		var emptyOutput model.PolicyHostSwitchProfilesListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostSwitchProfilesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "list", inputDataValue, executionContext)
	var emptyOutput model.PolicyHostSwitchProfilesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostSwitchProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.PolicyHostSwitchProfilesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Patch(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *data.StructValue) (*data.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostSwitchProfilesPatchInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	sv.AddStructField("PolicyBaseHostSwitchProfile", policyBaseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *data.StructValue
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostSwitchProfilesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "patch", inputDataValue, executionContext)
	var emptyOutput *data.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostSwitchProfilesPatchOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(*data.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostSwitchProfilesClient) Update(hostSwitchProfileIdParam string, policyBaseHostSwitchProfileParam *data.StructValue) (*data.StructValue, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostSwitchProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("HostSwitchProfileId", hostSwitchProfileIdParam)
	sv.AddStructField("PolicyBaseHostSwitchProfile", policyBaseHostSwitchProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *data.StructValue
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostSwitchProfilesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.host_switch_profiles", "update", inputDataValue, executionContext)
	var emptyOutput *data.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostSwitchProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(*data.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
