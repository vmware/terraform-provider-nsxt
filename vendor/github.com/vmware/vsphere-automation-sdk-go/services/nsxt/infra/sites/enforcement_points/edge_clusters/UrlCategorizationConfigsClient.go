/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Interface file for service: UrlCategorizationConfigs
 * Used by client-side stubs.
 */

package edge_clusters

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type UrlCategorizationConfigsClient interface {

    // Delete PolicyUrlCategorizationConfig. If deleted, the URL categorization will be disabled for that edge cluster.
    //
    // @param siteIdParam (required)
    // @param enforcementPointIdParam (required)
    // @param edgeClusterIdParam (required)
    // @param urlCategorizationConfigIdParam (required)
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Delete(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) error

    // Gets a PolicyUrlCategorizationConfig. This returns the details of the config like whether the URL categorization is enabled or disabled, the id of the context profiles which are used to filter the categories, and the update frequency of the data from the cloud.
    //
    // @param siteIdParam (required)
    // @param enforcementPointIdParam (required)
    // @param edgeClusterIdParam (required)
    // @param urlCategorizationConfigIdParam (required)
    // @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Get(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string) (model.PolicyUrlCategorizationConfig, error)

    // Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
    //
    // @param siteIdParam (required)
    // @param enforcementPointIdParam (required)
    // @param edgeClusterIdParam (required)
    // @param urlCategorizationConfigIdParam (required)
    // @param policyUrlCategorizationConfigParam (required)
    // @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Patch(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error)

    // Creates/Updates a PolicyUrlCategorizationConfig. Creating or updating the PolicyUrlCategorizationConfig will enable or disable URL categorization for the given edge cluster. If the context_profiles field is empty, the edge cluster will detect all the categories of URLs. If context_profiles field has any context profiles, the edge cluster will detect only the categories listed within those context profiles. The context profiles should have attribute type URL_CATEGORY. The update_frequency specifies how frequently in minutes, the edge cluster will get updates about the URL data from the URL categorization cloud service. If the update_frequency is not specified, the default update frequency will be 30 min.
    //
    // @param siteIdParam (required)
    // @param enforcementPointIdParam (required)
    // @param edgeClusterIdParam (required)
    // @param urlCategorizationConfigIdParam (required)
    // @param policyUrlCategorizationConfigParam (required)
    // @return com.vmware.nsx_policy.model.PolicyUrlCategorizationConfig
    // @throws InvalidRequest  Bad Request, Precondition Failed
    // @throws Unauthorized  Forbidden
    // @throws ServiceUnavailable  Service Unavailable
    // @throws InternalServerError  Internal Server Error
    // @throws NotFound  Not Found
	Update(siteIdParam string, enforcementPointIdParam string, edgeClusterIdParam string, urlCategorizationConfigIdParam string, policyUrlCategorizationConfigParam model.PolicyUrlCategorizationConfig) (model.PolicyUrlCategorizationConfig, error)
}
