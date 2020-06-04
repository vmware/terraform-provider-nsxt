/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GatewayQosProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GatewayQosProfilesClient interface {

    // Delete QoS profile
    //
    // @param qosProfileIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(qosProfileIdParam string) error

    // Read gateway QoS profile
    //
    // @param qosProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.GatewayQosProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(qosProfileIdParam string) (model.GatewayQosProfile, error)

    // Paginated list of all gateway QoS profle instances
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.GatewayQosProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GatewayQosProfileListResult, error)

    // If profile with the qos-profile-id is not already present, create a new gateway QoS profile instance. If it already exists, update the gateway QoS profile instance with specified attributes.
    //
    // @param qosProfileIdParam (required)
    // @param gatewayQosProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(qosProfileIdParam string, gatewayQosProfileParam model.GatewayQosProfile) error

    // If profile with the qos-profile-id is not already present, create a new gateway QoS profile instance. If it already exists, replace the gateway QoS profile instance with this object.
    //
    // @param qosProfileIdParam (required)
    // @param gatewayQosProfileParam (required)
    // @return com.vmware.nsx_policy.model.GatewayQosProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(qosProfileIdParam string, gatewayQosProfileParam model.GatewayQosProfile) (model.GatewayQosProfile, error)
}
