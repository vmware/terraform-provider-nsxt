/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Interfaces
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type InterfacesClient interface {

    // Delete Tier-1 interface
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param interfaceIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, localeServicesIdParam string, interfaceIdParam string) error

    // Read Tier-1 interface
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param interfaceIdParam (required)
    // @return com.vmware.nsx_global_policy.model.Tier1Interface
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServicesIdParam string, interfaceIdParam string) (model.Tier1Interface, error)

    // Paginated list of all Tier-1 interfaces
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.Tier1InterfaceListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, localeServicesIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Tier1InterfaceListResult, error)

    // If an interface with the interface-id is not already present, create a new interface. If it already exists, update the interface for specified attributes.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param interfaceIdParam (required)
    // @param tier1InterfaceParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, localeServicesIdParam string, interfaceIdParam string, tier1InterfaceParam model.Tier1Interface) error

    // If an interface with the interface-id is not already present, create a new interface. If it already exists, replace the interface with this object.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param interfaceIdParam (required)
    // @param tier1InterfaceParam (required)
    // @return com.vmware.nsx_global_policy.model.Tier1Interface
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, localeServicesIdParam string, interfaceIdParam string, tier1InterfaceParam model.Tier1Interface) (model.Tier1Interface, error)
}
