// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/custom_signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var idpsOversubscriptionValues = []string{"BYPASSED", "DROPPED"}

const (
	idpsSettingsID                   = "intrusion-services"
	idpsDefaultAutoUpdate            = false
	idpsDefaultEnableSyslog          = false
	idpsDefaultOversubscription      = "BYPASSED"
	idpsDefaultEnableCustomSignature = false
)

func resourceNsxtPolicyIdpsSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIdpsSettingsCreate,
		Read:   resourceNsxtPolicyIdpsSettingsRead,
		Update: resourceNsxtPolicyIdpsSettingsUpdate,
		Delete: resourceNsxtPolicyIdpsSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "Global IDPS settings configuration. This is a singleton resource managing IDS/IPS system settings including syslog, oversubscription, signature auto-update, and custom signature enablement.",

		Schema: map[string]*schema.Schema{
			"path": getPathSchema(),
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name (read-only for singleton resource)",
				Computed:    true,
			},
			"description": getDescriptionSchema(),
			"revision":    getRevisionSchema(),
			"tag":         getTagsSchema(),
			"auto_update_signatures": {
				Type:        schema.TypeBool,
				Description: "Enable automatic update of IDS/IPS signatures. When enabled, NSX-T will automatically download and install the latest signature updates.",
				Optional:    true,
				Default:     idpsDefaultAutoUpdate,
			},
			"enable_syslog": {
				Type:        schema.TypeBool,
				Description: "Enable sending IDS/IPS events to syslog server. When enabled, IDPS events will be forwarded to configured syslog servers.",
				Optional:    true,
				Default:     idpsDefaultEnableSyslog,
			},
			"oversubscription": {
				Type:         schema.TypeString,
				Description:  "Action to take when IDPS engine is oversubscribed. BYPASSED: oversubscribed packets bypass IDPS engine. DROPPED: oversubscribed packets are dropped.",
				Optional:     true,
				Default:      idpsDefaultOversubscription,
				ValidateFunc: validation.StringInSlice(idpsOversubscriptionValues, false),
			},
			"enable_custom_signatures": {
				Type:        schema.TypeBool,
				Description: "Enable custom signatures globally. When enabled, custom signatures will be active in IDPS policies. Requires a custom signature version to be configured.",
				Optional:    true,
				Default:     idpsDefaultEnableCustomSignature,
			},
			"custom_signature_version_id": {
				Type:        schema.TypeString,
				Description: "The custom signature version ID to use when enabling custom signatures. This should reference an existing custom signature version resource.",
				Optional:    true,
			},
		},
	}
}

func resourceNsxtPolicyIdpsSettingsCreate(d *schema.ResourceData, m interface{}) error {
	// IDPS Settings is a singleton resource - use fixed ID
	d.SetId(idpsSettingsID)
	log.Printf("[INFO] Creating IDPS Settings with ID %s", idpsSettingsID)
	return resourceNsxtPolicyIdpsSettingsUpdate(d, m)
}

func resourceNsxtPolicyIdpsSettingsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Settings ID")
	}

	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	// Read IdsSettings
	client := security.NewIntrusionServicesClient(connector)

	obj, err := client.Get()
	if err != nil {
		return handleReadError(d, "IdsSettings", id, err)
	}

	if obj.DisplayName != nil {
		d.Set("display_name", *obj.DisplayName)
	}
	if obj.Description != nil {
		d.Set("description", *obj.Description)
	}
	if obj.Path != nil {
		d.Set("path", *obj.Path)
	}
	if obj.Revision != nil {
		d.Set("revision", *obj.Revision)
	}

	setPolicyTagsInSchema(d, obj.Tags)

	// Set IdsSettings specific fields
	if obj.AutoUpdate != nil {
		d.Set("auto_update_signatures", *obj.AutoUpdate)
	}

	if obj.IdsEventsToSyslog != nil {
		d.Set("enable_syslog", *obj.IdsEventsToSyslog)
	}

	if obj.Oversubscription != nil {
		d.Set("oversubscription", *obj.Oversubscription)
	}

	// Read IdsCustomSignatureSettings if custom_signature_version_id is set
	customSigVersionID := d.Get("custom_signature_version_id").(string)
	if customSigVersionID != "" {
		customSigSettingsClient := custom_signature_versions.NewSettingsClient(connector)
		customSigSettings, err := customSigSettingsClient.Get(customSigVersionID)
		if err != nil {
			log.Printf("[WARN] Failed to read custom signature settings for version %s: %v", customSigVersionID, err)
		} else {
			if customSigSettings.EnableCustomSignatures != nil {
				d.Set("enable_custom_signatures", *customSigSettings.EnableCustomSignatures)
			}
		}
	}

	return nil
}

func resourceNsxtPolicyIdpsSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Settings ID")
	}

	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	// Update IdsSettings
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	autoUpdate := d.Get("auto_update_signatures").(bool)
	enableSyslog := d.Get("enable_syslog").(bool)
	oversubscription := d.Get("oversubscription").(string)

	resourceType := "IdsSettings"
	obj := model.IdsSettings{
		Description:       &description,
		Tags:              tags,
		ResourceType:      &resourceType,
		AutoUpdate:        &autoUpdate,
		IdsEventsToSyslog: &enableSyslog,
		Oversubscription:  &oversubscription,
	}

	// Get current revision
	client := security.NewIntrusionServicesClient(connector)

	existingObj, err := client.Get()
	if err == nil && existingObj.Revision != nil {
		obj.Revision = existingObj.Revision
	}

	log.Printf("[INFO] Updating IDPS Settings with ID %s", id)
	log.Printf("[DEBUG] IDPS Settings - auto_update: %v, enable_syslog: %v, oversubscription: %s", autoUpdate, enableSyslog, oversubscription)

	_, err = client.Update(obj)
	if err != nil {
		return handleUpdateError("IdsSettings", id, err)
	}

	// Update IdsCustomSignatureSettings if custom_signature_version_id is set
	customSigVersionID := d.Get("custom_signature_version_id").(string)
	enableCustomSig := d.Get("enable_custom_signatures").(bool)

	if customSigVersionID != "" {
		customSigResourceType := "IdsCustomSignatureSettings"
		customSigSettingsID := "settings"
		customSigSettings := model.IdsCustomSignatureSettings{
			ResourceType:           &customSigResourceType,
			Id:                     &customSigSettingsID,
			EnableCustomSignatures: &enableCustomSig,
		}

		customSigSettingsClient := custom_signature_versions.NewSettingsClient(connector)

		// Get current revision
		existingCustomSigSettings, err := customSigSettingsClient.Get(customSigVersionID)
		if err == nil && existingCustomSigSettings.Revision != nil {
			customSigSettings.Revision = existingCustomSigSettings.Revision
		}

		log.Printf("[INFO] Updating Custom Signature Settings for version %s - enable_custom_signatures: %v", customSigVersionID, enableCustomSig)

		err = customSigSettingsClient.Patch(customSigVersionID, customSigSettings)
		if err != nil {
			return handleUpdateError("IdsCustomSignatureSettings", customSigVersionID, err)
		}
	} else if enableCustomSig {
		log.Printf("[WARN] enable_custom_signatures is true but custom_signature_version_id is not set. Custom signatures will not be enabled.")
	}

	return resourceNsxtPolicyIdpsSettingsRead(d, m)
}

func resourceNsxtPolicyIdpsSettingsDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Settings ID")
	}

	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	// IDPS Settings cannot be truly deleted, only reset to defaults
	// Reset to default values
	log.Printf("[INFO] Resetting IDPS Settings to defaults with ID %s (DELETE operation)", id)

	client := security.NewIntrusionServicesClient(connector)

	// Get current object to preserve revision
	existingObj, err := client.Get()
	if err != nil {
		return handleDeleteError("IdsSettings", id, err)
	}

	// Reset to default values
	resourceType := "IdsSettings"
	defaultAutoUpdate := idpsDefaultAutoUpdate
	defaultEnableSyslog := idpsDefaultEnableSyslog
	defaultOversubscription := idpsDefaultOversubscription

	obj := model.IdsSettings{
		ResourceType:      &resourceType,
		AutoUpdate:        &defaultAutoUpdate,
		IdsEventsToSyslog: &defaultEnableSyslog,
		Oversubscription:  &defaultOversubscription,
		Revision:          existingObj.Revision,
	}

	_, err = client.Update(obj)
	if err != nil {
		return handleDeleteError("IdsSettings", id, err)
	}

	// Reset custom signature settings if version ID was set
	customSigVersionID := d.Get("custom_signature_version_id").(string)
	if customSigVersionID != "" {
		defaultEnableCustomSig := idpsDefaultEnableCustomSignature
		customSigResourceType := "IdsCustomSignatureSettings"
		customSigSettingsID := "settings"
		customSigSettings := model.IdsCustomSignatureSettings{
			ResourceType:           &customSigResourceType,
			Id:                     &customSigSettingsID,
			EnableCustomSignatures: &defaultEnableCustomSig,
		}

		customSigSettingsClient := custom_signature_versions.NewSettingsClient(connector)

		existingCustomSigSettings, err := customSigSettingsClient.Get(customSigVersionID)
		if err == nil && existingCustomSigSettings.Revision != nil {
			customSigSettings.Revision = existingCustomSigSettings.Revision
		}

		err = customSigSettingsClient.Patch(customSigVersionID, customSigSettings)
		if err != nil {
			log.Printf("[WARN] Failed to reset custom signature settings: %v", err)
		}
	}

	return nil
}
