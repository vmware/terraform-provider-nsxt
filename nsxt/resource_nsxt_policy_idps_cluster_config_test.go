// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
)

func getTestComputeCollectionName() string {
	// Returns the compute collection display_name from environment variable (e.g., "cl1")
	// This will be used in the data source to lookup and fetch the external_id
	// The external_id format is: "uuid:domain-cXX" which is what IDPS cluster config needs
	return getComputeCollectionName()
}

func TestAccResourceNsxtPolicyIdpsClusterConfig_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_cluster_config.test"
	displayName := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplateWithDFW(computeCollectionName, displayName, "terraform created", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", "VC_Cluster"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplateWithDFW(computeCollectionName, displayName, "terraform updated", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform updated"),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", "VC_Cluster"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsClusterConfig_basicLegacy(t *testing.T) {
	testResourceName := "nsxt_policy_idps_cluster_config.test"
	displayName := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
			testAccNSXVersionLessThan(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplate(computeCollectionName, displayName, "terraform created", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", "VC_Cluster"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplate(computeCollectionName, displayName, "terraform updated", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform updated"),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", "VC_Cluster"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsClusterConfig_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_cluster_config.test"
	displayName := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalisticWithDFW(computeCollectionName, displayName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyIdpsClusterConfig_importBasicLegacy(t *testing.T) {
	testResourceName := "nsxt_policy_idps_cluster_config.test"
	displayName := getAccTestResourceName()
	computeCollectionName := getTestComputeCollectionName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
			testAccNSXVersionLessThan(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalistic(computeCollectionName, displayName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyIdpsClusterConfigExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IdpsClusterConfig resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IdpsClusterConfig resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIdpsClusterConfigExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy IdpsClusterConfig %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	client := intrusion_services.NewClusterConfigsClient(connector)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_idps_cluster_config" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]

		// IDPS Cluster Config cannot be deleted, only disabled
		// So we check if ids_enabled is false after "deletion"
		obj, err := client.Get(resourceID, nil)
		if err != nil {
			// If we get a not found error, that's unexpected but acceptable
			if isNotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error checking IDPS Cluster Config %s: %v", resourceID, err)
		}

		// Verify that IDPS is disabled (ids_enabled should be false)
		if obj.IdsEnabled != nil && *obj.IdsEnabled {
			return fmt.Errorf("Policy IdpsClusterConfig %s still has IDPS enabled (expected disabled after destroy)", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIdpsClusterConfigTemplate(computeCollection, displayName, description, idsEnabled string) string {
	// For NSX 9.1.0+: Use testAccNsxtPolicyIdpsClusterConfigTemplateWithDFW instead
	// For NSX < 9.1.0: DFW is not a prerequisite for IDPS (different architecture)
	return fmt.Sprintf(`
data "nsxt_compute_collection" "test" {
  display_name = "%s"
}

resource "nsxt_policy_idps_cluster_config" "test" {
  display_name = "%s"
  description  = "%s"
  ids_enabled  = %s

  cluster {
    target_id   = data.nsxt_compute_collection.test.id
    target_type = "VC_Cluster"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, computeCollection, displayName, description, idsEnabled)
}

func testAccNsxtPolicyIdpsClusterConfigTemplateWithDFW(computeCollection, displayName, description, idsEnabled string) string {
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
  description  = "%s"
  ids_enabled  = %s

  cluster {
    target_id   = data.nsxt_compute_collection.test.id
    target_type = "VC_Cluster"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, computeCollection, displayName, description, idsEnabled)
}

func testAccNsxtPolicyIdpsClusterConfigMinimalistic(computeCollection, displayName string) string {
	// For NSX 9.1.0+: Use testAccNsxtPolicyIdpsClusterConfigMinimalisticWithDFW instead
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
}`, computeCollection, displayName)
}

func testAccNsxtPolicyIdpsClusterConfigMinimalisticWithDFW(computeCollection, displayName string) string {
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
}`, computeCollection, displayName)
}
