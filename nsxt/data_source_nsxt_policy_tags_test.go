// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The nsxt_policy_tags data source reads from a backend tag aggregation
// index that is synced asynchronously from realized objects (e.g. segments).
// Reading it in the same apply as the object that owns the tag is racy, since
// the sync may not have completed yet. Creating the segment in its own step
// with a settle delay before reading the tags avoids that race.
const testAccNSXPolicyTagsSyncDelay = 10 * time.Second

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
				Config: testAccNSXPolicyTagsCreateTemplate(tagName, emptyScopeTag, transportZone),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_segment.segment1", "path"),
					func(_ *terraform.State) error {
						time.Sleep(testAccNSXPolicyTagsSyncDelay)
						return nil
					},
				),
			},
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

func testAccNSXPolicyTagsCreateTemplate(tagName string, emptyScopeTag string, transportZone string) string {
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
    tag = "%s"
  }

}

data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name = "%s"
}
`, tagName, emptyScopeTag, transportZone)
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
    tag = "%s"
  }

}

data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name = "%s"
}

data "nsxt_policy_tags" "tags" {
  scope      = "scope-test"
  depends_on = [nsxt_policy_segment.segment1]
}

output "nsxt_tags" {
  value = data.nsxt_policy_tags.tags.items[0]
}

data "nsxt_policy_tags" "emptytags" {
  scope      = ""
  depends_on = [nsxt_policy_segment.segment1]
}

output "empty_nsxt_tags" {
  value = join("--", data.nsxt_policy_tags.emptytags.items)
}

data "nsxt_policy_tags" "wildcardscope" {
  scope      = "*test"
  depends_on = [nsxt_policy_segment.segment1]
}

output "wildcard_nsxt_tags" {
  value = data.nsxt_policy_tags.wildcardscope.items[0]
}


`, tagName, emptyScopeTag, transportZone)
}
