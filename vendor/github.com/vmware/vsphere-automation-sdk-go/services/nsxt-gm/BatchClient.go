/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Batch
 * Used by client-side stubs.
 */

package nsx_global_policy

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
)

type BatchClient interface {

    // Enables you to make multiple API requests using a single request. The batch API takes in an array of logical HTTP requests represented as JSON arrays. Each request has a method (GET, PUT, POST, or DELETE), a relative_url (the portion of the URL after https://<nsx-mgr>/api/), optional headers array (corresponding to HTTP headers) and an optional body (for POST and PUT requests). The batch API returns an array of logical HTTP responses represented as JSON arrays. Each response has a status code, an optional headers array and an optional body (which is a JSON-encoded string).
    //
    // @param batchRequestParam (required)
    // @param atomicParam transactional atomicity for the batch of requests embedded in the batch list (optional, default to false)
    // @return com.vmware.nsx_global_policy.model.BatchResponse
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Create(batchRequestParam model.BatchRequest, atomicParam *bool) (model.BatchResponse, error)
}
