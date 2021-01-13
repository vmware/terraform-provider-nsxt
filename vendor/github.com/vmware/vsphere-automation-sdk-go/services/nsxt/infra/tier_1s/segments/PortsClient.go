/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Ports
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type PortsClient interface {

    // Delete a Tier-1 segment port by giving ID.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param portIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, segmentIdParam string, portIdParam string) error

    // Get detail information on a Tier-1 segment port by giving ID.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param portIdParam (required)
    // @return com.vmware.nsx_policy.model.SegmentPort
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, segmentIdParam string, portIdParam string) (model.SegmentPort, error)

    // List all the ports for a Tier-1 segment.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.SegmentPortListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, segmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.SegmentPortListResult, error)

    // Create a Tier-1 segment port if it does not exist based on the IDs, or update existing port information by replacing the port object fields which presents in the request body.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param portIdParam (required)
    // @param segmentPortParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, segmentIdParam string, portIdParam string, segmentPortParam model.SegmentPort) error

    // Create a Tier-1 segment port if it does not exist based on the IDs, or update existing port information by replacing the port object already exists.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param portIdParam (required)
    // @param segmentPortParam (required)
    // @return com.vmware.nsx_policy.model.SegmentPort
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, segmentIdParam string, portIdParam string, segmentPortParam model.SegmentPort) (model.SegmentPort, error)
}
