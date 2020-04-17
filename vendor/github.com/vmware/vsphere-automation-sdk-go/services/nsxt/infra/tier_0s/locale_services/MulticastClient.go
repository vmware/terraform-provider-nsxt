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
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @return com.vmware.nsx_policy.model.PolicyMulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServicesIdParam string) (model.PolicyMulticastConfig, error)

    // Create or update a Tier-0 multicast configuration defining the multicast replication range, the IGMP or a PIM profile. It will update the configuration if there is already one in place.
    //
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @param policyMulticastConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServicesIdParam string, policyMulticastConfigParam model.PolicyMulticastConfig) error

    // Create or update a Tier-0 multicast configuration defining the multicast replication range, the IGMP or a PIM profile. It will update the configuration if there is already one in place.
    //
    // @param tier0IdParam tier0 id (required)
    // @param localeServicesIdParam locale services id (required)
    // @param policyMulticastConfigParam (required)
    // @return com.vmware.nsx_policy.model.PolicyMulticastConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServicesIdParam string, policyMulticastConfigParam model.PolicyMulticastConfig) (model.PolicyMulticastConfig, error)
}
