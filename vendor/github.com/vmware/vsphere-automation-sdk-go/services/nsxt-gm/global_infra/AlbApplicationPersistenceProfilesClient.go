/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbApplicationPersistenceProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbApplicationPersistenceProfilesClient interface {

    // Delete the ALBApplicationPersistenceProfile along with all the entities contained by this ALBApplicationPersistenceProfile.
    //
    // @param albApplicationpersistenceprofileIdParam ALBApplicationPersistenceProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albApplicationpersistenceprofileIdParam string, forceParam *bool) error

    // Read a ALBApplicationPersistenceProfile.
    //
    // @param albApplicationpersistenceprofileIdParam ALBApplicationPersistenceProfile ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBApplicationPersistenceProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albApplicationpersistenceprofileIdParam string) (model.ALBApplicationPersistenceProfile, error)

    // Paginated list of all ALBApplicationPersistenceProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBApplicationPersistenceProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBApplicationPersistenceProfileApiResponse, error)

    // If a ALBapplicationpersistenceprofile with the alb-applicationpersistenceprofile-id is not already present, create a new ALBapplicationpersistenceprofile. If it already exists, update the ALBapplicationpersistenceprofile. This is a full replace.
    //
    // @param albApplicationpersistenceprofileIdParam ALBapplicationpersistenceprofile ID (required)
    // @param aLBApplicationPersistenceProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albApplicationpersistenceprofileIdParam string, aLBApplicationPersistenceProfileParam model.ALBApplicationPersistenceProfile) error

    // If a ALBApplicationPersistenceProfile with the alb-ApplicationPersistenceProfile-id is not already present, create a new ALBApplicationPersistenceProfile. If it already exists, update the ALBApplicationPersistenceProfile. This is a full replace.
    //
    // @param albApplicationpersistenceprofileIdParam ALBApplicationPersistenceProfile ID (required)
    // @param aLBApplicationPersistenceProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBApplicationPersistenceProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albApplicationpersistenceprofileIdParam string, aLBApplicationPersistenceProfileParam model.ALBApplicationPersistenceProfile) (model.ALBApplicationPersistenceProfile, error)
}
