/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Traceflows
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type TraceflowsClient interface {

    // This will retrace even if current traceflow has observations. Current observations will be lost. Traceflow configuration will be cleaned up by the system after two hours of inactivity.
    //
    // @param traceflowIdParam (required)
    // @param actionParam Action to be performed (optional)
    // @return com.vmware.nsx_policy.model.TraceflowConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(traceflowIdParam string, actionParam *string) (model.TraceflowConfig, error)

    // Delete traceflow config with id traceflow-id
    //
    // @param traceflowIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(traceflowIdParam string) error

    // Read traceflow config with id traceflow-id. This configuration will be cleaned up by the system after two hours of inactivity.
    //
    // @param traceflowIdParam (required)
    // @return com.vmware.nsx_policy.model.TraceflowConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(traceflowIdParam string) (model.TraceflowConfig, error)

    // Paginated list of all TraceflowConfig for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.TraceflowConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.TraceflowConfigListResult, error)

    // If a traceflow config with the traceflow-id is not already present, create a new traceflow config. If it already exists, update the traceflow config. This is a full replace. This configuration will be cleaned up by the system after two hours of inactivity.
    //
    // @param traceflowIdParam (required)
    // @param traceflowConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(traceflowIdParam string, traceflowConfigParam model.TraceflowConfig) error

    // If a traceflow config with the traceflow-id is not already present, create a new traceflow config. If it already exists, update the traceflow config. This is a full replace. This configuration will be cleaned up by the system after two hours of inactivity.
    //
    // @param traceflowIdParam (required)
    // @param traceflowConfigParam (required)
    // @return com.vmware.nsx_policy.model.TraceflowConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(traceflowIdParam string, traceflowConfigParam model.TraceflowConfig) (model.TraceflowConfig, error)
}
