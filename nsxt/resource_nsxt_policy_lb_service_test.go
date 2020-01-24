/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

var accTestPolicyLBServiceCreateAttributes = map[string]string{
	"display_name":      "terra-test",
	"description":       "terraform created",
	"connectivity_path": "nsxt_policy_tier1_gateway.test1.path",
	"enabled":           "true",
	"error_log_level":   "ERROR",
	"size":              "SMALL",
}

var accTestPolicyLBServiceUpdateAttributes = map[string]string{
	"display_name":      "terra-test-updated",
	"description":       "terraform updated",
	"connectivity_path": "nsxt_policy_tier1_gateway.test2.path",
	"enabled":           "false",
	"error_log_level":   "EMERGENCY",
	"size":              "MEDIUM",
}

func TestAccResourceNsxtPolicyLBService_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_service.test"
	edgeClusterName := getEdgeClusterName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBServiceCheckDestroy(state, accTestPolicyLBServiceCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBServiceTemplate(edgeClusterName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBServiceCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBServiceCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "connectivity_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyLBServiceCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "error_log_level", accTestPolicyLBServiceCreateAttributes["error_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "size", accTestPolicyLBServiceCreateAttributes["size"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBServiceTemplate(edgeClusterName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBServiceUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBServiceUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "connectivity_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyLBServiceUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "error_log_level", accTestPolicyLBServiceUpdateAttributes["error_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "size", accTestPolicyLBServiceUpdateAttributes["size"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBServiceMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServiceExists(testResourceName),
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

func TestAccResourceNsxtPolicyLBService_importBasic(t *testing.T) {
	name := "terra-test-import"
	testResourceName := "nsxt_policy_lb_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBServiceMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBServiceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := infra.NewDefaultLbServicesClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBService resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBService resource ID not set in resources")
		}

		_, err := nsxClient.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy LBService ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyLBServiceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := infra.NewDefaultLbServicesClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy LBService %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBServiceTemplate(edgeClusterName string, createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBServiceCreateAttributes
	} else {
		attrMap = accTestPolicyLBServiceUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "terraform-lb-test"
}

resource "nsxt_policy_tier1_gateway" "test1" {
  display_name      = "terraform-lb-test-1"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  tier0_path        = nsxt_policy_tier0_gateway.test.path
}

resource "nsxt_policy_tier1_gateway" "test2" {
  display_name      = "terraform-lb-test-2"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  tier0_path        = nsxt_policy_tier0_gateway.test.path
}

resource "nsxt_policy_lb_service" "test" {
  display_name      = "%s"
  description       = "%s"
  connectivity_path = %s
  enabled           = %s
  error_log_level   = "%s"
  size              = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_service.test.path
}`, edgeClusterName, attrMap["display_name"], attrMap["description"], attrMap["connectivity_path"], attrMap["enabled"], attrMap["error_log_level"], attrMap["size"])
}

func testAccNsxtPolicyLBServiceMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_service" "test" {
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_service.test.path
}`, accTestPolicyLBServiceUpdateAttributes["display_name"])
}
