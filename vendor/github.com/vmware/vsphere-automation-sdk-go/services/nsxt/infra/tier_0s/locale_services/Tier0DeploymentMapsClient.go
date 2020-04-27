/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Tier0DeploymentMaps
 * Used by client-side stubs.
 */

package locale_services

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type Tier0DeploymentMapsClient interface {

    // Delete Tier-0 Deployment Map
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param tier0DeploymentMapIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, tier0DeploymentMapIdParam string) error

    // Read a Tier-0 Deployment Map
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param tier0DeploymentMapIdParam (required)
    // @return com.vmware.nsx_policy.model.Tier0DeploymentMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, tier0DeploymentMapIdParam string) (model.Tier0DeploymentMap, error)

    // Paginated list of all Tier-0 Deployment Entries.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.Tier0DeploymentMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.Tier0DeploymentMapListResult, error)

    // If the passed Tier-0 Deployment Map does not already exist, create a new Tier-0 Deployment Map. If it already exists, patch it.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param tier0DeploymentMapIdParam (required)
    // @param tier0DeploymentMapParam (required)
    // @return com.vmware.nsx_policy.model.Tier0DeploymentMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, tier0DeploymentMapIdParam string, tier0DeploymentMapParam model.Tier0DeploymentMap) (model.Tier0DeploymentMap, error)

    // If the passed Tier-0 Deployment Map does not already exist, create a new Tier-0 Deployment Map. If it already exists, replace it.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param tier0DeploymentMapIdParam (required)
    // @param tier0DeploymentMapParam (required)
    // @return com.vmware.nsx_policy.model.Tier0DeploymentMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, tier0DeploymentMapIdParam string, tier0DeploymentMapParam model.Tier0DeploymentMap) (model.Tier0DeploymentMap, error)
}
