/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Ipv6NdraProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type Ipv6NdraProfilesClient interface {

    // Delete IPv6 NDRA profile
    //
    // @param ndraProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ndraProfileIdParam string) error

    // Read IPv6 NDRA profile
    //
    // @param ndraProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.Ipv6NdraProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ndraProfileIdParam string) (model.Ipv6NdraProfile, error)

    // Paginated list of all IPv6 NDRA profile instances
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.Ipv6NdraProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Ipv6NdraProfileListResult, error)

    // If profile with the ndra-profile-id is not already present, create a new IPv6 NDRA profile instance. If it already exists, update the IPv6 NDRA profile instance with specified attributes.
    //
    // @param ndraProfileIdParam (required)
    // @param ipv6NdraProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ndraProfileIdParam string, ipv6NdraProfileParam model.Ipv6NdraProfile) error

    // If profile with the ndra-profile-id is not already present, create a new IPv6 NDRA profile instance. If it already exists, replace the IPv6 NDRA profile instance with this object.
    //
    // @param ndraProfileIdParam (required)
    // @param ipv6NdraProfileParam (required)
    // @return com.vmware.nsx_policy.model.Ipv6NdraProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ndraProfileIdParam string, ipv6NdraProfileParam model.Ipv6NdraProfile) (model.Ipv6NdraProfile, error)
}
