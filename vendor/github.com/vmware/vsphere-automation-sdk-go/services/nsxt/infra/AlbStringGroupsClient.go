/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbStringGroups
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbStringGroupsClient interface {

    // Delete the ALBStringGroup along with all the entities contained by this ALBStringGroup.
    //
    // @param albStringgroupIdParam ALBStringGroup ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albStringgroupIdParam string, forceParam *bool) error

    // Read a ALBStringGroup.
    //
    // @param albStringgroupIdParam ALBStringGroup ID (required)
    // @return com.vmware.nsx_policy.model.ALBStringGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albStringgroupIdParam string) (model.ALBStringGroup, error)

    // Paginated list of all ALBStringGroup for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBStringGroupApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBStringGroupApiResponse, error)

    // If a ALBstringgroup with the alb-stringgroup-id is not already present, create a new ALBstringgroup. If it already exists, update the ALBstringgroup. This is a full replace.
    //
    // @param albStringgroupIdParam ALBstringgroup ID (required)
    // @param aLBStringGroupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albStringgroupIdParam string, aLBStringGroupParam model.ALBStringGroup) error

    // If a ALBStringGroup with the alb-StringGroup-id is not already present, create a new ALBStringGroup. If it already exists, update the ALBStringGroup. This is a full replace.
    //
    // @param albStringgroupIdParam ALBStringGroup ID (required)
    // @param aLBStringGroupParam (required)
    // @return com.vmware.nsx_policy.model.ALBStringGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albStringgroupIdParam string, aLBStringGroupParam model.ALBStringGroup) (model.ALBStringGroup, error)
}
