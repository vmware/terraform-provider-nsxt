// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIdpsCustomSignature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsCustomSignatureRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the custom signature. Use format signature_version_id/signature_id (e.g. default/5000001).",
			},
			"signature_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom signature version ID (e.g. \"default\"). Required if id is only the signature_id.",
			},
			"path": {
				Type:        schema.TypeString,
				Description: "Absolute path of the custom signature.",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name.",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeInt,
				Description: "Revision of the custom signature.",
				Computed:    true,
			},
			"signature_id": {
				Type:        schema.TypeString,
				Description: "System-assigned signature ID.",
				Computed:    true,
			},
			"original_signature_id": {
				Type:        schema.TypeString,
				Description: "Original signature ID from the rule.",
				Computed:    true,
			},
			"original_signature": {
				Type:        schema.TypeString,
				Description: "Original raw signature (Snort rule).",
				Computed:    true,
			},
			"validation_status": {
				Type:        schema.TypeString,
				Description: "Validation status: VALID, INVALID, PENDING, WARNING.",
				Computed:    true,
			},
			"validation_message": {
				Type:        schema.TypeString,
				Description: "Validation message or error if invalid.",
				Computed:    true,
			},
			"severity": {
				Type:        schema.TypeString,
				Description: "Severity of the signature.",
				Computed:    true,
			},
			"signature_severity": {
				Type:        schema.TypeString,
				Description: "Vendor-assigned signature severity.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Signature name (from msg in rule).",
				Computed:    true,
			},
			"signature_revision": {
				Type:        schema.TypeString,
				Description: "Signature revision from the rule.",
				Computed:    true,
			},
			"enable": {
				Type:        schema.TypeBool,
				Description: "Whether the signature is enabled.",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsCustomSignatureRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	id := d.Get("id").(string)
	versionIDFromAttr := d.Get("signature_version_id").(string)

	var versionID, sigID string
	if strings.Contains(id, "/") {
		var err error
		versionID, sigID, err = parseCustomSignatureCompositeID(id)
		if err != nil {
			return err
		}
	} else {
		if versionIDFromAttr == "" {
			return fmt.Errorf("when id is only the signature_id, signature_version_id must be set")
		}
		versionID = versionIDFromAttr
		sigID = id
	}

	client := custom_signature_versions.NewCustomSignaturesClient(connector)
	obj, err := client.Get(versionID, sigID)
	if err != nil {
		if isNotFoundError(err) {
			// Preview signatures may not be gettable by Get(); use list fallback (retry without delay – condition-based)
			var found *model.IdsCustomSignature
			for attempt := 0; attempt < 3; attempt++ {
				found = dataSourceNsxtPolicyIdpsCustomSignatureFindByID(connector, versionID, sigID)
				if found != nil {
					break
				}
			}
			if found == nil {
				return handleDataSourceReadError(d, "IdsCustomSignature", id, err)
			}
			obj = *found
		} else {
			return handleDataSourceReadError(d, "IdsCustomSignature", id, err)
		}
	}

	d.SetId(versionID + "/" + sigID)
	d.Set("signature_version_id", versionID)
	d.Set("path", obj.Path)
	d.Set("display_name", obj.DisplayName)
	d.Set("revision", obj.Revision)
	d.Set("signature_id", obj.SignatureId)
	d.Set("original_signature_id", obj.OriginalSignatureId)
	d.Set("original_signature", obj.OriginalSignature)
	d.Set("validation_status", obj.ValidationStatus)
	d.Set("validation_message", obj.ValidationMessage)
	d.Set("severity", obj.Severity)
	d.Set("signature_severity", obj.SignatureSeverity)
	d.Set("name", obj.Name)
	d.Set("signature_revision", obj.SignatureRevision)
	d.Set("enable", obj.Enable)

	return nil
}

func dataSourceNsxtPolicyIdpsCustomSignatureFindByID(connector client.Connector, versionID, sigID string) *model.IdsCustomSignature {
	found, _ := idpsCustomSignatureFindByID(connector, versionID, sigID)
	return found
}
