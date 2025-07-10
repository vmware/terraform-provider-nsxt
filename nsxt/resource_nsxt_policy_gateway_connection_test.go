// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestPolicyGatewayConnectionCreateAttributes = map[string]string{
	"display_name":            getAccTestResourceName(),
	"description":             "terraform created",
	"aggregate_routes":        "192.168.240.0/24",
	"enable_snat":             "false",
	"logging_enabled":         "false",
	"inbound_remote_networks": "196.1.1.0/24",
}

var accTestPolicyGatewayConnectionUpdateAttributes = map[string]string{
	"display_name":            getAccTestResourceName(),
	"description":             "terraform updated",
	"aggregate_routes":        "192.168.241.0/24",
	"enable_snat":             "true",
	"logging_enabled":         "true",
	"inbound_remote_networks": "196.1.3.0/24",
}

func TestAccResourceNsxtPolicyGatewayConnection_basic(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_connection.test"
	testDataSourceName := "data.nsxt_policy_gateway_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayConnectionCheckDestroy(state, accTestPolicyGatewayConnectionUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayConnectionTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayConnectionExists(accTestPolicyGatewayConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyGatewayConnectionCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyGatewayConnectionCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestPolicyGatewayConnectionCreateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayConnectionTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayConnectionExists(accTestPolicyGatewayConnectionUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyGatewayConnectionUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyGatewayConnectionUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestPolicyGatewayConnectionUpdateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayConnectionMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayConnectionExists(accTestPolicyGatewayConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewayConnection_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_gateway_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayConnectionCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayConnectionMinimalistic(),
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

func TestAccResourceNsxtPolicyGatewayConnection_910(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")

		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayConnectionCheckDestroy(state, accTestPolicyGatewayConnectionUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayConnection910Template(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayConnectionExists(accTestPolicyGatewayConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyGatewayConnectionCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyGatewayConnectionCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestPolicyGatewayConnectionCreateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),
					resource.TestCheckResourceAttr(testResourceName, "inbound_remote_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "inbound_remote_networks.0", accTestPolicyGatewayConnectionCreateAttributes["inbound_remote_networks"]),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.0.enable_snat", accTestPolicyGatewayConnectionCreateAttributes["enable_snat"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nat_config.0.ip_block"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.0.logging_enabled", accTestPolicyGatewayConnectionCreateAttributes["logging_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayConnection910Template(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayConnectionExists(accTestPolicyGatewayConnectionUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyGatewayConnectionUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyGatewayConnectionUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestPolicyGatewayConnectionUpdateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),
					resource.TestCheckResourceAttr(testResourceName, "inbound_remote_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "inbound_remote_networks.0", accTestPolicyGatewayConnectionUpdateAttributes["inbound_remote_networks"]),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.0.enable_snat", accTestPolicyGatewayConnectionUpdateAttributes["enable_snat"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nat_config.0.ip_block"),
					resource.TestCheckResourceAttr(testResourceName, "nat_config.0.logging_enabled", accTestPolicyGatewayConnectionUpdateAttributes["logging_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func testAccNsxtPolicyGatewayConnectionExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy GatewayConnection resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy GatewayConnection resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyGatewayConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy GatewayConnection %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyGatewayConnectionCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_gateway_connection" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyGatewayConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy GatewayConnection %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyGatewayConnectionPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "terraformt0gw"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}`, getEdgeClusterName())
}

func testAccNsxtPolicyGatewayConnectionTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyGatewayConnectionCreateAttributes
	} else {
		attrMap = accTestPolicyGatewayConnectionUpdateAttributes
	}
	return testAccNsxtPolicyGatewayConnectionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_gateway_connection" "test" {
  display_name     = "%s"
  description      = "%s"
  tier0_path       = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_gateway_connection" "test" {
  tier0_path   = nsxt_policy_tier0_gateway.test.path
  display_name = "%s"
  depends_on   = [nsxt_policy_gateway_connection.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["aggregate_routes"], attrMap["display_name"])
}

func testAccNsxtPolicyGatewayConnection910Template(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyGatewayConnectionCreateAttributes
	} else {
		attrMap = accTestPolicyGatewayConnectionUpdateAttributes
	}
	return testAccNsxtPolicyGatewayConnectionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_ip_block" "defExtBlock" {
  display_name = "defExtBlock"
  visibility   = "EXTERNAL"
  cidr_list  = ["22.21.0.0/16"]
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name = "%s"
  description  = "%s"
  tier0_path   = nsxt_policy_tier0_gateway.test.path

  aggregate_routes        = ["%s"]
  inbound_remote_networks = ["%s"]

  advertise_outbound_networks {
    allow_external_blocks = [nsxt_policy_ip_block.defExtBlock.path]
  }

  nat_config {
    enable_snat     = %s
    ip_block        = nsxt_policy_ip_block.defExtBlock.path
    logging_enabled = %s
  }
}
`, attrMap["display_name"], attrMap["description"], attrMap["aggregate_routes"], attrMap["inbound_remote_networks"], attrMap["enable_snat"], attrMap["logging_enabled"])
}

func testAccNsxtPolicyGatewayConnectionMinimalistic() string {
	return testAccNsxtPolicyGatewayConnectionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_gateway_connection" "test" {
  tier0_path   = nsxt_policy_tier0_gateway.test.path
  display_name = "%s"
}

data "nsxt_policy_gateway_connection" "test" {
  display_name = "%s"
  depends_on   = [nsxt_policy_gateway_connection.test]
}`, accTestPolicyGatewayConnectionUpdateAttributes["display_name"], accTestPolicyGatewayConnectionUpdateAttributes["display_name"])
}
