/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ArpTable
 * Used by client-side stubs.
 */

package interfaces

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type ArpTableClient interface {

    // Returns ARP table (IPv4) or Neighbor Discovery table (IPv6) for the tier-1 interface, on a edge node specified in edge_path parameter. The edge_path parameter is mandatory.
    //
    // @param tier1IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param edgePathParam Policy path of edge node (optional)
    // @param enforcementPointPathParam Enforcement point path (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.InterfaceArpTable
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, localeServiceIdParam string, interfaceIdParam string, cursorParam *string, edgePathParam *string, enforcementPointPathParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.InterfaceArpTable, error)
}
