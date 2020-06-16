/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: PrefixLists
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type PrefixListsClient interface {

    // Delete a prefix list
    //
    // @param tier0IdParam (required)
    // @param prefixListIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, prefixListIdParam string) error

    // Read a prefix list
    //
    // @param tier0IdParam (required)
    // @param prefixListIdParam (required)
    // @return com.vmware.nsx_global_policy.model.PrefixList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, prefixListIdParam string) (model.PrefixList, error)

    // Paginated list of all prefix lists
    //
    // @param tier0IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.PrefixListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PrefixListResult, error)

    // If prefix list for prefix-list-id is not already present, create a prefix list. If it already exists, update prefix list for prefix-list-id.
    //
    // @param tier0IdParam (required)
    // @param prefixListIdParam (required)
    // @param prefixListParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, prefixListIdParam string, prefixListParam model.PrefixList) error

    // If prefix list for prefix-list-id is not already present, create a prefix list. If it already exists, replace the prefix list for prefix-list-id.
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @param prefixListIdParam Prefix List ID (required)
    // @param prefixListParam (required)
    // @return com.vmware.nsx_global_policy.model.PrefixList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, prefixListIdParam string, prefixListParam model.PrefixList) (model.PrefixList, error)
}
