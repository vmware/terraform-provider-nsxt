/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services"
	"testing"
)

var nsxtPolicyTier0GatewayName = "test"
var nsxtPolicyTier0GatewayID = "test"

func TestAccResourceNsxtPolicyTier0GatewayInterface_service(t *testing.T) {
	name := "test-nsx-policy-tier0-interface-basic"
	updatedName := fmt.Sprintf("%s-update", name)
	mtu := "1500"
	updatedMtu := "1800"
	subnet := "1.1.12.2/24"
	updatedSubnet := "1.2.12.2/24"
	ipAddress := "1.1.12.2"
	updatedIPAddress := "1.2.12.2"
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceServiceTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", mtu),
					resource.TestCheckResourceAttr(testResourceName, "type", "SERVICE"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceServiceTemplate(updatedName, updatedSubnet, updatedMtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", updatedMtu),
					resource.TestCheckResourceAttr(testResourceName, "type", "SERVICE"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", updatedSubnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", updatedIPAddress),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceThinTemplate(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "mtu", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0GatewayInterface_external(t *testing.T) {
	name := "test-nsx-policy-tier0-interface-basic"
	updatedName := fmt.Sprintf("%s-update", name)
	mtu := "1500"
	updatedMtu := "1800"
	subnet := "1.1.12.2/24"
	updatedSubnet := "1.2.12.2/24"
	ipAddress := "1.1.12.2"
	updatedIPAddress := "1.2.12.2"
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceExternalTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", mtu),
					resource.TestCheckResourceAttr(testResourceName, "type", "EXTERNAL"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					resource.TestCheckResourceAttr(testResourceName, "enable_pim", "true"),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "STRICT"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_node_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceExternalTemplate(updatedName, updatedSubnet, updatedMtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", updatedMtu),
					resource.TestCheckResourceAttr(testResourceName, "type", "EXTERNAL"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", updatedSubnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", updatedIPAddress),
					resource.TestCheckResourceAttr(testResourceName, "enable_pim", "true"),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "STRICT"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_node_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0GatewayInterface_withID(t *testing.T) {
	name := "test-nsx-policy-tier0-interface-basic"
	updatedName := fmt.Sprintf("%s-update", name)
	subnet := "1.1.12.2/24"
	// Update to 2 addresses
	ipv6Subnet := "4003::12/64"
	updatedSubnet := fmt.Sprintf("%s\",\"%s", subnet, ipv6Subnet)
	ipAddress := "1.1.12.2"
	ipv6Address := "4003::12"
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceTemplateWithID(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "test"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceTemplateWithID(updatedName, updatedSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "test"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "subnets.1", ipv6Subnet),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", ipAddress),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.1", ipv6Address),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier0InterfaceImporterGetID(s *terraform.State) (string, error) {
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"
	rs, ok := s.RootModule().Resources[testResourceName]
	if !ok {
		return "", fmt.Errorf("NSX Policy Tier0 Interface resource %s not found in resources", testResourceName)
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy Tier0 Interface resource ID not set in resources ")
	}
	gwPath := rs.Primary.Attributes["gateway_path"]
	if gwPath == "" {
		return "", fmt.Errorf("NSX Policy Interface Gateway Policy Path not set in resources ")
	}
	_, gwID := parseGatewayPolicyPath(gwPath)
	return fmt.Sprintf("%s/%s/%s", gwID, defaultPolicyLocaleServiceID, resourceID), nil
}

func TestAccResourceNsxtPolicyTier0GatewayInterface_importBasic(t *testing.T) {
	name := "test-nsx-policy-tier0-interface-import"
	testResourceName := "nsxt_policy_tier0_gateway_interface.test"
	subnet := "1.1.12.2/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceThinTemplate(name, subnet),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyTier0InterfaceImporterGetID,
			},
		},
	})
}

func testAccNsxtPolicyTier0InterfaceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := locale_services.NewDefaultInterfacesClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier0 Interface resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier0 Interface resource ID not set in resources")
		}

		_, err := nsxClient.Get(nsxtPolicyTier0GatewayID, defaultPolicyLocaleServiceID, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy Tier0 Interface ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyTier0InterfaceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := locale_services.NewDefaultInterfacesClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier0_gateway_interface" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(nsxtPolicyTier0GatewayID, defaultPolicyLocaleServiceID, resourceID)
		if err == nil {
			return fmt.Errorf("Policy Tier0 Interface %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyGatewayInterfaceDeps(vlans string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_vlan_segment" "test" {
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  display_name        = "interface_test"
  vlan_ids            = [%s]
  subnet {
      cidr = "10.2.2.2/24"
  }
}`, getEdgeClusterName(), getVlanTransportZoneName(), vlans)

}

func testAccNsxtPolicyTier0InterfaceServiceTemplate(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`

resource "nsxt_policy_tier0_gateway" "test" {
  nsx_id            = "%s"
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  type         = "SERVICE"
  mtu          = %s
  gateway_path = nsxt_policy_tier0_gateway.test.path
  segment_path = nsxt_policy_vlan_segment.test.path
  subnets      = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_interface.test.path
}`, nsxtPolicyTier0GatewayID, nsxtPolicyTier0GatewayName, name, mtu, subnet)
}

func testAccNsxtPolicyTier0InterfaceThinTemplate(name string, subnet string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  nsx_id            = "%s"
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name = "%s"
  type         = "SERVICE"
  gateway_path = nsxt_policy_tier0_gateway.test.path
  segment_path = nsxt_policy_vlan_segment.test.path
  subnets      = ["%s"]
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_interface.test.path
}`, nsxtPolicyTier0GatewayID, nsxtPolicyTier0GatewayName, name, subnet)
}

func testAccNsxtPolicyTier0InterfaceTemplateWithID(name string, subnet string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
data "nsxt_policy_ipv6_ndra_profile" "default" {
  display_name = "default"
}

resource "nsxt_policy_tier0_gateway" "test" {
  nsx_id            = "%s"
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  nsx_id                 = "test"
  display_name           = "%s"
  type                   = "SERVICE"
  description            = "Acceptance Test"
  gateway_path           = nsxt_policy_tier0_gateway.test.path
  segment_path           = nsxt_policy_vlan_segment.test.path
  subnets                = ["%s"]
  ipv6_ndra_profile_path = data.nsxt_policy_ipv6_ndra_profile.default.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_interface.test.path
}`, nsxtPolicyTier0GatewayID, nsxtPolicyTier0GatewayName, name, subnet)
}

func testAccNsxtPolicyTier0InterfaceExternalTemplate(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
data "nsxt_policy_edge_node" "EN" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

resource "nsxt_policy_tier0_gateway" "test" {
  nsx_id            = "%s"
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name   = "%s"
  description    = "Acceptance Test"
  type           = "EXTERNAL"
  mtu            = %s
  gateway_path   = nsxt_policy_tier0_gateway.test.path
  segment_path   = nsxt_policy_vlan_segment.test.path
  edge_node_path = data.nsxt_policy_edge_node.EN.path
  subnets        = ["%s"]
  enable_pim     = true
  urpf_mode      = "STRICT"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_interface.test.path
}`, nsxtPolicyTier0GatewayID, nsxtPolicyTier0GatewayName, name, mtu, subnet)
}
