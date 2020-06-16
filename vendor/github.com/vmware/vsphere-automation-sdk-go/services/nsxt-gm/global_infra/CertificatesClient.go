/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Certificates
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type CertificatesClient interface {

    // Removes the specified certificate. The private key associated with the certificate is also deleted.
    //
    // @param certificateIdParam ID of certificate to delete (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(certificateIdParam string) error

    // Returns information for the specified certificate ID, including the certificate's id; resource_type (for example, certificate_self_signed, certificate_ca, or certificate_signed); pem_encoded data; and history of the certificate (who created or modified it and when). For additional information, include the ?details=true modifier at the end of the request URI.
    //
    // @param certificateIdParam ID of certificate to read (required)
    // @param detailsParam whether to expand the pem data and show all its details (optional, default to false)
    // @return com.vmware.nsx_global_policy.model.TlsCertificate
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(certificateIdParam string, detailsParam *bool) (model.TlsCertificate, error)

    // Returns all certificate information viewable by the user, including each certificate's id; resource_type (for example, certificate_self_signed, certificate_ca, or certificate_signed); pem_encoded data; and history of the certificate (who created or modified it and when). For additional information, include the ?details=true modifier at the end of the request URI.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param detailsParam whether to expand the pem data and show all its details (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @param type_Param Type of certificate to return (optional)
    // @return com.vmware.nsx_global_policy.model.TlsCertificateList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, detailsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string, type_Param *string) (model.TlsCertificateList, error)

    // Adds a new private-public certificate and, optionally, a private key that can be applied to one of the user-facing components (appliance management or edge). The certificate and the key should be stored in PEM format. If no private key is provided, the certificate is used as a client certificate in the trust store. A certificate chain will not be expanded into separate certificate instances for reference, but would be pushed to the enforcement point as a single certificate. This patch method does not modify an existing certificate.
    //
    // @param certificateIdParam (required)
    // @param tlsTrustDataParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(certificateIdParam string, tlsTrustDataParam model.TlsTrustData) error

    // Adds a new private-public certificate and, optionally, a private key that can be applied to one of the user-facing components (appliance management or edge). The certificate and the key should be stored in PEM format. If no private key is provided, the certificate is used as a client certificate in the trust store. A certificate chain will not be expanded into separate certificate instances for reference, but would be pushed to the enforcement point as a single certificate.
    //
    // @param certificateIdParam (required)
    // @param tlsTrustDataParam (required)
    // @return com.vmware.nsx_global_policy.model.TlsCertificate
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(certificateIdParam string, tlsTrustDataParam model.TlsTrustData) (model.TlsCertificate, error)
}
