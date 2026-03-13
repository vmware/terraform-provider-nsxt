// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
)

func dataSourceNsxtPolicyIdpsSignatureDiff() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsSignatureDiffRead,

		Schema: map[string]*schema.Schema{
			"signature_version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the custom signature version (e.g. \"default\").",
			},
			"newly_added_signatures": {
				Type:        schema.TypeList,
				Description: "List of custom signature IDs that are newly added (unpublished).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deleted_signatures": {
				Type:        schema.TypeList,
				Description: "List of custom signature IDs that are deleted (removed from published).",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"existing_signatures": {
				Type:        schema.TypeList,
				Description: "List of custom signature IDs that exist in both published and current state.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsSignatureDiffRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	versionID := d.Get("signature_version_id").(string)
	client := custom_signature_versions.NewCustomSignaturesDiffClient(connector)
	diff, err := client.Get(versionID)
	if err != nil {
		return handleDataSourceReadError(d, "IdsCustomSignaturesDiff", versionID, err)
	}

	d.SetId(versionID)
	d.Set("newly_added_signatures", diff.NewlyAddedSignatures)
	d.Set("deleted_signatures", diff.DeletedSignatures)
	d.Set("existing_signatures", diff.ExistingSignatures)

	return nil
}
