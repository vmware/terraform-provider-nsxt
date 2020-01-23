/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DomainDeploymentMaps
 * Used by client-side stubs.
 */

package domains

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type DomainDeploymentMapsClient interface {

    // Delete Domain Deployment Map
    //
    // @param domainIdParam (required)
    // @param domainDeploymentMapIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(domainIdParam string, domainDeploymentMapIdParam string) error

    // Read a Domain Deployment Map
    //
    // @param domainIdParam (required)
    // @param domainDeploymentMapIdParam (required)
    // @return com.vmware.nsx_policy.model.DomainDeploymentMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(domainIdParam string, domainDeploymentMapIdParam string) (model.DomainDeploymentMap, error)

    // Paginated list of all Domain Deployment Entries for infra.
    //
    // @param domainIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.DomainDeploymentMapListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(domainIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.DomainDeploymentMapListResult, error)

    // If the passed Domain Deployment Map does not already exist, create a new Domain Deployment Map. If it already exist, patch it.
    //
    // @param domainIdParam (required)
    // @param domainDeploymentMapIdParam (required)
    // @param domainDeploymentMapParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(domainIdParam string, domainDeploymentMapIdParam string, domainDeploymentMapParam model.DomainDeploymentMap) error

    // If the passed Domain Deployment Map does not already exist, create a new Domain Deployment Map. If it already exist, replace it.
    //
    // @param domainIdParam (required)
    // @param domainDeploymentMapIdParam (required)
    // @param domainDeploymentMapParam (required)
    // @return com.vmware.nsx_policy.model.DomainDeploymentMap
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(domainIdParam string, domainDeploymentMapIdParam string, domainDeploymentMapParam model.DomainDeploymentMap) (model.DomainDeploymentMap, error)
}
