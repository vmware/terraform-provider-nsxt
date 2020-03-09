/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SegmentMonitoringProfileBindingMaps
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SegmentMonitoringProfileBindingMapsClient interface {

    // API will delete Infra Segment Monitoring Profile Binding Profile.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) error

    // API will get Infra Segment Monitoring Profile Binding Map.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
    // @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string) (model.SegmentMonitoringProfileBindingMap, error)

    // API will list all Infra Segment Monitoring Profile Binding Maps in current segment id.
    //
    // @param infraSegmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentMonitoringProfileBindingMapListResult, error)

    // API will create infra segment monitoring profile binding map.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
    // @param segmentMonitoringProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) error

    // API will update Infra Segment Monitoring Profile Binding Map.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentMonitoringProfileBindingMapIdParam Segment Monitoring Profile Binding Map ID (required)
    // @param segmentMonitoringProfileBindingMapParam (required)
    // @return com.vmware.nsx_policy.model.SegmentMonitoringProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(infraSegmentIdParam string, segmentMonitoringProfileBindingMapIdParam string, segmentMonitoringProfileBindingMapParam model.SegmentMonitoringProfileBindingMap) (model.SegmentMonitoringProfileBindingMap, error)
}
