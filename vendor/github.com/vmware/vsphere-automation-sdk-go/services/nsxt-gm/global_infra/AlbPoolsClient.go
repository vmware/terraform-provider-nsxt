/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbPools
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbPoolsClient interface {

    // Delete the ALBPool along with all the entities contained by this ALBPool.
    //
    // @param albPoolIdParam ALBPool ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albPoolIdParam string, forceParam *bool) error

    // Read a ALBPool.
    //
    // @param albPoolIdParam ALBPool ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albPoolIdParam string) (model.ALBPool, error)

    // Paginated list of all ALBPool for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBPoolApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPoolApiResponse, error)

    // If a ALBpool with the alb-pool-id is not already present, create a new ALBpool. If it already exists, update the ALBpool. This is a full replace.
    //
    // @param albPoolIdParam ALBpool ID (required)
    // @param aLBPoolParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albPoolIdParam string, aLBPoolParam model.ALBPool) error

    // If a ALBPool with the alb-Pool-id is not already present, create a new ALBPool. If it already exists, update the ALBPool. This is a full replace.
    //
    // @param albPoolIdParam ALBPool ID (required)
    // @param aLBPoolParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albPoolIdParam string, aLBPoolParam model.ALBPool) (model.ALBPool, error)
}
