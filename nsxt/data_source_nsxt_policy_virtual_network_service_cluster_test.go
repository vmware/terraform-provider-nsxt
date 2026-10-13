// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyVirtualNetworkServiceCluster_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_virtual_network_service_cluster.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVNSCluster(t)
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_MANAGER")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkServiceClusterCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkServiceClusterDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceName, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "node_form_factor", "MEDIUM"),
					resource.TestCheckResourceAttr(testDataSourceName, "number_of_nodes", "2"),
					resource.TestCheckResourceAttr(testDataSourceName, "vsphere_deployment_config.#", "1"),
					resource.TestCheckResourceAttr("data.nsxt_policy_virtual_network_service_cluster.by_id", "display_name", displayName),
					resource.TestCheckResourceAttr("data.nsxt_policy_virtual_network_service_cluster.by_id", "node_form_factor", "MEDIUM"),
				),
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkServiceClusterDataSourceTemplate(displayName string) string {
	return testAccNsxtPolicyVirtualNetworkServiceClusterTemplate(displayName, "MEDIUM", 2, false) + `
data "nsxt_policy_virtual_network_service_cluster" "test" {
  display_name = nsxt_policy_virtual_network_service_cluster.test.display_name
}

data "nsxt_policy_virtual_network_service_cluster" "by_id" {
  id = nsxt_policy_virtual_network_service_cluster.test.id
}
`
}
