/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LbVirtualServers
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbVirtualServersClient interface {

    // Delete the LBVirtualServer along with all the entities contained by this LBVirtualServer.
    //
    // @param lbVirtualServerIdParam LBVirtualServer ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(lbVirtualServerIdParam string, forceParam *bool) error

    // Read a LBVirtualServer.
    //
    // @param lbVirtualServerIdParam LBVirtualServer ID (required)
    // @return com.vmware.nsx_policy.model.LBVirtualServer
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(lbVirtualServerIdParam string) (model.LBVirtualServer, error)

    // Paginated list of all LBVirtualServers.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LBVirtualServerListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LBVirtualServerListResult, error)

    // If a LBVirtualServer with the lb-virtual-server-id is not already present, create a new LBVirtualServer. If it already exists, update the LBVirtualServer. This is a full replace.
    //
    // @param lbVirtualServerIdParam LBVirtualServer ID (required)
    // @param lbVirtualServerParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(lbVirtualServerIdParam string, lbVirtualServerParam model.LBVirtualServer) error

    // If a LBVirtualServer with the lb-virtual-server-id is not already present, create a new LBVirtualServer. If it already exists, update the LBVirtualServer. This is a full replace.
    //
    // @param lbVirtualServerIdParam LBVirtualServer ID (required)
    // @param lbVirtualServerParam (required)
    // @return com.vmware.nsx_policy.model.LBVirtualServer
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(lbVirtualServerIdParam string, lbVirtualServerParam model.LBVirtualServer) (model.LBVirtualServer, error)
}
