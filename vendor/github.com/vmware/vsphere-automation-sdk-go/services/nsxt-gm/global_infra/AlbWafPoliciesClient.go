/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbWafPolicies
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbWafPoliciesClient interface {

    // Delete the ALBWafPolicy along with all the entities contained by this ALBWafPolicy.
    //
    // @param albWafpolicyIdParam ALBWafPolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albWafpolicyIdParam string, forceParam *bool) error

    // Read a ALBWafPolicy.
    //
    // @param albWafpolicyIdParam ALBWafPolicy ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albWafpolicyIdParam string) (model.ALBWafPolicy, error)

    // Paginated list of all ALBWafPolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafPolicyApiResponse, error)

    // If a ALBwafpolicy with the alb-wafpolicy-id is not already present, create a new ALBwafpolicy. If it already exists, update the ALBwafpolicy. This is a full replace.
    //
    // @param albWafpolicyIdParam ALBwafpolicy ID (required)
    // @param aLBWafPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albWafpolicyIdParam string, aLBWafPolicyParam model.ALBWafPolicy) error

    // If a ALBWafPolicy with the alb-WafPolicy-id is not already present, create a new ALBWafPolicy. If it already exists, update the ALBWafPolicy. This is a full replace.
    //
    // @param albWafpolicyIdParam ALBWafPolicy ID (required)
    // @param aLBWafPolicyParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albWafpolicyIdParam string, aLBWafPolicyParam model.ALBWafPolicy) (model.ALBWafPolicy, error)
}
