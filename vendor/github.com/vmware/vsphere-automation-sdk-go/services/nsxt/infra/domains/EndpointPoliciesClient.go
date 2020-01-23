/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EndpointPolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EndpointPoliciesClient interface {

    // Delete Endpoint policy.
    //
    // @param domainIdParam Domain id (required)
    // @param endpointPolicyIdParam Endpoint policy id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, endpointPolicyIdParam string) error

    // Read Endpoint policy.
    //
    // @param domainIdParam Domain id (required)
    // @param endpointPolicyIdParam Endpoint policy id (required)
    // @return com.vmware.nsx_policy.model.EndpointPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, endpointPolicyIdParam string) (model.EndpointPolicy, error)

    // List all Endpoint policies across all domains ordered by precedence.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.EndpointPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.EndpointPolicyListResult, error)

    // Create or update the Endpoint policy.
    //
    // @param domainIdParam Domain id (required)
    // @param endpointPolicyIdParam Endpoint policy id (required)
    // @param endpointPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, endpointPolicyIdParam string, endpointPolicyParam model.EndpointPolicy) error

    // Create or update the Endpoint policy.
    //
    // @param domainIdParam Domain id (required)
    // @param endpointPolicyIdParam Endpoint policy id (required)
    // @param endpointPolicyParam (required)
    // @return com.vmware.nsx_policy.model.EndpointPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, endpointPolicyIdParam string, endpointPolicyParam model.EndpointPolicy) (model.EndpointPolicy, error)
}
