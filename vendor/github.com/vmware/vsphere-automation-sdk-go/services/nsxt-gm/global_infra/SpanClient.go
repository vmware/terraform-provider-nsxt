/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Span
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type SpanClient interface {

    // Get span for an entity with specified path.
    //
    // @param intentPathParam String Path of the intent object (required)
    // @param sitePathParam Policy Path of the site (optional)
    // @return com.vmware.nsx_global_policy.model.Span
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(intentPathParam string, sitePathParam *string) (model.Span, error)
}
