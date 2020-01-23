/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpsecVpnTunnelProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpsecVpnTunnelProfilesClient interface {

    // Delete custom IPSec tunnel Profile. Profile can not be deleted if profile has references to it.
    //
    // @param tunnelProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tunnelProfileIdParam string) error

    // Get custom IPSec tunnel Profile, given the particular id.
    //
    // @param tunnelProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnTunnelProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tunnelProfileIdParam string) (model.IPSecVpnTunnelProfile, error)

    // Get paginated list of all IPSec tunnel Profiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPSecVpnTunnelProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnTunnelProfileListResult, error)

    // Create or patch custom IPSec tunnel profile. IPSec tunnel profile is a reusable profile that captures phase two negotiation parameters and tunnel properties. System will be provisioned with system owned editable default IPSec tunnel profile. Any change in profile affects all sessions consuming this profile.
    //
    // @param tunnelProfileIdParam (required)
    // @param ipSecVpnTunnelProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tunnelProfileIdParam string, ipSecVpnTunnelProfileParam model.IPSecVpnTunnelProfile) error

    // Create or fully replace custom IPSec tunnel profile. IPSec tunnel profile is a reusable profile that captures phase two negotiation parameters and tunnel properties. System will be provisioned with system owned editable default IPSec tunnel profile. Any change in profile affects all sessions consuming this profile. Revision is optional for creation and required for update.
    //
    // @param tunnelProfileIdParam (required)
    // @param ipSecVpnTunnelProfileParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnTunnelProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tunnelProfileIdParam string, ipSecVpnTunnelProfileParam model.IPSecVpnTunnelProfile) (model.IPSecVpnTunnelProfile, error)
}
