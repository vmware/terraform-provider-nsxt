/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpSubnets
 * Used by client-side stubs.
 */

package ip_pools

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type IpSubnetsClient interface {

    // Delete the IpAddressPoolSubnet with the given id.
    //
    // @param ipPoolIdParam (required)
    // @param ipSubnetIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipPoolIdParam string, ipSubnetIdParam string) error

    // Read IpAddressPoolSubnet with given Id.
    //
    // @param ipPoolIdParam (required)
    // @param ipSubnetIdParam (required)
    // @return com.vmware.nsx_policy.model.IpAddressPoolSubnet
    // The return value will contain all the properties defined in model.IpAddressPoolSubnet.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipPoolIdParam string, ipSubnetIdParam string) (*data.StructValue, error)

    // Paginated list of IpAddressPoolSubnets.
    //
    // @param ipPoolIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IpAddressPoolSubnetListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(ipPoolIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IpAddressPoolSubnetListResult, error)

    // Creates a new IpAddressPoolSubnet with the specified ID if it does not already exist. If a IpAddressPoolSubnet of the given ID already exists, IpAddressPoolSubnet will be updated. This is a full replace.
    //
    // @param ipPoolIdParam (required)
    // @param ipSubnetIdParam (required)
    // @param ipAddressPoolSubnetParam (required)
    // The parameter must contain all the properties defined in model.IpAddressPoolSubnet.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipPoolIdParam string, ipSubnetIdParam string, ipAddressPoolSubnetParam *data.StructValue) error

    // Creates a new IpAddressPoolSubnet with the specified ID if it does not already exist. If a IpAddressPoolSubnet of the given ID already exists, IpAddressPoolSubnet will be updated. This is a full replace.
    //
    // @param ipPoolIdParam (required)
    // @param ipSubnetIdParam (required)
    // @param ipAddressPoolSubnetParam (required)
    // The parameter must contain all the properties defined in model.IpAddressPoolSubnet.
    // @return com.vmware.nsx_policy.model.IpAddressPoolSubnet
    // The return value will contain all the properties defined in model.IpAddressPoolSubnet.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipPoolIdParam string, ipSubnetIdParam string, ipAddressPoolSubnetParam *data.StructValue) (*data.StructValue, error)
}
