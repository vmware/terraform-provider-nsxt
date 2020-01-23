/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FloodProtectionProfileBindings
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type FloodProtectionProfileBindingsClient interface {

    // API will delete Flood Protection Profile Binding for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, floodProtectionProfileBindingIdParam string) error

    // API will get Flood Protection Profile Binding Map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @return com.vmware.nsx_policy.model.FloodProtectionProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, floodProtectionProfileBindingIdParam string) (model.FloodProtectionProfileBindingMap, error)

    // API will create or update Flood Protection profile binding map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @param floodProtectionProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) error

    // API will create or update Flood Protection profile binding map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @param floodProtectionProfileBindingMapParam (required)
    // @return com.vmware.nsx_policy.model.FloodProtectionProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) (model.FloodProtectionProfileBindingMap, error)
}
