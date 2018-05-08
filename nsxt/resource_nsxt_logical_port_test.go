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

func TestAccResourceNsxtLogicalPort_withProfiles(t *testing.T) {
	portName := fmt.Sprintf("test-nsx-logical-port-with-profiles")
	updatePortName := fmt.Sprintf("%s-update", portName)
	testResourceName := "nsxt_logical_port.test"
	transportZoneName := getOverlayTransportZoneName()
	customProfileName := "terraform_test_LP_profile"
	oobProfileName := "nsx-default-qos-switching-profile"
	profileType := "QosSwitchingProfile"
	var s *terraform.State

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalPortCheckDestroy(state, portName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// Create a custom switching profile
					if err := testAccDataSourceNsxtSwitchingProfileCreate(customProfileName, profileType); err != nil {
						panic(err)
					}
				},
				// Create a logical port to use the custom switching profile
				Config: testAccNSXLogicalPortCreateWithProfilesTemplate(portName, transportZoneName, customProfileName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalPortExists(portName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", portName),
					resource.TestCheckResourceAttr(testResourceName, "switching_profile_id.#", "1"),
				),
			},
			{
				// Replace the custom switching profile with OOB one
				Config:             testAccNSXLogicalPortUpdateWithProfilesTemplate(updatePortName, transportZoneName, customProfileName, oobProfileName),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalPortExists(updatePortName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePortName),
					// Counting only custom profiles so count should be 0
					resource.TestCheckResourceAttr(testResourceName, "switching_profile_id.#", "0"),
					copyStatePtr(&s),
				),
			},
			{
				PreConfig: func() {
					// Delete the configured switching profile
					if err := testAccDataSourceNsxtSwitchingProfileDelete(s, "data.nsxt_switching_profile.test1"); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXSwitchingNoProfileTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtLogicalPort_importBasic(t *testing.T) {
	portName := fmt.Sprintf("test-nsx-logical-port")
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
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXLogicalPortExists(displayName string, resourceName string) resource.TestCheckFunc {
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

		if displayName == logicalPort.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical port %s wasn't found", displayName)
	}
}

func testAccNSXLogicalPortCheckDestroy(state *terraform.State, displayName string) error {
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

		if displayName == logicalPort.DisplayName {
			return fmt.Errorf("NSX logical port %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXLogicalSwitchCreateForPort(transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "TZ1" {
  display_name = "%s"
}

resource "nsxt_logical_switch" "test" {
  display_name      = "test-nsx-switch"
  admin_state       = "UP"
  replication_mode  = "MTEP"
  transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}`, transportZoneName)
}

func testAccNSXLogicalPortCreateTemplate(portName string, transportZoneName string) string {
	return testAccNSXLogicalSwitchCreateForPort(transportZoneName) + fmt.Sprintf(`
resource "nsxt_logical_port" "test" {
  display_name      = "%s"
  admin_state       = "UP"
  description       = "Acceptance Test"
  logical_switch_id = "${nsxt_logical_switch.test.id}"

  tag {
  	scope = "scope1"
    tag   = "tag1"
  }
}`, portName)
}

func testAccNSXLogicalPortUpdateTemplate(portUpdatedName string, transportZoneName string) string {
	return testAccNSXLogicalSwitchCreateForPort(transportZoneName) + fmt.Sprintf(`
resource "nsxt_logical_port" "test" {
  display_name      = "%s"
  admin_state       = "UP"
  description       = "Acceptance Test Update"
  logical_switch_id = "${nsxt_logical_switch.test.id}"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, portUpdatedName)
}

func testAccNSXLogicalPortCreateWithProfilesTemplate(portName string, transportZoneName string, profileName string) string {
	return testAccNSXLogicalSwitchCreateForPort(transportZoneName) + fmt.Sprintf(`
data "nsxt_switching_profile" "test1" {
  display_name = "%s"
}

resource "nsxt_logical_port" "test" {
  display_name      = "%s"
  admin_state       = "UP"
  logical_switch_id = "${nsxt_logical_switch.test.id}"

  switching_profile_id {
    key   = "${data.nsxt_switching_profile.test1.resource_type}"
    value = "${data.nsxt_switching_profile.test1.id}"
  }
}`, profileName, portName)
}

func testAccNSXLogicalPortUpdateWithProfilesTemplate(portUpdatedName string, transportZoneName string, profileName1 string, profileName2 string) string {
	return testAccNSXLogicalSwitchCreateForPort(transportZoneName) + fmt.Sprintf(`
data "nsxt_switching_profile" "test1" {
  display_name = "%s"
}

data "nsxt_switching_profile" "test2" {
  display_name = "%s"
}

resource "nsxt_logical_port" "test" {
  display_name      = "%s"
  admin_state       = "UP"
  logical_switch_id = "${nsxt_logical_switch.test.id}"

  switching_profile_id {
    key   = "${data.nsxt_switching_profile.test2.resource_type}"
    value = "${data.nsxt_switching_profile.test2.id}"
  }

}`, profileName1, profileName2, portUpdatedName)
}
