/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbWafProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbWafProfilesClient interface {

    // Delete the ALBWafProfile along with all the entities contained by this ALBWafProfile.
    //
    // @param albWafprofileIdParam ALBWafProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albWafprofileIdParam string, forceParam *bool) error

    // Read a ALBWafProfile.
    //
    // @param albWafprofileIdParam ALBWafProfile ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albWafprofileIdParam string) (model.ALBWafProfile, error)

    // Paginated list of all ALBWafProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBWafProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafProfileApiResponse, error)

    // If a ALBwafprofile with the alb-wafprofile-id is not already present, create a new ALBwafprofile. If it already exists, update the ALBwafprofile. This is a full replace.
    //
    // @param albWafprofileIdParam ALBwafprofile ID (required)
    // @param aLBWafProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albWafprofileIdParam string, aLBWafProfileParam model.ALBWafProfile) error

    // If a ALBWafProfile with the alb-WafProfile-id is not already present, create a new ALBWafProfile. If it already exists, update the ALBWafProfile. This is a full replace.
    //
    // @param albWafprofileIdParam ALBWafProfile ID (required)
    // @param aLBWafProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albWafprofileIdParam string, aLBWafProfileParam model.ALBWafProfile) (model.ALBWafProfile, error)
}
