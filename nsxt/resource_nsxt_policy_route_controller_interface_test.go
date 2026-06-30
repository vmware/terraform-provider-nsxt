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

var accTestPolicyRCInterfaceCreateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"subnet":             "192.168.200.1/24",
	"floating_ip_subnet": "192.168.200.10/24",
	"urpf_mode":          "STRICT",
}

var accTestPolicyRCInterfaceUpdateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform updated",
	"subnet":             "192.168.201.1/24",
	"floating_ip_subnet": "192.168.201.10/24",
	"urpf_mode":          "NONE",
}

var accTestPolicyRCInterfaceImportAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"subnet":             "192.168.203.1/24",
	"floating_ip_subnet": "192.168.203.10/24",
}

func testAccNsxtPolicyRCInterfacePreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_NAME")
	testAccEnvDefined(t, "NSXT_TEST_PORTGROUP_ID")
}

func TestAccResourceNsxtPolicyRouteControllerInterface_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_interface.test"
	testDataSourceName := "data.nsxt_policy_route_controller_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCInterfacePreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRCInterfaceCheckDestroy(state, accTestPolicyRCInterfaceUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCInterfaceTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRCInterfaceExists(accTestPolicyRCInterfaceCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRCInterfaceCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRCInterfaceCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", accTestPolicyRCInterfaceCreateAttributes["urpf_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "floating_ip_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "interface_address.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "interface_address.0.subnets.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "interface_address.0.portgroup_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyRCInterfaceTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRCInterfaceExists(accTestPolicyRCInterfaceUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRCInterfaceUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRCInterfaceUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", accTestPolicyRCInterfaceUpdateAttributes["urpf_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "floating_ip_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "interface_address.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "interface_address.0.subnets.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "interface_address.0.portgroup_id"),
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

func TestAccResourceNsxtPolicyRouteControllerInterface_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCInterfacePreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRCInterfaceCheckDestroy(state, accTestPolicyRCInterfaceImportAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCInterfaceMinimalistic(),
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

func testAccNsxtPolicyRCInterfaceExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy RouteControllerInterface resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy RouteControllerInterface resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		parents, err := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
		if err != nil {
			return err
		}

		sessionContext := utl.SessionContext{ClientType: utl.Local}
		c := infra.NewRouteControllerInterfaceClient(sessionContext, connector)
		if c == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err = c.Get(parents[0], resourceID)
		if err != nil {
			return fmt.Errorf("Policy RouteControllerInterface %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyRCInterfaceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_route_controller_interface" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		parents, err := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample)
		if err != nil {
			return err
		}

		sessionContext := utl.SessionContext{ClientType: utl.Local}
		c := infra.NewRouteControllerInterfaceClient(sessionContext, connector)
		if c == nil {
			return fmt.Errorf("unsupported client type")
		}

		_, err = c.Get(parents[0], resourceID)
		if err == nil {
			return fmt.Errorf("Policy RouteControllerInterface %s still exists", displayName)
		}
	}
	return nil
}

// testAccNsxtPolicyRCInterfaceRouteControllerTemplate builds the prerequisite
// infrastructure for route controller interface tests:
//  1. Pre-created ROUTE_CONTROLLER VNA cluster (via NSXT_TEST_RC_VNA_CLUSTER_NAME)
//  2. Pre-created VNA appliance inside that cluster (via NSXT_TEST_RC_VNA_NAME)
//  3. An inline route controller bound to the cluster.
func testAccNsxtPolicyRCInterfaceRouteControllerTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
data "nsxt_policy_virtual_network_appliance" "vna" {
  display_name = "%s"
  cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

resource "nsxt_policy_route_controller" "rc" {
  display_name                           = "tf-acc-intf-rc"
  ha_mode                                = "ACTIVE_STANDBY"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}
`, getRCVNAName())
}

func testAccNsxtPolicyRCInterfaceTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyRCInterfaceCreateAttributes
	} else {
		attrMap = accTestPolicyRCInterfaceUpdateAttributes
	}
	return testAccNsxtPolicyRCInterfaceRouteControllerTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_interface" "test" {
  display_name        = "%s"
  description         = "%s"
  parent_path         = nsxt_policy_route_controller.rc.path
  urpf_mode           = "%s"
  floating_ip_subnets = ["%s"]

  interface_address {
    subnets                        = ["%s"]
    portgroup_id                   = "%s"
    virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance.vna.path
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_route_controller_interface" "test" {
  display_name = "%s"
  parent_path  = nsxt_policy_route_controller.rc.path

  depends_on = [nsxt_policy_route_controller_interface.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["urpf_mode"], attrMap["floating_ip_subnet"], attrMap["subnet"], getPortgroupID(), attrMap["display_name"])
}

func testAccNsxtPolicyRCInterfaceMinimalistic() string {
	attrMap := accTestPolicyRCInterfaceImportAttributes
	return testAccNsxtPolicyRCInterfaceRouteControllerTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_interface" "test" {
  display_name        = "%s"
  parent_path         = nsxt_policy_route_controller.rc.path
  floating_ip_subnets = ["%s"]

  interface_address {
    subnets                        = ["%s"]
    portgroup_id                   = "%s"
    virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance.vna.path
  }
}`, attrMap["display_name"], attrMap["floating_ip_subnet"], attrMap["subnet"], getPortgroupID())
}
