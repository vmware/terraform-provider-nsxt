/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SegmentDiscoveryProfileBindingMaps
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SegmentDiscoveryProfileBindingMapsClient interface {

    // API will delete Segment Discovery Profile Binding Profile
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) error

    // API will get Infra Segment Discovery Profile Binding Map
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
    // @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string) (model.SegmentDiscoveryProfileBindingMap, error)

    // API will list all Infra Segment Discovery Profile Binding Maps in current segment id.
    //
    // @param infraSegmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(infraSegmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentDiscoveryProfileBindingMapListResult, error)

    // API will create Infra Segment Discovery Profile Binding Map. For objects with no binding maps, default profile is applied.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
    // @param segmentDiscoveryProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) error

    // API will update Infra Segment Discovery Profile Binding Map. For objects with no binding maps, default profile is applied.
    //
    // @param infraSegmentIdParam Infra Segment ID (required)
    // @param segmentDiscoveryProfileBindingMapIdParam Segment Discovery Profile Binding Map ID (required)
    // @param segmentDiscoveryProfileBindingMapParam (required)
    // @return com.vmware.nsx_global_policy.model.SegmentDiscoveryProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(infraSegmentIdParam string, segmentDiscoveryProfileBindingMapIdParam string, segmentDiscoveryProfileBindingMapParam model.SegmentDiscoveryProfileBindingMap) (model.SegmentDiscoveryProfileBindingMap, error)
}
