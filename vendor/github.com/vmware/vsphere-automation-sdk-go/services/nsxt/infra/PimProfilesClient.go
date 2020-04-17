/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: PimProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type PimProfilesClient interface {

    // Delete Pim Profile.
    //
    // @param pimProfileIdParam pim profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(pimProfileIdParam string) error

    // Read Pim Profile.
    //
    // @param pimProfileIdParam pim profile id (required)
    // @return com.vmware.nsx_policy.model.PolicyPimProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(pimProfileIdParam string) (model.PolicyPimProfile, error)

    // List all pim profile.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyPimProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyPimProfileListResult, error)

    // Create a pim profile with the pim-profile-id is not already present, otherwise update the pim profile.
    //
    // @param pimProfileIdParam pim profile id (required)
    // @param policyPimProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(pimProfileIdParam string, policyPimProfileParam model.PolicyPimProfile) error

    // Create or update pim profile.
    //
    // @param pimProfileIdParam pim profile id (required)
    // @param policyPimProfileParam (required)
    // @return com.vmware.nsx_policy.model.PolicyPimProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(pimProfileIdParam string, policyPimProfileParam model.PolicyPimProfile) (model.PolicyPimProfile, error)
}
