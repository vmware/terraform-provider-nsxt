// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	enforcement_points "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
)

const testAccVNSClusterResourceName = "nsxt_policy_virtual_network_service_cluster.test"

func testAccPreCheckVNSCluster(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.2.0")
	testAccNsxtExtraCoverage(t)
}

func TestAccResourceNsxtPolicyVirtualNetworkServiceCluster_basic(t *testing.T) {
	displayName := getAccTestResourceName()
	updatedDisplayName := displayName + "-updated"

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
				Config: testAccNsxtPolicyVirtualNetworkServiceClusterTemplate(displayName, "MEDIUM", 2, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkServiceClusterExists(testAccVNSClusterResourceName),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "description", "Acceptance test cluster"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "node_form_factor", "MEDIUM"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "number_of_nodes", "2"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "default_project_service_allocation", "false"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "vsphere_deployment_config.#", "1"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "vsphere_deployment_config.0.cluster_ids.#", "1"),
					resource.TestCheckResourceAttrSet(testAccVNSClusterResourceName, "vsphere_deployment_config.0.compute_manager_id"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testAccVNSClusterResourceName, "path"),
					resource.TestCheckResourceAttrSet(testAccVNSClusterResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyVirtualNetworkServiceClusterTemplate(updatedDisplayName, "LARGE", 3, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkServiceClusterExists(testAccVNSClusterResourceName),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "node_form_factor", "LARGE"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "number_of_nodes", "3"),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "default_project_service_allocation", "true"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVirtualNetworkServiceCluster_importBasic(t *testing.T) {
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
				Config: testAccNsxtPolicyVirtualNetworkServiceClusterTemplate(displayName, "MEDIUM", 2, false),
			},
			{
				ResourceName:            testAccVNSClusterResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccResourceNsxtPolicyImportIDRetriever(testAccVNSClusterResourceName),
				ImportStateVerifyIgnore: []string{"revision"},
			},
		},
	})
}

// Verifies that omitting every optional attribute does not cause perpetual drift
// when NSX fills in server-side defaults.
func TestAccResourceNsxtPolicyVirtualNetworkServiceCluster_minimal(t *testing.T) {
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
				Config: testAccNsxtPolicyVirtualNetworkServiceClusterMinimalTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkServiceClusterExists(testAccVNSClusterResourceName),
					resource.TestCheckResourceAttr(testAccVNSClusterResourceName, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testAccVNSClusterResourceName, "node_form_factor"),
					resource.TestCheckResourceAttrSet(testAccVNSClusterResourceName, "number_of_nodes"),
				),
			},
			{
				Config:             testAccNsxtPolicyVirtualNetworkServiceClusterMinimalTemplate(displayName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkServiceClusterExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VirtualNetworkServiceCluster resource %s not found in resources", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("VirtualNetworkServiceCluster resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta())
		client := enforcement_points.NewVirtualNetworkServiceClustersClient(connector)
		_, err := client.Get("default", "default", id)
		if err != nil {
			return fmt.Errorf("error while retrieving VirtualNetworkServiceCluster ID %s: %v", id, err)
		}

		return nil
	}
}

func testAccNsxtPolicyVirtualNetworkServiceClusterCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta())
	client := enforcement_points.NewVirtualNetworkServiceClustersClient(connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_virtual_network_service_cluster" {
			continue
		}

		id := rs.Primary.ID
		_, err := client.Get("default", "default", id)
		if err == nil {
			return fmt.Errorf("VirtualNetworkServiceCluster %s still exists", id)
		}
		if !isNotFoundError(err) {
			return fmt.Errorf("unexpected error checking VirtualNetworkServiceCluster %s: %v", id, err)
		}
	}

	return nil
}

func testAccNsxtPolicyVirtualNetworkServiceClusterPrereqs() string {
	return fmt.Sprintf(`
data "nsxt_compute_manager" "test" {
  display_name = "%s"
}

data "nsxt_compute_collection" "test" {
  display_name = "%s"
}
`, os.Getenv("NSXT_TEST_COMPUTE_MANAGER"), os.Getenv("NSXT_TEST_COMPUTE_COLLECTION"))
}

func testAccNsxtPolicyVirtualNetworkServiceClusterTemplate(displayName, formFactor string, nodes int, allocation bool) string {
	return testAccNsxtPolicyVirtualNetworkServiceClusterPrereqs() + fmt.Sprintf(`
resource "nsxt_policy_virtual_network_service_cluster" "test" {
  display_name                       = "%s"
  description                        = "Acceptance test cluster"
  node_form_factor                   = "%s"
  number_of_nodes                    = %d
  default_project_service_allocation = %t

  vsphere_deployment_config {
    compute_manager_id = data.nsxt_compute_manager.test.id
    cluster_ids        = [data.nsxt_compute_collection.test.cm_local_id]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, displayName, formFactor, nodes, allocation)
}

func testAccNsxtPolicyVirtualNetworkServiceClusterMinimalTemplate(displayName string) string {
	return testAccNsxtPolicyVirtualNetworkServiceClusterPrereqs() + fmt.Sprintf(`
resource "nsxt_policy_virtual_network_service_cluster" "test" {
  display_name = "%s"

  vsphere_deployment_config {
    compute_manager_id = data.nsxt_compute_manager.test.id
    cluster_ids        = [data.nsxt_compute_collection.test.cm_local_id]
  }
}
`, displayName)
}
