/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbPriorityLabels
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbPriorityLabelsClient interface {

    // Delete the ALBPriorityLabels along with all the entities contained by this ALBPriorityLabels.
    //
    // @param albPrioritylabelsIdParam ALBPriorityLabels ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albPrioritylabelsIdParam string, forceParam *bool) error

    // Read a ALBPriorityLabels.
    //
    // @param albPrioritylabelsIdParam ALBPriorityLabels ID (required)
    // @return com.vmware.nsx_policy.model.ALBPriorityLabels
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albPrioritylabelsIdParam string) (model.ALBPriorityLabels, error)

    // Paginated list of all ALBPriorityLabels for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBPriorityLabelsApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPriorityLabelsApiResponse, error)

    // If a ALBprioritylabels with the alb-prioritylabels-id is not already present, create a new ALBprioritylabels. If it already exists, update the ALBprioritylabels. This is a full replace.
    //
    // @param albPrioritylabelsIdParam ALBprioritylabels ID (required)
    // @param aLBPriorityLabelsParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albPrioritylabelsIdParam string, aLBPriorityLabelsParam model.ALBPriorityLabels) error

    // If a ALBPriorityLabels with the alb-PriorityLabels-id is not already present, create a new ALBPriorityLabels. If it already exists, update the ALBPriorityLabels. This is a full replace.
    //
    // @param albPrioritylabelsIdParam ALBPriorityLabels ID (required)
    // @param aLBPriorityLabelsParam (required)
    // @return com.vmware.nsx_policy.model.ALBPriorityLabels
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albPrioritylabelsIdParam string, aLBPriorityLabelsParam model.ALBPriorityLabels) (model.ALBPriorityLabels, error)
}
