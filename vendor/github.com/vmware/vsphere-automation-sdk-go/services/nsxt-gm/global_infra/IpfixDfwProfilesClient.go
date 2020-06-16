/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixDfwProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type IpfixDfwProfilesClient interface {

    // API deletes IPFIX DFW Profile. Selected IPFIX Collectors will stop receiving flows.
    //
    // @param ipfixDfwProfileIdParam IPFIX DFW Profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipfixDfwProfileIdParam string) error

    // API will return details of IPFIX DFW profile.
    //
    // @param ipfixDfwProfileIdParam IPFIX DFW collection id (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXDFWProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipfixDfwProfileIdParam string) (model.IPFIXDFWProfile, error)

    // API provides list IPFIX DFW profiles available on selected logical DFW.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.IPFIXDFWProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXDFWProfileListResult, error)

    // Create a new IPFIX DFW profile if the IPFIX DFW profile with given id does not already exist. If the IPFIX DFW profile with the given id already exists, patch with the existing IPFIX DFW profile.
    //
    // @param ipfixDfwProfileIdParam IPFIX DFW Profile ID (required)
    // @param iPFIXDFWProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipfixDfwProfileIdParam string, iPFIXDFWProfileParam model.IPFIXDFWProfile) error

    // Create or replace IPFIX DFW profile. Config will start forwarding data to provided IPFIX DFW collector.
    //
    // @param ipfixDfwProfileIdParam IPFIX DFW Profile ID (required)
    // @param iPFIXDFWProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXDFWProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipfixDfwProfileIdParam string, iPFIXDFWProfileParam model.IPFIXDFWProfile) (model.IPFIXDFWProfile, error)
}
