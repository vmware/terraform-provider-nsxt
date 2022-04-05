/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var valueTypeValues = []string{"bios_id", "external_id", "instance_id"}

func dataSourceNsxtPolicyVMs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVMsRead,

		Schema: map[string]*schema.Schema{
			// TODO: add option to filter by display name regex
			"value_type": {
				Type:         schema.TypeString,
				Description:  "Type of data populated in map value",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(valueTypeValues, false),
				Default:      "bios_id",
			},
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of VM instance ID by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyVMsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	valueType := d.Get("value_type").(string)
	vmMap := make(map[string]interface{})

	allVMs, err := listAllPolicyVirtualMachines(connector, m)
	if err != nil {
		return fmt.Errorf("Error reading Virtual Machines: %v", err)
	}

	for _, vm := range allVMs {
		if vm.DisplayName == nil {
			continue
		}
		computeIDMap := collectSeparatedStringListToMap(vm.ComputeIds, ":")
		if valueType == "instance_id" {
			vmMap[*vm.DisplayName] = computeIDMap[nsxtPolicyInstanceUUIDKey]
		} else if valueType == "bios_id" {
			vmMap[*vm.DisplayName] = computeIDMap[nsxtPolicyBiosUUIDKey]
		} else if valueType == "external_id" {
			if vm.ExternalId == nil {
				// skip this vm
				continue
			}
			vmMap[*vm.DisplayName] = *vm.ExternalId
		}
	}

	d.SetId(newUUID())
	d.Set("items", vmMap)

	return nil
}
