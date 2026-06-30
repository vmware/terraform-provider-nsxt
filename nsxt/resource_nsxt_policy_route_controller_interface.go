// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliRCInterfaceClient = infra.NewRouteControllerInterfaceClient

var routeControllerInterfaceUrpfModeValues = []string{
	model.RouteControllerInterface_URPF_MODE_NONE,
	model.RouteControllerInterface_URPF_MODE_STRICT,
}

// routeControllerInterfacePathExample is the parent_path description used for parse validation.
// The parent_path for a route controller interface is the route controller path.
const routeControllerInterfacePathExample = routeControllerPathExample

func resourceNsxtPolicyRouteControllerInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerInterfaceCreate,
		Read:   resourceNsxtPolicyRouteControllerInterfaceRead,
		Update: resourceNsxtPolicyRouteControllerInterfaceUpdate,
		Delete: resourceNsxtPolicyRouteControllerInterfaceDelete,
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
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent Route Controller"),
			"mtu": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Maximum transmission unit specifies the size of the largest packet that a network protocol can transmit",
				ValidateFunc: validation.IntAtLeast(64),
			},
			"urpf_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Unicast Reverse Path Forwarding mode",
				ValidateFunc: validation.StringInSlice(routeControllerInterfaceUrpfModeValues, false),
				Default:      model.RouteControllerInterface_URPF_MODE_STRICT,
			},
			"floating_ip_subnets": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of floating IP subnets in CIDR notation (IP/prefix-length) for this interface. Required when the Route Controller HA mode is ACTIVE_STANDBY",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIPCidr(),
				},
			},
			"interface_address": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "List of interface address configurations",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnets": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of IP addresses and network prefixes in CIDR notation for this interface address",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateIPCidr(),
							},
						},
						"portgroup_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DV port group identifier discovered from vCenter",
						},
						"virtual_network_appliance_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Policy path for the virtual network appliance",
							ValidateFunc: validatePolicyPath(),
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyRouteControllerInterfaceExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
	if pathErr != nil {
		return false, pathErr
	}
	c := cliRCInterfaceClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := c.Get(parents[0], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving RouteControllerInterface", err)
}

// rcInterfaceSubnetsToStruct converts a list of CIDR strings to []model.InterfaceSubnet.
func rcInterfaceSubnetsToStruct(subnets []interface{}) []model.InterfaceSubnet {
	var result []model.InterfaceSubnet
	for _, s := range subnets {
		cidr := s.(string)
		parts := strings.Split(cidr, "/")
		if len(parts) != 2 {
			continue
		}
		prefix, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		prefix64 := int64(prefix)
		result = append(result, model.InterfaceSubnet{
			IpAddresses: []string{parts[0]},
			PrefixLen:   &prefix64,
		})
	}
	return result
}

// rcInterfaceSubnetsFromStruct converts []model.InterfaceSubnet back to CIDR strings.
func rcInterfaceSubnetsFromStruct(subnets []model.InterfaceSubnet) []string {
	var result []string
	for _, s := range subnets {
		if len(s.IpAddresses) > 0 && s.PrefixLen != nil {
			result = append(result, fmt.Sprintf("%s/%d", s.IpAddresses[0], *s.PrefixLen))
		}
	}
	return result
}

func rcInterfaceToStruct(d *schema.ResourceData, id string) model.RouteControllerInterface {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	urpfMode := d.Get("urpf_mode").(string)

	obj := model.RouteControllerInterface{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		UrpfMode:    &urpfMode,
		Id:          &id,
	}

	if mtu, ok := d.GetOk("mtu"); ok {
		mtu64 := int64(mtu.(int))
		obj.Mtu = &mtu64
	}

	floatingSubnets := d.Get("floating_ip_subnets").([]interface{})
	if len(floatingSubnets) > 0 {
		obj.FloatingIpSubnets = rcInterfaceSubnetsToStruct(floatingSubnets)
	}

	for _, raw := range d.Get("interface_address").([]interface{}) {
		data := raw.(map[string]interface{})
		addrObj := model.RouteControllerInterfaceAddress{
			InterfaceSubnet: rcInterfaceSubnetsToStruct(data["subnets"].([]interface{})),
		}
		if portgroupID, ok := data["portgroup_id"].(string); ok && portgroupID != "" {
			addrObj.PortgroupId = &portgroupID
		}
		if vnaPath, ok := data["virtual_network_appliance_path"].(string); ok && vnaPath != "" {
			addrObj.VirtualNetworkAppliancePath = &vnaPath
		}
		obj.InterfaceAddress = append(obj.InterfaceAddress, addrObj)
	}

	return obj
}

func setRCInterfaceInSchema(d *schema.ResourceData, obj model.RouteControllerInterface) {
	d.Set("urpf_mode", obj.UrpfMode)
	if obj.Mtu != nil {
		d.Set("mtu", int(*obj.Mtu))
	}

	d.Set("floating_ip_subnets", rcInterfaceSubnetsFromStruct(obj.FloatingIpSubnets))

	var addrList []interface{}
	for _, addr := range obj.InterfaceAddress {
		addrMap := make(map[string]interface{})
		addrMap["subnets"] = rcInterfaceSubnetsFromStruct(addr.InterfaceSubnet)

		portgroupID := ""
		if addr.PortgroupId != nil {
			portgroupID = *addr.PortgroupId
		}
		addrMap["portgroup_id"] = portgroupID

		vnaPath := ""
		if addr.VirtualNetworkAppliancePath != nil {
			vnaPath = *addr.VirtualNetworkAppliancePath
		}
		addrMap["virtual_network_appliance_path"] = vnaPath

		addrList = append(addrList, addrMap)
	}
	d.Set("interface_address", addrList)
}

func resourceNsxtPolicyRouteControllerInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyRouteControllerInterfaceExists)
	if err != nil {
		return err
	}

	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := rcInterfaceToStruct(d, id)

	log.Printf("[INFO] Creating RouteControllerInterface with ID %s", id)
	sessionContext := getSessionContext(d, m)
	c := cliRCInterfaceClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := c.Patch(parents[0], id, obj); err != nil {
		return handleCreateError("RouteControllerInterface", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
}

func resourceNsxtPolicyRouteControllerInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerInterface ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getSessionContext(d, m)
	c := cliRCInterfaceClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := c.Get(parents[0], id)
	if err != nil {
		return handleReadError(d, "RouteControllerInterface", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	setRCInterfaceInSchema(d, obj)

	return nil
}

func resourceNsxtPolicyRouteControllerInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerInterface ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := rcInterfaceToStruct(d, id)

	sessionContext := getSessionContext(d, m)
	c := cliRCInterfaceClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := c.Patch(parents[0], id, obj); err != nil {
		return handleUpdateError("RouteControllerInterface", id, err)
	}
	return resourceNsxtPolicyRouteControllerInterfaceRead(d, m)
}

func resourceNsxtPolicyRouteControllerInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerInterface ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getSessionContext(d, m)
	c := cliRCInterfaceClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := c.Delete(parents[0], id); err != nil {
		return handleDeleteError("RouteControllerInterface", id, err)
	}
	return nil
}
