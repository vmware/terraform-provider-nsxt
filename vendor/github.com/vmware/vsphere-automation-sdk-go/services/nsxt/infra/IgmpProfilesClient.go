/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IgmpProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IgmpProfilesClient interface {

    // Delete Igmp Profile.
    //
    // @param igmpProfileIdParam igmp profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(igmpProfileIdParam string) error

    // Read Igmp Profile.
    //
    // @param igmpProfileIdParam igmp profile id (required)
    // @return com.vmware.nsx_policy.model.PolicyIgmpProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(igmpProfileIdParam string) (model.PolicyIgmpProfile, error)

    // List all igmp profile.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyIgmpProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyIgmpProfileListResult, error)

    // Create a igmp profile with the igmp-profile-id is not already present, otherwise update the igmp profile.
    //
    // @param igmpProfileIdParam igmp profile id (required)
    // @param policyIgmpProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(igmpProfileIdParam string, policyIgmpProfileParam model.PolicyIgmpProfile) error

    // Create or update igmp profile.
    //
    // @param igmpProfileIdParam igmp profile id (required)
    // @param policyIgmpProfileParam (required)
    // @return com.vmware.nsx_policy.model.PolicyIgmpProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(igmpProfileIdParam string, policyIgmpProfileParam model.PolicyIgmpProfile) (model.PolicyIgmpProfile, error)
}
