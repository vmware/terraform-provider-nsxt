// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: MacSets
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

type MacSetsClient interface {

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param mACSetParam (required)
	// @return com.vmware.nsx.model.MACSet
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(mACSetParam nsxModel.MACSet) (nsxModel.MACSet, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param macSetIdParam MACSet Id (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(macSetIdParam string, forceParam *bool) error

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param macSetIdParam MACSet Id (required)
	// @return com.vmware.nsx.model.MACSet
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(macSetIdParam string) (nsxModel.MACSet, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.MACSetListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.MACSetListResult, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param macSetIdParam MACSet Id (required)
	// @param mACSetParam (required)
	// @return com.vmware.nsx.model.MACSet
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(macSetIdParam string, mACSetParam nsxModel.MACSet) (nsxModel.MACSet, error)
}

type macSetsClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewMacSetsClient(connector vapiProtocolClient_.Connector) *macSetsClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.mac_sets")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	mIface := macSetsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &mIface
}

func (mIface *macSetsClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := mIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (mIface *macSetsClient) Create(mACSetParam nsxModel.MACSet) (nsxModel.MACSet, error) {
	typeConverter := mIface.connector.TypeConverter()
	executionContext := mIface.connector.NewExecutionContext()
	operationRestMetaData := macSetsCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(macSetsCreateInputType(), typeConverter)
	sv.AddStructField("MACSet", mACSetParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MACSet
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := mIface.connector.GetApiProvider().Invoke("com.vmware.nsx.mac_sets", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.MACSet
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), MacSetsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MACSet), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), mIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (mIface *macSetsClient) Delete(macSetIdParam string, forceParam *bool) error {
	typeConverter := mIface.connector.TypeConverter()
	executionContext := mIface.connector.NewExecutionContext()
	operationRestMetaData := macSetsDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(macSetsDeleteInputType(), typeConverter)
	sv.AddStructField("MacSetId", macSetIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := mIface.connector.GetApiProvider().Invoke("com.vmware.nsx.mac_sets", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), mIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (mIface *macSetsClient) Get(macSetIdParam string) (nsxModel.MACSet, error) {
	typeConverter := mIface.connector.TypeConverter()
	executionContext := mIface.connector.NewExecutionContext()
	operationRestMetaData := macSetsGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(macSetsGetInputType(), typeConverter)
	sv.AddStructField("MacSetId", macSetIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MACSet
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := mIface.connector.GetApiProvider().Invoke("com.vmware.nsx.mac_sets", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.MACSet
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), MacSetsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MACSet), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), mIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (mIface *macSetsClient) List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.MACSetListResult, error) {
	typeConverter := mIface.connector.TypeConverter()
	executionContext := mIface.connector.NewExecutionContext()
	operationRestMetaData := macSetsListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(macSetsListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MACSetListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := mIface.connector.GetApiProvider().Invoke("com.vmware.nsx.mac_sets", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.MACSetListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), MacSetsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MACSetListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), mIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (mIface *macSetsClient) Update(macSetIdParam string, mACSetParam nsxModel.MACSet) (nsxModel.MACSet, error) {
	typeConverter := mIface.connector.TypeConverter()
	executionContext := mIface.connector.NewExecutionContext()
	operationRestMetaData := macSetsUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(macSetsUpdateInputType(), typeConverter)
	sv.AddStructField("MacSetId", macSetIdParam)
	sv.AddStructField("MACSet", mACSetParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.MACSet
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := mIface.connector.GetApiProvider().Invoke("com.vmware.nsx.mac_sets", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.MACSet
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), MacSetsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.MACSet), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), mIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
