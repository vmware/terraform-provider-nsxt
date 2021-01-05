/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IntrusionServicePolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IntrusionServicePoliciesClient interface {

    // Delete intrusion detection system security policy.
    //
    // @param domainIdParam Domain ID (required)
    // @param policyIdParam Policy ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, policyIdParam string) error

    // Read intrusion detection system security policy.
    //
    // @param domainIdParam Domain ID (required)
    // @param policyIdParam Policy ID (required)
    // @return com.vmware.nsx_policy.model.IdsSecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, policyIdParam string) (model.IdsSecurityPolicy, error)

    // List intrusion detection system security policies.
    //
    // @param domainIdParam Domain ID (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includeRuleCountParam Include the count of rules in policy (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IdsSecurityPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IdsSecurityPolicyListResult, error)

    // Patch intrusion detection system security policy for a domain.
    //
    // @param domainIdParam Domain ID (required)
    // @param policyIdParam Policy ID (required)
    // @param idsSecurityPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, policyIdParam string, idsSecurityPolicyParam model.IdsSecurityPolicy) error

    // This is used to set a precedence of a security policy w.r.t others.
    //
    // @param domainIdParam (required)
    // @param policyIdParam (required)
    // @param idsSecurityPolicyParam (required)
    // @param anchorPathParam The security policy/rule path if operation is 'insert_after' or 'insert_before' (optional)
    // @param operationParam Operation (optional, default to insert_top)
    // @return com.vmware.nsx_policy.model.IdsSecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Revise(domainIdParam string, policyIdParam string, idsSecurityPolicyParam model.IdsSecurityPolicy, anchorPathParam *string, operationParam *string) (model.IdsSecurityPolicy, error)

    // Update intrusion detection system security policy for a domain.
    //
    // @param domainIdParam Domain ID (required)
    // @param policyIdParam Policy ID (required)
    // @param idsSecurityPolicyParam (required)
    // @return com.vmware.nsx_policy.model.IdsSecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, policyIdParam string, idsSecurityPolicyParam model.IdsSecurityPolicy) (model.IdsSecurityPolicy, error)
}
