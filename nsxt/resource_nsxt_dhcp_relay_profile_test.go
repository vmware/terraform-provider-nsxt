/* Copyright © 2018 VMware, Inc. All Rights Reserved.
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

func TestAccResourceNsxtDhcpRelayProfile_basic(t *testing.T) {

	prfName := fmt.Sprintf("test-nsx-dhcp-relay-profile")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_dhcp_relay_profile.test"
	server_ip := "1.1.1.1"
	additional_ip := "2.1.1.1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXDhcpRelayProfileCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXDhcpRelayProfileCreateTemplate(prfName, server_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXDhcpRelayProfileExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "server_addresses.#", "1"),
				),
			},
			{
				Config: testAccNSXDhcpRelayProfileUpdateTemplate(updatePrfName, server_ip, additional_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXDhcpRelayProfileExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "server_addresses.#", "2"),
				),
			},
		},
	})
}

func testAccNSXDhcpRelayProfileExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Dhcp Relay Profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Dhcp Relay Profile resource ID not set in resources ")
		}

		profile, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadDhcpRelayProfile(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving dhcp relay profile ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if dhcp relay profile %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == profile.DisplayName {
			return nil
		}
		return fmt.Errorf("Dhcp Relay Profile %s wasn't found", display_name)
	}
}

func testAccNSXDhcpRelayProfileCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_port" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		profile, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadDhcpRelayProfile(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving dhcp relay profile ID %s. Error: %v", resourceID, err)
		}

		if display_name == profile.DisplayName {
			return fmt.Errorf("Dhcp Relay Profile %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXDhcpRelayProfileCreateTemplate(prfName string, server_ip string) string {
	return fmt.Sprintf(`
resource "nsxt_dhcp_relay_profile" "test" {
	display_name = "%s"
	description = "Acceptance Test"
	server_addresses = ["%s"]
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
}`, prfName, server_ip)
}

func testAccNSXDhcpRelayProfileUpdateTemplate(prfUpdatedName string, server_ip string, additional_ip string) string {
	return fmt.Sprintf(`
resource "nsxt_dhcp_relay_profile" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
	server_addresses = ["%s", "%s"]
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
    tag {
    	scope = "scope2"
        tag = "tag2"
    }
}`, prfUpdatedName, server_ip, additional_ip)
}
