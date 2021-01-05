/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbWafCrs
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbWafCrsClient interface {

    // Delete the ALBWafCRS along with all the entities contained by this ALBWafCRS.
    //
    // @param albWafcrsIdParam ALBWafCRS ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albWafcrsIdParam string, forceParam *bool) error

    // Read a ALBWafCRS.
    //
    // @param albWafcrsIdParam ALBWafCRS ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafCRS
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albWafcrsIdParam string) (model.ALBWafCRS, error)

    // Paginated list of all ALBWafCRS for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBWafCRSApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafCRSApiResponse, error)

    // If a ALBwafcrs with the alb-wafcrs-id is not already present, create a new ALBwafcrs. If it already exists, update the ALBwafcrs. This is a full replace.
    //
    // @param albWafcrsIdParam ALBwafcrs ID (required)
    // @param aLBWafCRSParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albWafcrsIdParam string, aLBWafCRSParam model.ALBWafCRS) error

    // If a ALBWafCRS with the alb-WafCRS-id is not already present, create a new ALBWafCRS. If it already exists, update the ALBWafCRS. This is a full replace.
    //
    // @param albWafcrsIdParam ALBWafCRS ID (required)
    // @param aLBWafCRSParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafCRS
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albWafcrsIdParam string, aLBWafCRSParam model.ALBWafCRS) (model.ALBWafCRS, error)
}
