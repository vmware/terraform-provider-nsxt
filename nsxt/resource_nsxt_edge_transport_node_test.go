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
)

func TestAccResourceNsxtEdgeTransportNode_basic(t *testing.T) {
	testResourceName := "nsxt_edge_transport_node.test"
	displayName := getAccTestResourceName()
	updatedDescription := "Updated acceptance test edge transport node"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEdgeTransportNodePreCheck(t)
			testAccNsxtExtraCoverage(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtEdgeTransportNodeCreateTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance test edge transport node"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
			{
				Config: testAccNsxtEdgeTransportNodeWithDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.nsxt_transport_node.test", "display_name",
						testResourceName, "display_name",
					),
					resource.TestCheckResourceAttrSet("data.nsxt_transport_node.test", "id"),
				),
			},
			{
				Config: testAccNsxtEdgeTransportNodeUpdateTemplate(displayName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", updatedDescription),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"deployment_config.0.node_user_settings.0.cli_password",
					"deployment_config.0.node_user_settings.0.root_password",
					// ip_addresses is assigned by NSX after deployment and may not
					// yet be populated in the initial apply state.
					"ip_addresses",
					// revision may be incremented by NSX between apply and import.
					"revision",
					// On a fresh import there is no prior state to restore policy
					// paths from (etNodeNormalizeHostSwitchesInState requires them).
					"standard_host_switch.0.uplink_profile",
					"standard_host_switch.0.transport_zone_endpoint.0.transport_zone",
				},
			},
		}),
	})
}

func TestAccResourceNsxtEdgeTransportNode_importBasic(t *testing.T) {
	testResourceName := "nsxt_edge_transport_node.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEdgeTransportNodePreCheck(t)
			testAccNsxtExtraCoverage(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtEdgeTransportNodeCreateTemplate(displayName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"deployment_config.0.node_user_settings.0.cli_password",
					"deployment_config.0.node_user_settings.0.root_password",
					"ip_addresses",
					"revision",
					"standard_host_switch.0.uplink_profile",
					"standard_host_switch.0.transport_zone_endpoint.0.transport_zone",
				},
			},
		}),
	})
}

func testAccNsxtEdgeTransportNodeCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	client := cliTransportNodesClient(connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_edge_transport_node" {
			continue
		}

		id := rs.Primary.ID
		_, err := client.Get(id)
		if err == nil {
			return fmt.Errorf("edge transport node %s (%s) still exists", displayName, id)
		}
		if !isNotFoundError(err) {
			return fmt.Errorf("unexpected error checking edge transport node %s: %v", id, err)
		}
	}
	return nil
}

func testAccNsxtEdgeTransportNodeInfraTemplate() string {
	return fmt.Sprintf(`
data "nsxt_compute_manager" "test" {
  display_name = "%s"
}

data "nsxt_compute_collection" "test" {
  display_name = "%s"
  origin_id    = data.nsxt_compute_manager.test.id
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_uplink_host_switch_profile" "uplink" {
  display_name = "test-edge-uplink"
  teaming {
    active {
      uplink_name = "uplink1"
      uplink_type = "PNIC"
    }
    policy = "FAILOVER_ORDER"
  }
}
`, getComputeManagerName(), getComputeCollectionName(), getOverlayTransportZoneName())
}

func testAccNsxtEdgeTransportNodeCreateTemplate(displayName string) string {
	return testAccNsxtEdgeTransportNodeInfraTemplate() + fmt.Sprintf(`
resource "nsxt_edge_transport_node" "test" {
  display_name = "%s"
  description  = "Acceptance test edge transport node"

  deployment_config {
    form_factor = "SMALL"
    node_user_settings {
      cli_password  = "%s"
      root_password = "%s"
    }
    vm_deployment_config {
      vc_id                 = data.nsxt_compute_manager.test.id
      compute_id            = data.nsxt_compute_collection.test.cm_local_id
      storage_id            = "%s"
      management_network_id = "%s"
      data_network_ids      = ["%s"]
    }
  }

  node_settings {
    hostname = "%s"
  }

  standard_host_switch {
    uplink_profile = nsxt_policy_uplink_host_switch_profile.uplink.path
    ip_assignment {
      assigned_by_dhcp = true
    }
    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.test.path
    }
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
}

data "nsxt_transport_node_realization" "test" {
  id      = nsxt_edge_transport_node.test.id
  timeout = 3600
}
`, displayName,
		os.Getenv("NSXT_PASSWORD"), os.Getenv("NSXT_PASSWORD"),
		getDatastoreID(), getMgtNetworkID(), getEdgeDataNetworkID(),
		getEdgeHostname())
}

func testAccNsxtEdgeTransportNodeWithDataSourceTemplate(displayName string) string {
	return testAccNsxtEdgeTransportNodeCreateTemplate(displayName) + `
data "nsxt_transport_node" "test" {
  display_name = nsxt_edge_transport_node.test.display_name
}
`
}

func testAccNsxtEdgeTransportNodeUpdateTemplate(displayName, description string) string {
	return testAccNsxtEdgeTransportNodeInfraTemplate() + fmt.Sprintf(`
resource "nsxt_edge_transport_node" "test" {
  display_name = "%s"
  description  = "%s"

  deployment_config {
    form_factor = "SMALL"
    node_user_settings {
      cli_password  = "%s"
      root_password = "%s"
    }
    vm_deployment_config {
      vc_id                 = data.nsxt_compute_manager.test.id
      compute_id            = data.nsxt_compute_collection.test.cm_local_id
      storage_id            = "%s"
      management_network_id = "%s"
      data_network_ids      = ["%s"]
    }
  }

  node_settings {
    hostname = "%s"
  }

  standard_host_switch {
    uplink_profile = nsxt_policy_uplink_host_switch_profile.uplink.path
    ip_assignment {
      assigned_by_dhcp = true
    }
    transport_zone_endpoint {
      transport_zone = data.nsxt_policy_transport_zone.test.path
    }
    pnic {
      device_name = "fp-eth0"
      uplink_name = "uplink1"
    }
  }
}
`, displayName, description,
		os.Getenv("NSXT_PASSWORD"), os.Getenv("NSXT_PASSWORD"),
		getDatastoreID(), getMgtNetworkID(), getEdgeDataNetworkID(),
		getEdgeHostname())
}
