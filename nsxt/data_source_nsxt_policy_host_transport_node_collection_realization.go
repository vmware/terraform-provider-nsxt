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
			"aggregate_progress_percentage": {
				Type:        schema.TypeInt,
				Description: "Aggregate percentage of compute collection deployment",
				Computed:    true,
			},
			"cluster_level_error": {
				Type:        schema.TypeString,
				Description: "Errors which needs cluster level to resolution",
				Optional:    true,
				Computed:    true,
			},
			"validation_errors": {
				Type:        schema.TypeList,
				Description: "Errors while applying transport node profile on discovered node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"discovered_node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"vlcm_transition_error": {
				Type:        schema.TypeString,
				Description: "Errors while enabling vLCM on the compute collection",
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

		if len(objList.Results) == 0 {
			return fmt.Errorf("HostTransportNodeCollection was not found on %s/%s", objSitePath, objPolicyEnforcementPoint)
		}
		d.SetId(*objList.Results[0].Id)
		d.Set("state", targetStates[3])
	}

	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			state, err := client.Get(objSitePath, getPolicyEnforcementPoint(m), objID)
			if err != nil {
				return state, model.ConfigurationState_STATE_ERROR, logAPIError("Error while waiting for realization of Transport Node Collection", err)
			}

			log.Printf("[DEBUG] Current realization state for Transport Node Collection %s is %s", objID, *state.State)

			d.SetId(objID)
			d.Set("state", state.State)
			d.Set("aggregate_progress_percentage", state.AggregateProgressPercentage)
			d.Set("cluster_level_error", state.ClusterLevelError)
			var validationErrorsList []map[string]interface{}
			for _, item := range state.ValidationErrors {
				data := make(map[string]interface{})
				data["discovered_node_id"] = item.DiscoveredNodeId
				data["error_message"] = item.ErrorMessage
				validationErrorsList = append(validationErrorsList, data)
			}
			d.Set("validation_errors", validationErrorsList)
			d.Set("vlcm_transition_error", state.VlcmTransitionError)

			return state, *state.State, nil
		},
		Timeout:    time.Duration(timeout) * time.Second,
		MinTimeout: 1 * time.Second,
		Delay:      time.Duration(delay) * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("failed to get realization information for %s: %v", objID, err)
	}
	return nil
}
