/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: CustomAttributes
 * Used by client-side stubs.
 */

package context_profiles

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type CustomAttributesClient interface {

    // This API adds/removes custom attribute values from list for a given attribute key.
    //
    // @param policyAttributesParam (required)
    // @param actionParam Add or Remove Custom Context Profile Attribute values. (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(policyAttributesParam model.PolicyAttributes, actionParam string) error

    // This API updates custom attribute value list for given key.
    //
    // @param policyAttributesParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(policyAttributesParam model.PolicyAttributes) error
}
