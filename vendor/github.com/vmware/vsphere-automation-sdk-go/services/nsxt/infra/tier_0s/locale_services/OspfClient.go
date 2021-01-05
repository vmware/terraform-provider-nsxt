/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Ospf
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type OspfClient interface {

    // Read OSPF routing config
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @return com.vmware.nsx_policy.model.OspfRoutingConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string) (model.OspfRoutingConfig, error)

    // If OSPF routing config is not already present, create OSPF routing config. If it already exists, replace the OSPF routing config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param ospfRoutingConfigParam (required)
    // @return com.vmware.nsx_policy.model.OspfRoutingConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, ospfRoutingConfigParam model.OspfRoutingConfig) (model.OspfRoutingConfig, error)

    // If OSPF routing config is not already present, create OSPF routing config. If it already exists, replace the OSPF routing config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param ospfRoutingConfigParam (required)
    // @return com.vmware.nsx_policy.model.OspfRoutingConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, ospfRoutingConfigParam model.OspfRoutingConfig) (model.OspfRoutingConfig, error)
}
