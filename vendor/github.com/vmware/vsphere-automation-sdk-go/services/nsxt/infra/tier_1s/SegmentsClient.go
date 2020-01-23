/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Segments
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SegmentsClient interface {

    // Delete segment
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, segmentIdParam string) error

    // Force delete bypasses validations during segment deletion. This may result in an inconsistent connectivity.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete0(tier1IdParam string, segmentIdParam string) error

    // Read segment
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @return com.vmware.nsx_policy.model.Segment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, segmentIdParam string) (model.Segment, error)

    // Paginated list of all segments under Tier-1 instance
    //
    // @param tier1IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.SegmentListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentListResult, error)

    // If segment with the segment-id is not already present, create a new segment. If it already exists, update the segment with specified attributes.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param segmentParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, segmentIdParam string, segmentParam model.Segment) error

    // If segment with the segment-id is not already present, create a new segment. If it already exists, replace the segment with this object.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param segmentParam (required)
    // @return com.vmware.nsx_policy.model.Segment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, segmentIdParam string, segmentParam model.Segment) (model.Segment, error)
}
