/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: StaticRoutes
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type StaticRoutesClient interface {

    // Delete Tier-1 static routes
    //
    // @param tier1IdParam (required)
    // @param routeIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, routeIdParam string) error

    // Read Tier-1 static routes
    //
    // @param tier1IdParam (required)
    // @param routeIdParam (required)
    // @return com.vmware.nsx_policy.model.StaticRoutes
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, routeIdParam string) (model.StaticRoutes, error)

    // Paginated list of all Tier-1 Static Routes
    //
    // @param tier1IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.StaticRoutesListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.StaticRoutesListResult, error)

    // If static routes for route-id are not already present, create static routes. If it already exists, update static routes for route-id.
    //
    // @param tier1IdParam (required)
    // @param routeIdParam (required)
    // @param staticRoutesParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, routeIdParam string, staticRoutesParam model.StaticRoutes) error

    // If static routes for route-id are not already present, create static routes. If it already exists, replace the static routes for route-id.
    //
    // @param tier1IdParam (required)
    // @param routeIdParam (required)
    // @param staticRoutesParam (required)
    // @return com.vmware.nsx_policy.model.StaticRoutes
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, routeIdParam string, staticRoutesParam model.StaticRoutes) (model.StaticRoutes, error)
}
