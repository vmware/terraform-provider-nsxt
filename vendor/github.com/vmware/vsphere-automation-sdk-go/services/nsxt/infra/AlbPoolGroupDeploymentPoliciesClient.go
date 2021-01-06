/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbPoolGroupDeploymentPolicies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbPoolGroupDeploymentPoliciesClient interface {

    // Delete the ALBPoolGroupDeploymentPolicy along with all the entities contained by this ALBPoolGroupDeploymentPolicy.
    //
    // @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albPoolgroupdeploymentpolicyIdParam string, forceParam *bool) error

    // Read a ALBPoolGroupDeploymentPolicy.
    //
    // @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
    // @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albPoolgroupdeploymentpolicyIdParam string) (model.ALBPoolGroupDeploymentPolicy, error)

    // Paginated list of all ALBPoolGroupDeploymentPolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBPoolGroupDeploymentPolicyApiResponse, error)

    // If a ALBpoolgroupdeploymentpolicy with the alb-poolgroupdeploymentpolicy-id is not already present, create a new ALBpoolgroupdeploymentpolicy. If it already exists, update the ALBpoolgroupdeploymentpolicy. This is a full replace.
    //
    // @param albPoolgroupdeploymentpolicyIdParam ALBpoolgroupdeploymentpolicy ID (required)
    // @param aLBPoolGroupDeploymentPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) error

    // If a ALBPoolGroupDeploymentPolicy with the alb-PoolGroupDeploymentPolicy-id is not already present, create a new ALBPoolGroupDeploymentPolicy. If it already exists, update the ALBPoolGroupDeploymentPolicy. This is a full replace.
    //
    // @param albPoolgroupdeploymentpolicyIdParam ALBPoolGroupDeploymentPolicy ID (required)
    // @param aLBPoolGroupDeploymentPolicyParam (required)
    // @return com.vmware.nsx_policy.model.ALBPoolGroupDeploymentPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albPoolgroupdeploymentpolicyIdParam string, aLBPoolGroupDeploymentPolicyParam model.ALBPoolGroupDeploymentPolicy) (model.ALBPoolGroupDeploymentPolicy, error)
}
