/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"strings"
)

var gatewayInterfaceTypeValues = []string{
	model.Tier0Interface_TYPE_SERVICE,
	model.Tier0Interface_TYPE_EXTERNAL,
	model.Tier0Interface_TYPE_LOOPBACK,
}

func resourceNsxtPolicyTier0GatewayInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTier0GatewayInterfaceCreate,
		Read:   resourceNsxtPolicyTier0GatewayInterfaceRead,
		Update: resourceNsxtPolicyTier0GatewayInterfaceUpdate,
		Delete: resourceNsxtPolicyTier0GatewayInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyTier0GatewayInterfaceImport,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":                 getNsxIDSchema(),
			"path":                   getPathSchema(),
			"display_name":           getDisplayNameSchema(),
			"description":            getDescriptionSchema(),
			"revision":               getRevisionSchema(),
			"tag":                    getTagsSchema(),
			"gateway_path":           getPolicyPathSchema(true, true, "Policy path for Tier0 gateway"),
			"segment_path":           getPolicyPathSchema(false, true, "Policy path for connected segment"),
			"subnets":                getGatewayInterfaceSubnetsSchema(),
			"mtu":                    getMtuSchema(),
			"ipv6_ndra_profile_path": getIPv6NDRAPathSchema(),
			"type": {
				Type:         schema.TypeString,
				Description:  "Interface Type",
				ValidateFunc: validation.StringInSlice(gatewayInterfaceTypeValues, false),
				Optional:     true,
				ForceNew:     true,
				Default:      model.Tier0Interface_TYPE_EXTERNAL,
			},
			"edge_node_path": getPolicyPathSchema(false, false, "Policy path for edge node"),
			"enable_pim": {
				Type:        schema.TypeBool,
				Description: "Enable Protocol Independent Multicast on Interface",
				Optional:    true,
				Default:     false,
			},
			"access_vlan_id": {
				Type:         schema.TypeInt,
				Description:  "Vlan ID",
				Optional:     true,
				ValidateFunc: validateVLANId,
			},
			"locale_service_id": {
				Type:        schema.TypeString,
				Description: "Id of associated Gateway Locale Service on NSX",
				Computed:    true,
			},
			"urpf_mode": getGatewayInterfaceUrpfModeSchema(),
			"ip_addresses": {
				Type:        schema.TypeList,
				Description: "Ip addresses",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func gatewayInterfaceVersionDepenantSet(d *schema.ResourceData, obj *model.Tier0Interface) {
	if nsxVersionLower("3.0.0") {
		return
	}
	interfaceType := d.Get("type").(string)
	// PIM config can only be configured on external interface
	if interfaceType == model.Tier0Interface_TYPE_EXTERNAL {
		enablePIM := d.Get("enable_pim").(bool)
		pimConfig := model.Tier0InterfacePimConfig{
			Enabled: &enablePIM,
		}
		obj.Multicast = &pimConfig
	}

	vlanID := int64(d.Get("access_vlan_id").(int))
	if vlanID > 0 {
		obj.AccessVlanId = &vlanID
	}

	urpfMode := d.Get("urpf_mode").(string)
	obj.UrpfMode = &urpfMode
}

func resourceNsxtPolicyTier0GatewayInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := locale_services.NewDefaultInterfacesClient(connector)

	id := d.Get("nsx_id").(string)
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)

	segmentPath := d.Get("segment_path").(string)
	ifType := d.Get("type").(string)
	if len(segmentPath) == 0 && ifType != model.Tier0Interface_TYPE_LOOPBACK {
		// segment_path in required for all interfaces other than loopback
		return fmt.Errorf("segment_path is mandatory for interface of type %s", ifType)
	}

	localeService, err := resourceNsxtPolicyTier0GatewayGetLocaleServiceEntry(tier0ID, connector)
	if err != nil {
		return err
	}
	if localeService == nil {
		return fmt.Errorf("Edge cluster is mandatory on gateway %s in order to create interfaces", tier0ID)
	}

	localeServiceID := *localeService.Id

	if id == "" {
		id = newUUID()
	} else {
		_, err := client.Get(tier0ID, localeServiceID, id)
		if err == nil {
			return fmt.Errorf("Interface with ID '%s' already exists on Tier0 Gateway %s", id, tier0ID)
		} else if !isNotFoundError(err) {
			return err
		}
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	interfaceSubnetList := getGatewayInterfaceSubnetList(d)
	var ipv6ProfilePaths []string
	if d.Get("ipv6_ndra_profile_path").(string) != "" {
		ipv6ProfilePaths = append(ipv6ProfilePaths, d.Get("ipv6_ndra_profile_path").(string))
	}
	mtu := int64(d.Get("mtu").(int))
	edgePath := d.Get("edge_node_path").(string)
	obj := model.Tier0Interface{
		Id:               &id,
		DisplayName:      &displayName,
		Description:      &description,
		Tags:             tags,
		Type_:            &ifType,
		Subnets:          interfaceSubnetList,
		Ipv6ProfilePaths: ipv6ProfilePaths,
	}

	if len(segmentPath) > 0 {
		obj.SegmentPath = &segmentPath
	}

	if mtu > 0 {
		obj.Mtu = &mtu
	}

	if edgePath != "" {
		obj.EdgePath = &edgePath
	}
	gatewayInterfaceVersionDepenantSet(d, &obj)

	// Create the resource using PATCH
	log.Printf("[INFO] Creating Tier0 interface with ID %s", id)
	err = client.Patch(tier0ID, localeServiceID, id, obj)
	if err != nil {
		return handleCreateError("Tier0 Interface", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	d.Set("locale_service_id", localeServiceID)

	return resourceNsxtPolicyTier0GatewayInterfaceRead(d, m)
}

func resourceNsxtPolicyTier0GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := locale_services.NewDefaultInterfacesClient(connector)

	id := d.Id()
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)
	if id == "" || tier0ID == "" {
		return fmt.Errorf("Error obtaining Tier0 Interface id")
	}

	obj, err := client.Get(tier0ID, defaultPolicyLocaleServiceID, id)
	if err != nil {
		return handleReadError(d, "Tier0 Interface", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("segment_path", obj.SegmentPath)
	d.Set("edge_node_path", obj.EdgePath)
	d.Set("type", obj.Type_)

	if obj.Ipv6ProfilePaths != nil {
		d.Set("ipv6_ndra_profile_path", obj.Ipv6ProfilePaths[0]) // only one supported for now
	}
	if obj.Mtu != nil {
		d.Set("mtu", *obj.Mtu)
	}

	if obj.Multicast != nil {
		d.Set("enable_pim", *obj.Multicast.Enabled)
	} else {
		d.Set("enable_pim", false)
	}

	if obj.AccessVlanId != nil {
		d.Set("access_vlan_id", *obj.AccessVlanId)
	}

	if obj.Subnets != nil {
		var subnetList []string
		var ipList []string
		for _, subnet := range obj.Subnets {
			cidr := fmt.Sprintf("%s/%d", subnet.IpAddresses[0], *subnet.PrefixLen)
			subnetList = append(subnetList, cidr)
			ipList = append(ipList, subnet.IpAddresses[0])
		}
		d.Set("subnets", subnetList)
		d.Set("ip_addresses", ipList)
	}

	if obj.UrpfMode != nil {
		d.Set("urpf_mode", obj.UrpfMode)
	} else {
		// assign default for version that is lower than 3.0.0
		d.Set("urpf_mode", model.Tier0Interface_URPF_MODE_STRICT)
	}

	return nil
}

func resourceNsxtPolicyTier0GatewayInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := locale_services.NewDefaultInterfacesClient(connector)

	id := d.Id()
	tier0Path := d.Get("gateway_path").(string)
	localeServiceID := d.Get("locale_service_id").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)
	if id == "" || tier0ID == "" || localeServiceID == "" {
		return fmt.Errorf("Error obtaining Tier0 id or Locale Service id")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	interfaceSubnetList := getGatewayInterfaceSubnetList(d)
	segmentPath := d.Get("segment_path").(string)
	var ipv6ProfilePaths []string
	if d.Get("ipv6_ndra_profile_path").(string) != "" {
		ipv6ProfilePaths = append(ipv6ProfilePaths, d.Get("ipv6_ndra_profile_path").(string))
	}
	mtu := int64(d.Get("mtu").(int))
	ifType := d.Get("type").(string)
	edgePath := d.Get("edge_node_path").(string)
	revision := int64(d.Get("revision").(int))
	obj := model.Tier0Interface{
		DisplayName:      &displayName,
		Description:      &description,
		Tags:             tags,
		Type_:            &ifType,
		Subnets:          interfaceSubnetList,
		Ipv6ProfilePaths: ipv6ProfilePaths,
		Revision:         &revision,
	}

	if mtu > 0 {
		obj.Mtu = &mtu
	}

	if edgePath != "" {
		obj.EdgePath = &edgePath
	}

	if len(segmentPath) > 0 {
		obj.SegmentPath = &segmentPath
	}

	gatewayInterfaceVersionDepenantSet(d, &obj)

	_, err := client.Update(tier0ID, localeServiceID, id, obj)
	if err != nil {
		return handleUpdateError("Tier0 Interface", id, err)
	}

	return resourceNsxtPolicyTier0GatewayInterfaceRead(d, m)
}

func resourceNsxtPolicyTier0GatewayInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := locale_services.NewDefaultInterfacesClient(connector)

	id := d.Id()
	tier0Path := d.Get("gateway_path").(string)
	tier0ID := getPolicyIDFromPath(tier0Path)
	localeServiceID := d.Get("locale_service_id").(string)
	if id == "" || tier0ID == "" || localeServiceID == "" {
		return fmt.Errorf("Error obtaining Tier0 id or Locale Service id")
	}

	err := client.Delete(tier0ID, localeServiceID, id)
	if err != nil {
		return handleDeleteError("Tier0 Interface", id, err)
	}

	return nil
}

func resourceNsxtPolicyTier0GatewayInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	if len(s) != 3 {
		return nil, fmt.Errorf("Please provide <gateway-id>/<locale-service-id>/<interface-id> as an input")
	}

	gwID := s[0]
	connector := getPolicyConnector(m)
	client := infra.NewDefaultTier0sClient(connector)
	tier0GW, err := client.Get(gwID)
	if err != nil {
		return nil, err
	}
	d.Set("gateway_path", tier0GW.Path)
	d.Set("locale_service_id", s[1])

	d.SetId(s[2])

	return []*schema.ResourceData{d}, nil

}
