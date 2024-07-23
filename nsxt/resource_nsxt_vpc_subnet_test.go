/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestVpcSubnetCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"ip_addresses":     "192.168.240.0/26",
	"access_mode":      "Isolated",
	"enabled":          "true",
	"enable_dhcp":      "true",
	"ipv4_subnet_size": "16",
}

var accTestVpcSubnetUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"ip_addresses":     "192.168.240.0/26",
	"access_mode":      "Isolated",
	"enabled":          "true",
	"enable_dhcp":      "true",
	"ipv4_subnet_size": "16",
}

func TestAccResourceNsxtVpcSubnet_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetCreateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.enable_dhcp", accTestVpcSubnetCreateAttributes["enable_dhcp"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetUpdateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.enable_dhcp", accTestVpcSubnetUpdateAttributes["enable_dhcp"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			//{
			//	Config: testAccNsxtVpcSubnetMinimalistic(),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
			//		resource.TestCheckResourceAttr(testResourceName, "description", ""),
			//		resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
			//		resource.TestCheckResourceAttrSet(testResourceName, "path"),
			//		resource.TestCheckResourceAttrSet(testResourceName, "revision"),
			//		resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
			//	),
			//},
		},
	})
}

func TestAccResourceNsxtVpcSubnet_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtVpcSubnetExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VpcSubnet resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VpcSubnet resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcSubnetExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VpcSubnet %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcSubnetCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_subnet" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcSubnetExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("VpcSubnet %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcSubnetTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetCreateAttributes
	} else {
		attrMap = accTestVpcSubnetUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  description  = "%s"

  ip_addresses = ["%s"]
  access_mode = "%s"
  ipv4_subnet_size = %s
  dhcp_config {
    enable_dhcp = %s
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["ip_addresses"], attrMap["access_mode"], attrMap["ipv4_subnet_size"], attrMap["enable_dhcp"])
}

func testAccNsxtVpcSubnetMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcSubnetUpdateAttributes["display_name"])
}
