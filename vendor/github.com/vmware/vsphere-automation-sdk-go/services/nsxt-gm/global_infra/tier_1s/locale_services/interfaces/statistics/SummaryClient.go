/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Summary
 * Used by client-side stubs.
 */

package statistics

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SummaryClient interface {

    // Get tier-1 interface statistics information.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param enforcementPointPathParam Enforcement point path (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.PolicyInterfaceStatisticsSummary
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServiceIdParam string, interfaceIdParam string, cursorParam *string, enforcementPointPathParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyInterfaceStatisticsSummary, error)
}
