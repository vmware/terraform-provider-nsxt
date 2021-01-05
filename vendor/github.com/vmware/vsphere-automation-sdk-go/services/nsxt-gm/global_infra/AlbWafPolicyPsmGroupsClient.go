/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbWafPolicyPsmGroups
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbWafPolicyPsmGroupsClient interface {

    // Delete the ALBWafPolicyPSMGroup along with all the entities contained by this ALBWafPolicyPSMGroup.
    //
    // @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albWafpolicypsmgroupIdParam string, forceParam *bool) error

    // Read a ALBWafPolicyPSMGroup.
    //
    // @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicyPSMGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albWafpolicypsmgroupIdParam string) (model.ALBWafPolicyPSMGroup, error)

    // Paginated list of all ALBWafPolicyPSMGroup for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicyPSMGroupApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBWafPolicyPSMGroupApiResponse, error)

    // If a ALBwafpolicypsmgroup with the alb-wafpolicypsmgroup-id is not already present, create a new ALBwafpolicypsmgroup. If it already exists, update the ALBwafpolicypsmgroup. This is a full replace.
    //
    // @param albWafpolicypsmgroupIdParam ALBwafpolicypsmgroup ID (required)
    // @param aLBWafPolicyPSMGroupParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) error

    // If a ALBWafPolicyPSMGroup with the alb-WafPolicyPSMGroup-id is not already present, create a new ALBWafPolicyPSMGroup. If it already exists, update the ALBWafPolicyPSMGroup. This is a full replace.
    //
    // @param albWafpolicypsmgroupIdParam ALBWafPolicyPSMGroup ID (required)
    // @param aLBWafPolicyPSMGroupParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBWafPolicyPSMGroup
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albWafpolicypsmgroupIdParam string, aLBWafPolicyPSMGroupParam model.ALBWafPolicyPSMGroup) (model.ALBWafPolicyPSMGroup, error)
}
