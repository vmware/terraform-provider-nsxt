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

func TestAccResourceNsxtPolicyParentIntrusionServiceGatewayPolicy_basic(t *testing.T) {
	testResourceName := "nsxt_policy_parent_intrusion_service_gateway_policy.test"

	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	locked := "true"
	updatedLocked := "false"
	seqNum := "1"
	updatedSeqNum := "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentIntrusionServiceGatewayPolicyCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServiceGatewayPolicyTemplate(name, locked, seqNum),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyParentIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "locked", locked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", seqNum),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyParentIntrusionServiceGatewayPolicyTemplate(updatedName, updatedLocked, updatedSeqNum),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyParentIntrusionServiceGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "locked", updatedLocked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", updatedSeqNum),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyParentIntrusionServiceGatewayPolicy_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_intrusion_service_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentIntrusionServiceGatewayPolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServiceGatewayPolicyTemplate(name, "true", "1"),
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

func testAccNsxtPolicyParentIntrusionServiceGatewayPolicyExists(resourceName string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Parent Intrusion Service Gateway Policy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Parent Intrusion Service Gateway Policy resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domainName, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving Parent Intrusion Service Gateway Policy resource ID %s", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyParentIntrusionServiceGatewayPolicyCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_parent_intrusion_service_gateway_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]
		exists, err := resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domain, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Parent Intrusion Service Gateway Policy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyParentIntrusionServiceGatewayPolicyTemplate(name, locked, seqNum string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_gateway_policy" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  locked          = %s
  sequence_number = %s

  tag {
    scope = "env"
    tag   = "acceptance-test"
  }
}`, name, locked, seqNum)
}
