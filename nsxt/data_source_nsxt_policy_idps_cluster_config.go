// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIdpsClusterConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsClusterConfigRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"ids_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether IDPS is enabled on this cluster",
				Computed:    true,
			},
			"cluster": {
				Type:        schema.TypeList,
				Description: "Cluster reference configuration",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_id": {
							Type:        schema.TypeString,
							Description: "Cluster target ID (e.g., domain-c123)",
							Computed:    true,
						},
						"target_type": {
							Type:        schema.TypeString,
							Description: "Target type (e.g., VC_Cluster)",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsClusterConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.IdsClusterConfig

	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	if objID != "" {
		// Get by ID
		objGet, err := client.Get(objID, nil)
		if err != nil {
			return handleDataSourceReadError(d, "IdsClusterConfig", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining IdsClusterConfig ID or name during read")
	} else {
		// List and search by display name
		objList, err := client.List(nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("IdsClusterConfig", err)
		}

		// Find matching object (prefer perfect match, then prefix match)
		var perfectMatch []model.IdsClusterConfig
		var prefixMatch []model.IdsClusterConfig
		for _, item := range objList.Results {
			if item.DisplayName != nil {
				if strings.HasPrefix(*item.DisplayName, objName) {
					prefixMatch = append(prefixMatch, item)
				}
				if *item.DisplayName == objName {
					perfectMatch = append(perfectMatch, item)
				}
			}
		}

		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple IdsClusterConfigs with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple IdsClusterConfigs with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("IdsClusterConfig with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("ids_enabled", obj.IdsEnabled)

	// Set cluster information
	if obj.Cluster != nil {
		clusterList := make([]map[string]interface{}, 0, 1)
		clusterMap := make(map[string]interface{})
		clusterMap["target_id"] = obj.Cluster.TargetId
		clusterMap["target_type"] = obj.Cluster.TargetType
		clusterList = append(clusterList, clusterMap)
		d.Set("cluster", clusterList)
	}

	return nil
}
