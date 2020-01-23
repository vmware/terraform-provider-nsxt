/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: RedirectionPolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type RedirectionPoliciesClient interface {

    // Delete redirection policy.
    //
    // @param domainIdParam Domain id (required)
    // @param redirectionPolicyIdParam Redirection map id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, redirectionPolicyIdParam string) error

    // Read redirection policy.
    //
    // @param domainIdParam Domain id (required)
    // @param redirectionPolicyIdParam Redirection map id (required)
    // @return com.vmware.nsx_policy.model.RedirectionPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, redirectionPolicyIdParam string) (model.RedirectionPolicy, error)

    // List all redirection policys across all domains ordered by precedence.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.RedirectionPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.RedirectionPolicyListResult, error)

    // List redirection policys for a domain
    //
    // @param domainIdParam Domain id (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.RedirectionPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List0(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.RedirectionPolicyListResult, error)

    // Create or update the redirection policy.
    //
    // @param domainIdParam Domain id (required)
    // @param redirectionPolicyIdParam Redirection map id (required)
    // @param redirectionPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, redirectionPolicyIdParam string, redirectionPolicyParam model.RedirectionPolicy) error

    // Create or update the redirection policy.
    //
    // @param domainIdParam Domain id (required)
    // @param redirectionPolicyIdParam Redirection map id (required)
    // @param redirectionPolicyParam (required)
    // @return com.vmware.nsx_policy.model.RedirectionPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, redirectionPolicyIdParam string, redirectionPolicyParam model.RedirectionPolicy) (model.RedirectionPolicy, error)
}
