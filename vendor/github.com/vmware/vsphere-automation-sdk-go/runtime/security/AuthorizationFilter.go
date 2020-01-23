/* Copyright © 2019 VMware, Inc. All Rights Reserved.
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
		return a.getInvalidAuthzMethodResult(nil)
	}

	if !a.pValidator.Validate(userName, groupNames, requiredPrivileges) {
		return a.getInvalidAuthzMethodResult(nil)
	}

	for _, authzHandler := range a.handlers {
		//TODO invoke only those handlers which support auth schemes.
		authzResult, authzError := authzHandler.Authorize(serviceID, operationID, ctx.SecurityContext())
		if authzError != nil {
			// authz failed.
			return a.getInvalidAuthzMethodResult(authzError)
		} else if authzResult {
			return a.provider.Invoke(serviceID, operationID, inputValue, ctx)
		}
	}
	return a.getInvalidAuthzMethodResult(nil)
}

func (a *AuthorizationFilter) getInvalidAuthzMethodResult(authzError error) core.MethodResult {
	if authzError != nil {
		errorValue := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHORIZED_ERROR_DEF,
			"vapi.security.authorization.invalid", nil)

		return core.NewMethodResult(nil, errorValue)
	}
	errorValue := bindings.CreateErrorValueFromMessageId(bindings.UNAUTHORIZED_ERROR_DEF,
		"vapi.security.authorization.invalid", nil)

	return core.NewMethodResult(nil, errorValue)
}
