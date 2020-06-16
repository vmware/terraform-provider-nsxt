/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: ErrorResolver
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type ErrorResolverClient interface {

    // Returns some metadata about the given error_id. This includes information of whether there is a resolver present for the given error_id and its associated user input data
    //
    // @param errorIdParam (required)
    // @return com.vmware.nsx_global_policy.model.ErrorResolverInfo
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(errorIdParam string) (model.ErrorResolverInfo, error)

    // Returns a list of metadata for all the error resolvers registered.
    // @return com.vmware.nsx_global_policy.model.ErrorResolverInfoList
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List() (model.ErrorResolverInfoList, error)

    // Invokes the corresponding error resolver for the given error(s) present in the payload
    //
    // @param errorResolverMetadataListParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Resolveerror(errorResolverMetadataListParam model.ErrorResolverMetadataList) error
}
