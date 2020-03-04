/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Multicast
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type MulticastClient interface {

    // Read Multicast Config.
    //
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @return com.vmware.nsx_policy.model.PolicyMulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServicesIdParam string) (model.PolicyMulticastConfig, error)

    // Create a multicast config with the multicast-config-id is not already present, otherwise update the multicast config.
    //
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @param policyMulticastConfigParam (required)
    // @param overrideParam Locally override the global object (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServicesIdParam string, policyMulticastConfigParam model.PolicyMulticastConfig, overrideParam *bool) error

    // Create or update multicast config.
    //
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @param policyMulticastConfigParam (required)
    // @param overrideParam Locally override the global object (optional, default to false)
    // @return com.vmware.nsx_policy.model.PolicyMulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServicesIdParam string, policyMulticastConfigParam model.PolicyMulticastConfig, overrideParam *bool) (model.PolicyMulticastConfig, error)
}
