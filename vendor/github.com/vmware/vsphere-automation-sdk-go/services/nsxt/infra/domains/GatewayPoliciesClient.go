/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GatewayPolicies
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GatewayPoliciesClient interface {

    // Delete GatewayPolicy
    //
    // @param domainIdParam (required)
    // @param gatewayPolicyIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, gatewayPolicyIdParam string) error

    // Read gateway policy for a domain.
    //
    // @param domainIdParam (required)
    // @param gatewayPolicyIdParam (required)
    // @return com.vmware.nsx_policy.model.GatewayPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, gatewayPolicyIdParam string) (model.GatewayPolicy, error)

    // List all gateway policies for specified Domain.
    //
    // @param domainIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.GatewayPolicyListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GatewayPolicyListResult, error)

    // Update the gateway policy for a domain. This is a full replace. All the rules are replaced.
    //
    // @param domainIdParam (required)
    // @param gatewayPolicyIdParam (required)
    // @param gatewayPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) error

    // This is used to set a precedence of a gateway policy w.r.t others.
    //
    // @param domainIdParam (required)
    // @param gatewayPolicyIdParam (required)
    // @param gatewayPolicyParam (required)
    // @param anchorPathParam The security policy/rule path if operation is 'insert_after' or 'insert_before' (optional)
    // @param operationParam Operation (optional, default to insert_top)
    // @return com.vmware.nsx_policy.model.GatewayPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Revise(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy, anchorPathParam *string, operationParam *string) (model.GatewayPolicy, error)

    // Update the gateway policy for a domain. This is a full replace. All the rules are replaced.
    //
    // @param domainIdParam (required)
    // @param gatewayPolicyIdParam (required)
    // @param gatewayPolicyParam (required)
    // @return com.vmware.nsx_policy.model.GatewayPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, gatewayPolicyIdParam string, gatewayPolicyParam model.GatewayPolicy) (model.GatewayPolicy, error)
}
