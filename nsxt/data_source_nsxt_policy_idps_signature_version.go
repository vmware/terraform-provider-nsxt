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

func dataSourceNsxtPolicyIdpsSignatureVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIdpsSignatureVersionRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"version_id": {
				Type:        schema.TypeString,
				Description: "Version identifier",
				Computed:    true,
			},
			"change_log": {
				Type:        schema.TypeString,
				Description: "Version change log",
				Computed:    true,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Description: "Time when version was downloaded and saved (epoch milliseconds)",
				Computed:    true,
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Version state (ACTIVE or NOTACTIVE)",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Version status (OUTDATED or LATEST)",
				Computed:    true,
			},
			"user_uploaded": {
				Type:        schema.TypeBool,
				Description: "Whether signature version was uploaded by user",
				Computed:    true,
			},
			"auto_update": {
				Type:        schema.TypeBool,
				Description: "Whether signature version came via auto update mechanism",
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
				Description: "Version name",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIdpsSignatureVersionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	// Get filter parameters
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	client := intrusion_services.NewSignatureVersionsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	var obj model.IdsSignatureVersion
	var err error

	if objID != "" {
		// Lookup by ID
		obj, err = client.Get(objID)
		if err != nil {
			return handleDataSourceReadError(d, "IdsSignatureVersion", objID, err)
		}
	} else if objName != "" {
		// List and find by display name
		versionsList, listErr := client.List(nil, nil, nil, nil, nil, nil)
		if listErr != nil {
			return handleListError("IdsSignatureVersion", listErr)
		}

		var matches []model.IdsSignatureVersion
		for _, version := range versionsList.Results {
			if version.DisplayName != nil && strings.EqualFold(*version.DisplayName, objName) {
				matches = append(matches, version)
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("IdsSignatureVersion with display_name '%s' not found", objName)
		}
		if len(matches) > 1 {
			return fmt.Errorf("found multiple IdsSignatureVersion with display_name '%s'", objName)
		}
		obj = matches[0]
	} else {
		return fmt.Errorf("Either 'id' or 'display_name' must be specified")
	}

	// Set all fields from API
	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
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
