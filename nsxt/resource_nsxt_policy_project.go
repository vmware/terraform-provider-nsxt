// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	infra2 "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
)

func resourceNsxtPolicyProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyProjectCreate,
		Read:   resourceNsxtPolicyProjectRead,
		Update: resourceNsxtPolicyProjectUpdate,
		Delete: resourceNsxtPolicyProjectDelete,
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
			"activate_default_dfw_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vc_folder": {
				Type:        schema.TypeBool,
				Description: "Flag to specify whether the DVPGs created for project segments are grouped under a folder on the VC",
				Optional:    true,
			},
			"external_ipv4_blocks": {
				Type:     schema.TypeList,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
			"tgw_external_connections": {
				Type:     schema.TypeList,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
			"default_security_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"north_south_firewall": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Flag that indicates whether north-south firewall (Gateway Firewall) is enabled",
									},
								},
							},
						},
					},
				},
			},
			"default_span_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the Cluster based default Span object of type NetworkSpan",
				ValidateFunc: validatePolicyPath(),
				Optional:     true,
			},
			"non_default_span_paths": {
				Type:        schema.TypeList,
				Description: "List of non default policy paths of the Span objects of type NetworkSpan",
				MaxItems:    10,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyPath(),
				},
			},
			"zone_external_ids": {
				Type:        schema.TypeList,
				Description: "An array of Zone object's external IDs",
				Optional:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	activateDefaultDFWRules := d.Get("activate_default_dfw_rules").(bool)
	vcFolder := d.Get("vc_folder").(bool)
	extIpv4BlocksList := getStringListFromSchemaList(d, "external_ipv4_blocks")
	tgwConnectionsList := getStringListFromSchemaList(d, "tgw_external_connections")
	quotasList := getStringListFromSchemaList(d, "quotas")

	obj := model.Project{
		DisplayName:        &displayName,
		Description:        &description,
		Tags:               tags,
		SiteInfos:          siteInfos,
		Tier0s:             tier0s,
		ExternalIpv4Blocks: extIpv4BlocksList,
	}

	if util.NsxVersionHigherOrEqual("4.2.0") {
		obj.ActivateDefaultDfwRules = &activateDefaultDFWRules
	}

	if util.NsxVersionHigherOrEqual("4.1.1") {
		obj.ExternalIpv4Blocks = extIpv4BlocksList
	}

	if util.NsxVersionHigherOrEqual("9.0.0") {
		obj.TgwExternalConnections = tgwConnectionsList
		obj.VcFolder = &vcFolder
		obj.Limits = quotasList
	}

	if shortID != "" {
		obj.ShortId = &shortID
	}

	if util.NsxVersionHigherOrEqual("9.1.0") {
		// There should be just one object here
		var spanReferences []model.SpanReference
		defaultSpanPathinterface, isDefaultSet := d.GetOkExists("default_span_path")
		var defaultSpanPath string
		if !isDefaultSet {
			var err error
			defaultSpanPath, err = getDefaultSpan(connector)
			if err != nil {
				return err
			}
		} else {
			defaultSpanPath = defaultSpanPathinterface.(string)
		}
		// default_span_path will never be empty, since it has a default value and the validator will make sure that
		// user will not assign an empty string or such.
		isDefault := true
		spanReferences = append(spanReferences, model.SpanReference{
			SpanPath:  &defaultSpanPath,
			IsDefault: &isDefault,
		})

		spanRefs := d.Get("non_default_span_paths").([]interface{})
		for _, spanRef := range spanRefs {
			spanPath := spanRef.(string)
			isDefault := false
			spanReferences = append(spanReferences, model.SpanReference{
				SpanPath:  &spanPath,
				IsDefault: &isDefault,
			})
		}
		zoneExternalIds := interfaceListToStringList(d.Get("zone_external_ids").([]interface{}))

		vpcDeploymentScope := model.VpcDeploymentScope{
			SpanReferences:  spanReferences,
			ZoneExternalIds: zoneExternalIds,
		}
		obj.VpcDeploymentScope = &vpcDeploymentScope
	}

	log.Printf("[INFO] Patching Project with ID %s", id)

	client := infra.NewProjectsClient(connector)
	err := client.Patch(defaultOrgID, id, obj)
	if err != nil {
		return err
	}

	if d.HasChanges("default_security_profile") && util.NsxVersionHigherOrEqual("9.0.0") {
		err = patchVpcSecurityProfile(d, connector, id)
	}
	return err
}

func getDefaultSpan(connector client.Connector) (string, error) {
	var cursor *string
	client := infra2.NewNetworkSpansClient(connector)
	spanList, err := client.List(cursor, nil, nil, nil, nil, nil)
	if err != nil {
		return "", err
	}
	defaultSpanPath := []string{}
	for _, spanObj := range spanList.Results {
		if *spanObj.IsDefault {
			defaultSpanPath = append(defaultSpanPath, *spanObj.Path)
		}
	}

	return defaultSpanPath[0], nil
}

func patchVpcSecurityProfile(d *schema.ResourceData, connector client.Connector, projectID string) error {
	enabled := false
	defaultSecurityProfile := d.Get("default_security_profile")
	if defaultSecurityProfile != nil {
		dsp := defaultSecurityProfile.([]interface{})
		if len(dsp) > 0 {
			northSouthFirewall := dsp[0].(map[string]interface{})["north_south_firewall"]
			if northSouthFirewall != nil {
				nsfw := northSouthFirewall.([]interface{})
				if len(nsfw) > 0 {
					elem := nsfw[0].(map[string]interface{})
					enabled = elem["enabled"].(bool)
				}
			}
		}
	}

	client := projects.NewVpcSecurityProfilesClient(connector)
	objSec, err := client.Get(defaultOrgID, projectID, "default")
	if isNotFoundError(err) {
		return fmt.Errorf("failed to fetch the details of the VPC security profile: %v", err)
	}
	// Default security profile is created by NSX, we can assume that it's there already
	obj := model.VpcSecurityProfile{
		DisplayName: objSec.DisplayName,
		Description: objSec.Description,
		NorthSouthFirewall: &model.NorthSouthFirewall{
			Enabled: &enabled,
		},
	}
	return client.Patch(defaultOrgID, projectID, "default", obj)
}

func setVpcSecurityProfileInSchema(d *schema.ResourceData, connector client.Connector, projectID string) error {
	client := projects.NewVpcSecurityProfilesClient(connector)
	obj, err := client.Get(defaultOrgID, projectID, "default")
	if isNotFoundError(err) {
		return nil
	}
	if err != nil {
		return err
	}

	enabled := false
	if obj.NorthSouthFirewall != nil && obj.NorthSouthFirewall.Enabled != nil {
		enabled = *obj.NorthSouthFirewall.Enabled
	}

	nsfw := map[string]interface{}{"enabled": &enabled}
	nsfws := []interface{}{nsfw}
	dsp := map[string]interface{}{"north_south_firewall": nsfws}
	dsps := []interface{}{dsp}

	d.Set("default_security_profile", dsps)
	return nil
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
		return fmt.Errorf("error obtaining Project ID")
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

	if util.NsxVersionHigherOrEqual("9.0.0") {
		err = setVpcSecurityProfileInSchema(d, connector, id)
	}
	if err != nil {
		return err
	}
	if util.NsxVersionHigherOrEqual("4.2.0") {
		d.Set("activate_default_dfw_rules", obj.ActivateDefaultDfwRules)
	}
	if util.NsxVersionHigherOrEqual("4.1.1") {
		d.Set("external_ipv4_blocks", obj.ExternalIpv4Blocks)
	}
	if util.NsxVersionHigherOrEqual("9.0.0") {
		d.Set("tgw_external_connections", obj.TgwExternalConnections)
		d.Set("vc_folder", obj.VcFolder)
		d.Set("quotas", obj.Limits)
	}

	if util.NsxVersionHigherOrEqual("9.1.0") && obj.VpcDeploymentScope != nil {
		var nonDefaultSpanPaths []interface{}
		var defaultSpanRefs *string
		for _, spanRef := range obj.VpcDeploymentScope.SpanReferences {
			if *spanRef.IsDefault {
				defaultSpanRefs = spanRef.SpanPath
			} else {
				nonDefaultSpanPaths = append(nonDefaultSpanPaths, *spanRef.SpanPath)
			}
		}
		d.Set("default_span_path", defaultSpanRefs)
		d.Set("non_default_span_paths", nonDefaultSpanPaths)
		d.Set("zone_external_ids", stringList2Interface(obj.VpcDeploymentScope.ZoneExternalIds))

	}
	return nil
}

func resourceNsxtPolicyProjectUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Project ID")
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
		return fmt.Errorf("error obtaining Project ID")
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
