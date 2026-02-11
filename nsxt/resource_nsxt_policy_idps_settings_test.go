// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security"
)

func TestAccResourceNsxtPolicyIdpsSettings_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSettingsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSettingsExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "auto_update_signatures", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_syslog", "false"),
					resource.TestCheckResourceAttr(testResourceName, "oversubscription", "BYPASSED"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsSettings_update(t *testing.T) {
	testResourceName := "nsxt_policy_idps_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSettingsBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSettingsExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "auto_update_signatures", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_syslog", "false"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsSettingsUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSettingsExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "auto_update_signatures", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_syslog", "true"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsSettings_oversubscription(t *testing.T) {
	testResourceName := "nsxt_policy_idps_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSettingsOversubscription("BYPASSED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSettingsExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "oversubscription", "BYPASSED"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsSettingsOversubscription("DROPPED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsSettingsExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "oversubscription", "DROPPED"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("IDPS Settings resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("IDPS Settings resource ID not set in resources")
		}

		client := security.NewIntrusionServicesClient(connector)
		_, err := client.Get()
		if err != nil {
			return fmt.Errorf("Error while retrieving IDPS Settings: %v", err)
		}

		return nil
	}
}

func testAccNsxtPolicyIdpsSettingsBasic() string {
	return `
resource "nsxt_policy_idps_settings" "test" {
  description            = "Acceptance test basic"
  auto_update_signatures = false
  enable_syslog          = false
  oversubscription       = "BYPASSED"
}`
}

func testAccNsxtPolicyIdpsSettingsUpdate() string {
	return `
resource "nsxt_policy_idps_settings" "test" {
  description            = "Acceptance test updated"
  auto_update_signatures = true
  enable_syslog          = true
  oversubscription       = "BYPASSED"
}`
}

func testAccNsxtPolicyIdpsSettingsOversubscription(action string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_idps_settings" "test" {
  description      = "Testing oversubscription"
  oversubscription = "%s"
}`, action)
}
