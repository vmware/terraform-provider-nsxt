/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Crls
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type CrlsClient interface {

    // Deletes an existing CRL.
    //
    // @param crlIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(crlIdParam string) error

    // Returns information about the specified CRL. For additional information, include the ?details=true modifier at the end of the request URI.
    //
    // @param crlIdParam (required)
    // @param detailsParam whether to expand the pem data and show all its details (optional, default to false)
    // @return com.vmware.nsx_global_policy.model.TlsCrl
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(crlIdParam string, detailsParam *bool) (model.TlsCrl, error)

    // Adds a new certificate revocation list (CRLs). The CRL is used to verify the client certificate status against the revocation lists published by the CA. For this reason, the administrator needs to add the CRL in certificate repository as well. The CRL can contain a single CRL or multiple CRLs depending on the PEM data. - Single CRL: a single CRL is created with the given id. - Composite CRL: multiple CRLs are generated. Each of the CRL is created with an id generated based on the given id. First CRL is created with crl-id, second with crl-id-1, third with crl-id-2, etc.
    //
    // @param crlIdParam (required)
    // @param tlsCrlParam (required)
    // @return com.vmware.nsx_global_policy.model.TlsCrlListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Importcrl(crlIdParam string, tlsCrlParam model.TlsCrl) (model.TlsCrlListResult, error)

    // Returns information about all CRLs. For additional information, include the ?details=true modifier at the end of the request URI.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param detailsParam whether to expand the pem data and show all its details (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @param type_Param Type of certificate to return (optional)
    // @return com.vmware.nsx_global_policy.model.TlsCrlListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, detailsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, type_Param *string) (model.TlsCrlListResult, error)

    // Create or patch a Certificate Revocation List for the given id. The CRL is used to verify the client certificate status against the revocation lists published by the CA. For this reason, the administrator needs to add the CRL in certificate repository as well. The CRL must contain PEM data for a single CRL.
    //
    // @param crlIdParam (required)
    // @param tlsCrlParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(crlIdParam string, tlsCrlParam model.TlsCrl) error

    // Create or replace a Certificate Revocation List for the given id. The CRL is used to verify the client certificate status against the revocation lists published by the CA. For this reason, the administrator needs to add the CRL in certificate repository as well. The CRL must contain PEM data for a single CRL. Revision is required.
    //
    // @param crlIdParam (required)
    // @param tlsCrlParam (required)
    // @return com.vmware.nsx_global_policy.model.TlsCrl
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(crlIdParam string, tlsCrlParam model.TlsCrl) (model.TlsCrl, error)
}
