/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: State
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type StateClient interface {

    // Returns
    //
    // @param tier1IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param enforcementPointPathParam Enforcement point path (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param interfacePathParam Interface path for interface specific state such as IPv6 DAD state (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.Tier1GatewayState
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, cursorParam *string, enforcementPointPathParam *string, includedFieldsParam *string, interfacePathParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Tier1GatewayState, error)
}
