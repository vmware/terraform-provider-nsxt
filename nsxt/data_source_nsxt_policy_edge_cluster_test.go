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
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
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
		return testAccNsxtGlobalPolicyEdgeClusterReadTemplate()
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtGlobalPolicyEdgeClusterReadTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_edge_cluster" "test" {
  site_path = data.nsxt_policy_site.test.path
}`, getTestSiteName())
}

func testAccNsxtPolicyNoEdgeClusterTemplate() string {
	return fmt.Sprintf(` `)
}
