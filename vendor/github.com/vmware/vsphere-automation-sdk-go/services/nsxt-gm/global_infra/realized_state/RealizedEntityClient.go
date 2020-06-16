/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: RealizedEntity
 * Used by client-side stubs.
 */

package realized_state


type RealizedEntityClient interface {

    // Refresh the status and statistics of all realized entities associated with given intent path synchronously. The vmw-async: True HTTP header cannot be used with this API.
    //
    // @param intentPathParam String Path of the intent object (required)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Refresh(intentPathParam string, enforcementPointPathParam *string) error
}
