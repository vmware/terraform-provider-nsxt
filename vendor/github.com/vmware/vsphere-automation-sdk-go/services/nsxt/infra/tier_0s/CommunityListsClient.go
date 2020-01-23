/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: CommunityLists
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type CommunityListsClient interface {

    // Delete a BGP community list
    //
    // @param tier0IdParam (required)
    // @param communityListIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, communityListIdParam string) error

    // Read a BGP community list
    //
    // @param tier0IdParam (required)
    // @param communityListIdParam (required)
    // @return com.vmware.nsx_policy.model.CommunityList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, communityListIdParam string) (model.CommunityList, error)

    // Paginated list of all community lists under a tier-0
    //
    // @param tier0IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.CommunityListListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.CommunityListListResult, error)

    // If a community list with the community-list-id is not already present, create a new community list. If it already exists, update the community list for specified attributes.
    //
    // @param tier0IdParam (required)
    // @param communityListIdParam (required)
    // @param communityListParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, communityListIdParam string, communityListParam model.CommunityList) error

    // If a community list with the community-list-id is not already present, create a new community list. If it already exists, replace the community list instance with the new object.
    //
    // @param tier0IdParam (required)
    // @param communityListIdParam (required)
    // @param communityListParam (required)
    // @return com.vmware.nsx_policy.model.CommunityList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, communityListIdParam string, communityListParam model.CommunityList) (model.CommunityList, error)
}
