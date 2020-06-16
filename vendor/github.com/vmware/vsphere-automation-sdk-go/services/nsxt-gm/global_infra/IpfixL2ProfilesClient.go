/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixL2Profiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type IpfixL2ProfilesClient interface {

    // API deletes IPFIX L2 Profile. Flow forwarding to selected collector will be stopped.
    //
    // @param ipfixL2ProfileIdParam IPFIX L2 Profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipfixL2ProfileIdParam string) error

    // API will return details of IPFIX L2 profile.
    //
    // @param ipfixL2ProfileIdParam IPFIX L2 profile id (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2Profile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipfixL2ProfileIdParam string) (model.IPFIXL2Profile, error)

    // API provides list IPFIX L2 Profiles available on selected logical l2.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2ProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXL2ProfileListResult, error)

    // Create a new IPFIX L2 profile if the IPFIX L2 profile with given id does not already exist. If the IPFIX L2 profile with the given id already exists, patch with the existing IPFIX L2 profile.
    //
    // @param ipfixL2ProfileIdParam IPFIX L2 Profile ID (required)
    // @param iPFIXL2ProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipfixL2ProfileIdParam string, iPFIXL2ProfileParam model.IPFIXL2Profile) error

    // Create or replace IPFIX L2 Profile. Profile is reusable entity. Single profile can attached multiple bindings e.g group, segment and port.
    //
    // @param ipfixL2ProfileIdParam IPFIX L2 Profile ID (required)
    // @param iPFIXL2ProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2Profile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipfixL2ProfileIdParam string, iPFIXL2ProfileParam model.IPFIXL2Profile) (model.IPFIXL2Profile, error)
}
