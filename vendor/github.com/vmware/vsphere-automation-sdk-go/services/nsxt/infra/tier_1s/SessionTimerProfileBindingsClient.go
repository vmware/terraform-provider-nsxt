/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: SessionTimerProfileBindings
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type SessionTimerProfileBindingsClient interface {

    // API will delete Session Timer Profile Binding for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param sessionTimerProfileBindingIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, sessionTimerProfileBindingIdParam string) error

    // API will get Session Timer Profile Binding Map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param sessionTimerProfileBindingIdParam (required)
    // @return com.vmware.nsx_policy.model.SessionTimerProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, sessionTimerProfileBindingIdParam string) (model.SessionTimerProfileBindingMap, error)

    // API will create or update Session Timer profile binding map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param sessionTimerProfileBindingIdParam (required)
    // @param sessionTimerProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, sessionTimerProfileBindingIdParam string, sessionTimerProfileBindingMapParam model.SessionTimerProfileBindingMap) error

    // API will create or update Session Timer profile binding map for Tier-1 Logical Router.
    //
    // @param tier1IdParam (required)
    // @param sessionTimerProfileBindingIdParam (required)
    // @param sessionTimerProfileBindingMapParam (required)
    // @return com.vmware.nsx_policy.model.SessionTimerProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, sessionTimerProfileBindingIdParam string, sessionTimerProfileBindingMapParam model.SessionTimerProfileBindingMap) (model.SessionTimerProfileBindingMap, error)
}
