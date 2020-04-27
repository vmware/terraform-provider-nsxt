/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: MacDiscoveryProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type MacDiscoveryProfilesClient interface {

    // API will delete Mac Discovery profile.
    //
    // @param macDiscoveryProfileIdParam Mac Discovery Profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(macDiscoveryProfileIdParam string) error

    // API will get Mac Discovery profile.
    //
    // @param macDiscoveryProfileIdParam Mac Discovery Profile ID (required)
    // @return com.vmware.nsx_policy.model.MacDiscoveryProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(macDiscoveryProfileIdParam string) (model.MacDiscoveryProfile, error)

    // API will list all Mac Discovery Profiles active in current discovery profile id.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.MacDiscoveryProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.MacDiscoveryProfileListResult, error)

    // API will create Mac Discovery profile.
    //
    // @param macDiscoveryProfileIdParam Mac Discovery Profile ID (required)
    // @param macDiscoveryProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(macDiscoveryProfileIdParam string, macDiscoveryProfileParam model.MacDiscoveryProfile) error

    // API will update Mac Discovery profile.
    //
    // @param macDiscoveryProfileIdParam Mac Discovery Profile ID (required)
    // @param macDiscoveryProfileParam (required)
    // @return com.vmware.nsx_policy.model.MacDiscoveryProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(macDiscoveryProfileIdParam string, macDiscoveryProfileParam model.MacDiscoveryProfile) (model.MacDiscoveryProfile, error)
}
