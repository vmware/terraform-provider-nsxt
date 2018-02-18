/* Copyright © 2017 VMware, Inc. All Rights Reserved.
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

func TestNSXLogicalRouterDownlinkPortBasic(t *testing.T) {

	portName := fmt.Sprintf("test-nsx-logical-router-downlink-port")
	updatePortName := fmt.Sprintf("%s-update", portName)
	testResourceName := "nsxt_logical_router_downlink_port.test"
	transportZoneName := overlayTransportZoneNamePrefix
	edgeClusterName := edgeClusterDefaultName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalRouterDownlinkPortCheckDestroy(state, portName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalRouterDownlinkPortCreateTemplate(portName, transportZoneName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterDownlinkPortExists(portName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", portName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "1"),
				),
			},
			{
				Config: testAccNSXLogicalRouterDownlinkPortUpdateTemplate(updatePortName, transportZoneName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterDownlinkPortExists(updatePortName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePortName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "2"),
				),
			},
		},
	})
}

func TestNSXLogicalRouterDownlinkPortWithRelay(t *testing.T) {

	portName := fmt.Sprintf("test-nsx-logical-router-downlink-port")
	updatePortName := fmt.Sprintf("%s-update", portName)
	testResourceName := "nsxt_logical_router_downlink_port.test"
	transportZoneName := overlayTransportZoneNamePrefix
	edgeClusterName := edgeClusterDefaultName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalRouterDownlinkPortCheckDestroy(state, portName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalRouterDownlinkPortCreateWithRelayTemplate(portName, transportZoneName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterDownlinkPortExists(portName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", portName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.0.target_type", "LogicalService"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.0.target_display_name", "srv"),
				),
			},
			{
				Config: testAccNSXLogicalRouterDownlinkPortUpdateWithRelayTemplate(updatePortName, transportZoneName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterDownlinkPortExists(updatePortName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatePortName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttrSet(testResourceName, "linked_logical_switch_port_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "logical_router_id"),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.0.target_type", "LogicalService"),
					resource.TestCheckResourceAttr(testResourceName, "service_binding.0.target_display_name", "srv"),
				),
			},
		},
	})
}

func testAccNSXLogicalRouterDownlinkPortExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical port resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical port resource ID not set in resources ")
		}

		logicalPort, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterDownLinkPort(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical port ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical port %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == logicalPort.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical port %s wasn't found", display_name)
	}
}

func testAccNSXLogicalRouterDownlinkPortCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_router_downlink_port" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		logicalPort, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterDownLinkPort(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving logical port ID %s. Error: %v", resourceID, err)
		}

		if display_name == logicalPort.DisplayName {
			return fmt.Errorf("NSX logical port %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXLogicalRouterDownlinkPortPreConditionTemplate(transportZoneName string, edgeClusterName string) string {
	return fmt.Sprintf(`
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

data "nsxt_edge_cluster" "EC" {
	display_name = "%s"
}

resource "nsxt_logical_tier1_router" "RTR1" {
	display_name = "downlink_test_router"
	edge_cluster_id = "${data.nsxt_edge_cluster.EC.id}"
}`, transportZoneName, edgeClusterName)
}

func testAccNSXLogicalRouterDownlinkPortRelayTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_dhcp_relay_profile" "DRP1" {
	display_name = "prf"
	server_addresses = ["1.1.1.1"]
}

resource "nsxt_dhcp_relay_service" "DRS1" {
	display_name = "srv"
	description = "Acceptance Test"
	dhcp_relay_profile_id = "${nsxt_dhcp_relay_profile.DRP1.id}"
}`)
}

func testAccNSXLogicalRouterDownlinkPortCreateTemplate(portName string, transportZoneName string, edgeClusterName string) string {
	return testAccNSXLogicalRouterDownlinkPortPreConditionTemplate(transportZoneName, edgeClusterName) + fmt.Sprintf(`
resource "nsxt_logical_router_downlink_port" "test" {
	display_name = "%s"
	description = "Acceptance Test"
	linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
	logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
	subnets = [{ip_addresses = ["8.0.0.1"],
    	        prefix_length = 24}]
	tags = [{scope = "scope1", tag = "tag1"}]
}`, portName)
}

func testAccNSXLogicalRouterDownlinkPortUpdateTemplate(portUpdatedName string, transportZoneName string, edgeClusterName string) string {
	return testAccNSXLogicalRouterDownlinkPortPreConditionTemplate(transportZoneName, edgeClusterName) + fmt.Sprintf(`
resource "nsxt_logical_router_downlink_port" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
	linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
	logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
	subnets = [{ip_addresses = ["8.0.0.1"],
            	prefix_length = 24}]
	tags = [{scope = "scope1", tag = "tag1"},
			{scope = "scope2", tag = "tag2"}]
}`, portUpdatedName)
}

func testAccNSXLogicalRouterDownlinkPortCreateWithRelayTemplate(portName string, transportZoneName string, edgeClusterName string) string {
	return testAccNSXLogicalRouterDownlinkPortPreConditionTemplate(transportZoneName, edgeClusterName) +
		testAccNSXLogicalRouterDownlinkPortRelayTemplate() +
		fmt.Sprintf(`
resource "nsxt_logical_router_downlink_port" "test" {
	display_name = "%s"
	description = "Acceptance Test"
	linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
	logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
	subnets = [{ip_addresses = ["8.0.0.1"],
	            prefix_length = 24}]
	service_binding {
		target_id = "${nsxt_dhcp_relay_service.DRS1.id}"
		target_type = "LogicalService"
	}
	tags = [{scope = "scope1"
	    	 tag = "tag1"}
	]
}`, portName)
}

func testAccNSXLogicalRouterDownlinkPortUpdateWithRelayTemplate(portUpdatedName string, transportZoneName string, edgeClusterName string) string {
	return testAccNSXLogicalRouterDownlinkPortPreConditionTemplate(transportZoneName, edgeClusterName) +
		testAccNSXLogicalRouterDownlinkPortRelayTemplate() +
		fmt.Sprintf(`
resource "nsxt_logical_router_downlink_port" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
	linked_logical_switch_port_id = "${nsxt_logical_port.PORT1.id}"
	logical_router_id = "${nsxt_logical_tier1_router.RTR1.id}"
	subnets = [{ip_addresses = ["8.0.0.1"],
	            prefix_length = 24}]
	service_binding {
		target_id = "${nsxt_dhcp_relay_service.DRS1.id}"
		target_type = "LogicalService"
	}
	tags = [{scope = "scope1"
	    	 tag = "tag1"},
		    {scope = "scope2"
	    	 tag = "tag2"}
	]
}`, portUpdatedName)
}
