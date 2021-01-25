/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Areas
 * Used by client-side stubs.
 */

package ospf

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type AreasClient interface {

    // Delete OSPF Area config
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param areaIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(tier0IdParam string, localeServiceIdParam string, areaIdParam string) error

    // Read OSPF Area config
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param areaIdParam (required)
    // @return com.vmware.nsx_policy.model.OspfAreaConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(tier0IdParam string, localeServiceIdParam string, areaIdParam string) (model.OspfAreaConfig, error)

    // List all OSPF area configurations.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.OspfAreaConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(tier0IdParam string, localeServiceIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.OspfAreaConfigListResult, error)

    // If OSPF Area config is not already present, create OSPF Area config. If it already exists, replace the OSPF Area config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param areaIdParam (required)
    // @param ospfAreaConfigParam (required)
    // @return com.vmware.nsx_policy.model.OspfAreaConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(tier0IdParam string, localeServiceIdParam string, areaIdParam string, ospfAreaConfigParam model.OspfAreaConfig) (model.OspfAreaConfig, error)

    // If OSPF Area config is not already present, create OSPF Area config. If it already exists, replace the OSPF Area config with this object.
    //
    // @param tier0IdParam (required)
    // @param localeServiceIdParam (required)
    // @param areaIdParam (required)
    // @param ospfAreaConfigParam (required)
    // @return com.vmware.nsx_policy.model.OspfAreaConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(tier0IdParam string, localeServiceIdParam string, areaIdParam string, ospfAreaConfigParam model.OspfAreaConfig) (model.OspfAreaConfig, error)
}
