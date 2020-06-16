/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Segments
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SegmentsClient interface {

    // Delete infra segment
    //
    // @param segmentIdParam Segment ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(segmentIdParam string) error

    // Force delete bypasses validations during segment deletion. This may result in an inconsistent connectivity.
    //
    // @param segmentIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete0(segmentIdParam string) error

    // Read infra segment
    //
    // @param segmentIdParam Segment ID (required)
    // @return com.vmware.nsx_global_policy.model.Segment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(segmentIdParam string) (model.Segment, error)

    // Paginated list of all segments under infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.SegmentListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentListResult, error)

    // If segment with the segment-id is not already present, create a new segment. If it already exists, update the segment with specified attributes.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(segmentIdParam string, segmentParam model.Segment) error

    // If segment with the segment-id is not already present, create a new segment. If it already exists, update the segment with specified attributes. Force parameter is required when workload connectivity is indirectly impacted with the current update.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch0(segmentIdParam string, segmentParam model.Segment) error

    // If segment with the segment-id is not already present, create a new segment. If it already exists, replace the segment with this object.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentParam (required)
    // @return com.vmware.nsx_global_policy.model.Segment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(segmentIdParam string, segmentParam model.Segment) (model.Segment, error)

    // If segment with the segment-id is not already present, create a new segment. If it already exists, replace the segment with this object. Force parameter is required when workload connectivity is indirectly impacted with the current replacement.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentParam (required)
    // @return com.vmware.nsx_global_policy.model.Segment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update0(segmentIdParam string, segmentParam model.Segment) (model.Segment, error)
}
