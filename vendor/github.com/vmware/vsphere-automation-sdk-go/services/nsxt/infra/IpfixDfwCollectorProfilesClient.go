/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixDfwCollectorProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type IpfixDfwCollectorProfilesClient interface {

    // API deletes IPFIX dfw collector profile. Flow forwarding to collector will be stopped.
    //
    // @param ipfixDfwCollectorProfileIdParam IPFIX dfw collector Profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipfixDfwCollectorProfileIdParam string) error

    // API will return details of IPFIX dfw collector profile. If profile does not exist, it will return 404.
    //
    // @param ipfixDfwCollectorProfileIdParam IPFIX dfw collector profile id (required)
    // @return com.vmware.nsx_policy.model.IPFIXDFWCollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipfixDfwCollectorProfileIdParam string) (model.IPFIXDFWCollectorProfile, error)

    // API will provide list of all IPFIX dfw collector profiles and their details.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.IPFIXDFWCollectorProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXDFWCollectorProfileListResult, error)

    // Create a new IPFIX dfw collector profile if the IPFIX dfw collector profile with given id does not already exist. If the IPFIX dfw collector profile with the given id already exists, patch with the existing IPFIX dfw collector profile.
    //
    // @param ipfixDfwCollectorProfileIdParam (required)
    // @param iPFIXDFWCollectorProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipfixDfwCollectorProfileIdParam string, iPFIXDFWCollectorProfileParam model.IPFIXDFWCollectorProfile) error

    // Create or Replace IPFIX dfw collector profile. IPFIX data will be sent to IPFIX collector port.
    //
    // @param ipfixDfwCollectorProfileIdParam IPFIX dfw collector profile id (required)
    // @param iPFIXDFWCollectorProfileParam (required)
    // @return com.vmware.nsx_policy.model.IPFIXDFWCollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipfixDfwCollectorProfileIdParam string, iPFIXDFWCollectorProfileParam model.IPFIXDFWCollectorProfile) (model.IPFIXDFWCollectorProfile, error)
}
