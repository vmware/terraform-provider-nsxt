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

var accTestPolicyTransitGatewayCommunityListCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"list":         `"NO_EXPORT_SUBCONFED", "65001:12"`,
}

var accTestPolicyTransitGatewayCommunityListUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"list":         `"562:00:12"`,
}

func TestAccResourceNsxtPolicyTransitGatewayCommunityList_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_community_list.test"
	prereqName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCommunityListCheckDestroy(state, accTestPolicyTransitGatewayCommunityListUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayCommunityListTemplate(true, prereqName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayCommunityListExists(accTestPolicyTransitGatewayCommunityListCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayCommunityListCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayCommunityListCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "communities.#", "2"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayCommunityListTemplate(false, prereqName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayCommunityListExists(accTestPolicyTransitGatewayCommunityListUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayCommunityListUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayCommunityListUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "communities.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayCommunityListMinimalistic(prereqName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayCommunityListExists(accTestPolicyTransitGatewayCommunityListCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyTransitGatewayCommunityList_importBasic(t *testing.T) {
	prereqName := getAccTestResourceName()
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_community_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCommunityListCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayCommunityListMinimalistic(prereqName),
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

func testAccNsxtPolicyTransitGatewayCommunityListExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayCommunityList resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayCommunityList resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayCommunityListExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayCommunityList %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayCommunityListCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_transit_gateway_community_list" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayCommunityListExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy TransitGatewayCommunityList %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayCommunityListPrerequisites(prereqName string) string {
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
}`, getEdgeClusterName(), prereqName, prereqName)
}

func testAccNsxtPolicyTransitGatewayCommunityListTemplate(createFlow bool, prereqName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTransitGatewayCommunityListCreateAttributes
	} else {
		attrMap = accTestPolicyTransitGatewayCommunityListUpdateAttributes
	}
	return testAccNsxtPolicyTransitGatewayCommunityListPrerequisites(prereqName) + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_community_list" "test" {
  display_name = "%s"
  description  = "%s"
  parent_path  = nsxt_policy_transit_gateway.test.path
  communities  = [%s]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["list"])
}

func testAccNsxtPolicyTransitGatewayCommunityListMinimalistic(prereqName string) string {
	return testAccNsxtPolicyTransitGatewayCommunityListPrerequisites(prereqName) + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_community_list" "test" {
  display_name = "%s"
  description  = ""
  parent_path  = nsxt_policy_transit_gateway.test.path
  communities  = [%s]
}`, accTestPolicyTransitGatewayCommunityListUpdateAttributes["display_name"],
		accTestPolicyTransitGatewayCommunityListUpdateAttributes["list"])
}
