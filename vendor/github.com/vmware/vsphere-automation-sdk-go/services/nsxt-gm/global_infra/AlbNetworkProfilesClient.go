/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbNetworkProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbNetworkProfilesClient interface {

    // Delete the ALBNetworkProfile along with all the entities contained by this ALBNetworkProfile.
    //
    // @param albNetworkprofileIdParam ALBNetworkProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albNetworkprofileIdParam string, forceParam *bool) error

    // Read a ALBNetworkProfile.
    //
    // @param albNetworkprofileIdParam ALBNetworkProfile ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBNetworkProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albNetworkprofileIdParam string) (model.ALBNetworkProfile, error)

    // Paginated list of all ALBNetworkProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBNetworkProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBNetworkProfileApiResponse, error)

    // If a ALBnetworkprofile with the alb-networkprofile-id is not already present, create a new ALBnetworkprofile. If it already exists, update the ALBnetworkprofile. This is a full replace.
    //
    // @param albNetworkprofileIdParam ALBnetworkprofile ID (required)
    // @param aLBNetworkProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albNetworkprofileIdParam string, aLBNetworkProfileParam model.ALBNetworkProfile) error

    // If a ALBNetworkProfile with the alb-NetworkProfile-id is not already present, create a new ALBNetworkProfile. If it already exists, update the ALBNetworkProfile. This is a full replace.
    //
    // @param albNetworkprofileIdParam ALBNetworkProfile ID (required)
    // @param aLBNetworkProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBNetworkProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albNetworkprofileIdParam string, aLBNetworkProfileParam model.ALBNetworkProfile) (model.ALBNetworkProfile, error)
}
