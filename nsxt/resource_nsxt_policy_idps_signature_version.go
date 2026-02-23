// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

// NOTE: IDS Signature Versions are system-managed resources in NSX.
// They are created automatically by NSX when signature updates are downloaded.
// This resource provides READ-ONLY access to signature versions for reference purposes.
// Users cannot CREATE or DELETE signature versions directly.
// The only supported operation is making a version ACTIVE via UPDATE.

func resourceNsxtPolicyIdpsSignatureVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIdpsSignatureVersionCreate,
		Read:   resourceNsxtPolicyIdpsSignatureVersionRead,
		Update: resourceNsxtPolicyIdpsSignatureVersionUpdate,
		Delete: resourceNsxtPolicyIdpsSignatureVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"version_id": {
				Type:        schema.TypeString,
				Description: "Version identifier (read-only, assigned by NSX)",
				Computed:    true,
			},
			"change_log": {
				Type:        schema.TypeString,
				Description: "Version change log (read-only)",
				Computed:    true,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Description: "Time when version was downloaded and saved in epoch milliseconds (read-only)",
				Computed:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Version state: ACTIVE (currently used) or NOTACTIVE (available but not used). Can be set to ACTIVE to activate this version.",
				Optional:    true,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Version status: OUTDATED (new version available) or LATEST (up to date) - read-only",
				Computed:    true,
			},
			"user_uploaded": {
				Type:        schema.TypeBool,
				Description: "Whether signature version was uploaded by user (read-only)",
				Computed:    true,
			},
			"auto_update": {
				Type:        schema.TypeBool,
				Description: "Whether signature version came via auto update mechanism (read-only)",
				Computed:    true,
			},
			"sites": {
				Type:        schema.TypeList,
				Description: "Sites mapped with this signature version",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version_name": {
				Type:        schema.TypeString,
				Description: "Version name (read-only)",
				Computed:    true,
			},
		},
	}
}

func resourceNsxtPolicyIdpsSignatureVersionExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	client := intrusion_services.NewSignatureVersionsClient(connector)
	_, err := client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving IDS Signature Version", err)
}

func resourceNsxtPolicyIdpsSignatureVersionCreate(d *schema.ResourceData, m interface{}) error {
	// NOTE: Signature versions are system-managed and cannot be created by users.
	// This "create" operation is actually an import/reference operation.
	// The user must specify an existing version ID via nsx_id.

	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	id := d.Get("nsx_id").(string)
	if id == "" {
		return fmt.Errorf("nsx_id is required for IDS Signature Version resource. Signature versions are system-managed and cannot be created. Use nsx_id to reference an existing version.")
	}

	// Verify the version exists
	client := intrusion_services.NewSignatureVersionsClient(connector)
	_, err := client.Get(id)
	if err != nil {
		return handleCreateError("IdsSignatureVersion", id, err)
	}

	log.Printf("[INFO] Referencing existing IDS Signature Version with ID %s", id)
	d.SetId(id)

	return resourceNsxtPolicyIdpsSignatureVersionRead(d, m)
}

func resourceNsxtPolicyIdpsSignatureVersionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDS Signature Version ID")
	}

	client := intrusion_services.NewSignatureVersionsClient(connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "IdsSignatureVersion", id, err)
	}

	// Set standard fields
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setPolicyTagsInSchema(d, obj.Tags)

	// Set IdsSignatureVersion specific fields
	d.Set("version_id", obj.VersionId)
	d.Set("change_log", obj.ChangeLog)
	d.Set("update_time", obj.UpdateTime)
	d.Set("state", obj.State)
	d.Set("status", obj.Status)
	d.Set("user_uploaded", obj.UserUploaded)
	d.Set("auto_update", obj.AutoUpdate)
	d.Set("sites", obj.Sites)
	d.Set("version_name", obj.VersionName)

	return nil
}

func resourceNsxtPolicyIdpsSignatureVersionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDS Signature Version ID")
	}

	// Only support changing state to ACTIVE
	if d.HasChange("state") {
		newState := d.Get("state").(string)
		switch newState {
		case "ACTIVE":
			log.Printf("[INFO] Making IDS Signature Version %s ACTIVE", id)

			client := intrusion_services.NewSignatureVersionsClient(connector)

			// Get current version object
			obj, err := client.Get(id)
			if err != nil {
				return handleUpdateError("IdsSignatureVersion", id, err)
			}

			// Call make_active_version action
			err = client.Makeactiveversion(obj)
			if err != nil {
				return handleUpdateError("IdsSignatureVersion", id, err)
			}

			log.Printf("[INFO] Successfully made IDS Signature Version %s ACTIVE", id)
		case "NOTACTIVE":
			return fmt.Errorf("Cannot set state to NOTACTIVE directly. To deactivate this version, activate a different version.")
		default:
			return fmt.Errorf("Invalid state value: %s. Must be ACTIVE or NOTACTIVE.", newState)
		}
	}

	// Other fields are read-only and cannot be updated
	if d.HasChange("version_id") || d.HasChange("change_log") || d.HasChange("update_time") ||
		d.HasChange("status") || d.HasChange("user_uploaded") || d.HasChange("auto_update") ||
		d.HasChange("version_name") {
		log.Printf("[WARN] Attempted to update read-only fields on IDS Signature Version %s. These changes will be ignored.", id)
	}

	return resourceNsxtPolicyIdpsSignatureVersionRead(d, m)
}

func resourceNsxtPolicyIdpsSignatureVersionDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDS Signature Version ID")
	}

	// NOTE: Signature versions are system-managed and cannot be deleted by users.
	// This delete operation only removes the resource from Terraform state.
	// The actual signature version remains in NSX.

	log.Printf("[INFO] Removing IDS Signature Version %s from Terraform state (version remains in NSX)", id)
	log.Printf("[INFO] IDS Signature Versions are system-managed and cannot be deleted. The version will remain in NSX.")

	return nil
}
