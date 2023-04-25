/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicySegmentSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySegmentSecurityProfileCreate,
		Read:   resourceNsxtPolicySegmentSecurityProfileRead,
		Update: resourceNsxtPolicySegmentSecurityProfileUpdate,
		Delete: resourceNsxtPolicySegmentSecurityProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"bpdu_filter_allow": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsMACAddress,
				},
				Optional: true,
			},
			"bpdu_filter_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dhcp_client_block_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dhcp_client_block_v6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dhcp_server_block_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dhcp_server_block_v6_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"non_ip_traffic_block_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ra_guard_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rate_limit": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rx_broadcast": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"rx_multicast": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"tx_broadcast": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"tx_multicast": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"rate_limits_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNsxtPolicySegmentSecurityProfileExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewSegmentSecurityProfilesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewSegmentSecurityProfilesClient(connector)
		_, err = client.Get(id)
	}
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicySegmentSecurityProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	bpduFilterAllow := getStringListFromSchemaSet(d, "bpdu_filter_allow")
	bpduFilterEnable := d.Get("bpdu_filter_enable").(bool)
	dhcpClientBlockEnabled := d.Get("dhcp_client_block_enabled").(bool)
	dhcpClientBlockV6Enabled := d.Get("dhcp_client_block_v6_enabled").(bool)
	dhcpServerBlockEnabled := d.Get("dhcp_server_block_enabled").(bool)
	dhcpServerBlockV6Enabled := d.Get("dhcp_server_block_v6_enabled").(bool)
	nonIpTrafficBlockEnabled := d.Get("non_ip_traffic_block_enabled").(bool)
	raGuardEnabled := d.Get("ra_guard_enabled").(bool)
	rateLimitsList := d.Get("rate_limit").([]interface{})
	var rateLimits *model.TrafficRateLimits
	for _, item := range rateLimitsList {
		data := item.(map[string]interface{})
		rxBroadcast := int64(data["rx_broadcast"].(int))
		rxMulticast := int64(data["rx_multicast"].(int))
		txBroadcast := int64(data["tx_broadcast"].(int))
		txMulticast := int64(data["tx_multicast"].(int))
		obj := model.TrafficRateLimits{
			RxBroadcast: &rxBroadcast,
			RxMulticast: &rxMulticast,
			TxBroadcast: &txBroadcast,
			TxMulticast: &txMulticast,
		}
		rateLimits = &obj
		break
	}
	rateLimitsEnabled := d.Get("rate_limits_enabled").(bool)

	obj := model.SegmentSecurityProfile{
		DisplayName:              &displayName,
		Description:              &description,
		Tags:                     tags,
		BpduFilterAllow:          bpduFilterAllow,
		BpduFilterEnable:         &bpduFilterEnable,
		DhcpClientBlockEnabled:   &dhcpClientBlockEnabled,
		DhcpClientBlockV6Enabled: &dhcpClientBlockV6Enabled,
		DhcpServerBlockEnabled:   &dhcpServerBlockEnabled,
		DhcpServerBlockV6Enabled: &dhcpServerBlockV6Enabled,
		NonIpTrafficBlockEnabled: &nonIpTrafficBlockEnabled,
		RaGuardEnabled:           &raGuardEnabled,
		RateLimits:               rateLimits,
		RateLimitsEnabled:        &rateLimitsEnabled,
	}

	log.Printf("[INFO] Sending SegmentSecurityProfile with ID %s", id)
	if isPolicyGlobalManager(m) {
		gmObj, convErr := convertModelBindingType(obj, model.SegmentSecurityProfileBindingType(), gm_model.SegmentSecurityProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewSegmentSecurityProfilesClient(connector)
		return client.Patch(id, gmObj.(gm_model.SegmentSecurityProfile), nil)
	}

	client := infra.NewSegmentSecurityProfilesClient(connector)
	return client.Patch(id, obj, nil)
}

func resourceNsxtPolicySegmentSecurityProfileCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicySegmentSecurityProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtPolicySegmentSecurityProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("SegmentSecurityProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySegmentSecurityProfileRead(d, m)
}

func resourceNsxtPolicySegmentSecurityProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	var obj model.SegmentSecurityProfile
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewSegmentSecurityProfilesClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "SegmentSecurityProfile", id, err)
		}

		lmObj, err := convertModelBindingType(gmObj, gm_model.SegmentSecurityProfileBindingType(), model.SegmentSecurityProfileBindingType())
		if err != nil {
			return err
		}
		obj = lmObj.(model.SegmentSecurityProfile)
	} else {
		client := infra.NewSegmentSecurityProfilesClient(connector)
		var err error
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "SegmentSecurityProfile", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("bpdu_filter_allow", obj.BpduFilterAllow)
	d.Set("bpdu_filter_enable", obj.BpduFilterEnable)
	d.Set("dhcp_client_block_enabled", obj.DhcpClientBlockEnabled)
	d.Set("dhcp_client_block_v6_enabled", obj.DhcpClientBlockV6Enabled)
	d.Set("dhcp_server_block_enabled", obj.DhcpServerBlockEnabled)
	d.Set("dhcp_server_block_v6_enabled", obj.DhcpServerBlockV6Enabled)
	d.Set("non_ip_traffic_block_enabled", obj.NonIpTrafficBlockEnabled)
	d.Set("ra_guard_enabled", obj.RaGuardEnabled)
	d.Set("rate_limits_enabled", obj.RateLimitsEnabled)

	var rateLimitsList []map[string]interface{}
	if obj.RateLimits != nil {
		item := obj.RateLimits
		data := make(map[string]interface{})
		data["rx_broadcast"] = item.RxBroadcast
		data["rx_multicast"] = item.RxMulticast
		data["tx_broadcast"] = item.TxBroadcast
		data["tx_multicast"] = item.TxMulticast
		rateLimitsList = append(rateLimitsList, data)
	}
	d.Set("rate_limit", rateLimitsList)

	return nil
}

func resourceNsxtPolicySegmentSecurityProfileUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	err := resourceNsxtPolicySegmentSecurityProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("SegmentSecurityProfile", id, err)
	}

	return resourceNsxtPolicySegmentSecurityProfileRead(d, m)
}

func resourceNsxtPolicySegmentSecurityProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SegmentSecurityProfile ID")
	}

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewSegmentSecurityProfilesClient(connector)
		err = client.Delete(id, nil)
	} else {
		client := infra.NewSegmentSecurityProfilesClient(connector)
		err = client.Delete(id, nil)
	}

	if err != nil {
		return handleDeleteError("SegmentSecurityProfile", id, err)
	}

	return nil
}
