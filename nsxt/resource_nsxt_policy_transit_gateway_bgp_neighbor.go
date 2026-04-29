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

	tgwbgp "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/bgp"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliTGWBgpNeighborsClient = tgwbgp.NewNeighborsClient

func resourceNsxtPolicyTransitGatewayBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayBgpNeighborCreate,
		Read:   resourceNsxtPolicyTransitGatewayBgpNeighborRead,
		Update: resourceNsxtPolicyTransitGatewayBgpNeighborUpdate,
		Delete: resourceNsxtPolicyTransitGatewayBgpNeighborDelete,
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
			"source_attachment": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Policy path of the transit gateway attachment or IPSec route-based VPN session used as the BGP peering source",
				MaxItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyPath(),
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
			"neighbor_local_as_config": {
				Type:        schema.TypeList,
				Description: "BGP neighbor local-as configuration",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"as_path_modifier_type": {
							Type:         schema.TypeString,
							Description:  "AS_PATH modifier type for BGP local AS",
							Optional:     true,
							ValidateFunc: validation.StringInSlice(bgpNeighborLocalAsConfigAsPathModifierTypeValues, false),
						},
						"local_as_num": {
							Type:         schema.TypeString,
							Description:  "BGP neighbor local-as number in ASPLAIN/ASDOT Format",
							Required:     true,
							ValidateFunc: validateASPlainOrDot,
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyTransitGatewayBgpNeighborExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliTGWBgpNeighborsClient(sessionContext, connector)
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
	return false, logAPIError("Error retrieving TransitGatewayBgpNeighbor", err)
}

func tgwBgpNeighborToStruct(d *schema.ResourceData, id string) (model.BgpNeighborConfig, error) {
	var obj model.BgpNeighborConfig

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	allowAsIn := d.Get("allow_as_in").(bool)
	gracefulRestartMode := d.Get("graceful_restart_mode").(string)
	holdDownTime := int64(d.Get("hold_down_time").(int))
	keepAliveTime := int64(d.Get("keep_alive_time").(int))
	maximumHopLimit := int64(d.Get("maximum_hop_limit").(int))
	neighborAddress := d.Get("neighbor_address").(string)
	password := d.Get("password").(string)
	remoteAsNum := d.Get("remote_as_num").(string)
	sourceAttachment := interface2StringList(d.Get("source_attachment").([]interface{}))

	var bfdConfig *model.BgpBfdConfig
	for _, bfd := range d.Get("bfd_config").([]interface{}) {
		data := bfd.(map[string]interface{})
		enabled := data["enabled"].(bool)
		interval := int64(data["interval"].(int))
		multiple := int64(data["multiple"].(int))
		bfdConfig = &model.BgpBfdConfig{
			Enabled:  &enabled,
			Interval: &interval,
			Multiple: &multiple,
		}
		break
	}

	var rFilters []model.BgpRouteFiltering
	for _, filter := range d.Get("route_filtering").([]interface{}) {
		data := filter.(map[string]interface{})
		addrFamily := data["address_family"].(string)
		enabled := data["enabled"].(bool)

		filterStruct := model.BgpRouteFiltering{
			AddressFamily: &addrFamily,
			Enabled:       &enabled,
		}

		inFilter := data["in_route_filter"].(string)
		outFilter := data["out_route_filter"].(string)
		if len(inFilter) > 0 {
			filterStruct.InRouteFilters = []string{inFilter}
		}
		if len(outFilter) > 0 {
			filterStruct.OutRouteFilters = []string{outFilter}
		}
		if data["maximum_routes"] != 0 {
			maxRoutes := int64(data["maximum_routes"].(int))
			filterStruct.MaximumRoutes = &maxRoutes
		}

		rFilters = append(rFilters, filterStruct)
	}

	var neighborLocalAsConfig *model.BgpNeighborLocalAsConfig
	for _, localAsConfig := range d.Get("neighbor_local_as_config").([]interface{}) {
		data := localAsConfig.(map[string]interface{})
		localAsNum := data["local_as_num"].(string)
		neighborLocalAsConfig = &model.BgpNeighborLocalAsConfig{
			LocalAsNum: &localAsNum,
		}
		if asPathModifierType, ok := data["as_path_modifier_type"].(string); ok && asPathModifierType != "" {
			neighborLocalAsConfig.AsPathModifierType = &asPathModifierType
		}
		break
	}

	obj = model.BgpNeighborConfig{
		DisplayName:           &displayName,
		Description:           &description,
		Tags:                  tags,
		AllowAsIn:             &allowAsIn,
		Bfd:                   bfdConfig,
		GracefulRestartMode:   &gracefulRestartMode,
		HoldDownTime:          &holdDownTime,
		KeepAliveTime:         &keepAliveTime,
		MaximumHopLimit:       &maximumHopLimit,
		NeighborAddress:       &neighborAddress,
		RemoteAsNum:           &remoteAsNum,
		RouteFiltering:        rFilters,
		SourceAttachment:      sourceAttachment,
		NeighborLocalAsConfig: neighborLocalAsConfig,
		Id:                    &id,
	}

	if d.HasChange("password") {
		obj.Password = &password
	}

	return obj, nil
}

func resourceNsxtPolicyTransitGatewayBgpNeighborCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayBgpNeighborExists)
	if err != nil {
		return err
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := tgwBgpNeighborToStruct(d, id)
	if err != nil {
		return handleCreateError("TransitGatewayBgpNeighbor", id, err)
	}

	log.Printf("[INFO] Creating TransitGatewayBgpNeighbor with ID %s", id)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBgpNeighborsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Patch(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleCreateError("TransitGatewayBgpNeighbor", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayBgpNeighborRead(d, m)
}

func resourceNsxtPolicyTransitGatewayBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBgpNeighbor ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBgpNeighborsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayBgpNeighbor", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	// Reuse shared helper; source_attachment is always set for TGW neighbors.
	setBgpNeighborConfigInSchema(d, obj, true)
	// Clear source_addresses — not applicable for TGW neighbors.
	d.Set("source_addresses", nil)

	return nil
}

func resourceNsxtPolicyTransitGatewayBgpNeighborUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBgpNeighbor ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := tgwBgpNeighborToStruct(d, id)
	if err != nil {
		return handleUpdateError("TransitGatewayBgpNeighbor", id, err)
	}
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBgpNeighborsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if _, err := client.Update(parents[0], parents[1], parents[2], id, obj); err != nil {
		return handleUpdateError("TransitGatewayBgpNeighbor", id, err)
	}
	return resourceNsxtPolicyTransitGatewayBgpNeighborRead(d, m)
}

func resourceNsxtPolicyTransitGatewayBgpNeighborDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayBgpNeighbor ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, transitGatewayPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliTGWBgpNeighborsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	if err := client.Delete(parents[0], parents[1], parents[2], id); err != nil {
		return handleDeleteError("TransitGatewayBgpNeighbor", id, err)
	}
	return nil
}
