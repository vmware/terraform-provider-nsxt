/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbAutoScaleLaunchConfigs
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbAutoScaleLaunchConfigsClient interface {

    // Delete the ALBAutoScaleLaunchConfig along with all the entities contained by this ALBAutoScaleLaunchConfig.
    //
    // @param albAutoscalelaunchconfigIdParam ALBAutoScaleLaunchConfig ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albAutoscalelaunchconfigIdParam string, forceParam *bool) error

    // Read a ALBAutoScaleLaunchConfig.
    //
    // @param albAutoscalelaunchconfigIdParam ALBAutoScaleLaunchConfig ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBAutoScaleLaunchConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albAutoscalelaunchconfigIdParam string) (model.ALBAutoScaleLaunchConfig, error)

    // Paginated list of all ALBAutoScaleLaunchConfig for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBAutoScaleLaunchConfigApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBAutoScaleLaunchConfigApiResponse, error)

    // If a ALBautoscalelaunchconfig with the alb-autoscalelaunchconfig-id is not already present, create a new ALBautoscalelaunchconfig. If it already exists, update the ALBautoscalelaunchconfig. This is a full replace.
    //
    // @param albAutoscalelaunchconfigIdParam ALBautoscalelaunchconfig ID (required)
    // @param aLBAutoScaleLaunchConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albAutoscalelaunchconfigIdParam string, aLBAutoScaleLaunchConfigParam model.ALBAutoScaleLaunchConfig) error

    // If a ALBAutoScaleLaunchConfig with the alb-AutoScaleLaunchConfig-id is not already present, create a new ALBAutoScaleLaunchConfig. If it already exists, update the ALBAutoScaleLaunchConfig. This is a full replace.
    //
    // @param albAutoscalelaunchconfigIdParam ALBAutoScaleLaunchConfig ID (required)
    // @param aLBAutoScaleLaunchConfigParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBAutoScaleLaunchConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albAutoscalelaunchconfigIdParam string, aLBAutoScaleLaunchConfigParam model.ALBAutoScaleLaunchConfig) (model.ALBAutoScaleLaunchConfig, error)
}
