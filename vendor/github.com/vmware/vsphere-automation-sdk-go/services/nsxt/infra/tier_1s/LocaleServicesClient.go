/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: LocaleServices
 * Used by client-side stubs.
 */

package tier_1s

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LocaleServicesClient interface {

    // Delete Tier-1 locale-services
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier1IdParam string, localeServicesIdParam string) error

    // Read Tier-1 locale-services
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @return com.vmware.nsx_policy.model.LocaleServices
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier1IdParam string, localeServicesIdParam string) (model.LocaleServices, error)

    // Paginated list of all Tier-1 locale-services
    //
    // @param tier1IdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.LocaleServicesListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier1IdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.LocaleServicesListResult, error)

    // If a Tier-1 locale services with the locale-services-id is not already present, create a new locale services. If it already exists, update Tier-1 locale services with specified attributes.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param localeServicesParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier1IdParam string, localeServicesIdParam string, localeServicesParam model.LocaleServices) error

    // If a Tier-1 locale services with the locale-services-id is not already present, create a new locale-services. If it already exists, replace the Tier-1 locale services instance with the new object. Modifying edge cluster re-allocates NSX service routers of TIER1 to a different edge cluster. You can also place edge nodes manually and provide maximum two indices for HA mode ACTIVE_STANDBY. To re-allocating edge cluster is a disruptive operation and all existing statistics of Tier1 gateway will be removed.
    //
    // @param tier1IdParam (required)
    // @param localeServicesIdParam (required)
    // @param localeServicesParam (required)
    // @return com.vmware.nsx_policy.model.LocaleServices
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier1IdParam string, localeServicesIdParam string, localeServicesParam model.LocaleServices) (model.LocaleServices, error)
}
