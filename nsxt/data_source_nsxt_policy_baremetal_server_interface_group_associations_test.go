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

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyBareMetalServerInterfaceGroupAssociationsTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
					testAccBMSInterfaceGroupAssociationsConditionalCheck(),
				),
			},
		},
	})
}

func testAccBMSInterfaceGroupAssociationsConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get the discovery data source
		rs, ok := s.RootModule().Resources["data.nsxt_policy_baremetal_server_interfaces.all"]
		if !ok {
			return fmt.Errorf("Not found: data.nsxt_policy_baremetal_server_interfaces.all")
		}

		// Check if we have interfaces
		resultsAttr := rs.Primary.Attributes["results.#"]
		if resultsAttr == "0" {
			// No interfaces found - skip group associations test
			return nil
		}

		// We have interfaces, check that the group associations data source exists
		groupAssocRs, ok := s.RootModule().Resources["data.nsxt_policy_baremetal_server_interface_group_associations.test.0"]
		if !ok {
			return fmt.Errorf("BMS interfaces exist but group associations data source not found")
		}

		if groupAssocRs.Primary.Attributes["external_id"] == "" {
			return fmt.Errorf("external_id not set for group associations")
		}

		return nil
	}
}

func testAccNSXPolicyBareMetalServerInterfaceGroupAssociationsTemplate() string {
	return `
# Discovery - find available BMS interfaces
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  first_interface_id = local.has_interfaces ? sort([for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id])[0] : ""
}

# Test group associations for first available interface
data "nsxt_policy_baremetal_server_interface_group_associations" "test" {
  count       = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
}`
}
