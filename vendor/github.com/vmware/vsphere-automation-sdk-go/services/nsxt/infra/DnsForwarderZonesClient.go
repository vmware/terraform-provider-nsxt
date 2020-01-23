/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DnsForwarderZones
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type DnsForwarderZonesClient interface {

    // Delete the DNS Forwarder Zone
    //
    // @param dnsForwarderZoneIdParam DNS Forwarder Zone ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(dnsForwarderZoneIdParam string) error

    // Read the DNS Forwarder Zone
    //
    // @param dnsForwarderZoneIdParam DNS Forwarder Zone ID (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarderZone
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(dnsForwarderZoneIdParam string) (model.PolicyDnsForwarderZone, error)

    // Paginated list of all Dns Forwarder Zones
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarderZoneListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyDnsForwarderZoneListResult, error)

    // Create or update the DNS Forwarder Zone
    //
    // @param dnsForwarderZoneIdParam DNS Forwarder Zone ID (required)
    // @param policyDnsForwarderZoneParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(dnsForwarderZoneIdParam string, policyDnsForwarderZoneParam model.PolicyDnsForwarderZone) error

    // Create or update the DNS Forwarder Zone
    //
    // @param dnsForwarderZoneIdParam DNS Forwarder Zone ID (required)
    // @param policyDnsForwarderZoneParam (required)
    // @return com.vmware.nsx_policy.model.PolicyDnsForwarderZone
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(dnsForwarderZoneIdParam string, policyDnsForwarderZoneParam model.PolicyDnsForwarderZone) (model.PolicyDnsForwarderZone, error)
}
