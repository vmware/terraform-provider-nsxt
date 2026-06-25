// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyVirtualNetworkApplianceCluster_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_virtual_network_appliance_cluster.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.1")
			testAccEnvDefined(t, "NSXT_TEST_OVERLAY_TRANSPORT_ZONE")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "appliance_form_factor", "MEDIUM"),
					resource.TestCheckResourceAttr(testResourceName, "service_type", "VPC_SERVICES"),
				),
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterDataSourceTemplate(displayName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s"
  description           = "Acceptance test cluster for data source"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.test.path
  }
}

data "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name = nsxt_policy_virtual_network_appliance_cluster.test.display_name
}
`, getOverlayTransportZoneName(), displayName)
}
