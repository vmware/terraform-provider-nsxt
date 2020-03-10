/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ServiceSegments
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ServiceSegmentsClient interface {

    // Delete Service Segment with given ID
    //
    // @param serviceSegmentIdParam Service Segment ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(serviceSegmentIdParam string) error

    // Read a Service Segment with the given id
    //
    // @param serviceSegmentIdParam Service Segment ID (required)
    // @return com.vmware.nsx_policy.model.ServiceSegment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(serviceSegmentIdParam string) (model.ServiceSegment, error)

    // Paginated list of all Service Segments
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ServiceSegmentListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ServiceSegmentListResult, error)

    // A service segment with the service-segment-id is created. Modification of service segment is not supported.
    //
    // @param serviceSegmentIdParam Service Segment ID (required)
    // @param serviceSegmentParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(serviceSegmentIdParam string, serviceSegmentParam model.ServiceSegment) error

    // A service segment with the service-segment-id is created. Modification of service segment is not supported.
    //
    // @param serviceSegmentIdParam Service Segment ID (required)
    // @param serviceSegmentParam (required)
    // @return com.vmware.nsx_policy.model.ServiceSegment
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(serviceSegmentIdParam string, serviceSegmentParam model.ServiceSegment) (model.ServiceSegment, error)
}
