/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: PortMirroringProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type PortMirroringProfilesClient interface {

    // API will delete port mirroring profile. Mirroring from source to destination ports will be stopped.
    //
    // @param portMirroringProfileIdParam Port Mirroring Profile Id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(portMirroringProfileIdParam string) error

    // API will return details of port mirroring profile.
    //
    // @param portMirroringProfileIdParam Port Mirroring Profile Id (required)
    // @return com.vmware.nsx_policy.model.PortMirroringProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(portMirroringProfileIdParam string) (model.PortMirroringProfile, error)

    // API will list all port mirroring profiles group.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PortMirroringProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PortMirroringProfileListResult, error)

    // Create a new Port Mirroring Profile if the Port Mirroring Profile with given id does not already exist. If the Port Mirroring Profile with the given id already exists, patch with the existing Port Mirroring Profile. Realized entities of this API can be found using the path of monitoring profile binding map that is used to apply this profile.
    //
    // @param portMirroringProfileIdParam Port Mirroring Profile Id (required)
    // @param portMirroringProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(portMirroringProfileIdParam string, portMirroringProfileParam model.PortMirroringProfile) error

    // Create or Replace port mirroring profile. Packets will be mirrored from source group, segment, port to destination group. Realized entities of this API can be found using the path of monitoring profile binding map that is used to apply this profile.
    //
    // @param portMirroringProfileIdParam Port Mirroring Profiles Id (required)
    // @param portMirroringProfileParam (required)
    // @return com.vmware.nsx_policy.model.PortMirroringProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(portMirroringProfileIdParam string, portMirroringProfileParam model.PortMirroringProfile) (model.PortMirroringProfile, error)
}
