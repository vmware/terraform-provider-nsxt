/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: BfdProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type BfdProfilesClient interface {

    // Delete BFD Config and all the entities contained by this BfdProfile.
    //
    // @param bfdProfileIdParam BfdProfile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(bfdProfileIdParam string) error

    // Read a BfdProfile.
    //
    // @param bfdProfileIdParam BfdProfile ID (required)
    // @return com.vmware.nsx_policy.model.BfdProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(bfdProfileIdParam string) (model.BfdProfile, error)

    // Paginated list of all BfdProfiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.BfdProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.BfdProfileListResult, error)

    // If a BfdProfile with the bfd-profile-id is not already present, create a new BfdProfile. If it already exists, update the BfdProfile. This operation will fully replace the object.
    //
    // @param bfdProfileIdParam BfdProfile ID (required)
    // @param bfdProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(bfdProfileIdParam string, bfdProfileParam model.BfdProfile) error

    // If a BfdProfile with the bfd-profile-id is not already present, create a new BfdProfile. If it already exists, update the BfdProfile. This operation will fully replace the object.
    //
    // @param bfdProfileIdParam BfdProfile ID (required)
    // @param bfdProfileParam (required)
    // @return com.vmware.nsx_policy.model.BfdProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(bfdProfileIdParam string, bfdProfileParam model.BfdProfile) (model.BfdProfile, error)
}
