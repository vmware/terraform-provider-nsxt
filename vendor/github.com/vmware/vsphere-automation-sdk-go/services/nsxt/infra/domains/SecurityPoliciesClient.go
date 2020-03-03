/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SecurityPolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SecurityPoliciesClient interface {

    // Deletes the security policy along with all the rules
    //
    // @param domainIdParam (required)
    // @param securityPolicyIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, securityPolicyIdParam string) error

    // Read security policy for a domain.
    //
    // @param domainIdParam (required)
    // @param securityPolicyIdParam (required)
    // @return com.vmware.nsx_policy.model.SecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, securityPolicyIdParam string) (model.SecurityPolicy, error)

    // List all security policies for a domain.
    //
    // @param domainIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includeRuleCountParam Include the count of rules in policy (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.SecurityPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includeRuleCountParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SecurityPolicyListResult, error)

    // Patch the security policy for a domain. If a security policy for the given security-policy-id is not present, the object will get created and if it is present it will be updated. This is a full replace
    //
    // @param domainIdParam (required)
    // @param securityPolicyIdParam (required)
    // @param securityPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, securityPolicyIdParam string, securityPolicyParam model.SecurityPolicy) error

    // This is used to set a precedence of a security policy w.r.t others.
    //
    // @param domainIdParam (required)
    // @param securityPolicyIdParam (required)
    // @param securityPolicyParam (required)
    // @param anchorPathParam The security policy/rule path if operation is 'insert_after' or 'insert_before' (optional)
    // @param operationParam Operation (optional, default to insert_top)
    // @return com.vmware.nsx_policy.model.SecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Revise(domainIdParam string, securityPolicyIdParam string, securityPolicyParam model.SecurityPolicy, anchorPathParam *string, operationParam *string) (model.SecurityPolicy, error)

    // Create or Update the security policy for a domain. This is a full replace. All the rules are replaced.
    //
    // @param domainIdParam (required)
    // @param securityPolicyIdParam (required)
    // @param securityPolicyParam (required)
    // @return com.vmware.nsx_policy.model.SecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, securityPolicyIdParam string, securityPolicyParam model.SecurityPolicy) (model.SecurityPolicy, error)
}
