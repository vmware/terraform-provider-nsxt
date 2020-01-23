/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpsecVpnIkeProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpsecVpnIkeProfilesClient interface {

    // Delete custom IKE Profile. Profile can not be deleted if profile has references to it.
    //
    // @param ikeProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ikeProfileIdParam string) error

    // Get custom IKE Profile, given the particular id.
    //
    // @param ikeProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnIkeProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ikeProfileIdParam string) (model.IPSecVpnIkeProfile, error)

    // Get paginated list of all IKE Profiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPSecVpnIkeProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnIkeProfileListResult, error)

    // Create or patch custom internet key exchange (IKE) Profile. IKE Profile is a reusable profile that captures IKE and phase one negotiation parameters. System will be pre provisioned with system owned editable default IKE profile and suggested set of profiles that can be used for peering with popular remote peers like AWS VPN. User can create custom profiles as needed. Any change in profile affects all sessions consuming this profile.
    //
    // @param ikeProfileIdParam (required)
    // @param ipSecVpnIkeProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ikeProfileIdParam string, ipSecVpnIkeProfileParam model.IPSecVpnIkeProfile) error

    // Create or fully replace custom internet key exchange (IKE) Profile. IKE Profile is a reusable profile that captures IKE and phase one negotiation parameters. System will be pre provisioned with system owned editable default IKE profile and suggested set of profiles that can be used for peering with popular remote peers like AWS VPN. User can create custom profiles as needed. Any change in profile affects all sessions consuming this profile. Revision is optional for creation and required for update.
    //
    // @param ikeProfileIdParam (required)
    // @param ipSecVpnIkeProfileParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnIkeProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ikeProfileIdParam string, ipSecVpnIkeProfileParam model.IPSecVpnIkeProfile) (model.IPSecVpnIkeProfile, error)
}
