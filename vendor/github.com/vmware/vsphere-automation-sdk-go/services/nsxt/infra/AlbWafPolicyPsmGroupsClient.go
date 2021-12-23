// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: AlbWafPolicyPsmGroups
// Used by client-side stubs.

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type AlbWafPolicyPsmGroupsClient interface {

	// Delete the ALBWafPolicyPSMGroup along with all the entities contained by this ALBWafPolicyPSMGroup.
	//
	// @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(albWafpolicypsmgroupIdParam string, forceParam *bool) error

	// Read a ALBWafPolicyPSMGroup.
	//
	// @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
	// @return com.vmware.nsx_policy.model.ALBWafPolicyPSMGroup
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(albWafpolicypsmgroupIdParam string) (model.ALBWafPolicyPSMGroup, error)

	// Paginated list of all ALBWafPolicyPSMGroup for infra.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.ALBWafPolicyPSMGroupApiResponse
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafPolicyPSMGroupApiResponse, error)

	// If a ALBwafpolicypsmgroup with the alb-wafpolicypsmgroup-id is not already present, create a new ALBwafpolicypsmgroup. If it already exists, update the ALBwafpolicypsmgroup. This is a full replace.
	//
	// @param albWafpolicypsmgroupIdParam ALBwafpolicypsmgroup ID (required)
	// @param aLBWafPolicyPSMGroupParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) error

	// If a ALBWafPolicyPSMGroup with the alb-WafPolicyPSMGroup-id is not already present, create a new ALBWafPolicyPSMGroup. If it already exists, update the ALBWafPolicyPSMGroup. This is a full replace.
	//
	// @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
	// @param aLBWafPolicyPSMGroupParam (required)
	// @return com.vmware.nsx_policy.model.ALBWafPolicyPSMGroup
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) (model.ALBWafPolicyPSMGroup, error)
}

type albWafPolicyPsmGroupsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewAlbWafPolicyPsmGroupsClient(connector client.Connector) *albWafPolicyPsmGroupsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	aIface := albWafPolicyPsmGroupsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &aIface
}

func (aIface *albWafPolicyPsmGroupsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := aIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (aIface *albWafPolicyPsmGroupsClient) Delete(albWafpolicypsmgroupIdParam string, forceParam *bool) error {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albWafPolicyPsmGroupsDeleteInputType(), typeConverter)
	sv.AddStructField("AlbWafpolicypsmgroupId", albWafpolicypsmgroupIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albWafPolicyPsmGroupsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (aIface *albWafPolicyPsmGroupsClient) Get(albWafpolicypsmgroupIdParam string) (model.ALBWafPolicyPSMGroup, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albWafPolicyPsmGroupsGetInputType(), typeConverter)
	sv.AddStructField("AlbWafpolicypsmgroupId", albWafpolicypsmgroupIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBWafPolicyPSMGroup
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albWafPolicyPsmGroupsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups", "get", inputDataValue, executionContext)
	var emptyOutput model.ALBWafPolicyPSMGroup
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albWafPolicyPsmGroupsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBWafPolicyPSMGroup), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *albWafPolicyPsmGroupsClient) List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafPolicyPSMGroupApiResponse, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albWafPolicyPsmGroupsListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBWafPolicyPSMGroupApiResponse
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albWafPolicyPsmGroupsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups", "list", inputDataValue, executionContext)
	var emptyOutput model.ALBWafPolicyPSMGroupApiResponse
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albWafPolicyPsmGroupsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBWafPolicyPSMGroupApiResponse), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *albWafPolicyPsmGroupsClient) Patch(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) error {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albWafPolicyPsmGroupsPatchInputType(), typeConverter)
	sv.AddStructField("AlbWafpolicypsmgroupId", albWafpolicypsmgroupIdParam)
	sv.AddStructField("ALBWafPolicyPSMGroup", aLBWafPolicyPSMGroupParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albWafPolicyPsmGroupsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (aIface *albWafPolicyPsmGroupsClient) Update(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) (model.ALBWafPolicyPSMGroup, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albWafPolicyPsmGroupsUpdateInputType(), typeConverter)
	sv.AddStructField("AlbWafpolicypsmgroupId", albWafpolicypsmgroupIdParam)
	sv.AddStructField("ALBWafPolicyPSMGroup", aLBWafPolicyPSMGroupParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBWafPolicyPSMGroup
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albWafPolicyPsmGroupsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_waf_policy_psm_groups", "update", inputDataValue, executionContext)
	var emptyOutput model.ALBWafPolicyPSMGroup
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albWafPolicyPsmGroupsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBWafPolicyPSMGroup), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
