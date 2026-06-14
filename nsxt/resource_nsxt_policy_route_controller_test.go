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

var accTestPolicyRouteControllerCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"ha_mode":      "ACTIVE_STANDBY",
	"local_as_num": "65001",
}

var accTestPolicyRouteControllerUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"ha_mode":      "ACTIVE_STANDBY",
	"local_as_num": "65002",
}

// accTestPolicyRouteControllerWithBgpName is a dedicated name for the
// _withBgp test so it never shares display_name with the _importBasic test
// (which uses updateAttributes["display_name"]) when both run in parallel.
var accTestPolicyRouteControllerWithBgpName = getAccTestResourceName()

func testAccNsxtPolicyRouteControllerPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
}

// testAccNsxtPolicyRouteControllerVnaTemplate returns HCL that looks up a
// pre-created ROUTE_CONTROLLER VNA cluster (via NSXT_TEST_RC_VNA_CLUSTER_NAME).
// Exposed reference:
//   - data.nsxt_policy_virtual_network_appliance_cluster.vna.path
func testAccNsxtPolicyRouteControllerVnaTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_virtual_network_appliance_cluster" "vna" {
  display_name = "%s"
}
`, getRCVNAClusterName())
}

func TestAccResourceNsxtPolicyRouteController_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"
	testDataSourceName := "data.nsxt_policy_route_controller.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerCheckDestroy(state, accTestPolicyRouteControllerUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRouteControllerCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", accTestPolicyRouteControllerCreateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRouteControllerUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", accTestPolicyRouteControllerUpdateAttributes["ha_mode"]),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyRouteController_withBgp(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerCheckDestroy(state, accTestPolicyRouteControllerWithBgpName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerWithBgpTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerWithBgpName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerWithBgpName),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", accTestPolicyRouteControllerCreateAttributes["local_as_num"]),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_timer", "180"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyRouteControllerWithBgpTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRouteControllerExists(accTestPolicyRouteControllerWithBgpName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRouteControllerWithBgpName),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", accTestPolicyRouteControllerUpdateAttributes["local_as_num"]),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "false"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "DISABLE"),
					resource.TestCheckResourceAttrSet(testResourceName, "virtual_network_appliance_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyRouteController_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRouteControllerCheckDestroy(state, accTestPolicyRouteControllerUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerMinimalistic(),
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

func testAccNsxtPolicyRouteControllerExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteController resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy RouteController resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyRouteControllerExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy RouteController %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyRouteControllerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyRouteControllerExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy RouteController %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyRouteControllerTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyRouteControllerCreateAttributes
	} else {
		attrMap = accTestPolicyRouteControllerUpdateAttributes
	}
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "%s"
  ha_mode                                = "%s"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_route_controller" "test" {
  display_name = "%s"

  depends_on = [nsxt_policy_route_controller.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["ha_mode"], attrMap["display_name"])
}

func testAccNsxtPolicyRouteControllerWithBgpTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyRouteControllerCreateAttributes
	} else {
		attrMap = accTestPolicyRouteControllerUpdateAttributes
	}
	ecmp := "true"
	restartMode := "HELPER_ONLY"
	if !createFlow {
		ecmp = "false"
		restartMode = "DISABLE"
	}
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  description                            = "%s"
  ha_mode                                = "%s"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path

  bgp_config {
    local_as_num           = "%s"
    ecmp                   = %s
    graceful_restart_mode  = "%s"
    graceful_restart_timer = 180
  }
}`, accTestPolicyRouteControllerWithBgpName, attrMap["description"], attrMap["ha_mode"], attrMap["local_as_num"], ecmp, restartMode)
}

func testAccNsxtPolicyRouteControllerMinimalistic() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller" "test" {
  display_name = "%s"

  depends_on = [nsxt_policy_route_controller.test]
}`, accTestPolicyRouteControllerUpdateAttributes["display_name"], accTestPolicyRouteControllerUpdateAttributes["display_name"])
}
