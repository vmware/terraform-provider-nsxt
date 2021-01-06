/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbTrafficCloneProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbTrafficCloneProfilesClient interface {

    // Delete the ALBTrafficCloneProfile along with all the entities contained by this ALBTrafficCloneProfile.
    //
    // @param albTrafficcloneprofileIdParam ALBTrafficCloneProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albTrafficcloneprofileIdParam string, forceParam *bool) error

    // Read a ALBTrafficCloneProfile.
    //
    // @param albTrafficcloneprofileIdParam ALBTrafficCloneProfile ID (required)
    // @return com.vmware.nsx_policy.model.ALBTrafficCloneProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albTrafficcloneprofileIdParam string) (model.ALBTrafficCloneProfile, error)

    // Paginated list of all ALBTrafficCloneProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBTrafficCloneProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBTrafficCloneProfileApiResponse, error)

    // If a ALBtrafficcloneprofile with the alb-trafficcloneprofile-id is not already present, create a new ALBtrafficcloneprofile. If it already exists, update the ALBtrafficcloneprofile. This is a full replace.
    //
    // @param albTrafficcloneprofileIdParam ALBtrafficcloneprofile ID (required)
    // @param aLBTrafficCloneProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albTrafficcloneprofileIdParam string, aLBTrafficCloneProfileParam model.ALBTrafficCloneProfile) error

    // If a ALBTrafficCloneProfile with the alb-TrafficCloneProfile-id is not already present, create a new ALBTrafficCloneProfile. If it already exists, update the ALBTrafficCloneProfile. This is a full replace.
    //
    // @param albTrafficcloneprofileIdParam ALBTrafficCloneProfile ID (required)
    // @param aLBTrafficCloneProfileParam (required)
    // @return com.vmware.nsx_policy.model.ALBTrafficCloneProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albTrafficcloneprofileIdParam string, aLBTrafficCloneProfileParam model.ALBTrafficCloneProfile) (model.ALBTrafficCloneProfile, error)
}
