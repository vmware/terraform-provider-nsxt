// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyEdgeNode_basic(t *testing.T) {
	edgeClusterName := getEdgeClusterName()
	testResourceName := "data.nsxt_policy_edge_node.test"

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
		resource.TestCheckResourceAttr(testResourceName, "member_index", "0"),
		resource.TestCheckResourceAttrSet(testResourceName, "path"),
		resource.TestCheckResourceAttrSet(testResourceName, "unique_id"),
	}
	if !testAccIsGlobalManager() {
		// Global Manager does not populate realization_id for PolicyEdgeNode
		checks = append(checks, resource.TestCheckResourceAttrSet(testResourceName, "realization_id"))
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXGlobalManagerSitePrecheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeNodeReadTemplate(edgeClusterName),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccNsxtPolicyEdgeNodeReadTemplate(name string) string {
	return testAccNsxtPolicyEdgeClusterReadTemplate(name) + `

data "nsxt_policy_edge_node" "test" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  member_index      = 0
}`
}
