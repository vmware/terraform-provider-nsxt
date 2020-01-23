/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: GlobalManagers
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type GlobalManagersClient interface {

    // Switch over from Active to Standby Global Manager. This operation will fail if there is no Standby Global Manager.
    //
    // @param actionParam Indicates whether it is managed switchover or due to failure (required)
    // @return com.vmware.nsx_policy.model.GlobalManager
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(actionParam string) (model.GlobalManager, error)

    // Delete a particular global manager under Infra. Global Manager id 'self' is reserved and can be used for referring to local logged in Global Manager. Example - /infra/global-managers/self
    //
    // @param globalManagerIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(globalManagerIdParam string) error

    // Retrieve information about a particular configured global manager. Global Manager id 'self' is reserved and can be used for referring to local logged in Global Manager. Example - /infra/global-managers/self
    //
    // @param globalManagerIdParam (required)
    // @return com.vmware.nsx_policy.model.GlobalManager
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(globalManagerIdParam string) (model.GlobalManager, error)

    // List Global Managers under Infra.
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.GlobalManagerListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.GlobalManagerListResult, error)

    // Create or patch a Global Manager under Infra. Global Manager id 'self' is reserved and can be used for referring to local logged in Global Manager. Example - /infra/global-managers/self
    //
    // @param globalManagerIdParam (required)
    // @param globalManagerParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(globalManagerIdParam string, globalManagerParam model.GlobalManager) error

    // Create or fully replace Global Manager under Infra. Revision is optional for creation and required for update. Global Manager id 'self' is reserved and can be used for referring to local logged in Global Manager. Example - /infra/global-managers/self
    //
    // @param globalManagerIdParam (required)
    // @param globalManagerParam (required)
    // @return com.vmware.nsx_policy.model.GlobalManager
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(globalManagerIdParam string, globalManagerParam model.GlobalManager) (model.GlobalManager, error)
}
