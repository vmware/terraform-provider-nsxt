/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtManagementCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtManagementClusterRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique identifier of this cluster.",
				Computed:    true,
			},
			"node_sha256_thumbprint": {
				Type:        schema.TypeString,
				Description: "SHA256 of certificate thumbprint of this manager node.",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtManagementClusterRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	clusterObj, resp, err := nsxClient.NsxComponentAdministrationApi.ReadClusterConfig(nsxClient.Context)
	if err != nil {
		return fmt.Errorf("Error while reading cluster configuration: %v", err)
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected Response while reading cluster configuration. Status Code: %d", resp.StatusCode)
	}

	nodeList, resp, err := nsxClient.NsxComponentAdministrationApi.ListClusterNodeConfigs(nsxClient.Context, nil)
	if err != nil {
		return fmt.Errorf("Error while reading cluster node configuration: %v", err)
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected Response while reading cluster node configuration. Status Code: %d", resp.StatusCode)
	}
	for _, nodeConfig := range nodeList.Results {
		if nodeConfig.ManagerRole != nil && nodeConfig.ManagerRole.ApiListenAddr != nil && nodeConfig.ManagerRole.ApiListenAddr.IpAddress == m.(nsxtClients).Host[len("https://"):] {
			if nodeConfig.ManagerRole.ApiListenAddr.CertificateSha256Thumbprint == "" {
				return fmt.Errorf("Manager node thumbprint not found while reading cluster node configuration")
			}
			d.Set("node_sha256_thumbprint", nodeConfig.ManagerRole.ApiListenAddr.CertificateSha256Thumbprint)
		}
	}

	if clusterObj.ClusterId == "" {
		return fmt.Errorf("Cluster id not found")
	}
	if d.Get("node_sha256_thumbprint").(string) == "" {
		return fmt.Errorf("Cluster node sha256 thumbprint not found")
	}

	d.SetId(clusterObj.ClusterId)
	return nil
}
