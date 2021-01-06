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

    // Read Multicast Configuration.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @return com.vmware.nsx_policy.model.PolicyTier1MulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServicesIdParam string) (model.PolicyTier1MulticastConfig, error)

    // Create or update a Tier-1 multicast configuration defining the multicast replication range. It will update the configuration if there is already one in place.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param policyTier1MulticastConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, localeServicesIdParam string, policyTier1MulticastConfigParam model.PolicyTier1MulticastConfig) error

    // Create or update a Tier-1 multicast configuration defining the multicast replication range. It will update the configuration if there is already one in place.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param policyTier1MulticastConfigParam (required)
    // @return com.vmware.nsx_policy.model.PolicyTier1MulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, localeServicesIdParam string, policyTier1MulticastConfigParam model.PolicyTier1MulticastConfig) (model.PolicyTier1MulticastConfig, error)
}
