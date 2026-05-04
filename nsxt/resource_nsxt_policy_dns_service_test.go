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

var accTestPolicyDnsServiceCreateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform created",
	"allocated_listener_ips": "/orgs/default/projects/${data.nsxt_policy_project.test.id}/ip-address-allocations/listener-ip1",
}

var accTestPolicyDnsServiceUpdateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform updated",
	"allocated_listener_ips": "/orgs/default/projects/${data.nsxt_policy_project.test.id}/ip-address-allocations/listener-ip1",
}

func TestAccResourceNsxtPolicyDnsService_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsServiceCheckDestroy(state, accTestPolicyDnsServiceUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsServiceTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsServiceCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsServiceCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "allocated_listener_ips.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "vns_clusters.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDnsServiceTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsServiceUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsServiceUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsService_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsServiceCheckDestroy(state, accTestPolicyDnsServiceCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsServiceTemplate(true),
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

func testAccNsxtPolicyDnsServiceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DNS Service resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DNS Service resource ID not set in resources")
		}
		exists, err := resourceNsxtPolicyDnsServiceExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DNS Service %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyDnsServiceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_dns_service" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDnsServiceExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy DNS Service %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDnsServiceTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDnsServiceCreateAttributes
	} else {
		attrMap = accTestPolicyDnsServiceUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_project" "test" {
  %s
  display_name = "%s"
}

resource "nsxt_policy_dns_service" "test" {
  %s
  display_name = "%s"
  description  = "%s"

  allocated_listener_ips = [var.dns_listener_ip]
  vns_clusters           = [var.vns_cluster_path]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`,
		testAccNsxtProjectContext(), getAccTestResourceName(),
		testAccNsxtProjectContext(),
		attrMap["display_name"],
		attrMap["description"],
	)
}
