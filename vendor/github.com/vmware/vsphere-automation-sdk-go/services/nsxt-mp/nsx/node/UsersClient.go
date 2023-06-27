// Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause

// Auto generated code. DO NOT EDIT.

// Interface file for service: Users
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

type UsersClient interface {

	// Activates the account for this user. When an account is successfully activated, the \"status\" field in the response is \"ACTIVE\". This API is not supported for userid 0 and userid 10000.
	//
	// @param useridParam User id of the user (required)
	// @param nodeUserPasswordPropertyParam (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Activate(useridParam string, nodeUserPasswordPropertyParam nsxModel.NodeUserPasswordProperty) (nsxModel.NodeUserProperties, error)

	// Create new user account to log in to the NSX web-based user interface or access API. ``username`` is required field in case of creating new user, further following usernames - ``root, admin, audit`` are reserved and can not be used to create new user account unless for local audit user. In case of local audit account when username not specified in request by default account will be created with ``audit`` username, although administrators are allowed to use any other non-duplicate usernames during creation.
	//
	// @param nodeUserPropertiesParam (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Createaudituser(nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error)

	// Create new user account to log in to the NSX web-based user interface or access API. ``username`` is required field in case of creating new user, further following usernames - ``root, admin, audit`` are reserved and can not be used to create new user account unless for local audit user. In case of local audit account when username not specified in request by default account will be created with ``audit`` username, although administrators are allowed to use any other non-duplicate usernames during creation.
	//
	// @param nodeUserPropertiesParam (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Createuser(nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error)

	// Deactivates the account for this user. Deactivating an account is permanent, unlike an account that is temporarily locked because of too many password failures. A deactivated account has to be explicitly activated. When an account is successfully deactivated, the \"status\" field in the response is \"NOT_ACTIVATED\". This API is not supported for userid 0 and userid 10000.
	//
	// @param useridParam User id of the user (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Deactivate(useridParam string) (nsxModel.NodeUserProperties, error)

	// Delete specified user who is configured to log in to the NSX appliance. Whereas local users root and administrator are not allowed to be deleted, but local user audit is deletable on-demand.
	//
	// **Caution**, users deleted from following node types cannot be recovered, kindly plan the removal of user accounts accordingly.
	//
	// * Autonomous Edge
	// * Cloud Service Manager
	// * Edge
	// * Public Cloud Gateway
	//
	// @param useridParam User id of the user (required)
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Delete(useridParam string) error

	// Returns information about a specified user who is configured to log in to the NSX appliance. The valid user IDs are: 0, 10000, 10002 or other users managed by administrators.
	//
	// @param useridParam User id of the user (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Get(useridParam string) (nsxModel.NodeUserProperties, error)

	// Returns the list of users configured to log in to the NSX appliance.
	// @return com.vmware.nsx.model.NodeUserPropertiesListResult
	//
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	List() (nsxModel.NodeUserPropertiesListResult, error)

	// Enables a user to reset their own password.
	//
	// @param resetNodeUserOwnPasswordPropertiesParam (required)
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Resetownpassword(resetNodeUserOwnPasswordPropertiesParam nsxModel.ResetNodeUserOwnPasswordProperties) error

	// Unlike the PUT version of this call (PUT /node/users/<userid>), this API does not require that the current password for the user be provided. The account of the target user must be \"ACTIVE\" for the call to succeed. This API is not supported for userid 0 and userid 10000.
	//
	// @param useridParam User id of the user (required)
	// @param nodeUserPasswordPropertyParam (required)
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Resetpassword(useridParam string, nodeUserPasswordPropertyParam nsxModel.NodeUserPasswordProperty) error

	//
	//
	// Updates attributes of an existing NSX appliance user. This method cannot be used to add a new user. Modifiable attributes include the username, full name of the user, and password. If you specify a password in a PUT request, it is not returned in the response. Nor is it returned in a GET request.
	//
	// The specified password does not meet the following (default) complexity requirements: - minimum 12 characters in length - minimum 128 characters in length - minimum 1 uppercase character - minimum 1 lowercase character - minimum 1 numeric character - minimum 1 special character - minimum 5 unique characters - default password complexity rules as enforced by the Linux PAM module *the configured password complexity may vary as per defined Authentication and Password policies, which shall be available at: [GET]: /api/v1/node/aaa/auth-policy*
	//
	//
	// The valid user IDs are: 0, 10000, 10002 or other users managed by administrators. Note that invoking this API does not update any user-related properties of existing objects in the system and does not modify the username field in existing audit log entries.
	//
	// @param useridParam User id of the user (required)
	// @param nodeUserPropertiesParam (required)
	// @return com.vmware.nsx.model.NodeUserProperties
	//
	// @throws ConcurrentChange  Conflict
	// @throws InvalidRequest  Bad Request, Precondition Failed
	// @throws Unauthorized  Forbidden
	// @throws ServiceUnavailable  Service Unavailable
	// @throws InternalServerError  Internal Server Error
	// @throws NotFound  Not Found
	Update(useridParam string, nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error)
}

type usersClient struct {
	connector           vapiProtocolClient_.Connector
	interfaceDefinition vapiCore_.InterfaceDefinition
	errorsBindingMap    map[string]vapiBindings_.BindingType
}

func NewUsersClient(connector vapiProtocolClient_.Connector) *usersClient {
	interfaceIdentifier := vapiCore_.NewInterfaceIdentifier("com.vmware.nsx.node.users")
	methodIdentifiers := map[string]vapiCore_.MethodIdentifier{
		"activate":         vapiCore_.NewMethodIdentifier(interfaceIdentifier, "activate"),
		"createaudituser":  vapiCore_.NewMethodIdentifier(interfaceIdentifier, "createaudituser"),
		"createuser":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "createuser"),
		"deactivate":       vapiCore_.NewMethodIdentifier(interfaceIdentifier, "deactivate"),
		"delete":           vapiCore_.NewMethodIdentifier(interfaceIdentifier, "delete"),
		"get":              vapiCore_.NewMethodIdentifier(interfaceIdentifier, "get"),
		"list":             vapiCore_.NewMethodIdentifier(interfaceIdentifier, "list"),
		"resetownpassword": vapiCore_.NewMethodIdentifier(interfaceIdentifier, "resetownpassword"),
		"resetpassword":    vapiCore_.NewMethodIdentifier(interfaceIdentifier, "resetpassword"),
		"update":           vapiCore_.NewMethodIdentifier(interfaceIdentifier, "update"),
	}
	interfaceDefinition := vapiCore_.NewInterfaceDefinition(interfaceIdentifier, methodIdentifiers)
	errorsBindingMap := make(map[string]vapiBindings_.BindingType)

	uIface := usersClient{interfaceDefinition: interfaceDefinition, errorsBindingMap: errorsBindingMap, connector: connector}
	return &uIface
}

func (uIface *usersClient) GetErrorBindingType(errorName string) vapiBindings_.BindingType {
	if entry, ok := uIface.errorsBindingMap[errorName]; ok {
		return entry
	}
	return vapiStdErrors_.ERROR_BINDINGS_MAP[errorName]
}

func (uIface *usersClient) Activate(useridParam string, nodeUserPasswordPropertyParam nsxModel.NodeUserPasswordProperty) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersActivateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersActivateInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	sv.AddStructField("NodeUserPasswordProperty", nodeUserPasswordPropertyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "activate", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersActivateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) Createaudituser(nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersCreateaudituserRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersCreateaudituserInputType(), typeConverter)
	sv.AddStructField("NodeUserProperties", nodeUserPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "createaudituser", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersCreateaudituserOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) Createuser(nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersCreateuserRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersCreateuserInputType(), typeConverter)
	sv.AddStructField("NodeUserProperties", nodeUserPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "createuser", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersCreateuserOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) Deactivate(useridParam string) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersDeactivateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersDeactivateInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "deactivate", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersDeactivateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) Delete(useridParam string) error {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersDeleteRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersDeleteInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "delete", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (uIface *usersClient) Get(useridParam string) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersGetRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersGetInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "get", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersGetOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) List() (nsxModel.NodeUserPropertiesListResult, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersListRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersListInputType(), typeConverter)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserPropertiesListResult
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "list", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserPropertiesListResult
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersListOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserPropertiesListResult), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}

func (uIface *usersClient) Resetownpassword(resetNodeUserOwnPasswordPropertiesParam nsxModel.ResetNodeUserOwnPasswordProperties) error {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersResetownpasswordRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersResetownpasswordInputType(), typeConverter)
	sv.AddStructField("ResetNodeUserOwnPasswordProperties", resetNodeUserOwnPasswordPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "resetownpassword", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (uIface *usersClient) Resetpassword(useridParam string, nodeUserPasswordPropertyParam nsxModel.NodeUserPasswordProperty) error {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersResetpasswordRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersResetpasswordInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	sv.AddStructField("NodeUserPasswordProperty", nodeUserPasswordPropertyParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		return vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "resetpassword", inputDataValue, executionContext)
	if methodResult.IsSuccess() {
		return nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return methodError.(error)
	}
}

func (uIface *usersClient) Update(useridParam string, nodeUserPropertiesParam nsxModel.NodeUserProperties) (nsxModel.NodeUserProperties, error) {
	typeConverter := uIface.connector.TypeConverter()
	executionContext := uIface.connector.NewExecutionContext()
	operationRestMetaData := usersUpdateRestMetadata()
	executionContext.SetConnectionMetadata(vapiCore_.RESTMetadataKey, operationRestMetaData)
	executionContext.SetConnectionMetadata(vapiCore_.ResponseTypeKey, vapiCore_.NewResponseType(true, false))

	sv := vapiBindings_.NewStructValueBuilder(usersUpdateInputType(), typeConverter)
	sv.AddStructField("Userid", useridParam)
	sv.AddStructField("NodeUserProperties", nodeUserPropertiesParam)
	inputDataValue, inputError := sv.GetStructValue()
	if inputError != nil {
		var emptyOutput nsxModel.NodeUserProperties
		return emptyOutput, vapiBindings_.VAPIerrorsToError(inputError)
	}

	methodResult := uIface.connector.GetApiProvider().Invoke("com.vmware.nsx.node.users", "update", inputDataValue, executionContext)
	var emptyOutput nsxModel.NodeUserProperties
	if methodResult.IsSuccess() {
		output, errorInOutput := typeConverter.ConvertToGolang(methodResult.Output(), UsersUpdateOutputType())
		if errorInOutput != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInOutput)
		}
		return output.(nsxModel.NodeUserProperties), nil
	} else {
		methodError, errorInError := typeConverter.ConvertToGolang(methodResult.Error(), uIface.GetErrorBindingType(methodResult.Error().Name()))
		if errorInError != nil {
			return emptyOutput, vapiBindings_.VAPIerrorsToError(errorInError)
		}
		return emptyOutput, methodError.(error)
	}
}
