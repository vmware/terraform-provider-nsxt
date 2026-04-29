// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliTGWRouteMapsClient = transitgateways.NewRouteMapsClient

func resourceNsxtPolicyTransitGatewayRouteMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayRouteMapCreate,
		Read:   resourceNsxtPolicyTransitGatewayRouteMapRead,
		Update: resourceNsxtPolicyTransitGatewayRouteMapUpdate,
		Delete: resourceNsxtPolicyTransitGatewayRouteMapDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent Transit Gateway"),
			"entry": {
				Type:        schema.TypeList,
				Description: "List of Route Map entries",
				Required:    true,
				Elem:        getPolicyRouteMapEntrySchema(),
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayRouteMapExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving TransitGatewayRouteMap", err)
}

func tgwRouteMapBuildObj(d *schema.ResourceData, id string) model.TransitGatewayRouteMap {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	schemaEntries := d.Get("entry").([]interface{})
	var entries []model.RouteMapEntry
	for entryNo, entry := range schemaEntries {
		entries = append(entries, policyGatewayRouteMapBuildEntry(d, entryNo, entry.(map[string]interface{})))
	}

	return model.TransitGatewayRouteMap{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Entries:     entries,
	}
}

func resourceNsxtPolicyTransitGatewayRouteMapCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayRouteMapExists)
	if err != nil {
		return err
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := tgwRouteMapBuildObj(d, id)

	log.Printf("[INFO] Creating TransitGatewayRouteMap with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayRouteMap", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayRouteMapRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRouteMapRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayRouteMap", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setRouteMapEntriesInSchema(d, obj.Entries)

	return nil
}

func resourceNsxtPolicyTransitGatewayRouteMapUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := tgwRouteMapBuildObj(d, id)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating TransitGatewayRouteMap with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayRouteMap", id, err)
	}
	return resourceNsxtPolicyTransitGatewayRouteMapRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRouteMapDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayRouteMap ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWRouteMapsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayRouteMap", id, err)
	}
	return nil
}
