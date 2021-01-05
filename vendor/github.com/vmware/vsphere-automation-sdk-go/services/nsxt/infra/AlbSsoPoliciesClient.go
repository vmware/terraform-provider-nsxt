/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbSsoPolicies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbSsoPoliciesClient interface {

    // Delete the ALBSSOPolicy along with all the entities contained by this ALBSSOPolicy.
    //
    // @param albSsopolicyIdParam ALBSSOPolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albSsopolicyIdParam string, forceParam *bool) error

    // Read a ALBSSOPolicy.
    //
    // @param albSsopolicyIdParam ALBSSOPolicy ID (required)
    // @return com.vmware.nsx_policy.model.ALBSSOPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albSsopolicyIdParam string) (model.ALBSSOPolicy, error)

    // Paginated list of all ALBSSOPolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBSSOPolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBSSOPolicyApiResponse, error)

    // If a ALBssopolicy with the alb-ssopolicy-id is not already present, create a new ALBssopolicy. If it already exists, update the ALBssopolicy. This is a full replace.
    //
    // @param albSsopolicyIdParam ALBssopolicy ID (required)
    // @param aLBSSOPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albSsopolicyIdParam string, aLBSSOPolicyParam model.ALBSSOPolicy) error

    // If a ALBSSOPolicy with the alb-SSOPolicy-id is not already present, create a new ALBSSOPolicy. If it already exists, update the ALBSSOPolicy. This is a full replace.
    //
    // @param albSsopolicyIdParam ALBSSOPolicy ID (required)
    // @param aLBSSOPolicyParam (required)
    // @return com.vmware.nsx_policy.model.ALBSSOPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albSsopolicyIdParam string, aLBSSOPolicyParam model.ALBSSOPolicy) (model.ALBSSOPolicy, error)
}
