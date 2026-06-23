// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyEdgeTransportNode_basic(t *testing.T) {
	testResourceName := "nsxt_policy_edge_transport_node.test"
	displayName := getAccTestResourceName()
	updatedDescription := "Updated acceptance test policy edge transport node"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPolicyEdgeTransportNodePreCheck(t)
			testAccNsxtExtraCoverage(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: withImportIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeTransportNodeCreateTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance test policy edge transport node"),
					resource.TestCheckResourceAttr(testResourceName, "form_factor", "SMALL"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEdgeTransportNodeWithDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.nsxt_policy_edge_transport_node.test", "display_name",
						testResourceName, "display_name",
					),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_edge_transport_node.test", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_edge_transport_node.test", "unique_id"),
				),
			},
			{
				Config: testAccNsxtPolicyEdgeTransportNodeUpdateTemplate(displayName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", updatedDescription),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
				ImportStateVerifyIgnore: []string{
					"credentials.0.cli_password",
					"credentials.0.root_password",
					// vm_deployment_config is not returned by the API after
					// deployment completes.
					"vm_deployment_config",
					"revision",
				},
			},
		}),
	})
}

func TestAccResourceNsxtPolicyEdgeTransportNode_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_edge_transport_node.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPolicyEdgeTransportNodePreCheck(t)
			testAccNsxtExtraCoverage(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: withImportIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeTransportNodeCreateTemplate(displayName),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
				ImportStateVerifyIgnore: []string{
					"credentials.0.cli_password",
					"credentials.0.root_password",
					"vm_deployment_config",
					"revision",
				},
			},
		}),
	})
}

func testAccNsxtPolicyEdgeTransportNodeCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliEdgeTransportNodesClient(sessionContext, connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_edge_transport_node" {
			continue
		}

		tnPath := rs.Primary.Attributes["path"]
		siteID, epID, id := getEdgeTransportNodeKeysFromPath(tnPath)
		_, err := client.Get(siteID, epID, id)
		if err == nil {
			return fmt.Errorf("policy edge transport node %s (%s) still exists", displayName, id)
		}
		if !isNotFoundError(err) {
			return fmt.Errorf("unexpected error checking policy edge transport node %s: %v", id, err)
		}
	}
	return nil
}

func testAccNsxtPolicyEdgeTransportNodeInfraTemplate() string {
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
  display_name = "test-policy-edge-uplink"
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

func testAccNsxtPolicyEdgeTransportNodeCreateTemplate(displayName string) string {
	return testAccNsxtPolicyEdgeTransportNodeInfraTemplate() + fmt.Sprintf(`
resource "nsxt_policy_edge_transport_node" "test" {
  display_name = "%s"
  description  = "Acceptance test policy edge transport node"
  form_factor  = "SMALL"
  hostname     = "%s"

  credentials {
    cli_password  = "%s"
    root_password = "%s"
  }

  management_interface {
    network_id = "%s"
    ip_assignment {
      dhcp_v4 = true
    }
  }

  switch {
    overlay_transport_zone_path     = data.nsxt_policy_transport_zone.test.path
    uplink_host_switch_profile_path = nsxt_policy_uplink_host_switch_profile.uplink.path
    pnic {
      device_name         = "fp-eth0"
      uplink_name         = "uplink1"
      datapath_network_id = "%s"
    }
    tunnel_endpoint {
      ip_assignment {
        dhcp_v4 = true
      }
    }
  }

  vm_deployment_config {
    compute_manager_id = data.nsxt_compute_manager.test.id
    compute_id         = data.nsxt_compute_collection.test.cm_local_id
    storage_id         = "%s"
  }
}

data "nsxt_policy_edge_transport_node_realization" "test" {
  path    = nsxt_policy_edge_transport_node.test.path
  timeout = 3600
}
`, displayName, getEdgeHostname(),
		os.Getenv("NSXT_PASSWORD"), os.Getenv("NSXT_PASSWORD"),
		getMgtNetworkID(), getEdgeDataNetworkID(),
		getDatastoreID())
}

func testAccNsxtPolicyEdgeTransportNodeWithDataSourceTemplate(displayName string) string {
	return testAccNsxtPolicyEdgeTransportNodeCreateTemplate(displayName) + `
data "nsxt_policy_edge_transport_node" "test" {
  display_name = nsxt_policy_edge_transport_node.test.display_name
}
`
}

func testAccNsxtPolicyEdgeTransportNodeUpdateTemplate(displayName, description string) string {
	return testAccNsxtPolicyEdgeTransportNodeInfraTemplate() + fmt.Sprintf(`
resource "nsxt_policy_edge_transport_node" "test" {
  display_name = "%s"
  description  = "%s"
  form_factor  = "SMALL"
  hostname     = "%s"

  credentials {
    cli_password  = "%s"
    root_password = "%s"
  }

  management_interface {
    network_id = "%s"
    ip_assignment {
      dhcp_v4 = true
    }
  }

  switch {
    overlay_transport_zone_path     = data.nsxt_policy_transport_zone.test.path
    uplink_host_switch_profile_path = nsxt_policy_uplink_host_switch_profile.uplink.path
    pnic {
      device_name         = "fp-eth0"
      uplink_name         = "uplink1"
      datapath_network_id = "%s"
    }
    tunnel_endpoint {
      ip_assignment {
        dhcp_v4 = true
      }
    }
  }

  vm_deployment_config {
    compute_manager_id = data.nsxt_compute_manager.test.id
    compute_id         = data.nsxt_compute_collection.test.cm_local_id
    storage_id         = "%s"
  }
}
`, displayName, description, getEdgeHostname(),
		os.Getenv("NSXT_PASSWORD"), os.Getenv("NSXT_PASSWORD"),
		getMgtNetworkID(), getEdgeDataNetworkID(),
		getDatastoreID())
}
