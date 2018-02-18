/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtLogicalPort_basic(t *testing.T) {

	portName := fmt.Sprintf("test-nsx-logical-port")
	updatePortName := fmt.Sprintf("%s-update", portName)
	testResourceName := "nsxt_logical_port.test"
	transportZoneName := getOverlayTransportZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalPortCheckDestroy(state, portName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalPortCreateTemplate(portName, transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalPortExists(portName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", portName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXLogicalPortUpdateTemplate(updatePortName, transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalPortExists(updatePortName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePortName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "admin_state", "UP"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
				),
			},
		},
	})
}

func testAccNSXLogicalPortExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical port resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical port resource ID not set in resources ")
		}

		logicalPort, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical port ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical port %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == logicalPort.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical port %s wasn't found", display_name)
	}
}

func testAccNSXLogicalPortCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_port" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		logicalPort, responseCode, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving logical port ID %s. Error: %v", resourceID, err)
		}

		if display_name == logicalPort.DisplayName {
			return fmt.Errorf("NSX logical port %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXLogicalPortCreateTemplate(portName string, transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
     display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
	display_name = "test_switch"
	admin_state = "UP"
	replication_mode = "MTEP"
	transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}

resource "nsxt_logical_port" "test" {
	display_name = "%s"
	admin_state = "UP"
	description = "Acceptance Test"
	logical_switch_id = "${nsxt_logical_switch.test.id}"
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
}`, transportZoneName, portName)
}

func testAccNSXLogicalPortUpdateTemplate(portUpdatedName string, transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
     display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
	display_name = "test_switch"
	admin_state = "UP"
	replication_mode = "MTEP"
	transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}

resource "nsxt_logical_port" "test" {
	display_name = "%s"
	admin_state = "UP"
	description = "Acceptance Test Update"
	logical_switch_id = "${nsxt_logical_switch.test.id}"
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
    tag {
    	scope = "scope2"
        tag = "tag2"
    }
}`, transportZoneName, portUpdatedName)
}
