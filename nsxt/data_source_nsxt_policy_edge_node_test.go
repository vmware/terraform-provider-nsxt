/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyEdgeNode_basic(t *testing.T) {
	edgeClusterName := getEdgeClusterName()
	testResourceName := "data.nsxt_policy_edge_node.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeNodeReadTemplate(edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "member_index", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoEdgeNodeTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyEdgeNodeReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
    display_name = "%s"
}
data "nsxt_policy_edge_node" "test" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  member_index       = 0
}`, name)
}

func testAccNsxtPolicyNoEdgeNodeTemplate() string {
	return fmt.Sprintf(` `)
}
