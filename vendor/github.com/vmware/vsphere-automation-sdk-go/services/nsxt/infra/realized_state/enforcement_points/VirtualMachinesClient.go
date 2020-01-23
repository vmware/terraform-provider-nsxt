/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: VirtualMachines
 * Used by client-side stubs.
 */

package enforcement_points

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type VirtualMachinesClient interface {

    // This API filters objects of type virtual machines from the specified NSX Manager. This API has been deprecated. Please use the new API GET /infra/realized-state/virtual-machines
    //
    // @param enforcementPointNameParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param dslParam Search DSL (domain specific language) query (optional)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param queryParam Search query (optional)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.SearchResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(enforcementPointNameParam string, cursorParam *string, dslParam *string, includedFieldsParam *string, pageSizeParam *int64, queryParam *string, sortAscendingParam *bool, sortByParam *string) (model.SearchResponse, error)

    // Allows an admin to apply multiple tags to a virtual machine. This operation does not store the intent on the policy side. It applies the tag directly on the specified enforcement point. This operation will replace the existing tags on the virtual machine with the ones that have been passed. If the application of tag fails on the enforcement point, then an error is reported. The admin will have to retry the operation again. Policy framework does not perform a retry. Failure could occur due to multiple reasons. For e.g enforcement point is down, Enforcement point could not apply the tag due to constraints like max tags limit exceeded, etc.
    //
    // @param enforcementPointNameParam (required)
    // @param virtualMachineTagsUpdateParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Updatetags(enforcementPointNameParam string, virtualMachineTagsUpdateParam model.VirtualMachineTagsUpdate) error
}
