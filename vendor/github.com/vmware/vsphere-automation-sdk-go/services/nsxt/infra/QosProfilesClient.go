// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: QosProfiles
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

type QosProfilesClient interface {

	// API will delete QoS profile.
	//
	// @param qosProfileIdParam QoS profile Id (required)
	// @param overrideParam Locally override the global object (optional, default to false)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(qosProfileIdParam string, overrideParam *bool) error

	// API will return details of QoS profile.
	//
	// @param qosProfileIdParam QoS profile Id (required)
	// @return com.vmware.nsx_policy.model.QoSProfile
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(qosProfileIdParam string) (model.QosProfile, error)

	// API will list all QoS profiles.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.QoSProfileListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.QosProfileListResult, error)

	// Create a new QoS profile if the QoS profile with given id does not already exist. If the QoS profile with the given id already exists, patch with the existing QoS profile.
	//
	// @param qosProfileIdParam QoS profile Id (required)
	// @param qosProfileParam (required)
	// @param overrideParam Locally override the global object (optional, default to false)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(qosProfileIdParam string, qosProfileParam model.QosProfile, overrideParam *bool) error

	// Create or Replace QoS profile.
	//
	// @param qosProfileIdParam QoS profile Id (required)
	// @param qosProfileParam (required)
	// @param overrideParam Locally override the global object (optional, default to false)
	// @return com.vmware.nsx_policy.model.QoSProfile
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(qosProfileIdParam string, qosProfileParam model.QosProfile, overrideParam *bool) (model.QosProfile, error)
}

type qosProfilesClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewQosProfilesClient(connector client.Connector) *qosProfilesClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.qos_profiles")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":   core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	qIface := qosProfilesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &qIface
}

func (qIface *qosProfilesClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := qIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (qIface *qosProfilesClient) Delete(qosProfileIdParam string, overrideParam *bool) error {
	typeConverter := qIface.connector.TypeConverter()
	executionContext := qIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(qosProfilesDeleteInputType(), typeConverter)
	sv.AddStructField("QosProfileId", qosProfileIdParam)
	sv.AddStructField("Override", overrideParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := qosProfilesDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	qIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := qIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.qos_profiles", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), qIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (qIface *qosProfilesClient) Get(qosProfileIdParam string) (model.QosProfile, error) {
	typeConverter := qIface.connector.TypeConverter()
	executionContext := qIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(qosProfilesGetInputType(), typeConverter)
	sv.AddStructField("QosProfileId", qosProfileIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.QosProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := qosProfilesGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	qIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := qIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.qos_profiles", "get", inputDataValue, executionContext)
	var emptyOutput model.QosProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), qosProfilesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.QosProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), qIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (qIface *qosProfilesClient) List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.QosProfileListResult, error) {
	typeConverter := qIface.connector.TypeConverter()
	executionContext := qIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(qosProfilesListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.QosProfileListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := qosProfilesListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	qIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := qIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.qos_profiles", "list", inputDataValue, executionContext)
	var emptyOutput model.QosProfileListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), qosProfilesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.QosProfileListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), qIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (qIface *qosProfilesClient) Patch(qosProfileIdParam string, qosProfileParam model.QosProfile, overrideParam *bool) error {
	typeConverter := qIface.connector.TypeConverter()
	executionContext := qIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(qosProfilesPatchInputType(), typeConverter)
	sv.AddStructField("QosProfileId", qosProfileIdParam)
	sv.AddStructField("QosProfile", qosProfileParam)
	sv.AddStructField("Override", overrideParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := qosProfilesPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	qIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := qIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.qos_profiles", "patch", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), qIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (qIface *qosProfilesClient) Update(qosProfileIdParam string, qosProfileParam model.QosProfile, overrideParam *bool) (model.QosProfile, error) {
	typeConverter := qIface.connector.TypeConverter()
	executionContext := qIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(qosProfilesUpdateInputType(), typeConverter)
	sv.AddStructField("QosProfileId", qosProfileIdParam)
	sv.AddStructField("QosProfile", qosProfileParam)
	sv.AddStructField("Override", overrideParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.QosProfile
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := qosProfilesUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	qIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := qIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.qos_profiles", "update", inputDataValue, executionContext)
	var emptyOutput model.QosProfile
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), qosProfilesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.QosProfile), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), qIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
