/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbClientSslProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbClientSslProfilesClient interface {

    // Delete the LBClientSslProfile along with all the entities contained by this LBClientSslProfile.
    //
    // @param lbClientSslProfileIdParam LBClientSslProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(lbClientSslProfileIdParam string, forceParam *bool) error

    // Read a LBClientSslProfile.
    //
    // @param lbClientSslProfileIdParam LBClientSslProfile ID (required)
    // @return com.vmware.nsx_policy.model.LBClientSslProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(lbClientSslProfileIdParam string) (model.LBClientSslProfile, error)

    // Paginated list of all LBClientSslProfiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LBClientSslProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LBClientSslProfileListResult, error)

    // If a LBClientSslProfile with the lb-client-ssl-profile-id is not already present, create a new LBClientSslProfile. If it already exists, update the LBClientSslProfile. This is a full replace.
    //
    // @param lbClientSslProfileIdParam LBClientSslProfile ID (required)
    // @param lbClientSslProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(lbClientSslProfileIdParam string, lbClientSslProfileParam model.LBClientSslProfile) error

    // If a LBClientSslProfile with the lb-client-ssl-profile-id is not already present, create a new LBClientSslProfile. If it already exists, update the LBClientSslProfile. This is a full replace.
    //
    // @param lbClientSslProfileIdParam LBClientSslProfile ID (required)
    // @param lbClientSslProfileParam (required)
    // @return com.vmware.nsx_policy.model.LBClientSslProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(lbClientSslProfileIdParam string, lbClientSslProfileParam model.LBClientSslProfile) (model.LBClientSslProfile, error)
}
