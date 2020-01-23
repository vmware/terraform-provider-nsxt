/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DnsForwarder
 * Used by client-side stubs.
 */

package tier_0s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type DnsForwarderClient interface {

    // Perform the specified action for Tier0 DNS forwarder on specified enforcement point.
    //
    // @param tier0IdParam (required)
    // @param actionParam An action to be performed for DNS forwarder on EP (required)
    // @param enforcementPointPathParam An enforcement point path, on which the action is to be performed (optional, default to /infra/sites/default/enforcement-points/default)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(tier0IdParam string, actionParam string, enforcementPointPathParam *string) error

    // Delete DNS configuration for tier-0 instance
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string) error

    // Read the DNS Forwarder for the given tier-0 instance
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarder
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string) (model.PolicyDnsForwarder, error)

    // Update the DNS Forwarder
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @param policyDnsForwarderParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, policyDnsForwarderParam model.PolicyDnsForwarder) error

    // Update the DNS Forwarder
    //
    // @param tier0IdParam Tier-0 ID (required)
    // @param policyDnsForwarderParam (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarder
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, policyDnsForwarderParam model.PolicyDnsForwarder) (model.PolicyDnsForwarder, error)
}
