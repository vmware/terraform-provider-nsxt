/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbSecurityPolicies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbSecurityPoliciesClient interface {

    // Delete the ALBSecurityPolicy along with all the entities contained by this ALBSecurityPolicy.
    //
    // @param albSecuritypolicyIdParam ALBSecurityPolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albSecuritypolicyIdParam string, forceParam *bool) error

    // Read a ALBSecurityPolicy.
    //
    // @param albSecuritypolicyIdParam ALBSecurityPolicy ID (required)
    // @return com.vmware.nsx_policy.model.ALBSecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albSecuritypolicyIdParam string) (model.ALBSecurityPolicy, error)

    // Paginated list of all ALBSecurityPolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBSecurityPolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBSecurityPolicyApiResponse, error)

    // If a ALBsecuritypolicy with the alb-securitypolicy-id is not already present, create a new ALBsecuritypolicy. If it already exists, update the ALBsecuritypolicy. This is a full replace.
    //
    // @param albSecuritypolicyIdParam ALBsecuritypolicy ID (required)
    // @param aLBSecurityPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albSecuritypolicyIdParam string, aLBSecurityPolicyParam model.ALBSecurityPolicy) error

    // If a ALBSecurityPolicy with the alb-SecurityPolicy-id is not already present, create a new ALBSecurityPolicy. If it already exists, update the ALBSecurityPolicy. This is a full replace.
    //
    // @param albSecuritypolicyIdParam ALBSecurityPolicy ID (required)
    // @param aLBSecurityPolicyParam (required)
    // @return com.vmware.nsx_policy.model.ALBSecurityPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albSecuritypolicyIdParam string, aLBSecurityPolicyParam model.ALBSecurityPolicy) (model.ALBSecurityPolicy, error)
}
