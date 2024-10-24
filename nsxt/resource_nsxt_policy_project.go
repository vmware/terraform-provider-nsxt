/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs"
)

func resourceNsxtPolicyProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyProjectCreate,
		Read:   resourceNsxtPolicyProjectRead,
		Update: resourceNsxtPolicyProjectUpdate,
		Delete: resourceNsxtPolicyProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"short_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"site_info": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edge_cluster_paths": {
							Type:     schema.TypeList,
							Elem:     getElemPolicyPathSchemaWithFlags(false, false, false),
							Optional: true,
						},
						"site_path": getElemPolicyPathSchemaWithFlags(true, true, false),
					},
				},
				Optional: true,
				Computed: true,
			},
			"tier0_gateway_paths": {
				Type:     schema.TypeList,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
			"external_ipv4_blocks": {
				Type:     schema.TypeList,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
		},
	}
}

func resourceNsxtPolicyProjectExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error
	client := infra.NewProjectsClient(connector)
	_, err = client.Get(defaultOrgID, id, nil)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyProjectPatch(connector client.Connector, d *schema.ResourceData, m interface{}, id string) error {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	shortID := d.Get("short_id").(string)
	siteInfosList := d.Get("site_info").([]interface{})
	var siteInfos []model.SiteInfo
	for _, item := range siteInfosList {
		data := item.(map[string]interface{})
		var edgeClusterPaths []string
		for _, ecp := range data["edge_cluster_paths"].([]interface{}) {
			edgeClusterPaths = append(edgeClusterPaths, ecp.(string))
		}
		sitePath := data["site_path"].(string)
		obj := model.SiteInfo{
			EdgeClusterPaths: edgeClusterPaths,
			SitePath:         &sitePath,
		}
		siteInfos = append(siteInfos, obj)
	}
	tier0s := getStringListFromSchemaList(d, "tier0_gateway_paths")
	externalIPV4Blocks := getStringListFromSchemaList(d, "external_ipv4_blocks")

	obj := model.Project{
		DisplayName:        &displayName,
		Description:        &description,
		Tags:               tags,
		SiteInfos:          siteInfos,
		Tier0s:             tier0s,
		ExternalIpv4Blocks: externalIPV4Blocks,
	}

	if shortID != "" {
		obj.ShortId = &shortID
	}

	log.Printf("[INFO] Patching Project with ID %s", id)

	client := infra.NewProjectsClient(connector)
	return client.Patch(defaultOrgID, id, obj)
}

func resourceNsxtPolicyProjectCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyProjectExists)
	if err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	err = resourceNsxtPolicyProjectPatch(connector, d, m, id)
	if err != nil {
		return handleCreateError("Project", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyProjectRead(d, m)
}

func resourceNsxtPolicyProjectRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Project ID")
	}

	var obj model.Project
	client := infra.NewProjectsClient(connector)
	var err error
	obj, err = client.Get(defaultOrgID, id, nil)
	if err != nil {
		return handleReadError(d, "Project", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("short_id", obj.ShortId)
	var siteInfosList []map[string]interface{}
	for _, item := range obj.SiteInfos {
		data := make(map[string]interface{})
		data["edge_cluster_paths"] = item.EdgeClusterPaths
		data["site_path"] = item.SitePath
		siteInfosList = append(siteInfosList, data)
	}
	d.Set("site_info", siteInfosList)
	d.Set("tier0_gateway_paths", obj.Tier0s)
	d.Set("external_ipv4_blocks", obj.ExternalIpv4Blocks)

	return nil
}

func resourceNsxtPolicyProjectUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Project ID")
	}

	connector := getPolicyConnector(m)
	err := resourceNsxtPolicyProjectPatch(connector, d, m, id)
	if err != nil {
		return handleUpdateError("Project", id, err)
	}

	return resourceNsxtPolicyProjectRead(d, m)
}

func resourceNsxtPolicyProjectDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Project ID")
	}

	connector := getPolicyConnector(m)
	var err error
	client := infra.NewProjectsClient(connector)
	err = client.Delete(defaultOrgID, id, nil)

	if err != nil {
		return handleDeleteError("Project", id, err)
	}

	return nil
}
