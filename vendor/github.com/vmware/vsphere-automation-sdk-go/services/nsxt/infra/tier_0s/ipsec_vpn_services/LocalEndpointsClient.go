// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: LocalEndpoints
// Used by client-side stubs.

package ipsec_vpn_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type LocalEndpointsClient interface {

	// Delete IPSec VPN local endpoint for a given ipsec vpn service under Tier-0.
	//
	// @param tier0IdParam (required)
	// @param serviceIdParam (required)
	// @param localEndpointIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(tier0IdParam string, serviceIdParam string, localEndpointIdParam string) error

	// Get IPSec VPN local endpoint for a given ipsec vpn service under Tier-0.
	//
	// @param tier0IdParam (required)
	// @param serviceIdParam (required)
	// @param localEndpointIdParam (required)
	// @return com.vmware.nsx_policy.model.IPSecVpnLocalEndpoint
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(tier0IdParam string, serviceIdParam string, localEndpointIdParam string) (model.IPSecVpnLocalEndpoint, error)

	// Get paginated list of all IPSec VPN local endpoints for a given ipsec vpn service under Tier-0.
	//
	// @param tier0IdParam (required)
	// @param serviceIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.IPSecVpnLocalEndpointListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(tier0IdParam string, serviceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnLocalEndpointListResult, error)

	// Create or patch a custom IPSec VPN local endpoint under Tier-0.
	//
	// @param tier0IdParam (required)
	// @param serviceIdParam (required)
	// @param localEndpointIdParam (required)
	// @param ipSecVpnLocalEndpointParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(tier0IdParam string, serviceIdParam string, localEndpointIdParam string, ipSecVpnLocalEndpointParam model.IPSecVpnLocalEndpoint) error

	// Create or fully replace IPSec VPN local endpoint for a given ipsec vpn service under Tier-0. Revision is optional for creation and required for update.
	//
	// @param tier0IdParam (required)
	// @param serviceIdParam (required)
	// @param localEndpointIdParam (required)
	// @param ipSecVpnLocalEndpointParam (required)
	// @return com.vmware.nsx_policy.model.IPSecVpnLocalEndpoint
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(tier0IdParam string, serviceIdParam string, localEndpointIdParam string, ipSecVpnLocalEndpointParam model.IPSecVpnLocalEndpoint) (model.IPSecVpnLocalEndpoint, error)
}

type localEndpointsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewLocalEndpointsClient(connector client.Connector) *localEndpointsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	lIface := localEndpointsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *localEndpointsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *localEndpointsClient) Delete(tier0IdParam string, serviceIdParam string, localEndpointIdParam string) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(localEndpointsDeleteInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("LocalEndpointId", localEndpointIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := localEndpointsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (lIface *localEndpointsClient) Get(tier0IdParam string, serviceIdParam string, localEndpointIdParam string) (model.IPSecVpnLocalEndpoint, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(localEndpointsGetInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("LocalEndpointId", localEndpointIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPSecVpnLocalEndpoint
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := localEndpointsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints", "get", inputDataValue, executionContext)
	var emptyOutput model.IPSecVpnLocalEndpoint
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), localEndpointsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPSecVpnLocalEndpoint), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *localEndpointsClient) List(tier0IdParam string, serviceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnLocalEndpointListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(localEndpointsListInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPSecVpnLocalEndpointListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := localEndpointsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints", "list", inputDataValue, executionContext)
	var emptyOutput model.IPSecVpnLocalEndpointListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), localEndpointsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPSecVpnLocalEndpointListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *localEndpointsClient) Patch(tier0IdParam string, serviceIdParam string, localEndpointIdParam string, ipSecVpnLocalEndpointParam model.IPSecVpnLocalEndpoint) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(localEndpointsPatchInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("LocalEndpointId", localEndpointIdParam)
	sv.AddStructField("IpSecVpnLocalEndpoint", ipSecVpnLocalEndpointParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := localEndpointsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (lIface *localEndpointsClient) Update(tier0IdParam string, serviceIdParam string, localEndpointIdParam string, ipSecVpnLocalEndpointParam model.IPSecVpnLocalEndpoint) (model.IPSecVpnLocalEndpoint, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(localEndpointsUpdateInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("LocalEndpointId", localEndpointIdParam)
	sv.AddStructField("IpSecVpnLocalEndpoint", ipSecVpnLocalEndpointParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.IPSecVpnLocalEndpoint
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := localEndpointsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.ipsec_vpn_services.local_endpoints", "update", inputDataValue, executionContext)
	var emptyOutput model.IPSecVpnLocalEndpoint
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), localEndpointsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.IPSecVpnLocalEndpoint), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
