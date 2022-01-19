// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: VirtualMachines
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

type VirtualMachinesClient interface {

	// This API filters objects of type virtual machines from the specified NSX Manager. This API has been deprecated. Please use the new API GET /infra/realized-state/virtual-machines
	//
	// @param enforcementPointNameParam (required)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param dslParam Search DSL (domain specific language) query (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param queryParam Search query (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.SearchResponse
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(enforcementPointNameParam string, cursorParam *string, dslParam *string, includedFieldsParam *string, pageSizeParam *int64, queryParam *string, sortAscendingParam *bool, sortByParam *string) (model.SearchResponse, error)

	// Allows an admin to apply multiple tags to a virtual machine. This operation does not store the intent on the policy side. It applies the tag directly on the specified enforcement point. This operation will replace the existing tags on the virtual machine with the ones that have been passed. If the application of tag fails on the enforcement point, then an error is reported. The admin will have to retry the operation again. Policy framework does not perform a retry. Failure could occur due to multiple reasons. For e.g enforcement point is down, Enforcement point could not apply the tag due to constraints like max tags limit exceeded, etc.
	//
	// @param enforcementPointNameParam (required)
	// @param virtualMachineTagsUpdateParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Updatetags(enforcementPointNameParam string, virtualMachineTagsUpdateParam model.VirtualMachineTagsUpdate) error
}

type virtualMachinesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewVirtualMachinesClient(connector client.Connector) *virtualMachinesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.realized_state.enforcement_points.virtual_machines")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"list":       core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"updatetags": core.NewMethodIdentifier(interfaceIdentifier, "updatetags"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	vIface := virtualMachinesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &vIface
}

func (vIface *virtualMachinesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := vIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (vIface *virtualMachinesClient) List(enforcementPointNameParam string, cursorParam *string, dslParam *string, includedFieldsParam *string, pageSizeParam *int64, queryParam *string, sortAscendingParam *bool, sortByParam *string) (model.SearchResponse, error) {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(virtualMachinesListInputType(), typeConverter)
	sv.AddStructField("EnforcementPointName", enforcementPointNameParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("Dsl", dslParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("Query", queryParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.SearchResponse
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := virtualMachinesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	vIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.realized_state.enforcement_points.virtual_machines", "list", inputDataValue, executionContext)
	var emptyOutput model.SearchResponse
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), virtualMachinesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.SearchResponse), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (vIface *virtualMachinesClient) Updatetags(enforcementPointNameParam string, virtualMachineTagsUpdateParam model.VirtualMachineTagsUpdate) error {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(virtualMachinesUpdatetagsInputType(), typeConverter)
	sv.AddStructField("EnforcementPointName", enforcementPointNameParam)
	sv.AddStructField("VirtualMachineTagsUpdate", virtualMachineTagsUpdateParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := virtualMachinesUpdatetagsRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	vIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.realized_state.enforcement_points.virtual_machines", "updatetags", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}
