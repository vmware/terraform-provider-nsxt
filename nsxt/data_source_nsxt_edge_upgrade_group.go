/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
)

var (
	edgeUpgradeGroup = "EDGE"
	hostUpgradeGroup = "HOST"

	timeoutUUGCreat  = 1 * time.Minute
	intervalUUGCreat = 5 * time.Second
)

func dataSourceNsxtEdgeUpgradeGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtEdgeUpgradeGroupRead,

		Schema: map[string]*schema.Schema{
			"upgrade_prepare_id": {
				Type:        schema.TypeString,
				Description: "ID of corresponding nsxt_upgrade_prepare resource",
				Required:    true,
			},
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
		},
	}
}

func dataSourceNsxtEdgeUpgradeGroupRead(d *schema.ResourceData, m interface{}) error {
	return upgradeGroupRead(d, m, edgeUpgradeGroup)
}

func upgradeGroupRead(d *schema.ResourceData, m interface{}, groupType string) error {
	connector := getPolicyConnector(m)
	client := upgrade.NewUpgradeUnitGroupsClient(connector)

	err := waitUpgradeGroupCreate(client, groupType)
	if err != nil {
		return fmt.Errorf("unable to fetch UpgradeUnitGroup: %v", err)
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.UpgradeUnitGroup
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID, nil)
		if isNotFoundError(err) {
			return fmt.Errorf("%s UpgradeUnitGroup with ID %s was not found", groupType, objID)
		}

		if err != nil {
			return fmt.Errorf("error while reading %s UpgradeUnitGroup %s: %v", groupType, objID, err)
		}

		if *objGet.Type_ != groupType {
			return fmt.Errorf("%s UpgradeUnitGroup with ID %s was not found", groupType, objID)
		}
		obj = objGet

	} else if objName == "" {
		return fmt.Errorf("error obtaining %s UpgradeUnitGroup ID or name during read", groupType)
	} else {
		// Get by full name/prefix
		objList, err := client.List(&groupType, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error while reading %s UpgradeUnitGroup: %v", groupType, err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.UpgradeUnitGroup
		var prefixMatch []model.UpgradeUnitGroup
		for _, objInList := range objList.Results {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("found multiple %s UpgradeUnitGroup with name '%s'", groupType, objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("found multiple %s UpgradeUnitGroup with name starting with '%s'", groupType, objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("%s UpgradeUnitGroup with name '%s' was not found", groupType, objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)

	return nil
}

func waitUpgradeGroupCreate(client upgrade.UpgradeUnitGroupsClient, groupType string) error {

	resultChan := make(chan error, 1)

	poll := func(resultChan chan error) {
		preLen := 0
		for {
			objList, err := client.List(&groupType, nil, nil, nil, nil, nil, nil, nil)
			if err != nil && !isNotFoundError(err) {
				resultChan <- err
				return
			}
			if err == nil && len(objList.Results) > 0 {
				if preLen == 0 {
					preLen = len(objList.Results)
				} else if preLen == len(objList.Results) {
					resultChan <- nil
					return
				}
			}
			time.Sleep(intervalUUGCreat)
		}
	}

	go poll(resultChan)

	select {
	case res := <-resultChan:
		return res
	case <-time.After(timeoutUUGCreat):
		return fmt.Errorf("timeout reached")
	}
}
