/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Groups
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GroupsClient interface {

    // Delete the Group under Tier-0.
    //
    // @param tier0IdParam (required)
    // @param groupIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, groupIdParam string) error

    // Read Tier-0 Group
    //
    // @param tier0IdParam (required)
    // @param groupIdParam (required)
    // @return com.vmware.nsx_policy.model.Group
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, groupIdParam string) (model.Group, error)

    // Paginated list of all Groups for Tier-0.
    //
    // @param tier0IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param memberTypesParam Comma Seperated Member types (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.GroupListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, memberTypesParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GroupListResult, error)

    // If a Group with the group-id is not already present, create a new Group under the tier-0-id. Update if exists. The API valiates that Tier-0 is present before creating the Group.
    //
    // @param tier0IdParam (required)
    // @param groupIdParam (required)
    // @param groupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, groupIdParam string, groupParam model.Group) error

    // If a Group with the group-id is not already present, create a new Group under the tier-0-id. Update if exists. The API valiates that Tier-0 is present before creating the Group.
    //
    // @param tier0IdParam (required)
    // @param groupIdParam (required)
    // @param groupParam (required)
    // @return com.vmware.nsx_policy.model.Group
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, groupIdParam string, groupParam model.Group) (model.Group, error)
}
