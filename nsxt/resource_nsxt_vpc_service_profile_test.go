/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyVpcServiceProfileCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"ntp_servers":      "5.5.5.5",
	"lease_time":       "50840",
	"mode":             "DHCP_IP_ALLOCATION_BY_PORT",
	"dns_server_ips":   "7.7.7.7",
	"server_addresses": "11.11.11.11",
	"cache_size":       "2048",
	"log_level":        "DEBUG",
}

var accTestPolicyVpcServiceProfileUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"ntp_servers":      "5.5.5.7",
	"lease_time":       "148000",
	"mode":             "DHCP_IP_ALLOCATION_BY_MAC",
	"dns_server_ips":   "7.7.7.2",
	"server_addresses": "11.11.11.111",
	"cache_size":       "1024",
	"log_level":        "INFO",
}

func TestAccResourceNsxtVpcServiceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, accTestPolicyVpcServiceProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(accTestPolicyVpcServiceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.ntp_servers.0", accTestPolicyVpcServiceProfileCreateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.lease_time", accTestPolicyVpcServiceProfileCreateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", accTestPolicyVpcServiceProfileCreateAttributes["mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfileCreateAttributes["dns_server_ips"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.cache_size", accTestPolicyVpcServiceProfileCreateAttributes["cache_size"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.log_level", accTestPolicyVpcServiceProfileCreateAttributes["log_level"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(accTestPolicyVpcServiceProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.ntp_servers.0", accTestPolicyVpcServiceProfileUpdateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.lease_time", accTestPolicyVpcServiceProfileUpdateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", accTestPolicyVpcServiceProfileUpdateAttributes["mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfileUpdateAttributes["dns_server_ips"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.cache_size", accTestPolicyVpcServiceProfileUpdateAttributes["cache_size"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.log_level", accTestPolicyVpcServiceProfileUpdateAttributes["log_level"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(accTestPolicyVpcServiceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcServiceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtVpcServiceProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcServiceProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcServiceProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcServiceProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcServiceProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcServiceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_service_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcServiceProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcServiceProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcServiceProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyVpcServiceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyVpcServiceProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"
  description  = "%s"

  dhcp_config {
    ntp_servers = ["%s"]

    lease_time = %s
    mode = "%s"

    dns_client_config {
      dns_server_ips = ["%s"]
    }
  }

  dns_forwarder_config {
    cache_size = %s
    log_level = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["ntp_servers"], attrMap["lease_time"], attrMap["mode"], attrMap["dns_server_ips"], attrMap["cache_size"], attrMap["log_level"])
}

func testAccNsxtVpcServiceProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"
  
  dhcp_config {
    mode = "%s"
  }
}`, testAccNsxtPolicyMultitenancyContext(), accTestPolicyVpcServiceProfileUpdateAttributes["display_name"], accTestPolicyVpcServiceProfileUpdateAttributes["mode"])
}
