/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: OnboardingCheckCompatibility
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type OnboardingCheckCompatibilityClient interface {

    // Create or fully replace a Site under Infra. Revision is optional for creation and required for update.
    //
    // @param siteNodeConnectionInfoParam (required)
    // @return com.vmware.nsx_policy.model.CompatibilityCheckResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(siteNodeConnectionInfoParam model.SiteNodeConnectionInfo) (model.CompatibilityCheckResult, error)
}
