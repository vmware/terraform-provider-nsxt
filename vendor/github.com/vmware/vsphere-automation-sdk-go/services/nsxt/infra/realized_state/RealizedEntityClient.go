/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: RealizedEntity
 * Used by client-side stubs.
 */

package realized_state

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type RealizedEntityClient interface {

    // Get realized entity uniquely identified by realized path, specified by query parameter
    //
    // @param realizedPathParam String Path of the realized object (required)
    // @return com.vmware.nsx_policy.model.GenericPolicyRealizedResource
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(realizedPathParam string) (model.GenericPolicyRealizedResource, error)

    // Refresh the status and statistics of all realized entities associated with given intent path synchronously. The vmw-async: True HTTP header cannot be used with this API.
    //
    // @param intentPathParam String Path of the intent object (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Refresh(intentPathParam string, enforcementPointPathParam *string) error
}
