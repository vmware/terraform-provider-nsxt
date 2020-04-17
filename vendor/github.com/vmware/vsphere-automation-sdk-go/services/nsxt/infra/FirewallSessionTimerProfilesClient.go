/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: FirewallSessionTimerProfiles
 * Used by client-side stubs.
 */

package infra

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type FirewallSessionTimerProfilesClient interface {

    // API will delete Firewall Session Timer Profile
    //
    // @param firewallSessionTimerProfileIdParam Firewall Session Timer Profile ID (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(firewallSessionTimerProfileIdParam string) error

    // API will get Firewall Session Timer Profile
    //
    // @param firewallSessionTimerProfileIdParam Firewall Session Timer Profile ID (required)
    // @return com.vmware.nsx_policy.model.PolicyFirewallSessionTimerProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(firewallSessionTimerProfileIdParam string) (model.PolicyFirewallSessionTimerProfile, error)

    // API will list all Firewall Session Timer Profiles
    //
    // @param cursorParam Opaque cursor to be used for getting next page of records (supplied by current result page) (optional)
    // @param includeMarkForDeleteObjectsParam Include objects that are marked for deletion in results (optional, default to false)
    // @param includedFieldsParam Comma separated list of fields that should be included in query result (optional)
    // @param pageSizeParam Maximum number of results to return in this page (server may return fewer) (optional, default to 1000)
    // @param sortAscendingParam (optional)
    // @param sortByParam Field by which records are sorted (optional)
    // @return com.vmware.nsx_policy.model.PolicyFirewallSessionTimerProfileListResult
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	List(cursorParam *string, includeMarkForDeleteObjectsParam *bool, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (model.PolicyFirewallSessionTimerProfileListResult, error)

    // API will create/update Firewall Session Timer Profile
    //
    // @param firewallSessionTimerProfileIdParam Firewall Session Timer Profile ID (required)
    // @param policyFirewallSessionTimerProfileParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(firewallSessionTimerProfileIdParam string, policyFirewallSessionTimerProfileParam model.PolicyFirewallSessionTimerProfile) error

    // API will update Firewall Session Timer Profile
    //
    // @param firewallSessionTimerProfileIdParam Firewall Session Timer Profile ID (required)
    // @param policyFirewallSessionTimerProfileParam (required)
    // @return com.vmware.nsx_policy.model.PolicyFirewallSessionTimerProfile
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(firewallSessionTimerProfileIdParam string, policyFirewallSessionTimerProfileParam model.PolicyFirewallSessionTimerProfile) (model.PolicyFirewallSessionTimerProfile, error)
}
