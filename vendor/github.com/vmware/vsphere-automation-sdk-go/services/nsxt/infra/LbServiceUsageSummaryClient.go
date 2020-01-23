/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbServiceUsageSummary
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbServiceUsageSummaryClient interface {

    // API is used to retrieve the load balancer usage summary for all load balancer services. If the parameter ?include_usages=true exists, the property service-usages is included in the response. By default, service-usages is not included in the response.
    //
    // @param includeUsagesParam Whether to include usages (optional)
    // @return com.vmware.nsx_policy.model.LBServiceUsageSummary
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(includeUsagesParam *bool) (model.LBServiceUsageSummary, error)
}
