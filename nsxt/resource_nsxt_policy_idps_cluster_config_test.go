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

var accTestPolicyIdpsClusterConfigCreateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"target_type":        "VC_Cluster",
	"ids_enabled":        "true",
	"compute_collection": getTestComputeCollectionName(),
}

var accTestPolicyIdpsClusterConfigUpdateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform updated",
	"target_type":        "VC_Cluster",
	"ids_enabled":        "false",
	"compute_collection": getTestComputeCollectionName(),
}

func TestAccResourceNsxtPolicyIdpsClusterConfig_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_cluster_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, accTestPolicyIdpsClusterConfigUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(accTestPolicyIdpsClusterConfigCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIdpsClusterConfigCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIdpsClusterConfigCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", accTestPolicyIdpsClusterConfigCreateAttributes["ids_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", accTestPolicyIdpsClusterConfigCreateAttributes["target_type"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIdpsClusterConfigTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIdpsClusterConfigExists(accTestPolicyIdpsClusterConfigUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIdpsClusterConfigUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIdpsClusterConfigUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ids_enabled", accTestPolicyIdpsClusterConfigUpdateAttributes["ids_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "cluster.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "cluster.0.target_id"),
					resource.TestCheckResourceAttr(testResourceName, "cluster.0.target_type", accTestPolicyIdpsClusterConfigUpdateAttributes["target_type"]),

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
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_idps_cluster_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsClusterConfigCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsClusterConfigMinimalistic(),
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

func testAccNsxtPolicyIdpsClusterConfigTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIdpsClusterConfigCreateAttributes
	} else {
		attrMap = accTestPolicyIdpsClusterConfigUpdateAttributes
	}

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
    target_type = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["compute_collection"], attrMap["display_name"], attrMap["description"], attrMap["ids_enabled"], attrMap["target_type"])
}

func testAccNsxtPolicyIdpsClusterConfigMinimalistic() string {
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
}`, getTestComputeCollectionName(), getAccTestResourceName())
}
