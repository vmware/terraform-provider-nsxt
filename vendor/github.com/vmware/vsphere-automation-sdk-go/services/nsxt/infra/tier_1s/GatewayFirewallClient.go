/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GatewayFirewall
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GatewayFirewallClient interface {

    // Get filtered view of Gateway Firewall rules associated with the Tier-1. The gateway policies are returned in the order of category and sequence number.
    //
    // @param tier1IdParam (required)
    // @return com.vmware.nsx_policy.model.GatewayPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string) (model.GatewayPolicyListResult, error)
}
