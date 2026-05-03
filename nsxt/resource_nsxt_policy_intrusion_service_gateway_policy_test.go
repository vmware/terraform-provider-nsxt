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

func TestAccResourceNsxtPolicyIntrusionServiceGatewayPolicy_basic(t *testing.T) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	tierInfraName := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy"
	comments1 := "Acceptance test create"
	comments2 := "Acceptance test update"
	direction1 := "IN"
	direction2 := "OUT"
	proto1 := "IPV4"
	proto2 := "IPV4_IPV6"
	defaultAction := "DETECT"
	tag1 := "ids-gw-detect-label"
	tag2 := "ids-gw-detect-label-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServiceGatewayPolicyCheckDestroy(state, updatedName, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyBasic(name, comments1),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", comments1),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyBasic(updatedName, comments2),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", comments2),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyWithRule(updatedName, direction1, proto1, tag1, tierInfraName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "ids-gw-detect-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", direction1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", proto1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", tag1),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyWithRule(updatedName, direction2, proto2, tag2, tierInfraName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "comments", ""),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "ids-gw-detect-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", direction2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", proto2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", tag2),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIntrusionServiceGatewayPolicy_withDependencies(t *testing.T) {
	name := getAccTestResourceName()
	tierInfraName := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy"
	defaultAction := "DETECT"
	defaultDirection := "IN_OUT"
	defaultProtocol := "IPV4_IPV6"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServiceGatewayPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDepsCreate(name, tierInfraName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "locked", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "ids-gw-rule-with-exclusions"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "ids-gw-rule-multi-source"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action", defaultAction),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.source_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.services.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.ids_profiles.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDepsUpdate(name, tierInfraName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "locked", "true"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "ids-gw-rule-with-exclusions"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", defaultDirection),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", defaultProtocol),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT_PREVENT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIntrusionServiceGatewayPolicy_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServiceGatewayPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyBasic(name, "import"),
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

func testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(resourceName string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Intrusion Service Gateway Policy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Intrusion Service Gateway Policy resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domainName, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving Intrusion Service Gateway Policy resource ID %s", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyCheckDestroy(state *terraform.State, displayName string, domainName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_intrusion_service_gateway_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domainName, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Intrusion Service Gateway Policy %s still exists", displayName)
		}
	}
	return nil
}

// testAccNsxtPolicyIntrusionServiceGatewayPolicyT1Prerequisites creates Tier-0/Tier-1
// used as IDS gateway policy rule scope (same pattern as other policy acceptance tests).
func testAccNsxtPolicyIntrusionServiceGatewayPolicyT1Prerequisites(uniquePrefix string) string {
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) + fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "t0_for_ids" {
  display_name = "%[1]s-t0"
  description    = "Terraform acceptance IDS gateway policy T0"
  ha_mode        = "ACTIVE_STANDBY"
}

resource "nsxt_policy_tier1_gateway" "t1_gw" {
  display_name = "%[1]s-t1"
  description    = "Terraform acceptance IDS gateway policy T1"
  locale_service {
    edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  }
  tier0_path = nsxt_policy_tier0_gateway.t0_for_ids.path
}
`, uniquePrefix)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyBasic(name string, comments string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  comments        = "%s"
  locked          = true
  sequence_number = 3

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }
}`, name, comments)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyWithRule(name string, direction string, protocol string, ruleTag string, tierInfraPrefix string) string {
	return testAccNsxtPolicyIntrusionServiceGatewayPolicyT1Prerequisites(tierInfraPrefix) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }

  rule {
    display_name = "ids-gw-detect-rule"
    direction    = "%s"
    ip_version   = "%s"
    scope        = [nsxt_policy_tier1_gateway.t1_gw.path]
    log_label    = "%s"
    ids_profiles = ["%s"]

    tag {
      scope = "rule-env"
      tag   = "acceptance-test"
    }
  }
}`, name, direction, protocol, ruleTag, policyDefaultIdsProfilePath)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyDeps(tierInfraPrefix string) string {
	return testAccNsxtPolicyIntrusionServiceGatewayPolicyT1Prerequisites(tierInfraPrefix) + `
resource "nsxt_policy_group" "src_group" {
  display_name = "tf-intrusion-svc-gw-policy-src-group"
}

resource "nsxt_policy_group" "src_group_2" {
  display_name = "tf-intrusion-svc-gw-policy-src-group_2"
}

resource "nsxt_policy_group" "dst_group" {
  display_name = "tf-intrusion-svc-gw-policy-dst-group"
}

resource "nsxt_policy_group" "dst_group_2" {
  display_name = "tf-intrusion-svc-gw-policy-dst-group_2"
}

resource "nsxt_policy_service" "icmp_svc" {
  display_name = "tf-intrusion-svc-gw-policy-icmp"
  icmp_entry {
    protocol = "ICMPv4"
  }
}

resource "nsxt_policy_service" "tcp_svc" {
  display_name = "tf-intrusion-svc-gw-policy-tcp778"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["778"]
  }
}`
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDepsCreate(name string, tierInfraPrefix string) string {
	return testAccNsxtPolicyIntrusionServiceGatewayPolicyDeps(tierInfraPrefix) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 3

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }

  rule {
    display_name          = "ids-gw-rule-with-exclusions"
    source_groups         = [nsxt_policy_group.src_group.path]
    destination_groups    = [nsxt_policy_group.dst_group.path]
    sources_excluded      = true
    destinations_excluded = true
    scope                 = [nsxt_policy_tier1_gateway.t1_gw.path]
    services              = [nsxt_policy_service.icmp_svc.path, nsxt_policy_service.tcp_svc.path]
    ids_profiles          = ["%s"]
  }

  rule {
    display_name          = "ids-gw-rule-multi-source"
    source_groups         = [nsxt_policy_group.src_group.path, nsxt_policy_group.src_group_2.path]
    destination_groups    = [nsxt_policy_group.dst_group.path]
    sources_excluded      = false
    destinations_excluded = false
    scope                 = [nsxt_policy_tier1_gateway.t1_gw.path]
    ids_profiles          = ["%s"]
  }
}`, name, policyDefaultIdsProfilePath, policyDefaultIdsProfilePath)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDepsUpdate(name string, tierInfraPrefix string) string {
	return testAccNsxtPolicyIntrusionServiceGatewayPolicyDeps(tierInfraPrefix) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  locked          = true
  sequence_number = 3

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }

  rule {
    display_name          = "ids-gw-rule-with-exclusions"
    destination_groups    = [nsxt_policy_group.dst_group.path, nsxt_policy_group.dst_group_2.path]
    sources_excluded      = false
    destinations_excluded = false
    disabled              = true
    action                = "DETECT_PREVENT"
    logged                = true
    scope                 = [nsxt_policy_tier1_gateway.t1_gw.path]
    ids_profiles          = ["%s"]
  }
}`, name, policyDefaultIdsProfilePath)
}
