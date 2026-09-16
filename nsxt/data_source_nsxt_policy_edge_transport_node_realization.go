// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	edge_transport_nodes "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points/edge_transport_nodes"
)

var cliEdgeTransportNodeStateClient = edge_transport_nodes.NewStateClient

func dataSourceNsxtPolicyEdgeTransportNodeRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyEdgeTransportNodeRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The path for the policy resource",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
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
		},
	}
}

func getEdgeTransportNodeKeysFromPath(path string) (string, string, string) {
	siteID := getResourceIDFromResourcePath(path, "sites")
	epID := getResourceIDFromResourcePath(path, "enforcement-points")
	id := getResourceIDFromResourcePath(path, "edge-transport-nodes")

	return siteID, epID, id
}

func getEdgeTnRealizationFailureMessage(edgeState *model.EdgeTnState) string {
	if edgeState == nil {
		return ""
	}

	var messages []string
	addMsg := func(msg *string) {
		if msg != nil {
			trimmed := strings.TrimSpace(*msg)
			if trimmed != "" {
				for _, m := range messages {
					if m == trimmed {
						return
					}
				}
				messages = append(messages, trimmed)
			}
		}
	}

	addDetails := func(details []model.ConfigurationStateElement) {
		for _, elem := range details {
			addMsg(elem.FailureMessage)
		}
	}

	addConfigState := func(cfgState *model.ConfigurationState) {
		if cfgState != nil {
			addMsg(cfgState.FailureMessage)
			addDetails(cfgState.Details)
		}
	}

	addMsg(edgeState.FailureMessage)
	addDetails(edgeState.Details)
	addConfigState(edgeState.DeploymentState)
	addConfigState(edgeState.TransportNodeState)

	return strings.Join(messages, ", ")
}

func dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliEdgeTransportNodeStateClient(sessionContext, connector)

	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	tnPath := d.Get("path").(string)
	delay := d.Get("delay").(int)
	timeout := d.Get("timeout").(int)

	siteID, epID, edgeID := getEdgeTransportNodeKeysFromPath(tnPath)

	pendingStates := []string{
		model.EdgeTnState_CONSOLIDATED_STATUS_UNINITIALIZED,
		model.EdgeTnState_CONSOLIDATED_STATUS_IN_PROGRESS,
		model.EdgeTnState_CONSOLIDATED_STATUS_SANDBOXED_REALIZATION_PENDING,
		model.EdgeTnState_CONSOLIDATED_STATUS_UNKNOWN,
	}
	targetStates := []string{
		model.EdgeTnState_CONSOLIDATED_STATUS_SUCCESS,
		model.EdgeTnState_CONSOLIDATED_STATUS_ERROR,
	}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			state, err := client.Get(siteID, epID, edgeID)
			if err != nil {
				return state, model.EdgeTnState_CONSOLIDATED_STATUS_ERROR, logAPIError("Error while waiting for realization of Transport Node", err)
			}

			if state.EdgeTnState == nil {
				log.Printf("[DEBUG] Realization of Transport Node is unknown")
				return state, model.EdgeTnState_CONSOLIDATED_STATUS_UNKNOWN, nil
			} else {
				log.Printf("[DEBUG] Current realization state for Transport Node %s is %s", id, *state.EdgeTnState.ConsolidatedStatus)
				if *state.EdgeTnState.ConsolidatedStatus == model.EdgeTnState_CONSOLIDATED_STATUS_ERROR {
					failureMsg := getEdgeTnRealizationFailureMessage(state.EdgeTnState)
					if failureMsg != "" {
						return state, *state.EdgeTnState.ConsolidatedStatus, fmt.Errorf("transport node %s failed to realize. %s", edgeID, failureMsg)
					}
					return state, *state.EdgeTnState.ConsolidatedStatus, fmt.Errorf("transport node %s failed to realize", edgeID)
				}

				return state, *state.EdgeTnState.ConsolidatedStatus, nil
			}
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
