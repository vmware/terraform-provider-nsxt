/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbPools
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbPoolsClient interface {

    // Delete the LBPool along with all the entities contained by this LBPool.
    //
    // @param lbPoolIdParam LBPool ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(lbPoolIdParam string, forceParam *bool) error

    // Read a LBPool.
    //
    // @param lbPoolIdParam LBPool ID (required)
    // @return com.vmware.nsx_policy.model.LBPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(lbPoolIdParam string) (model.LBPool, error)

    // Paginated list of all LBPools.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LBPoolListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LBPoolListResult, error)

    // If a LBPool with the lb-pool-id is not already present, create a new LBPool. If it already exists, update the LBPool. This is a full replace.
    //
    // @param lbPoolIdParam LBPool ID (required)
    // @param lbPoolParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(lbPoolIdParam string, lbPoolParam model.LBPool) error

    // If a LBPool with the lb-pool-id is not already present, create a new LBPool. If it already exists, update the LBPool. This is a full replace.
    //
    // @param lbPoolIdParam LBPool ID (required)
    // @param lbPoolParam (required)
    // @return com.vmware.nsx_policy.model.LBPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(lbPoolIdParam string, lbPoolParam model.LBPool) (model.LBPool, error)
}
