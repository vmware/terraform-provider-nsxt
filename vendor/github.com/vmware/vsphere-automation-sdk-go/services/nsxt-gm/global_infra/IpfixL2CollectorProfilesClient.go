/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: IpfixL2CollectorProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type IpfixL2CollectorProfilesClient interface {

    // API deletes IPFIX collector profile. Flow forwarding to collector will be stopped.
    //
    // @param ipfixL2CollectorProfileIdParam IPFIX collector Profile id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(ipfixL2CollectorProfileIdParam string) error

    // API will return details of IPFIX collector profile.
    //
    // @param ipfixL2CollectorProfileIdParam IPFIX collector profile id (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2CollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(ipfixL2CollectorProfileIdParam string) (model.IPFIXL2CollectorProfile, error)

    // API will provide list of all IPFIX collector profiles and their details.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2CollectorProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.IPFIXL2CollectorProfileListResult, error)

    // Create a new IPFIX collector profile if the IPFIX collector profile with given id does not already exist. If the IPFIX collector profile with the given id already exists, patch with the existing IPFIX collector profile.
    //
    // @param ipfixL2CollectorProfileIdParam IPFIX collector profile id (required)
    // @param iPFIXL2CollectorProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(ipfixL2CollectorProfileIdParam string, iPFIXL2CollectorProfileParam model.IPFIXL2CollectorProfile) error

    // Create or Replace IPFIX collector profile. IPFIX data will be sent to IPFIX collector.
    //
    // @param ipfixL2CollectorProfileIdParam IPFIX collector profile id (required)
    // @param iPFIXL2CollectorProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.IPFIXL2CollectorProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(ipfixL2CollectorProfileIdParam string, iPFIXL2CollectorProfileParam model.IPFIXL2CollectorProfile) (model.IPFIXL2CollectorProfile, error)
}
