/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbVirtualServices
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbVirtualServicesClient interface {

    // Delete the ALBVirtualService along with all the entities contained by this ALBVirtualService.
    //
    // @param albVirtualserviceIdParam ALBVirtualService ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albVirtualserviceIdParam string, forceParam *bool) error

    // Read a ALBVirtualService.
    //
    // @param albVirtualserviceIdParam ALBVirtualService ID (required)
    // @return com.vmware.nsx_policy.model.ALBVirtualService
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albVirtualserviceIdParam string) (model.ALBVirtualService, error)

    // Paginated list of all ALBVirtualService for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBVirtualServiceApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBVirtualServiceApiResponse, error)

    // If a ALBvirtualservice with the alb-virtualservice-id is not already present, create a new ALBvirtualservice. If it already exists, update the ALBvirtualservice. This is a full replace.
    //
    // @param albVirtualserviceIdParam ALBvirtualservice ID (required)
    // @param aLBVirtualServiceParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albVirtualserviceIdParam string, aLBVirtualServiceParam model.ALBVirtualService) error

    // If a ALBVirtualService with the alb-VirtualService-id is not already present, create a new ALBVirtualService. If it already exists, update the ALBVirtualService. This is a full replace.
    //
    // @param albVirtualserviceIdParam ALBVirtualService ID (required)
    // @param aLBVirtualServiceParam (required)
    // @return com.vmware.nsx_policy.model.ALBVirtualService
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albVirtualserviceIdParam string, aLBVirtualServiceParam model.ALBVirtualService) (model.ALBVirtualService, error)
}
