/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: RouteMaps
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type RouteMapsClient interface {

    // Delete a route map
    //
    // @param tier0IdParam (required)
    // @param routeMapIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, routeMapIdParam string) error

    // Read a route map
    //
    // @param tier0IdParam (required)
    // @param routeMapIdParam (required)
    // @return com.vmware.nsx_global_policy.model.Tier0RouteMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, routeMapIdParam string) (model.Tier0RouteMap, error)

    // Paginated list of all route maps under a tier-0
    //
    // @param tier0IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.Tier0RouteMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Tier0RouteMapListResult, error)

    // If a route map with the route-map-id is not already present, create a new route map. If it already exists, update the route map for specified attributes.
    //
    // @param tier0IdParam (required)
    // @param routeMapIdParam (required)
    // @param tier0RouteMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, routeMapIdParam string, tier0RouteMapParam model.Tier0RouteMap) error

    // If a route map with the route-map-id is not already present, create a new route map. If it already exists, replace the route map instance with the new object.
    //
    // @param tier0IdParam (required)
    // @param routeMapIdParam (required)
    // @param tier0RouteMapParam (required)
    // @return com.vmware.nsx_global_policy.model.Tier0RouteMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, routeMapIdParam string, tier0RouteMapParam model.Tier0RouteMap) (model.Tier0RouteMap, error)
}
