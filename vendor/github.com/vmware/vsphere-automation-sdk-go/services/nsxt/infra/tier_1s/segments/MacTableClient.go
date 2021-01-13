/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: MacTable
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type MacTableClient interface {

    // Returns MAC table for a segment
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @param sourceParam Data source type. (optional)
    // @param transportNodeIdParam TransportNode Id (optional)
    // @return com.vmware.nsx_policy.model.SegmentMacAddressListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, segmentIdParam string, cursorParam *string, enforcementPointPathParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, sourceParam *string, transportNodeIdParam *string) (model.SegmentMacAddressListResult, error)
}
