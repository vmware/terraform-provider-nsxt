/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

// AuthorizationFilter in API Provider chain enforces the authorization
// schemes specified in the authorization metadata file

type AuthorizationFilter struct {
	handlers      []AuthorizationHandler
	provider      core.APIProvider
	privilegeProv PrivilegeProvider
	pValidator    PermissionValidator
}

func NewAuthorizationFilter(provider core.APIProvider, privilegeProv PrivilegeProvider, pValidator PermissionValidator) (*AuthorizationFilter, error) {
	aFilter := AuthorizationFilter{provider: provider, handlers: []AuthorizationHandler{}, privilegeProv: privilegeProv, pValidator: pValidator}
	return &aFilter, nil
}

func (a *AuthorizationFilter) AddHandler(handler AuthorizationHandler) {
	a.handlers = append(a.handlers, handler)
}

func (a *AuthorizationFilter) Invoke(serviceID string, operationID string,
	inputValue data.DataValue, ctx *core.ExecutionContext) core.MethodResult {

	fullyQualifiedOperName := serviceID + "." + operationID

	userId, err := RetrieveUserIdentity(ctx)
	if userId == nil {
		// No security context, or not authn data (because method doesn't require authn) => no authentication;
		// without user information, no authorization as well
		log.Debugf("Skipping authorization checks, because there is no authentication data for: " + fullyQualifiedOperName)
		return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
	}
	userName := RetrieveUserName(userId)
	groupNames := RetrieveUserGroups(userId)

	requiredPrivileges, err := a.privilegeProv.GetPrivilegeInfo(fullyQualifiedOperName, inputValue)
	if err != nil {
		return a.getInvalidAuthzMethodResult(err)
	}

	isValid, err := a.pValidator.Validate(userName, groupNames, requiredPrivileges)
	if !isValid || err != nil {
		return a.getInvalidAuthzMethodResult(err)
	}

	for _, authzHandler := range a.handlers {
		//TODO invoke only those handlers which support auth schemes.
		authzResult, err := authzHandler.Authorize(serviceID, operationID, ctx.SecurityContext())
		if err != nil {
			// authz failed.
			return a.getInvalidAuthzMethodResult(err)
		} else if authzResult {
			return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
		}
	}
	return a.getInvalidAuthzMethodResult(nil)
}

func (a *AuthorizationFilter) getInvalidAuthzMethodResult(err error) core.MethodResult {
	var errorValue *data.ErrorValue
	if err != nil {
		if vapiError, isVapiError := err.(bindings.Structure); isVapiError {
			dataVal, err := vapiError.GetDataValue__()
			if dataVal != nil && err == nil {
				errorValue = dataVal.(*data.ErrorValue)
			}
		}

		if errorValue == nil {
			args := map[string]string{"err": err.Error()}
			errorValue = bindings.CreateErrorValueFromMessageId(bindings.INTERNAL_SERVER_ERROR_DEF,
				"vapi.security.authorization.internal_server_error", args)
		}

		return core.NewMethodResult(nil, errorValue)
	}

	errorValue = bindings.CreateErrorValueFromMessageId(bindings.UNAUTHORIZED_ERROR_DEF,
		"vapi.security.authorization.invalid", nil)

	return core.NewMethodResult(nil, errorValue)
}
