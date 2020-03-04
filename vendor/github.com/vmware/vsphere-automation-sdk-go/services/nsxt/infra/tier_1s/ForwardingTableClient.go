/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ForwardingTable
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ForwardingTableClient interface {

    // Get forwarding table from tier-1
    //
    // @param tier1IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param edgeIdParam UUID of edge node (optional)
    // @param edgePathParam Policy path of edge node (optional)
    // @param enforcementPointPathParam Enforcement point path (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param networkPrefixParam IPv4 or IPv6 CIDR Block (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param routeSourceParam Filter routes based on the source from which route is learned (optional)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.RoutingTableListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, edgeIdParam *string, edgePathParam *string, enforcementPointPathParam *string, includedFieldsParam *string, networkPrefixParam *string, pageSizeParam *int64, routeSourceParam *string, sortAscendingParam *bool, sortByParam *string) (model.RoutingTableListResult, error)
}
