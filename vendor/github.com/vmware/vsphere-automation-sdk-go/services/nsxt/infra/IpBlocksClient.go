/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpBlocks
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpBlocksClient interface {

    // Delete the IpAddressBlock with the given id.
    //
    // @param ipBlockIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipBlockIdParam string) error

    // Read IpAddressBlock with given Id.
    //
    // @param ipBlockIdParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressBlock
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipBlockIdParam string) (model.IpAddressBlock, error)

    // Paginated list of IpAddressBlocks.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IpAddressBlockListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IpAddressBlockListResult, error)

    // Creates a new IpAddressBlock with specified ID if not already present. If IpAddressBlock of given ID is already present, then the instance is updated with specified attributes.
    //
    // @param ipBlockIdParam (required)
    // @param ipAddressBlockParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipBlockIdParam string, ipAddressBlockParam model.IpAddressBlock) error

    // Create a new IpAddressBlock with given ID if it does not exist. If IpAddressBlock with given ID already exists, it will update existing instance. This is a full replace.
    //
    // @param ipBlockIdParam (required)
    // @param ipAddressBlockParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressBlock
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipBlockIdParam string, ipAddressBlockParam model.IpAddressBlock) (model.IpAddressBlock, error)
}
