/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbHealthMonitors
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbHealthMonitorsClient interface {

    // Delete the ALBHealthMonitor along with all the entities contained by this ALBHealthMonitor.
    //
    // @param albHealthmonitorIdParam ALBHealthMonitor ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albHealthmonitorIdParam string, forceParam *bool) error

    // Read a ALBHealthMonitor.
    //
    // @param albHealthmonitorIdParam ALBHealthMonitor ID (required)
    // @return com.vmware.nsx_policy.model.ALBHealthMonitor
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albHealthmonitorIdParam string) (model.ALBHealthMonitor, error)

    // Paginated list of all ALBHealthMonitor for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBHealthMonitorApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBHealthMonitorApiResponse, error)

    // If a ALBhealthmonitor with the alb-healthmonitor-id is not already present, create a new ALBhealthmonitor. If it already exists, update the ALBhealthmonitor. This is a full replace.
    //
    // @param albHealthmonitorIdParam ALBhealthmonitor ID (required)
    // @param aLBHealthMonitorParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albHealthmonitorIdParam string, aLBHealthMonitorParam model.ALBHealthMonitor) error

    // If a ALBHealthMonitor with the alb-HealthMonitor-id is not already present, create a new ALBHealthMonitor. If it already exists, update the ALBHealthMonitor. This is a full replace.
    //
    // @param albHealthmonitorIdParam ALBHealthMonitor ID (required)
    // @param aLBHealthMonitorParam (required)
    // @return com.vmware.nsx_policy.model.ALBHealthMonitor
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albHealthmonitorIdParam string, aLBHealthMonitorParam model.ALBHealthMonitor) (model.ALBHealthMonitor, error)
}
