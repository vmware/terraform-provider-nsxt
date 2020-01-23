/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package security


import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

func RetrieveUserIdentity(ctx *core.ExecutionContext) (*UserIdentity, error) {

	securityCtx := ctx.SecurityContext()
	if securityCtx == nil {
		log.Error("securityCtx is empty when trying to retrieve UserIdentity")
		return nil, l10n.NewRuntimeError("vapi.security.authentication.failed", make(map[string]string))
	}
	if temp := securityCtx.Property(AUTHN_IDENTITY); temp != nil {
		if userId, ok := temp.(*UserIdentity); ok {
			return userId, nil
		} else {
			log.Error("SecurityCtx property AUTHN_IDENTITY failed to assert (*UserIdentity)")
		}
	} else {
		log.Error("SecurityCtx doesn't have property AUTHN_IDENTITY")
	}

	return nil, l10n.NewRuntimeError("vapi.security.authentication.failed", make(map[string]string))
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
