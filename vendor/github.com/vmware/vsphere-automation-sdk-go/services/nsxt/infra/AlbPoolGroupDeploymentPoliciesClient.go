// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: AlbPoolGroupDeploymentPolicies
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

type AlbPoolGroupDeploymentPoliciesClient interface {

	// Delete the ALBPoolGroupDeploymentPolicy along with all the entities contained by this ALBPoolGroupDeploymentPolicy.
	//
	// @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(albPoolgroupdeploymentpolicyIdParam string, forceParam *bool) error

	// Read a ALBPoolGroupDeploymentPolicy.
	//
	// @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
	// @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicy
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(albPoolgroupdeploymentpolicyIdParam string) (model.ALBPoolGroupDeploymentPolicy, error)

	// Paginated list of all ALBPoolGroupDeploymentPolicy for infra.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicyApiResponse
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPoolGroupDeploymentPolicyApiResponse, error)

	// If a ALBpoolgroupdeploymentpolicy with the alb-poolgroupdeploymentpolicy-id is not already present, create a new ALBpoolgroupdeploymentpolicy. If it already exists, update the ALBpoolgroupdeploymentpolicy. This is a full replace.
	//
	// @param albPoolgroupdeploymentpolicyIdParam ALBpoolgroupdeploymentpolicy ID (required)
	// @param aLBPoolGroupDeploymentPolicyParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) error

	// If a ALBPoolGroupDeploymentPolicy with the alb-PoolGroupDeploymentPolicy-id is not already present, create a new ALBPoolGroupDeploymentPolicy. If it already exists, update the ALBPoolGroupDeploymentPolicy. This is a full replace.
	//
	// @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
	// @param aLBPoolGroupDeploymentPolicyParam (required)
	// @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicy
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) (model.ALBPoolGroupDeploymentPolicy, error)
}

type albPoolGroupDeploymentPoliciesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewAlbPoolGroupDeploymentPoliciesClient(connector client.Connector) *albPoolGroupDeploymentPoliciesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	aIface := albPoolGroupDeploymentPoliciesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &aIface
}

func (aIface *albPoolGroupDeploymentPoliciesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := aIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (aIface *albPoolGroupDeploymentPoliciesClient) Delete(albPoolgroupdeploymentpolicyIdParam string, forceParam *bool) error {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albPoolGroupDeploymentPoliciesDeleteInputType(), typeConverter)
	sv.AddStructField("AlbPoolgroupdeploymentpolicyId", albPoolgroupdeploymentpolicyIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albPoolGroupDeploymentPoliciesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies", "delete", inputDataValue, executionContext)
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

func (aIface *albPoolGroupDeploymentPoliciesClient) Get(albPoolgroupdeploymentpolicyIdParam string) (model.ALBPoolGroupDeploymentPolicy, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albPoolGroupDeploymentPoliciesGetInputType(), typeConverter)
	sv.AddStructField("AlbPoolgroupdeploymentpolicyId", albPoolgroupdeploymentpolicyIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBPoolGroupDeploymentPolicy
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albPoolGroupDeploymentPoliciesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies", "get", inputDataValue, executionContext)
	var emptyOutput model.ALBPoolGroupDeploymentPolicy
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albPoolGroupDeploymentPoliciesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBPoolGroupDeploymentPolicy), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *albPoolGroupDeploymentPoliciesClient) List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPoolGroupDeploymentPolicyApiResponse, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albPoolGroupDeploymentPoliciesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBPoolGroupDeploymentPolicyApiResponse
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albPoolGroupDeploymentPoliciesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies", "list", inputDataValue, executionContext)
	var emptyOutput model.ALBPoolGroupDeploymentPolicyApiResponse
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albPoolGroupDeploymentPoliciesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBPoolGroupDeploymentPolicyApiResponse), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (aIface *albPoolGroupDeploymentPoliciesClient) Patch(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) error {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albPoolGroupDeploymentPoliciesPatchInputType(), typeConverter)
	sv.AddStructField("AlbPoolgroupdeploymentpolicyId", albPoolgroupdeploymentpolicyIdParam)
	sv.AddStructField("ALBPoolGroupDeploymentPolicy", aLBPoolGroupDeploymentPolicyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albPoolGroupDeploymentPoliciesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies", "patch", inputDataValue, executionContext)
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

func (aIface *albPoolGroupDeploymentPoliciesClient) Update(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) (model.ALBPoolGroupDeploymentPolicy, error) {
	typeConverter := aIface.connector.TypeConverter()
	executionContext := aIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(albPoolGroupDeploymentPoliciesUpdateInputType(), typeConverter)
	sv.AddStructField("AlbPoolgroupdeploymentpolicyId", albPoolgroupdeploymentpolicyIdParam)
	sv.AddStructField("ALBPoolGroupDeploymentPolicy", aLBPoolGroupDeploymentPolicyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.ALBPoolGroupDeploymentPolicy
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := albPoolGroupDeploymentPoliciesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	aIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := aIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.alb_pool_group_deployment_policies", "update", inputDataValue, executionContext)
	var emptyOutput model.ALBPoolGroupDeploymentPolicy
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), albPoolGroupDeploymentPoliciesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.ALBPoolGroupDeploymentPolicy), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), aIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
