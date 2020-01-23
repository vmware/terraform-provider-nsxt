/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ServiceChains
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ServiceChainsClient interface {

    // This API can be user to delete service chain with given service-chain-id.
    //
    // @param serviceChainIdParam Id of Service chain (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(serviceChainIdParam string) error

    // This API can be used to read service chain with given service-chain-id.
    //
    // @param serviceChainIdParam Id of Service chain (required)
    // @return com.vmware.nsx_policy.model.PolicyServiceChain
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(serviceChainIdParam string) (model.PolicyServiceChain, error)

    // List all the service chains available for service insertion
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyServiceChainListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyServiceChainListResult, error)

    // Create Service chain representing the sequence in which 3rd party services must be consumed.
    //
    // @param serviceChainIdParam Service chain id (required)
    // @param policyServiceChainParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(serviceChainIdParam string, policyServiceChainParam model.PolicyServiceChain) error

    // Create or update Service chain representing the sequence in which 3rd party services must be consumed.
    //
    // @param serviceChainIdParam Service chain id (required)
    // @param policyServiceChainParam (required)
    // @return com.vmware.nsx_policy.model.PolicyServiceChain
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(serviceChainIdParam string, policyServiceChainParam model.PolicyServiceChain) (model.PolicyServiceChain, error)
}
