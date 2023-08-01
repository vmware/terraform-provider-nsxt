// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: NsServices
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

type NsServicesClient interface {

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceParam (required)
	// @return com.vmware.nsx.model.NSService
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(nsServiceParam nsxModel.NSService) (nsxModel.NSService, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceIdParam NSService Id (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(nsServiceIdParam string, forceParam *bool) error

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceIdParam NSService Id (required)
	// @return com.vmware.nsx.model.NSService
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(nsServiceIdParam string) (nsxModel.NSService, error)

	// Returns paginated list of NSServices
	//  This API is deprecated. Please use policy API /infra/services/nsxt-mp to fetch all services.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param defaultServiceParam Fetch all default NSServices (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.NSServiceListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, defaultServiceParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSServiceListResult, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsServiceIdParam NSService Id (required)
	// @param nsServiceParam (required)
	// @return com.vmware.nsx.model.NSService
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(nsServiceIdParam string, nsServiceParam nsxModel.NSService) (nsxModel.NSService, error)
}

type nsServicesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewNsServicesClient(connector vapiProtocolClient_.Connector) *nsServicesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.ns_services")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	nIface := nsServicesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &nIface
}

func (nIface *nsServicesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := nIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (nIface *nsServicesClient) Create(nsServiceParam nsxModel.NSService) (nsxModel.NSService, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServicesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServicesCreateInputType(), typeConverter)
	sv.AddStructField("NsService", nsServiceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSService
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_services", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSService
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServicesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSService), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServicesClient) Delete(nsServiceIdParam string, forceParam *bool) error {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServicesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServicesDeleteInputType(), typeConverter)
	sv.AddStructField("NsServiceId", nsServiceIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_services", "delete", inputDataValue, executionContext)
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

func (nIface *nsServicesClient) Get(nsServiceIdParam string) (nsxModel.NSService, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServicesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServicesGetInputType(), typeConverter)
	sv.AddStructField("NsServiceId", nsServiceIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSService
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_services", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSService
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServicesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSService), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServicesClient) List(cursorParam *string, defaultServiceParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSServiceListResult, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServicesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServicesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("DefaultService", defaultServiceParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSServiceListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_services", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSServiceListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServicesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSServiceListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsServicesClient) Update(nsServiceIdParam string, nsServiceParam nsxModel.NSService) (nsxModel.NSService, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsServicesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsServicesUpdateInputType(), typeConverter)
	sv.AddStructField("NsServiceId", nsServiceIdParam)
	sv.AddStructField("NsService", nsServiceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSService
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_services", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSService
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsServicesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSService), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
