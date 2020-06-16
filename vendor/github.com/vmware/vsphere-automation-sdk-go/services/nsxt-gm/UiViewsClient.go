/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: UiViews
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type UiViewsClient interface {

    // Creates a new View.
    //
    // @param viewParam (required)
    // @return com.vmware.nsx_global_policy.model.View
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(viewParam model.View) (model.View, error)

    // Delete View
    //
    // @param viewIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(viewIdParam string) error

    // If no query params are specified then all the views entitled for the user are returned. The views to which a user is entitled to include the views created by the user and the shared views.
    //
    // @param tagParam The tag for which associated views to be queried. (optional)
    // @param viewIdsParam Ids of the Views (optional)
    // @param widgetIdParam Id of widget configuration (optional)
    // @return com.vmware.nsx_global_policy.model.ViewList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tagParam *string, viewIdsParam *string, widgetIdParam *string) (model.ViewList, error)

    // Returns Information about a specific View.
    //
    // @param viewIdParam (required)
    // @return com.vmware.nsx_global_policy.model.View
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get0(viewIdParam string) (model.View, error)

    // Update View
    //
    // @param viewIdParam (required)
    // @param viewParam (required)
    // @return com.vmware.nsx_global_policy.model.View
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(viewIdParam string, viewParam model.View) (model.View, error)
}
