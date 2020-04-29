/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"testing"
)

func TestAccResourceNsxtPolicySecurityPolicy_basic(t *testing.T) {
	name := "terraform-test"
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_security_policy.test"
	comments1 := "Acceptance test create"
	comments2 := "Acceptance test update"
	direction1 := "IN"
	direction2 := "OUT"
	proto1 := "IPV4"
	proto2 := "IPV4_IPV6"
	defaultAction := "ALLOW"
	tag1 := "abc"
	tag2 := "def"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyBasic(name, comments1),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", comments1),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyBasic(updatedName, comments2),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", comments2),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyWithRule(updatedName, direction1, proto1, tag1),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", direction1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", proto1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", tag1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyWithRule(updatedName, direction2, proto2, tag2),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", direction2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", proto2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", tag2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicySecurityPolicy_withDependencies(t *testing.T) {
	name := "terraform-test"
	testResourceName := "nsxt_policy_security_policy.test"
	defaultAction := "ALLOW"
	defaultDirection := "IN_OUT"
	defaultProtocol := "IPV4_IPV6"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyWithDepsCreate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "rule1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "rule2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.source_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.destination_groups.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.services.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyWithDepsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Application"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "rule1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DROP"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "0"),
				),
			},
		},
	})
}
func TestAccResourceNsxtPolicySecurityPolicy_importBasic(t *testing.T) {
	name := fmt.Sprintf("terraform-test-import")
	testResourceName := "nsxt_policy_security_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyBasic(name, "import"),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicySecurityPolicyExists(resourceName string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := domains.NewDefaultSecurityPoliciesClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy SecurityPolicy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy SecurityPolicy resource ID not set in resources")
		}

		_, err := nsxClient.Get(domainName, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy SecurityPolicy ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicySecurityPolicyCheckDestroy(state *terraform.State, displayName string, domainName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := domains.NewDefaultSecurityPoliciesClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_security_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(domainName, resourceID)
		if err == nil {
			return fmt.Errorf("Policy SecurityPolicy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicySecurityPolicyBasic(name string, comments string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_security_policy" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "Application"
  comments        = "%s"
  locked          = true
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

}`, name, comments)
}

func testAccNsxtPolicySecurityPolicyWithRule(name string, direction string, protocol string, ruleTag string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_security_policy" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "Application"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name = "%s"
    direction    = "%s"
    ip_version   = "%s"
    log_label    = "%s"

    tag {
      scope = "color"
      tag   = "blue"
    }
  }
}`, name, name, direction, protocol, ruleTag)
}

func testAccNsxtPolicySecurityPolicyDeps() string {
	return `
resource "nsxt_policy_group" "group1" {
  display_name = "terraform testacc 1"
}

resource "nsxt_policy_group" "group2" {
  display_name = "terraform testacc 2"
}

resource "nsxt_policy_service" "icmp" {
    display_name = "security-policy-test-icmp"
    icmp_entry {
        protocol = "ICMPv4"
    }
}

resource "nsxt_policy_service" "tcp778" {
    display_name = "security-policy-test-tcp"
    l4_port_set_entry {
        protocol          = "TCP"
        destination_ports = [ "778" ]
    }
}`
}

//TODO: add profiles when available
func testAccNsxtPolicySecurityPolicyWithDepsCreate(name string) string {
	return testAccNsxtPolicySecurityPolicyDeps() + fmt.Sprintf(`
resource "nsxt_policy_security_policy" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "Application"
  locked          = false
  sequence_number = 3
  stateful        = true
  tcp_strict      = false
  scope           = [nsxt_policy_group.group1.path]

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name          = "rule1"
    source_groups         = [nsxt_policy_group.group1.path]
    destination_groups    = [nsxt_policy_group.group2.path]
    sources_excluded      = true
    destinations_excluded = true
    services              = [nsxt_policy_service.icmp.path, nsxt_policy_service.tcp778.path]
  }

  rule {
    display_name          = "rule2"
    source_groups         = [nsxt_policy_group.group1.path, nsxt_policy_group.group2.path]
    sources_excluded      = false
    destinations_excluded = false
  }
}`, name)
}

func testAccNsxtPolicySecurityPolicyWithDepsUpdate(name string) string {
	return testAccNsxtPolicySecurityPolicyDeps() + fmt.Sprintf(`
resource "nsxt_policy_security_policy" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "Application"
  locked          = true
  sequence_number = 3
  stateful        = true
  tcp_strict      = false
  scope           = [nsxt_policy_group.group2.path]

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name          = "rule1"
    destination_groups    = [nsxt_policy_group.group1.path, nsxt_policy_group.group2.path]
    sources_excluded      = false
    destinations_excluded = false
    disabled              = true
    action                = "DROP"
    logged                = true
  }

}`, name)
}
