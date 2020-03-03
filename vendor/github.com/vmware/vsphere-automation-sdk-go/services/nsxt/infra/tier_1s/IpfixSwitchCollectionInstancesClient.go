/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixSwitchCollectionInstances
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpfixSwitchCollectionInstancesClient interface {

    // API deletes IPFIX Switch Collection Instance.Flow forwarding to selected collector will be stopped. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) error

    // API will return details of IPFIX switch collection. If instance does not exist, it will return 404. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param ipfixSwitchCollectionInstanceIdParam IPFIX switch collection id (required)
    // @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstance
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string) (model.IPFIXSwitchCollectionInstance, error)

    // API provides list IPFIX Switch collection instances available on selected logical switch. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstanceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXSwitchCollectionInstanceListResult, error)

    // Create a new IPFIX switch collection instance if the IPFIX switch collection instance with given id does not already exist. If the IPFIX switch collection instance with the given id already exists, patch with the existing IPFIX switch collection instance. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
    // @param iPFIXSwitchCollectionInstanceParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) error

    // Create or replace IPFIX switch collection instance. Instance will start forwarding data to provided IPFIX collector. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-profiles
    //
    // @param tier1IdParam Tier-1 ID (required)
    // @param ipfixSwitchCollectionInstanceIdParam IPFIX Switch Collection Instance ID (required)
    // @param iPFIXSwitchCollectionInstanceParam (required)
    // @return com.vmware.nsx_policy.model.IPFIXSwitchCollectionInstance
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, ipfixSwitchCollectionInstanceIdParam string, iPFIXSwitchCollectionInstanceParam model.IPFIXSwitchCollectionInstance) (model.IPFIXSwitchCollectionInstance, error)
}
