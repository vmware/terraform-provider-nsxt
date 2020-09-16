/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gm_locale_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s/locale_services"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Note: this resource is supported for global manager only
// TODO: consider supporting it for local manager and deprecating bgp config on T0 gateway
func resourceNsxtPolicyBgpConfig() *schema.Resource {
	bgpSchema := getPolicyBGPConfigSchema()
	bgpSchema["gateway_path"] = getPolicyPathSchema(true, true, "Gateway for this BGP config")
	bgpSchema["site_path"] = getPolicyPathSchema(true, true, "Site Path for this BGP config")

	return &schema.Resource{
		Create: resourceNsxtPolicyBgpConfigCreate,
		Read:   resourceNsxtPolicyBgpConfigRead,
		Update: resourceNsxtPolicyBgpConfigUpdate,
		Delete: resourceNsxtPolicyBgpConfigDelete,

		Schema: bgpSchema,
	}
}

func resourceNsxtPolicyBgpConfigRead(d *schema.ResourceData, m interface{}) error {
	if !isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}
	connector := getPolicyConnector(m)

	gwPath := d.Get("gateway_path").(string)
	gwID := getPolicyIDFromPath(gwPath)
	sitePath := d.Get("site_path").(string)
	if !isPolicyGlobalManager(m) {
		return fmt.Errorf("This resource is not supported for local manager")
	}

	serviceID, err := findTier0LocaleServiceForSite(connector, gwID, sitePath)
	if err != nil {
		return err
	}
	client := gm_locale_services.NewDefaultBgpClient(connector)
	gmObj, err := client.Get(gwID, serviceID)
	if err != nil {
		return handleReadError(d, "BGP Config", serviceID, err)
	}
	lmObj, convErr := convertModelBindingType(gmObj, gm_model.BgpRoutingConfigBindingType(), model.BgpRoutingConfigBindingType())
	if convErr != nil {
		return convErr
	}
	lmRoutingConfig := lmObj.(model.BgpRoutingConfig)

	data := initPolicyTier0BGPConfigMap(&lmRoutingConfig)

	for key, value := range data {
		d.Set(key, value)
	}

	return nil
}

func resourceNsxtPolicyBgpConfigToStruct(d *schema.ResourceData) (*gm_model.BgpRoutingConfig, error) {
	ecmp := d.Get("ecmp").(bool)
	enabled := d.Get("enabled").(bool)
	interSrIbgp := d.Get("inter_sr_ibgp").(bool)
	localAsNum := d.Get("local_as_num").(string)
	multipathRelax := d.Get("multipath_relax").(bool)
	restartMode := d.Get("graceful_restart_mode").(string)
	restartTimer := int64(d.Get("graceful_restart_timer").(int))
	staleTimer := int64(d.Get("graceful_restart_stale_route_timer").(int))
	tags := getPolicyGlobalManagerTagsFromSchema(d)

	var aggregationStructs []gm_model.RouteAggregationEntry
	routeAggregations := d.Get("route_aggregation").([]interface{})
	if len(routeAggregations) > 0 {
		for _, agg := range routeAggregations {
			data := agg.(map[string]interface{})
			prefix := data["prefix"].(string)
			summary := data["summary_only"].(bool)
			elem := gm_model.RouteAggregationEntry{
				Prefix:      &prefix,
				SummaryOnly: &summary,
			}

			aggregationStructs = append(aggregationStructs, elem)
		}
	}

	restartTimerStruct := gm_model.BgpGracefulRestartTimer{
		RestartTimer:    &restartTimer,
		StaleRouteTimer: &staleTimer,
	}

	restartConfigStruct := gm_model.BgpGracefulRestartConfig{
		Mode:  &restartMode,
		Timer: &restartTimerStruct,
	}

	routeStruct := gm_model.BgpRoutingConfig{
		Ecmp:                  &ecmp,
		Enabled:               &enabled,
		RouteAggregations:     aggregationStructs,
		Tags:                  tags,
		InterSrIbgp:           &interSrIbgp,
		LocalAsNum:            &localAsNum,
		MultipathRelax:        &multipathRelax,
		GracefulRestartConfig: &restartConfigStruct,
	}

	return &routeStruct, nil
}

func resourceNsxtPolicyBgpConfigCreate(d *schema.ResourceData, m interface{}) error {
	if !isPolicyGlobalManager(m) {
		return globalManagerOnlyError()
	}
	// This is not a create operation on NSX, since BGP config us auto created
	connector := getPolicyConnector(m)

	gwPath := d.Get("gateway_path").(string)
	gwID := getPolicyIDFromPath(gwPath)
	sitePath := d.Get("site_path").(string)

	obj, err := resourceNsxtPolicyBgpConfigToStruct(d)
	if err != nil {
		return handleCreateError("BgpRoutingConfig", gwID, err)
	}

	serviceID, err := findTier0LocaleServiceForSite(connector, gwID, sitePath)
	if err != nil {
		return handleCreateError("BgpRoutingConfig", gwID, err)
	}

	client := gm_locale_services.NewDefaultBgpClient(connector)
	err = client.Patch(gwID, serviceID, *obj)
	if err != nil {
		return handleCreateError("BgpRoutingConfig", gwID, err)
	}
	d.SetId("bgp")

	return resourceNsxtPolicyBgpConfigRead(d, m)
}

func resourceNsxtPolicyBgpConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	gwPath := d.Get("gateway_path").(string)
	gwID := getPolicyIDFromPath(gwPath)
	sitePath := d.Get("site_path").(string)
	revision := int64(d.Get("revision").(int))

	obj, err := resourceNsxtPolicyBgpConfigToStruct(d)
	if err != nil {
		return handleUpdateError("BgpRoutingConfig", gwID, err)
	}

	serviceID, err := findTier0LocaleServiceForSite(connector, gwID, sitePath)
	if err != nil {
		return handleUpdateError("BgpRoutingConfig", gwID, err)
	}

	obj.Revision = &revision

	client := gm_locale_services.NewDefaultBgpClient(connector)
	_, err = client.Update(gwID, serviceID, *obj)
	if err != nil {
		return handleUpdateError("BgpRoutingConfig", gwID, err)
	}

	return resourceNsxtPolicyBgpConfigRead(d, m)
}

func resourceNsxtPolicyBgpConfigDelete(d *schema.ResourceData, m interface{}) error {
	// BGP object can not be deleted as long as locale service exist
	// Delete call on NSX should revert settings to default, but this is
	// not supported on platform as of today
	// In order to revert manually, NSX requires us to disable BGP first
	// It makes more sense to avoid any action here (which mimics the computed
	// bgp behavior on T0 for local manager)

	return nil
}
