/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GmOperationalState
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type GmOperationalStateClient interface {

    // Global Manager operation state includes the current status, switchover status of global manager nodes if any, errors if any and consolidated status of the operation.
    // @return com.vmware.nsx_global_policy.model.GmOperationalState
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get() (model.GmOperationalState, error)
}
