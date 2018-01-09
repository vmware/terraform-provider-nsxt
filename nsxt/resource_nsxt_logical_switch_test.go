/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"regexp"
	"testing"
)

func TestNSXLogicalSwitchBasic(t *testing.T) {
	// Test without verification for realization state

	testNSXLogicalSwitchBasic(t, false)
}

func TestNSXLogicalSwitchBasicWithRealization(t *testing.T) {
	// Test with verification for realization state

	testNSXLogicalSwitchBasic(t, true)
}

func TestNSXLogicalSwitchVlan(t *testing.T) {
	// Test without verification for realization state

	testNSXLogicalSwitchVlan(t, false)
}

func TestNSXLogicalSwitchVlanWithRealization(t *testing.T) {
	// Test with verification for realization state

	testNSXLogicalSwitchVlan(t, true)
}

func testNSXLogicalSwitchBasic(t *testing.T, verifyRealization bool) {

	switchName := fmt.Sprintf("test-nsx-logical-switch-overlay")
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	testResourceName := "nsxt_logical_switch.test"
	novlan := "0"
	replicationMode := "MTEP"
	transportZoneName := OverlayTransportZoneNamePrefix

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoTZIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccNSXLogicalSwitchCreateTemplate(switchName, transportZoneName, novlan, replicationMode, verifyRealization),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "replication_mode", replicationMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", novlan),
				),
			},
			{
				Config: testAccNSXLogicalSwitchUpdateTemplate(updateSwitchName, transportZoneName, novlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "replication_mode", replicationMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", novlan),
				),
			},
		},
	})
}

func testNSXLogicalSwitchVlan(t *testing.T, verifyRealization bool) {

	switchName := fmt.Sprintf("test-nsx-logical-switch-vlan")
	updateSwitchName := fmt.Sprintf("%s-update", switchName)
	transportZoneName := VlanTransportZoneName
	testResourceName := "nsxt_logical_switch.test"

	origvlan := "1"
	updatedvlan := "2"
	replicationMode := ""

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoVlanTemplate(switchName, transportZoneName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config: testAccNSXLogicalSwitchCreateTemplate(switchName, transportZoneName, origvlan, replicationMode, verifyRealization),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", origvlan),
				),
			},
			{
				Config: testAccNSXLogicalSwitchUpdateTemplate(updateSwitchName, transportZoneName, updatedvlan, replicationMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "vlan", updatedvlan),
				),
			},
		},
	})

}

func testAccNSXLogicalSwitchExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical switch resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical switch resource ID not set in resources ")
		}

		logicalSwitch, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical switch %s exists. HTTP return code was %d", resourceID, responseCode)
		}

		if display_name == logicalSwitch.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical switch %s wasn't found", display_name)
	}
}

func testAccNSXLogicalSwitchCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_switch" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		logicalSwitch, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}

		if display_name == logicalSwitch.DisplayName {
			return fmt.Errorf("NSX logical switch %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXLogicalSwitchNoTZIDTemplate(switchName string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
replication_mode = "MTEP"
}`, switchName)
}

func testAccNSXLogicalSwitchNoVlanTemplate(switchName string, transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
     display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}`, transportZoneName, switchName)
}

func testAccNSXLogicalSwitchCreateTemplate(switchName string, transportZoneName string, vlan string, replicationMode string, verifyRealization bool) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
     display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test"
transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
replication_mode = "%s"
vlan = "%s"
verify_realization = "%t"
tags = [
    {
	scope = "scope1"
        tag = "tag1"
    }
]
}`, transportZoneName, switchName, replicationMode, vlan, verifyRealization)
}

func testAccNSXLogicalSwitchUpdateTemplate(switchUpdateName string, transportZoneName string, vlan string, replicationMode string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
     display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
display_name = "%s"
admin_state = "UP"
description = "Acceptance Test Update"
transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
replication_mode = "%s"
vlan = "%s"
tags = [
	{
		scope = "scope1"
    	tag = "tag1"
    },
	{
		scope = "scope2"
    	tag = "tag2"
    },
]
}`, transportZoneName, switchUpdateName, replicationMode, vlan)
}
