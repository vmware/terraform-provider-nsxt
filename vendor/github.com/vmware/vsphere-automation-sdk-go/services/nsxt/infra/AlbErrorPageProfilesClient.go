/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbErrorPageProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbErrorPageProfilesClient interface {

    // Delete the ALBErrorPageProfile along with all the entities contained by this ALBErrorPageProfile.
    //
    // @param albErrorpageprofileIdParam ALBErrorPageProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albErrorpageprofileIdParam string, forceParam *bool) error

    // Read a ALBErrorPageProfile.
    //
    // @param albErrorpageprofileIdParam ALBErrorPageProfile ID (required)
    // @return com.vmware.nsx_policy.model.ALBErrorPageProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albErrorpageprofileIdParam string) (model.ALBErrorPageProfile, error)

    // Paginated list of all ALBErrorPageProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBErrorPageProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBErrorPageProfileApiResponse, error)

    // If a ALBerrorpageprofile with the alb-errorpageprofile-id is not already present, create a new ALBerrorpageprofile. If it already exists, update the ALBerrorpageprofile. This is a full replace.
    //
    // @param albErrorpageprofileIdParam ALBerrorpageprofile ID (required)
    // @param aLBErrorPageProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albErrorpageprofileIdParam string, aLBErrorPageProfileParam model.ALBErrorPageProfile) error

    // If a ALBErrorPageProfile with the alb-ErrorPageProfile-id is not already present, create a new ALBErrorPageProfile. If it already exists, update the ALBErrorPageProfile. This is a full replace.
    //
    // @param albErrorpageprofileIdParam ALBErrorPageProfile ID (required)
    // @param aLBErrorPageProfileParam (required)
    // @return com.vmware.nsx_policy.model.ALBErrorPageProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albErrorpageprofileIdParam string, aLBErrorPageProfileParam model.ALBErrorPageProfile) (model.ALBErrorPageProfile, error)
}
