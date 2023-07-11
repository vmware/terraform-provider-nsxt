// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: VirtualMachines
// Used by client-side stubs.

package fabric

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type VirtualMachinesClient interface {

	// Perform action on a specific virtual machine. External id of the virtual machine needs to be provided in the request body. Some of the actions that can be performed are update tags, add tags, remove tags. To add tags to existing list of tag, use action parameter add_tags. To remove tags from existing list of tag, use action parameter remove_tags. To replace existing tags with new tags, use action parameter update_tags. To clear all tags, provide an empty list and action parameter as update_tags. The vmw-async: True HTTP header cannot be used with this API.
	//
	// @param virtualMachineTagUpdateParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Addtags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error

	// Returns information about all virtual machines. If you have not added NSX tags on the VM or removed all the NSX tags that were earlier added to the VM, then tags property is not returned in the API response.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param displayNameParam Display Name of the virtual machine (optional)
	// @param excludeVmTypeParam VM types to be excluded (optional)
	// @param externalIdParam External id of the virtual machine (optional)
	// @param hostIdParam Id of the host where this vif is located (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.VirtualMachineListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, displayNameParam *string, excludeVmTypeParam *string, externalIdParam *string, hostIdParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.VirtualMachineListResult, error)

	// Perform action on a specific virtual machine. External id of the virtual machine needs to be provided in the request body. Some of the actions that can be performed are update tags, add tags, remove tags. To add tags to existing list of tag, use action parameter add_tags. To remove tags from existing list of tag, use action parameter remove_tags. To replace existing tags with new tags, use action parameter update_tags. To clear all tags, provide an empty list and action parameter as update_tags. The vmw-async: True HTTP header cannot be used with this API.
	//
	// @param virtualMachineTagUpdateParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Removetags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error

	// Perform action on a specific virtual machine. External id of the virtual machine needs to be provided in the request body. Some of the actions that can be performed are update tags, add tags, remove tags. To add tags to existing list of tag, use action parameter add_tags. To remove tags from existing list of tag, use action parameter remove_tags. To replace existing tags with new tags, use action parameter update_tags. To clear all tags, provide an empty list and action parameter as update_tags. The vmw-async: True HTTP header cannot be used with this API.
	//
	// @param virtualMachineTagUpdateParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Updatetags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error
}

type virtualMachinesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewVirtualMachinesClient(connector vapiProtocolClient_.Connector) *virtualMachinesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.fabric.virtual_machines")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"addtags":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "addtags"),
		"list":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"removetags": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "removetags"),
		"updatetags": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "updatetags"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	vIface := virtualMachinesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &vIface
}

func (vIface *virtualMachinesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := vIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (vIface *virtualMachinesClient) Addtags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	operationRestMetaData := virtualMachinesAddtagsRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(virtualMachinesAddtagsInputType(), typeConverter)
	sv.AddStructField("VirtualMachineTagUpdate", virtualMachineTagUpdateParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.virtual_machines", "addtags", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (vIface *virtualMachinesClient) List(cursorParam *string, displayNameParam *string, excludeVmTypeParam *string, externalIdParam *string, hostIdParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.VirtualMachineListResult, error) {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	operationRestMetaData := virtualMachinesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(virtualMachinesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DisplayName", displayNameParam)
	sv.AddStructField("ExcludeVmType", excludeVmTypeParam)
	sv.AddStructField("ExternalId", externalIdParam)
	sv.AddStructField("HostId", hostIdParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.VirtualMachineListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.virtual_machines", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.VirtualMachineListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), VirtualMachinesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.VirtualMachineListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (vIface *virtualMachinesClient) Removetags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	operationRestMetaData := virtualMachinesRemovetagsRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(virtualMachinesRemovetagsInputType(), typeConverter)
	sv.AddStructField("VirtualMachineTagUpdate", virtualMachineTagUpdateParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.virtual_machines", "removetags", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (vIface *virtualMachinesClient) Updatetags(virtualMachineTagUpdateParam nsxModel.VirtualMachineTagUpdate) error {
	typeConverter := vIface.connector.TypeConverter()
	executionContext := vIface.connector.NewExecutionContext()
	operationRestMetaData := virtualMachinesUpdatetagsRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(virtualMachinesUpdatetagsInputType(), typeConverter)
	sv.AddStructField("VirtualMachineTagUpdate", virtualMachineTagUpdateParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := vIface.connector.GetApiProvider().Invoke("com.vmware.nsx.fabric.virtual_machines", "updatetags", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), vIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}
