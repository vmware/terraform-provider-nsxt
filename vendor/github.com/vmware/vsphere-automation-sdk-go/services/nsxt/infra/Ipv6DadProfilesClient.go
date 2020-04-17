/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Ipv6DadProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type Ipv6DadProfilesClient interface {

    // Delete IPv6 DAD profile
    //
    // @param dadProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(dadProfileIdParam string) error

    // Read IPv6 DAD profile
    //
    // @param dadProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.Ipv6DadProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(dadProfileIdParam string) (model.Ipv6DadProfile, error)

    // Paginated list of all IPv6 DAD profile instances
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.Ipv6DadProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Ipv6DadProfileListResult, error)

    // If profile with the dad-profile-id is not already present, create a new IPv6 DAD profile instance. If it already exists, update the IPv6 DAD profile instance with specified attributes.
    //
    // @param dadProfileIdParam (required)
    // @param ipv6DadProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(dadProfileIdParam string, ipv6DadProfileParam model.Ipv6DadProfile) error

    // If profile with the dad-profile-id is not already present, create a new IPv6 DAD profile instance. If it already exists, replace the IPv6 DAD profile instance with this object.
    //
    // @param dadProfileIdParam (required)
    // @param ipv6DadProfileParam (required)
    // @return com.vmware.nsx_policy.model.Ipv6DadProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(dadProfileIdParam string, ipv6DadProfileParam model.Ipv6DadProfile) (model.Ipv6DadProfile, error)
}
