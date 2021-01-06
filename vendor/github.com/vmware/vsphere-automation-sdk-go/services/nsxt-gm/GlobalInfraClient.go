/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GlobalInfra
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type GlobalInfraClient interface {

    // Read infra. Returns only the infra related properties. Inner object are not populated.
    //
    // @param basePathParam Base Path for retrieving hierarchical intent (optional)
    // @param filterParam Filter string as java regex (optional)
    // @param typeFilterParam Filter string to retrieve hierarchy. (optional)
    // @return com.vmware.nsx_global_policy.model.Infra
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(basePathParam *string, filterParam *string, typeFilterParam *string) (model.Infra, error)

    // Patch API at infra level can be used in two flavours 1. Like a regular API to update Infra object 2. Hierarchical API: To create/update/delete entire or part of intent hierarchy Hierarchical API: Provides users a way to create entire or part of intent in single API invocation. Input is expressed in a tree format. Each node in tree can have multiple children of different types. System will resolve the dependecies of nodes within the intent tree and will create the model. Children for any node can be specified using ChildResourceReference or ChildPolicyConfigResource. If a resource is specified using ChildResourceReference then it will not be updated only its children will be updated. If Object is specified using ChildPolicyConfigResource, object along with its children will be updated. Hierarchical API can also be used to delete any sub-branch of entire tree.
    //
    // @param infraParam (required)
    // @param enforceRevisionCheckParam Force revision check (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(infraParam model.Infra, enforceRevisionCheckParam *bool) error

    // Update the infra including all the nested entities
    //
    // @param infraParam (required)
    // @return com.vmware.nsx_global_policy.model.Infra
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(infraParam model.Infra) (model.Infra, error)
}
