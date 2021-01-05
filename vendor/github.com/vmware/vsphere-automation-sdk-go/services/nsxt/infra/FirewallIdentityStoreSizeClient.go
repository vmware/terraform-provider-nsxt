/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FirewallIdentityStoreSize
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type FirewallIdentityStoreSizeClient interface {

    // This call scans the size of a directory domain. It may be very | expensive to run this call in some AD domain deployments. Please | use it with caution.
    //
    // @param directoryDomainParam (required)
    // The parameter must contain all the properties defined in model.DirectoryDomain.
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_policy.model.DirectoryDomainSize
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(directoryDomainParam *data.StructValue, enforcementPointPathParam *string) (model.DirectoryDomainSize, error)
}
