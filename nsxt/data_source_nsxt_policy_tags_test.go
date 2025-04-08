// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyTags_basic(t *testing.T) {
	tagName := "testTag"
	emptyScopeTag := "testEmptyScopeTag"
	transportZone := getOverlayTransportZoneName()
	re, _ := regexp.Compile(`.*EmptyScope.*`)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTagsReadTemplate(tagName, emptyScopeTag, transportZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("nsxt_tags", tagName),
					resource.TestMatchOutput("empty_nsxt_tags", re),
					resource.TestCheckOutput("wildcard_nsxt_tags", tagName),
				),
			},
		},
	})
}

func testAccNSXPolicyTagsReadTemplate(tagName string, emptyScopeTag string, transportZone string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_segment" "segment1" {
  display_name        = "segment1"
  description         = "Terraform provisioned Segment"
  transport_zone_path = data.nsxt_policy_transport_zone.overlay_transport_zone.path
  tag {
    scope = "scope-test"
    tag   = "%s"
  }
  tag {
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

data "nsxt_policy_tags" "emptytags" {
  scope = ""
  depends_on = [nsxt_policy_segment.segment1]
}

output "empty_nsxt_tags" {
  value = join("--",data.nsxt_policy_tags.emptytags.items)
}

data "nsxt_policy_tags" "wildcardscope" {
  scope = "*test"
  depends_on = [nsxt_policy_segment.segment1]
}

output "wildcard_nsxt_tags" {
  value = data.nsxt_policy_tags.wildcardscope.items[0]
}

`, tagName, emptyScopeTag, transportZone)
}
