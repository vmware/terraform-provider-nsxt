/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EvpnTenantConfigs
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type EvpnTenantConfigsClient interface {

    // Create a global evpn tenant configuration if it is not already present, otherwise update the evpn tenant configuration.
    //
    // @param configIdParam Evpn Tenant config id (required)
    // @param evpnTenantConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(configIdParam string, evpnTenantConfigParam model.EvpnTenantConfig) error

    // Create or update Evpn Tenant configuration.
    //
    // @param configIdParam Evpn Tenant config id (required)
    // @param evpnTenantConfigParam (required)
    // @return com.vmware.nsx_global_policy.model.EvpnTenantConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(configIdParam string, evpnTenantConfigParam model.EvpnTenantConfig) (model.EvpnTenantConfig, error)
}
