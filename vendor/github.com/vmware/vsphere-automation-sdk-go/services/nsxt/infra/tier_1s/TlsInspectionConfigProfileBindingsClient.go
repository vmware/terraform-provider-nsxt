// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: TlsInspectionConfigProfileBindings
// Used by client-side stubs.

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const _ = core.SupportedByRuntimeVersion1

type TlsInspectionConfigProfileBindingsClient interface {

	// API will delete TLS Config Profile Binding for Tier-1 Logical Router.
	//
	// @param tier1IdParam (required)
	// @param tlsInspectionConfigProfileBindingIdParam (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string) error

	// API will get TLS Config Profile Binding Map for Tier-1 Logical Router.
	//
	// @param tier1IdParam (required)
	// @param tlsInspectionConfigProfileBindingIdParam (required)
	// @return com.vmware.nsx_policy.model.TlsConfigProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string) (model.TlsConfigProfileBindingMap, error)

	// API will create or update TLS Config profile binding map for Tier-1 Logical Router.
	//
	// @param tier1IdParam (required)
	// @param tlsInspectionConfigProfileBindingIdParam (required)
	// @param tlsConfigProfileBindingMapParam (required)
	// @return com.vmware.nsx_policy.model.TlsConfigProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Patch(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string, tlsConfigProfileBindingMapParam model.TlsConfigProfileBindingMap) (model.TlsConfigProfileBindingMap, error)

	// API will create or update TLS Config profile binding map for Tier-1 Logical Router.
	//
	// @param tier1IdParam (required)
	// @param tlsInspectionConfigProfileBindingIdParam (required)
	// @param tlsConfigProfileBindingMapParam (required)
	// @return com.vmware.nsx_policy.model.TlsConfigProfileBindingMap
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string, tlsConfigProfileBindingMapParam model.TlsConfigProfileBindingMap) (model.TlsConfigProfileBindingMap, error)
}

type tlsInspectionConfigProfileBindingsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewTlsInspectionConfigProfileBindingsClient(connector client.Connector) *tlsInspectionConfigProfileBindingsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.tier_1s.tls_inspection_config_profile_bindings")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"delete": core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":    core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"patch":  core.NewMethodIdentifier(interfaceIdentifier, "patch"),
		"update": core.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	tIface := tlsInspectionConfigProfileBindingsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &tIface
}

func (tIface *tlsInspectionConfigProfileBindingsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := tIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (tIface *tlsInspectionConfigProfileBindingsClient) Delete(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string) error {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(tlsInspectionConfigProfileBindingsDeleteInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("TlsInspectionConfigProfileBindingId", tlsInspectionConfigProfileBindingIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := tlsInspectionConfigProfileBindingsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.tls_inspection_config_profile_bindings", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (tIface *tlsInspectionConfigProfileBindingsClient) Get(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string) (model.TlsConfigProfileBindingMap, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(tlsInspectionConfigProfileBindingsGetInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("TlsInspectionConfigProfileBindingId", tlsInspectionConfigProfileBindingIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsConfigProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := tlsInspectionConfigProfileBindingsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.tls_inspection_config_profile_bindings", "get", inputDataValue, executionContext)
	var emptyOutput model.TlsConfigProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), tlsInspectionConfigProfileBindingsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsConfigProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *tlsInspectionConfigProfileBindingsClient) Patch(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string, tlsConfigProfileBindingMapParam model.TlsConfigProfileBindingMap) (model.TlsConfigProfileBindingMap, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(tlsInspectionConfigProfileBindingsPatchInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("TlsInspectionConfigProfileBindingId", tlsInspectionConfigProfileBindingIdParam)
	sv.AddStructField("TlsConfigProfileBindingMap", tlsConfigProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsConfigProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := tlsInspectionConfigProfileBindingsPatchRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.tls_inspection_config_profile_bindings", "patch", inputDataValue, executionContext)
	var emptyOutput model.TlsConfigProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), tlsInspectionConfigProfileBindingsPatchOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsConfigProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (tIface *tlsInspectionConfigProfileBindingsClient) Update(tier1IdParam string, tlsInspectionConfigProfileBindingIdParam string, tlsConfigProfileBindingMapParam model.TlsConfigProfileBindingMap) (model.TlsConfigProfileBindingMap, error) {
	typeConverter := tIface.connector.TypeConverter()
	executionContext := tIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(tlsInspectionConfigProfileBindingsUpdateInputType(), typeConverter)
	sv.AddStructField("Tier1Id", tier1IdParam)
	sv.AddStructField("TlsInspectionConfigProfileBindingId", tlsInspectionConfigProfileBindingIdParam)
	sv.AddStructField("TlsConfigProfileBindingMap", tlsConfigProfileBindingMapParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsConfigProfileBindingMap
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := tlsInspectionConfigProfileBindingsUpdateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	tIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := tIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.tier_1s.tls_inspection_config_profile_bindings", "update", inputDataValue, executionContext)
	var emptyOutput model.TlsConfigProfileBindingMap
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), tlsInspectionConfigProfileBindingsUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsConfigProfileBindingMap), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), tIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
