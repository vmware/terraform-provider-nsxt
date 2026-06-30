// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var cliVNAClusterStateClient = enforcement_points.NewVirtualNetworkApplianceClusterStateClient
var cliVNAsClient = enforcement_points.NewVirtualNetworkAppliancesClient
var cliVNAStateClient = enforcement_points.NewVirtualNetworkApplianceStateClient

// vnaClusterPendingStates are ConsolidatedStatus values indicating that
// provisioning is still in progress at the cluster or per-VNA level.
var vnaClusterPendingStates = []string{
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNINITIALIZED,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_IN_PROGRESS,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SANDBOXED_REALIZATION_PENDING,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNKNOWN,
}

// vnaClusterTerminalStates are status values that stop polling (success or
// any failure/degraded condition).
var vnaClusterTerminalStates = []string{
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SUCCESS,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DOWN,
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DEGRAGED, //nolint:misspell
	model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_DISABLED,
}

func dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The policy path of the VirtualNetworkApplianceCluster resource",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Current realization state of the cluster",
				Computed:    true,
			},
			"vna_paths": {
				Type:        schema.TypeList,
				Description: "Policy paths of the deployed Virtual Network Appliance nodes within the cluster",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timeout": {
				Type:         schema.TypeInt,
				Description:  "Realization timeout in seconds",
				Optional:     true,
				Default:      1800,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"delay": {
				Type:         schema.TypeInt,
				Description:  "Initial delay before starting realization checks, in seconds",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	clusterPath := d.Get("path").(string)
	delay := d.Get("delay").(int)
	timeout := d.Get("timeout").(int)

	siteID := getResourceIDFromResourcePath(clusterPath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining site ID from path %s", clusterPath)
	}
	epID := getResourceIDFromResourcePath(clusterPath, "enforcement-points")
	if epID == "" {
		return fmt.Errorf("error obtaining enforcement-point ID from path %s", clusterPath)
	}
	clusterID := getResourceIDFromResourcePath(clusterPath, "virtual-network-appliance-clusters")
	if clusterID == "" {
		return fmt.Errorf("error obtaining cluster ID from path %s", clusterPath)
	}

	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	clusterStateClient := cliVNAClusterStateClient(sessionContext, connector)
	if clusterStateClient == nil {
		return policyResourceNotSupportedError()
	}

	// Stage 1: poll the cluster-level state endpoint until ConsolidatedStatus
	// reaches a terminal value. This confirms the cluster configuration has
	// been accepted and VNA deployment has been attempted.
	stateConf := &resource.StateChangeConf{
		Pending: vnaClusterPendingStates,
		Target:  vnaClusterTerminalStates,
		Refresh: func() (interface{}, string, error) {
			clusterState, err := clusterStateClient.Get(siteID, epID, clusterID)
			if err != nil {
				return clusterState, model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR,
					logAPIError("Error while waiting for realization of VirtualNetworkApplianceCluster", err)
			}
			if clusterState.ConsolidatedStatus == nil {
				log.Printf("[DEBUG] VirtualNetworkApplianceCluster %s realization state is unknown", clusterID)
				return clusterState, model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_UNKNOWN, nil
			}
			status := *clusterState.ConsolidatedStatus
			log.Printf("[DEBUG] VirtualNetworkApplianceCluster %s realization state: %s", clusterID, status)
			d.Set("state", status)
			return clusterState, status, nil
		},
		Timeout:    time.Duration(timeout) * time.Second,
		MinTimeout: 5 * time.Second,
		Delay:      time.Duration(delay) * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed to get realization information for VirtualNetworkApplianceCluster %s: %v", clusterPath, err)
	}

	// Stage 2: for each VNA appliance that exists within the cluster, wait for
	// its individual state to reach a terminal value via the per-appliance
	// state endpoint (.../virtual-network-appliances/{vna-id}/state).
	// VNAs are only present after NSX has deployed them; if none exist this
	// stage is a no-op, which is the correct behaviour for environments where
	// VNA deployment has not yet occurred.
	vnaPaths, err := waitForVNAAppliancesRealization(siteID, epID, clusterID, timeout, sessionContext, connector)
	if err != nil {
		return err
	}
	d.Set("vna_paths", vnaPaths)

	return nil
}

// waitForVNAAppliancesRealization lists all deployed VNA appliances within the
// cluster, polls each one's /state endpoint until it reaches a terminal
// ConsolidatedStatus, and returns the policy paths of all deployed appliances.
// If no appliances are deployed yet the function is a no-op and returns nil.
func waitForVNAAppliancesRealization(siteID, epID, clusterID string, timeoutSec int, sessionContext utl.SessionContext, connector vapiProtocolClient.Connector) ([]string, error) {
	vnasClient := cliVNAsClient(sessionContext, connector)
	if vnasClient == nil {
		return nil, policyResourceNotSupportedError()
	}

	result, err := vnasClient.List(siteID, epID, clusterID)
	if err != nil {
		return nil, logAPIError("Error listing VNA appliances in cluster "+clusterID, err)
	}
	if len(result.Results) == 0 {
		log.Printf("[DEBUG] VirtualNetworkApplianceCluster %s has no deployed appliances yet; skipping per-VNA state check", clusterID)
		return nil, nil
	}

	vnaStateClient := cliVNAStateClient(sessionContext, connector)
	if vnaStateClient == nil {
		return nil, policyResourceNotSupportedError()
	}

	var vnaPaths []string
	for _, vna := range result.Results {
		if vna.Id == nil {
			continue
		}
		vnaID := *vna.Id
		if err := waitForSingleVNARealization(siteID, epID, clusterID, vnaID, timeoutSec, vnaStateClient); err != nil {
			return nil, err
		}
		if vna.Path != nil {
			vnaPaths = append(vnaPaths, *vna.Path)
		}
	}
	return vnaPaths, nil
}

// waitForSingleVNARealization polls the per-appliance state endpoint for a
// single VNA until its VirtualNetworkApplianceConfigurationState reaches a
// terminal ConsolidatedStatus.
func waitForSingleVNARealization(siteID, epID, clusterID, vnaID string, timeoutSec int, vnaStateClient *enforcement_points.VirtualNetworkApplianceStateClientContext) error {
	// The per-VNA ConfigurationState shares the same status constants as the
	// cluster-level state, so reuse the same pending/terminal slices.
	pendingStates := []string{
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_UNINITIALIZED,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_IN_PROGRESS,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_SANDBOXED_REALIZATION_PENDING,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_UNKNOWN,
	}
	targetStates := []string{
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_SUCCESS,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_ERROR,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_DOWN,
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_DEGRAGED, //nolint:misspell
		model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_DISABLED,
	}

	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			applianceState, err := vnaStateClient.Get(siteID, epID, clusterID, vnaID)
			if err != nil {
				return applianceState, model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_ERROR,
					logAPIError(fmt.Sprintf("Error while waiting for realization of VNA %s", vnaID), err)
			}
			vnaState := applianceState.VirtualNetworkApplianceState
			if vnaState == nil || vnaState.ConfigurationState == nil || vnaState.ConfigurationState.ConsolidatedStatus == nil {
				log.Printf("[DEBUG] VNA %s realization state is unknown", vnaID)
				return applianceState, model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_UNKNOWN, nil
			}
			status := *vnaState.ConfigurationState.ConsolidatedStatus
			log.Printf("[DEBUG] VNA %s realization state: %s", vnaID, status)
			return applianceState, status, nil
		},
		Timeout:    time.Duration(timeoutSec) * time.Second,
		MinTimeout: 5 * time.Second,
		Delay:      2 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed to get realization information for VNA %s in cluster %s: %v", vnaID, clusterID, err)
	}
	return nil
}
