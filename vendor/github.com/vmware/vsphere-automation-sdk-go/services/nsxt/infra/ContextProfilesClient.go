/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ContextProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ContextProfilesClient interface {

    // Deletes the specified Policy Context Profile. If the Policy Context Profile is consumed in a firewall rule, it won't get deleted.
    //
    // @param contextProfileIdParam Policy Context Profile Id (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(contextProfileIdParam string, forceParam *bool) error

    // Get a single PolicyContextProfile by id
    //
    // @param contextProfileIdParam (required)
    // @return com.vmware.nsx_policy.model.PolicyContextProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(contextProfileIdParam string) (model.PolicyContextProfile, error)

    // Get all PolicyContextProfiles
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyContextProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyContextProfileListResult, error)

    // Creates/Updates a PolicyContextProfile, which encapsulates attribute and sub-attributes of network services. Rules for using attributes and sub-attributes in single PolicyContextProfile 1. One type of attribute can't have multiple occurrences. ( Eg. - Attribute type APP_ID can be used only once per PolicyContextProfile.) 2. For specifying multiple values for an attribute, provide them in an array. 3. If sub-attribtes are mentioned for an attribute, then only single value is allowed for that attribute. 4. To get a list of supported attributes and sub-attributes fire the following REST API GET https://<policy-mgr>/policy/api/v1/infra/context-profiles/attributes
    //
    // @param contextProfileIdParam (required)
    // @param policyContextProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(contextProfileIdParam string, policyContextProfileParam model.PolicyContextProfile) error

    // Creates/Updates a PolicyContextProfile, which encapsulates attribute and sub-attributes of network services. Rules for using attributes and sub-attributes in single PolicyContextProfile 1. One type of attribute can't have multiple occurrences. ( Eg. - Attribute type APP_ID can be used only once per PolicyContextProfile.) 2. For specifying multiple values for an attribute, provide them in an array. 3. If sub-attribtes are mentioned for an attribute, then only single value is allowed for that attribute. 4. To get a list of supported attributes and sub-attributes fire the following REST API GET https://<policy-mgr>/policy/api/v1/infra/context-profiles/attributes
    //
    // @param contextProfileIdParam (required)
    // @param policyContextProfileParam (required)
    // @return com.vmware.nsx_policy.model.PolicyContextProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(contextProfileIdParam string, policyContextProfileParam model.PolicyContextProfile) (model.PolicyContextProfile, error)
}
