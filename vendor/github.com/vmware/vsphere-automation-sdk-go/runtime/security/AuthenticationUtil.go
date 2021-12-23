/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security

import (
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

func RetrieveUserIdentity(ctx *core.ExecutionContext) (*UserIdentity, error) {
	securityCtx := ctx.SecurityContext()
	if securityCtx == nil {
		errMsg := "securityCtx is empty when trying to retrieve UserIdentity"
		log.Error(errMsg)
		return nil, l10n.NewRuntimeError(
			"vapi.security.authentication.failed.missing_security_context",
			nil)
	}

	authnIdentity := securityCtx.Property(AUTHN_IDENTITY)
	if authnIdentity == nil {
		log.Infof("SecurityCtx doesn't contain authentication information."+
			"Missing '%s' property.", AUTHN_IDENTITY)
		return nil, nil
	}

	var (
		ok           bool
		userIdentity *UserIdentity
	)
	if userIdentity, ok = authnIdentity.(*UserIdentity); !ok {
		log.Info(fmt.Sprintf("SecurityCtx property %s failed to "+
			"assert (*UserIdentity)", AUTHN_IDENTITY))
		return nil, l10n.NewRuntimeError(
			"vapi.security.authentication.failed.identity_convert_error",
			nil)
	}

	return userIdentity, nil
}

func RetrieveUserName(userId *UserIdentity) string {
	name := userId.UserName()
	domain := userId.Domain()
	if domain != nil {
		log.Debug("User name and domain: ", name, domain)
		return *domain + "\\" + name
	} else {
		log.Debug("User name: ", name)
		return name
	}
}

func RetrieveUserGroups(userId *UserIdentity) []string {
	return userId.Groups()
}
