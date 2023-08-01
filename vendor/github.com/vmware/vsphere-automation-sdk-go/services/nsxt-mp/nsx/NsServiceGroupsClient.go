// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: NsServiceGroups
// Used by client-side stubs.

package nsx

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type NsServiceGroupsClient interface {

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceGroupParam (required)
	// @return com.vmware.nsx.model.NSServiceGroup
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(nsServiceGroupParam nsxModel.NSServiceGroup) (nsxModel.NSServiceGroup, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceGroupIdParam NSServiceGroup Id (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(nsServiceGroupIdParam string, forceParam *bool) error

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceGroupIdParam NSServiceGroup Id (required)
	// @return com.vmware.nsx.model.NSServiceGroup
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(nsServiceGroupIdParam string) (nsxModel.NSServiceGroup, error)

	// Returns paginated list of NSServiceGroups
	//  This API is deprecated. Please use policy API /infra/services/nsxt-mp to fetch all services.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param defaultServiceParam Fetch all default NSServiceGroups (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.NSServiceGroupListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, defaultServiceParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSServiceGroupListResult, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceGroupIdParam NSServiceGroup Id (required)
	// @param nsServiceGroupParam (required)
	// @return com.vmware.nsx.model.NSServiceGroup
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(nsServiceGroupIdParam string, nsServiceGroupParam nsxModel.NSServiceGroup) (nsxModel.NSServiceGroup, error)
}

type nsServiceGroupsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewNsServiceGroupsClient(connector vapiProtocolClient_.Connector) *nsServiceGroupsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.ns_service_groups")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	nIface := nsServiceGroupsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &nIface
}

func (nIface *nsServiceGroupsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := nIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (nIface *nsServiceGroupsClient) Create(nsServiceGroupParam nsxModel.NSServiceGroup) (nsxModel.NSServiceGroup, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServiceGroupsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServiceGroupsCreateInputType(), typeConverter)
	sv.AddStructField("NsServiceGroup", nsServiceGroupParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSServiceGroup
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_service_groups", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSServiceGroup
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServiceGroupsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSServiceGroup), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServiceGroupsClient) Delete(nsServiceGroupIdParam string, forceParam *bool) error {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServiceGroupsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServiceGroupsDeleteInputType(), typeConverter)
	sv.AddStructField("NsServiceGroupId", nsServiceGroupIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_service_groups", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (nIface *nsServiceGroupsClient) Get(nsServiceGroupIdParam string) (nsxModel.NSServiceGroup, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServiceGroupsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServiceGroupsGetInputType(), typeConverter)
	sv.AddStructField("NsServiceGroupId", nsServiceGroupIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSServiceGroup
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_service_groups", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSServiceGroup
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServiceGroupsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSServiceGroup), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServiceGroupsClient) List(cursorParam *string, defaultServiceParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSServiceGroupListResult, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServiceGroupsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServiceGroupsListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DefaultService", defaultServiceParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSServiceGroupListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_service_groups", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSServiceGroupListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServiceGroupsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSServiceGroupListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServiceGroupsClient) Update(nsServiceGroupIdParam string, nsServiceGroupParam nsxModel.NSServiceGroup) (nsxModel.NSServiceGroup, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServiceGroupsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServiceGroupsUpdateInputType(), typeConverter)
	sv.AddStructField("NsServiceGroupId", nsServiceGroupIdParam)
	sv.AddStructField("NsServiceGroup", nsServiceGroupParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSServiceGroup
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_service_groups", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSServiceGroup
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServiceGroupsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSServiceGroup), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
