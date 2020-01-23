/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Vifs
 * Used by client-side stubs.
 */

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type VifsClient interface {

    // This API lists VIFs from the specified NSX Manager.
    //
    // @param enforcementPointNameParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param lportAttachmentIdParam LPort attachment ID of the VIF. (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.VirtualNetworkInterfaceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(enforcementPointNameParam string, cursorParam *string, includedFieldsParam *string, lportAttachmentIdParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.VirtualNetworkInterfaceListResult, error)
}
