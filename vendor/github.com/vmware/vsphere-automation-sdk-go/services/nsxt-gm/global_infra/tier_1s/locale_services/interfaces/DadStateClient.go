/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DadState
 * Used by client-side stubs.
 */

package interfaces

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type DadStateClient interface {

    // Get tier-1 interface DAD state information.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_global_policy.model.InterfaceDADState
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServiceIdParam string, interfaceIdParam string, enforcementPointPathParam *string) (model.InterfaceDADState, error)
}
