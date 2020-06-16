/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Drafts
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type DraftsClient interface {

    // Delete a manual draft.
    //
    // @param draftIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(draftIdParam string) error

    // Read a draft for a given draft identifier.
    //
    // @param draftIdParam (required)
    // @return com.vmware.nsx_global_policy.model.PolicyDraft
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(draftIdParam string) (model.PolicyDraft, error)

    // List policy drafts.
    //
    // @param autoDraftsParam Fetch list of draft based on is_auto_draft flag (optional)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.PolicyDraftListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(autoDraftsParam *bool, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyDraftListResult, error)

    // Create a new manual draft if the specified draft id does not correspond to an existing draft. Update the manual draft otherwise. Auto draft can not be updated.
    //
    // @param draftIdParam (required)
    // @param policyDraftParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(draftIdParam string, policyDraftParam model.PolicyDraft) error

    // Read a draft and publish it by applying changes onto current configuration.
    //
    // @param draftIdParam (required)
    // @param infraParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Publish(draftIdParam string, infraParam model.Infra) error

    // Create a new manual draft if the specified draft id does not correspond to an existing draft. Update the manual draft otherwise. Auto draft can not be updated.
    //
    // @param draftIdParam (required)
    // @param policyDraftParam (required)
    // @return com.vmware.nsx_global_policy.model.PolicyDraft
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(draftIdParam string, policyDraftParam model.PolicyDraft) (model.PolicyDraft, error)
}
