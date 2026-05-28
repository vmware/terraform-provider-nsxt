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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyVirtualNetworkApplianceRealization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"path": {
				Type:         schema.TypeString,
				Description:  "The policy path of the VirtualNetworkAppliance resource",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"state": {
				Type:        schema.TypeString,
				Description: "Current realization state of the appliance",
				Computed:    true,
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

func getVNAApplianceKeysFromPath(path string) (siteID, epID, clusterID, vnaID string) {
	siteID = getResourceIDFromResourcePath(path, "sites")
	epID = getResourceIDFromResourcePath(path, "enforcement-points")
	clusterID = getResourceIDFromResourcePath(path, "virtual-network-appliance-clusters")
	vnaID = getResourceIDFromResourcePath(path, "virtual-network-appliances")
	return
}

func dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	vnaPath := d.Get("path").(string)
	delay := d.Get("delay").(int)
	timeout := d.Get("timeout").(int)

	siteID, epID, clusterID, vnaID := getVNAApplianceKeysFromPath(vnaPath)
	if siteID == "" {
		return fmt.Errorf("error obtaining site ID from path %s", vnaPath)
	}
	if epID == "" {
		return fmt.Errorf("error obtaining enforcement-point ID from path %s", vnaPath)
	}
	if clusterID == "" {
		return fmt.Errorf("error obtaining cluster ID from path %s", vnaPath)
	}
	if vnaID == "" {
		return fmt.Errorf("error obtaining VNA ID from path %s", vnaPath)
	}

	id := d.Get("id").(string)
	if id == "" {
		d.SetId(newUUID())
	}

	vnaStateClient := cliVNAStateClient(sessionContext, connector)
	if vnaStateClient == nil {
		return policyResourceNotSupportedError()
	}

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
			d.Set("state", status)
			return applianceState, status, nil
		},
		Timeout:    time.Duration(timeout) * time.Second,
		MinTimeout: 5 * time.Second,
		Delay:      time.Duration(delay) * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed to get realization information for VNA %s: %v", vnaPath, err)
	}

	return nil
}
