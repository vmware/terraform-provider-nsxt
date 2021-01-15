/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GlobalSignatures
 * Used by client-side stubs.
 */

package intrusion_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GlobalSignaturesClient interface {

    // Delete global intrusion detection signature.
    //
    // @param signatureIdParam Signature ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(signatureIdParam string) error

    // Read global intrusion detection signature
    //
    // @param signatureIdParam Signature ID (required)
    // @return com.vmware.nsx_policy.model.GlobalIdsSignature
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(signatureIdParam string) (model.GlobalIdsSignature, error)

    // List global intrusion detection signatures.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.GlobalIdsSignatureListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GlobalIdsSignatureListResult, error)

    // Patch global intrusion detection system signature.
    //
    // @param signatureIdParam Signature ID (required)
    // @param globalIdsSignatureParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(signatureIdParam string, globalIdsSignatureParam model.GlobalIdsSignature) error

    // Update global intrusion detection signature.
    //
    // @param signatureIdParam Signature ID (required)
    // @param globalIdsSignatureParam (required)
    // @return com.vmware.nsx_policy.model.GlobalIdsSignature
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(signatureIdParam string, globalIdsSignatureParam model.GlobalIdsSignature) (model.GlobalIdsSignature, error)
}
