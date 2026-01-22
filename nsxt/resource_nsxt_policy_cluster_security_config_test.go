// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/security"
)

func TestAccResourceNsxtPolicyClusterSecurityConfig_basic(t *testing.T) {
	testResourceName := "nsxt_policy_cluster_security_config.test"
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyClusterSecurityConfigCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyClusterSecurityConfigCreateTemplate(computeCollectionName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyClusterSecurityConfigExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "dfw_enabled", "true"),
				),
			},
			{
				Config: testAccNsxtPolicyClusterSecurityConfigUpdateTemplate(computeCollectionName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyClusterSecurityConfigExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster_id"),
					resource.TestCheckResourceAttr(testResourceName, "dfw_enabled", "false"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyClusterSecurityConfig_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_cluster_security_config.test"
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyClusterSecurityConfigCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyClusterSecurityConfigCreateTemplate(computeCollectionName),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revision"},
			},
		},
	})
}

func testAccNsxtPolicyClusterSecurityConfigExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Cluster Security Config resource %s not found in resources", resourceName)
		}

		clusterID := rs.Primary.ID
		if clusterID == "" {
			return fmt.Errorf("Cluster Security Config resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta())
		client := security.NewClusterConfigsClient(connector)
		_, err := client.Get(clusterID, nil)
		if err != nil {
			return fmt.Errorf("Error while retrieving cluster security config ID %s: %v", clusterID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyClusterSecurityConfigCheckDestroy(state *terraform.State) error {
	// Cluster security config cannot be deleted, so we don't check for deletion
	// We just verify that all features are disabled
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_cluster_security_config" {
			continue
		}

		clusterID := rs.Primary.ID
		connector := getPolicyConnector(testAccProvider.Meta())
		client := security.NewClusterConfigsClient(connector)
		config, err := client.Get(clusterID, nil)
		if err != nil {
			// If we can't read it, that's fine for destroy check
			continue
		}

		// Check if all features are disabled
		if config.Features != nil {
			for _, feature := range config.Features {
				if feature.Enabled != nil && *feature.Enabled {
					return fmt.Errorf("Cluster security config %s still has enabled features after destroy", clusterID)
				}
			}
		}
	}

	return nil
}

func testAccNsxtPolicyClusterSecurityConfigCreateTemplate(computeCollectionName string) string {
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

resource "nsxt_policy_cluster_security_config" "test" {
  cluster_id  = data.nsxt_compute_collection.test.id
  dfw_enabled = true
}
`, computeCollectionName)
}

func testAccNsxtPolicyClusterSecurityConfigUpdateTemplate(computeCollectionName string) string {
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

resource "nsxt_policy_cluster_security_config" "test" {
  cluster_id  = data.nsxt_compute_collection.test.id
  dfw_enabled = false
}
`, computeCollectionName)
}
