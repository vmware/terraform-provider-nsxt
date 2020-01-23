/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GatewayFirewall
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GatewayFirewallClient interface {

    // Get filtered view of gateway rules associated with the Tier-0. The gateay policies are returned in the order of category and precedence.
    //
    // @param tier0IdParam (required)
    // @return com.vmware.nsx_policy.model.GatewayPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string) (model.GatewayPolicyListResult, error)
}
