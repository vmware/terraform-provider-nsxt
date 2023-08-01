// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Licenses
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

type LicensesClient interface {

	// This will add a license key to the system. The API supports adding only one license key for each license edition type - Standard, Advanced or Enterprise. If a new license key is tried to add for an edition for which the license key already exists, then this API will return an error.
	//
	// @param licenseParam (required)
	// @return com.vmware.nsx.model.License
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(licenseParam nsxModel.License) (nsxModel.License, error)

	// Remove a license identified by the license-key.
	//  Deprecated. Use POST /licenses?action=delete API instead.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param licenseKeyParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(licenseKeyParam string) error

	// This will delete the license key identified in the request body by \"license_key\" and its properties from the system. Attempting to delete the last license key will result in an error.
	//
	// @param licenseParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete0(licenseParam nsxModel.License) error

	// Return the Enterprise License.
	//  Deprecated. Use the GET /licenses API instead.
	//
	// Deprecated: This API element is deprecated.
	// @return com.vmware.nsx.model.License
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get() (nsxModel.License, error)

	// Get license properties for license identified by the license-key.
	//  Use GET /licenses API instead.
	//
	// Deprecated: This API element is deprecated.
	//
	// @param licenseKeyParam (required)
	// @return com.vmware.nsx.model.License
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Getlicensebykey(licenseKeyParam string) (nsxModel.License, error)

	// Returns all licenses.
	// @return com.vmware.nsx.model.LicensesListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List() (nsxModel.LicensesListResult, error)

	// Assign an Updated Enterprise License Key.
	//  Deprecated. Use the POST /licenses API instead
	//
	// Deprecated: This API element is deprecated.
	//
	// @param licenseParam (required)
	// @return com.vmware.nsx.model.License
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(licenseParam nsxModel.License) (nsxModel.License, error)
}

type licensesClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewLicensesClient(connector vapiProtocolClient_.Connector) *licensesClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.licenses")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"create":          vapiCore_.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete":          vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"delete_0":        vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete_0"),
		"get":             vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"getlicensebykey": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "getlicensebykey"),
		"list":            vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"update":          vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	lIface := licensesClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &lIface
}

func (lIface *licensesClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := lIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (lIface *licensesClient) Create(licenseParam nsxModel.License) (nsxModel.License, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesCreateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesCreateInputType(), typeConverter)
	sv.AddStructField("License", licenseParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.License
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "create", inputDataValue, executionContext)
	var emptyOutput nsxModel.License
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LicensesCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.License), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *licensesClient) Delete(licenseKeyParam string) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesDeleteInputType(), typeConverter)
	sv.AddStructField("LicenseKey", licenseKeyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (lIface *licensesClient) Delete0(licenseParam nsxModel.License) error {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesDelete0RestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesDelete0InputType(), typeConverter)
	sv.AddStructField("License", licenseParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "delete_0", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (lIface *licensesClient) Get() (nsxModel.License, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesGetInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.License
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.License
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LicensesGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.License), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *licensesClient) Getlicensebykey(licenseKeyParam string) (nsxModel.License, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesGetlicensebykeyRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesGetlicensebykeyInputType(), typeConverter)
	sv.AddStructField("LicenseKey", licenseKeyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.License
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "getlicensebykey", inputDataValue, executionContext)
	var emptyOutput nsxModel.License
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LicensesGetlicensebykeyOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.License), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *licensesClient) List() (nsxModel.LicensesListResult, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesListInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.LicensesListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.LicensesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LicensesListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.LicensesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (lIface *licensesClient) Update(licenseParam nsxModel.License) (nsxModel.License, error) {
	typeConverter := lIface.connector.TypeConverter()
	executionContext := lIface.connector.NewExecutionContext()
	operationRestMetaData := licensesUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(licensesUpdateInputType(), typeConverter)
	sv.AddStructField("License", licenseParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.License
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := lIface.connector.GetApiProvider().Invoke("com.vmware.nsx.licenses", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.License
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), LicensesUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.License), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), lIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
