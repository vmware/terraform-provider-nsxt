/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpAllocations
 * Used by client-side stubs.
 */

package ip_pools

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpAllocationsClient interface {

    // Releases the IP that was allocated for this allocation request
    //
    // @param ipPoolIdParam (required)
    // @param ipAllocationIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipPoolIdParam string, ipAllocationIdParam string) error

    // Read a previously created allocation
    //
    // @param ipPoolIdParam (required)
    // @param ipAllocationIdParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressAllocation
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipPoolIdParam string, ipAllocationIdParam string) (model.IpAddressAllocation, error)

    // Returns information about which addresses have been allocated from a specified IP address pool in policy.
    //
    // @param ipPoolIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IpAddressAllocationListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(ipPoolIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IpAddressAllocationListResult, error)

    // If allocation of the same ID is found, this is a no-op. If no allocation of the specified ID is found, then a new allocation is created. An allocation cannot be updated once created. When an allocation is requested from an IpAddressPool, the IP could be allocated from any subnet in the pool that has the available capacity. Request to allocate an IP will fail if no subnet was previously created. If specific IP was requested, the status of allocation is reflected in the realized state. If any IP is requested, the IP finally allocated is obtained by polling on the realized state until the allocated IP is returned in the extended attributes.
    //
    // @param ipPoolIdParam (required)
    // @param ipAllocationIdParam (required)
    // @param ipAddressAllocationParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipPoolIdParam string, ipAllocationIdParam string, ipAddressAllocationParam model.IpAddressAllocation) error

    // If allocation of the same ID is found, this is a no-op. If no allocation of the specified ID is found, then a new allocation is created. An allocation cannot be updated once created. When an IP allocation is requested from an IpAddressPool, the IP could be allocated from any subnet in the pool that has the available capacity. Request to allocate an IP will fail if no subnet was previously created. If specific IP was requested, the status of allocation is reflected in the realized state. If any IP is requested, the IP finally allocated is obtained by polling on the realized state until the allocated IP is returned in the extended attributes. An allocation cannot be updated once created.
    //
    // @param ipPoolIdParam (required)
    // @param ipAllocationIdParam (required)
    // @param ipAddressAllocationParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressAllocation
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipPoolIdParam string, ipAllocationIdParam string, ipAddressAllocationParam model.IpAddressAllocation) (model.IpAddressAllocation, error)
}
