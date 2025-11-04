// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
)

func resourceNsxtUpgradePrecheckAcknowledge() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtUpgradePrecheckAcknowledgeCreate,
		Read:   resourceNsxtUpgradePrecheckAcknowledgeRead,
		Update: resourceNsxtUpgradePrecheckAcknowledgeUpdate,
		Delete: resourceNsxtUpgradePrecheckAcknowledgeDelete,

		Schema: map[string]*schema.Schema{
			"precheck_ids": {
				Type:        schema.TypeList,
				Description: "IDs of precheck warnings that need to be acknowledged",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"precheck_warnings": {
				Type:        schema.TypeList,
				Description: "List of warnings from precheck",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the precheck warning",
							Computed:    true,
						},
						"message": {
							Type:        schema.TypeString,
							Description: "Message of the precheck warning",
							Computed:    true,
						},
						"is_acknowledged": {
							Type:        schema.TypeBool,
							Description: "Boolean value which identifies if the precheck warning has been acknowledged.",
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"target_version": {
				Type:        schema.TypeString,
				Description: "Target system version",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func validatePrecheckIDs(m interface{}, precheckIDs []string) error {
	var validChecks []string
	connector := getPolicyConnector(m)
	client := upgrade.NewUpgradeChecksInfoClient(connector)
	checkInfoResults, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	for _, checkInfo := range checkInfoResults.Results {
		for _, ci := range checkInfo.PreUpgradeChecksInfo {
			validChecks = append(validChecks, *ci.Id)
		}
	}

	for _, precheckID := range precheckIDs {
		// Here we are checking if the precheck ID is in the list of valid checks.
		// If not we are also checking if the ID has a "-" and the string before it is in the list.
		// This is done because, in some cases NSX adds the precheck in <precheckID>-<taskName> format and only the precheck ID is considered while acknowledging/resolving
		if !slices.Contains(validChecks, precheckID) && !slices.Contains(validChecks, strings.Split(precheckID, "-")[0]) {
			return fmt.Errorf("precheck ID %s is not valid", precheckID)
		}
	}
	return nil
}

func resourceNsxtUpgradePrecheckAcknowledgeCreate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		id = newUUID()
	}
	d.SetId(id)
	precheckIDs := interface2StringList(d.Get("precheck_ids").([]interface{}))
	err := validatePrecheckIDs(m, precheckIDs)
	if err != nil {
		return handleCreateError("NsxtPrecheckAcknowledge", id, err)
	}
	err = acknowledgePrecheckWarnings(m, precheckIDs)
	if err != nil {
		return handleCreateError("NsxtPrecheckAcknowledge", id, err)
	}
	return resourceNsxtUpgradePrecheckAcknowledgeRead(d, m)
}

func acknowledgePrecheckWarnings(m interface{}, precheckIDs []string) error {
	connector := getPolicyConnector(m)

	typeParam := nsxModel.UpgradeCheckFailure_TYPE_WARNING
	precheckWarnings, err := getPrecheckErrors(m, &typeParam)
	if err != nil {
		return err
	}
	var warns []string
	for _, warn := range precheckWarnings {
		warns = append(warns, *warn.Id)
	}

	client := upgrade.NewPreUpgradeChecksClient(connector)
	for _, precheckID := range precheckIDs {
		if slices.Contains(warns, precheckID) {
			err := client.Acknowledge(precheckID)
			if err != nil {
				msg := fmt.Sprintf("Failed to acknowledge precheck warning with ID %s", precheckID)
				return logAPIError(msg, err)
			}
		} else {
			log.Printf("[INFO] Precheck warning with ID %s has not been triggered by NSX", precheckID)
		}
	}
	return nil
}

func resourceNsxtUpgradePrecheckAcknowledgeRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	typeParam := nsxModel.UpgradeCheckFailure_TYPE_WARNING
	precheckWarnings, err := getPrecheckErrors(m, &typeParam)
	if err != nil {
		return handleReadError(d, "NsxtUpgradePrecheckAcknowledge", id, err)
	}
	err = setPrecheckWarningsInSchema(d, precheckWarnings)
	if err != nil {
		return handleReadError(d, "NsxtUpgradePrecheckAcknowledge", id, err)
	}
	return nil
}

func setPrecheckWarningsInSchema(d *schema.ResourceData, precheckWarnings []nsxModel.UpgradeCheckFailure) error {
	var precheckWarningList []map[string]interface{}
	for _, precheckWarning := range precheckWarnings {
		elem := make(map[string]interface{})
		elem["id"] = precheckWarning.Id
		elem["message"] = precheckWarning.Message.Message
		elem["is_acknowledged"] = precheckWarning.Acked
		precheckWarningList = append(precheckWarningList, elem)
	}
	return d.Set("precheck_warnings", precheckWarningList)
}

func resourceNsxtUpgradePrecheckAcknowledgeUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	precheckIDs := interface2StringList(d.Get("precheck_ids").([]interface{}))
	if d.HasChange("precheck_ids") {
		err := validatePrecheckIDs(m, precheckIDs)
		if err != nil {
			return handleCreateError("NsxtPrecheckAcknowledge", id, err)
		}
	}
	err := acknowledgePrecheckWarnings(m, precheckIDs)
	if err != nil {
		return handleUpdateError("NsxtPrecheckAcknowledge", id, err)
	}
	return resourceNsxtUpgradePrecheckAcknowledgeRead(d, m)
}

func resourceNsxtUpgradePrecheckAcknowledgeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getAcknowledgedPrecheckIDs(m interface{}) ([]string, error) {
	var result []string
	typeParam := nsxModel.UpgradeCheckFailure_TYPE_WARNING
	precheckWarnings, err := getPrecheckErrors(m, &typeParam)
	if err != nil {
		return result, err
	}
	for _, warning := range precheckWarnings {
		if warning.Acked != nil && *warning.Acked {
			result = append(result, *warning.Id)
		}
	}
	return result, nil
}
