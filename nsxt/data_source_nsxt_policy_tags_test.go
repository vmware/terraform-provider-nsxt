// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyTags_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyTagsBasic(t, func() {
		testAccPreCheck(t)
	})
}

func testAccDataSourceNsxtPolicyTagsBasic(t *testing.T, preCheck func()) {
	tagName := "testTag"
	transportZone := getOverlayTransportZoneName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTagsReadTemplate(tagName, transportZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("nsxt_tags", tagName),
				),
			},
		},
	})
}

func testAccNSXPolicyTagsReadTemplate(tagName string, transportZone string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_segment" "segment1" {
  display_name        = "segment1"
  description         = "Terraform provisioned Segment"
  transport_zone_path = data.nsxt_policy_transport_zone.overlay_transport_zone.path
  tag {
    scope = "scope-test"
    tag   = "%s"
  }

}

data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name = "%s"
}

data "nsxt_policy_tags" "tags" {
  scope = "scope-test"
  depends_on = [nsxt_policy_segment.segment1]
}

output "nsxt_tags" {
  value = data.nsxt_policy_tags.tags.items[0]
}

`, tagName, transportZone)
}
