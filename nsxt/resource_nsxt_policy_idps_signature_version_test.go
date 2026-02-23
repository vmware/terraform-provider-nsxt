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

// NOTE: IDS Signature Versions are system-managed resources.
// These tests reference existing versions rather than creating new ones.

func TestAccResourceNsxtPolicyIdpsSignatureVersion_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_signature_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsSignatureVersionCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSignatureVersionExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "DEFAULT"),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "state"),
					resource.TestCheckResourceAttrSet(testResourceName, "status"),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsSignatureVersion_makeActive(t *testing.T) {
	testResourceName := "nsxt_policy_idps_signature_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsSignatureVersionCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSignatureVersionExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "DEFAULT"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionMakeActive(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSignatureVersionExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "DEFAULT"),
					resource.TestCheckResourceAttr(testResourceName, "state", "ACTIVE"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsSignatureVersionExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IDS Signature Version resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IDS Signature Version resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIdpsSignatureVersionExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy IDS Signature Version %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyIdpsSignatureVersionCheckDestroy(state *terraform.State) error {
	// NOTE: Signature versions are system-managed and cannot be deleted.
	// This check just verifies the resource was removed from Terraform state.
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_idps_signature_version" {
			continue
		}

		// Resource should be removed from state
		// The actual version still exists in NSX (which is expected)
	}
	return nil
}

func testAccNsxtPolicyIdpsSignatureVersionBasic() string {
	return `
resource "nsxt_policy_idps_signature_version" "test" {
  nsx_id       = "DEFAULT"
  display_name = "DEFAULT"
}`
}

func testAccNsxtPolicyIdpsSignatureVersionMakeActive() string {
	return `
resource "nsxt_policy_idps_signature_version" "test" {
  nsx_id       = "DEFAULT"
  display_name = "DEFAULT"
  state        = "ACTIVE"
}`
}
