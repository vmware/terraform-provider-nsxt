/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpsecVpnServices
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpsecVpnServicesClient interface {

    // Delete IPSec VPN service for a given locale service under Tier-1.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param serviceIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, localeServiceIdParam string, serviceIdParam string) error

    // Get IPSec VPN service for a given locale service under Tier-1.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param serviceIdParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnService
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServiceIdParam string, serviceIdParam string) (model.IPSecVpnService, error)

    // Get paginated list of all IPSec VPN services for a given locale service under Tier-1.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPSecVpnServiceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPSecVpnServiceListResult, error)

    // Create or patch IPSec VPN service for a given locale service under Tier-1.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param serviceIdParam (required)
    // @param ipSecVpnServiceParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, ipSecVpnServiceParam model.IPSecVpnService) error

    // Create or fully replace IPSec VPN service for a given locale service under Tier-1. Revision is optional for creation and required for update.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param serviceIdParam (required)
    // @param ipSecVpnServiceParam (required)
    // @return com.vmware.nsx_policy.model.IPSecVpnService
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, localeServiceIdParam string, serviceIdParam string, ipSecVpnServiceParam model.IPSecVpnService) (model.IPSecVpnService, error)
}
