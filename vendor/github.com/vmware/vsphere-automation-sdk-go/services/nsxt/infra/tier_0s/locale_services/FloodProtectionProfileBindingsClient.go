/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FloodProtectionProfileBindings
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type FloodProtectionProfileBindingsClient interface {

    // API will delete Flood Protection Profile Binding for Tier-0 Logical Router LocaleServices.
    //
    // @param tier0IdParam (required)
    // @param localeServicesIdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServicesIdParam string, floodProtectionProfileBindingIdParam string) error

    // API will get Flood Protection Profile Binding Map for Tier-0 Logical Router LocaleServices.
    //
    // @param tier0IdParam (required)
    // @param localeServicesIdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @return com.vmware.nsx_policy.model.FloodProtectionProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServicesIdParam string, floodProtectionProfileBindingIdParam string) (model.FloodProtectionProfileBindingMap, error)

    // API will create or update Flood Protection profile binding map for Tier-0 Logical Router LocaleServices.
    //
    // @param tier0IdParam (required)
    // @param localeServicesIdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @param floodProtectionProfileBindingMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServicesIdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) error

    // API will create or update Flood Protection profile binding map for Tier-0 Logical Router LocaleServices.
    //
    // @param tier0IdParam (required)
    // @param localeServicesIdParam (required)
    // @param floodProtectionProfileBindingIdParam (required)
    // @param floodProtectionProfileBindingMapParam (required)
    // @return com.vmware.nsx_policy.model.FloodProtectionProfileBindingMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServicesIdParam string, floodProtectionProfileBindingIdParam string, floodProtectionProfileBindingMapParam model.FloodProtectionProfileBindingMap) (model.FloodProtectionProfileBindingMap, error)
}
