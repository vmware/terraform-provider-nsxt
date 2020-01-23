/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EdgeBridgeProfiles
 * Used by client-side stubs.
 */

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EdgeBridgeProfilesClient interface {

    // API will delete L2 bridge profile with ID profile-id
    //
    // @param siteIdParam site ID (required)
    // @param enforcementPointIdParam enforcement point ID (required)
    // @param profileIdParam profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(siteIdParam string, enforcementPointIdParam string, profileIdParam string) error

    // Read L2 bridge profile with ID profile-id
    //
    // @param siteIdParam site ID (required)
    // @param enforcementPointIdParam enforcement point ID (required)
    // @param profileIdParam profile ID (required)
    // @return com.vmware.nsx_policy.model.L2BridgeEndpointProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(siteIdParam string, enforcementPointIdParam string, profileIdParam string) (model.L2BridgeEndpointProfile, error)

    // List all L2 bridge profiles
    //
    // @param siteIdParam site ID (required)
    // @param enforcementPointIdParam enforcement point ID (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.L2BridgeEndpointProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(siteIdParam string, enforcementPointIdParam string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L2BridgeEndpointProfileListResult, error)

    // API will create or update L2 bridge profile with ID profile-id
    //
    // @param siteIdParam site ID (required)
    // @param enforcementPointIdParam enforcement point ID (required)
    // @param profileIdParam profile ID (required)
    // @param l2BridgeEndpointProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementPointIdParam string, profileIdParam string, l2BridgeEndpointProfileParam model.L2BridgeEndpointProfile) error

    // API will create or update L2 bridge profile with ID profile-id
    //
    // @param siteIdParam site ID (required)
    // @param enforcementPointIdParam enforcement point ID (required)
    // @param profileIdParam profile ID (required)
    // @param l2BridgeEndpointProfileParam (required)
    // @return com.vmware.nsx_policy.model.L2BridgeEndpointProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(siteIdParam string, enforcementPointIdParam string, profileIdParam string, l2BridgeEndpointProfileParam model.L2BridgeEndpointProfile) (model.L2BridgeEndpointProfile, error)
}
