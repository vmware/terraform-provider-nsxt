/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

func TestAccResourceNsxtPolicyTier0Gateway_basic(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-basic")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier0_gateway.test"
	failoverMode := "NON_PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0CreateTemplate(name, failoverMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "false"),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", "ACTIVE_STANDBY"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "false"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0UpdateTemplate(updateName, failoverMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					// TODO: file a bug for description updates not sticking
					//resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "true"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ha_mode", "ACTIVE_ACTIVE"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_withId(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-id")
	id := fmt.Sprintf("test-id")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier0_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0SetTemplateWithID(name, id),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "id", id),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0SetTemplateWithID(updateName, id),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "id", id),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_withSubnets(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-subnets")
	testResourceName := "nsxt_policy_tier0_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0SubnetsTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "internal_transit_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_withDHCP(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-dhcp")
	testResourceName := "nsxt_policy_tier0_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0CreateWithDHCPTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "internal_transit_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "dhcp_config_path"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0CreateWithDHCPRemovedTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "internal_transit_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config_path", ""),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_withEdgeCluster(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-ec")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier0_gateway.test"
	edgeClusterName := getEdgeClusterName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0CreateWithEcTemplate(name, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.path"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.inter_sr_ibgp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", "65000"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.multipath_relax", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.route_aggregation.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_timer", "180"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_stale_route_timer", "600"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0UpdateWithEcTemplate(updateName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.path"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.inter_sr_ibgp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", "60000"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.multipath_relax", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.route_aggregation.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_timer", "180"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_stale_route_timer", "600"),
				),
			},
			/* TODO: enable when 2472726 is resolved
			{
				Config: testAccNsxtPolicyTier0CreateWithEcRemovedTemplate(updateName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "edge_cluster_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.path"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.inter_sr_ibgp", "false"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", "65000"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.multipath_relax", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.route_aggregation.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_timer", "180"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_stale_route_timer", "600"),
				),
			},
			*/
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_createWithBGP(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-bgp")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier0_gateway.test"
	edgeClusterName := getEdgeClusterName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0UpdateWithEcTemplate(updateName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "bgp_config.0.path"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.ecmp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.inter_sr_ibgp", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.local_as_num", "60000"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.multipath_relax", "true"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.route_aggregation.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_mode", "HELPER_ONLY"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_timer", "180"),
					resource.TestCheckResourceAttr(testResourceName, "bgp_config.0.graceful_restart_stale_route_timer", "600"),
				),
			},
		},
	})
}

// TODO: add route_distinguisher when VNI pool DS is exposed
func TestAccResourceNsxtPolicyTier0Gateway_withVRF(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-vrf")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier0_gateway.test"
	testInterfaceName := "nsxt_policy_tier0_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0WithVRFTemplate(name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "vrf_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "vrf_config.0.gateway_path"),
					resource.TestCheckResourceAttr(testResourceName, "vrf_config.0.route_target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "vrf_config.0.route_target.0.address_family"),
					resource.TestCheckResourceAttr(testResourceName, "vrf_config.0.route_target.0.export_targets.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "vrf_config.0.route_target.0.import_targets.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testInterfaceName, "display_name", name),
					resource.TestCheckResourceAttr(testInterfaceName, "access_vlan_id", "12"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0WithVRFTemplate(updateName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "vrf_config.0.route_target.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "vrf_config.0.gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testInterfaceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testInterfaceName, "access_vlan_id", "12"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0WithVRFTearDown(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier0Gateway_importBasic(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier0-import")
	testResourceName := "nsxt_policy_tier0_gateway.test"
	failoverMode := "PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0CreateTemplate(name, failoverMode),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyTier0Exists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := infra.NewDefaultTier0sClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier0 resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier0 resource ID not set in resources")
		}

		_, err := nsxClient.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy Tier0 ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyTier0CheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := infra.NewDefaultTier0sClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier0_gateway" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy Tier0 %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTier0CreateWithEcTemplate(name string, edgeClusterName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  description       = "Acceptance Test"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, edgeClusterName, name)
}

func testAccNsxtPolicyTier0UpdateWithEcTemplate(name string, edgeClusterName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  description       = "Acceptance Test"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path

  bgp_config {
    local_as_num    = "60000"
    multipath_relax = true
    route_aggregation {
      prefix = "12.12.12.0/24"
    }
    route_aggregation {
      prefix = "13.12.12.0/24"
    }

    tag {
      scope = "bgp-scope"
      tag   = "bgp-tag"
    }
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

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}

data "nsxt_policy_realization_info" "bgp_realization_info" {
  path = nsxt_policy_tier0_gateway.test.bgp_config.0.path
}`, edgeClusterName, name)
}

func testAccNsxtPolicyTier0CreateWithDHCPTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_dhcp_relay" "test" {
  display_name      = "terraform-dhcp-relay"
  server_addresses  = ["88.9.9.2"]
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  dhcp_config_path  = nsxt_policy_dhcp_relay.test.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, getEdgeClusterName(), name)
}

func testAccNsxtPolicyTier0CreateWithDHCPRemovedTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_dhcp_relay" "test" {
  display_name      = "terraform-dhcp-relay"
  server_addresses  = ["88.9.9.2"]
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, getEdgeClusterName(), name)
}

func testAccNsxtPolicyTier0CreateTemplate(name string, failoverMode string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test"
  failover_mode             = "%s"
  default_rule_logging      = "true"
  enable_firewall           = "false"
  force_whitelisting        = "false"
  ha_mode                   = "ACTIVE_STANDBY"
  ipv6_ndra_profile_path    = "/infra/ipv6-ndra-profiles/default"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }

}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, name, failoverMode)
}

func testAccNsxtPolicyTier0UpdateTemplate(name string, failoverMode string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test"
  failover_mode             = "%s"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  force_whitelisting        = "true"
  ha_mode                   = "ACTIVE_ACTIVE"
  ipv6_dad_profile_path     = "/infra/ipv6-dad-profiles/default"

  tag {
    scope = "scope3"
    tag   = "tag3"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, name, failoverMode)
}

func testAccNsxtPolicyTier0SetTemplateWithID(name string, id string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_tier0_gateway" "test" {
  nsx_id       = "%s"
  display_name = "%s"
  description  = "Acceptance Test"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, id, name)
}

func testAccNsxtPolicyTier0SubnetsTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test"
  failover_mode             = "NON_PREEMPTIVE"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  force_whitelisting        = "true"
  ha_mode                   = "ACTIVE_STANDBY"
  ipv6_dad_profile_path     = "/infra/ipv6-dad-profiles/default"
  internal_transit_subnets  = ["102.64.0.0/16"]
  transit_subnets           = ["101.64.0.0/16"]

  tag {
    scope = "scope3"
    tag   = "tag3"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, name)
}

// TODO: add vrf_config tags when bug 2557096 is resolved
func testAccNsxtPolicyTier0WithVRFTemplate(name string, targets bool) string {

	var routeTargets string
	if targets {
		routeTargets = `
        route_target {
            auto_mode      = "false"
            import_targets = ["2:12"]
            export_targets = ["8999:123", "2:14"]
        }
            `
	}
	return testAccNsxtPolicyGatewayInterfaceDeps("11, 12") + fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "parent" {
  nsx_id            = "vrf-parent"
  display_name      = "parent"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  vrf_config {
    gateway_path = nsxt_policy_tier0_gateway.parent.path
    %s
  }
}

resource "nsxt_policy_tier0_gateway_interface" "parent-loopback" {
  display_name   = "parent interface"
  type           = "LOOPBACK"
  gateway_path   = nsxt_policy_tier0_gateway.parent.path
  edge_node_path = data.nsxt_policy_edge_node.EN.path
  subnets        = ["4.4.4.12/24"]
}

data "nsxt_policy_edge_node" "EN" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name   = "%s"
  type           = "EXTERNAL"
  gateway_path   = nsxt_policy_tier0_gateway.test.path
  segment_path   = nsxt_policy_vlan_segment.test.path
  edge_node_path = data.nsxt_policy_edge_node.EN.path
  subnets        = ["4.4.4.1/24"]
  access_vlan_id = 12
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway.test.path
}`, name, routeTargets, name)
}

func testAccNsxtPolicyTier0WithVRFTearDown() string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11, 12") + `
data "nsxt_policy_edge_node" "EN" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  member_index      = 0
}

resource "nsxt_policy_tier0_gateway" "parent" {
  nsx_id            = "vrf-parent"
  display_name      = "parent"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_policy_tier0_gateway_interface" "parent-loopback" {
  display_name   = "parent interface"
  type           = "LOOPBACK"
  gateway_path   = nsxt_policy_tier0_gateway.parent.path
  edge_node_path = data.nsxt_policy_edge_node.EN.path
  subnets        = ["4.4.4.12/24"]
}`
}
