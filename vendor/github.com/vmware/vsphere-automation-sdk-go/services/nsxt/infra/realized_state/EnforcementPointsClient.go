/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EnforcementPoints
 * Used by client-side stubs.
 */

package realized_state

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EnforcementPointsClient interface {

    // Read a Enforcement Point and the complete tree underneath. Returns the populated enforcement point object.
    //
    // @param enforcementPointNameParam Enforcement Point Name (required)
    // @return com.vmware.nsx_policy.model.RealizedEnforcementPoint
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(enforcementPointNameParam string) (model.RealizedEnforcementPoint, error)

    // Paginated list of all enforcement points. Returns the populated enforcement points.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.RealizedEnforcementPointListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.RealizedEnforcementPointListResult, error)
}
