/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbSslProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbSslProfilesClient interface {

    // Delete the ALBSSLProfile along with all the entities contained by this ALBSSLProfile.
    //
    // @param albSslprofileIdParam ALBSSLProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albSslprofileIdParam string, forceParam *bool) error

    // Read a ALBSSLProfile.
    //
    // @param albSslprofileIdParam ALBSSLProfile ID (required)
    // @return com.vmware.nsx_policy.model.ALBSSLProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albSslprofileIdParam string) (model.ALBSSLProfile, error)

    // Paginated list of all ALBSSLProfile for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBSSLProfileApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBSSLProfileApiResponse, error)

    // If a ALBsslprofile with the alb-sslprofile-id is not already present, create a new ALBsslprofile. If it already exists, update the ALBsslprofile. This is a full replace.
    //
    // @param albSslprofileIdParam ALBsslprofile ID (required)
    // @param aLBSSLProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albSslprofileIdParam string, aLBSSLProfileParam model.ALBSSLProfile) error

    // If a ALBSSLProfile with the alb-SSLProfile-id is not already present, create a new ALBSSLProfile. If it already exists, update the ALBSSLProfile. This is a full replace.
    //
    // @param albSslprofileIdParam ALBSSLProfile ID (required)
    // @param aLBSSLProfileParam (required)
    // @return com.vmware.nsx_policy.model.ALBSSLProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albSslprofileIdParam string, aLBSSLProfileParam model.ALBSSLProfile) (model.ALBSSLProfile, error)
}
