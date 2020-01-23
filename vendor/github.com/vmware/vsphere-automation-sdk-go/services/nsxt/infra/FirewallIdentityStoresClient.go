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
)

type FirewallIdentityStoresClient interface {

    // If the firewall identity store is removed, it will stop the identity store synchronization. User will not be able to define new IDFW rules
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(firewallIdentityStoreIdParam string) error

    // Return a firewall identity store based on the store identifier
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @return com.vmware.nsx_policy.model.FirewallIdentityStore
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(firewallIdentityStoreIdParam string) (model.FirewallIdentityStore, error)

    // List all firewall identity stores
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.FirewallIdentityStoreListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.FirewallIdentityStoreListResult, error)

    // If a firewall identity store with the firewall-identity-store-id is not already present, create a new firewall identity store. If it already exists, update the firewall identity store with specified attributes.
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param firewallIdentityStoreParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(firewallIdentityStoreIdParam string, firewallIdentityStoreParam model.FirewallIdentityStore) error

    // If a firewall identity store with the firewall-identity-store-id is not already present, create a new firewall identity store. If it already exists, replace the firewall identity store instance with the new object.
    //
    // @param firewallIdentityStoreIdParam firewall identity store ID (required)
    // @param firewallIdentityStoreParam (required)
    // @return com.vmware.nsx_policy.model.FirewallIdentityStore
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(firewallIdentityStoreIdParam string, firewallIdentityStoreParam model.FirewallIdentityStore) (model.FirewallIdentityStore, error)
}
