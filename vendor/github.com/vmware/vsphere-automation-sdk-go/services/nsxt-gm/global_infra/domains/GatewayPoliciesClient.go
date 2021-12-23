// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: GatewayPolicies
// Used by client-side stubs.

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

const _ = core.SupportedByRuntimeVersion1

type GatewayPoliciesClient interface {

	// Delete GatewayPolicy
	//
	// @param domainIdParam (required)
	// @param gatewayPolicyIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(domainIdParam string, gatewayPolicyIdParam string) error

	// Read gateway policy for a domain.
	//
	// @param domainIdParam (required)
	// @param gatewayPolicyIdParam (required)
	// @return com.vmware.nsx_global_policy.model.GatewayPolicy
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(domainIdParam string, gatewayPolicyIdParam string) (model.GatewayPolicy, error)

	// List all gateway policies for specified Domain.
	//
	// @param domainIdParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
	// @param includeRuleCountParam Include the count of rules in policy (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_global_policy.model.GatewayPolicyListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GatewayPolicyListResult, error)

	// Update the gateway policy for a domain. This is a full replace. All the rules are replaced. Performance Note: If you want to edit several rules in a gateway policy use this API. It will perform better than several individual rule APIs. Just pass all the rules which you wish to edit as embedded rules to it.
	//
	// @param domainIdParam (required)
	// @param gatewayPolicyIdParam (required)
	// @param gatewayPolicyParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) error

	// This is used to set a precedence of a gateway policy w.r.t others.
	//
	// @param domainIdParam (required)
	// @param gatewayPolicyIdParam (required)
	// @param gatewayPolicyParam (required)
	// @param anchorPathParam The security policy/rule path if operation is 'insert_after' or 'insert_before' (optional)
	// @param operationParam Operation (optional, default to insert_top)
	// @return com.vmware.nsx_global_policy.model.GatewayPolicy
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Revise(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy, anchorPathParam *string, operationParam *string) (model.GatewayPolicy, error)

	// Update the gateway policy for a domain. This is a full replace. All the rules are replaced. Performance Note: If you want to edit several rules in a gateway policy, use this API. It will perform better than several individual rule APIs. Just pass all the rules which you wish to edit as embedded rules to it.
	//
	// @param domainIdParam (required)
	// @param gatewayPolicyIdParam (required)
	// @param gatewayPolicyParam (required)
	// @return com.vmware.nsx_global_policy.model.GatewayPolicy
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) (model.GatewayPolicy, error)
}

type gatewayPoliciesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewGatewayPoliciesClient(connector client.Connector) *gatewayPoliciesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"revise": core.NewMethodIdentifier(interfaceIdentifier, "revise"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	gIface := gatewayPoliciesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &gIface
}

func (gIface *gatewayPoliciesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := gIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (gIface *gatewayPoliciesClient) Delete(domainIdParam string, gatewayPolicyIdParam string) error {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesDeleteInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("GatewayPolicyId", gatewayPolicyIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (gIface *gatewayPoliciesClient) Get(domainIdParam string, gatewayPolicyIdParam string) (model.GatewayPolicy, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesGetInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("GatewayPolicyId", gatewayPolicyIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.GatewayPolicy
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "get", inputDataValue, executionContext)
	var emptyOutput model.GatewayPolicy
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), gatewayPoliciesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.GatewayPolicy), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *gatewayPoliciesClient) List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GatewayPolicyListResult, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesListInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeMarkForDeleteObjects", includeMarkForDeleteObjectsParam)
	sv.AddStructField("IncludeRuleCount", includeRuleCountParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.GatewayPolicyListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "list", inputDataValue, executionContext)
	var emptyOutput model.GatewayPolicyListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), gatewayPoliciesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.GatewayPolicyListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *gatewayPoliciesClient) Patch(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) error {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesPatchInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("GatewayPolicyId", gatewayPolicyIdParam)
	sv.AddStructField("GatewayPolicy", gatewayPolicyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (gIface *gatewayPoliciesClient) Revise(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy, anchorPathParam *string, operationParam *string) (model.GatewayPolicy, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesReviseInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("GatewayPolicyId", gatewayPolicyIdParam)
	sv.AddStructField("GatewayPolicy", gatewayPolicyParam)
	sv.AddStructField("AnchorPath", anchorPathParam)
	sv.AddStructField("Operation", operationParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.GatewayPolicy
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesReviseRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "revise", inputDataValue, executionContext)
	var emptyOutput model.GatewayPolicy
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), gatewayPoliciesReviseOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.GatewayPolicy), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (gIface *gatewayPoliciesClient) Update(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) (model.GatewayPolicy, error) {
	typeConverter := gIface.connector.TypeConverter()
	executionContext := gIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(gatewayPoliciesUpdateInputType(), typeConverter)
	sv.AddStructField("DomainId", domainIdParam)
	sv.AddStructField("GatewayPolicyId", gatewayPolicyIdParam)
	sv.AddStructField("GatewayPolicy", gatewayPolicyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.GatewayPolicy
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := gatewayPoliciesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	gIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := gIface.connector.GetApiProvider().Invoke("com.vmware.nsx_global_policy.global_infra.domains.gateway_policies", "update", inputDataValue, executionContext)
	var emptyOutput model.GatewayPolicy
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), gatewayPoliciesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.GatewayPolicy), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), gIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
