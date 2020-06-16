/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FloodProtectionProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type FloodProtectionProfilesClient interface {

    // API will delete Flood Protection Profile
    //
    // @param floodProtectionProfileIdParam Flood Protection Profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(floodProtectionProfileIdParam string) error

    // API will get Flood Protection Profile
    //
    // @param floodProtectionProfileIdParam Flood Protection Profile ID (required)
    // @return com.vmware.nsx_global_policy.model.FloodProtectionProfile
    // The return value will contain all the properties defined in model.FloodProtectionProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(floodProtectionProfileIdParam string) (*data.StructValue, error)

    // API will list all Flood Protection Profiles
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.FloodProtectionProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.FloodProtectionProfileListResult, error)

    // API will create/update Flood Protection Profile
    //
    // @param floodProtectionProfileIdParam Firewall Flood Protection Profile ID (required)
    // @param floodProtectionProfileParam (required)
    // The parameter must contain all the properties defined in model.FloodProtectionProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(floodProtectionProfileIdParam string, floodProtectionProfileParam *data.StructValue) error

    // API will update Firewall Flood Protection Profile
    //
    // @param floodProtectionProfileIdParam Flood Protection Profile ID (required)
    // @param floodProtectionProfileParam (required)
    // The parameter must contain all the properties defined in model.FloodProtectionProfile.
    // @return com.vmware.nsx_global_policy.model.FloodProtectionProfile
    // The return value will contain all the properties defined in model.FloodProtectionProfile.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(floodProtectionProfileIdParam string, floodProtectionProfileParam *data.StructValue) (*data.StructValue, error)
}
