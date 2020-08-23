/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SegmentSecurityProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SegmentSecurityProfilesClient interface {

    // API will delete segment security profile with the given id.
    //
    // @param segmentSecurityProfileIdParam Segment security profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(segmentSecurityProfileIdParam string) error

    // API will return details of the segment security profile with given id. If the profile does not exist, it will return 404.
    //
    // @param segmentSecurityProfileIdParam Segment security profile id (required)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(segmentSecurityProfileIdParam string) (model.SegmentSecurityProfile, error)

    // API will list all segment security profiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentSecurityProfileListResult, error)

    // Create a new segment security profile if the segment security profile with given id does not exist. Otherwise, PATCH the existing segment security profile
    //
    // @param segmentSecurityProfileIdParam Segment security profile id (required)
    // @param segmentSecurityProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(segmentSecurityProfileIdParam string, segmentSecurityProfileParam model.SegmentSecurityProfile) error

    // Create or replace a segment security profile
    //
    // @param segmentSecurityProfileIdParam Segment security profile id (required)
    // @param segmentSecurityProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.SegmentSecurityProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(segmentSecurityProfileIdParam string, segmentSecurityProfileParam model.SegmentSecurityProfile) (model.SegmentSecurityProfile, error)
}
