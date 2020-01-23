/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: MetadataProxies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type MetadataProxiesClient interface {

    // API will delete Metadata Proxy Config with ID profile-id
    //
    // @param metadataProxyIdParam Metadata Proxy ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(metadataProxyIdParam string) error

    // API will read Metadata Proxy Config with ID profile-id
    //
    // @param metadataProxyIdParam Metadata Proxy ID (required)
    // @return com.vmware.nsx_policy.model.MetadataProxyConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(metadataProxyIdParam string) (model.MetadataProxyConfig, error)

    // List all L2 Metadata Proxy Configurations
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.MetadataProxyConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.MetadataProxyConfigListResult, error)

    // API will create or update Metadata Proxy Config with ID profile-id
    //
    // @param metadataProxyIdParam Metadata Proxy ID (required)
    // @param metadataProxyConfigParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(metadataProxyIdParam string, metadataProxyConfigParam model.MetadataProxyConfig) error

    // API will create or update Metadata Proxy Config with ID profile-id
    //
    // @param metadataProxyIdParam Metadata Proxy ID (required)
    // @param metadataProxyConfigParam (required)
    // @return com.vmware.nsx_policy.model.MetadataProxyConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(metadataProxyIdParam string, metadataProxyConfigParam model.MetadataProxyConfig) (model.MetadataProxyConfig, error)
}
