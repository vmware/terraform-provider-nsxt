// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNsxtProviderContext() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtProviderContextCreate,
		Read:   resourceNsxtProviderContextRead,
		Delete: resourceNsxtProviderContextDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNsxtProviderContextCreate(d *schema.ResourceData, m interface{}) error {
	configID := uuid.NewString()
	if err := d.Set("config_id", configID); err != nil {
		return err
	}
	// Singleton instance ID.
	d.SetId("nsxt_provider_context")
	InitConfigIDOnce(m, configID)
	return resourceNsxtProviderContextRead(d, m)
}

func resourceNsxtProviderContextRead(d *schema.ResourceData, m interface{}) error {
	// Source of truth is Terraform state.
	configID := d.Get("config_id").(string)
	if configID == "" {
		configID = uuid.NewString()
		if err := d.Set("config_id", configID); err != nil {
			return err
		}
	}
	if d.Id() == "" {
		d.SetId("nsxt_provider_context")
	}
	return InitConfigIDOnce(m, configID)
}

func resourceNsxtProviderContextDelete(d *schema.ResourceData, m interface{}) error {
	// Removing from state removes the persisted config_id.
	d.SetId("")
	return nil
}
