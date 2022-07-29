// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: L3vpns
// Used by client-side stubs.

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type L3vpnsClient interface {

	// Delete the L3Vpn with the given id. This API is deprecated. Please use the following APIs instead: - DELETE /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/L3VPN_<l3vpn-id> to delete the associated IPSecVpnSession. - DELETE /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints/<local-endpoint-id> to delete the associated IPSecVpnLocalEndpoint. - DELETE /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id> to delete the associated IPSecVpnTunnelProfile. - DELETE /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id> to delete the associated IPSecVpnIkeProfile. - DELETE /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id> to delete the associated IPSecVpnDpdProfile. If used, this deprecated API will result in the following objects being internally deleted: - IPSecVpnSession: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/sessions/L3VPN_<l3vpn-id>. - IPSecVpnLocalEndpoint: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/local-endpoints/<local-endpoint-id> when not used by other IPSecVpnSessions. - IPSecVpnTunnelProfile: /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id>. - IPSecVpnIkeProfile: /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id>. - IPSecVpnDpdProfile: /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id>.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param l3vpnIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) error

	// Read the L3Vpn with the given id. No sensitive data is returned as part of the response. This API is deprecated. Please use the following APIs instead: - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/L3VPN_<l3vpn-id> to get the associated IPSecVpnSession. - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints/<local-endpoint-id> to get the associated IPSecVpnLocalEndpoint. - GET /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnTunnelProfile. - GET /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnIkeProfile. - GET /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnDpdProfile. If used, this deprecated API will not return L3Vpn with <l3vpn-id> id unless the associated IPSecVpnSession with L3VPN_<l3vpn-id> id exists. For example, if the IPSecVpnSession gets deleted using DELETE /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/L3VPN_<l3vpn-id>, the deprecated API will throw an ObjectNotFoundException.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param l3vpnIdParam (required)
	// @return com.vmware.nsx_policy.model.L3Vpn
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error)

	// Paginated list of L3Vpns. This API is deprecated. Please use the following APIs instead: - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions to list all IPSecVpnSessions. - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints to list all IPSecVpnLocalEndpoints. - GET /infra/ipsec-vpn-tunnel-profiles to list all IPSecVpnTunnelProfiles. - GET /infra/ipsec-vpn-ike-profiles to list all IPSecVpnIkeProfiles. - GET /infra/ipsec-vpn-dpd-profiles to list all IPSecVpnDpdProfiles. If used, this deprecated API will only return L3Vpns that were created through the deprecated PATCH and PUT /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/l3vpns/<l3vpn-id> APIs.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param l3vpnSessionParam Resource type of L3Vpn Session (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.L3VpnListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, l3vpnSessionParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L3VpnListResult, error)

	// Create the new L3Vpn if it does not exist. If the L3Vpn already exists, merge with the the existing one. This is a patch. - If the passed L3Vpn is a policy-based one and has new L3VpnRules, add them to the existing L3VpnRules. - If the passed L3Vpn is a policy-based one and also has existing L3VpnRules, update the existing L3VpnRules. This API is deprecated. Please use the following APIs instead: - PATCH /infra/ipsec-vpn-tunnel-profiles/<tunnel-profile-id> to patch the IPSecVpnTunnelProfile. - PATCH /infra/ipsec-vpn-ike-profiles/<ike-profile-id> to patch the IPSecVpnIkeProfile. - PATCH /infra/ipsec-vpn-dpd-profiles/<dpd-profile-id> to patch the IPSecVpnDpdProfile. - PATCH /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints/<local-endpoint-id> to patch the IPSecVpnLocalEndpoint. - PATCH /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/<l3vpn-id> to patch the IPSecVpnSession. If used, this deprecated API will result in the following objects being internally created/patched: - IPSecVpnTunnelProfile: /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id>. - IPSecVpnIkeProfile: /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id>. - IPSecVpnDpdProfile: /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id>. - IPSecVpnLocalEndpoint: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/local-endpoints/<local-endpoint-id>. If an object with the same \"local_address\" already exists, then it will be re-used. - IPSecVpnSession: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/sessions/L3VPN_<l3vpn-id>.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param l3vpnIdParam (required)
	// @param l3VpnParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) error

	// Read the L3Vpn with the given id. Sensitive data is returned as part of the response. This API is deprecated. Please use the following APIs instead: - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/L3VPN_<l3vpn-id>?action=show_sensitive_data to get the associated IPSecVpnSession. - GET /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints/<local-endpoint-id> to get the associated IPSecVpnLocalEndpoint. - GET /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnTunnelProfile. - GET /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnIkeProfile. - GET /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id> to get the associated IPSecVpnDpdProfile. If used, this deprecated API will not return L3Vpn with <l3vpn-id> id unless the associated IPSecVpnSession with L3VPN_<l3vpn-id> id exists. For example, if the IPSecVpnSession gets deleted using DELETE /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/L3VPN_<l3vpn-id>, the deprecated API will throw an ObjectNotFoundException.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param l3vpnIdParam (required)
	// @return com.vmware.nsx_policy.model.L3Vpn
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Showsensitivedata(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error)

	// Create a new L3Vpn if the L3Vpn with given id does not already exist. If the L3Vpn with the given id already exists, replace the existing L3Vpn. This a full replace. This API is deprecated. Please use the following APIs instead: - PUT /infra/ipsec-vpn-tunnel-profiles/<tunnel-profile-id> to update the IPSecVpnTunnelProfile. - PUT /infra/ipsec-vpn-ike-profiles/<ike-profile-id> to update the IPSecVpnIkeProfile. - PUT /infra/ipsec-vpn-dpd-profiles/<dpd-profile-id> to update the IPSecVpnDpdProfile. - PUT /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/local-endpoints/<local-endpoint-id> to update the IPSecVpnLocalEndpoint. - PUT /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ipsec-vpn-services/ default/sessions/<l3vpn-id> to update the IPSecVpnSession. If used, this deprecated API will result in the following objects being internally created/updated: - IPSecVpnTunnelProfile: /infra/ipsec-vpn-tunnel-profiles/L3VPN_<l3vpn-id>. - IPSecVpnIkeProfile: /infra/ipsec-vpn-ike-profiles/L3VPN_<l3vpn-id>. - IPSecVpnDpdProfile: /infra/ipsec-vpn-dpd-profiles/L3VPN_<l3vpn-id>. - IPSecVpnLocalEndpoint: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/local-endpoints/<local-endpoint-id>. If an object with the same \"local_address\" already exists, then it will be re-used. - IPSecVpnSession: /infra/tier-0s/<tier-0-id>/locale-services/<locale-service-id>/ ipsec-vpn-services/default/sessions/L3VPN_<l3vpn-id>.
	//
	// @param tier0IdParam (required)
	// @param localeServiceIdParam (required)
	// @param l3vpnIdParam (required)
	// @param l3VpnParam (required)
	// @return com.vmware.nsx_policy.model.L3Vpn
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) (model.L3Vpn, error)
}

type l3vpnsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewL3vpnsClient(connector client.Connector) *l3vpnsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete":            core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":               core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":              core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":             core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"showsensitivedata": core.NewMethodIdentifier(interfaceIdentifier, "showsensitivedata"),
		"update":            core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	lIface := l3vpnsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *l3vpnsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *l3vpnsClient) Delete(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsDeleteInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("L3vpnId", l3vpnIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "delete", inputDataValue, executionContext)
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

func (lIface *l3vpnsClient) Get(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsGetInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("L3vpnId", l3vpnIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L3Vpn
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "get", inputDataValue, executionContext)
	var emptyOutput model.L3Vpn
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), l3vpnsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L3Vpn), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *l3vpnsClient) List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, l3vpnSessionParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L3VpnListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsListInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("L3vpnSession", l3vpnSessionParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L3VpnListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "list", inputDataValue, executionContext)
	var emptyOutput model.L3VpnListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), l3vpnsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L3VpnListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *l3vpnsClient) Patch(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsPatchInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("L3vpnId", l3vpnIdParam)
	sv.AddStructField("L3Vpn", l3VpnParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "patch", inputDataValue, executionContext)
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

func (lIface *l3vpnsClient) Showsensitivedata(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsShowsensitivedataInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("L3vpnId", l3vpnIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L3Vpn
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsShowsensitivedataRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "showsensitivedata", inputDataValue, executionContext)
	var emptyOutput model.L3Vpn
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), l3vpnsShowsensitivedataOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L3Vpn), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *l3vpnsClient) Update(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) (model.L3Vpn, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(l3vpnsUpdateInputType(), typeConverter)
	sv.AddStructField("Tier0Id", tier0IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("L3vpnId", l3vpnIdParam)
	sv.AddStructField("L3Vpn", l3VpnParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L3Vpn
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := l3vpnsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	lIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_0s.locale_services.l3vpns", "update", inputDataValue, executionContext)
	var emptyOutput model.L3Vpn
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), l3vpnsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L3Vpn), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
