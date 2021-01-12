/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: Stats
 * Used by client-side stubs.
 */

package intrusion_services


type StatsClient interface {

    // Sets IDS/IPS rule statistics counter to zero. - no enforcement point path specified: Reset of stats will be executed for each enforcement point. - {enforcement_point_path}: Reset of stats will be executed only for the given enforcement point.
    //
    // @param categoryParam Aggregation statistic category (optional, default to IDPSDFW)
    // @param enforcementPointPathParam String Path of the enforcement point (optional)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Reset(categoryParam *string, enforcementPointPathParam *string) error
}
