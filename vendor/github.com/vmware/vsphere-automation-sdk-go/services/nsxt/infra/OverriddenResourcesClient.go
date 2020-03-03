/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: OverriddenResources
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type OverriddenResourcesClient interface {

    // List overridden resources
    //
    // @param intentPathParam Global resource path (optional)
    // @param sitePathParam Site path (optional)
    // @return com.vmware.nsx_policy.model.OverriddenResourceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(intentPathParam *string, sitePathParam *string) (model.OverriddenResourceListResult, error)
}
