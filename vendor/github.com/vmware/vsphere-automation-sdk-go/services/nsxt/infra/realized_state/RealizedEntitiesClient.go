/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: RealizedEntities
 * Used by client-side stubs.
 */

package realized_state

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type RealizedEntitiesClient interface {

    // Get list of realized entities associated with intent object, specified by path in query parameter
    //
    // @param intentPathParam String Path of the intent object (required)
    // @param sitePathParam Policy Path of the site (optional)
    // @return com.vmware.nsx_policy.model.GenericPolicyRealizedResourceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(intentPathParam string, sitePathParam *string) (model.GenericPolicyRealizedResourceListResult, error)
}
