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

var cliTGWCommunityListsClient = transitgateways.NewCommunityListsClient

func resourceNsxtPolicyTransitGatewayCommunityList() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayCommunityListCreate,
		Read:   resourceNsxtPolicyTransitGatewayCommunityListRead,
		Update: resourceNsxtPolicyTransitGatewayCommunityListUpdate,
		Delete: resourceNsxtPolicyTransitGatewayCommunityListDelete,
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
			"communities": {
				Type:        schema.TypeSet,
				Description: "List of BGP community entries",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyBGPCommunity,
				},
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayCommunityListExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWCommunityListsClient(sessionContext, connector)
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
	return false, logAPIError("Error retrieving TransitGatewayCommunityList", err)
}

func resourceNsxtPolicyTransitGatewayCommunityListCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayCommunityListExists)
	if err != nil {
		return err
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	communities := getStringListFromSchemaSet(d, "communities")

	obj := model.CommunityList{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Communities: communities,
	}

	log.Printf("[INFO] Creating TransitGatewayCommunityList with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWCommunityListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayCommunityList", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayCommunityListRead(d, m)
}

func resourceNsxtPolicyTransitGatewayCommunityListRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayCommunityList ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWCommunityListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayCommunityList", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("communities", obj.Communities)

	return nil
}

func resourceNsxtPolicyTransitGatewayCommunityListUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayCommunityList ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	communities := getStringListFromSchemaSet(d, "communities")
	revision := int64(d.Get("revision").(int))

	obj := model.CommunityList{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Communities: communities,
		Revision:    &revision,
	}

	log.Printf("[INFO] Updating TransitGatewayCommunityList with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWCommunityListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if _, err := client.Update(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayCommunityList", id, err)
	}
	return resourceNsxtPolicyTransitGatewayCommunityListRead(d, m)
}

func resourceNsxtPolicyTransitGatewayCommunityListDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayCommunityList ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWCommunityListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayCommunityList", id, err)
	}
	return nil
}
