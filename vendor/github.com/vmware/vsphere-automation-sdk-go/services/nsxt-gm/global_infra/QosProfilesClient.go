/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: QosProfiles
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type QosProfilesClient interface {

    // API will delete QoS profile.
    //
    // @param qosProfileIdParam QoS profile Id (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(qosProfileIdParam string) error

    // API will return details of QoS profile.
    //
    // @param qosProfileIdParam QoS profile Id (required)
    // @return com.vmware.nsx_global_policy.model.QoSProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(qosProfileIdParam string) (model.QosProfile, error)

    // API will list all QoS profiles.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.QoSProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.QosProfileListResult, error)

    // Create a new QoS profile if the QoS profile with given id does not already exist. If the QoS profile with the given id already exists, patch with the existing QoS profile.
    //
    // @param qosProfileIdParam QoS profile Id (required)
    // @param qosProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(qosProfileIdParam string, qosProfileParam model.QosProfile) error

    // Create or Replace QoS profile.
    //
    // @param qosProfileIdParam QoS profile Id (required)
    // @param qosProfileParam (required)
    // @return com.vmware.nsx_global_policy.model.QoSProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(qosProfileIdParam string, qosProfileParam model.QosProfile) (model.QosProfile, error)
}
