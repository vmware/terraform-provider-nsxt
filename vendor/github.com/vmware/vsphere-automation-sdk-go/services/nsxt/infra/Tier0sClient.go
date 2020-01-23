/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Tier0s
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type Tier0sClient interface {

    // Delete Tier-0
    //
    // @param tier0IdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string) error

    // Read Tier-0
    //
    // @param tier0IdParam (required)
    // @return com.vmware.nsx_policy.model.Tier0
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string) (model.Tier0, error)

    // Paginated list of all Tier-0s
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.Tier0ListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Tier0ListResult, error)

    // If a Tier-0 with the tier-0-id is not already present, create a new Tier-0. If it already exists, update the Tier-0 for specified attributes.
    //
    // @param tier0IdParam (required)
    // @param tier0Param (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, tier0Param model.Tier0) error

    // Reprocess Tier0 gateway configuration and configuration of related entities like Tier0 interfaces and static routes, etc. Any missing Updates are published to NSX controller.
    //
    // @param tier0IdParam (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Reprocess(tier0IdParam string, enforcementPointPathParam *string) error

    // If a Tier-0 with the tier-0-id is not already present, create a new Tier-0. If it already exists, replace the Tier-0 instance with the new object.
    //
    // @param tier0IdParam (required)
    // @param tier0Param (required)
    // @return com.vmware.nsx_policy.model.Tier0
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, tier0Param model.Tier0) (model.Tier0, error)
}
