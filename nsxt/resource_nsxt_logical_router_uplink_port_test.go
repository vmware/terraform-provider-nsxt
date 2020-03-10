package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

func TestAccResourceNsxtLogicalRouterUpLinkPort_basic(t *testing.T) {
	portName := "test-nsx-uplink-port"
	updatedPortName := "test-nsx-uplink-port-updated"
	tier0RouterName := getTier0RouterName()
	transportZoneName := getVlanTransportZoneName()
	edgeClusterMemberIndex, err := getEdgeClusterMemberIndex()
	if err != nil {
		t.Skip(err)
	}
	testResourceName := "nsxt_logical_router_uplink_port.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtLogicalRouterUpLinkPortCreateTemplate(portName, tier0RouterName, transportZoneName, edgeClusterMemberIndex),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtLogicalRouterUpLinkPortExists(portName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", portName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "NONE"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "edge_cluster_member_index.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
				),
			},
			{
				Config: testAccNsxtLogicalRouterUpLinkPortUpdateTemplate(updatedPortName, tier0RouterName, transportZoneName, edgeClusterMemberIndex),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtLogicalRouterUpLinkPortExists(updatedPortName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedPortName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Updated"),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "STRICT"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "edge_cluster_member_index.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
				),
			},
		},
	})
}

func testAccNsxtLogicalRouterUpLinkPortExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical router uplink port resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical router uplink port resource ID not set in resources ")
		}

		logicalPort, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterUpLinkPort(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical router uplink port ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical router uplink port %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if displayName == logicalPort.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical port %s wasn't found", displayName)
	}
}

func testAccNsxtLogicalRouterUpLinkPortPreconditionsTemplate(tier0RouterName string, transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_logical_tier0_router" "tier0rtr" {
  display_name = "%s"
}

data "nsxt_transport_zone" "tz1" {
  display_name = "%s"
}

resource "nsxt_logical_switch" "ls1" {
  display_name     = "uplink_test_switch"
  admin_state      = "UP"
  replication_mode = "MTEP"
  transport_zone_id = "${data.nsxt_transport_zone.tz1.id}"
}

resource "nsxt_logical_port" "port1" {
  display_name      = "LP"
  admin_state       = "UP"
  description       = "Acceptance Test"
  logical_switch_id = "${nsxt_logical_switch.ls1.id}"
}`, tier0RouterName, transportZoneName)
}

func testAccNsxtLogicalRouterUpLinkPortCreateTemplate(portName string, tier0RouterName string, transportZoneName string, edgeClusterMemberIndex int64) string {
	return testAccNsxtLogicalRouterUpLinkPortPreconditionsTemplate(tier0RouterName, transportZoneName) + fmt.Sprintf(`
resource "nsxt_logical_router_uplink_port" "test" {
  display_name                  = "%s"
  description                   = "Acceptance Test"
  logical_router_id             = "${data.nsxt_logical_tier0_router.tier0rtr.id}"
  linked_logical_switch_port_id = "${nsxt_logical_port.port1.id}"
  edge_cluster_member_index     = [%d]
  urpf_mode                     = "NONE"

  subnets {
    ip_addresses = ["100.101.102.103"]
    prefix_length = 24
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, portName, edgeClusterMemberIndex)
}

func testAccNsxtLogicalRouterUpLinkPortUpdateTemplate(portName string, tier0RouterName string, transportZoneName string, edgeClusterMemberIndex int64) string {
	return testAccNsxtLogicalRouterUpLinkPortPreconditionsTemplate(tier0RouterName, transportZoneName) + fmt.Sprintf(`
resource "nsxt_logical_router_uplink_port" "test" {
  display_name                  = "%s"
  description                   = "Acceptance Test Updated"
  logical_router_id             = "${data.nsxt_logical_tier0_router.tier0rtr.id}"
  linked_logical_switch_port_id = "${nsxt_logical_port.port1.id}"
  edge_cluster_member_index     = [%d]
  urpf_mode                     = "STRICT"

  subnets {
    ip_addresses = ["100.101.102.103"]
    prefix_length = 24
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, portName, edgeClusterMemberIndex)
}
