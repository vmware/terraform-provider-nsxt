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

var accTestPolicyDnsRuleCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"domain_patterns":  "*.create.example.com",
	"upstream_servers": "8.8.8.8",
}

var accTestPolicyDnsRuleUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"domain_patterns":  "*.update.example.com",
	"upstream_servers": "8.8.4.4",
}

func TestAccResourceNsxtPolicyDnsRule_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRuleCheckDestroy(state, accTestPolicyDnsRuleUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRuleTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRuleCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRuleCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "action_type", "FORWARD"),
					resource.TestCheckResourceAttr(testResourceName, "domain_patterns.0", accTestPolicyDnsRuleCreateAttributes["domain_patterns"]),
					resource.TestCheckResourceAttr(testResourceName, "upstream_servers.0", accTestPolicyDnsRuleCreateAttributes["upstream_servers"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDnsRuleTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDnsRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDnsRuleUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDnsRuleUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "domain_patterns.0", accTestPolicyDnsRuleUpdateAttributes["domain_patterns"]),
					resource.TestCheckResourceAttr(testResourceName, "upstream_servers.0", accTestPolicyDnsRuleUpdateAttributes["upstream_servers"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDnsRule_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_dns_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.2.0"); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDnsRuleCheckDestroy(state, accTestPolicyDnsRuleCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDnsRuleTemplate(true),
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

func testAccNsxtPolicyDnsRuleExistsFn(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	return resourceNsxtPolicyDnsRuleExists(sessionContext, parentPath, id, connector)
}

func testAccNsxtPolicyDnsRuleExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DNS Rule resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DNS Rule resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := testAccNsxtPolicyDnsRuleExistsFn(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DNS Rule %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyDnsRuleCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_dns_rule" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := testAccNsxtPolicyDnsRuleExistsFn(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy DNS Rule %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDnsRuleTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDnsRuleCreateAttributes
	} else {
		attrMap = accTestPolicyDnsRuleUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_dns_service" "parent" {
  %s
  display_name           = "%s-svc"
  allocated_listener_ips = [var.dns_listener_ip]
  vns_clusters           = [var.vns_cluster_path]
}

resource "nsxt_policy_dns_rule" "test" {
  parent_path      = nsxt_policy_dns_service.parent.path
  display_name     = "%s"
  description      = "%s"
  action_type      = "FORWARD"
  domain_patterns  = ["%s"]
  upstream_servers = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`,
		testAccNsxtProjectContext(),
		attrMap["display_name"],
		attrMap["display_name"],
		attrMap["description"],
		attrMap["domain_patterns"],
		attrMap["upstream_servers"],
	)
}
