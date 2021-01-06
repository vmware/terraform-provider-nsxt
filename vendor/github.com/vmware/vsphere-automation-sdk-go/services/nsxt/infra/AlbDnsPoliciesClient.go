/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbDnsPolicies
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AlbDnsPoliciesClient interface {

    // Delete the ALBDnsPolicy along with all the entities contained by this ALBDnsPolicy.
    //
    // @param albDnspolicyIdParam ALBDnsPolicy ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albDnspolicyIdParam string, forceParam *bool) error

    // Read a ALBDnsPolicy.
    //
    // @param albDnspolicyIdParam ALBDnsPolicy ID (required)
    // @return com.vmware.nsx_policy.model.ALBDnsPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albDnspolicyIdParam string) (model.ALBDnsPolicy, error)

    // Paginated list of all ALBDnsPolicy for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ALBDnsPolicyApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBDnsPolicyApiResponse, error)

    // If a ALBdnspolicy with the alb-dnspolicy-id is not already present, create a new ALBdnspolicy. If it already exists, update the ALBdnspolicy. This is a full replace.
    //
    // @param albDnspolicyIdParam ALBdnspolicy ID (required)
    // @param aLBDnsPolicyParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albDnspolicyIdParam string, aLBDnsPolicyParam model.ALBDnsPolicy) error

    // If a ALBDnsPolicy with the alb-DnsPolicy-id is not already present, create a new ALBDnsPolicy. If it already exists, update the ALBDnsPolicy. This is a full replace.
    //
    // @param albDnspolicyIdParam ALBDnsPolicy ID (required)
    // @param aLBDnsPolicyParam (required)
    // @return com.vmware.nsx_policy.model.ALBDnsPolicy
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albDnspolicyIdParam string, aLBDnsPolicyParam model.ALBDnsPolicy) (model.ALBDnsPolicy, error)
}
