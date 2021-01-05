/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbHardwareSecurityModuleGroups
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbHardwareSecurityModuleGroupsClient interface {

    // Delete the ALBHardwareSecurityModuleGroup along with all the entities contained by this ALBHardwareSecurityModuleGroup.
    //
    // @param albHardwaresecuritymodulegroupIdParam ALBHardwareSecurityModuleGroup ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albHardwaresecuritymodulegroupIdParam string, forceParam *bool) error

    // Read a ALBHardwareSecurityModuleGroup.
    //
    // @param albHardwaresecuritymodulegroupIdParam ALBHardwareSecurityModuleGroup ID (required)
    // @return com.vmware.nsx_policy.model.ALBHardwareSecurityModuleGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albHardwaresecuritymodulegroupIdParam string) (model.ALBHardwareSecurityModuleGroup, error)

    // Paginated list of all ALBHardwareSecurityModuleGroup for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBHardwareSecurityModuleGroupApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBHardwareSecurityModuleGroupApiResponse, error)

    // If a ALBhardwaresecuritymodulegroup with the alb-hardwaresecuritymodulegroup-id is not already present, create a new ALBhardwaresecuritymodulegroup. If it already exists, update the ALBhardwaresecuritymodulegroup. This is a full replace.
    //
    // @param albHardwaresecuritymodulegroupIdParam ALBhardwaresecuritymodulegroup ID (required)
    // @param aLBHardwareSecurityModuleGroupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albHardwaresecuritymodulegroupIdParam string, aLBHardwareSecurityModuleGroupParam model.ALBHardwareSecurityModuleGroup) error

    // If a ALBHardwareSecurityModuleGroup with the alb-HardwareSecurityModuleGroup-id is not already present, create a new ALBHardwareSecurityModuleGroup. If it already exists, update the ALBHardwareSecurityModuleGroup. This is a full replace.
    //
    // @param albHardwaresecuritymodulegroupIdParam ALBHardwareSecurityModuleGroup ID (required)
    // @param aLBHardwareSecurityModuleGroupParam (required)
    // @return com.vmware.nsx_policy.model.ALBHardwareSecurityModuleGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albHardwaresecuritymodulegroupIdParam string, aLBHardwareSecurityModuleGroupParam model.ALBHardwareSecurityModuleGroup) (model.ALBHardwareSecurityModuleGroup, error)
}
