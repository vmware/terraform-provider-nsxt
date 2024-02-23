/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func dataSourceNsxtUpgradePrepareReady() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtUpgradePrepareReadyRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"upgrade_prepare_id": {
				Type:        schema.TypeString,
				Description: "ID of corresponding nsxt_upgrade_prepare resource",
				Required:    true,
			},
		},
	}
}

func dataSourceNsxtUpgradePrepareReadyRead(d *schema.ResourceData, m interface{}) error {
	precheckErrors, err := getPrecheckErrors(m, nil)
	if err != nil {
		return fmt.Errorf("Error while reading precheck failures: %v", err)
	}
	var precheckFailureIDs []string
	var unacknowledgedWarningIDs []string
	for _, precheckError := range precheckErrors {
		errType := *precheckError.Type_
		if errType == nsxModel.UpgradeCheckFailure_TYPE_FAILURE {
			precheckFailureIDs = append(precheckFailureIDs, *precheckError.Id)
		} else {
			if !(*precheckError.Acked) {
				unacknowledgedWarningIDs = append(unacknowledgedWarningIDs, *precheckError.Id)
			}
		}
	}
	var errMessage string
	if len(precheckFailureIDs) > 0 {
		errMessage += fmt.Sprintf("There are failures in precheck: %s, please check their status from nsxt_upgrade_prepare resource and address these failures on NSX", precheckFailureIDs)
	}
	if len(unacknowledgedWarningIDs) > 0 {
		errMessage += fmt.Sprintf("\nThere are unacknowledged warnings in precheck: %s, please address these errors from NSX or using nsxt_upgrade_precheck_acknowledge resource", unacknowledgedWarningIDs)
	}
	if len(errMessage) > 0 {
		return fmt.Errorf(errMessage)
	}
	objID := d.Get("id").(string)
	if objID == "" {
		objID = newUUID()
	}
	d.SetId(objID)

	return nil
}
