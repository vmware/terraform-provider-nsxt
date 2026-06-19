// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cliBareMetalServerGroupAssociationsClient "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServerGroupAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerGroupAssociationsRead,

		Schema: map[string]*schema.Schema{
			"id":          getDataSourceIDSchema(),
			"external_id": getBMSExternalIDDataSourceSchema("External ID of the bare metal server", true),
			"enforcement_point_path": {
				Type:        schema.TypeString,
				Description: "Path of the enforcement point",
				Optional:    true,
			},
			"groups": {
				Type:        schema.TypeList,
				Description: "List of groups this bare metal server is a member of",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Description: "Policy path of the group",
							Computed:    true,
						},
						"display_name": {
							Type:        schema.TypeString,
							Description: "Display name of the group",
							Computed:    true,
						},
						"target_type": {
							Type:        schema.TypeString,
							Description: "Type of the target resource",
							Computed:    true,
						},
						"is_valid": {
							Type:        schema.TypeBool,
							Description: "Indicates if the referenced NSX resource is valid",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyBareMetalServerGroupAssociationsRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)

	externalID := d.Get("external_id").(string)
	if externalID == "" {
		return fmt.Errorf("external_id is required")
	}

	enforcementPointPath := d.Get("enforcement_point_path").(string)
	var enforcementPointPathPtr *string
	if enforcementPointPath != "" {
		enforcementPointPathPtr = &enforcementPointPath
	}

	sessionContext := utl.SessionContext{
		ClientType: utl.Local,
	}

	client := cliBareMetalServerGroupAssociationsClient.NewBareMetalServerGroupAssociationsClient(sessionContext, connector)

	associations, err := client.List(externalID, nil, enforcementPointPathPtr, nil, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("Failed to read group associations for bare metal server %s: %v. Verify the server exists and is accessible via Policy API", externalID, err)
	}

	var groups []map[string]interface{}
	for _, result := range associations.Results {
		group := map[string]interface{}{
			"path":         getStringValue(result.TargetId),
			"display_name": getStringValue(result.TargetDisplayName),
			"target_type":  getStringValue(result.TargetType),
		}

		if result.IsValid != nil {
			group["is_valid"] = *result.IsValid
		}

		groups = append(groups, group)
	}

	d.SetId(externalID)
	d.Set("external_id", externalID)
	d.Set("groups", groups)

	return nil
}
