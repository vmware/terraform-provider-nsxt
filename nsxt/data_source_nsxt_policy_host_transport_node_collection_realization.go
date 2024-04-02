/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/transport_node_collections"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyHostTransportNodeCollectionRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"timeout": {
				Type:         schema.TypeInt,
				Description:  "Realization timeout in seconds",
				Optional:     true,
				Default:      1200,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"delay": {
				Type:         schema.TypeInt,
				Description:  "Initial delay to start realization checks in seconds",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path of the site this Transport Node Collection belongs to",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Application state of transport node profile on compute collection",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyHostTransportNodeCollectionRealizationRead(d *schema.ResourceData, m interface{}) error {
	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}
	connector := getPolicyConnector(m)
	client := transport_node_collections.NewStateClient(connector)

	objID := d.Get("id").(string)
	delay := d.Get("delay").(int)
	timeout := d.Get("timeout").(int)
	objSitePath := d.Get("site_path").(string)
	objPolicyEnforcementPoint := getPolicyEnforcementPoint(m)

	pendingStates := []string{model.TransportNodeCollectionState_STATE_IN_PROGRESS}
	targetStates := []string{
		model.TransportNodeCollectionState_STATE_FAILED_TO_CREATE,
		model.TransportNodeCollectionState_STATE_FAILED_TO_REALIZE,
		model.TransportNodeCollectionState_STATE_PROFILE_MISMATCH,
		model.TransportNodeCollectionState_STATE_SUCCESS,
	}

	// For local manager, if site path is not provided, use default site
	if objSitePath == "" {
		objSitePath = defaultSite
	}
	if objID == "" {
		// Find the host transport node collection if exists
		tncClient := enforcement_points.NewTransportNodeCollectionsClient(connector)
		objList, err := tncClient.List(objSitePath, objPolicyEnforcementPoint, nil, nil, nil)
		if err != nil {
			return handleListError("HostTransportNodeCollection", err)
		}

		if len(objList.Results) != 1 {
			return fmt.Errorf("Single HostTransportNodeCollection was not identified on %s/%s", objSitePath, objPolicyEnforcementPoint)
		}
		objID = *objList.Results[0].Id
	}

	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			state, err := client.Get(objSitePath, getPolicyEnforcementPoint(m), objID)
			if err != nil {
				return state, model.TransportNodeCollectionState_STATE_IN_PROGRESS, logAPIError("Error while waiting for realization of Transport Node Collection", err)
			}

			log.Printf("[DEBUG] Current realization state for Transport Node Collection %s is %s", objID, *state.State)

			if *state.State != model.TransportNodeCollectionState_STATE_SUCCESS && *state.State != model.TransportNodeCollectionState_STATE_IN_PROGRESS {
				return state, *state.State, getErrorFromState(&state)
			}

			d.SetId(objID)
			d.Set("state", state.State)
			return state, *state.State, nil
		},
		Timeout:    time.Duration(timeout) * time.Second,
		MinTimeout: 1 * time.Second,
		Delay:      time.Duration(delay) * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return err
	}
	return nil
}

func getErrorFromState(state *model.TransportNodeCollectionState) error {
	var result string
	if state.ClusterLevelError != nil {
		result += fmt.Sprintf("cluster level error: %v\n", *state.ClusterLevelError)
	}
	for _, item := range state.ValidationErrors {
		result += fmt.Sprintf("validation error for node %s: %v\n", item.DiscoveredNodeId, item.ErrorMessage)
	}
	if state.VlcmTransitionError != nil {
		result += fmt.Sprintf("VCLM transition error: %v\n", *state.VlcmTransitionError)
	}

	return fmt.Errorf(result)
}
