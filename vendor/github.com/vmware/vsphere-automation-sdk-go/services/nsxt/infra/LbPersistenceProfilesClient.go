/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbPersistenceProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type LbPersistenceProfilesClient interface {

    // Delete the LBPersistenceProfile along with all the entities contained by this LBPersistenceProfile.
    //
    // @param lbPersistenceProfileIdParam LBPersistenceProfile ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(lbPersistenceProfileIdParam string, forceParam *bool) error

    // Read a LBPersistenceProfile.
    //
    // @param lbPersistenceProfileIdParam LBPersistenceProfile ID (required)
    // @return com.vmware.nsx_policy.model.LBPersistenceProfile
    // The return value will contain all the properties defined in model.LBPersistenceProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(lbPersistenceProfileIdParam string) (*data.StructValue, error)

    // Paginated list of all LBPersistenceProfiles for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LBPersistenceProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LBPersistenceProfileListResult, error)

    // If a LBPersistenceProfile with the lb-persistence-profile-id is not already present, create a new LBPersistenceProfile. If it already exists, update the LBPersistenceProfile. This is a full replace.
    //
    // @param lbPersistenceProfileIdParam LBPersistenceProfile ID (required)
    // @param lbPersistenceProfileParam (required)
    // The parameter must contain all the properties defined in model.LBPersistenceProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(lbPersistenceProfileIdParam string, lbPersistenceProfileParam *data.StructValue) error

    // If a LBPersistenceProfile with the lb-persistence-profile-id is not already present, create a new LBPersistenceProfile. If it already exists, update the LBPersistenceProfile. This is a full replace.
    //
    // @param lbPersistenceProfileIdParam LBPersistenceProfile ID (required)
    // @param lbPersistenceProfileParam (required)
    // The parameter must contain all the properties defined in model.LBPersistenceProfile.
    // @return com.vmware.nsx_policy.model.LBPersistenceProfile
    // The return value will contain all the properties defined in model.LBPersistenceProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(lbPersistenceProfileIdParam string, lbPersistenceProfileParam *data.StructValue) (*data.StructValue, error)
}
