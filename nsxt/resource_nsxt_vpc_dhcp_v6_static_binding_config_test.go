// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"mac_address":     "10:0e:00:11:22:02",
	"lease_time":      "7200",
	"preferred_time":  "5760",
	"ip_addresses":    "fd00:240:2400::41",
	"domain_names":    "client-create.example.org",
	"dns_nameservers": "fd00:240:2400::2",
	"ntp_servers":     "ntp-create.example.org",
	"sntp_servers":    "fd00:240:2400::3",
}

var accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"mac_address":     "10:ff:22:11:cc:02",
	"lease_time":      "10800",
	"preferred_time":  "8640",
	"ip_addresses":    "fd00:240:2400::42",
	"domain_names":    "client-update.example.org",
	"dns_nameservers": "fd00:240:2400::4",
	"ntp_servers":     "ntp-update.example.org",
	"sntp_servers":    "fd00:240:2400::5",
}

func TestAccResourceNsxtVpcSubnetDhcpV6StaticBindingConfig920_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_dhcp_v6_static_binding.test"
	subnetResourceName := "nsxt_vpc_subnet.test"
	// One stable name for the VPC fixture across all steps; NSX rejects changing vpc short_id after create.
	fixtureBase := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetDhcpV6StaticBindingConfigCheckDestroy(state, accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetDhcpV6StaticBindingConfigTemplate(true, fixtureBase),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetDhcpV6StaticBindingConfigExists(accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(subnetResourceName, "ip_address_type", "IPV4_IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["display_name"]),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_address", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["mac_address"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "preferred_time", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["preferred_time"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "domain_names.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_names.0", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["domain_names"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.0", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["dns_nameservers"]),
					resource.TestCheckResourceAttr(testResourceName, "ntp_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ntp_servers.0", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "sntp_servers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "sntp_servers.0", accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["sntp_servers"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetDhcpV6StaticBindingConfigTemplate(false, fixtureBase),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetDhcpV6StaticBindingConfigExists(accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(subnetResourceName, "ip_address_type", "IPV4_IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_address", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["mac_address"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "preferred_time", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["preferred_time"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "domain_names.0", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["domain_names"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.0", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["dns_nameservers"]),
					resource.TestCheckResourceAttr(testResourceName, "ntp_servers.0", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "sntp_servers.0", accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["sntp_servers"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetDhcpV6StaticBindingConfigMinimalistic(fixtureBase),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetDhcpV6StaticBindingConfigExists(accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(subnetResourceName, "ip_address_type", "IPV4_IPV6"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcSubnetDhcpV6StaticBindingConfig920_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	fixtureBase := getAccTestResourceName()
	testResourceName := "nsxt_vpc_dhcp_v6_static_binding.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetDhcpV6StaticBindingConfigCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetDhcpV6StaticBindingConfigMinimalisticImport(name, fixtureBase),
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

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("DhcpV6StaticBindingConfig resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("DhcpV6StaticBindingConfig resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]

		exists, err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigExists(getSessionContextFromParentPath(testAccProvider.Meta(), parentPath), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("DhcpV6StaticBindingConfig %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_dhcp_v6_static_binding" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]

		exists, err := resourceNsxtVpcSubnetDhcpV6StaticBindingConfigExists(getSessionContextFromParentPath(testAccProvider.Meta(), parentPath), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("DhcpV6StaticBindingConfig %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigPrerequisites(base string) string {
	// CI default VPC service profile may omit dhcpv6_server_config; NSX 641028 requires it when the
	// subnet uses DHCPv6 DHCP_SERVER. Provision a dedicated VPC + profile with dhcpv6_server_config.
	shortID := base
	if len(shortID) > 8 {
		shortID = shortID[len(shortID)-8:]
	}
	projCtx := testAccNsxtMultitenancyContext(false)
	vpcProjectID := os.Getenv("NSXT_VPC_PROJECT_ID")

	return fmt.Sprintf(`
data "nsxt_policy_transit_gateway" "dhcpv6_bind_tgw" {
%s
  is_default = true
}

resource "nsxt_vpc_service_profile" "dhcpv6_bind_sp" {
%s
  display_name = "%s-sp"
  dhcp_config {
    dhcp_server_config {
      lease_time = 86400
    }
  }
  dhcpv6_config {
    dhcpv6_server_config {
      lease_time     = 86400
      preferred_time = 3600
    }
  }
}

resource "nsxt_vpc_connectivity_profile" "dhcpv6_bind_cp" {
%s
  display_name         = "%s-cp"
  transit_gateway_path = data.nsxt_policy_transit_gateway.dhcpv6_bind_tgw.path
  service_gateway {
    enable = false
  }
}

resource "nsxt_vpc" "dhcpv6_bind_vpc" {
%s
  display_name        = "%s-vpc"
  description         = "tf acc dhcpv6 static binding vpc"
  private_ips         = ["192.168.246.0/24"]
  short_id            = "%s"
  vpc_service_profile = nsxt_vpc_service_profile.dhcpv6_bind_sp.path
  load_balancer_vpc_endpoint {
    enabled = false
  }
}

resource "nsxt_vpc_attachment" "dhcpv6_bind_att" {
  display_name             = "%s-att"
  parent_path              = nsxt_vpc.dhcpv6_bind_vpc.path
  vpc_connectivity_profile = nsxt_vpc_connectivity_profile.dhcpv6_bind_cp.path
}

resource "nsxt_vpc_subnet" "test" {
  context {
    project_id = "%s"
    vpc_id     = nsxt_vpc.dhcpv6_bind_vpc.id
  }
  display_name = "test-subnet-dhcpv6-acc"
  depends_on   = [nsxt_vpc_attachment.dhcpv6_bind_att]

  # Dual-stack: both CIDRs in ip_addresses. NSX 9.2 allows one IPv4 and one IPv6 across ip_addresses
  # and gateway_addresses combined (610716); do not repeat IPv4 in gateway_addresses.
  ip_address_type = "IPV4_IPV6"
  ip_addresses    = ["10.203.240.0/26", "fd00:240:2400::/64"]
  advanced_config {
    connectivity_state = "DISCONNECTED"
    # NSX fills gateway and DHCP server addresses from the subnet CIDRs; omitting them causes
    # post-apply plan drift (non-empty plan after refresh) in acceptance tests.
    gateway_addresses     = ["10.203.240.1/26", "fd00:240:2400::1/64"]
    dhcp_server_addresses = ["10.203.240.2/26", "fd00:240:2400:0:0:0:0:2/64"]
  }
  access_mode = "Isolated"
  dhcp_config {
    mode = "DHCP_SERVER"
    dhcp_server_additional_config {
      reserved_ip_ranges = ["10.203.240.40-10.203.240.60"]
    }
  }
  subnet_dhcpv6_config {
    mode = "DHCP_SERVER"
    dhcpv6_server_additional_config {
      domain_names       = []
      reserved_ip_ranges = ["fd00:240:2400::40-fd00:240:2400::5f"]
    }
  }
}
`, projCtx, projCtx, base, projCtx, base, projCtx, base, shortID, base, vpcProjectID)
}

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigTemplate(createFlow bool, fixtureBase string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes
	} else {
		attrMap = accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes
	}
	return testAccNsxtVpcSubnetDhcpV6StaticBindingConfigPrerequisites(fixtureBase) + fmt.Sprintf(`
resource "nsxt_vpc_dhcp_v6_static_binding" "test" {
  parent_path       = nsxt_vpc_subnet.test.path
  display_name      = "%s"
  description       = "%s"
  mac_address       = "%s"
  lease_time        = %s
  preferred_time    = %s
  ip_addresses      = ["%s"]
  domain_names      = ["%s"]
  dns_nameservers   = ["%s"]
  ntp_servers       = ["%s"]
  sntp_servers      = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["mac_address"],
		attrMap["lease_time"], attrMap["preferred_time"], attrMap["ip_addresses"],
		attrMap["domain_names"], attrMap["dns_nameservers"], attrMap["ntp_servers"], attrMap["sntp_servers"])
}

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigMinimalistic(fixtureBase string) string {
	attrMap := accTestVpcSubnetDhcpV6StaticBindingConfigCreateAttributes
	return testAccNsxtVpcSubnetDhcpV6StaticBindingConfigPrerequisites(fixtureBase) + fmt.Sprintf(`
resource "nsxt_vpc_dhcp_v6_static_binding" "test" {
  display_name = "%s"
  parent_path  = nsxt_vpc_subnet.test.path
  ip_addresses = ["%s"]
  mac_address  = "%s"
}`, accTestVpcSubnetDhcpV6StaticBindingConfigUpdateAttributes["display_name"], attrMap["ip_addresses"], attrMap["mac_address"])
}

func testAccNsxtVpcSubnetDhcpV6StaticBindingConfigMinimalisticImport(displayName string, fixtureBase string) string {
	return testAccNsxtVpcSubnetDhcpV6StaticBindingConfigPrerequisites(fixtureBase) + fmt.Sprintf(`
resource "nsxt_vpc_dhcp_v6_static_binding" "test" {
  display_name = "%s"
  parent_path  = nsxt_vpc_subnet.test.path
  ip_addresses = ["fd00:240:2400::41"]
  mac_address  = "10:0e:00:11:22:03"
}`, displayName)
}
