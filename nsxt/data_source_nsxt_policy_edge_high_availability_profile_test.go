// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyEdgeHighAvailabilityProfile_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_edge_high_availability_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXEdgeHighAvailabilityProfileReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccNSXEdgeHighAvailabilityProfileReadTemplate() string {
	return testAccNsxtPolicyEdgeHighAvailabilityProfileMinimalistic() + `
data "nsxt_policy_edge_high_availability_profile" "test" {
  display_name = nsxt_policy_edge_high_availability_profile.test.display_name 

  depends_on = [nsxt_policy_edge_high_availability_profile.test]
}`
}
