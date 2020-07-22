/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Neighbors
 * Used by client-side stubs.
 */

package bgp

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type NeighborsClient interface {

    // Delete BGP neighbor config
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param neighborIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, neighborIdParam string) error

    // Read BGP neighbor config
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param neighborIdParam (required)
    // @return com.vmware.nsx_global_policy.model.BgpNeighborConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, neighborIdParam string) (model.BgpNeighborConfig, error)

    // Paginated list of all BGP neighbor configurations
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.BgpNeighborConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.BgpNeighborConfigListResult, error)

    // If BGP neighbor config with the neighbor-id is not already present, create a new neighbor config. If it already exists, replace the BGP neighbor config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param neighborIdParam (required)
    // @param bgpNeighborConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, neighborIdParam string, bgpNeighborConfigParam model.BgpNeighborConfig) error

    // If BGP neighbor config with the neighbor-id is not already present, create a new neighbor config. If it already exists, replace the BGP neighbor config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param neighborIdParam (required)
    // @param bgpNeighborConfigParam (required)
    // @return com.vmware.nsx_global_policy.model.BgpNeighborConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, neighborIdParam string, bgpNeighborConfigParam model.BgpNeighborConfig) (model.BgpNeighborConfig, error)
}
