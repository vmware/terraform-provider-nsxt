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
		testAccNSXVersion(t, "9.1.0")
		testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
	})
}

func TestAccDataSourceNsxtPolicyIdpsClusterConfig_basicLegacy(t *testing.T) {
	testAccDataSourceNsxtPolicyIdpsClusterConfigBasicLegacy(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "4.2.0")
		testAccNSXVersionLessThan(t, "9.1.0")
		testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
	})
}

func testAccDataSourceNsxtPolicyIdpsClusterConfigBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()
	testResourceName := "data.nsxt_policy_idps_cluster_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSourceWithDFW(computeCollectionName, name) + testAccNsxtPolicyIdpsClusterConfigReadTemplate("nsxt_policy_idps_cluster_config.test.display_name", withContext),
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

func testAccDataSourceNsxtPolicyIdpsClusterConfigBasicLegacy(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()
	testResourceName := "data.nsxt_policy_idps_cluster_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSource(computeCollectionName, name) + testAccNsxtPolicyIdpsClusterConfigReadTemplate("nsxt_policy_idps_cluster_config.test.display_name", withContext),
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

func testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSource(computeCollection, name string) string {
	// For NSX < 9.1.0: DFW is not a prerequisite for IDPS (different architecture)
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
}`, computeCollection, name)
}

func testAccNsxtPolicyIdpsClusterConfigMinimalisticForDataSourceWithDFW(computeCollection, name string) string {
	// For NSX 9.1.0+: DFW must be enabled before IDPS (enforced by NSX)
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

# Enable DFW on the cluster (required for IDPS on NSX 9.1.0+)
resource "nsxt_policy_cluster_security_config" "test" {
  cluster_id  = data.nsxt_compute_collection.test.id
  dfw_enabled = true
}

resource "nsxt_policy_idps_cluster_config" "test" {
  depends_on = [nsxt_policy_cluster_security_config.test]
  
  display_name = "%s"
  ids_enabled  = true

  cluster {
    target_id   = data.nsxt_compute_collection.test.id
    target_type = "VC_Cluster"
  }
}`, computeCollection, name)
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
