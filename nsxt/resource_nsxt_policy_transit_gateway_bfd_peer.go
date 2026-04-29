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

var cliTGWBfdPeersClient = transitgateways.NewBfdPeersClient

func resourceNsxtPolicyTransitGatewayBfdPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayBfdPeerCreate,
		Read:   resourceNsxtPolicyTransitGatewayBfdPeerRead,
		Update: resourceNsxtPolicyTransitGatewayBfdPeerUpdate,
		Delete: resourceNsxtPolicyTransitGatewayBfdPeerDelete,
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
			"bfd_profile_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the BFD profile to use for timing parameters of this BFD session. If not specified, default BFD timing parameters are used.",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable or disable this BFD peer configuration",
				Optional:    true,
				Default:     true,
			},
			"peer_address": {
				Type:         schema.TypeString,
				Description:  "IP address (IPv4 or IPv6) of the static route next-hop peer for BFD liveness detection",
				Required:     true,
				ValidateFunc: validateSingleIP(),
			},
			"source_attachment": {
				Type:        schema.TypeList,
				Description: "Policy path of the transit gateway attachment (with a cna_path) that serves as the BFD session source. Exactly one path must be specified.",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyPath(),
				},
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayBfdPeerExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWBfdPeersClient(sessionContext, connector)
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
	return false, logAPIError("Error retrieving TransitGatewayBfdPeer", err)
}

func resourceNsxtPolicyTransitGatewayBfdPeerCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	parentPath := d.Get("parent_path").(string)
	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayBfdPeerExists)
	if err != nil {
		return err
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := tgwBfdPeerFromSchema(d)

	log.Printf("[INFO] Creating TransitGatewayBfdPeer with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBfdPeersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayBfdPeer", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransitGatewayBfdPeerRead(d, m)
}

func resourceNsxtPolicyTransitGatewayBfdPeerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBfdPeer ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBfdPeersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayBfdPeer", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("bfd_profile_path", obj.BfdProfilePath)
	d.Set("enabled", obj.Enabled)
	d.Set("peer_address", obj.PeerAddress)
	d.Set("source_attachment", obj.SourceAttachment)

	return nil
}

func resourceNsxtPolicyTransitGatewayBfdPeerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBfdPeer ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	revision := int64(d.Get("revision").(int))
	obj := tgwBfdPeerFromSchema(d)
	obj.Revision = &revision

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBfdPeersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if _, err := client.Update(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayBfdPeer", id, err)
	}
	return resourceNsxtPolicyTransitGatewayBfdPeerRead(d, m)
}

func resourceNsxtPolicyTransitGatewayBfdPeerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBfdPeer ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBfdPeersClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayBfdPeer", id, err)
	}
	return nil
}

func tgwBfdPeerFromSchema(d *schema.ResourceData) model.TransitGatewayBfdPeer {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	enabled := d.Get("enabled").(bool)
	peerAddress := d.Get("peer_address").(string)
	sourceAttachment := getStringListFromSchemaList(d, "source_attachment")

	obj := model.TransitGatewayBfdPeer{
		DisplayName:      &displayName,
		Description:      &description,
		Tags:             tags,
		Enabled:          &enabled,
		PeerAddress:      &peerAddress,
		SourceAttachment: sourceAttachment,
	}

	if bfdProfilePath := d.Get("bfd_profile_path").(string); bfdProfilePath != "" {
		obj.BfdProfilePath = &bfdProfilePath
	}

	return obj
}
