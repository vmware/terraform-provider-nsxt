// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// site_path is required on Global Manager, since the realized entity lookup
// is site-scoped there. Without it, GM's realized-entities API silently
// returns no results and the data source times out.
func TestAccDataSourceNsxtPolicyGatewayInterfaceRealization_sitePathRequiredOnGM(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyGlobalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyGatewayInterfaceRealizationMissingSitePathTemplate(),
				ExpectError: regexp.MustCompile("requires site_path configuration for NSX Global Manager"),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGatewayInterfaceRealization_sitePathNotAllowedOnLocalManager(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyGatewayInterfaceRealizationExtraSitePathTemplate(),
				ExpectError: regexp.MustCompile("only supported with NSX Global Manager"),
			},
		},
	})
}

func testAccNsxtPolicyGatewayInterfaceRealizationMissingSitePathTemplate() string {
	return `
data "nsxt_policy_gateway_interface_realization" "test" {
  gateway_interface_path = "/global-infra/tier-0s/dummy/locale-services/dummy/interfaces/dummy"
}`
}

func testAccNsxtPolicyGatewayInterfaceRealizationExtraSitePathTemplate() string {
	return `
data "nsxt_policy_gateway_interface_realization" "test" {
  gateway_interface_path = "/infra/tier-0s/dummy/locale-services/dummy/interfaces/dummy"
  site_path              = "/infra/sites/default"
}`
}
