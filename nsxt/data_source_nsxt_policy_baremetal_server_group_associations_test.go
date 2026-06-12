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

func TestAccDataSourceNsxtPolicyBareMetalServerGroupAssociations_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyBaremetalServerGroupAssociationsTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					testAccBMSServerGroupAssociationsConditionalCheck(),
				),
			},
		},
	})
}

func testAccBMSServerGroupAssociationsConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get the discovery data source
		rs, ok := s.RootModule().Resources["data.nsxt_policy_baremetal_servers.all"]
		if !ok {
			return fmt.Errorf("Not found: data.nsxt_policy_baremetal_servers.all")
		}

		// Check if we have servers
		resultsAttr := rs.Primary.Attributes["results.#"]
		if resultsAttr == "0" {
			// No servers found - skip group associations test
			return nil
		}

		// We have servers, check that the group associations data source exists
		groupAssocRs, ok := s.RootModule().Resources["data.nsxt_policy_baremetal_server_group_associations.test.0"]
		if !ok {
			return fmt.Errorf("BMS servers exist but group associations data source not found")
		}

		if groupAssocRs.Primary.Attributes["external_id"] == "" {
			return fmt.Errorf("external_id not set for group associations")
		}

		return nil
	}
}

func testAccNsxtPolicyBaremetalServerGroupAssociationsTemplate() string {
	return `
# Discovery - find available BMS servers
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  first_server_id = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id])[0] : ""
}

# Test group associations for first available server
data "nsxt_policy_baremetal_server_group_associations" "test" {
  count       = local.has_servers ? 1 : 0
  external_id = local.first_server_id
}`
}
