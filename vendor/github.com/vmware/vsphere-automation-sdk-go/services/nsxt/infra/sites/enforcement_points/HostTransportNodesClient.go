// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: HostTransportNodes
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

type HostTransportNodesClient interface {

	// Deletes the specified transport node. Query param force can be used to force delete the host nodes. It also removes the specified host node from system. If unprepare_host option is set to false, then host will be deleted without uninstalling the NSX components from the host.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	// @param unprepareHostParam Uninstall NSX components from host while deleting (optional, default to true)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, forceParam *bool, unprepareHostParam *bool) error

	// Returns information about a specified transport node.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @return com.vmware.nsx_policy.model.HostTransportNode
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) (model.HostTransportNode, error)

	// Returns information about all host transport nodes along with underlying host details. A transport node is a host that contains hostswitches. A hostswitch can have virtual machines connected to them. Because each transport node has hostswitches, transport nodes can also have virtual tunnel endpoints, which means that they can be part of the overlay.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param discoveredNodeIdParam discovered node id (optional)
	// @param inMaintenanceModeParam maintenance mode flag (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param nodeIpParam Transport node IP address (optional)
	// @param nodeTypesParam a list of node types separated by comma or a single type (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @param transportZonePathParam Transport zone path (optional)
	// @return com.vmware.nsx_policy.model.HostTransportNodeListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(siteIdParam string, enforcementpointIdParam string, cursorParam *string, discoveredNodeIdParam *string, inMaintenanceModeParam *bool, includedFieldsParam *string, nodeIpParam *string, nodeTypesParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, transportZonePathParam *string) (model.HostTransportNodeListResult, error)

	// Migrates all NVDS to VDS on given TransportNode. Upgrade precheck apis should have been run prior to invoking this API on transport node and a migration topology should be created. Please refer to Migration guide for details about migration APIs.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @param skipMaintmodeParam Skip Maintenance mode check (optional, default to false)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Migratetovds(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, skipMaintmodeParam *bool) error

	// Transport nodes are hypervisor hosts that will participate in an NSX-T overlay. For a hypervisor host, this means that it hosts VMs that will communicate over NSX-T logical switches. This API creates transport node for a host node (hypervisor) in the transport network. When you run this command for a host, NSX Manager attempts to install the NSX kernel modules, which are packaged as VIB, RPM, or DEB files. For the installation to succeed, you must provide the host login credentials and the host thumbprint. To get the ESXi host thumbprint, SSH to the host and run the **openssl x509 -in /etc/vmware/ssl/rui.crt -fingerprint -sha256 -noout** command. To generate host key thumbprint using SHA-256 algorithm please follow the steps below. Log into the host, making sure that the connection is not vulnerable to a man in the middle attack. Check whether a public key already exists. Host public key is generally located at '/etc/ssh/ssh_host_rsa_key.pub'. If the key is not present then generate a new key by running the following command and follow the instructions. **ssh-keygen -t rsa** Now generate a SHA256 hash of the key using the following command. Please make sure to pass the appropriate file name if the public key is stored with a different file name other than the default 'id_rsa.pub'. **awk '{print $2}' id_rsa.pub | base64 -d | sha256sum -b | sed 's/ .\*$//' | xxd -r -p | base64** Additional documentation on creating a transport node can be found in the NSX-T Installation Guide. In order for the transport node to forward packets, the host_switch_spec property must be specified. Host switches (called bridges in OVS on KVM hypervisors) are the individual switches within the host virtual switch. Virtual machines are connected to the host switches. When creating a transport node, you need to specify if the host switches are already manually preconfigured on the node, or if NSX should create and manage the host switches. You specify this choice by the type of host switches you pass in the host_switch_spec property of the TransportNode request payload. For a KVM host, you can preconfigure the host switch, or you can have NSX Manager perform the configuration. For an ESXi host NSX Manager always configures the host switch. To preconfigure the host switches on a KVM host, pass an array of PreconfiguredHostSwitchSpec objects that describes those host switches. In the current NSX-T release, only one prefonfigured host switch can be specified. See the PreconfiguredHostSwitchSpec schema definition for documentation on the properties that must be provided. Preconfigured host switches are only supported on KVM hosts, not on ESXi hosts. To allow NSX to manage the host switch configuration on KVM hosts, ESXi hosts, pass an array of StandardHostSwitchSpec objects in the host_switch_spec property, and NSX will automatically create host switches with the properties you provide. In the current NSX-T release, up to 16 host switches can be automatically managed. See the StandardHostSwitchSpec schema definition for documentation on the properties that must be provided. The request should provide node_deployement_info.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @param hostTransportNodeParam (required)
	// @param esxMgmtIfMigrationDestParam The network ids to which the ESX vmk interfaces will be migrated (optional)
	// @param ifIdParam The ESX vmk interfaces to migrate (optional)
	// @param pingIpParam IP Addresses to ping right after ESX vmk interfaces were migrated. (optional)
	// @param skipValidationParam Whether to skip front-end validation for vmk/vnic/pnic migration (optional, default to false)
	// @param vnicParam The ESX vmk interfaces and/or VM NIC to migrate (optional)
	// @param vnicMigrationDestParam The migration destinations of ESX vmk interfaces and/or VM NIC (optional)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, hostTransportNodeParam model.HostTransportNode, esxMgmtIfMigrationDestParam *string, ifIdParam *string, pingIpParam *string, skipValidationParam *bool, vnicParam *string, vnicMigrationDestParam *string) error

	// A host can be overridden to have different configuration than Transport Node Profile(TNP) on cluster. This action will restore such overridden host back to cluster level TNP. This API can be used in other case. When TNP is applied to a cluster, if any validation fails (e.g. VMs running on host) then existing transport node (TN) is not updated. In that case after the issue is resolved manually (e.g. VMs powered off), you can call this API to update TN as per cluster level TNP.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Restoreclusterconfig(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) error

	// Resync the TransportNode configuration on a host. It is similar to updating the TransportNode with existing configuration, but force synce these configurations to the host (no backend optimizations).
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Resynchostconfig(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) error

	// Transport nodes are hypervisor hosts that will participate in an NSX-T overlay. For a hypervisor host, this means that it hosts VMs that will communicate over NSX-T logical switches. This API creates transport node for a host node (hypervisor) in the transport network. When you run this command for a host, NSX Manager attempts to install the NSX kernel modules, which are packaged as VIB, RPM, or DEB files. For the installation to succeed, you must provide the host login credentials and the host thumbprint. To get the ESXi host thumbprint, SSH to the host and run the **openssl x509 -in /etc/vmware/ssl/rui.crt -fingerprint -sha256 -noout** command. To generate host key thumbprint using SHA-256 algorithm please follow the steps below. Log into the host, making sure that the connection is not vulnerable to a man in the middle attack. Check whether a public key already exists. Host public key is generally located at '/etc/ssh/ssh_host_rsa_key.pub'. If the key is not present then generate a new key by running the following command and follow the instructions. **ssh-keygen -t rsa** Now generate a SHA256 hash of the key using the following command. Please make sure to pass the appropriate file name if the public key is stored with a different file name other than the default 'id_rsa.pub'. **awk '{print $2}' id_rsa.pub | base64 -d | sha256sum -b | sed 's/ .\*$//' | xxd -r -p | base64** Additional documentation on creating a transport node can be found in the NSX-T Installation Guide. In order for the transport node to forward packets, the host_switch_spec property must be specified. Host switches (called bridges in OVS on KVM hypervisors) are the individual switches within the host virtual switch. Virtual machines are connected to the host switches. When creating a transport node, you need to specify if the host switches are already manually preconfigured on the node, or if NSX should create and manage the host switches. You specify this choice by the type of host switches you pass in the host_switch_spec property of the TransportNode request payload. For a KVM host, you can preconfigure the host switch, or you can have NSX Manager perform the configuration. For an ESXi host NSX Manager always configures the host switch. To preconfigure the host switches on a KVM host, pass an array of PreconfiguredHostSwitchSpec objects that describes those host switches. In the current NSX-T release, only one prefonfigured host switch can be specified. See the PreconfiguredHostSwitchSpec schema definition for documentation on the properties that must be provided. Preconfigured host switches are only supported on KVM hosts, not on ESXi hosts. To allow NSX to manage the host switch configuration on KVM hosts, ESXi hosts, pass an array of StandardHostSwitchSpec objects in the host_switch_spec property, and NSX will automatically create host switches with the properties you provide. In the current NSX-T release, up to 16 host switches can be automatically managed. See the StandardHostSwitchSpec schema definition for documentation on the properties that must be provided. The request should provide node_deployement_info.
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @param hostTransportNodeParam (required)
	// @param esxMgmtIfMigrationDestParam The network ids to which the ESX vmk interfaces will be migrated (optional)
	// @param ifIdParam The ESX vmk interfaces to migrate (optional)
	// @param pingIpParam IP Addresses to ping right after ESX vmk interfaces were migrated. (optional)
	// @param skipValidationParam Whether to skip front-end validation for vmk/vnic/pnic migration (optional, default to false)
	// @param vnicParam The ESX vmk interfaces and/or VM NIC to migrate (optional)
	// @param vnicMigrationDestParam The migration destinations of ESX vmk interfaces and/or VM NIC (optional)
	// @return com.vmware.nsx_policy.model.HostTransportNode
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, hostTransportNodeParam model.HostTransportNode, esxMgmtIfMigrationDestParam *string, ifIdParam *string, pingIpParam *string, skipValidationParam *bool, vnicParam *string, vnicMigrationDestParam *string) (model.HostTransportNode, error)

	// Put transport node into maintenance mode or exit from maintenance mode. When HostTransportNode is in maintenance mode no configuration changes are allowed
	//
	// @param siteIdParam (required)
	// @param enforcementpointIdParam (required)
	// @param hostTransportNodeIdParam (required)
	// @param actionParam (optional)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Updatemaintenancemode(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, actionParam *string) error
}

type hostTransportNodesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewHostTransportNodesClient(connector client.Connector) *hostTransportNodesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete":                core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                   core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":                  core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"migratetovds":          core.NewMethodIdentifier(interfaceIdentifier, "migratetovds"),
		"patch":                 core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"restoreclusterconfig":  core.NewMethodIdentifier(interfaceIdentifier, "restoreclusterconfig"),
		"resynchostconfig":      core.NewMethodIdentifier(interfaceIdentifier, "resynchostconfig"),
		"update":                core.NewMethodIdentifier(interfaceIdentifier, "update"),
		"updatemaintenancemode": core.NewMethodIdentifier(interfaceIdentifier, "updatemaintenancemode"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	hIface := hostTransportNodesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &hIface
}

func (hIface *hostTransportNodesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := hIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (hIface *hostTransportNodesClient) Delete(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, forceParam *bool, unprepareHostParam *bool) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesDeleteInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	sv.AddStructField("Force", forceParam)
	sv.AddStructField("UnprepareHost", unprepareHostParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "delete", inputDataValue, executionContext)
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

func (hIface *hostTransportNodesClient) Get(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) (model.HostTransportNode, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesGetInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNode
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "get", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNode
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostTransportNodesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNode), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostTransportNodesClient) List(siteIdParam string, enforcementpointIdParam string, cursorParam *string, discoveredNodeIdParam *string, inMaintenanceModeParam *bool, includedFieldsParam *string, nodeIpParam *string, nodeTypesParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, transportZonePathParam *string) (model.HostTransportNodeListResult, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesListInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DiscoveredNodeId", discoveredNodeIdParam)
	sv.AddStructField("InMaintenanceMode", inMaintenanceModeParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("NodeIp", nodeIpParam)
	sv.AddStructField("NodeTypes", nodeTypesParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	sv.AddStructField("TransportZonePath", transportZonePathParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNodeListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "list", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNodeListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostTransportNodesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNodeListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostTransportNodesClient) Migratetovds(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, skipMaintmodeParam *bool) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesMigratetovdsInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	sv.AddStructField("SkipMaintmode", skipMaintmodeParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesMigratetovdsRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "migratetovds", inputDataValue, executionContext)
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

func (hIface *hostTransportNodesClient) Patch(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, hostTransportNodeParam model.HostTransportNode, esxMgmtIfMigrationDestParam *string, ifIdParam *string, pingIpParam *string, skipValidationParam *bool, vnicParam *string, vnicMigrationDestParam *string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesPatchInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	sv.AddStructField("HostTransportNode", hostTransportNodeParam)
	sv.AddStructField("EsxMgmtIfMigrationDest", esxMgmtIfMigrationDestParam)
	sv.AddStructField("IfId", ifIdParam)
	sv.AddStructField("PingIp", pingIpParam)
	sv.AddStructField("SkipValidation", skipValidationParam)
	sv.AddStructField("Vnic", vnicParam)
	sv.AddStructField("VnicMigrationDest", vnicMigrationDestParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "patch", inputDataValue, executionContext)
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

func (hIface *hostTransportNodesClient) Restoreclusterconfig(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesRestoreclusterconfigInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesRestoreclusterconfigRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "restoreclusterconfig", inputDataValue, executionContext)
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

func (hIface *hostTransportNodesClient) Resynchostconfig(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesResynchostconfigInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesResynchostconfigRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "resynchostconfig", inputDataValue, executionContext)
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

func (hIface *hostTransportNodesClient) Update(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, hostTransportNodeParam model.HostTransportNode, esxMgmtIfMigrationDestParam *string, ifIdParam *string, pingIpParam *string, skipValidationParam *bool, vnicParam *string, vnicMigrationDestParam *string) (model.HostTransportNode, error) {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesUpdateInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	sv.AddStructField("HostTransportNode", hostTransportNodeParam)
	sv.AddStructField("EsxMgmtIfMigrationDest", esxMgmtIfMigrationDestParam)
	sv.AddStructField("IfId", ifIdParam)
	sv.AddStructField("PingIp", pingIpParam)
	sv.AddStructField("SkipValidation", skipValidationParam)
	sv.AddStructField("Vnic", vnicParam)
	sv.AddStructField("VnicMigrationDest", vnicMigrationDestParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.HostTransportNode
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "update", inputDataValue, executionContext)
	var emptyOutput model.HostTransportNode
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), hostTransportNodesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.HostTransportNode), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), hIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (hIface *hostTransportNodesClient) Updatemaintenancemode(siteIdParam string, enforcementpointIdParam string, hostTransportNodeIdParam string, actionParam *string) error {
	typeConverter := hIface.connector.TypeConverter()
	executionContext := hIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(hostTransportNodesUpdatemaintenancemodeInputType(), typeConverter)
	sv.AddStructField("SiteId", siteIdParam)
	sv.AddStructField("EnforcementpointId", enforcementpointIdParam)
	sv.AddStructField("HostTransportNodeId", hostTransportNodeIdParam)
	sv.AddStructField("Action", actionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := hostTransportNodesUpdatemaintenancemodeRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	hIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := hIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.sites.enforcement_points.host_transport_nodes", "updatemaintenancemode", inputDataValue, executionContext)
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
