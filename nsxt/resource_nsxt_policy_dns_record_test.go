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

var accTestPolicyDnsRecordCreateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform created",
	"record_name":   "www",
	"record_values": "192.168.1.10",
	"ttl":           "600",
}

var accTestPolicyDnsRecordUpdateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform updated",
	"record_name":   "www",
	"record_values": "192.168.1.20",
	"ttl":           "3600",
}

func TestAccResourceNsxtPolicyDnsRecord_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_record.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRecordCheckDestroy(state, accTestPolicyDnsRecordUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRecordTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRecordExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRecordCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRecordCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "record_name", accTestPolicyDnsRecordCreateAttributes["record_name"]),
					resource.TestCheckResourceAttr(testResourceName, "record_type", "A"),
					resource.TestCheckResourceAttr(testResourceName, "record_values.0", accTestPolicyDnsRecordCreateAttributes["record_values"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsRecordCreateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "zone_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDnsRecordTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRecordExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRecordUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRecordUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "record_values.0", accTestPolicyDnsRecordUpdateAttributes["record_values"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsRecordUpdateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsRecord_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_record.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRecordCheckDestroy(state, accTestPolicyDnsRecordCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRecordTemplate(true),
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

func testAccNsxtPolicyDnsRecordExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DNS Record resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DNS Record resource ID not set in resources")
		}
		exists, err := resourceNsxtPolicyDnsRecordExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DNS Record %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyDnsRecordCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_dns_record" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDnsRecordExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy DNS Record %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDnsRecordTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDnsRecordCreateAttributes
	} else {
		attrMap = accTestPolicyDnsRecordUpdateAttributes
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
  dns_domain_name = "record.example.com"
  ttl             = 300
}

resource "nsxt_policy_dns_record" "test" {
  %s
  display_name  = "%s"
  description   = "%s"
  record_name   = "%s"
  record_type   = "A"
  record_values = ["%s"]
  zone_path     = nsxt_policy_dns_zone.parent_zone.path
  ttl           = %s

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
		attrMap["record_name"],
		attrMap["record_values"],
		attrMap["ttl"],
	)
}
