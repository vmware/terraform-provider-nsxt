/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/realized_state"
	"time"
)

func dataSourceNsxtPolicyRealizationInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRealizationInfoRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The path for the policy resource",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"entity_type": {
				Type:        schema.TypeString,
				Description: "The entity type of the realized resource",
				Computed:    true,
				Optional:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the realized resource",
				Computed:    true,
			},
			"realized_id": {
				Type:        schema.TypeString,
				Description: "The ID of the realized resource",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyRealizationInfoRead(d *schema.ResourceData, m interface{}) error {
	// Read the realization info by the path, and wait till it is valid
	connector := getPolicyConnector(m)
	client := realized_state.NewDefaultRealizedEntitiesClient(connector)

	// Get the realization info of this resource
	path := d.Get("path").(string)
	entityType := d.Get("entity_type").(string)

	// Dummy id, just because each data source needs one
	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	pendingStates := []string{"UNKNOWN", "UNREALIZED"}
	targetStates := []string{"REALIZED", "ERROR"}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {

			realizationResult, realizationError := client.List(path, &policySite)
			state := "UNKNOWN"
			if realizationError == nil {
				// Find the right entry
				for _, objInList := range realizationResult.Results {
					if objInList.State != nil {
						state = *objInList.State
					}
					if entityType == "" {
						// Take the first one
						d.Set("state", state)
						d.Set("entity_type", *objInList.EntityType)
						if objInList.RealizationSpecificIdentifier == nil {
							d.Set("realized_id", "")
						} else {
							d.Set("realized_id", *objInList.RealizationSpecificIdentifier)
						}
						return realizationResult, state, nil
					} else if (objInList.EntityType != nil) && (*objInList.EntityType == entityType) {
						d.Set("state", state)
						if objInList.RealizationSpecificIdentifier == nil {
							d.Set("realized_id", "")
						} else {
							d.Set("realized_id", *objInList.RealizationSpecificIdentifier)
						}
						return realizationResult, state, nil
					}
				}
				// Realization info not found yet
				d.Set("state", "UNKNOWN")
				d.Set("realized_id", "")
				return realizationResult, "UNKNOWN", nil
			}
			return realizationResult, "", realizationError
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 1 * time.Second,
		Delay:      1 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Failed to get realization information for %s: %v", path, err)
	}
	return nil
}
