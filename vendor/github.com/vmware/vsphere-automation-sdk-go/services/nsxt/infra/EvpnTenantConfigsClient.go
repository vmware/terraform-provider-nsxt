/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EvpnTenantConfigs
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EvpnTenantConfigsClient interface {

    // Delete evpn tunnel endpoint configuration.
    //
    // @param configIdParam tier0 id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(configIdParam string) error

    // Read Evpn Tenant Configuration.
    //
    // @param configIdParam config id (required)
    // @return com.vmware.nsx_policy.model.EvpnTenantConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(configIdParam string) (model.EvpnTenantConfig, error)

    // List all evpn tunnel endpoint configuration.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.EvpnTenantConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.EvpnTenantConfigListResult, error)

    // Create a global evpn tenant configuration if it is not already present, otherwise update the evpn tenant configuration.
    //
    // @param configIdParam Evpn Tenant config id (required)
    // @param evpnTenantConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(configIdParam string, evpnTenantConfigParam model.EvpnTenantConfig) error

    // Create or update Evpn Tenant configuration.
    //
    // @param configIdParam Evpn Tenant config id (required)
    // @param evpnTenantConfigParam (required)
    // @return com.vmware.nsx_policy.model.EvpnTenantConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(configIdParam string, evpnTenantConfigParam model.EvpnTenantConfig) (model.EvpnTenantConfig, error)
}
