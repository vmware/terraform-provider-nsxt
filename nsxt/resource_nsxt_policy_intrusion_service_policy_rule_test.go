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

func TestAccResourceNsxtPolicyIntrusionServicePolicyRule_basic(t *testing.T) {
	testAccResourceNsxtPolicyIntrusionServicePolicyRuleBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func TestAccResourceNsxtPolicyIntrusionServicePolicyRule_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIntrusionServicePolicyRuleBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIntrusionServicePolicyRuleBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_policy_rule.rule1"
	action1 := "DETECT"
	action2 := "DETECT_PREVENT"
	direction1 := "IN"
	direction2 := "OUT"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServicePolicyRuleCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleBasicTemplate(withContext, name, action1, direction1),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "action", action1),
					resource.TestCheckResourceAttr(testResourceName, "direction", direction1),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleBasicTemplate(withContext, updatedName, action2, direction2),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "action", action2),
					resource.TestCheckResourceAttr(testResourceName, "direction", direction2),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIntrusionServicePolicyRule_withAllFields(t *testing.T) {
	testAccResourceNsxtPolicyIntrusionServicePolicyRuleWithAllFields(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func testAccResourceNsxtPolicyIntrusionServicePolicyRuleWithAllFields(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_policy_rule.rule1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServicePolicyRuleCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleWithAllFieldsCreate(withContext, name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Standalone rule with all fields"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "direction", "IN"),
					resource.TestCheckResourceAttr(testResourceName, "ip_version", "IPV4"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "log_label", "intrusion-svc-rule-detect"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sources_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "destinations_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "1"),
					resource.TestCheckResourceAttr(testResourceName, "source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleWithAllFieldsUpdate(withContext, name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyRuleExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Updated standalone rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT_PREVENT"),
					resource.TestCheckResourceAttr(testResourceName, "direction", "IN_OUT"),
					resource.TestCheckResourceAttr(testResourceName, "ip_version", "IPV4_IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "log_label", "intrusion-svc-rule-prevent"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", "2"),
					resource.TestCheckResourceAttr(testResourceName, "source_groups.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "destination_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "services.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIntrusionServicePolicyRule_importBasic(t *testing.T) {
	testAccResourceNsxtPolicyIntrusionServicePolicyRuleImportBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func TestAccResourceNsxtPolicyIntrusionServicePolicyRule_importBasic_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIntrusionServicePolicyRuleImportBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIntrusionServicePolicyRuleImportBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_policy_rule.rule1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServicePolicyRuleCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyRuleBasicTemplate(withContext, name, "DETECT", "IN"),
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

func testAccNsxtPolicyIntrusionServicePolicyRuleExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Intrusion Service Policy Rule resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Intrusion Service Policy Rule resource ID not set in resources")
		}

		policyPath := rs.Primary.Attributes["policy_path"]
		exists, err := resourceNsxtPolicyIntrusionServicePolicyRuleExists(testAccGetSessionContext(), resourceID, policyPath, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving Intrusion Service Policy Rule ID %s under policy %s", resourceID, policyPath)
		}
		return nil
	}
}

func testAccNsxtPolicyIntrusionServicePolicyRuleCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_intrusion_service_policy_rule" {
			continue
		}
		resourceID := rs.Primary.ID
		policyPath := rs.Primary.Attributes["policy_path"]
		exists, err := resourceNsxtPolicyIntrusionServicePolicyRuleExists(testAccGetSessionContext(), resourceID, policyPath, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Intrusion Service Policy Rule %s still exists", resourceID)
		}
	}
	return nil
}

func testAccNsxtPolicyIntrusionServicePolicyRuleParentDeps(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "parent" {
%s
  display_name    = "tf-intrusion-svc-parent-for-rule-test"
  locked          = false
  sequence_number = 3
  stateful        = true
}`, context)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleAllFieldsDeps(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return testAccNsxtPolicyIntrusionServicePolicyRuleParentDeps(withContext) + fmt.Sprintf(`
resource "nsxt_policy_group" "src_group" {
%s
  display_name = "tf-intrusion-svc-rule-src-group"
}

resource "nsxt_policy_group" "dst_group1" {
%s
  display_name = "tf-intrusion-svc-rule-dst-group1"
}

resource "nsxt_policy_group" "dst_group2" {
%s
  display_name = "tf-intrusion-svc-rule-dst-group2"
}

resource "nsxt_policy_service" "icmp_svc" {
  display_name = "tf-intrusion-svc-rule-icmp"
  icmp_entry {
    protocol = "ICMPv4"
  }
}

resource "nsxt_policy_service" "tcp_svc" {
  display_name = "tf-intrusion-svc-rule-tcp8443"
  l4_port_set_entry {
    protocol          = "TCP"
    destination_ports = ["8443"]
  }
}`, context, context, context)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleBasicTemplate(withContext bool, name string, action string, direction string) string {
	profilePath := fmt.Sprintf("\"%s\"", policyDefaultIdsProfilePath)
	if withContext {
		profilePath = "nsxt_policy_intrusion_service_profile.test.path"
		return testAccNsxtPolicyIntrusionServiceProfileMinimalistic(name, withContext) +
			testAccNsxtPolicyIntrusionServicePolicyRuleParentDeps(withContext) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_policy_rule" "rule1" {
%s
  display_name    = "%s"
  policy_path     = nsxt_policy_parent_intrusion_service_policy.parent.path
  action          = "%s"
  direction       = "%s"
  sequence_number = 1
  ids_profiles    = [%s]
}`, testAccNsxtPolicyMultitenancyContext(), name, action, direction, profilePath)
	}
	return testAccNsxtPolicyIntrusionServicePolicyRuleParentDeps(withContext) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_policy_rule" "rule1" {
  display_name    = "%s"
  policy_path     = nsxt_policy_parent_intrusion_service_policy.parent.path
  action          = "%s"
  direction       = "%s"
  sequence_number = 1
  ids_profiles    = [%s]
}`, name, action, direction, profilePath)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleWithAllFieldsCreate(withContext bool, name string) string {
	context := ""
	profilePath := fmt.Sprintf("\"%s\"", policyDefaultIdsProfilePath)
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
		profilePath = "nsxt_policy_intrusion_service_profile.test.path"
	}
	profile := ""
	if withContext {
		profile = testAccNsxtPolicyIntrusionServiceProfileMinimalistic(name, withContext)
	}
	return profile + testAccNsxtPolicyIntrusionServicePolicyRuleAllFieldsDeps(withContext) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_policy_rule" "rule1" {
%s
  display_name          = "%s"
  description           = "Standalone rule with all fields"
  policy_path           = nsxt_policy_parent_intrusion_service_policy.parent.path
  action                = "DETECT"
  direction             = "IN"
  ip_version            = "IPV4"
  logged                = true
  log_label             = "intrusion-svc-rule-detect"
  disabled              = false
  sources_excluded      = true
  destinations_excluded = true
  sequence_number       = 1
  source_groups         = [nsxt_policy_group.src_group.path]
  destination_groups    = [nsxt_policy_group.dst_group1.path]
  services              = [nsxt_policy_service.icmp_svc.path, nsxt_policy_service.tcp_svc.path]
  ids_profiles          = [%s]

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }
}`, context, name, profilePath)
}

func testAccNsxtPolicyIntrusionServicePolicyRuleWithAllFieldsUpdate(withContext bool, name string) string {
	context := ""
	profilePath := fmt.Sprintf("\"%s\"", policyDefaultIdsProfilePath)
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
		profilePath = "nsxt_policy_intrusion_service_profile.test.path"
	}
	profile := ""
	if withContext {
		profile = testAccNsxtPolicyIntrusionServiceProfileMinimalistic(name, withContext)
	}
	return profile + testAccNsxtPolicyIntrusionServicePolicyRuleAllFieldsDeps(withContext) + fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_policy_rule" "rule1" {
%s
  display_name          = "%s"
  description           = "Updated standalone rule"
  policy_path           = nsxt_policy_parent_intrusion_service_policy.parent.path
  action                = "DETECT_PREVENT"
  direction             = "IN_OUT"
  ip_version            = "IPV4_IPV6"
  logged                = true
  log_label             = "intrusion-svc-rule-prevent"
  disabled              = true
  sources_excluded      = false
  destinations_excluded = false
  sequence_number       = 2
  destination_groups    = [nsxt_policy_group.dst_group1.path, nsxt_policy_group.dst_group2.path]
  ids_profiles          = [%s]

  tag {
    scope = "env"
    tag   = "acceptance-test-updated"
  }
}`, context, name, profilePath)
}
