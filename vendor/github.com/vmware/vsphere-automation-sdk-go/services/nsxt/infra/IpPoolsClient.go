/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpPools
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpPoolsClient interface {

    // Delete the IpAddressPool with the given id.
    //
    // @param ipPoolIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipPoolIdParam string) error

    // Read IpAddressPool with given Id.
    //
    // @param ipPoolIdParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipPoolIdParam string) (model.IpAddressPool, error)

    // Paginated list of IpAddressPools.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IpAddressPoolListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IpAddressPoolListResult, error)

    // Creates a new IpAddressPool with specified ID if not already present. If IpAddressPool of given ID is already present, then the instance is updated. This is a full replace.
    //
    // @param ipPoolIdParam (required)
    // @param ipAddressPoolParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipPoolIdParam string, ipAddressPoolParam model.IpAddressPool) error

    // Create a new IpAddressPool with given ID if it does not exist. If IpAddressPool with given ID already exists, it will update existing instance. This is a full replace.
    //
    // @param ipPoolIdParam (required)
    // @param ipAddressPoolParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressPool
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipPoolIdParam string, ipAddressPoolParam model.IpAddressPool) (model.IpAddressPool, error)
}
