// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Sessions
// Used by client-side stubs.

package l2vpn_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type SessionsClient interface {

	// Create or patch an L2VPN session under Tier-1 from Peer Codes. In addition to the L2VPN Session, the IPSec VPN Session, along with the IKE, Tunnel, and DPD Profiles are created and owned by the system. IPSec VPN Service and Local Endpoint are created only when required, i.e., an IPSec VPN Service does not already exist, or an IPSec VPN Local Endpoint with same local address does not already exist. Updating the L2VPN Session can be performed only through this API by specifying new peer codes. Use of specific APIs to update the L2VPN Session and the different resources associated with it is not allowed, except for IPSec VPN Service and Local Endpoint, resources that are not system owned. API supported only when L2VPN Service is in Client Mode. This API is deprecated. Please use POST /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/sessions/<session-id>?action=create_with_peer_code instead. Note: Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource. Also VPN path returned in the Alarm, GPRR payload may include the new VPN path.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param sessionIdParam (required)
	// @param l2VPNSessionDataParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Createwithpeercode(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionDataParam model.L2VPNSessionData) error

	// Delete L2VPN session under Tier-1. When L2VPN Service is in CLIENT Mode, the L2VPN Session is deleted along with its transpot tunnels and related resources. This API is deprecated. Please use DELETE /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/ sessions/<session-id> instead. Note: Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource. Also VPN path returned in the Alarm, GPRR payload may include the new VPN path.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param sessionIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string) error

	// Get L2VPN session under Tier-1. This API is deprecated. Please use GET /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/ sessions/<session-id> instead. Note: The API will return a new VPN path for \"transport_tunnels\" in the response payload instead of the deprecated API path Both paths refer to the same object. Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param sessionIdParam (required)
	// @return com.vmware.nsx_policy.model.L2VPNSession
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string) (model.L2VPNSession, error)

	// Get paginated list of all L2VPN sessions under Tier-1. This API is deprecated. Please use GET /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/sessions instead. Note: The API will return a new VPN path for \"transport_tunnels\" in the response payload instead of the deprecated API path Both paths refer to the same object. Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.L2VPNSessionListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L2VPNSessionListResult, error)

	// Create or patch an L2VPN session under Tier-1. API supported only when L2VPN Service is in Server Mode. This API is deprecated. Please use PATCH /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/ sessions/<session-id> instead. Note: Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource. Also VPN path returned in the Alarm, GPRR payload may include the new VPN path.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param sessionIdParam (required)
	// @param l2VPNSessionParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionParam model.L2VPNSession) error

	// Create or fully replace L2VPN session under Tier-1. API supported only when L2VPN Service is in Server Mode. Revision is optional for creation and required for update. This API is deprecated. Please use PUT /infra/tier-1s/<tier-1-id>/l2vpn-services/<service-id>/ sessions/<session-id> instead. Note: The API will return a new VPN path for \"transport_tunnels\" in the response payload instead of the deprecated API path Both paths refer to the same object. Please note that request is validated and any error messages returned from validation may include the new VPN path instead of the deprecated path. Both new path and old path refer to same resource. Also VPN path returned in the Alarm, GPRR payload may include the new VPN path.
	//
	// @param tier1IdParam (required)
	// @param localeServiceIdParam (required)
	// @param serviceIdParam (required)
	// @param sessionIdParam (required)
	// @param l2VPNSessionParam (required)
	// @return com.vmware.nsx_policy.model.L2VPNSession
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionParam model.L2VPNSession) (model.L2VPNSession, error)
}

type sessionsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewSessionsClient(connector client.Connector) *sessionsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"createwithpeercode": core.NewMethodIdentifier(interfaceIdentifier, "createwithpeercode"),
		"delete":             core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":               core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":              core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update":             core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	sIface := sessionsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &sIface
}

func (sIface *sessionsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := sIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (sIface *sessionsClient) Createwithpeercode(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionDataParam model.L2VPNSessionData) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsCreatewithpeercodeInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("SessionId", sessionIdParam)
	sv.AddStructField("L2VPNSessionData", l2VPNSessionDataParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsCreatewithpeercodeRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "createwithpeercode", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (sIface *sessionsClient) Delete(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsDeleteInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("SessionId", sessionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (sIface *sessionsClient) Get(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string) (model.L2VPNSession, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsGetInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("SessionId", sessionIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L2VPNSession
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "get", inputDataValue, executionContext)
	var emptyOutput model.L2VPNSession
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), sessionsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L2VPNSession), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *sessionsClient) List(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L2VPNSessionListResult, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsListInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L2VPNSessionListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "list", inputDataValue, executionContext)
	var emptyOutput model.L2VPNSessionListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), sessionsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L2VPNSessionListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (sIface *sessionsClient) Patch(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionParam model.L2VPNSession) error {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsPatchInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("SessionId", sessionIdParam)
	sv.AddStructField("L2VPNSession", l2VPNSessionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (sIface *sessionsClient) Update(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, sessionIdParam string, l2VPNSessionParam model.L2VPNSession) (model.L2VPNSession, error) {
	typeConverter := sIface.connector.TypeConverter()
	executionContext := sIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(sessionsUpdateInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("LocaleServiceId", localeServiceIdParam)
	sv.AddStructField("ServiceId", serviceIdParam)
	sv.AddStructField("SessionId", sessionIdParam)
	sv.AddStructField("L2VPNSession", l2VPNSessionParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.L2VPNSession
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := sessionsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	sIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := sIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.locale_services.l2vpn_services.sessions", "update", inputDataValue, executionContext)
	var emptyOutput model.L2VPNSession
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), sessionsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.L2VPNSession), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), sIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
