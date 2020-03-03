/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ServiceInterfaces
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type ServiceInterfacesClient interface {

    // Delete service interface
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param overrideParam Locally override the global object (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, interfaceIdParam string, overrideParam *bool) error

    // Read service interface
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @return com.vmware.nsx_policy.model.ServiceInterface
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, interfaceIdParam string) (model.ServiceInterface, error)

    // Paginated list of all Service Interfaces
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.ServiceInterfaceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.ServiceInterfaceListResult, error)

    // If an interface with the interface-id is not already present, create a new interface. If it already exists, update the interface for specified attributes.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param serviceInterfaceParam (required)
    // @param overrideParam Locally override the global object (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, interfaceIdParam string, serviceInterfaceParam model.ServiceInterface, overrideParam *bool) error

    // If an interface with the interface-id is not already present, create a new interface. Modification of service interface is not allowed.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param interfaceIdParam (required)
    // @param serviceInterfaceParam (required)
    // @param overrideParam Locally override the global object (optional, default to false)
    // @return com.vmware.nsx_policy.model.ServiceInterface
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, interfaceIdParam string, serviceInterfaceParam model.ServiceInterface, overrideParam *bool) (model.ServiceInterface, error)
}
