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

var plPrereqResourceName = getAccTestResourceName()

var accTestPolicyTransitGatewayPrefixListCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"network":      "10.0.0.0/8",
	"action":       "PERMIT",
}

var accTestPolicyTransitGatewayPrefixListUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"network":      "172.16.0.0/12",
	"action":       "DENY",
}

func TestAccResourceNsxtPolicyTransitGatewayPrefixList_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_prefix_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayPrefixListCheckDestroy(state, accTestPolicyTransitGatewayPrefixListUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayPrefixListTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayPrefixListExists(accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayPrefixListCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "prefix.0.network", accTestPolicyTransitGatewayPrefixListCreateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix.0.action", accTestPolicyTransitGatewayPrefixListCreateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayPrefixListTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayPrefixListExists(accTestPolicyTransitGatewayPrefixListUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayPrefixListUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayPrefixListUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "prefix.0.network", accTestPolicyTransitGatewayPrefixListUpdateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix.0.action", accTestPolicyTransitGatewayPrefixListUpdateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayPrefixListMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayPrefixListExists(accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"], testResourceName),
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

// TestAccResourceNsxtPolicyTransitGatewayPrefixList_withNsxID verifies that creating the
// resource with an explicit nsx_id inside a project context succeeds. Regression test for
// BZ#3713681: the existence check in getOrGenerateIDWithParent was deriving a Local client
// context instead of Multitenancy, causing "unsupported client type" when nsx_id was set.
func TestAccResourceNsxtPolicyTransitGatewayPrefixList_withNsxID(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_prefix_list.test"
	nsxID := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayPrefixListCheckDestroy(state, accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayPrefixListWithNsxID(nsxID),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayPrefixListExists(accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", nsxID),
					resource.TestCheckResourceAttr(testResourceName, "id", nsxID),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGatewayPrefixList_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_prefix_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayPrefixListCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayPrefixListMinimalistic(),
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

func testAccNsxtPolicyTransitGatewayPrefixListExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayPrefixList resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayPrefixList resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayPrefixListExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayPrefixList %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayPrefixListCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_transit_gateway_prefix_list" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayPrefixListExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy TransitGatewayPrefixList %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayPrefixListPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_project" "test" {
  display_name = "%s"
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name    = "%s"
  transit_subnets = ["100.64.0.0/16"]
  centralized_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}`, getEdgeClusterName(), plPrereqResourceName, plPrereqResourceName)
}

func testAccNsxtPolicyTransitGatewayPrefixListTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTransitGatewayPrefixListCreateAttributes
	} else {
		attrMap = accTestPolicyTransitGatewayPrefixListUpdateAttributes
	}
	return testAccNsxtPolicyTransitGatewayPrefixListPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_prefix_list" "test" {
  display_name = "%s"
  description  = "%s"
  parent_path  = nsxt_policy_transit_gateway.test.path

  prefix {
    action  = "%s"
    network = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["action"], attrMap["network"])
}

func testAccNsxtPolicyTransitGatewayPrefixListMinimalistic() string {
	return testAccNsxtPolicyTransitGatewayPrefixListPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_prefix_list" "test" {
  display_name = "%s"
  description  = ""
  parent_path  = nsxt_policy_transit_gateway.test.path

  prefix {
    action  = "%s"
    network = "%s"
  }
}`, accTestPolicyTransitGatewayPrefixListUpdateAttributes["display_name"],
		accTestPolicyTransitGatewayPrefixListUpdateAttributes["action"],
		accTestPolicyTransitGatewayPrefixListUpdateAttributes["network"])
}

func testAccNsxtPolicyTransitGatewayPrefixListWithNsxID(nsxID string) string {
	return testAccNsxtPolicyTransitGatewayPrefixListPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_prefix_list" "test" {
  display_name = "%s"
  nsx_id       = "%s"
  parent_path  = nsxt_policy_transit_gateway.test.path

  prefix {
    action  = "%s"
    network = "%s"
  }
}`, accTestPolicyTransitGatewayPrefixListCreateAttributes["display_name"],
		nsxID,
		accTestPolicyTransitGatewayPrefixListCreateAttributes["action"],
		accTestPolicyTransitGatewayPrefixListCreateAttributes["network"])
}
