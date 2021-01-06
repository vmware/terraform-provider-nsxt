/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbPoolGroups
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbPoolGroupsClient interface {

    // Delete the ALBPoolGroup along with all the entities contained by this ALBPoolGroup.
    //
    // @param albPoolgroupIdParam ALBPoolGroup ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albPoolgroupIdParam string, forceParam *bool) error

    // Read a ALBPoolGroup.
    //
    // @param albPoolgroupIdParam ALBPoolGroup ID (required)
    // @return com.vmware.nsx_policy.model.ALBPoolGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albPoolgroupIdParam string) (model.ALBPoolGroup, error)

    // Paginated list of all ALBPoolGroup for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBPoolGroupApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPoolGroupApiResponse, error)

    // If a ALBpoolgroup with the alb-poolgroup-id is not already present, create a new ALBpoolgroup. If it already exists, update the ALBpoolgroup. This is a full replace.
    //
    // @param albPoolgroupIdParam ALBpoolgroup ID (required)
    // @param aLBPoolGroupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albPoolgroupIdParam string, aLBPoolGroupParam model.ALBPoolGroup) error

    // If a ALBPoolGroup with the alb-PoolGroup-id is not already present, create a new ALBPoolGroup. If it already exists, update the ALBPoolGroup. This is a full replace.
    //
    // @param albPoolgroupIdParam ALBPoolGroup ID (required)
    // @param aLBPoolGroupParam (required)
    // @return com.vmware.nsx_policy.model.ALBPoolGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albPoolgroupIdParam string, aLBPoolGroupParam model.ALBPoolGroup) (model.ALBPoolGroup, error)
}
