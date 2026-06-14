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

// Each parallel test uses its own VNI and name to avoid pool-overlap conflicts.
var accTestPolicyDistributedVxlanConnectionCreateAttributes = map[string]string{
	"display_name":        getAccTestResourceName(),
	"description":         "terraform created",
	"l3_vni":              "5000",
	"route_distinguisher": "65000:5000",
}

var accTestPolicyDistributedVxlanConnectionUpdateAttributes = map[string]string{
	"display_name":        getAccTestResourceName(),
	"description":         "terraform updated",
	"l3_vni":              "5000",
	"route_distinguisher": "65000:5000",
}

// Separate VNI (6000) so _importBasic does not collide with _basic (VNI 5000).
var accTestPolicyDistributedVxlanConnectionImportAttributes = map[string]string{
	"display_name":        getAccTestResourceName(),
	"l3_vni":              "6000",
	"route_distinguisher": "65000:6000",
}

// Separate VNI (5001) so _withRouteTarget does not collide with the other tests.
var accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes = map[string]string{
	"display_name":        getAccTestResourceName(),
	"description":         "terraform created with route target",
	"l3_vni":              "5001",
	"route_distinguisher": "65001:5001",
}

func testAccNsxtPolicyDistributedVxlanConnectionPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_RC_VNA_CLUSTER_NAME")
}

// testAccNsxtPolicyDistributedVxlanConnectionRcTemplate returns HCL that looks
// up the pre-created ROUTE_CONTROLLER VNA cluster and creates an inline
// nsxt_policy_route_controller named "rc" on top of it. All DVXC templates
// prepend this block and reference nsxt_policy_route_controller.rc.path as
// route_controller_path.
func testAccNsxtPolicyDistributedVxlanConnectionRcTemplate() string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + `
resource "nsxt_policy_route_controller" "rc" {
  display_name                           = "tf-acc-dvxc-rc"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}
`
}

func TestAccResourceNsxtPolicyDistributedVxlanConnection_basic(t *testing.T) {
	testResourceName := "nsxt_policy_distributed_vxlan_connection.test"
	testDataSourceName := "data.nsxt_policy_distributed_vxlan_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyDistributedVxlanConnectionPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDistributedVxlanConnectionCheckDestroy(state, accTestPolicyDistributedVxlanConnectionUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDistributedVxlanConnectionTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDistributedVxlanConnectionExists(accTestPolicyDistributedVxlanConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDistributedVxlanConnectionCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDistributedVxlanConnectionCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "l3_vni", accTestPolicyDistributedVxlanConnectionCreateAttributes["l3_vni"]),
					resource.TestCheckResourceAttr(testResourceName, "connectivity_type", "L3_EVPN"),
					resource.TestCheckResourceAttr(testResourceName, "route_distinguisher", accTestPolicyDistributedVxlanConnectionCreateAttributes["route_distinguisher"]),
					resource.TestCheckResourceAttrSet(testResourceName, "route_controller_path"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.address_family", "L2VPN_EVPN"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.import_targets.0", "65000:5000"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.export_targets.0", "65000:5000"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyDistributedVxlanConnectionTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDistributedVxlanConnectionExists(accTestPolicyDistributedVxlanConnectionUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDistributedVxlanConnectionUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDistributedVxlanConnectionUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "l3_vni", accTestPolicyDistributedVxlanConnectionUpdateAttributes["l3_vni"]),
					resource.TestCheckResourceAttr(testResourceName, "route_distinguisher", accTestPolicyDistributedVxlanConnectionUpdateAttributes["route_distinguisher"]),
					resource.TestCheckResourceAttrSet(testResourceName, "route_controller_path"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyDistributedVxlanConnectionMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDistributedVxlanConnectionExists(accTestPolicyDistributedVxlanConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "route_controller_path"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.#", "1"),
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

func TestAccResourceNsxtPolicyDistributedVxlanConnection_withRouteTarget(t *testing.T) {
	testResourceName := "nsxt_policy_distributed_vxlan_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyDistributedVxlanConnectionPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDistributedVxlanConnectionCheckDestroy(state, accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDistributedVxlanConnectionWithRouteTarget(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDistributedVxlanConnectionExists(accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "l3_vni", accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes["l3_vni"]),
					resource.TestCheckResourceAttr(testResourceName, "route_distinguisher", accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes["route_distinguisher"]),
					resource.TestCheckResourceAttrSet(testResourceName, "route_controller_path"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.address_family", "L2VPN_EVPN"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.import_targets.0", "65001:100"),
					resource.TestCheckResourceAttr(testResourceName, "route_target.0.export_targets.0", "65001:200"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDistributedVxlanConnection_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_distributed_vxlan_connection.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyDistributedVxlanConnectionPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDistributedVxlanConnectionCheckDestroy(state, accTestPolicyDistributedVxlanConnectionImportAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDistributedVxlanConnectionImportMinimalistic(),
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

func testAccNsxtPolicyDistributedVxlanConnectionExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DistributedVxlanConnection resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DistributedVxlanConnection resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyDistributedVxlanConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DistributedVxlanConnection %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyDistributedVxlanConnectionCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_distributed_vxlan_connection" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDistributedVxlanConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy DistributedVxlanConnection %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDistributedVxlanConnectionTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDistributedVxlanConnectionCreateAttributes
	} else {
		attrMap = accTestPolicyDistributedVxlanConnectionUpdateAttributes
	}
	return testAccNsxtPolicyDistributedVxlanConnectionRcTemplate() + fmt.Sprintf(`
resource "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name          = "%s"
  description           = "%s"
  l3_vni                = %s
  connectivity_type     = "L3_EVPN"
  route_controller_path = nsxt_policy_route_controller.rc.path
  route_distinguisher   = "%s"

  route_target {
    address_family = "L2VPN_EVPN"
    import_targets = ["65000:5000"]
    export_targets = ["65000:5000"]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name = "%s"

  depends_on = [nsxt_policy_distributed_vxlan_connection.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["l3_vni"],
		attrMap["route_distinguisher"], attrMap["display_name"])
}

// testAccNsxtPolicyDistributedVxlanConnectionMinimalistic is used as step 3 of _basic
// (same VNI 5000 as the create/update steps, minimal attributes).
func testAccNsxtPolicyDistributedVxlanConnectionMinimalistic() string {
	attrMap := accTestPolicyDistributedVxlanConnectionCreateAttributes
	return testAccNsxtPolicyDistributedVxlanConnectionRcTemplate() + fmt.Sprintf(`
resource "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name          = "%s"
  l3_vni                = %s
  route_controller_path = nsxt_policy_route_controller.rc.path
  route_distinguisher   = "%s"

  route_target {
    import_targets = ["65000:5000"]
    export_targets = ["65000:5000"]
  }
}

data "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name = "%s"

  depends_on = [nsxt_policy_distributed_vxlan_connection.test]
}`, attrMap["display_name"], attrMap["l3_vni"], attrMap["route_distinguisher"],
		attrMap["display_name"])
}

// testAccNsxtPolicyDistributedVxlanConnectionImportMinimalistic is used by _importBasic
// and intentionally uses a different VNI (6000) to avoid parallel conflicts with _basic (VNI 5000).
func testAccNsxtPolicyDistributedVxlanConnectionImportMinimalistic() string {
	attrMap := accTestPolicyDistributedVxlanConnectionImportAttributes
	return testAccNsxtPolicyDistributedVxlanConnectionRcTemplate() + fmt.Sprintf(`
resource "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name          = "%s"
  l3_vni                = %s
  route_controller_path = nsxt_policy_route_controller.rc.path
  route_distinguisher   = "%s"

  route_target {
    import_targets = ["65000:6000"]
    export_targets = ["65000:6000"]
  }
}`, attrMap["display_name"], attrMap["l3_vni"], attrMap["route_distinguisher"])
}

func testAccNsxtPolicyDistributedVxlanConnectionWithRouteTarget() string {
	attrMap := accTestPolicyDistributedVxlanConnectionWithRouteTargetAttributes
	return testAccNsxtPolicyDistributedVxlanConnectionRcTemplate() + fmt.Sprintf(`
resource "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name          = "%s"
  description           = "%s"
  l3_vni                = %s
  route_controller_path = nsxt_policy_route_controller.rc.path
  route_distinguisher   = "%s"

  route_target {
    address_family = "L2VPN_EVPN"
    import_targets = ["65001:100"]
    export_targets = ["65001:200"]
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["l3_vni"],
		attrMap["route_distinguisher"])
}
