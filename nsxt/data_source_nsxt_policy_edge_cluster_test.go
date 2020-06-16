/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyEdgeCluster_basic(t *testing.T) {
	edgeClusterName := getEdgeClusterName()
	testResourceName := "data.nsxt_policy_edge_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXGlobalManagerSitePrecheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeClusterReadTemplate(edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", edgeClusterName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoEdgeClusterTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyEdgeClusterReadTemplate(name string) string {
	if testAccIsGlobalManager() {
		return testAccNsxtGlobalPolicyEdgeClusterReadTemplate(name)
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtGlobalPolicyEdgeClusterReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
  site_path = data.nsxt_policy_site.test.path
}`, getTestSiteName(), name)
}

func testAccNsxtPolicyNoEdgeClusterTemplate() string {
	return fmt.Sprintf(` `)
}
