/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: L3vpns
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type L3vpnsClient interface {

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param l3vpnIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) error

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param l3vpnIdParam (required)
    // @return com.vmware.nsx_policy.model.L3Vpn
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error)

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param l3vpnSessionParam Resource type of L3Vpn Session (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.L3VpnListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, l3vpnSessionParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.L3VpnListResult, error)

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param l3vpnIdParam (required)
    // @param l3VpnParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) error

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param l3vpnIdParam (required)
    // @return com.vmware.nsx_policy.model.L3Vpn
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Showsensitivedata(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string) (model.L3Vpn, error)

    //
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param l3vpnIdParam (required)
    // @param l3VpnParam (required)
    // @return com.vmware.nsx_policy.model.L3Vpn
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, l3vpnIdParam string, l3VpnParam model.L3Vpn) (model.L3Vpn, error)
}
