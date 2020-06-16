/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: DhcpStaticBindingConfigs
 * Used by client-side stubs.
 */

package segments

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type DhcpStaticBindingConfigsClient interface {

    // Delete DHCP static binding
    //
    // @param segmentIdParam (required)
    // @param bindingIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(segmentIdParam string, bindingIdParam string) error

    // Read DHCP static binding
    //
    // @param segmentIdParam (required)
    // @param bindingIdParam (required)
    // @return com.vmware.nsx_global_policy.model.DhcpStaticBindingConfig
    // The return value will contain all the properties defined in model.DhcpStaticBindingConfig.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(segmentIdParam string, bindingIdParam string) (*data.StructValue, error)

    // Paginated list of all DHCP static binding instances
    //
    // @param segmentIdParam (required)
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_global_policy.model.DhcpStaticBindingConfigListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(segmentIdParam string, cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.DhcpStaticBindingConfigListResult, error)

    // If binding with the binding-id is not already present, create a new DHCP static binding instance. If it already exists, replace the existing DHCP static binding instance with specified attributes.
    //
    // @param segmentIdParam (required)
    // @param bindingIdParam (required)
    // @param dhcpStaticBindingConfigParam (required)
    // The parameter must contain all the properties defined in model.DhcpStaticBindingConfig.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(segmentIdParam string, bindingIdParam string, dhcpStaticBindingConfigParam *data.StructValue) error

    // If binding with the binding-id is not already present, create a new DHCP static binding instance. If it already exists, replace the existing DHCP static binding instance with this object.
    //
    // @param segmentIdParam (required)
    // @param bindingIdParam (required)
    // @param dhcpStaticBindingConfigParam (required)
    // The parameter must contain all the properties defined in model.DhcpStaticBindingConfig.
    // @return com.vmware.nsx_global_policy.model.DhcpStaticBindingConfig
    // The return value will contain all the properties defined in model.DhcpStaticBindingConfig.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(segmentIdParam string, bindingIdParam string, dhcpStaticBindingConfigParam *data.StructValue) (*data.StructValue, error)
}
