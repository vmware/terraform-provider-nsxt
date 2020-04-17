/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Groups
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GroupsClient interface {

    // Delete Group
    //
    // @param domainIdParam Domain ID (required)
    // @param groupIdParam Group ID (required)
    // @param failIfSubtreeExistsParam Do not delete if the group subtree has any entities (optional, default to false)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, groupIdParam string, failIfSubtreeExistsParam *bool, forceParam *bool) error

    // Read group
    //
    // @param domainIdParam Domain ID (required)
    // @param groupIdParam Group ID (required)
    // @return com.vmware.nsx_policy.model.Group
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, groupIdParam string) (model.Group, error)

    // List Groups for a domain. Groups can be filtered using member_types query parameter, which returns the groups that contains the specified member types. Multiple member types can be provided as comma separated values. The API also return groups having member type that are subset of provided member_types.
    //
    // @param domainIdParam Domain ID (required)
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
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, memberTypesParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GroupListResult, error)

    // If a group with the group-id is not already present, create a new group. If it already exists, patch the group.
    //
    // @param domainIdParam Domain ID (required)
    // @param groupIdParam Group ID (required)
    // @param groupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, groupIdParam string, groupParam model.Group) error

    // If a group with the group-id is not already present, create a new group. If it already exists, update the group. Avoid creating groups with multiple MACAddressExpression and IPAddressExpression. In future releases, group will be restricted to contain a single MACAddressExpression and IPAddressExpression along with other expressions. To group IPAddresses or MACAddresses, use nested groups instead of multiple IPAddressExpressions/MACAddressExpression.
    //
    // @param domainIdParam Domain ID (required)
    // @param groupIdParam Group ID (required)
    // @param groupParam (required)
    // @return com.vmware.nsx_policy.model.Group
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, groupIdParam string, groupParam model.Group) (model.Group, error)
}
