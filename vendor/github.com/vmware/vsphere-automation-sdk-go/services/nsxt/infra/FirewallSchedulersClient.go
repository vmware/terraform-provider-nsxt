/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FirewallSchedulers
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type FirewallSchedulersClient interface {

    // Deletes the specified PolicyFirewallScheduler. If scheduler is consumed in a security policy, it won't get deleted.
    //
    // @param firewallSchedulerIdParam (required)
    // @param forceParam Force delete the resource even if it is being used somewhere (optional, default to false)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(firewallSchedulerIdParam string, forceParam *bool) error

    // Get a PolicyFirewallScheduler by id
    //
    // @param firewallSchedulerIdParam (required)
    // @return com.vmware.nsx_policy.model.PolicyFirewallScheduler
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(firewallSchedulerIdParam string) (model.PolicyFirewallScheduler, error)

    // Get all PolicyFirewallSchedulers
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyFirewallSchedulerListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyFirewallSchedulerListResult, error)

    // Creates/Updates a PolicyFirewallScheduler, which can be set at security policy. Note that at least one property out of \"days\", \"start_date\", \"time_interval\", \"end_date\" is required if \"recurring\" field is true. Also \"start_time\" and \"end_time\" should not be present. And if \"recurring\" field is false then \"start_date\" and \"end_date\" is mandatory, \"start_time\" and \"end_time\" is optional. Also the fields \"days\" and \"time_interval\" should not be present.
    //
    // @param firewallSchedulerIdParam (required)
    // @param policyFirewallSchedulerParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(firewallSchedulerIdParam string, policyFirewallSchedulerParam model.PolicyFirewallScheduler) error

    // Updates a PolicyFirewallScheduler, which can be set at security policy. Note that at least one property out of \"days\", \"start_date\", \"time_interval\", \"end_date\" is required if \"recurring\" field is true. Also \"start_time\" and \"end_time\" should not be present. And if \"recurring\" field is false then \"start_date\" and \"end_date\" is mandatory, \"start_time\" and \"end_time\" is optional. Also the fields \"days\" and \"time_interval\" should not be present.
    //
    // @param firewallSchedulerIdParam (required)
    // @param policyFirewallSchedulerParam (required)
    // @return com.vmware.nsx_policy.model.PolicyFirewallScheduler
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(firewallSchedulerIdParam string, policyFirewallSchedulerParam model.PolicyFirewallScheduler) (model.PolicyFirewallScheduler, error)
}
