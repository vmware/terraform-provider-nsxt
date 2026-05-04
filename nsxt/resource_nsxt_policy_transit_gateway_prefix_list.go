// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliTGWPrefixListsClient = transitgateways.NewPrefixListsClient

func resourceNsxtPolicyTransitGatewayPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayPrefixListCreate,
		Read:   resourceNsxtPolicyTransitGatewayPrefixListRead,
		Update: resourceNsxtPolicyTransitGatewayPrefixListUpdate,
		Delete: resourceNsxtPolicyTransitGatewayPrefixListDelete,
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
			"prefix": {
				Type:        schema.TypeList,
				Description: "Ordered list of network prefixes",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Action for the prefix list entry",
							ValidateFunc: validation.StringInSlice(gatewayPrefixListActionTypeValues, false),
							Default:      model.PrefixEntry_ACTION_PERMIT,
						},
						"ge": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Prefix length greater than or equal to",
							ValidateFunc: validation.IntBetween(0, 128),
							Default:      0,
						},
						"le": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Prefix length less than or equal to",
							ValidateFunc: validation.IntBetween(0, 128),
							Default:      0,
						},
						"network": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Network prefix in CIDR format. If not set, matches ANY network",
							ValidateFunc: validateCidr(),
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayPrefixListExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWPrefixListsClient(sessionContext, connector)
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
	return false, logAPIError("Error retrieving TransitGatewayPrefixList", err)
}

func resourceNsxtPolicyTransitGatewayPrefixListCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayPrefixListExists)
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
	prefixes := getPrefixesFromSchema(d)

	obj := model.PrefixList{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Prefixes:    prefixes,
	}

	log.Printf("[INFO] Creating TransitGatewayPrefixList with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWPrefixListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayPrefixList", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayPrefixListRead(d, m)
}

func resourceNsxtPolicyTransitGatewayPrefixListRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayPrefixList ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWPrefixListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayPrefixList", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setPrefixesInSchema(d, obj.Prefixes)

	return nil
}

func resourceNsxtPolicyTransitGatewayPrefixListUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayPrefixList ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	prefixes := getPrefixesFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.PrefixList{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Prefixes:    prefixes,
		Revision:    &revision,
	}

	log.Printf("[INFO] Updating TransitGatewayPrefixList with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWPrefixListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayPrefixList", id, err)
	}
	return resourceNsxtPolicyTransitGatewayPrefixListRead(d, m)
}

func resourceNsxtPolicyTransitGatewayPrefixListDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayPrefixList ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWPrefixListsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayPrefixList", id, err)
	}
	return nil
}
