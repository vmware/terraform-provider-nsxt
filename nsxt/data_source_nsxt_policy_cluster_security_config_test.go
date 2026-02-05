// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyClusterSecurityConfig_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_cluster_security_config.test"
	resourceName := "nsxt_policy_cluster_security_config.test"
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyClusterSecurityConfigReadTemplate(computeCollectionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "dfw_enabled", "true"),
					resource.TestCheckResourceAttrPair(testResourceName, "cluster_id", resourceName, "cluster_id"),
					resource.TestCheckResourceAttrPair(testResourceName, "dfw_enabled", resourceName, "dfw_enabled"),
					resource.TestCheckResourceAttrPair(testResourceName, "display_name", resourceName, "display_name"),
					resource.TestCheckResourceAttrPair(testResourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(testResourceName, "path", resourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyClusterSecurityConfigReadTemplate(computeCollectionName string) string {
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

resource "nsxt_policy_cluster_security_config" "test" {
  cluster_id  = data.nsxt_compute_collection.test.id
  dfw_enabled = true
}

data "nsxt_policy_cluster_security_config" "test" {
  cluster_id = nsxt_policy_cluster_security_config.test.cluster_id
}
`, computeCollectionName)
}
