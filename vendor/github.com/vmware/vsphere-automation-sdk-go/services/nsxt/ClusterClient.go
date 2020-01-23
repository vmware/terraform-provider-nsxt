/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Cluster
 * Used by client-side stubs.
 */

package nsx_policy


type ClusterClient interface {

    // Request one-time backup. The backup will be uploaded using the same server configuration as for automatic backup.
    //
    // @param frameTypeParam Frame type (optional, default to LOCAL_LOCAL_MANAGER)
    // @param siteIdParam Site ID (optional, default to localhost)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Backuptoremote(frameTypeParam *string, siteIdParam *string) error

    // Request one-time inventory summary. The backup will be uploaded using the same server configuration as for an automatic backup.
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Summarizeinventorytoremote() error
}
