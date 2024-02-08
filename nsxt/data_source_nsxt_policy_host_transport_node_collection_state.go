/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/transport_node_collections"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var policyHostTransportNodeCollectionState = []string{
	model.TransportNodeCollectionState_STATE_FAILED_TO_CREATE,
	model.TransportNodeCollectionState_STATE_FAILED_TO_REALIZE,
	model.TransportNodeCollectionState_STATE_IN_PROGRESS,
	model.TransportNodeCollectionState_STATE_PROFILE_MISMATCH,
	model.TransportNodeCollectionState_STATE_SUCCESS,
}

func dataSourceNsxtPolicyHostTransportNodeCollectionState() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyHostTransportNodeCollectionStateRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path of the site this Transport Node Collection belongs to",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
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
			"state": {
				Type:        schema.TypeString,
				Description: "Application state of transport node profile on compute collection",
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

func dataSourceNsxtPolicyHostTransportNodeCollectionStateRead(d *schema.ResourceData, m interface{}) error {
	if isPolicyGlobalManager(m) {
		return localManagerOnlyError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	objSitePath := d.Get("site_path").(string)
	objState := model.TransportNodeCollectionState{State: &policyHostTransportNodeCollectionState[4]}
	// For local manager, if site path is not provided, use default site
	if objSitePath == "" {
		objSitePath = defaultSite
	}

	connector := getPolicyConnector(m)
	client := transport_node_collections.NewStateClient(connector)
	var err error
	if objID == "" {
		tncClient := enforcement_points.NewTransportNodeCollectionsClient(connector)
		objList, err := tncClient.List(objSitePath, getPolicyEnforcementPoint(m), nil, nil, nil)
		if err != nil {
			return handleListError("HostTransportNodeCollection", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == objName || strings.HasPrefix(*objInList.DisplayName, objName) || objName == "" {
				objID = *objInList.Id
				objName = *objInList.DisplayName
			}
		}
		if objID == "" {
			return fmt.Errorf("HostTransportNodeCollection '%s' was not found", objName)
		}
	} else {
		objState, err = client.Get(objSitePath, getPolicyEnforcementPoint(m), objID)
		if err != nil {
			return handleDataSourceReadError(d, "TransportNodeCollectionState", objID, err)
		}
	}

	d.SetId(objID)
	d.Set("display_name", objName)
	d.Set("site_path", objSitePath)
	// Computed fields
	d.Set("aggregate_progress_percentage", objState.AggregateProgressPercentage)
	d.Set("cluster_level_error", objState.ClusterLevelError)
	d.Set("state", objState.State)

	var validationErrorsList []map[string]interface{}
	for _, item := range objState.ValidationErrors {
		data := make(map[string]interface{})
		data["discovered_node_id"] = item.DiscoveredNodeId
		data["error_message"] = item.ErrorMessage
		validationErrorsList = append(validationErrorsList, data)
	}
	d.Set("validation_errors", validationErrorsList)
	d.Set("vlcm_transition_error", objState.VlcmTransitionError)
	return nil
}
