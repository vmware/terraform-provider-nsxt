// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtEdgeTransportNode_basic(t *testing.T) {
	testResourceName := "nsxt_edge_transport_node.test"
	displayName := getAccTestResourceName()
	updatedDescription := "Updated acceptance test edge transport node"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEdgeTransportNodePreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
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
					// node_settings are not read back from the NSX API (Read skips
					// them to avoid MPA-sync overwrite drift).  Users must supply
					// node_settings in their config after import.
					"node_settings",
				},
			},
		},
	})
}

func TestAccResourceNsxtEdgeTransportNode_importBasic(t *testing.T) {
	testResourceName := "nsxt_edge_transport_node.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccEdgeTransportNodePreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
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
					// node_settings are not read back from the NSX API.
					"node_settings",
				},
			},
		},
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

// TestAccResourceNsxtEdgeTransportNodeNodeSettings verifies that updating
// node_settings does not produce perpetual drift (bug 3689817).
//
// Root cause: when a freshly-deployed edge VM boots, its management plane
// agent (MPA) performs an initial sync that reports the VM's current running
// state to NSX.  NSX writes those MPA-reported values back into
// node_deployment_info.node_settings, temporarily overwriting the desired
// state that Terraform had just PUT.  Reading those stale values in the
// subsequent resourceNsxtEdgeTransportNodeRead would cause the same changes
// to re-appear on every plan.
//
// Fix: Read intentionally skips node_settings; the config is the authoritative
// source.  Empirical verification (nsxedge.1 on 10.162.50.158) confirmed that
// for an already-initialized edge GET reflects desired state immediately after
// PUT — the overwrite is specific to the MPA initial-boot sync window.
//
// This test exercises the two-step sequence from the failing FVT:
//  1. Create edge TN with minimal node_settings.
//  2. Update node_settings to all supported attributes.
//     A 60-second PreConfig delay simulates the MPA sync window.
//  3. Plan-only after another 60-second delay must show 0 changes (no drift).
//
// Passed against NSX 10.162.50.158 / VC 10.162.51.249 (824 s).
func TestAccResourceNsxtEdgeTransportNodeNodeSettings(t *testing.T) {
	testResourceName := "nsxt_edge_transport_node.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccEdgeTransportNodePreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtEdgeTransportNodeCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				// Step 1: deploy edge TN with minimal node_settings.
				Config: testAccNsxtEdgeTransportNodeNodeSettingsTemplate(displayName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.hostname", "edge-node-tf"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.allow_ssh_root_login", "false"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.enable_ssh", "false"),
				),
			},
			{
				// Step 2: update node_settings to cover all supported fields.
				// PreConfig: wait for the MPA initial-boot sync to complete
				// before issuing the update.
				PreConfig: func() { time.Sleep(60 * time.Second) },
				Config:    testAccNsxtEdgeTransportNodeNodeSettingsTemplate(displayName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.hostname", "edge-node-tf-updated"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.allow_ssh_root_login", "true"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.enable_ssh", "true"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.ntp_servers.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "node_settings.0.syslog_server.#", "2"),
				),
			},
			{
				// Step 3 (no-drift check): re-plan with the same config.
				// PreConfig: give the MPA sync window time to pass before the
				// plan.  Without the fix, Read would return MPA-overwritten
				// hostname/SSH/NTP/syslog values, causing perpetual drift.
				PreConfig:          func() { time.Sleep(60 * time.Second) },
				Config:             testAccNsxtEdgeTransportNodeNodeSettingsTemplate(displayName, true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccNsxtEdgeTransportNodeNodeSettingsTemplate(displayName string, updated bool) string {
	var nodeSettings string
	if updated {
		nodeSettings = `
  node_settings {
    hostname             = "edge-node-tf-updated"
    allow_ssh_root_login = true
    enable_ssh           = true
    ntp_servers          = ["time.vmware.com", "pool.ntp.org"]
    syslog_server {
      server    = "10.20.30.40"
      log_level = "INFO"
      port      = "514"
      protocol  = "UDP"
    }
    syslog_server {
      server    = "10.20.30.41"
      log_level = "WARNING"
      port      = "6514"
      protocol  = "TLS"
    }
  }`
	} else {
		nodeSettings = `
  node_settings {
    hostname = "edge-node-tf"
  }`
	}

	return testAccNsxtEdgeTransportNodeInfraTemplate() + fmt.Sprintf(`
resource "nsxt_edge_transport_node" "test" {
  display_name = "%s"
  description  = "node_settings drift acceptance test (bug 3689817)"

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
%s
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
		nodeSettings)
}
