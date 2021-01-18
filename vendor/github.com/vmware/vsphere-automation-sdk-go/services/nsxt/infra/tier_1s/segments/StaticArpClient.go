/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: StaticArp
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type StaticArpClient interface {

    // Delete static ARP config
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, segmentIdParam string) error

    // Read static ARP config
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @return com.vmware.nsx_policy.model.StaticARPConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, segmentIdParam string) (model.StaticARPConfig, error)

    // Create static ARP config with Tier-1 and segment IDs provided if it doesn't exist, update with provided config if it's already created.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param staticARPConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, segmentIdParam string, staticARPConfigParam model.StaticARPConfig) error

    // Create static ARP config with Tier-1 and segment IDs provided if it doesn't exist, update with provided config if it's already created.
    //
    // @param tier1IdParam (required)
    // @param segmentIdParam (required)
    // @param staticARPConfigParam (required)
    // @return com.vmware.nsx_policy.model.StaticARPConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, segmentIdParam string, staticARPConfigParam model.StaticARPConfig) (model.StaticARPConfig, error)
}
