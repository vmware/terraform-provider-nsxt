// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

var accTestPolicyDnsZoneCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"dns_domain_name": "create.example.com",
	"ttl":             "600",
}

var accTestPolicyDnsZoneUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"dns_domain_name": "create.example.com",
	"ttl":             "3600",
}

func TestAccResourceNsxtPolicyDnsZone_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsZoneCheckDestroy(state, accTestPolicyDnsZoneUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsZoneTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsZoneExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsZoneCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsZoneCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_domain_name", accTestPolicyDnsZoneCreateAttributes["dns_domain_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsZoneCreateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDnsZoneTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsZoneExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsZoneUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsZoneUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ttl", accTestPolicyDnsZoneUpdateAttributes["ttl"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsZone_withSoa(t *testing.T) {
	testResourceName := "nsxt_policy_dns_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsZoneCheckDestroy(state, accTestPolicyDnsZoneCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsZoneWithSoaTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsZoneExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "soa.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "soa.0.primary_nameserver", "ns1.example.com."),
					resource.TestCheckResourceAttr(testResourceName, "soa.0.responsible_party", "admin.example.com."),
					resource.TestCheckResourceAttr(testResourceName, "soa.0.refresh_interval", "3600"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsZone_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsZoneCheckDestroy(state, accTestPolicyDnsZoneCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsZoneTemplate(true),
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

func testAccNsxtPolicyDnsZoneExistsFn(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	return resourceNsxtPolicyDnsZoneExists(sessionContext, parentPath, id, connector)
}

func testAccNsxtPolicyDnsZoneExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DNS Zone resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DNS Zone resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := testAccNsxtPolicyDnsZoneExistsFn(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DNS Zone %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyDnsZoneCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_dns_zone" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := testAccNsxtPolicyDnsZoneExistsFn(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy DNS Zone %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDnsZoneTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDnsZoneCreateAttributes
	} else {
		attrMap = accTestPolicyDnsZoneUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_dns_service" "parent" {
  %s
  display_name           = "%s-svc"
  allocated_listener_ips = [var.dns_listener_ip]
  vns_clusters           = [var.vns_cluster_path]
}

resource "nsxt_policy_dns_zone" "test" {
  parent_path     = nsxt_policy_dns_service.parent.path
  display_name    = "%s"
  description     = "%s"
  dns_domain_name = "%s"
  ttl             = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`,
		testAccNsxtProjectContext(),
		attrMap["display_name"],
		attrMap["display_name"],
		attrMap["description"],
		attrMap["dns_domain_name"],
		attrMap["ttl"],
	)
}

func testAccNsxtPolicyDnsZoneWithSoaTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_dns_service" "parent" {
  %s
  display_name           = "%s-svc"
  allocated_listener_ips = [var.dns_listener_ip]
  vns_clusters           = [var.vns_cluster_path]
}

resource "nsxt_policy_dns_zone" "test" {
  parent_path     = nsxt_policy_dns_service.parent.path
  display_name    = "%s"
  dns_domain_name = "soa.example.com"
  ttl             = 600

  soa {
    primary_nameserver = "ns1.example.com."
    responsible_party  = "admin.example.com."
    refresh_interval   = 3600
    retry_interval     = 900
    expire_time        = 604800
    negative_cache_ttl = 300
  }
}`,
		testAccNsxtProjectContext(),
		accTestPolicyDnsZoneCreateAttributes["display_name"],
		accTestPolicyDnsZoneCreateAttributes["display_name"],
	)
}
