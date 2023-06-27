// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: FileStore
// Used by client-side stubs.

package node

import (
	vapiStdErrors_ "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiBindings_ "github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiCore_ "github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	vapiProtocolClient_ "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

const _ = vapiCore_.SupportedByRuntimeVersion2

type FileStoreClient interface {

	// Copy a remote file to the file store. If you use scp or sftp, you must provide the remote server's SSH fingerprint. See the *NSX-T Administration Guide* for information and instructions about finding the SSH fingerprint.
	//
	// @param fileNameParam Destination filename (required)
	// @param copyFromRemoteFilePropertiesParam (required)
	// @return com.vmware.nsx.model.FileProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Copyfromremotefile(fileNameParam string, copyFromRemoteFilePropertiesParam nsxModel.CopyFromRemoteFileProperties) (nsxModel.FileProperties, error)

	// Copy a file in the file store to a remote server. If you use scp or sftp, you must provide the remote server's SSH fingerprint. See the *NSX-T Administration Guide* for information and instructions about finding the SSH fingerprint.
	//
	// @param fileNameParam Destination filename (required)
	// @param copyToRemoteFilePropertiesParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Copytoremotefile(fileNameParam string, copyToRemoteFilePropertiesParam nsxModel.CopyToRemoteFileProperties) error

	// Create a directory on the remote remote server. Supports only SFTP. You must provide the remote server's SSH fingerprint. See the *NSX Administration Guide* for information and instructions about finding the SSH fingerprint.
	//
	// @param createRemoteDirectoryPropertiesParam (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws TimedOut  Gateway Timeout
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Createremotedirectory(createRemoteDirectoryPropertiesParam nsxModel.CreateRemoteDirectoryProperties) error

	// Delete file
	//
	// @param fileNameParam Name of the file to delete (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(fileNameParam string) error

	// Read file properties
	//
	// @param fileNameParam Name of the file to retrieve information about (required)
	// @return com.vmware.nsx.model.FileProperties
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(fileNameParam string) (nsxModel.FileProperties, error)

	// List node files
	// @return com.vmware.nsx.model.FilePropertiesListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List() (nsxModel.FilePropertiesListResult, error)

	// Retrieve ssh fingerprint for a given remote server and port.
	//
	// @param sshFingerprintPropertiesParam (required)
	// @return com.vmware.nsx.model.SshFingerprintProperties
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Retrievesshfingerprint(sshFingerprintPropertiesParam nsxModel.SshFingerprintProperties) (nsxModel.SshFingerprintProperties, error)
}

type fileStoreClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewFileStoreClient(connector vapiProtocolClient_.Connector) *fileStoreClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.node.file_store")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"copyfromremotefile":     vapiCore_.NewMethodIdentifier(interfaceIdentifier, "copyfromremotefile"),
		"copytoremotefile":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "copytoremotefile"),
		"createremotedirectory":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "createremotedirectory"),
		"delete":                 vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":                    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":                   vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"retrievesshfingerprint": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "retrievesshfingerprint"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	fIface := fileStoreClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &fIface
}

func (fIface *fileStoreClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := fIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (fIface *fileStoreClient) Copyfromremotefile(fileNameParam string, copyFromRemoteFilePropertiesParam nsxModel.CopyFromRemoteFileProperties) (nsxModel.FileProperties, error) {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreCopyfromremotefileRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreCopyfromremotefileInputType(), typeConverter)
	sv.AddStructField("FileName", fileNameParam)
	sv.AddStructField("CopyFromRemoteFileProperties", copyFromRemoteFilePropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.FileProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "copyfromremotefile", inputDataValue, executionContext)
	var emptyOutput nsxModel.FileProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), FileStoreCopyfromremotefileOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.FileProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (fIface *fileStoreClient) Copytoremotefile(fileNameParam string, copyToRemoteFilePropertiesParam nsxModel.CopyToRemoteFileProperties) error {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreCopytoremotefileRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreCopytoremotefileInputType(), typeConverter)
	sv.AddStructField("FileName", fileNameParam)
	sv.AddStructField("CopyToRemoteFileProperties", copyToRemoteFilePropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "copytoremotefile", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (fIface *fileStoreClient) Createremotedirectory(createRemoteDirectoryPropertiesParam nsxModel.CreateRemoteDirectoryProperties) error {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreCreateremotedirectoryRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreCreateremotedirectoryInputType(), typeConverter)
	sv.AddStructField("CreateRemoteDirectoryProperties", createRemoteDirectoryPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "createremotedirectory", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (fIface *fileStoreClient) Delete(fileNameParam string) error {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreDeleteInputType(), typeConverter)
	sv.AddStructField("FileName", fileNameParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (fIface *fileStoreClient) Get(fileNameParam string) (nsxModel.FileProperties, error) {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreGetInputType(), typeConverter)
	sv.AddStructField("FileName", fileNameParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.FileProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.FileProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), FileStoreGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.FileProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (fIface *fileStoreClient) List() (nsxModel.FilePropertiesListResult, error) {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreListInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.FilePropertiesListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.FilePropertiesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), FileStoreListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.FilePropertiesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (fIface *fileStoreClient) Retrievesshfingerprint(sshFingerprintPropertiesParam nsxModel.SshFingerprintProperties) (nsxModel.SshFingerprintProperties, error) {
	typeConverter := fIface.connector.TypeConverter()
	executionContext := fIface.connector.NewExecutionContext()
	operationRestMetaData := fileStoreRetrievesshfingerprintRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(fileStoreRetrievesshfingerprintInputType(), typeConverter)
	sv.AddStructField("SshFingerprintProperties", sshFingerprintPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.SshFingerprintProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := fIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.file_store", "retrievesshfingerprint", inputDataValue, executionContext)
	var emptyOutput nsxModel.SshFingerprintProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), FileStoreRetrievesshfingerprintOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.SshFingerprintProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), fIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
