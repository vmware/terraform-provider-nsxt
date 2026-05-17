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

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

const cnaPathExample = "/orgs/default/projects/[project]/centralized-network-attachments/[id]"

var cliCNAClient = apiprojects.NewCentralizedNetworkAttachmentsClient

func resourceNsxtPolicyProjectCentralizedNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyProjectCNACreate,
		Read:   resourceNsxtPolicyProjectCNARead,
		Update: resourceNsxtPolicyProjectCNAUpdate,
		Delete: resourceNsxtPolicyProjectCNADelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(cnaPathExample),
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"subnet_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the VPC subnet to which this centralized network attachment is connected. Required and immutable after creation.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"advertise_outbound_networks": {
				Type:        schema.TypeList,
				Description: "Outbound route advertisement configuration",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_private": {
							Type:        schema.TypeBool,
							Description: "When true, disables VPC auto-SNAT/EIP translation and enables TGW_PRIVATE route redistribution on this connection.",
							Optional:    true,
							Default:     false,
						},
						"allow_external_blocks": {
							Type:        schema.TypeList,
							Description: "External IP blocks (CIDRs) used as advertisement filter for prefixes from the transit gateway.",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateIPCidr(),
							},
						},
					},
				},
			},
			"interface_subnet": {
				Type:        schema.TypeList,
				Description: "Manual IP assignment for per-node interfaces. Supports 1-2 entries (one IPv4 and/or one IPv6 subnet). When set, CNAs with manual IP assignment cannot be shared to other centralized transit gateways.",
				Optional:    true,
				MaxItems:    2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_ip_address": {
							Type:        schema.TypeList,
							Description: "IP addresses assigned to each edge node interface. Count must match the number of edge nodes.",
							Required:    true,
							MinItems:    1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateSingleIP(),
							},
						},
						"prefix_length": {
							Type:         schema.TypeInt,
							Description:  "Prefix length of the subnet (1-128). For IPv4: 1-32, for IPv6: 1-128.",
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 128),
						},
						"ha_vip_ip_address": {
							Type:         schema.TypeString,
							Description:  "Floating VIP assigned to the active edge node in ACTIVE_STANDBY HA mode.",
							Optional:     true,
							ValidateFunc: validateSingleIP(),
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyProjectCNAExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	client := cliCNAClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := client.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving CentralizedNetworkAttachment", err)
}

func getCNAAdvertiseOutboundFromSchema(d *schema.ResourceData) *model.AdvertiseOutboundNetworks {
	list := d.Get("advertise_outbound_networks").([]interface{})
	if len(list) == 0 {
		return nil
	}
	raw, ok := list[0].(map[string]interface{})
	if !ok || raw == nil {
		return nil
	}
	allowPrivate := raw["allow_private"].(bool)
	aon := &model.AdvertiseOutboundNetworks{
		AllowPrivate:        &allowPrivate,
		AllowExternalBlocks: interfaceListToStringList(raw["allow_external_blocks"].([]interface{})),
	}
	return aon
}

func setCNAAdvertiseOutboundInSchema(d *schema.ResourceData, aon *model.AdvertiseOutboundNetworks) error {
	if aon == nil {
		return d.Set("advertise_outbound_networks", nil)
	}
	m := map[string]interface{}{
		"allow_external_blocks": aon.AllowExternalBlocks,
		"allow_private":         false,
	}
	if aon.AllowPrivate != nil {
		m["allow_private"] = *aon.AllowPrivate
	}
	return d.Set("advertise_outbound_networks", []interface{}{m})
}

func getCNAInterfaceSubnetsFromSchema(d *schema.ResourceData) []model.CentralizedNetworkAttachmentInterfaceSubnet {
	list := d.Get("interface_subnet").([]interface{})
	var subnets []model.CentralizedNetworkAttachmentInterfaceSubnet
	for _, item := range list {
		raw, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		prefixLen := int64(raw["prefix_length"].(int))
		subnet := model.CentralizedNetworkAttachmentInterfaceSubnet{
			InterfaceIpAddress: interfaceListToStringList(raw["interface_ip_address"].([]interface{})),
			PrefixLength:       &prefixLen,
		}
		if haVip, ok := raw["ha_vip_ip_address"].(string); ok && haVip != "" {
			subnet.HaVipIpAddress = &haVip
		}
		subnets = append(subnets, subnet)
	}
	return subnets
}

func setCNAInterfaceSubnetsInSchema(d *schema.ResourceData, subnets []model.CentralizedNetworkAttachmentInterfaceSubnet) error {
	var result []interface{}
	for _, s := range subnets {
		m := map[string]interface{}{
			"interface_ip_address": s.InterfaceIpAddress,
			"prefix_length":        0,
			"ha_vip_ip_address":    "",
		}
		if s.PrefixLength != nil {
			m["prefix_length"] = int(*s.PrefixLength)
		}
		if s.HaVipIpAddress != nil {
			m["ha_vip_ip_address"] = *s.HaVipIpAddress
		}
		result = append(result, m)
	}
	return d.Set("interface_subnet", result)
}

func resourceNsxtPolicyProjectCNACreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyProjectCNAExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	subnetPath := d.Get("subnet_path").(string)

	obj := model.CentralizedNetworkAttachment{
		DisplayName:               &displayName,
		Description:               &description,
		Tags:                      tags,
		SubnetPath:                &subnetPath,
		AdvertiseOutboundNetworks: getCNAAdvertiseOutboundFromSchema(d),
		InterfaceSubnets:          getCNAInterfaceSubnetsFromSchema(d),
	}

	client := cliCNAClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	log.Printf("[INFO] Creating CentralizedNetworkAttachment with ID %s", id)
	if err := client.Patch(parents[0], parents[1], id, obj); err != nil {
		return handleCreateError("CentralizedNetworkAttachment", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyProjectCNARead(d, m)
}

func resourceNsxtPolicyProjectCNARead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining CentralizedNetworkAttachment ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	client := cliCNAClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	obj, err := client.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "CentralizedNetworkAttachment", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("subnet_path", obj.SubnetPath)

	if err := setCNAAdvertiseOutboundInSchema(d, obj.AdvertiseOutboundNetworks); err != nil {
		return handleReadError(d, "CentralizedNetworkAttachment", id, err)
	}
	return setCNAInterfaceSubnetsInSchema(d, obj.InterfaceSubnets)
}

func resourceNsxtPolicyProjectCNAUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining CentralizedNetworkAttachment ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	subnetPath := d.Get("subnet_path").(string)

	obj := model.CentralizedNetworkAttachment{
		DisplayName:               &displayName,
		Description:               &description,
		Tags:                      tags,
		Revision:                  &revision,
		SubnetPath:                &subnetPath,
		AdvertiseOutboundNetworks: getCNAAdvertiseOutboundFromSchema(d),
		InterfaceSubnets:          getCNAInterfaceSubnetsFromSchema(d),
	}

	client := cliCNAClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	if _, err := client.Update(parents[0], parents[1], id, obj); err != nil {
		return handleUpdateError("CentralizedNetworkAttachment", id, err)
	}
	return resourceNsxtPolicyProjectCNARead(d, m)
}

func resourceNsxtPolicyProjectCNADelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining CentralizedNetworkAttachment ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	client := cliCNAClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	log.Printf("[INFO] Deleting CentralizedNetworkAttachment with ID %s", id)
	if err := client.Delete(parents[0], parents[1], id); err != nil {
		return handleDeleteError("CentralizedNetworkAttachment", id, err)
	}
	return nil
}
