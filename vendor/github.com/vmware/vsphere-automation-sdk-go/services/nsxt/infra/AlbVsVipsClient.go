/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbVsVips
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbVsVipsClient interface {

    // Delete the ALBVsVip along with all the entities contained by this ALBVsVip.
    //
    // @param albVsvipIdParam ALBVsVip ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albVsvipIdParam string, forceParam *bool) error

    // Read a ALBVsVip.
    //
    // @param albVsvipIdParam ALBVsVip ID (required)
    // @return com.vmware.nsx_policy.model.ALBVsVip
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albVsvipIdParam string) (model.ALBVsVip, error)

    // Paginated list of all ALBVsVip for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBVsVipApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBVsVipApiResponse, error)

    // If a ALBvsvip with the alb-vsvip-id is not already present, create a new ALBvsvip. If it already exists, update the ALBvsvip. This is a full replace.
    //
    // @param albVsvipIdParam ALBvsvip ID (required)
    // @param aLBVsVipParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albVsvipIdParam string, aLBVsVipParam model.ALBVsVip) error

    // If a ALBVsVip with the alb-VsVip-id is not already present, create a new ALBVsVip. If it already exists, update the ALBVsVip. This is a full replace.
    //
    // @param albVsvipIdParam ALBVsVip ID (required)
    // @param aLBVsVipParam (required)
    // @return com.vmware.nsx_policy.model.ALBVsVip
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albVsvipIdParam string, aLBVsVipParam model.ALBVsVip) (model.ALBVsVip, error)
}
