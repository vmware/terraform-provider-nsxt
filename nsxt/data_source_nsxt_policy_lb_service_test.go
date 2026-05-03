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

func TestAccDataSourceNsxtPolicyLBService(t *testing.T) {
	lbServiceName := getAccTestResourceName()
	tier1Name := getAccTestResourceName()
	tier0Name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_lb_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBServiceCheckDestroy(state, lbServiceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBServiceDataSourceTemplate(lbServiceName, tier1Name, tier0Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", lbServiceName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyLBServiceDataSourceTemplate(lbServiceName, tier1Name, tier0Name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "lb_svc_t1" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  tier0_path        = nsxt_policy_tier0_gateway.test.path
}

resource "nsxt_policy_lb_service" "created" {
  display_name        = "%s"
  connectivity_path   = nsxt_policy_tier1_gateway.lb_svc_t1.path
}

data "nsxt_policy_lb_service" "test" {
  display_name = "%s"

  depends_on = [nsxt_policy_lb_service.created]
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_service.created.path
}
`, getEdgeClusterName(), tier0Name, tier1Name, lbServiceName, lbServiceName)
}
