// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: TransportzoneProfiles
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

type TransportzoneProfilesClient interface {

	// Creates a transport zone profile. The resource_type is required.
	//  This api is now deprecated. Please use new api - PUT /policy/api/v1/infra/transport-zone-profiles/<tz-profile>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param transportZoneProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.TransportZoneProfile.
	// @return com.vmware.nsx.model.TransportZoneProfile
	// The return value will contain all the properties defined in nsxModel.TransportZoneProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(transportZoneProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)

	// Deletes a specified transport zone profile.
	//  This api is now deprecated. Please use new api - DELETE /policy/api/v1/infra/transport-zone-profiles/<tz-profile>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param transportzoneProfileIdParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(transportzoneProfileIdParam string) error

	// Returns information about a specified transport zone profile.
	//  This api is now deprecated. Please use new api - /policy/api/v1/infra/transport-zone-profiles/<tz-profile>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param transportzoneProfileIdParam (required)
	// @return com.vmware.nsx.model.TransportZoneProfile
	// The return value will contain all the properties defined in nsxModel.TransportZoneProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(transportzoneProfileIdParam string) (*vapiData_.StructValue, error)

	// Returns information about the configured transport zone profiles. Transport zone profiles define networking policies for transport zones and transport zone endpoints.
	//  This api is now deprecated. Please use new api - /policy/api/v1/infra/transport-zone-profiles
	//
	// Deprecated: This API element is deprecated.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includeSystemOwnedParam Whether the list result contains system resources (optional, default to false)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param resourceTypeParam comma-separated list of transport zone profile types, e.g. ?resource_type=BfdHealthMonitoringProfile (optional)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx.model.TransportZoneProfileListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includeSystemOwnedParam *bool, includedFieldsParam *string, pageSizeParam *int64, resourceTypeParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.TransportZoneProfileListResult, error)

	// Modifies a specified transport zone profile. The body of the PUT request must include the resource_type.
	//  This api is now deprecated. Please use new api - PATCH /policy/api/v1/infra/transport-zone-profiles/<tz-profile>
	//
	// Deprecated: This API element is deprecated.
	//
	// @param transportzoneProfileIdParam (required)
	// @param transportZoneProfileParam (required)
	// The parameter must contain all the properties defined in nsxModel.TransportZoneProfile.
	// @return com.vmware.nsx.model.TransportZoneProfile
	// The return value will contain all the properties defined in nsxModel.TransportZoneProfile.
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(transportzoneProfileIdParam string, transportZoneProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error)
}

type transportzoneProfilesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewTransportzoneProfilesClient(connector vapiProtocolClient_.Connector) *transportzoneProfilesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.transportzone_profiles")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	tIface := transportzoneProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &tIface
}

func (tIface *transportzoneProfilesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := tIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (tIface *transportzoneProfilesClient) Create(transportZoneProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	operationRestMetaData := transportzoneProfilesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(transportzoneProfilesCreateInputType(), typeConverter)
	sv.AddStructField("TransportZoneProfile", transportZoneProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx.transportzone_profiles", "create", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), TransportzoneProfilesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *transportzoneProfilesClient) Delete(transportzoneProfileIdParam string) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	operationRestMetaData := transportzoneProfilesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(transportzoneProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("TransportzoneProfileId", transportzoneProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx.transportzone_profiles", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *transportzoneProfilesClient) Get(transportzoneProfileIdParam string) (*vapiData_.StructValue, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	operationRestMetaData := transportzoneProfilesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(transportzoneProfilesGetInputType(), typeConverter)
	sv.AddStructField("TransportzoneProfileId", transportzoneProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx.transportzone_profiles", "get", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), TransportzoneProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *transportzoneProfilesClient) List(cursorParam *string, includeSystemOwnedParam *bool, includedFieldsParam *string, pageSizeParam *int64, resourceTypeParam *string, sortAscendingParam *bool, sortByParam *string) (nsxModel.TransportZoneProfileListResult, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	operationRestMetaData := transportzoneProfilesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(transportzoneProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludeSystemOwned", includeSystemOwnedParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("ResourceType", resourceTypeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.TransportZoneProfileListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx.transportzone_profiles", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.TransportZoneProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), TransportzoneProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.TransportZoneProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *transportzoneProfilesClient) Update(transportzoneProfileIdParam string, transportZoneProfileParam *vapiData_.StructValue) (*vapiData_.StructValue, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	operationRestMetaData := transportzoneProfilesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(transportzoneProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("TransportzoneProfileId", transportzoneProfileIdParam)
	sv.AddStructField("TransportZoneProfile", transportZoneProfileParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput *vapiData_.StructValue
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx.transportzone_profiles", "update", inputDataValue, executionContext)
	var emptyOutput *vapiData_.StructValue
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), TransportzoneProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(*vapiData_.StructValue), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
