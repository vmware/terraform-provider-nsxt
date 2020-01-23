/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Tasks
 * Used by client-side stubs.
 */

package nsx_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type TasksClient interface {

    // Get information about the specified task
    //
    // @param taskIdParam ID of task to read (required)
    // @return com.vmware.nsx_policy.model.TaskProperties
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(taskIdParam string) (model.TaskProperties, error)

    // Get information about all tasks
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param requestUriParam Request URI(s) to include in query result (optional)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @param statusParam Status(es) to include in query result (optional)
    // @param userParam Names of users to include in query result (optional)
    // @return com.vmware.nsx_policy.model.TaskListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, requestUriParam *string, sortAscendingParam *bool, sortByParam *string, statusParam *string, userParam *string) (model.TaskListResult, error)
}
