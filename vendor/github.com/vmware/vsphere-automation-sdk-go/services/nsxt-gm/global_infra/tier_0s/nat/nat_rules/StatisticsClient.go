/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Statistics
 * Used by client-side stubs.
 */

package nat_rules

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type StatisticsClient interface {

    // Get NAT Rule Statistics from Tier-0 denoted by Tier-0 ID, under NAT section denoted by <nat-id>. Under tier-0 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @param natIdParam NAT id (required)
    // @param natRuleIdParam Rule ID (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_global_policy.model.PolicyNatRuleStatisticsListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, natIdParam string, natRuleIdParam string, enforcementPointPathParam *string) (model.PolicyNatRuleStatisticsListResult, error)
}
