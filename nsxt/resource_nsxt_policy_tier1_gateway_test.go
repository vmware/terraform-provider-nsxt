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

func TestAccResourceNsxtPolicyTier1Gateway_basic(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-basic")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"
	failoverMode := "NON_PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1CreateTemplate(name, failoverMode, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "true"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "false"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1UpdateTemplate(updateName, failoverMode, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "false"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1Update2Template(updateName, failoverMode, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "false"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withPoolAllocation(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-basic")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"
	failoverMode := "NON_PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1CreateTemplate(name, failoverMode, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "true"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "false"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "pool_allocation", "ROUTING"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1UpdateTemplate(updateName, failoverMode, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "false"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "pool_allocation", "ROUTING"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				// ForceNew due to pool allocation change
				Config: testAccNsxtPolicyTier1Update2Template(updateName, failoverMode, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "default_rule_logging", "false"),
					resource.TestCheckResourceAttr(testResourceName, "enable_firewall", "true"),
					resource.TestCheckResourceAttr(testResourceName, "enable_standby_relocation", "false"),
					resource.TestCheckResourceAttr(testResourceName, "force_whitelisting", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_types.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", "/infra/ipv6-ndra-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_dad_profile_path", "/infra/ipv6-dad-profiles/default"),
					resource.TestCheckResourceAttr(testResourceName, "pool_allocation", "LB_SMALL"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withDHCP(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-dhcp")
	testResourceName := "nsxt_policy_tier1_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1CreateWithDHCPTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttrSet(testResourceName, "dhcp_config_path"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1CreateWithDHCPRemovedTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config_path", ""),
				),
			},
			{
				Config: testAccNsxtPolicyTier1CreateDHCPTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withEdgeCluster(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-ec")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"
	edgeClusterName := getEdgeClusterName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1CreateWithEcTemplate(name, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1CreateWithEcTemplate(updateName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1CreateWithEcRemovedTemplate(updateName, edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "edge_cluster_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withId(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-id")
	id := fmt.Sprintf("test-id")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1SetTemplateWithID(name, id),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "id", id),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1SetTemplateWithID(updateName, id),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "id", id),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
		},
	})
}
func TestAccResourceNsxtPolicyTier1Gateway_withQos(t *testing.T) {
	name := "test-nsx-policy-tier1"
	profileName := "test-nsx-qos-profile"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccDataSourceNsxtPolicyGatewayQosProfileDeleteByName(profileName); err != nil {
				return err
			}
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGatewayQosProfileCreate(profileName); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyTier1TemplateWithQos(name, profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "ingress_qos_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "egress_qos_profile_path"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1TemplateRemoveQos(updateName, profileName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "ingress_qos_profile_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "egress_qos_profile_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withRules(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-rule")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1CreateTemplateWithRules(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.0.prefix_operator", "EQ"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1UpdateTemplateWithRules(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "route_advertisement_rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_withTier0(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-t0")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_tier1_gateway.test"
	tier0Name := "tier-0-test-for-tier1"
	failoverMode := "NON_PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			testAccDataSourceNsxtPolicyTier0DeleteByName(tier0Name)
			return testAccNsxtPolicyTier1CheckDestroy(state, name)

		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier0Create(tier0Name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyTier1CreateWithTier0Template(name, tier0Name, failoverMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1UpdateTemplate(updateName, failoverMode, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "tier0_path", ""),
					resource.TestCheckResourceAttr(realizationResourceName, "state", "REALIZED"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1Gateway_importBasic(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-tier1-import")
	testResourceName := "nsxt_policy_tier1_gateway.test"
	failoverMode := "PREEMPTIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1CheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1ImportTemplate(name, failoverMode),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyTier1Exists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := infra.NewDefaultTier1sClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier1 resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier1 resource ID not set in resources")
		}

		_, err := nsxClient.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy Tier1 ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyTier1CheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := infra.NewDefaultTier1sClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier1_gateway" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy Tier1 %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTier1CreateDHCPTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_dhcp_relay" "test" {
  display_name      = "terraform-dhcp-relay"
  server_addresses  = ["88.9.9.2"]
}
`)
}

func testAccNsxtPolicyTier1CreateWithDHCPTemplate(name string) string {
	return testAccNsxtPolicyTier1CreateDHCPTemplate() + fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  dhcp_config_path  = nsxt_policy_dhcp_relay.test.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, getEdgeClusterName(), name)
}

func testAccNsxtPolicyTier1CreateWithDHCPRemovedTemplate(name string) string {
	return testAccNsxtPolicyTier1CreateDHCPTemplate() + fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, getEdgeClusterName(), name)
}

func testAccNsxtPolicyTier1CreateTemplate(name string, failoverMode string, withPoolAllocation bool) string {
	poolAllocation := ""
	if withPoolAllocation {
		poolAllocation = `pool_allocation = "ROUTING"`
	}

	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test"
  failover_mode             = "%s"
  default_rule_logging      = "true"
  enable_firewall           = "false"
  enable_standby_relocation = "true"
  force_whitelisting        = "false"
  route_advertisement_types = ["TIER1_STATIC_ROUTES", "TIER1_CONNECTED"]
  ipv6_ndra_profile_path    = "/infra/ipv6-ndra-profiles/default"
  %s

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
  path = nsxt_policy_tier1_gateway.test.path
}`, name, failoverMode, poolAllocation)
}

func testAccNsxtPolicyTier1UpdateTemplate(name string, failoverMode string, withPoolAllocation bool) string {
	poolAllocation := ""
	if withPoolAllocation {
		poolAllocation = `pool_allocation = "ROUTING"`
	}
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test Update"
  failover_mode             = "%s"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  enable_standby_relocation = "false"
  force_whitelisting        = "true"
  route_advertisement_types = ["TIER1_CONNECTED"]
  ipv6_dad_profile_path     = "/infra/ipv6-dad-profiles/default"
  %s

  tag {
    scope = "scope3"
    tag   = "tag3"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, name, failoverMode, poolAllocation)
}

func testAccNsxtPolicyTier1Update2Template(name string, failoverMode string, withPoolAllocation bool) string {
	poolAllocation := ""
	if withPoolAllocation {
		poolAllocation = `pool_allocation = "LB_SMALL"`
	}
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test Update"
  failover_mode             = "%s"
  default_rule_logging      = "false"
  enable_firewall           = "true"
  enable_standby_relocation = "false"
  force_whitelisting        = "true"
  route_advertisement_types = ["TIER1_CONNECTED"]
  ipv6_dad_profile_path     = "/infra/ipv6-dad-profiles/default"
  %s

  tag {
    scope = "scope3"
    tag   = "tag3"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, name, failoverMode, poolAllocation)
}

func testAccNsxtPolicyTier1ImportTemplate(name string, failoverMode string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name              = "%s"
  description               = "Acceptance Test"
  failover_mode             = "%s"
  default_rule_logging      = "true"
  enable_firewall           = "false"
  enable_standby_relocation = "true"
  force_whitelisting        = "false"
  route_advertisement_types = ["TIER1_STATIC_ROUTES", "TIER1_CONNECTED"]
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
  path = nsxt_policy_tier1_gateway.test.path
}`, name, failoverMode)
}

func testAccNsxtPolicyTier1CreateWithEcTemplate(name string, edgeClusterName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
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
  path = nsxt_policy_tier1_gateway.test.path
}`, edgeClusterName, name)
}

func testAccNsxtPolicyTier1CreateWithEcRemovedTemplate(name string, edgeClusterName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  description       = "Acceptance Test"

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
  path = nsxt_policy_tier1_gateway.test.path
}`, edgeClusterName, name)
}

func testAccNsxtPolicyTier1CreateWithTier0Template(name string, tier0Name string, failoverMode string) string {
	return fmt.Sprintf(`
data "nsxt_policy_tier0_gateway" "T0" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name  = "%s"
  description   = "Acceptance Test"
  tier0_path    = data.nsxt_policy_tier0_gateway.T0.path
  failover_mode = "%s"

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
  path = nsxt_policy_tier1_gateway.test.path
}`, tier0Name, name, failoverMode)
}

func testAccNsxtPolicyTier1SetTemplateWithID(name string, id string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_tier1_gateway" "test" {
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
  path = nsxt_policy_tier1_gateway.test.path
}`, id, name)
}

func testAccNsxtPolicyTier1CreateTemplateWithRules(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  route_advertisement_rule {
    name                      = "rule1"
    action                    = "DENY"
    subnets                   = ["20.0.0.0/24", "21.0.0.0/24"]
    route_advertisement_types = ["TIER1_CONNECTED"]
    prefix_operator           = "EQ"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, name)
}

func testAccNsxtPolicyTier1UpdateTemplateWithRules(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  route_advertisement_rule {
    name                      = "rule1"
    action                    = "DENY"
    subnets                   = ["20.0.0.0/24", "21.0.0.0/24"]
    route_advertisement_types = ["TIER1_CONNECTED"]
  }

  route_advertisement_rule {
    name            = "rule2"
    action          = "PERMIT"
    subnets         = ["30.0.0.0/24", "31.0.0.0/24"]
    prefix_operator = "GE"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, name)
}

func testAccNsxtPolicyTier1TemplateWithQos(name string, profileName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name             = "%s"
  ingress_qos_profile_path = data.nsxt_policy_gateway_qos_profile.test.path
  egress_qos_profile_path  = data.nsxt_policy_gateway_qos_profile.test.path
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, profileName, name)
}

func testAccNsxtPolicyTier1TemplateRemoveQos(name string, profileName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name             = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway.test.path
}`, profileName, name)
}
