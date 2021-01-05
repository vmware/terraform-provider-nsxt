/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FirewallIdentityStoreLdapServer
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type FirewallIdentityStoreLdapServerClient interface {

    // This API tests a LDAP server connectivity before the actual domain or LDAP server is configured. If the connectivity is good, the response will be HTTP status 200. Otherwise the response will be HTTP status 500 and corresponding error message will be returned.
    //
    // @param directoryLdapServerParam (required)
    // @param actionParam LDAP server test requested (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_policy.model.DirectoryLdapServerStatus
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(directoryLdapServerParam model.DirectoryLdapServer, actionParam string, enforcementPointPathParam *string) (model.DirectoryLdapServerStatus, error)
}
