// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Csrs
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

type CsrsClient interface {

	// Creates a new certificate signing request (CSR). A CSR is encrypted text that contains information about your organization (organization name, country, and so on) and your Web server's public key, which is a public certificate the is generated on the server that can be used to forward this request to a certificate authority (CA). A private key is also usually created at the same time as the CSR.
	//
	// @param csrIdParam ID of CSR to create (required)
	// @param tlsCsrParam (required)
	// @return com.vmware.nsx_policy.model.TlsCsr
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Create(csrIdParam string, tlsCsrParam model.TlsCsr) (model.TlsCsr, error)

	// Removes a specified CSR. If a CSR is not used for verification, you can delete it. Note that the CSR import and upload POST actions automatically delete the associated CSR.
	//
	// @param csrIdParam ID of CSR to delete (required)
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(csrIdParam string) error

	// Returns information about the specified CSR.
	//
	// @param csrIdParam ID of CSR to read (required)
	// @return com.vmware.nsx_policy.model.TlsCsr
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(csrIdParam string) (model.TlsCsr, error)

	// Imports a certificate authority (CA)-signed certificate for a CSR. This action links the certificate to the private key created by the CSR. The pem_encoded string in the request body is the signed certificate provided by your CA in response to the CSR that you provide to them. The import POST action automatically deletes the associated CSR.
	//
	// @param csrIdParam CSR this certificate is associated with (required)
	// @param tlsTrustDataParam (required)
	// @return com.vmware.nsx_policy.model.TlsCertificate
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Importcsr(csrIdParam string, tlsTrustDataParam model.TlsTrustData) (model.TlsCertificate, error)

	// Returns information about all of the CSRs that have been created.
	//
	// @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
	// @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
	// @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
	// @param sortAscendingParam (optional)
	// @param sortByParam Field by which records are sorted (optional)
	// @return com.vmware.nsx_policy.model.TlsCsrListResult
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.TlsCsrListResult, error)

	// Self-signs the previously generated CSR. This action is similar to the import certificate action, but instead of using a public certificate signed by a CA, the self_sign POST action uses a certificate that is signed with NSX's own private key. For validity of non-CA certificates, if a value greater than 825 days is provided, it will be set to 825 days. No limit is set for CA certificates.
	//
	// @param csrIdParam CSR this certificate is associated with (required)
	// @param daysValidParam Number of days the certificate will be valid, default 825 days (required)
	// @return com.vmware.nsx_policy.model.TlsCertificate
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Selfsign(csrIdParam string, daysValidParam int64) (model.TlsCertificate, error)

	// Creates a new self-signed certificate. A private key is also created at the same time. This is convenience call that will generate a CSR and then self-sign it. For validity of non-CA certificates, if a value greater than 825 days is provided, it will be set to 825 days. No limit is set for CA certificates.
	//
	// @param tlsCsrWithDaysValidParam (required)
	// @return com.vmware.nsx_policy.model.TlsCertificate
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Selfsign0(tlsCsrWithDaysValidParam model.TlsCsrWithDaysValid) (model.TlsCertificate, error)
}

type csrsClient struct {
	connector           client.Connector
	interfaceDefinition core.InterfaceDefinition
	errorsBindingMap    map[string]bindings.BindingType
}

func NewCsrsClient(connector client.Connector) *csrsClient {
	interfaceIdentifier := core.NewInterfaceIdentifier("com.vmware.nsx_policy.infra.csrs")
	methodIdentifiers := map[string]core.MethodIdentifier{
		"create":     core.NewMethodIdentifier(interfaceIdentifier, "create"),
		"delete":     core.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":        core.NewMethodIdentifier(interfaceIdentifier, "get"),
		"importcsr":  core.NewMethodIdentifier(interfaceIdentifier, "importcsr"),
		"list":       core.NewMethodIdentifier(interfaceIdentifier, "list"),
		"selfsign":   core.NewMethodIdentifier(interfaceIdentifier, "selfsign"),
		"selfsign_0": core.NewMethodIdentifier(interfaceIdentifier, "selfsign_0"),
	}
	interfaceDefinition := core.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]bindings.BindingType)

	cIface := csrsClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &cIface
}

func (cIface *csrsClient) GetErrorBindingType(errorName string) bindings.BindingType {
	if entry, ok := cIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return errors.ERROR_BINDINGS_MAP[errorName]
}

func (cIface *csrsClient) Create(csrIdParam string, tlsCsrParam model.TlsCsr) (model.TlsCsr, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsCreateInputType(), typeConverter)
	sv.AddStructField("CsrId", csrIdParam)
	sv.AddStructField("TlsCsr", tlsCsrParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCsr
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsCreateRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "create", inputDataValue, executionContext)
	var emptyOutput model.TlsCsr
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsCreateOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCsr), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *csrsClient) Delete(csrIdParam string) error {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsDeleteInputType(), typeConverter)
	sv.AddStructField("CsrId", csrIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsDeleteRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return bindings.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (cIface *csrsClient) Get(csrIdParam string) (model.TlsCsr, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsGetInputType(), typeConverter)
	sv.AddStructField("CsrId", csrIdParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCsr
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsGetRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "get", inputDataValue, executionContext)
	var emptyOutput model.TlsCsr
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCsr), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *csrsClient) Importcsr(csrIdParam string, tlsTrustDataParam model.TlsTrustData) (model.TlsCertificate, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsImportcsrInputType(), typeConverter)
	sv.AddStructField("CsrId", csrIdParam)
	sv.AddStructField("TlsTrustData", tlsTrustDataParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCertificate
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsImportcsrRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "importcsr", inputDataValue, executionContext)
	var emptyOutput model.TlsCertificate
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsImportcsrOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCertificate), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *csrsClient) List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.TlsCsrListResult, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsListInputType(), typeConverter)
	sv.AddStructField("Cursor", cursorParam)
	sv.AddStructField("IncludedFields", includedFieldsParam)
	sv.AddStructField("PageSize", pageSizeParam)
	sv.AddStructField("SortAscending", sortAscendingParam)
	sv.AddStructField("SortBy", sortByParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCsrListResult
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsListRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "list", inputDataValue, executionContext)
	var emptyOutput model.TlsCsrListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsListOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCsrListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *csrsClient) Selfsign(csrIdParam string, daysValidParam int64) (model.TlsCertificate, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsSelfsignInputType(), typeConverter)
	sv.AddStructField("CsrId", csrIdParam)
	sv.AddStructField("DaysValid", daysValidParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCertificate
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsSelfsignRestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "selfsign", inputDataValue, executionContext)
	var emptyOutput model.TlsCertificate
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsSelfsignOutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCertificate), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (cIface *csrsClient) Selfsign0(tlsCsrWithDaysValidParam model.TlsCsrWithDaysValid) (model.TlsCertificate, error) {
	typeConverter := cIface.connector.TypeConverter()
	executionContext := cIface.connector.NewExecutionContext()
	sv := bindings.NewStructValueBuilder(csrsSelfsign0InputType(), typeConverter)
	sv.AddStructField("TlsCsrWithDaysValid", tlsCsrWithDaysValidParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput model.TlsCertificate
		return emptyOutput, bindings.VAPIerrorsToError(inputError)
	}
	operationRestMetaData := csrsSelfsign0RestMetadata()
	connectionMetadata := map[string]interface{}{lib.REST_METADATA: operationRestMetaData}
	connectionMetadata["isStreamingResponse"] = false
	cIface.connector.SetConnectionMetadata(connectionMetadata)
	methodResult := cIface.connector.GetApiProvider().Invoke("com.vmware.nsx_policy.infra.csrs", "selfsign_0", inputDataValue, executionContext)
	var emptyOutput model.TlsCertificate
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), csrsSelfsign0OutputType())
		if errorInOutput != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInOutput)
		}
		return output.(model.TlsCertificate), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), cIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, bindings.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
