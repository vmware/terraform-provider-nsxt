// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIdpsClusterConfig_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyIdpsClusterConfigBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "4.2.0")
		testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
	})
}

func testAccDataSourceNsxtPolicyIdpsClusterConfigBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_idps_cluster_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSource(name) + testAccNsxtPolicyIdpsClusterConfigReadTemplate("nsxt_policy_idps_cluster_config.test.display_name", withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "ids_enabled"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_type"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSource(name string) string {
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

resource "nsxt_policy_idps_cluster_config" "test" {
  display_name = "%s"
  ids_enabled  = true

  cluster {
    target_id   = data.nsxt_compute_collection.test.id
    target_type = "VC_Cluster"
  }
}`, getTestComputeCollectionName(), name)
}

func testAccNsxtPolicyIdpsClusterConfigReadTemplate(name string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "nsxt_policy_idps_cluster_config" "test" {
%s
  display_name = %s
}`, context, name)
}
