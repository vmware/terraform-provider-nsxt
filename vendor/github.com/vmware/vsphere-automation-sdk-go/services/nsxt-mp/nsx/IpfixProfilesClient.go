// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: IpfixProfiles
// Used by client-side stubs.

package nsx

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type IpfixProfilesClient interface {

	// Create a new IPFIX profile with essential properties.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param ipfixUpmProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.IpfixUpmProfile.
	// @return com.vmware.nsx.model.IpfixUpmProfile
	// The return value will contain all the properties defined in nsxModel.IpfixUpmProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(ipfixUpmProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Delete an existing IPFIX profile by ID.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param ipfixProfileIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(ipfixProfileIdParam string) error

	// Get an existing IPFIX profile by profile ID.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param ipfixProfileIdParam (required)
	// @return com.vmware.nsx.model.IpfixUpmProfile
	// The return value will contain all the properties defined in nsxModel.IpfixUpmProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(ipfixProfileIdParam string) (*vapiData_.StructValue, error)

	// Query IPFIX profiles with list parameters. List result can be filtered by profile type defined by IpfixUpmProfileType.
	//
	// @param appliedToEntityIdParam ID of Entity Applied with Profile (optional)
	// @param appliedToEntityTypeParam Supported Entity Types (optional)
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param profileTypesParam IPFIX Profile Type List (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.IpfixUpmProfileListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(appliedToEntityIdParam *string, appliedToEntityTypeParam *string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, profileTypesParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.IpfixUpmProfileListResult, error)

	// Update an existing IPFIX profile with profile ID and modified properties.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param ipfixProfileIdParam (required)
	// @param ipfixUpmProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.IpfixUpmProfile.
	// @return com.vmware.nsx.model.IpfixUpmProfile
	// The return value will contain all the properties defined in nsxModel.IpfixUpmProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(ipfixProfileIdParam string, ipfixUpmProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type ipfixProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewIpfixProfilesClient(connector vapiProtocolClient_.Connector) *ipfixProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.ipfix_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	iIface := ipfixProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &iIface
}

func (iIface *ipfixProfilesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := iIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (iIface *ipfixProfilesClient) Create(ipfixUpmProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	operationRestMetaData := ipfixProfilesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(ipfixProfilesCreateInputType(), typeConverter)
	sv.AddStructField("IpfixUpmProfile", ipfixUpmProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ipfix_profiles", "create", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), IpfixProfilesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *ipfixProfilesClient) Delete(ipfixProfileIdParam string) error {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	operationRestMetaData := ipfixProfilesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(ipfixProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("IpfixProfileId", ipfixProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ipfix_profiles", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (iIface *ipfixProfilesClient) Get(ipfixProfileIdParam string) (*vapiData_.StructValue, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	operationRestMetaData := ipfixProfilesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(ipfixProfilesGetInputType(), typeConverter)
	sv.AddStructField("IpfixProfileId", ipfixProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ipfix_profiles", "get", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), IpfixProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *ipfixProfilesClient) List(appliedToEntityIdParam *string, appliedToEntityTypeParam *string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, profileTypesParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.IpfixUpmProfileListResult, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	operationRestMetaData := ipfixProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(ipfixProfilesListInputType(), typeConverter)
	sv.AddStructField("AppliedToEntityId", appliedToEntityIdParam)
	sv.AddStructField("AppliedToEntityType", appliedToEntityTypeParam)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("ProfileTypes", profileTypesParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.IpfixUpmProfileListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ipfix_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.IpfixUpmProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), IpfixProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.IpfixUpmProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (iIface *ipfixProfilesClient) Update(ipfixProfileIdParam string, ipfixUpmProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := iIface.connector.TypeConverter()
	executionContext := iIface.connector.NewExecutionContext()
	operationRestMetaData := ipfixProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(ipfixProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("IpfixProfileId", ipfixProfileIdParam)
	sv.AddStructField("IpfixUpmProfile", ipfixUpmProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := iIface.connector.GetApiProvider().Invoke("com.vmware.nsx.ipfix_profiles", "update", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), IpfixProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), iIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
