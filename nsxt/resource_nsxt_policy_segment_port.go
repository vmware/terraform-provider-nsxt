package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNsxtPolicySegmentPort() *schema.Resource {
	return &schema.Resource{
		Create:   resourceNsxtPolicySegmentPortCreate,
		Read:     resourceNsxtPolicySegmentPortRead,
		Update:   resourceNsxtPolicySegmentPortUpdate,
		Delete:   resourceNsxtPolicySegmentPortDelete,
		Importer: &schema.ResourceImporter{}, // TODO: Add importer
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Path of the segment",
				Optional:    true,
			},
			"attachment": {
				Type:        schema.TypeList,
				Description: "VIF attachment",
				Optional:    true,
				Elem:        getPolicySegmentPortSchema(),
				MaxItems:    1,
			},
			"vif_id": {
				Type:        schema.TypeString,
				Description: "Segment Port attachment id",
				Optional:    true,
			},
			"discovery_profile": {
				Type:        schema.TypeList,
				Description: "IP and MAC discovery profiles for this segment",
				Elem:        getPolicySegmentDiscoveryProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"qos_profile": {
				Type:        schema.TypeList,
				Description: "QoS profiles for this segment",
				Elem:        getPolicySegmentQosProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"security_profile": {
				Type:        schema.TypeList,
				Description: "Security profiles for this segment",
				Elem:        getPolicySegmentSecurityProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
		},
	}
}

var allocateAddresses = []string{"IP_POOL", "MAC_POOL", "BOTH", "DHCP", "DHCPV6", "SLAAC", "NONE"}

func getPolicySegmentPortSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allocate_addresses": {
				Type:         schema.TypeString,
				Description:  "Indicate how IP will be allocated for the port. Allowed values are IP_POOL, MAC_POOL, BOTH, DHCP, DHCPV6, SLAAC, NONE",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(allocateAddresses, false),
			},
			"app_id": {
				Type:        schema.TypeString,
				Description: "ID used to identify/look up a child attachment behind a parent attachment",
				Optional:    true,
			},
			"evpn_vlans": {
				Type:        schema.TypeList,
				Description: "Evpn tenant VLAN IDs the Parent logical-port serves.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				MaxItems: 1000,
			},
			"hyperbus_mode": {
				Type:         schema.TypeString,
				Description:  "ID used to identify/look up a child attachment behind a parent attachment",
				Optional:     true,
				Default:      "DISABLE",
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Type of port attachment. PARENT type is automatically set if evpn_vlans or hyperbus_mode is configured. INDEPENDENT type is automatically set for ports that belong to Segment of type DVPortgroup. STATIC type is deprecated.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PARENT", "CHILD", "INDEPENDENT", "STATIC"}, false),
			},
			"id": {
				Type:        schema.TypeString,
				Description: "VIF UUID on NSX Manager. If the attachement type is PARENT, this property is required.",
				Optional:    true,
			},
			"traffic_tag": {
				Type:        schema.TypeInt,
				Description: "VIF UUID on NSX Manager. If the attachement type is PARENT, this property is required.",
				Optional:    true,
			},
		},
	}
}

func resourceNsxtPolicySegmentPortCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicySegmentPortExists(d, context, connector))
	if err != nil {
		return err
	}

	obj, err := policySegmentPortResourceToInfraStruct(context, id, d, false)
	if err != nil {
		return err
	}

	err = policyInfraPatch(context, obj, connector, false)
	if err != nil {
		return handleCreateError("SegmentPort", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySegmentPortRead(d, m)
}

func resourceNsxtPolicySegmentPortRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	segmentPath := d.Get("segment_path").(string)
	id := d.Id()
	segPort, err := getSegmentPort(segmentPath, id, getSessionContext(d, m), connector)
	if err != nil {
		return fmt.Errorf("Error getting Segment Port : %v", err)
	}

	d.Set("display_name", segPort.DisplayName)
	d.Set("description", segPort.Description)
	setPolicyTagsInSchema(d, segPort.Tags)
	d.Set("nsx_id", id)
	d.Set("path", segPort.Path)
	d.Set("revision", segPort.Revision)
	err = nsxtPolicySegmentPortProfilesRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceNsxtPolicySegmentPortUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Segment ID")
	}

	obj, err := policySegmentPortResourceToInfraStruct(context, id, d, false)
	if err != nil {
		return err
	}

	err = policyInfraPatch(context, obj, connector, false)
	if err != nil {
		return handleCreateError("SegmentPort", id, err)
	}

	return resourceNsxtPolicySegmentPortRead(d, m)
}

func resourceNsxtPolicySegmentPortDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Segment Port ID")
	}
	obj, err := policySegmentPortResourceToInfraStruct(getSessionContext(d, m), id, d, true)
	if err != nil {
		return err
	}

	err = policyInfraPatch(getSessionContext(d, m), obj, getPolicyConnector(m), false)
	if err != nil {
		return handleCreateError("SegmentPort", id, err)
	}

	return nil
}
