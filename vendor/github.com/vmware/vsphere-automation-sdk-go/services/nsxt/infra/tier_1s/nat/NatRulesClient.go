/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: NatRules
 * Used by client-side stubs.
 */

package nat

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type NatRulesClient interface {

    // Delete NAT Rule from Tier-1 denoted by Tier-1 ID, under NAT section denoted by <nat-id>. Under tier-1 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param natIdParam NAT id (required)
    // @param natRuleIdParam Rule ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, natIdParam string, natRuleIdParam string) error

    // Get NAT Rule from Tier-1 denoted by Tier-1 ID, under NAT section denoted by <nat-id>. Under tier-1 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param natIdParam NAT id (required)
    // @param natRuleIdParam Rule ID (required)
    // @return com.vmware.nsx_policy.model.PolicyNatRule
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, natIdParam string, natRuleIdParam string) (model.PolicyNatRule, error)

    // List NAT Rules from Tier-1 denoted by Tier-1 ID, under NAT section denoted by <nat-id>. Under tier-1 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param natIdParam NAT id (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyNatRuleListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, natIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyNatRuleListResult, error)

    // If a NAT Rule is not already present on Tier-1 denoted by Tier-1 ID, under NAT section denoted by <nat-id>, create a new NAT Rule. If it already exists, update the NAT Rule. Under tier-1 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param natIdParam NAT id (required)
    // @param natRuleIdParam Rule ID (required)
    // @param policyNatRuleParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, natIdParam string, natRuleIdParam string, policyNatRuleParam model.PolicyNatRule) error

    // Update NAT Rule on Tier-1 denoted by Tier-1 ID, under NAT section denoted by <nat-id>. Under tier-1 there will be 3 different NATs(sections). (INTERNAL, USER and DEFAULT) For more details related to NAT section please refer to PolicyNAT schema.
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param natIdParam NAT id (required)
    // @param natRuleIdParam Rule ID (required)
    // @param policyNatRuleParam (required)
    // @return com.vmware.nsx_policy.model.PolicyNatRule
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, natIdParam string, natRuleIdParam string, policyNatRuleParam model.PolicyNatRule) (model.PolicyNatRule, error)
}
