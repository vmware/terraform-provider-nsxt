// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/virtual_network_appliance_clusters"
)

func testAccVNAPreCheck(t *testing.T) {
	testAccPreCheck(t)
	testAccOnlyLocalManager(t)
	testAccNSXVersion(t, "9.1.1")
	testAccEnvDefined(t, "NSXT_TEST_EDGE_TRANSPORT_NODE")
	testAccEnvDefined(t, "NSXT_TEST_OVERLAY_TRANSPORT_ZONE")
	testAccEnvDefined(t, "NSXT_TEST_MGT_NETWORK")
	testAccEnvDefined(t, "NSXT_TEST_COMPUTE_MANAGER")
	testAccEnvDefined(t, "NSXT_TEST_COMPUTE_COLLECTION")
	testAccEnvDefined(t, "NSXT_TEST_DATASTORE_ID")
}

func TestAccResourceNsxtPolicyVirtualNetworkAppliance_basic(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance.test"
	testDataSourceName := "data.nsxt_policy_virtual_network_appliance.test"
	testRealizationName := "data.nsxt_policy_virtual_network_appliance_realization.test"
	displayName := getAccTestResourceName()
	updatedDisplayName := displayName + "-updated"
	edgeTransportNodeName := getEdgeTransportNodeName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVNAPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceCreateTemplate(displayName, edgeTransportNodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance test VNA"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "id", testResourceName, "id"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "path", testResourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testRealizationName, "state"),
				),
			},
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceUpdateTemplate(updatedDisplayName, edgeTransportNodeName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance test VNA - updated"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVirtualNetworkAppliance_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance.test"
	displayName := getAccTestResourceName()
	edgeTransportNodeName := getEdgeTransportNodeName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccVNAPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceCreateTemplate(displayName, edgeTransportNodeName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
				// credentials are write-only (not returned by GET) so exclude from import verify
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VirtualNetworkAppliance resource %s not found in resources", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("VirtualNetworkAppliance resource ID not set in resources")
		}

		clusterPath := rs.Primary.Attributes["cluster_path"]
		siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
		if err != nil {
			return err
		}

		connector := getPolicyConnector(testAccProvider.Meta())
		client := virtual_network_appliance_clusters.NewVirtualNetworkAppliancesClient(connector)
		_, err = client.Get(siteID, epID, clusterID, id)
		if err != nil {
			return fmt.Errorf("error while retrieving VirtualNetworkAppliance ID %s: %v", id, err)
		}

		return nil
	}
}

func testAccNsxtPolicyVirtualNetworkApplianceCheckDestroy(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta())
		client := virtual_network_appliance_clusters.NewVirtualNetworkAppliancesClient(connector)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "nsxt_policy_virtual_network_appliance" {
				continue
			}

			id := rs.Primary.ID
			clusterPath := rs.Primary.Attributes["cluster_path"]
			siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
			if err != nil {
				return err
			}

			_, err = client.Get(siteID, epID, clusterID, id)
			if err == nil {
				return fmt.Errorf("VirtualNetworkAppliance %s still exists", id)
			}
			if !isNotFoundError(err) {
				return fmt.Errorf("unexpected error checking VirtualNetworkAppliance %s: %v", id, err)
			}
		}

		return nil
	}
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterTemplate(displayName, edgeTransportNodeName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_transport_node" "test" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

data "nsxt_compute_manager" "test" {
  display_name = "%s"
}

data "nsxt_compute_collection" "test" {
  display_name = "%s"
  origin_id    = data.nsxt_compute_manager.test.id
}

resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s-cluster"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  member {
    edge_transport_node_path = data.nsxt_policy_edge_transport_node.test.path
  }

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.test.path
  }
}
`, edgeTransportNodeName, getOverlayTransportZoneName(), getComputeManagerName(), getComputeCollectionName(), displayName)
}

func testAccNsxtPolicyVirtualNetworkApplianceCreateTemplate(displayName, edgeTransportNodeName string) string {
	return testAccNsxtPolicyVirtualNetworkApplianceClusterTemplate(displayName, edgeTransportNodeName) + fmt.Sprintf(`
resource "nsxt_policy_virtual_network_appliance" "test" {
  display_name = "%s"
  description  = "Acceptance test VNA"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.test.path
  hostname     = "%s"

  management_interface {
    network_id = "%s"

    ip_assignment {
      dhcp_v4 = true
    }
  }

  vm_deployment_config {
    compute_manager_id          = data.nsxt_compute_manager.test.id
    cluster_or_resource_pool_id = data.nsxt_compute_collection.test.cm_local_id
    datastore_id                = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_virtual_network_appliance" "test" {
  display_name = nsxt_policy_virtual_network_appliance.test.display_name
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.test.path

  depends_on = [nsxt_policy_virtual_network_appliance.test]
}

data "nsxt_policy_virtual_network_appliance_realization" "test" {
  path    = nsxt_policy_virtual_network_appliance.test.path
  timeout = 5400

  depends_on = [nsxt_policy_virtual_network_appliance.test]
}
`, displayName, getVNAHostname(), getMgtNetworkID(), getDatastoreID())
}

func testAccNsxtPolicyVirtualNetworkApplianceUpdateTemplate(displayName, edgeTransportNodeName string) string {
	return testAccNsxtPolicyVirtualNetworkApplianceClusterTemplate(displayName, edgeTransportNodeName) + fmt.Sprintf(`
resource "nsxt_policy_virtual_network_appliance" "test" {
  display_name = "%s"
  description  = "Acceptance test VNA - updated"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.test.path
  hostname     = "%s"

  management_interface {
    network_id = "%s"

    ip_assignment {
      dhcp_v4 = true
    }
  }

  vm_deployment_config {
    compute_manager_id          = data.nsxt_compute_manager.test.id
    cluster_or_resource_pool_id = data.nsxt_compute_collection.test.cm_local_id
    datastore_id                = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, displayName, getVNAHostname(), getMgtNetworkID(), getDatastoreID())
}
