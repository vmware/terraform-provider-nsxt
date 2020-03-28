/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func TestAccDataSourceNsxtLogicalSwitch_basic(t *testing.T) {
	switchName := "terraform_test_switch"
	testResourceName := "data.nsxt_logical_switch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtLogicalSwitchDeleteByName(switchName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtLogicalSwitchCreate(switchName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXLogicalSwitchReadTemplate(switchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", switchName),
				),
			},
			{
				Config: testAccNSXNoLogicalSwitchTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtLogicalSwitchCreate(switchName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := switchName
	description := switchName
	adminState := "UP"
	replicationMode := "MTEP"
	logicalSwitch := manager.LogicalSwitch{
		AdminState:      adminState,
		Description:     description,
		DisplayName:     displayName,
		TransportZoneId: getTransportZoneID(),
		ReplicationMode: replicationMode,
	}

	logicalSwitch, resp, err := nsxClient.LogicalSwitchingApi.CreateLogicalSwitch(nsxClient.Context, logicalSwitch)
	if err != nil {
		return fmt.Errorf("Error during logical switch creation: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during logical switch creation: %v", resp.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtLogicalSwitchDeleteByName(switchName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name
	objList, _, err := nsxClient.LogicalSwitchingApi.ListLogicalSwitches(nsxClient.Context, nil)
	if err != nil {
		return fmt.Errorf("Error while reading logical switches: %v", err)
	}
	// go over the list to find the correct one
	for _, objInList := range objList.Results {
		if objInList.DisplayName == switchName {
			localVarOptionals := make(map[string]interface{})
			responseCode, err := nsxClient.LogicalSwitchingApi.DeleteLogicalSwitch(nsxClient.Context, objInList.Id, localVarOptionals)
			if err != nil {
				return fmt.Errorf("Error during logical switch deletion: %v", err)
			}

			if responseCode.StatusCode != http.StatusOK {
				return fmt.Errorf("Unexpected status returned during logical switch deletion: %v", responseCode.StatusCode)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting logical switch '%s': service not found", switchName)
}

func testAccNSXLogicalSwitchReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_logical_switch" "test" {
  display_name = "%s"
}`, name)
}

func testAccNSXNoLogicalSwitchTemplate() string {
	return fmt.Sprintf(` `)
}
