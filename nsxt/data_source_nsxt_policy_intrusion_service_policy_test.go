// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIntrusionServicePolicy_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Check that rules are present
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "0"),
					// Verify oversubscription is present for DFW IDPS rules
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.oversubscription"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicy_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate embedded rules
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-rule-id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicy_byCategory(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy.by_category"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyReadByCategory(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate embedded rules
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-rule-category"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicy_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyReadByNameMultitenancy(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate embedded rules
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-rule-multitenancy"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServicePolicyReadByName(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-rule"
    notes              = "Test rule for data source"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = []
    destination_groups = []
    services           = []
    scope              = []
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy" "by_name" {
  display_name = nsxt_policy_intrusion_service_policy.ids_policy_for_ds.display_name
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyIntrusionServicePolicyReadByCategory(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-rule-category"
    notes              = "Test rule for category-based lookup"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = []
    destination_groups = []
    services           = []
    scope              = []
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy" "by_category" {
  display_name = nsxt_policy_intrusion_service_policy.ids_policy_for_ds.display_name
  category     = "ThreatRules"
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyIntrusionServicePolicyReadByID(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-rule-id"
    notes              = "Test rule for ID-based lookup"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = []
    destination_groups = []
    services           = []
    scope              = []
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy" "by_id" {
  id         = nsxt_policy_intrusion_service_policy.ids_policy_for_ds.id
  depends_on = [nsxt_policy_intrusion_service_policy.ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyIntrusionServicePolicyReadByNameMultitenancy(name string) string {
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

resource "nsxt_policy_intrusion_service_policy" "ids_policy_for_ds" {
%s
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "ThreatRules"
  sequence_number = 3
  
  rule {
    display_name = "test-rule-multitenancy"
    notes        = "Test rule for multitenancy"
    action       = "DETECT"
    direction    = "IN_OUT"
    ids_profiles = [nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy" "by_name" {
%s
  display_name = nsxt_policy_intrusion_service_policy.ids_policy_for_ds.display_name
  category     = nsxt_policy_intrusion_service_policy.ids_policy_for_ds.category
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_for_ds]
}`, context, name, context, name, context)
}

func TestAccDataSourceNsxtPolicyIntrusionServicePolicy_withDetailedRules(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_policy.with_detailed_rules"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyWithDetailedRules(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Comprehensive IDPS policy test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Check that rules are present with detailed configuration
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					// First rule - comprehensive configuration
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "comprehensive-test-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT_PREVENT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", "IPV4"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.oversubscription", "DROPPED"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.sequence_number"),
					// Validate specific field arrays have expected number of items
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
					// Second rule - simpler configuration for comparison
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "simple-test-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.direction", "IN_OUT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.services.#", "1"),
					// Verify rule_id, nsx_id, and scope are present (from our latest alignment)
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.1.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.1.nsx_id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServicePolicyWithDetailedRules(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_group" "web_servers" {
  display_name = "test-web-servers-%s"
  description  = "Test web servers group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["192.168.1.10", "192.168.1.11", "192.168.1.12"]
    }
  }
}

resource "nsxt_policy_group" "db_servers" {
  display_name = "test-db-servers-%s"
  description  = "Test database servers group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.0.1.100-10.0.1.110"]
    }
  }
}

resource "nsxt_policy_group" "dmz_group" {
  display_name = "test-dmz-group-%s"
  description  = "Test DMZ group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.1.0/24"]
    }
  }
}

data "nsxt_policy_service" "ssh" {
  display_name = "SSH"
}

data "nsxt_policy_service" "https" {
  display_name = "HTTPS"
}

resource "nsxt_policy_service" "custom_app" {
  display_name = "test-custom-app-%s"
  description  = "Test custom application service"

  l4_port_set_entry {
    display_name      = "custom-app-port"
    description       = "Custom app on port 8080"
    protocol          = "TCP"
    destination_ports = ["8080"]
  }
}

resource "nsxt_policy_intrusion_service_policy" "ids_policy_detailed" {
  display_name    = "%s"
  description     = "Comprehensive IDPS policy test"
  category        = "ThreatRules"
  sequence_number = 5
  locked          = false

  rule {
    display_name          = "comprehensive-test-rule"
    description           = "Comprehensive rule testing all fields"
    notes                 = "Test rule with realistic configuration"
    action                = "DETECT_PREVENT"
    direction             = "IN"
    logged                = true
    disabled              = false
    sources_excluded      = false
    destinations_excluded = true
    ip_version            = "IPV4"
    oversubscription      = "DROPPED"
    sequence_number       = 10
    
    source_groups      = [nsxt_policy_group.web_servers.path, nsxt_policy_group.db_servers.path]
    destination_groups = [nsxt_policy_group.dmz_group.path]
    services          = [data.nsxt_policy_service.ssh.path, data.nsxt_policy_service.https.path]
    scope             = [nsxt_policy_group.dmz_group.path]
    ids_profiles      = [data.nsxt_policy_intrusion_service_profile.test.path]
    
    tag {
      scope = "environment"
      tag   = "production"
    }
    
    tag {
      scope = "team"
      tag   = "security"
    }
  }

  rule {
    display_name        = "simple-test-rule"
    description         = "Simple rule for comparison using custom service"
    action              = "DETECT"
    direction           = "IN_OUT"
    logged              = false
    sequence_number     = 20
    
    source_groups      = [nsxt_policy_group.web_servers.path]
    destination_groups = [nsxt_policy_group.db_servers.path] 
    services           = [nsxt_policy_service.custom_app.path]
    scope              = [nsxt_policy_group.dmz_group.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_policy" "with_detailed_rules" {
  display_name = nsxt_policy_intrusion_service_policy.ids_policy_detailed.display_name
  depends_on   = [nsxt_policy_intrusion_service_policy.ids_policy_detailed]
}`, name, name, name, name, name)
}
