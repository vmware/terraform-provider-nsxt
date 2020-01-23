/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpsecVpnDpdProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpsecVpnDpdProfilesClient interface {

    // Delete custom dead peer detection (DPD) profile. Profile can not be deleted if profile has references to it.
    //
    // @param dpdProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(dpdProfileIdParam string) error

    // Get custom dead peer detection (DPD) profile, given the particular id.
    //
    // @param dpdProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnDpdProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(dpdProfileIdParam string) (model.IPSecVpnDpdProfile, error)

    // Get paginated list of all DPD Profiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPSecVpnDpdProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnDpdProfileListResult, error)

    // Create or patch dead peer detection (DPD) profile. Any change in profile affects all sessions consuming this profile. System will be provisioned with system owned editable default DPD profile. Any change in profile affects all sessions consuming this profile.
    //
    // @param dpdProfileIdParam (required)
    // @param ipSecVpnDpdProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(dpdProfileIdParam string, ipSecVpnDpdProfileParam model.IPSecVpnDpdProfile) error

    // Create or patch dead peer detection (DPD) profile. Any change in profile affects all sessions consuming this profile. System will be provisioned with system owned editable default DPD profile. Any change in profile affects all sessions consuming this profile. Revision is optional for creation and required for update.
    //
    // @param dpdProfileIdParam (required)
    // @param ipSecVpnDpdProfileParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnDpdProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(dpdProfileIdParam string, ipSecVpnDpdProfileParam model.IPSecVpnDpdProfile) (model.IPSecVpnDpdProfile, error)
}
