/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbServerAutoScalePolicies
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbServerAutoScalePoliciesClient interface {

    // Delete the ALBServerAutoScalePolicy along with all the entities contained by this ALBServerAutoScalePolicy.
    //
    // @param albServerautoscalepolicyIdParam ALBServerAutoScalePolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albServerautoscalepolicyIdParam string, forceParam *bool) error

    // Read a ALBServerAutoScalePolicy.
    //
    // @param albServerautoscalepolicyIdParam ALBServerAutoScalePolicy ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBServerAutoScalePolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albServerautoscalepolicyIdParam string) (model.ALBServerAutoScalePolicy, error)

    // Paginated list of all ALBServerAutoScalePolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBServerAutoScalePolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBServerAutoScalePolicyApiResponse, error)

    // If a ALBserverautoscalepolicy with the alb-serverautoscalepolicy-id is not already present, create a new ALBserverautoscalepolicy. If it already exists, update the ALBserverautoscalepolicy. This is a full replace.
    //
    // @param albServerautoscalepolicyIdParam ALBserverautoscalepolicy ID (required)
    // @param aLBServerAutoScalePolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albServerautoscalepolicyIdParam string, aLBServerAutoScalePolicyParam model.ALBServerAutoScalePolicy) error

    // If a ALBServerAutoScalePolicy with the alb-ServerAutoScalePolicy-id is not already present, create a new ALBServerAutoScalePolicy. If it already exists, update the ALBServerAutoScalePolicy. This is a full replace.
    //
    // @param albServerautoscalepolicyIdParam ALBServerAutoScalePolicy ID (required)
    // @param aLBServerAutoScalePolicyParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBServerAutoScalePolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albServerautoscalepolicyIdParam string, aLBServerAutoScalePolicyParam model.ALBServerAutoScalePolicy) (model.ALBServerAutoScalePolicy, error)
}
