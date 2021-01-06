/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbAuthProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbAuthProfilesClient interface {

    // Delete the ALBAuthProfile along with all the entities contained by this ALBAuthProfile.
    //
    // @param albAuthprofileIdParam ALBAuthProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albAuthprofileIdParam string, forceParam *bool) error

    // Read a ALBAuthProfile.
    //
    // @param albAuthprofileIdParam ALBAuthProfile ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBAuthProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albAuthprofileIdParam string) (model.ALBAuthProfile, error)

    // Paginated list of all ALBAuthProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBAuthProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBAuthProfileApiResponse, error)

    // If a ALBauthprofile with the alb-authprofile-id is not already present, create a new ALBauthprofile. If it already exists, update the ALBauthprofile. This is a full replace.
    //
    // @param albAuthprofileIdParam ALBauthprofile ID (required)
    // @param aLBAuthProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albAuthprofileIdParam string, aLBAuthProfileParam model.ALBAuthProfile) error

    // If a ALBAuthProfile with the alb-AuthProfile-id is not already present, create a new ALBAuthProfile. If it already exists, update the ALBAuthProfile. This is a full replace.
    //
    // @param albAuthprofileIdParam ALBAuthProfile ID (required)
    // @param aLBAuthProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBAuthProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albAuthprofileIdParam string, aLBAuthProfileParam model.ALBAuthProfile) (model.ALBAuthProfile, error)
}
