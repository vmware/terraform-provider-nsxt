/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SegmentSecurityProfileBindingMaps
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SegmentSecurityProfileBindingMapsClient interface {

    // API will delete segment security profile binding map.
    //
    // @param segmentIdParam segment id (required)
    // @param segmentSecurityProfileBindingMapIdParam segment security profile binding map id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(segmentIdParam string, segmentSecurityProfileBindingMapIdParam string) error

    // API will return details of the segment security profile binding map. If the binding map does not exist, it will return 404.
    //
    // @param segmentIdParam segment id (required)
    // @param segmentSecurityProfileBindingMapIdParam segment security profile binding map id (required)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(segmentIdParam string, segmentSecurityProfileBindingMapIdParam string) (model.SegmentSecurityProfileBindingMap, error)

    // API will list all segment security profile binding maps.
    //
    // @param segmentIdParam segment id (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfileBindingMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(segmentIdParam string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentSecurityProfileBindingMapListResult, error)

    // Create a new segment security profile binding map if the given security profile binding map does not exist. Otherwise, patch the existing segment security profile binding map.
    //
    // @param segmentIdParam segment id (required)
    // @param segmentSecurityProfileBindingMapIdParam segment security profile binding map id (required)
    // @param segmentSecurityProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(segmentIdParam string, segmentSecurityProfileBindingMapIdParam string, segmentSecurityProfileBindingMapParam model.SegmentSecurityProfileBindingMap) error

    // API will create or replace segment security profile binding map.
    //
    // @param segmentIdParam segment id (required)
    // @param segmentSecurityProfileBindingMapIdParam segment security profile binding map id (required)
    // @param segmentSecurityProfileBindingMapParam (required)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(segmentIdParam string, segmentSecurityProfileBindingMapIdParam string, segmentSecurityProfileBindingMapParam model.SegmentSecurityProfileBindingMap) (model.SegmentSecurityProfileBindingMap, error)
}
