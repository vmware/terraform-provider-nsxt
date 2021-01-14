/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ClusterConfigs
 * Used by client-side stubs.
 */

package intrusion_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ClusterConfigsClient interface {

    // Read intrusion detection system cluster config
    //
    // @param clusterIdParam User entered ID (required)
    // @return com.vmware.nsx_policy.model.IdsClusterConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(clusterIdParam string) (model.IdsClusterConfig, error)

    // List intrusion detection system cluster configs.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IdsClusterConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IdsClusterConfigListResult, error)

    // Patch intrusion detection system on cluster level.
    //
    // @param clusterIdParam User entered ID (required)
    // @param idsClusterConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(clusterIdParam string, idsClusterConfigParam model.IdsClusterConfig) error

    // Update intrusion detection system on cluster level.
    //
    // @param clusterIdParam User entered ID (required)
    // @param idsClusterConfigParam (required)
    // @return com.vmware.nsx_policy.model.IdsClusterConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(clusterIdParam string, idsClusterConfigParam model.IdsClusterConfig) (model.IdsClusterConfig, error)
}
