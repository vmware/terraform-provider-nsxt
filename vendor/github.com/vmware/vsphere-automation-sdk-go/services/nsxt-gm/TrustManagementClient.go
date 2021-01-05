/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: TrustManagement
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type TrustManagementClient interface {

    // Returns information about the supported algorithms and key sizes.
    // @return com.vmware.nsx_global_policy.model.TrustManagementData
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get() (model.TrustManagementData, error)
}
