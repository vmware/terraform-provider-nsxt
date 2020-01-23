/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: EdgeClusters
 * Used by client-side stubs.
 */

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type EdgeClustersClient interface {

    // Read a Edge Cluster under an Enforcement Point
    //
    // @param siteIdParam (required)
    // @param enforcementpointIdParam (required)
    // @param edgeClusterIdParam (required)
    // @return com.vmware.nsx_policy.model.PolicyEdgeCluster
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(siteIdParam string, enforcementpointIdParam string, edgeClusterIdParam string) (model.PolicyEdgeCluster, error)

    // Paginated list of all Edge Clusters under an Enforcement Point
    //
    // @param siteIdParam (required)
    // @param enforcementpointIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyEdgeClusterListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(siteIdParam string, enforcementpointIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyEdgeClusterListResult, error)
}
