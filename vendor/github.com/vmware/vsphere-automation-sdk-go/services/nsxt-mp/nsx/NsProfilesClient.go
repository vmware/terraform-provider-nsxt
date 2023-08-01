// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: NsProfiles
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

type NsProfilesClient interface {

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsProfileParam (required)
	// @return com.vmware.nsx.model.NSProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(nsProfileParam nsxModel.NSProfile) (nsxModel.NSProfile, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsProfileIdParam NSProfile Id (required)
	// @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(nsProfileIdParam string, forceParam *bool) error

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsProfileIdParam NSProfile Id (required)
	// @return com.vmware.nsx.model.NSProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(nsProfileIdParam string) (nsxModel.NSProfile, error)

	// List the NSProfiles created in a paginated format.The page size is restricted to 50 NSProfiles, so that the size of the response remains small even when there are high number of NSProfiles with multiple attributes and multiple attribute values for each attribute.
	//  This API is deprecated. Use GET /infra/context-profiles
	//
	// Deprecated: This API element is deprecated.
	//
	// @param attributeTypeParam Fetch NSProfiles for the given attribute type (optional)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.NSProfileListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(attributeTypeParam *string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSProfileListResult, error)

	//
	//
	// Deprecated: This API element is deprecated.
	//
	// @param nsProfileIdParam NSProfile Id (required)
	// @param nsProfileParam (required)
	// @return com.vmware.nsx.model.NSProfile
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(nsProfileIdParam string, nsProfileParam nsxModel.NSProfile) (nsxModel.NSProfile, error)
}

type nsProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewNsProfilesClient(connector vapiProtocolClient_.Connector) *nsProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.ns_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	nIface := nsProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &nIface
}

func (nIface *nsProfilesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := nIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (nIface *nsProfilesClient) Create(nsProfileParam nsxModel.NSProfile) (nsxModel.NSProfile, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsProfilesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsProfilesCreateInputType(), typeConverter)
	sv.AddStructField("NsProfile", nsProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_profiles", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsProfilesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsProfilesClient) Delete(nsProfileIdParam string, forceParam *bool) error {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsProfilesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("NsProfileId", nsProfileIdParam)
	sv.AddStructField("Force", forceParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_profiles", "delete", inputDataValue, executionContext)
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

func (nIface *nsProfilesClient) Get(nsProfileIdParam string) (nsxModel.NSProfile, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsProfilesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsProfilesGetInputType(), typeConverter)
	sv.AddStructField("NsProfileId", nsProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_profiles", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsProfilesClient) List(attributeTypeParam *string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.NSProfileListResult, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsProfilesListInputType(), typeConverter)
	sv.AddStructField("AttributeType", attributeTypeParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSProfileListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (nIface *nsProfilesClient) Update(nsProfileIdParam string, nsProfileParam nsxModel.NSProfile) (nsxModel.NSProfile, error) {
	typeConverter := nIface.connector.TypeConverter()
	executionContext := nIface.connector.NewExecutionContext()
	operationRestMetaData := nsProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(nsProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("NsProfileId", nsProfileIdParam)
	sv.AddStructField("NsProfile", nsProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NSProfile
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := nIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ns_profiles", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.NSProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), NsProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NSProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), nIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
