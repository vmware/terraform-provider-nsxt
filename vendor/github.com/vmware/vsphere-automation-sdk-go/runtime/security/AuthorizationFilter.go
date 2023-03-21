/* Copyright © 2019-2022 VMware, Inc. All Rights Reserved.
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
	if err != nil {
		log.Errorf("Could not find user identity for operation '%s' : %s", fullyQualifiedOperName, err)
		return newInternalServerErrorResult("vapi.security.authorization.user.missing", nil)
	}
	if userId == nil {
		// No security context, or not authn data (because method doesn't require authn) => no authentication;
		// without user information, no authorization as well
		log.Debugf("Skipping authorization checks, "+
			"because there is no authentication data for: %s", fullyQualifiedOperName)
		return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
	}
	userName := RetrieveUserName(userId)
	groupNames := RetrieveUserGroups(userId)

	requiredPrivileges, err := a.privilegeProv.GetPrivilegeInfo(fullyQualifiedOperName, inputValue)
	if err != nil {
		log.Errorf("Could not retrieve privilege information for operation '%s' : %s", fullyQualifiedOperName, err)
		// Allow the PrivilegeProvider to return std vapi error directly to the client, e.g. ServiceUnavailable
		if errorValue := getErrorValue(err); errorValue != nil {
			return core.NewErrorResult(errorValue)
		}
		return newInternalServerErrorResult("vapi.security.authorization.privilege.error", nil)
	}

	isValid, err := a.pValidator.Validate(userName, groupNames, requiredPrivileges)
	if err != nil {
		log.Errorf("Could not validate permission for operation '%s' : %s", fullyQualifiedOperName, err)
		// Allow the PermissionValidator to return std vapi error directly to the client, e.g. ServiceUnavailable
		if errorValue := getErrorValue(err); errorValue != nil {
			return core.NewErrorResult(errorValue)
		}
		return newInternalServerErrorResult("vapi.security.authorization.permission.error", nil)
	}
	if !isValid {
		log.Debugf("Permission denied for operation '%s'", fullyQualifiedOperName)
		return newUnauthorizedResult("vapi.security.authorization.permission.denied", nil)
	}

	if len(a.handlers) == 0 {
		return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
	}

	for _, authzHandler := range a.handlers {
		//TODO invoke only those handlers which support auth schemes.
		authzResult, err := authzHandler.Authorize(serviceID, operationID, ctx.SecurityContext())
		if err != nil {
			log.Errorf("AuthorizationHandler failed for operation '%s' : %s", fullyQualifiedOperName, err)
			return newInternalServerErrorResult("vapi.security.authorization.handler.error", nil)
		} else if authzResult {
			return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
		}
	}
	return newUnauthorizedResult("vapi.security.authorization.invalid", nil)
}

func newUnauthorizedResult(msgId string, args map[string]string) core.MethodResult {
	return core.NewErrorResult(bindings.CreateErrorValueFromMessageId(bindings.UNAUTHORIZED_ERROR_DEF, msgId, args))
}

func newInternalServerErrorResult(msgId string, args map[string]string) core.MethodResult {
	return core.NewErrorResult(bindings.CreateErrorValueFromMessageId(bindings.INTERNAL_SERVER_ERROR_DEF, msgId, args))
}

func getErrorValue(e error) *data.ErrorValue {
	vapiError, ok := e.(bindings.Structure)
	if !ok {
		return nil
	}
	dataValue, err := vapiError.GetDataValue__()
	if err != nil {
		return nil
	}
	errorValue, ok := dataValue.(*data.ErrorValue)
	if !ok {
		return nil
	}
	return errorValue
}
