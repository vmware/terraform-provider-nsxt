// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var accTestPolicyRCBgpNeighborCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"neighbor_address": "192.168.100.1",
	"remote_as_num":    "65001",
	"source_address":   "192.168.200.1",
}

var accTestPolicyRCBgpNeighborUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"neighbor_address": "192.168.100.2",
	"remote_as_num":    "65002",
	"source_address":   "192.168.200.1",
}

func testAccNsxtPolicyRCBgpNeighborPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_NAME")
	testAccEnvDefined(t, "NSXT_TEST_PORTGROUP_ID")
}

func TestAccResourceNsxtPolicyRouteControllerBgpNeighbor_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_bgp_neighbor.test"
	testDataSourceName := "data.nsxt_policy_route_controller_bgp_neighbor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRCBgpNeighborCheckDestroy(state, accTestPolicyRCBgpNeighborUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCBgpNeighborTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRCBgpNeighborExists(accTestPolicyRCBgpNeighborCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRCBgpNeighborCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRCBgpNeighborCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "neighbor_address", accTestPolicyRCBgpNeighborCreateAttributes["neighbor_address"]),
					resource.TestCheckResourceAttr(testResourceName, "remote_as_num", accTestPolicyRCBgpNeighborCreateAttributes["remote_as_num"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "hold_down_time", "180"),
					resource.TestCheckResourceAttr(testResourceName, "keep_alive_time", "60"),
					resource.TestCheckResourceAttr(testResourceName, "maximum_hop_limit", "1"),
					resource.TestCheckResourceAttr(testResourceName, "source_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "source_addresses.0", accTestPolicyRCBgpNeighborCreateAttributes["source_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyRCBgpNeighborTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRCBgpNeighborExists(accTestPolicyRCBgpNeighborUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRCBgpNeighborUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRCBgpNeighborUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "neighbor_address", accTestPolicyRCBgpNeighborUpdateAttributes["neighbor_address"]),
					resource.TestCheckResourceAttr(testResourceName, "remote_as_num", accTestPolicyRCBgpNeighborUpdateAttributes["remote_as_num"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "source_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "source_addresses.0", accTestPolicyRCBgpNeighborUpdateAttributes["source_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyRouteControllerBgpNeighbor_withRouteFiltering(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_bgp_neighbor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRCBgpNeighborCheckDestroy(state, accTestPolicyRCBgpNeighborCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCBgpNeighborWithRouteFilteringTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRCBgpNeighborExists(accTestPolicyRCBgpNeighborCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "neighbor_address", accTestPolicyRCBgpNeighborCreateAttributes["neighbor_address"]),
					resource.TestCheckResourceAttr(testResourceName, "remote_as_num", accTestPolicyRCBgpNeighborCreateAttributes["remote_as_num"]),
					resource.TestCheckResourceAttr(testResourceName, "route_filtering.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_filtering.0.address_family", "IPV4"),
					resource.TestCheckResourceAttr(testResourceName, "route_filtering.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyRouteControllerBgpNeighbor_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_bgp_neighbor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRCBgpNeighborCheckDestroy(state, accTestPolicyRCBgpNeighborCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCBgpNeighborMinimalistic(),
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

func testAccNsxtPolicyRCBgpNeighborExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		parents, err := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
		if err != nil {
			return err
		}

		sessionContext := utl.SessionContext{ClientType: utl.Local}
		c := infra.NewRouteControllerBgpNeighborClient(sessionContext, connector)
		if c == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err = c.Get(parents[0], resourceID)
		if err != nil {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyRCBgpNeighborCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller_bgp_neighbor" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		parents, err := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerBgpPathExample)
		if err != nil {
			return err
		}

		sessionContext := utl.SessionContext{ClientType: utl.Local}
		c := infra.NewRouteControllerBgpNeighborClient(sessionContext, connector)
		if c == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err = c.Get(parents[0], resourceID)
		if err == nil {
			return fmt.Errorf("Policy RouteControllerBgpNeighbor %s still exists", displayName)
		}
	}
	return nil
}

// testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate generates HCL for a route controller
// with BGP enabled, which is the prerequisite for BGP neighbors.
func testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
data "nsxt_policy_virtual_network_appliance" "vna" {
  display_name = "%s"
  cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller" "rc" {
  display_name                           = "tf-acc-bgp-nbr-rc"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path

  bgp_config {
    ecmp                   = true
    graceful_restart_mode  = "HELPER_ONLY"
    graceful_restart_timer = 180
    local_as_num           = "%s"
  }
}
`, getRCVNAName(), "65000")
}

func testAccNsxtPolicyRCBgpNeighborInterfaceTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_route_controller_interface" "bgp_src" {
  display_name        = "tf-acc-bgp-nbr-intf"
  parent_path         = nsxt_policy_route_controller.rc.path
  floating_ip_subnets = ["192.168.200.1/24"]

  interface_address {
    subnets                        = ["192.168.201.1/24"]
    portgroup_id                   = "%s"
    virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance.vna.path
  }
}
`, getPortgroupID())
}

func testAccNsxtPolicyRCBgpNeighborTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyRCBgpNeighborCreateAttributes
	} else {
		attrMap = accTestPolicyRCBgpNeighborUpdateAttributes
	}
	enabled := "true"
	if !createFlow {
		enabled = "false"
	}
	return testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate() + testAccNsxtPolicyRCBgpNeighborInterfaceTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name     = "%s"
  description      = "%s"
  parent_path      = "${nsxt_policy_route_controller.rc.path}/bgp"
  neighbor_address = "%s"
  remote_as_num    = "%s"
  enabled          = %s
  source_addresses = ["192.168.200.1"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  depends_on = [nsxt_policy_route_controller_interface.bgp_src]
}

data "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name = "%s"
  parent_path  = "${nsxt_policy_route_controller.rc.path}/bgp"

  depends_on = [nsxt_policy_route_controller_bgp_neighbor.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["neighbor_address"], attrMap["remote_as_num"], enabled, attrMap["display_name"])
}

func testAccNsxtPolicyRCBgpNeighborWithRouteFilteringTemplate() string {
	attrMap := accTestPolicyRCBgpNeighborCreateAttributes
	return testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate() + testAccNsxtPolicyRCBgpNeighborInterfaceTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name     = "%s"
  description      = "%s"
  parent_path      = "${nsxt_policy_route_controller.rc.path}/bgp"
  neighbor_address = "%s"
  remote_as_num    = "%s"
  source_addresses = ["192.168.200.1"]

  route_filtering {
    address_family = "IPV4"
    enabled        = true
  }

  depends_on = [nsxt_policy_route_controller_interface.bgp_src]
}`, attrMap["display_name"], attrMap["description"], attrMap["neighbor_address"], attrMap["remote_as_num"])
}

func testAccNsxtPolicyRCBgpNeighborMinimalistic() string {
	attrMap := accTestPolicyRCBgpNeighborCreateAttributes
	return testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate() + testAccNsxtPolicyRCBgpNeighborInterfaceTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name     = "%s"
  parent_path      = "${nsxt_policy_route_controller.rc.path}/bgp"
  neighbor_address = "%s"
  remote_as_num    = "%s"
  source_addresses = ["192.168.200.1"]

  depends_on = [nsxt_policy_route_controller_interface.bgp_src]
}`, attrMap["display_name"], attrMap["neighbor_address"], attrMap["remote_as_num"])
}
