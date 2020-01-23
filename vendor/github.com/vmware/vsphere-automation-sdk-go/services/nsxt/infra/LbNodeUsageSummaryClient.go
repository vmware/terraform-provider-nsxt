/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbNodeUsageSummary
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbNodeUsageSummaryClient interface {

    //
    //
    // @param enforcementPointPathParam enforcement point path (optional)
    // @param includeUsagesParam Whether to include usages (optional)
    // @return com.vmware.nsx_policy.model.AggregateLBNodeUsageSummary
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(enforcementPointPathParam *string, includeUsagesParam *bool) (model.AggregateLBNodeUsageSummary, error)
}
