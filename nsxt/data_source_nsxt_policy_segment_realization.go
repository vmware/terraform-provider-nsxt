/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"time"
)

func dataSourceNsxtPolicySegmentRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySegmentRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The path for the policy segment",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the realized resource on hypervisors",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicySegmentRealizationRead(d *schema.ResourceData, m interface{}) error {
	// Read the realization info by the path, and wait till it is valid
	connector := getPolicyConnector(m)

	// Get the realization info of this resource
	path := d.Get("path").(string)

	// Dummy id, just because each data source needs one
	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	// verifying segment realization on hypervisor
	segmentID := getPolicyIDFromPath(path)
	enforcementPointPath := getPolicyEnforcementPointPath()
	client := segments.NewDefaultStateClient(connector)
	pendingStates := []string{model.SegmentConfigurationState_STATE_PENDING,
		model.SegmentConfigurationState_STATE_IN_PROGRESS,
		model.SegmentConfigurationState_STATE_IN_SYNC,
		model.SegmentConfigurationState_STATE_UNKNOWN}
	targetStates := []string{model.SegmentConfigurationState_STATE_SUCCESS,
		model.SegmentConfigurationState_STATE_FAILED,
		model.SegmentConfigurationState_STATE_ERROR,
		model.SegmentConfigurationState_STATE_ORPHANED}
	if toleratePartialSuccess {
		targetStates = append(targetStates, model.SegmentConfigurationState_STATE_PARTIAL_SUCCESS)
	}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			state, err := client.Get(segmentID, nil, nil, &enforcementPointPath, nil, nil, nil, nil, nil)
			if err != nil {
				return state, model.SegmentConfigurationState_STATE_ERROR, fmt.Errorf("Error while waiting for realization of segment %s: %v", segmentID, err)
			}

			d.Set("state", state.State)
			return state, *state.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 1 * time.Second,
		Delay:      1 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Failed to get realization information for %s: %v", path, err)
	}

	return nil
}
