/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbWebhooks
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbWebhooksClient interface {

    // Delete the ALBWebhook along with all the entities contained by this ALBWebhook.
    //
    // @param albWebhookIdParam ALBWebhook ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albWebhookIdParam string, forceParam *bool) error

    // Read a ALBWebhook.
    //
    // @param albWebhookIdParam ALBWebhook ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBWebhook
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albWebhookIdParam string) (model.ALBWebhook, error)

    // Paginated list of all ALBWebhook for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBWebhookApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWebhookApiResponse, error)

    // If a ALBwebhook with the alb-webhook-id is not already present, create a new ALBwebhook. If it already exists, update the ALBwebhook. This is a full replace.
    //
    // @param albWebhookIdParam ALBwebhook ID (required)
    // @param aLBWebhookParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albWebhookIdParam string, aLBWebhookParam model.ALBWebhook) error

    // If a ALBWebhook with the alb-Webhook-id is not already present, create a new ALBWebhook. If it already exists, update the ALBWebhook. This is a full replace.
    //
    // @param albWebhookIdParam ALBWebhook ID (required)
    // @param aLBWebhookParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBWebhook
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albWebhookIdParam string, aLBWebhookParam model.ALBWebhook) (model.ALBWebhook, error)
}
