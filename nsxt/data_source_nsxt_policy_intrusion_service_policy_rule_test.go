// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIntrusionServicePolicyRule_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy_rule.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicyRule_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy_rule.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicyRule_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy_rule.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleReadByNameMultitenancy(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicyRule_withDetailedConfiguration(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy_rule.detailed_rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0") // EXEMPT requires 9.1.0+
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleWithDetailedConfiguration(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Comprehensive standalone IDPS rule test"),
					resource.TestCheckResourceAttr(testResourceName, "action", "EXEMPT"),
					resource.TestCheckResourceAttr(testResourceName, "direction", "IN"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sources_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "ip_version", "IPV4_IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "oversubscription", "INHERIT_GLOBAL"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate detailed field configuration
					resource.TestCheckResourceAttr(testResourceName, "source_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "destination_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					// Verify rule_id is present
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServicePolicyRuleReadByName(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_rule_ds" {
  display_name    = "%s-policy"
  description     = "Acceptance Test Policy"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "%s"
    description        = "Acceptance Test Rule"
    action             = "DETECT"
    logged             = true
    disabled           = false
    sequence_number    = 1
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy_rule" "by_name" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds.path
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds]
}`, name, name, name)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleReadByID(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_rule_ds" {
  display_name    = "%s-policy"
  description     = "Acceptance Test Policy"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "%s"
    description        = "Acceptance Test Rule"
    action             = "DETECT"
    logged             = true
    disabled           = false
    source_groups      = []
    destination_groups = []
    services           = []
    scope              = []
    sequence_number    = 1
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy_rule" "by_id" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds.path
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds]
}`, name, name, name)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleReadByNameMultitenancy(name string) string {
	context := testAccNsxtPolicyMultitenancyContext()
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_profile" "test" {
%s
  display_name = "%s-profile"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    products_affected = ["Linux"]
  }
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_rule_ds" {
%s
  display_name    = "%s-policy"
  description     = "Acceptance Test Policy"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "%s"
    description        = "Acceptance Test Rule"
    action             = "DETECT"
    logged             = true
    disabled           = false
    sequence_number    = 1
    ids_profiles       = [nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy_rule" "by_name" {
%s
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds.path
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_rule_ds]
}`, context, name, context, name, name, context, name)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleWithDetailedConfiguration(name string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_group" "app_servers" {
  display_name = "test-app-servers-%s"
  description  = "Test application servers"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.1.1.10-10.1.1.20"]
    }
  }
}

resource "nsxt_policy_group" "db_tier" {
  display_name = "test-db-tier-%s"
  description  = "Test database tier"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.1.2.96/28"]
    }
  }
}

resource "nsxt_policy_group" "web_dmz" {
  display_name = "test-web-dmz-%s"
  description  = "Test web DMZ group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["192.168.100.0/24"]
    }
  }
}

resource "nsxt_policy_group" "mgmt_network" {
  display_name = "test-mgmt-network-%s"
  description  = "Test management network"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.200.0/24"]
    }
  }
}

data "nsxt_policy_service" "mysql" {
  display_name = "MySQL"
}

data "nsxt_policy_service" "postgresql" {
  display_name = "PostgreSQL"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_detailed_rule" {
  display_name    = "%s-policy"
  description     = "Parent policy for detailed rule test"
  category        = "ThreatRules"
  sequence_number = 7
  
  rule {
    display_name          = "%s"
    description           = "Comprehensive standalone IDPS rule test"
    notes                 = "Testing all available IDPS rule fields with realistic values"
    action                = "EXEMPT"
    direction             = "IN"
    logged                = true
    log_label             = "security-audit"
    disabled              = false
    sources_excluded      = true
    destinations_excluded = false
    ip_version            = "IPV4_IPV6"
    oversubscription      = "INHERIT_GLOBAL"
    sequence_number       = 30
    
    source_groups      = [nsxt_policy_group.app_servers.path, nsxt_policy_group.web_dmz.path]
    destination_groups = [nsxt_policy_group.db_tier.path, nsxt_policy_group.mgmt_network.path]
    services          = [data.nsxt_policy_service.mysql.path, data.nsxt_policy_service.postgresql.path]
    scope             = [nsxt_policy_group.mgmt_network.path]
    ids_profiles      = ["/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile"]
    
    tag {
      scope = "criticality"
      tag   = "high"
    }
    
    tag {
      scope = "compliance"  
      tag   = "pci_dss"
    }
  }
}

data "nsxt_policy_intrusion_service_policy_rule" "detailed_rule" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_policy.ids_policy_for_detailed_rule.path
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_detailed_rule]
}`, name, name, name, name, name, name, name)
}
