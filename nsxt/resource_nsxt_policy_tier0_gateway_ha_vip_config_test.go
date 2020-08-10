/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gm_locale_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s/locale_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services"
	"testing"
)

func TestAccResourceNsxtPolicyTier0GatewayHaVip_basic(t *testing.T) {
	name := "test-nsx-policy-tier0-ha-vip"
	subnet := "1.1.12.2/24"
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0HAVipConfigCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0HAVipConfigTemplate(subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0HAVipConfigExists(testResourceName),
					// resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					// resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					// resource.TestCheckResourceAttr(testResourceName, "mtu", mtu),
					// resource.TestCheckResourceAttr(testResourceName, "type", "EXTERNAL"),
					// resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					// resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					// resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					// resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					// resource.TestCheckResourceAttr(testResourceName, "enable_pim", "true"),
					// resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "STRICT"),
					// resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					// resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					// resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					// resource.TestCheckResourceAttrSet(testResourceName, "edge_node_path"),
					// resource.TestCheckResourceAttrSet(testResourceName, "path"),
					// resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					// resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					// resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTier0HAVipConfigExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		//connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier0 HA vip config resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier0 HA vip config resource ID not set in resources")
		}

		// var err error
		// localeServiceID := rs.Primary.Attributes["locale_service_id"]
		// gwID := getPolicyIDFromPath(rs.Primary.Attributes["gateway_path"])
		// if testAccIsGlobalManager() {
		// 	nsxClient := gm_locale_services.NewDefaultInterfacesClient(connector)
		// 	_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		// } else {
		// 	nsxClient := locale_services.NewDefaultInterfacesClient(connector)
		// 	_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		// }
		// if err != nil {
		// 	return fmt.Errorf("Error while retrieving policy Tier0 Interface ID %s. Error: %v", resourceID, err)
		// }

		return nil
	}
}

func testAccNsxtPolicyTier0HAVipConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier0_gateway_interface" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		localeServiceID := rs.Primary.Attributes["locale_service_id"]
		gwID := getPolicyIDFromPath(rs.Primary.Attributes["gateway_path"])

		var err error
		if testAccIsGlobalManager() {
			nsxClient := gm_locale_services.NewDefaultInterfacesClient(connector)
			_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		} else {
			nsxClient := locale_services.NewDefaultInterfacesClient(connector)
			_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		}
		if err == nil {
			return fmt.Errorf("Policy Tier0 Interface %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTier0HAVipConfigSiteTemplate() string {
	if testAccIsGlobalManager() {
		return fmt.Sprintf("site_path = data.nsxt_policy_site.test.path")
	}
	return ""
}

func testAccNsxtPolicyTier0HAVipConfigTemplate(subnet string) string {
	return testAccNsxtPolicyGatewayFabricInterfaceDeps() + fmt.Sprintf(`

resource "nsxt_policy_vlan_segment" "test1" {
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  display_name        = "interface_test1"
  vlan_ids            = [11]
  subnet {
      cidr = "10.2.2.2/24"
  }
}

resource "nsxt_policy_vlan_segment" "test2" {
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  display_name        = "interface_test2"
  vlan_ids            = [11]
  subnet {
      cidr = "10.2.2.2/24"
  }
}

data "nsxt_policy_edge_node" "EN0" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

data "nsxt_policy_edge_node" "EN1" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 1
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "ha_vip_test"
  ha_mode           = "ACTIVE_STANDBY"
  %s
}

resource "nsxt_policy_tier0_gateway_interface" "test1" {
  display_name   = "interface1"
  description    = "Acceptance Test"
  type           = "EXTERNAL"
  gateway_path   = nsxt_policy_tier0_gateway.test.path
  segment_path   = nsxt_policy_vlan_segment.test1.path
  edge_node_path = data.nsxt_policy_edge_node.EN0.path
  subnets        = ["%s"]
  urpf_mode      = "STRICT"
  %s
}

resource "nsxt_policy_tier0_gateway_interface" "test2" {
  display_name   = "interface2"
  description    = "Acceptance Test"
  type           = "EXTERNAL"
  gateway_path   = nsxt_policy_tier0_gateway.test.path
  segment_path   = nsxt_policy_vlan_segment.test2.path
  edge_node_path = data.nsxt_policy_edge_node.EN1.path
  subnets        = ["%s"]
  urpf_mode      = "STRICT"
  %s
}`, testAccNsxtPolicyTier0EdgeClusterTemplate(),
    subnet, testAccNsxtPolicyTier0HAVipConfigSiteTemplate(),
    subnet, testAccNsxtPolicyTier0HAVipConfigSiteTemplate())
}
