/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DhcpServerConfigs
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type DhcpServerConfigsClient interface {

    // Delete DHCP server configuration
    //
    // @param dhcpServerConfigIdParam DHCP server config ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(dhcpServerConfigIdParam string) error

    // Read DHCP server configuration
    //
    // @param dhcpServerConfigIdParam DHCP server config ID (required)
    // @return com.vmware.nsx_global_policy.model.DhcpServerConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(dhcpServerConfigIdParam string) (model.DhcpServerConfig, error)

    // Paginated list of all DHCP server config instances
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.DhcpServerConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.DhcpServerConfigListResult, error)

    // If DHCP server config with the dhcp-server-config-id is not already present, create a new DHCP server config instance. If it already exists, update the DHCP server config instance with specified attributes. Realized entities of this API can be found using the path of Tier-0, Tier1, or Segment where this config is applied on. Modification of edge_cluster_path in DhcpServerConfig will lose all existing DHCP leases. If both the preferred_edge_paths in the DhcpServerConfig are changed in a same PATCH API, e.g. change from [a,b] to [x,y], the current DHCP server leases will be lost, which could cause network connectivity issues. It is recommended to change only one member index in an update call, e.g. from [a, b] to [a,y]. Clearing preferred_edge_paths will not reassign edge nodes from the edge cluster. Instead, the previously-allocated edge nodes will be retained to avoid loss of leases.
    //
    // @param dhcpServerConfigIdParam DHCP server config ID (required)
    // @param dhcpServerConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(dhcpServerConfigIdParam string, dhcpServerConfigParam model.DhcpServerConfig) error

    // If DHCP server config with the dhcp-server-config-id is not already present, create a new DHCP server config instance. If it already exists, replace the DHCP server config instance with this object. Realized entities of this API can be found using the path of Tier-0, Tier1, or Segment where this config is applied on. Modification of edge_cluster_path in DhcpServerConfig will lose all existing DHCP leases. If both the preferred_edge_paths in the DhcpServerConfig are changed in a same PUT API, e.g. change from [a,b] to [x,y], the current DHCP server leases will be lost, which could cause network connectivity issues. It is recommended to change only one member index in an update call, e.g. from [a, b] to [a,y]. Clearing preferred_edge_paths will not reassign edge nodes from the edge cluster. Instead, the previously-allocated edge nodes will be retained to avoid loss of leases.
    //
    // @param dhcpServerConfigIdParam DHCP server config ID (required)
    // @param dhcpServerConfigParam (required)
    // @return com.vmware.nsx_global_policy.model.DhcpServerConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(dhcpServerConfigIdParam string, dhcpServerConfigParam model.DhcpServerConfig) (model.DhcpServerConfig, error)
}
