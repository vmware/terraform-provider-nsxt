/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DnsForwarder
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type DnsForwarderClient interface {

    // Perform the specified action for Tier0 DNS forwarder on specified enforcement point.
    //
    // @param tier1IdParam (required)
    // @param actionParam An action to be performed for DNS forwarder on EP (required)
    // @param enforcementPointPathParam An enforcement point path, on which the action is to be performed (optional, default to /infra/sites/default/enforcement-points/default)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(tier1IdParam string, actionParam string, enforcementPointPathParam *string) error

    // Delete DNS configuration for tier-1 instance
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string) error

    // Read the DNS Forwarder for the given tier-1 instance
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarder
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string) (model.PolicyDnsForwarder, error)

    // Create or update the DNS Forwarder
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param policyDnsForwarderParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, policyDnsForwarderParam model.PolicyDnsForwarder) error

    // Create or update the DNS Forwarder
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param policyDnsForwarderParam (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarder
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, policyDnsForwarderParam model.PolicyDnsForwarder) (model.PolicyDnsForwarder, error)
}
