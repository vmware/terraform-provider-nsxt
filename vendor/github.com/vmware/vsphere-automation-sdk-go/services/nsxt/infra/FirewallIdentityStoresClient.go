/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FirewallIdentityStores
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type FirewallIdentityStoresClient interface {

    // Invoke full sync or delta sync for a specific domain, with additional delay in seconds if needed. Stop sync will try to stop any pending sync if any to return to idle state.
    //
    // @param firewallIdentityStoreIdParam Firewall identity store identifier (required)
    // @param actionParam Sync type requested (required)
    // @param delayParam Request to execute the sync with some delay in seconds (optional, default to 0)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(firewallIdentityStoreIdParam string, actionParam string, delayParam *int64, enforcementPointPathParam *string) error

    // If the firewall identity store is removed, it will stop the identity store synchronization. User will not be able to define new IDFW rules
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(firewallIdentityStoreIdParam string, enforcementPointPathParam *string) error

    // Return a firewall identity store based on the store identifier
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_policy.model.DirectoryDomain
    // The return value will contain all the properties defined in model.DirectoryDomain.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(firewallIdentityStoreIdParam string, enforcementPointPathParam *string) (*data.StructValue, error)

    // List all firewall identity stores
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.DirectoryDomainListResults
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, enforcementPointPathParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.DirectoryDomainListResults, error)

    // If a firewall identity store with the firewall-identity-store-id is not already present, create a new firewall identity store. If it already exists, update the firewall identity store with specified attributes.
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param directoryDomainParam (required)
    // The parameter must contain all the properties defined in model.DirectoryDomain.
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(firewallIdentityStoreIdParam string, directoryDomainParam *data.StructValue, enforcementPointPathParam *string) error

    // If a firewall identity store with the firewall-identity-store-id is not already present, create a new firewall identity store. If it already exists, replace the firewall identity store instance with the new object.
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param directoryDomainParam (required)
    // The parameter must contain all the properties defined in model.DirectoryDomain.
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @return com.vmware.nsx_policy.model.DirectoryDomain
    // The return value will contain all the properties defined in model.DirectoryDomain.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(firewallIdentityStoreIdParam string, directoryDomainParam *data.StructValue, enforcementPointPathParam *string) (*data.StructValue, error)
}
