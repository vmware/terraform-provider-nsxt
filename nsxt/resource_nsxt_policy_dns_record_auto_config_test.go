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

var accTestPolicyDnsRecordAutoConfigCreateAttributes = map[string]string{
	"display_name":   getAccTestResourceName(),
	"description":    "terraform created",
	"naming_pattern": "{vm_name}",
	"ttl":            "600",
}

var accTestPolicyDnsRecordAutoConfigUpdateAttributes = map[string]string{
	"display_name":   getAccTestResourceName(),
	"description":    "terraform updated",
	"naming_pattern": "{vm_name}-{index}",
	"ttl":            "3600",
}

func TestAccResourceNsxtPolicyDnsRecordAutoConfig_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_record_auto_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRecordAutoConfigCheckDestroy(state, accTestPolicyDnsRecordAutoConfigUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRecordAutoConfigTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRecordAutoConfigExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRecordAutoConfigCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRecordAutoConfigCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "naming_pattern", accTestPolicyDnsRecordAutoConfigCreateAttributes["naming_pattern"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsRecordAutoConfigCreateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_block_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "zone_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDnsRecordAutoConfigTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRecordAutoConfigExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRecordAutoConfigUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRecordAutoConfigUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "naming_pattern", accTestPolicyDnsRecordAutoConfigUpdateAttributes["naming_pattern"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsRecordAutoConfigUpdateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsRecordAutoConfig_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_record_auto_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRecordAutoConfigCheckDestroy(state, accTestPolicyDnsRecordAutoConfigCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRecordAutoConfigTemplate(true),
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

func testAccNsxtPolicyDnsRecordAutoConfigExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DNS Record Auto Config resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DNS Record Auto Config resource ID not set in resources")
		}
		exists, err := resourceNsxtPolicyDnsRecordAutoConfigExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DNS Record Auto Config %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyDnsRecordAutoConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_dns_record_auto_config" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDnsRecordAutoConfigExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy DNS Record Auto Config %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDnsRecordAutoConfigTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDnsRecordAutoConfigCreateAttributes
	} else {
		attrMap = accTestPolicyDnsRecordAutoConfigUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_dns_service" "parent_svc" {
  %s
  display_name           = "%s-svc"
  allocated_listener_ips = [var.dns_listener_ip]
  vns_clusters           = [var.vns_cluster_path]
}

resource "nsxt_policy_dns_zone" "parent_zone" {
  parent_path     = nsxt_policy_dns_service.parent_svc.path
  display_name    = "%s-zone"
  dns_domain_name = "auto.example.com"
  ttl             = 300
}

resource "nsxt_policy_dns_record_auto_config" "test" {
  %s
  display_name   = "%s"
  description    = "%s"
  ip_block_path  = var.dns_ip_block_path
  zone_path      = nsxt_policy_dns_zone.parent_zone.path
  naming_pattern = "%s"
  ttl            = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`,
		testAccNsxtProjectContext(),
		attrMap["display_name"],
		attrMap["display_name"],
		testAccNsxtProjectContext(),
		attrMap["display_name"],
		attrMap["description"],
		attrMap["naming_pattern"],
		attrMap["ttl"],
	)
}
