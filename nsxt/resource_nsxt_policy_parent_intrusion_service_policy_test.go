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

func TestAccResourceNsxtPolicyParentIntrusionServicePolicy_basic(t *testing.T) {
	testAccResourceNsxtPolicyParentIntrusionServicePolicyBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "4.2.0")
	})
}

func TestAccResourceNsxtPolicyParentIntrusionServicePolicy_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyParentIntrusionServicePolicyBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyParentIntrusionServicePolicyBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_parent_intrusion_service_policy.test"

	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	locked := "true"
	updatedLocked := "false"
	seqNum := "1"
	updatedSeqNum := "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentIntrusionServicePolicyCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyTemplate(withContext, name, locked, seqNum),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "locked", locked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", seqNum),
				),
			},
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyTemplate(withContext, updatedName, updatedLocked, updatedSeqNum),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "locked", updatedLocked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", updatedSeqNum),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyParentIntrusionServicePolicy_importBasic(t *testing.T) {
	testAccResourceNsxtPolicyParentIntrusionServicePolicyImportBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "4.2.0")
	})
}

func TestAccResourceNsxtPolicyParentIntrusionServicePolicy_importBasic_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyParentIntrusionServicePolicyImportBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyParentIntrusionServicePolicyImportBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_intrusion_service_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentIntrusionServicePolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyTemplate(withContext, name, "true", "1"),
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

func testAccNsxtPolicyParentIntrusionServicePolicyCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_parent_intrusion_service_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]
		exists, err := resourceNsxtPolicyIntrusionServicePolicyExistsInDomain(testAccGetSessionContext(), resourceID, domain, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Intrusion Service Policy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyParentIntrusionServicePolicyTemplate(withContext bool, name, locked, seqNum string) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "test" {
%s
  display_name    = "%s"
  description     = "Acceptance Test"
  domain          = "default"
  locked          = %s
  sequence_number = %s

  tag {
    scope = "color"
    tag   = "orange"
  }
}`, context, name, locked, seqNum)
}
