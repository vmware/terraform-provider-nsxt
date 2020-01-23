/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Span
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SpanClient interface {

    // Get site span for Tier-0 Gateway.
    //
    // @param tier0IdParam (required)
    // @return com.vmware.nsx_policy.model.Sites
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string) (model.Sites, error)
}
