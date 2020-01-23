/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: L3vpnContext
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type L3vpnContextClient interface {

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @return com.vmware.nsx_policy.model.L3VpnContext
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string) (model.L3VpnContext, error)
}
