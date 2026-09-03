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

func TestAccResourceNsxtPolicyParentSecurityPolicy_basic(t *testing.T) {
	testAccResourceNsxtPolicyParentSecurityPolicyBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func TestAccResourceNsxtPolicyParentSecurityPolicy_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyParentSecurityPolicyBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyParentSecurityPolicyBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_parent_security_policy.test"

	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	locked := "true"
	updatedLocked := "false"
	seqNum := "1"
	updatedSeqNum := "2"
	tcpStrict := "true"
	updatedTCPStrict := "false"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentSecurityPolicyCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentSecurityPolicyTemplate(withContext, name, locked, seqNum, tcpStrict),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "locked", locked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", seqNum),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", tcpStrict),
				),
			},
			{
				Config: testAccNsxtPolicyParentSecurityPolicyTemplate(withContext, updatedName, updatedLocked, updatedSeqNum, updatedTCPStrict),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "locked", updatedLocked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", updatedSeqNum),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", updatedTCPStrict),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyParentSecurityPolicy_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_security_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentSecurityPolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentSecurityPolicyTemplate(false, name, "true", "1", "true"),
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

func TestAccResourceNsxtPolicyParentSecurityPolicy_importBasic_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_security_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentSecurityPolicyTemplate(true, name, "true", "1", "true"),
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

// TestAccResourceNsxtPolicyParentSecurityPolicy_predefinedSiblingNoDrift
// reproduces Bugzilla 3761445: a nsxt_policy_predefined_security_policy
// resource pointed at the same path as this parent resource can PATCH
// description/tag on the shared underlying object. Since the parent
// resource's own config never sets those fields, its next refresh must not
// show drift wanting to null them back out.
func TestAccResourceNsxtPolicyParentSecurityPolicy_predefinedSiblingNoDrift(t *testing.T) {
	name := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentSecurityPolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentSecurityPolicyPredefinedSiblingTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nsxt_policy_parent_security_policy.parent", "description", ""),
					resource.TestCheckResourceAttr("nsxt_policy_parent_security_policy.parent", "tag.#", "0"),
					resource.TestCheckResourceAttr("nsxt_policy_predefined_security_policy.predef", "description", "predefined description"),
				),
			},
		},
	})
}

func testAccNsxtPolicyParentSecurityPolicyPredefinedSiblingTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_security_policy" "parent" {
  display_name    = "%s"
  domain          = "default"
  category        = "Application"
  sequence_number = 10
}

resource "nsxt_policy_predefined_security_policy" "predef" {
  path        = nsxt_policy_parent_security_policy.parent.path
  description = "predefined description"

  tag {
    scope = "color"
    tag   = "orange"
  }
}`, name)
}

func testAccNsxtPolicyParentSecurityPolicyCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_parent_security_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]
		exists, err := resourceNsxtPolicySecurityPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domain, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy SecurityPolicy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyParentSecurityPolicyTemplate(withContext bool, name, locked, seqNum, tcpStrict string) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_parent_security_policy" "test" {
%s
  display_name    = "%s"
  description     = "Acceptance Test"
  domain          = "default"
  category        = "Application"
  locked          = %s
  sequence_number = %s
  stateful        = true
  tcp_strict      = %s

  tag {
    scope = "color"
    tag   = "orange"
  }
}`, context, name, locked, seqNum, tcpStrict)
}
