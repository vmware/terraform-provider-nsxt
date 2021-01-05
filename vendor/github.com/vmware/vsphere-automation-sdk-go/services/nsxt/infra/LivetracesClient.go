/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Livetraces
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LivetracesClient interface {

    // Restart a livetrace session with the same set of parameters used in creating or updating of a livetrace config.
    //
    // @param livetraceIdParam (required)
    // @param actionParam Action to be performed (optional)
    // @return com.vmware.nsx_policy.model.LiveTraceConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(livetraceIdParam string, actionParam *string) (model.LiveTraceConfig, error)

    // Delete livetrace config with the specified identifier.
    //
    // @param livetraceIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(livetraceIdParam string) error

    // Read livetrace config with the specified identifier.
    //
    // @param livetraceIdParam (required)
    // @return com.vmware.nsx_policy.model.LiveTraceConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(livetraceIdParam string) (model.LiveTraceConfig, error)

    // Get a paginated list of all livetrace config entities.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LiveTraceConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LiveTraceConfigListResult, error)

    // If a livetrace config with the specified identifier is not present, then create a new livetrace config. If it already exists, update the livetrace config with a full replacement.
    //
    // @param livetraceIdParam (required)
    // @param liveTraceConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(livetraceIdParam string, liveTraceConfigParam model.LiveTraceConfig) error

    // If a livetrace config with the specified identifier is not present, then create a new livetrace config. If it already exists, update the livetrace config with a full replacement.
    //
    // @param livetraceIdParam (required)
    // @param liveTraceConfigParam (required)
    // @return com.vmware.nsx_policy.model.LiveTraceConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(livetraceIdParam string, liveTraceConfigParam model.LiveTraceConfig) (model.LiveTraceConfig, error)
}
