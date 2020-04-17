/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Evpn
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EvpnClient interface {

    // Read Evpn Configuration.
    //
    // @param tier0IdParam tier0 id (required)
    // @return com.vmware.nsx_policy.model.EvpnConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string) (model.EvpnConfig, error)

    // Create a evpn configuration if it is not already present, otherwise update the evpn configuration.
    //
    // @param tier0IdParam tier0 id (required)
    // @param evpnConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, evpnConfigParam model.EvpnConfig) error

    // Create or update evpn configuration.
    //
    // @param tier0IdParam tier0 id (required)
    // @param evpnConfigParam (required)
    // @return com.vmware.nsx_policy.model.EvpnConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, evpnConfigParam model.EvpnConfig) (model.EvpnConfig, error)
}
