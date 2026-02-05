// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
)

func dataSourceNsxtPolicyIdpsSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsSettingsRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"auto_update_signatures": {
				Type:        schema.TypeBool,
				Description: "Whether automatic update of IDS/IPS signatures is enabled",
				Computed:    true,
			},
			"enable_syslog": {
				Type:        schema.TypeBool,
				Description: "Whether sending IDS/IPS events to syslog server is enabled",
				Computed:    true,
			},
			"oversubscription": {
				Type:        schema.TypeString,
				Description: "Action to take when IDPS engine is oversubscribed (BYPASSED or DROPPED)",
				Computed:    true,
			},
			"enable_custom_signatures": {
				Type:        schema.TypeBool,
				Description: "Whether custom signatures are enabled globally",
				Computed:    true,
			},
			"custom_signature_version_id": {
				Type:        schema.TypeString,
				Description: "The custom signature version ID used for custom signatures",
				Optional:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsSettingsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	// Read IdsSettings
	client := security.NewIntrusionServicesClient(connector)

	obj, err := client.Get()
	if err != nil {
		return fmt.Errorf("Failed to read IdsSettings: %v", err)
	}

	d.SetId("intrusion-services")
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)

	// Set IdsSettings specific fields
	d.Set("auto_update_signatures", obj.AutoUpdate)
	d.Set("enable_syslog", obj.IdsEventsToSyslog)
	d.Set("oversubscription", obj.Oversubscription)

	// Read IdsCustomSignatureSettings if custom_signature_version_id is provided
	customSigVersionID := d.Get("custom_signature_version_id").(string)
	if customSigVersionID != "" {
		customSigSettingsClient := custom_signature_versions.NewSettingsClient(connector)
		customSigSettings, err := customSigSettingsClient.Get(customSigVersionID)
		if err != nil {
			log.Printf("[WARN] Failed to read custom signature settings for version %s: %v", customSigVersionID, err)
			d.Set("enable_custom_signatures", false)
		} else {
			if customSigSettings.EnableCustomSignatures != nil {
				d.Set("enable_custom_signatures", *customSigSettings.EnableCustomSignatures)
			}
		}
	} else {
		d.Set("enable_custom_signatures", false)
	}

	return nil
}
