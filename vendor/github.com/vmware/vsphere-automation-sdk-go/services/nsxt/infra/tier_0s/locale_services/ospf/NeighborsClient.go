/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Neighbors
 * Used by client-side stubs.
 */

package ospf

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type NeighborsClient interface {

    // Get OSPF neighbor information.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param edgePathParam Policy path of edge (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param neighborAddressParam IPv4 or IPv6 address (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.OspfNeighborsStatusListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, edgePathParam *string, includedFieldsParam *string, neighborAddressParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.OspfNeighborsStatusListResult, error)
}
