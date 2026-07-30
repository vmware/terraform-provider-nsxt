// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyTransitGatewayBfdPeer_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_bfd_peer.test"
	prereqName := getAccTestResourceName()
	displayName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayBfdPeerCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayBfdPeerTemplate(prereqName, displayName, "100.64.0.5", true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayBfdPeerExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", "100.64.0.5"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "source_attachment.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayBfdPeerTemplate(prereqName, displayName, "100.64.0.6", false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayBfdPeerExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", "100.64.0.6"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "source_attachment.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGatewayBfdPeer_importBasic(t *testing.T) {
	prereqName := getAccTestResourceName()
	displayName := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_bfd_peer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayBfdPeerCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayBfdPeerTemplate(prereqName, displayName, "100.64.0.5", true),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyTransitGatewayBfdPeerExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayBfdPeer resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayBfdPeer resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayBfdPeerExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayBfdPeer %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayBfdPeerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_transit_gateway_bfd_peer" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayBfdPeerExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy TransitGatewayBfdPeer %s still exists", displayName)
		}
	}
	return nil
}

// testAccNsxtPolicyTransitGatewayBfdPeerTemplate builds a Transit Gateway, a Centralized
// Network Attachment backed by a VPC subnet, and a Transit Gateway Attachment using that CNA,
// which serves as the required source_attachment for the BFD peer.
func testAccNsxtPolicyTransitGatewayBfdPeerTemplate(prereqName string, displayName string, peerAddress string, enabled bool) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_project" "test" {
  display_name = "%s"
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name    = "%s"
  transit_subnets = ["100.64.0.0/16"]
  centralized_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

resource "nsxt_vpc" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  private_ips  = ["192.168.240.0/24"]
}

resource "nsxt_vpc_subnet" "test" {
  context {
    project_id = nsxt_policy_project.test.id
    vpc_id     = nsxt_vpc.test.id
  }
  display_name = "%s"
  ip_addresses = ["192.168.240.0/26"]
  access_mode  = "Private"
}

resource "nsxt_policy_project_centralized_network_attachment" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  subnet_path  = nsxt_vpc_subnet.test.path
}

resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path  = nsxt_policy_transit_gateway.test.path
  display_name = "%s"
  cna_path     = nsxt_policy_project_centralized_network_attachment.test.path
}

resource "nsxt_policy_transit_gateway_bfd_peer" "test" {
  parent_path  = nsxt_policy_transit_gateway.test.path
  display_name = "%s"
  peer_address = "%s"
  enabled      = %t

  source_attachment = [nsxt_policy_transit_gateway_attachment.test.path]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getEdgeClusterName(), prereqName, prereqName, prereqName, prereqName, prereqName, prereqName, displayName, peerAddress, enabled)
}
