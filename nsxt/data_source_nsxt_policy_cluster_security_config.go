// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/security"
)

func dataSourceNsxtPolicyClusterSecurityConfig() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNsxtPolicyClusterSecurityConfigRead,
		Description: "Data source to read Cluster Security Configuration for NSX 9.1+. Retrieves DFW (Distributed Firewall) status at the cluster level.",

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster external ID (e.g., uuid:domain-c20)",
				Required:    true,
			},
			"dfw_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether Distributed Firewall (DFW) is enabled on the cluster",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the cluster security configuration",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the cluster security configuration",
				Computed:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "NSX path of the cluster security configuration",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyClusterSecurityConfigRead(d *schema.ResourceData, m interface{}) error {
	// Version check: This data source requires NSX 9.1.0 or higher
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Cluster Security Config data source requires NSX version 9.1.0 or higher (current version: %s)", util.NsxVersion)
	}

	connector := getPolicyConnector(m)
	client := security.NewClusterConfigsClient(connector)

	clusterID := d.Get("cluster_id").(string)
	if clusterID == "" {
		return fmt.Errorf("cluster_id is required")
	}

	config, err := client.Get(clusterID, nil)
	if err != nil {
		return fmt.Errorf("Error reading Cluster Security Config for cluster %s: %v", clusterID, err)
	}

	d.SetId(clusterID)
	d.Set("cluster_id", clusterID)

	return setClusterSecurityConfigInSchema(d, config, true)
}
