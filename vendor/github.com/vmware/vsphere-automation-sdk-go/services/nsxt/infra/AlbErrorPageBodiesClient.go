/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbErrorPageBodies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbErrorPageBodiesClient interface {

    // Delete the ALBErrorPageBody along with all the entities contained by this ALBErrorPageBody.
    //
    // @param albErrorpagebodyIdParam ALBErrorPageBody ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albErrorpagebodyIdParam string, forceParam *bool) error

    // Read a ALBErrorPageBody.
    //
    // @param albErrorpagebodyIdParam ALBErrorPageBody ID (required)
    // @return com.vmware.nsx_policy.model.ALBErrorPageBody
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albErrorpagebodyIdParam string) (model.ALBErrorPageBody, error)

    // Paginated list of all ALBErrorPageBody for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBErrorPageBodyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBErrorPageBodyApiResponse, error)

    // If a ALBerrorpagebody with the alb-errorpagebody-id is not already present, create a new ALBerrorpagebody. If it already exists, update the ALBerrorpagebody. This is a full replace.
    //
    // @param albErrorpagebodyIdParam ALBerrorpagebody ID (required)
    // @param aLBErrorPageBodyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albErrorpagebodyIdParam string, aLBErrorPageBodyParam model.ALBErrorPageBody) error

    // If a ALBErrorPageBody with the alb-ErrorPageBody-id is not already present, create a new ALBErrorPageBody. If it already exists, update the ALBErrorPageBody. This is a full replace.
    //
    // @param albErrorpagebodyIdParam ALBErrorPageBody ID (required)
    // @param aLBErrorPageBodyParam (required)
    // @return com.vmware.nsx_policy.model.ALBErrorPageBody
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albErrorpagebodyIdParam string, aLBErrorPageBodyParam model.ALBErrorPageBody) (model.ALBErrorPageBody, error)
}
