/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: AlbSslKeyAndCertificates
 * Used by client-side stubs.
 */

package global_infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type AlbSslKeyAndCertificatesClient interface {

    // Delete the ALBSSLKeyAndCertificate along with all the entities contained by this ALBSSLKeyAndCertificate.
    //
    // @param albSslkeyandcertificateIdParam ALBSSLKeyAndCertificate ID (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(albSslkeyandcertificateIdParam string, forceParam *bool) error

    // Read a ALBSSLKeyAndCertificate.
    //
    // @param albSslkeyandcertificateIdParam ALBSSLKeyAndCertificate ID (required)
    // @return com.vmware.nsx_global_policy.model.ALBSSLKeyAndCertificate
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(albSslkeyandcertificateIdParam string) (model.ALBSSLKeyAndCertificate, error)

    // Paginated list of all ALBSSLKeyAndCertificate for infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.ALBSSLKeyAndCertificateApiResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ALBSSLKeyAndCertificateApiResponse, error)

    // If a ALBsslkeyandcertificate with the alb-sslkeyandcertificate-id is not already present, create a new ALBsslkeyandcertificate. If it already exists, update the ALBsslkeyandcertificate. This is a full replace.
    //
    // @param albSslkeyandcertificateIdParam ALBsslkeyandcertificate ID (required)
    // @param aLBSSLKeyAndCertificateParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(albSslkeyandcertificateIdParam string, aLBSSLKeyAndCertificateParam model.ALBSSLKeyAndCertificate) error

    // If a ALBSSLKeyAndCertificate with the alb-SSLKeyAndCertificate-id is not already present, create a new ALBSSLKeyAndCertificate. If it already exists, update the ALBSSLKeyAndCertificate. This is a full replace.
    //
    // @param albSslkeyandcertificateIdParam ALBSSLKeyAndCertificate ID (required)
    // @param aLBSSLKeyAndCertificateParam (required)
    // @return com.vmware.nsx_global_policy.model.ALBSSLKeyAndCertificate
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(albSslkeyandcertificateIdParam string, aLBSSLKeyAndCertificateParam model.ALBSSLKeyAndCertificate) (model.ALBSSLKeyAndCertificate, error)
}
