/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GlobalManagerConfig
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type GlobalManagerConfigClient interface {

    // Create or patch a Global Manager Config
    //
    // @param globalManagerConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(globalManagerConfigParam model.GlobalManagerConfig) error

    // Read a Global Manager config along with sensitive data. For example - rtep_config.ibgp_password
    // @return com.vmware.nsx_global_policy.model.GlobalManagerConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Showsensitivedata() (model.GlobalManagerConfig, error)

    // Create or fully replace a Global Manager Config. Revision is optional for creation and required for update.
    //
    // @param globalManagerConfigParam (required)
    // @return com.vmware.nsx_global_policy.model.GlobalManagerConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(globalManagerConfigParam model.GlobalManagerConfig) (model.GlobalManagerConfig, error)
}
