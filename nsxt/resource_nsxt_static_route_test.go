/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtStaticRoute_basic(t *testing.T) {
	name := fmt.Sprintf("test-nsx-static-route")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_static_route.test"
	edgeClusterName := getEdgeClusterName()
	transportZoneName := getOverlayTransportZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXStaticRouteCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXStaticRouteCreateTemplate(name, edgeClusterName, transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXStaticRouteCheckExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "network", "4.4.4.0/24"),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.administrative_distance", "1"),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.ip_address", "8.0.0.10"),
					resource.TestCheckResourceAttrSet(testResourceName, "next_hop.0.logical_router_port_id"),
				),
			},
			{
				Config: testAccNSXStaticRouteUpdateTemplate(updateName, edgeClusterName, transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXStaticRouteCheckExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "network", "5.5.5.0/24"),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "2"),
				),
			},
		},
	})
}

func testAccNSXStaticRouteCheckExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX static route resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX static route resource ID not set in resources ")
		}
		routerID := rs.Primary.Attributes["logical_router_id"]
		if routerID == "" {
			return fmt.Errorf("NSX static route routerID not set in resources ")
		}

		staticRoute, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadStaticRoute(nsxClient.Context, routerID, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving static route ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if static route %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if displayName == staticRoute.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX static route %s wasn't found", displayName)
	}
}

func testAccNSXStaticRouteCheckDestroy(state *terraform.State, displayName string) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_static_route" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		routerID := rs.Primary.Attributes["logical_router_id"]
		staticRoute, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadStaticRoute(nsxClient.Context, routerID, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving static route ID %s. Error: %v", resourceID, err)
		}

		if displayName == staticRoute.DisplayName {
			return fmt.Errorf("NSX static route %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXStaticRoutePreConditionTemplate(edgeClusterName string, tzName string) string {
	return fmt.Sprintf(`
data "nsxt_edge_cluster" "EC" {
	display_name = "%s"
}

resource "nsxt_logical_tier1_router" "RTR1" {
	display_name = "tier1_router"
	edge_cluster_id = "${data.nsxt_edge_cluster.EC.id}"
}

data "nsxt_transport_zone" "TZ1" {
        display_name = "%s"
}

resource "nsxt_logical_switch" "LS1" {
	display_name = "downlink_test_switch"
	admin_state = "UP"
	replication_mode = "MTEP"
	vlan = "0"
	transport_zone_id = "${data.nsxt_transport_zone.TZ1.id}"
}

resource "nsxt_logical_port" "PORT1" {
	display_name = "LP"
	admin_state = "UP"
	description = "Acceptance Test"
	logical_switch_id = "${nsxt_logical_switch.LS1.id}"
}

resource "nsxt_logical_router_downlink_port" "LRP1" {
	display_name = "LRP"
	description = "Acceptance Test"
	linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
	logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
	ip_address = "8.0.0.1/24"
}
`, edgeClusterName, tzName)
}

func testAccNSXStaticRouteCreateTemplate(name string, edgeClusterName string, tzName string) string {
	return testAccNSXStaticRoutePreConditionTemplate(edgeClusterName, tzName) + fmt.Sprintf(`
resource "nsxt_static_route" "test" {
    logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
    display_name = "%s"
    description = "Acceptance Test"
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
    network = "4.4.4.0/24"
    next_hop {
	ip_address = "8.0.0.10"
	administrative_distance = "1"
	logical_router_port_id = "${nsxt_logical_router_downlink_port.LRP1.id}"
    }
}`, name)
}

func testAccNSXStaticRouteUpdateTemplate(name string, edgeClusterName string, tzName string) string {
	return testAccNSXStaticRoutePreConditionTemplate(edgeClusterName, tzName) + fmt.Sprintf(`
resource "nsxt_static_route" "test" {
    logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
    display_name = "%s"
    description = "Acceptance Test Update"
    network = "5.5.5.0/24"
    next_hop {
	ip_address = "8.0.0.10"
	administrative_distance = "1"
	logical_router_port_id = "${nsxt_logical_router_downlink_port.LRP1.id}"
    }
    next_hop {
	ip_address = "2.2.2.2"
	administrative_distance = "2"
    }
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
    tag {
    	scope = "scope2"
        tag = "tag2"
    }
}`, name)
}
