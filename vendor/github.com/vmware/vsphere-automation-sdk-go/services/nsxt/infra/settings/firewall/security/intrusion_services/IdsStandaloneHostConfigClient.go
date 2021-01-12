/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IdsStandaloneHostConfig
 * Used by client-side stubs.
 */

package intrusion_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IdsStandaloneHostConfigClient interface {

    // Read intrusion detection system config of standalone hosts.
    // @return com.vmware.nsx_policy.model.IdsStandaloneHostConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get() (model.IdsStandaloneHostConfig, error)

    // Patch intrusion detection system configuration on standalone hosts.
    //
    // @param idsStandaloneHostConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(idsStandaloneHostConfigParam model.IdsStandaloneHostConfig) error

    // Update intrusion detection system configuration on standalone hosts.
    //
    // @param idsStandaloneHostConfigParam (required)
    // @return com.vmware.nsx_policy.model.IdsStandaloneHostConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(idsStandaloneHostConfigParam model.IdsStandaloneHostConfig) (model.IdsStandaloneHostConfig, error)
}
