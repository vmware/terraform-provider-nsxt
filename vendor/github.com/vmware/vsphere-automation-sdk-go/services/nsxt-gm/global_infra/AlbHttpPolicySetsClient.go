/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbHttpPolicySets
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbHttpPolicySetsClient interface {

    // Delete the ALBHTTPPolicySet along with all the entities contained by this ALBHTTPPolicySet.
    //
    // @param albHttppolicysetIdParam ALBHTTPPolicySet ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albHttppolicysetIdParam string, forceParam *bool) error

    // Read a ALBHTTPPolicySet.
    //
    // @param albHttppolicysetIdParam ALBHTTPPolicySet ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBHTTPPolicySet
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albHttppolicysetIdParam string) (model.ALBHTTPPolicySet, error)

    // Paginated list of all ALBHTTPPolicySet for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBHTTPPolicySetApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBHTTPPolicySetApiResponse, error)

    // If a ALBhttppolicyset with the alb-httppolicyset-id is not already present, create a new ALBhttppolicyset. If it already exists, update the ALBhttppolicyset. This is a full replace.
    //
    // @param albHttppolicysetIdParam ALBhttppolicyset ID (required)
    // @param aLBHTTPPolicySetParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albHttppolicysetIdParam string, aLBHTTPPolicySetParam model.ALBHTTPPolicySet) error

    // If a ALBHTTPPolicySet with the alb-HTTPPolicySet-id is not already present, create a new ALBHTTPPolicySet. If it already exists, update the ALBHTTPPolicySet. This is a full replace.
    //
    // @param albHttppolicysetIdParam ALBHTTPPolicySet ID (required)
    // @param aLBHTTPPolicySetParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBHTTPPolicySet
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albHttppolicysetIdParam string, aLBHTTPPolicySetParam model.ALBHTTPPolicySet) (model.ALBHTTPPolicySet, error)
}
