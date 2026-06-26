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

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliRCBgpNeighborClient = infra.NewRouteControllerBgpNeighborClient

// routeControllerBgpPathExample is the policy path of a route controller's BGP config.
// It is used as parent_path for BGP neighbors.
const routeControllerBgpPathExample = "/infra/route-controllers/[controller]/bgp"

func resourceNsxtPolicyRouteControllerBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerBgpNeighborCreate,
		Read:   resourceNsxtPolicyRouteControllerBgpNeighborRead,
		Update: resourceNsxtPolicyRouteControllerBgpNeighborUpdate,
		Delete: resourceNsxtPolicyRouteControllerBgpNeighborDelete,
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
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent Route Controller BGP config"),
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Flag to enable or disable BGP peering with this neighbor",
			},
			"allow_as_in": {
				Description: "Flag to enable allow_as_in option for BGP neighbor",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"bfd_config": {
				Type:        schema.TypeList,
				Description: "BFD configuration for failure detection",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Flag to enable/disable BFD configuration",
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      500,
							Description:  "Time interval between heartbeat packets in milliseconds",
							ValidateFunc: validation.IntBetween(50, 60000),
						},
						"multiple": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3,
							Description:  "Number of times heartbeat packet is missed before BFD declares the neighbor is down",
							ValidateFunc: validation.IntBetween(2, 16),
						},
					},
				},
			},
			"gateway_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Next-hop gateway IP addresses used to reach non-directly connected BGP peers",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			"graceful_restart_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(bgpNeighborConfigGracefulRestartModeValues, false),
				Optional:     true,
				Description:  "BGP Graceful Restart Configuration Mode",
				Default:      model.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY,
			},
			"hold_down_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      180,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Wait time in seconds before declaring peer dead",
			},
			"keep_alive_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Interval between keep alive messages sent to peer",
			},
			"maximum_hop_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 255),
				Description:  "Maximum number of hops allowed to reach BGP neighbor",
			},
			"neighbor_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Neighbor IP Address",
				ValidateFunc: validateSingleIP(),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Password for BGP neighbor authentication",
				ValidateFunc: validation.StringLenBetween(0, 32),
				Sensitive:    true,
			},
			"remote_as_num": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ASN of the neighbor in ASPLAIN or ASDOT Format",
				ValidateFunc: validateASPlainOrDot,
			},
			"source_addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Source IP Addresses for BGP peering",
				MaxItems:    8,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			"route_filtering": {
				Type:        schema.TypeList,
				Description: "Enable address families and route filtering in each direction",
				Optional:    true,
				Computed:    true,
				MaxItems:    2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_family": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Address family type",
							ValidateFunc: validation.StringInSlice(bgpNeighborConfigRouteFilteringAddressFamilyValues, false),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Flag to enable/disable address family",
						},
						"in_route_filter": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Prefix-list or route map path for IN direction",
							ValidateFunc: validatePolicyPath(),
						},
						"maximum_routes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1000000),
							Description:  "Maximum number of routes for the address family",
						},
						"out_route_filter": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Prefix-list or route map path for OUT direction",
							ValidateFunc: validatePolicyPath(),
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyRouteControllerBgpNeighborExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	c := cliRCBgpNeighborClient(sessionContext, connector)
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
	return false, logAPIError("Error retrieving RouteControllerBgpNeighbor", err)
}

func rcBgpNeighborToStruct(d *schema.ResourceData, id string) (model.RouteControllerBgpNeighborConfig, error) {
	var obj model.RouteControllerBgpNeighborConfig

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	enabled := d.Get("enabled").(bool)
	allowAsIn := d.Get("allow_as_in").(bool)
	gracefulRestartMode := d.Get("graceful_restart_mode").(string)
	holdDownTime := int64(d.Get("hold_down_time").(int))
	keepAliveTime := int64(d.Get("keep_alive_time").(int))
	maximumHopLimit := int64(d.Get("maximum_hop_limit").(int))
	neighborAddress := d.Get("neighbor_address").(string)
	password := d.Get("password").(string)
	remoteAsNum := d.Get("remote_as_num").(string)
	sourceAddresses := interface2StringList(d.Get("source_addresses").([]interface{}))
	gatewayIPs := interface2StringList(d.Get("gateway_ips").([]interface{}))

	var bfdConfig *model.BgpBfdConfig
	for _, bfd := range d.Get("bfd_config").([]interface{}) {
		data := bfd.(map[string]interface{})
		bfdEnabled := data["enabled"].(bool)
		interval := int64(data["interval"].(int))
		multiple := int64(data["multiple"].(int))
		bfdConfig = &model.BgpBfdConfig{
			Enabled:  &bfdEnabled,
			Interval: &interval,
			Multiple: &multiple,
		}
		break
	}

	var rFilters []model.BgpRouteFiltering
	for _, filter := range d.Get("route_filtering").([]interface{}) {
		data := filter.(map[string]interface{})
		addrFamily := data["address_family"].(string)
		filterEnabled := data["enabled"].(bool)

		filterStruct := model.BgpRouteFiltering{
			AddressFamily: &addrFamily,
			Enabled:       &filterEnabled,
		}

		inFilter := data["in_route_filter"].(string)
		outFilter := data["out_route_filter"].(string)
		if len(inFilter) > 0 {
			filterStruct.InRouteFilters = []string{inFilter}
		}
		if len(outFilter) > 0 {
			filterStruct.OutRouteFilters = []string{outFilter}
		}
		if v, ok := data["maximum_routes"]; ok && v.(int) != 0 {
			maxRoutes := int64(v.(int))
			filterStruct.MaximumRoutes = &maxRoutes
		}

		rFilters = append(rFilters, filterStruct)
	}

	obj = model.RouteControllerBgpNeighborConfig{
		DisplayName:         &displayName,
		Description:         &description,
		Tags:                tags,
		Enabled:             &enabled,
		AllowAsIn:           &allowAsIn,
		Bfd:                 bfdConfig,
		GracefulRestartMode: &gracefulRestartMode,
		HoldDownTime:        &holdDownTime,
		KeepAliveTime:       &keepAliveTime,
		MaximumHopLimit:     &maximumHopLimit,
		NeighborAddress:     &neighborAddress,
		RemoteAsNum:         &remoteAsNum,
		RouteFiltering:      rFilters,
		Id:                  &id,
	}

	if len(sourceAddresses) > 0 {
		obj.SourceAddresses = sourceAddresses
	}
	if len(gatewayIPs) > 0 {
		obj.GatewayIps = gatewayIPs
	}

	if d.HasChange("password") {
		obj.Password = &password
	}

	return obj, nil
}

func setRCBgpNeighborConfigInSchema(d *schema.ResourceData, obj model.RouteControllerBgpNeighborConfig) {
	d.Set("enabled", obj.Enabled)
	d.Set("allow_as_in", obj.AllowAsIn)
	d.Set("graceful_restart_mode", obj.GracefulRestartMode)
	if obj.HoldDownTime != nil {
		d.Set("hold_down_time", int(*obj.HoldDownTime))
	}
	if obj.KeepAliveTime != nil {
		d.Set("keep_alive_time", int(*obj.KeepAliveTime))
	}
	if obj.MaximumHopLimit != nil {
		d.Set("maximum_hop_limit", int(*obj.MaximumHopLimit))
	}
	d.Set("neighbor_address", obj.NeighborAddress)
	d.Set("remote_as_num", obj.RemoteAsNum)
	d.Set("source_addresses", obj.SourceAddresses)
	d.Set("gateway_ips", obj.GatewayIps)

	var bfdConfigs []interface{}
	if obj.Bfd != nil {
		bfd := make(map[string]interface{})
		bfd["enabled"] = obj.Bfd.Enabled
		if obj.Bfd.Interval != nil {
			bfd["interval"] = int(*obj.Bfd.Interval)
		}
		if obj.Bfd.Multiple != nil {
			bfd["multiple"] = int(*obj.Bfd.Multiple)
		}
		bfdConfigs = append(bfdConfigs, bfd)
	}
	d.Set("bfd_config", bfdConfigs)

	var rFilters []interface{}
	for _, filter := range obj.RouteFiltering {
		rf := make(map[string]interface{})
		rf["address_family"] = filter.AddressFamily
		rf["enabled"] = filter.Enabled

		inFilter := ""
		outFilter := ""
		if len(filter.InRouteFilters) > 0 {
			inFilter = filter.InRouteFilters[0]
		}
		if len(filter.OutRouteFilters) > 0 {
			outFilter = filter.OutRouteFilters[0]
		}
		rf["in_route_filter"] = inFilter
		rf["out_route_filter"] = outFilter
		if filter.MaximumRoutes != nil {
			rf["maximum_routes"] = int(*filter.MaximumRoutes)
		}
		rFilters = append(rFilters, rf)
	}
	d.Set("route_filtering", rFilters)
}

func resourceNsxtPolicyRouteControllerBgpNeighborCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyRouteControllerBgpNeighborExists)
	if err != nil {
		return err
	}

	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := rcBgpNeighborToStruct(d, id)
	if err != nil {
		return handleCreateError("RouteControllerBgpNeighbor", id, err)
	}

	log.Printf("[INFO] Creating RouteControllerBgpNeighbor with ID %s", id)
	sessionContext := getSessionContext(d, m)
	c := cliRCBgpNeighborClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := c.Patch(parents[0], id, obj); err != nil {
		return handleCreateError("RouteControllerBgpNeighbor", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerBgpNeighbor ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getSessionContext(d, m)
	c := cliRCBgpNeighborClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := c.Get(parents[0], id)
	if err != nil {
		return handleReadError(d, "RouteControllerBgpNeighbor", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	setRCBgpNeighborConfigInSchema(d, obj)

	return nil
}

func resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerBgpNeighbor ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := rcBgpNeighborToStruct(d, id)
	if err != nil {
		return handleUpdateError("RouteControllerBgpNeighbor", id, err)
	}
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	sessionContext := getSessionContext(d, m)
	c := cliRCBgpNeighborClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if _, err := c.Update(parents[0], id, obj); err != nil {
		return handleUpdateError("RouteControllerBgpNeighbor", id, err)
	}
	return resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
}

func resourceNsxtPolicyRouteControllerBgpNeighborDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining RouteControllerBgpNeighbor ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getSessionContext(d, m)
	c := cliRCBgpNeighborClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := c.Delete(parents[0], id); err != nil {
		return handleDeleteError("RouteControllerBgpNeighbor", id, err)
	}
	return nil
}
