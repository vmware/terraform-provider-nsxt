/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Status
 * Used by client-side stubs.
 */

package realized_state

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type StatusClient interface {

    //
    //
    // @param intentPathParam Policy Path of the intent object (required)
    // @param includeEnforcedStatusParam Include Enforced Status Flag (optional, default to false)
    // @param sitePathParam Policy Path of the site from where the realization status needs to be fetched (optional)
    // @return com.vmware.nsx_policy.model.ConsolidatedRealizedStatus
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(intentPathParam string, includeEnforcedStatusParam *bool, sitePathParam *string) (model.ConsolidatedRealizedStatus, error)
}
