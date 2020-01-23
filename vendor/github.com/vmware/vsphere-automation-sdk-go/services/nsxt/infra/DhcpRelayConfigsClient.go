/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DhcpRelayConfigs
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type DhcpRelayConfigsClient interface {

    // Delete DHCP relay configuration
    //
    // @param dhcpRelayConfigIdParam DHCP relay config ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(dhcpRelayConfigIdParam string) error

    // Read DHCP relay configuration
    //
    // @param dhcpRelayConfigIdParam DHCP relay config ID (required)
    // @return com.vmware.nsx_policy.model.DhcpRelayConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(dhcpRelayConfigIdParam string) (model.DhcpRelayConfig, error)

    // Paginated list of all DHCP relay config instances
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.DhcpRelayConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.DhcpRelayConfigListResult, error)

    // If DHCP relay config with the dhcp-relay-config-id is not already present, create a new DHCP relay config instance. If it already exists, update the DHCP relay config instance with specified attributes.
    //
    // @param dhcpRelayConfigIdParam DHCP relay config ID (required)
    // @param dhcpRelayConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(dhcpRelayConfigIdParam string, dhcpRelayConfigParam model.DhcpRelayConfig) error

    // If DHCP relay config with the dhcp-relay-config-id is not already present, create a new DHCP relay config instance. If it already exists, replace the DHCP relay config instance with this object.
    //
    // @param dhcpRelayConfigIdParam DHCP relay config ID (required)
    // @param dhcpRelayConfigParam (required)
    // @return com.vmware.nsx_policy.model.DhcpRelayConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(dhcpRelayConfigIdParam string, dhcpRelayConfigParam model.DhcpRelayConfig) (model.DhcpRelayConfig, error)
}
