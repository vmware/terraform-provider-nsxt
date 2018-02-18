/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtFirewallSection_basic(t *testing.T) {

	prfName := fmt.Sprintf("test-nsx-firewall-section-basic")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_firewall_section.test"
	tags := singleTag
	updatedTags := doubleTags
	tos := string("[]")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXFirewallSectionCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXFirewallSectionCreateEmptyTemplate(prfName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
				),
			},
			{
				Config: testAccNSXFirewallSectionUpdateEmptyTemplate(updatePrfName, updatedTags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtFirewallSection_withTos(t *testing.T) {
	prfName := fmt.Sprintf("test-nsx-firewall-section-tos")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_firewall_section.test"
	tags := singleTag
	tos := string("[{target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP1.id}\"}]")
	updatedTos := string("[{target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP1.id}\"}, {target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP2.id}\"}]")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXFirewallSectionCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXFirewallSectionCreateEmptyTemplate(prfName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "1"),
				),
			},
			{
				Config: testAccNSXFirewallSectionUpdateEmptyTemplate(updatePrfName, tags, updatedTos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtFirewallSection_withRules(t *testing.T) {

	prfName := fmt.Sprintf("test-nsx-firewall-section-rules")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_firewall_section.test"
	ruleName := "rule1.0"
	updatedRuleName := "rule1.1"
	tags := singleTag
	tos := string("[]")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXFirewallSectionCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXFirewallSectionCreateTemplate(prfName, ruleName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", ruleName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
				),
			},
			{
				Config: testAccNSXFirewallSectionUpdateTemplate(updatePrfName, updatedRuleName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", updatedRuleName),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func skipFirewallSectionTests() (bool, error) {
	// TODO(asarfaty): this currently doesn't work since the provider configure func was not called yet
	//testAccProvider.ConfigureFunc([]*ResourceData{})
	//nsx_ver := getNSXVersion(testAccProvider.Meta())
	//if nsx_ver <= "2.1.0" {
	//	fmt.Printf("Skipping test because it is not supported on nsx version %s", nsx_ver)
	//	return true, nil
	//}
	return false, nil
}

func TestAccResourceNsxtFirewallSection_withRulesAndTags(t *testing.T) {

	// Note: this test will not pass with NSX 2.1 because of an NSX bug.
	// This test should be skipped based on the nsx version
	prfName := fmt.Sprintf("test-nsx-firewall-section-tags")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_firewall_section.test"
	ruleName := "rule1.0"
	updatedRuleName := "rule1.1"
	tags := singleTag
	updatedTags := doubleTags
	tos := string("[]")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXFirewallSectionCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				SkipFunc: skipFirewallSectionTests,
				Config:   testAccNSXFirewallSectionCreateTemplate(prfName, ruleName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", ruleName),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				SkipFunc: skipFirewallSectionTests,
				Config:   testAccNSXFirewallSectionUpdateTemplate(updatePrfName, updatedRuleName, updatedTags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", updatedRuleName),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtFirewallSection_withRulesAndTos(t *testing.T) {
	prfName := fmt.Sprintf("test-nsx-firewall-section-rules_and_tos")
	updatePrfName := fmt.Sprintf("%s-update", prfName)
	testResourceName := "nsxt_firewall_section.test"
	ruleName := "rule1.0"
	updatedRuleName := "rule1.1"
	tags := singleTag
	tos := string("[{target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP1.id}\"}]")
	updatedTos := string("[{target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP1.id}\"}, {target_type = \"NSGroup\", target_id = \"${nsxt_ns_group.GRP2.id}\"}]")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXFirewallSectionCheckDestroy(state, prfName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXFirewallSectionCreateTemplate(prfName, ruleName, tags, tos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(prfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", prfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", ruleName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "1"),
				),
			},
			{
				Config: testAccNSXFirewallSectionUpdateTemplate(updatePrfName, updatedRuleName, tags, updatedTos),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXFirewallSectionExists(updatePrfName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePrfName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "section_type", "LAYER3"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", updatedRuleName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "applied_to.#", "2"),
				),
			},
		},
	})
}

func testAccNSXFirewallSectionExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Firewall Section resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Firewall Section resource ID not set in resources ")
		}

		profile, responseCode, err := nsxClient.ServicesApi.GetSection(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving firewall section ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if firewall section %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == profile.DisplayName {
			return nil
		}
		return fmt.Errorf("Firewall Section %s wasn't found", display_name)
	}
}

func testAccNSXFirewallSectionCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_port" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		profile, responseCode, err := nsxClient.ServicesApi.GetSection(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving firewall section ID %s. Error: %v", resourceID, err)
		}

		if display_name == profile.DisplayName {
			return fmt.Errorf("Firewall Section %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXFirewallSectionNSGroups() string {
	return fmt.Sprintf(`
resource "nsxt_ns_group" "GRP1" {
    display_name = "grp1"
}
resource "nsxt_ns_group" "GRP2" {
    display_name = "grp2"
}`)
}

func testAccNSXFirewallSectionCreateTemplate(name string, ruleName string, tags string, tos string) string {
	return testAccNSXFirewallSectionNSGroups() + fmt.Sprintf(`
resource "nsxt_firewall_section" "test" {
	display_name = "%s"
	description = "Acceptance Test"
    section_type = "LAYER3"
    stateful = true
	tag = %s
	rule {
		display_name = "%s",
	    description = "rule1",
    	action = "ALLOW",
     	logged = "true",
    	ip_protocol = "IPV4",
    	direction = "IN"
    }
	applied_to = %s 
}`, name, tags, ruleName, tos)
}

func testAccNSXFirewallSectionUpdateTemplate(updatedName string, updatedRuleName string, tags string, tos string) string {
	return testAccNSXFirewallSectionNSGroups() + fmt.Sprintf(`
resource "nsxt_firewall_section" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
    section_type = "LAYER3"
    stateful = true
	tag = %s
	rule {
		display_name = "%s",
	    description = "rule1",
    	action = "ALLOW",
      	logged = "true",
      	ip_protocol = "IPV4",
      	direction = "IN"
    }
	rule {
		display_name = "rule2",
      	description = "rule2",
      	action = "ALLOW",
      	logged = "true",
      	ip_protocol = "IPV6",
      	direction = "OUT"
    }
	applied_to = %s 
}`, updatedName, tags, updatedRuleName, tos)
}

func testAccNSXFirewallSectionCreateEmptyTemplate(name string, tags string, tos string) string {
	return testAccNSXFirewallSectionNSGroups() + fmt.Sprintf(`
resource "nsxt_firewall_section" "test" {
	display_name = "%s"
	description = "Acceptance Test"
    section_type = "LAYER3"
    stateful = true
	tag = %s
	applied_to = %s
}`, name, tags, tos)
}

func testAccNSXFirewallSectionUpdateEmptyTemplate(updatedName string, tags string, tos string) string {
	return testAccNSXFirewallSectionNSGroups() + fmt.Sprintf(`
resource "nsxt_firewall_section" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
    section_type = "LAYER3"
    stateful = true
	tag = %s
	applied_to = %s
}`, updatedName, tags, tos)
}
