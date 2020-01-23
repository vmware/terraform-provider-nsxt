/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbNodeUsage
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type LbNodeUsageClient interface {

    //
    //
    // @param nodePathParam The node path for load balancer node usage (required)
    // @return com.vmware.nsx_policy.model.LBNodeUsage
    // The return value will contain all the properties defined in model.LBNodeUsage.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(nodePathParam string) (*data.StructValue, error)
}
