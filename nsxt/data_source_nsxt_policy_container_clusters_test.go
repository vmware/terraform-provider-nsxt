// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyContainerClusters_basic(t *testing.T) {
	containerClusterName := os.Getenv("NSXT_TEST_CONTAINER_CLUSTER")
	testResourceName := "data.nsxt_policy_container_clusters.test"
	checkResourceName := "data.nsxt_policy_container_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_CONTAINER_CLUSTER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyContainerClustersReadTemplate(containerClusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(checkResourceName, "display_name", containerClusterName),
				),
			},
		},
	})
}

func testAccNSXPolicyContainerClustersReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_container_clusters" "test" {
}

locals {
  // Get id from path
  path_split = split("/", data.nsxt_policy_container_clusters.test.items["%s"])
  container_cluster_id   = element(local.path_split, length(local.path_split) - 1)
}

data "nsxt_policy_container_cluster" "test" {
  id = local.container_cluster_id
}
`, name)
}
