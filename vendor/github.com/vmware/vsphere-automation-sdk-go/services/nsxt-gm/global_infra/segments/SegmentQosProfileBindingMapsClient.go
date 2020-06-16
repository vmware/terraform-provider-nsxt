/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SegmentQosProfileBindingMaps
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SegmentQosProfileBindingMapsClient interface {

    // API will delete Segment QoS Profile Binding Profile.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentQosProfileBindingMapIdParam Segment QoS Profile Binding Map ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(segmentIdParam string, segmentQosProfileBindingMapIdParam string) error

    // API will get Segment QoS Profile Binding Map.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentQosProfileBindingMapIdParam Segment QoS Profile Binding Map ID (required)
    // @return com.vmware.nsx_global_policy.model.SegmentQoSProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(segmentIdParam string, segmentQosProfileBindingMapIdParam string) (model.SegmentQosProfileBindingMap, error)

    // API will list all Segment QoS Profile Binding Maps in current segment id.
    //
    // @param segmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.SegmentQoSProfileBindingMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(segmentIdParam string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentQosProfileBindingMapListResult, error)

    // API will create segment QoS profile binding map.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentQosProfileBindingMapIdParam Segment QoS Profile Binding Map ID (required)
    // @param segmentQosProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(segmentIdParam string, segmentQosProfileBindingMapIdParam string, segmentQosProfileBindingMapParam model.SegmentQosProfileBindingMap) error

    // API will update Segment QoS Profile Binding Map.
    //
    // @param segmentIdParam Segment ID (required)
    // @param segmentQosProfileBindingMapIdParam Segment QoS Profile Binding Map ID (required)
    // @param segmentQosProfileBindingMapParam (required)
    // @return com.vmware.nsx_global_policy.model.SegmentQoSProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(segmentIdParam string, segmentQosProfileBindingMapIdParam string, segmentQosProfileBindingMapParam model.SegmentQosProfileBindingMap) (model.SegmentQosProfileBindingMap, error)
}
