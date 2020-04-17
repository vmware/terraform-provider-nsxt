/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixCollectorProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpfixCollectorProfilesClient interface {

    // API deletes IPFIX collector profile. Flow forwarding to collector will be stopped. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-collector-profiles
    //
    // @param ipfixCollectorProfileIdParam IPFIX collector Profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipfixCollectorProfileIdParam string) error

    // API will return details of IPFIX collector profile. If profile does not exist, it will return 404. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-collector-profiles
    //
    // @param ipfixCollectorProfileIdParam IPFIX collector profile id (required)
    // @return com.vmware.nsx_policy.model.IPFIXCollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipfixCollectorProfileIdParam string) (model.IPFIXCollectorProfile, error)

    // API will provide list of all IPFIX collector profiles and their details. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-collector-profiles
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPFIXCollectorProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXCollectorProfileListResult, error)

    // Create a new IPFIX collector profile if the IPFIX collector profile with given id does not already exist. If the IPFIX collector profile with the given id already exists, patch with the existing IPFIX collector profile. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-collector-profiles
    //
    // @param ipfixCollectorProfileIdParam IPFIX collector profile id (required)
    // @param iPFIXCollectorProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipfixCollectorProfileIdParam string, iPFIXCollectorProfileParam model.IPFIXCollectorProfile) error

    // Create or Replace IPFIX collector profile. IPFIX data will be sent to IPFIX collector port. This API is deprecated. Please use the following API: https://<policy-mgr>/policy/api/v1/infra/ipfix-l2-collector-profiles
    //
    // @param ipfixCollectorProfileIdParam IPFIX collector profile id (required)
    // @param iPFIXCollectorProfileParam (required)
    // @return com.vmware.nsx_policy.model.IPFIXCollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipfixCollectorProfileIdParam string, iPFIXCollectorProfileParam model.IPFIXCollectorProfile) (model.IPFIXCollectorProfile, error)
}
