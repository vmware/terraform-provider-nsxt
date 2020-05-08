/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func testAccDataSourceNsxtPolicySegmentRealization(t *testing.T, vlan bool) {
	testResourceName := "data.nsxt_policy_segment_realization.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySegmentRealizationTemplate(vlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "success"),
					resource.TestCheckResourceAttr(testResourceName, "network_name", "terra-test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicySegmentRealization_basic(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentRealization(t, false)
}

func TestAccDataSourceNsxtPolicySegmentRealization_vlan(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentRealization(t, true)
}

func testAccNsxtPolicySegmentRealizationTemplate(vlan bool) string {
	resource := "nsxt_policy_segment"
	tz := getOverlayTransportZoneName()
	extra := ""
	if vlan {
		resource = "nsxt_policy_vlan_segment"
		tz = getVlanTransportZoneName()
		extra = "vlan_ids = [12]"
	}
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "%s" "test" {
  display_name        = "terra-test"
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  %s
}

data "nsxt_policy_segment_realization" "test" {
  path = %s.test.path
}`, tz, resource, extra, resource)
}
