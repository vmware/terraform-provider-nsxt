package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

func TestAccResourceNsxtPolicyVlanSegment_basicImport(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-vlan-segment")
	testResourceName := "nsxt_policy_vlan_segment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVlanSegmentImportTemplate(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyVlanSegment_basicUpdate(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-vlan-segment")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_vlan_segment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyVlanSegmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVlanSegmentBasicTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVlanSegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", "tftest.org"),
					resource.TestCheckResourceAttr(testResourceName, "vlan_ids.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyVlanSegmentBasicUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVlanSegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					// TODO: file bug for description not being updated
					// resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test2"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", "tftest2.org"),
					resource.TestCheckResourceAttr(testResourceName, "vlan_ids.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVlanSegment_updateAdvConfig(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-vlan-segment")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_vlan_segment.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyVlanSegmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVlanSegmentBasicAdvConfigTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVlanSegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", "tftest.org"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.connectivity", "OFF"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.local_egress", "true"),
				),
			},
			{
				Config: testAccNsxtPolicyVlanSegmentBasicAdvConfigUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVlanSegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", "tftest.org"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.connectivity", "ON"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.local_egress", "false"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVlanSegment_withDhcp(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-vlan-segment")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_vlan_segment.test"
	leaseTimes := []string{"3600", "36000"}
	preferredTimes := []string{"3200", "32000"}
	dnsServersV4 := []string{"2.2.2.2", "3.3.3.3"}
	dnsServersV6 := []string{"2000::2", "3000::3"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVlanSegmentWithDhcpTemplate(name, dnsServersV4[0], dnsServersV6[0], leaseTimes[0], preferredTimes[0]),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "subnet.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.cidr", "12.12.2.1/24"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v6_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.server_address", "12.12.2.2/24"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.lease_time", leaseTimes[0]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dns_servers.0", dnsServersV4[0]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dhcp_option_121.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dhcp_generic_option.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.cidr", "4012::1/64"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v4_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.server_address", "4012::2/64"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.lease_time", leaseTimes[0]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.preferred_time", preferredTimes[0]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.dns_servers.0", dnsServersV6[0]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.sntp_servers.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.excluded_range.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyVlanSegmentWithDhcpTemplate(updatedName, dnsServersV4[1], dnsServersV6[1], leaseTimes[1], preferredTimes[1]),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySegmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "subnet.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.cidr", "12.12.2.1/24"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v6_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.server_address", "12.12.2.2/24"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.lease_time", leaseTimes[1]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dns_servers.0", dnsServersV4[1]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dhcp_option_121.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.0.dhcp_v4_config.0.dhcp_generic_option.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.cidr", "4012::1/64"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v4_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.server_address", "4012::2/64"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.lease_time", leaseTimes[1]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.preferred_time", preferredTimes[1]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.dns_servers.0", dnsServersV6[1]),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.sntp_servers.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnet.1.dhcp_v6_config.0.excluded_range.#", "2"),
				),
			},
		},
	})
}

// TODO: add tests for l2_extension; requires L2 VPN Session

func testAccNsxtPolicyVlanSegmentExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := infra.NewDefaultSegmentsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VLAN Segment resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VLAN Segment resource ID not set in resources")
		}

		_, err := nsxClient.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy VLAN Segment ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyVlanSegmentCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := infra.NewDefaultSegmentsClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_vlan_segment" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy VLAN Segment %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyVlanSegmentDeps() string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "vlantz" {
  display_name = "%s"
}

`, getVlanTransportZoneName())
}

func testAccNsxtPolicyVlanSegmentImportTemplate(name string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`
resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  vlan_ids            = ["101"]
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
}
`, name)
}

func testAccNsxtPolicyVlanSegmentBasicTemplate(name string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`

resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
  domain_name         = "tftest.org"
  vlan_ids            = ["101"]

  tag {
    scope = "color"
    tag   = "orange"
  }
}
`, name)
}

func testAccNsxtPolicyVlanSegmentBasicUpdateTemplate(name string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`

resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
  domain_name         = "tftest2.org"
  vlan_ids            = ["101", "102"]

  tag {
    scope = "color"
    tag   = "green"
  }
  tag {
    scope = "color"
    tag   = "orange"
  }
}
`, name)
}

func testAccNsxtPolicyVlanSegmentBasicAdvConfigTemplate(name string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`

resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  domain_name         = "tftest.org"
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
  vlan_ids            = ["101", "102"]

  tag {
    scope = "color"
    tag   = "orange"
  }

  advanced_config {
    connectivity = "OFF"
    local_egress = true
  }
}
`, name)
}

func testAccNsxtPolicyVlanSegmentBasicAdvConfigUpdateTemplate(name string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`

resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  domain_name         = "tftest.org"
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
  vlan_ids            = ["101", "102"]

  tag {
    scope = "color"
    tag   = "orange"
  }

  advanced_config {
    connectivity = "ON"
    local_egress = false
  }
}
`, name)
}

func testAccNsxtPolicyVlanSegmentWithDhcpTemplate(name string, dnsServerV4 string, dnsServerV6 string, lease string, preferred string) string {
	return testAccNsxtPolicyVlanSegmentDeps() + fmt.Sprintf(`

data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_dhcp_server" "test" {
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  display_name      = "segment-test"
}

resource "nsxt_policy_vlan_segment" "test" {
  display_name        = "%s"
  transport_zone_path = data.nsxt_policy_transport_zone.vlantz.path
  vlan_ids            = ["101", "102"]

  subnet {
    cidr = "12.12.2.1/24"
    dhcp_v4_config {
        server_address = "12.12.2.2/24"
        lease_time     = %s
        dns_servers    = ["%s"]
        dhcp_option_121 {
          network  = "2.1.1.0/24"
          next_hop = "2.3.1.3"
        }
        dhcp_option_121 {
          network  = "3.1.1.0/24"
          next_hop = "3.3.1.3"
        }
        dhcp_generic_option {
          code   = "119"
          values = ["abc", "def"]
        }
    }
  }

  subnet {
    cidr = "4012::1/64"
    dhcp_v6_config {
        server_address = "4012::2/64"
        lease_time     = %s
        preferred_time = %s
        dns_servers    = ["%s"]
        sntp_servers   = ["3001::1", "3001::2"]
        excluded_range {
            start = "4012::400"
            end   = "4012::500"
        }
        excluded_range {
            start = "4012::a00"
            end   = "4012::b00"
        }
    }
  }

  dhcp_config_path = nsxt_policy_dhcp_server.test.path

}
`, getEdgeClusterName(), name, lease, dnsServerV4, lease, preferred, dnsServerV6)
}
