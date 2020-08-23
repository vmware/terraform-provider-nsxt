/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ForwardingPolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type ForwardingPoliciesClient interface {

    // Read forwarding policy.
    //
    // @param domainIdParam Domain id (required)
    // @param forwardingPolicyIdParam Forwarding map id (required)
    // @return com.vmware.nsx_global_policy.model.ForwardingPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, forwardingPolicyIdParam string) (model.ForwardingPolicy, error)

    // List all forwarding policies for the given domain ordered by precedence.
    //
    // @param domainIdParam Domain id (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includeRuleCountParam Include the count of rules in policy (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ForwardingPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ForwardingPolicyListResult, error)
}
