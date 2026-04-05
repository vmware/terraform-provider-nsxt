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

var accTestPolicyVpcServiceProfileCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"ntp_servers":      "5.5.5.5",
	"lease_time":       "50840",
	"dns_server_ips":   "7.7.7.7",
	"server_addresses": "11.11.11.11",
}

var accTestPolicyVpcServiceProfileUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"ntp_servers":      "5.5.5.7",
	"lease_time":       "148000",
	"dns_server_ips":   "7.7.7.2",
	"server_addresses": "11.11.11.111",
}

func TestAccResourceNsxtVpcServiceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_service_profile.test"
	testDataSourceName := "data.nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, accTestPolicyVpcServiceProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfileTemplate(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
					resource.TestCheckResourceAttr(testDataSourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.ntp_servers.0", accTestPolicyVpcServiceProfileCreateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.lease_time", accTestPolicyVpcServiceProfileCreateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfileCreateAttributes["dns_server_ips"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				// add profiles
				Config: testAccNsxtVpcServiceProfileTemplate(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
					resource.TestCheckResourceAttr(testDataSourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.ntp_servers.0", accTestPolicyVpcServiceProfileCreateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.lease_time", accTestPolicyVpcServiceProfileCreateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfileCreateAttributes["dns_server_ips"]),
					resource.TestCheckResourceAttrSet(testResourceName, "qos_profile"),
					resource.TestCheckResourceAttrSet(testResourceName, "spoof_guard_profile"),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_discovery_profile"),
					resource.TestCheckResourceAttrSet(testResourceName, "mac_discovery_profile"),
					resource.TestCheckResourceAttrSet(testResourceName, "security_profile"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				// update values, no profiles
				Config: testAccNsxtVpcServiceProfileTemplate(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.ntp_servers.0", accTestPolicyVpcServiceProfileUpdateAttributes["ntp_servers"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.lease_time", accTestPolicyVpcServiceProfileUpdateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfileUpdateAttributes["dns_server_ips"]),

					resource.TestCheckResourceAttr(testResourceName, "qos_profile", ""),
					resource.TestCheckResourceAttr(testResourceName, "spoof_guard_profile", ""),
					resource.TestCheckResourceAttr(testResourceName, "ip_discovery_profile", ""),
					resource.TestCheckResourceAttr(testResourceName, "mac_discovery_profile", ""),
					resource.TestCheckResourceAttr(testResourceName, "security_profile", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceProfileDhcpRelay(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config.0.server_addresses.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceProfileMinimalistic(accTestPolicyVpcServiceProfileCreateAttributes["display_name"]),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
				),
			},
		},
	})
}

var accTestPolicyVpcServiceProfile920Ipv6Create = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform ipv6 create",
	"v6_lease":     "86400",
	"v6_pref":      "48000",
	"v6_ntp":       "2001:db8:acc1::1",
	"v6_sntp":      "2001:db8:acc1::2",
	"v6_dns":       "2001:db8:acc1::53",
	"fwd_log":      "INFO",
	"svc_cidr":     "10.210.240.0/28",
}

var accTestPolicyVpcServiceProfile920Ipv6Update = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform ipv6 update",
	"v6_lease":     "7200",
	"v6_pref":      "5760",
	"v6_ntp":       "2001:db8:acc2::1",
	"v6_sntp":      "2001:db8:acc2::2",
	"v6_dns":       "2001:db8:acc2::53",
	"fwd_log":      "WARNING",
	"svc_cidr":     "10.210.240.16/28",
}

func TestAccResourceNsxtVpcServiceProfile920_ipv6Extensions(t *testing.T) {
	testResourceName := "nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, accTestPolicyVpcServiceProfile920Ipv6Update["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfile920Ipv6Template(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfile920Ipv6Create["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfile920Ipv6Create["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.lease_time", accTestPolicyVpcServiceProfile920Ipv6Create["v6_lease"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.preferred_time", accTestPolicyVpcServiceProfile920Ipv6Create["v6_pref"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.ntp_servers.0", accTestPolicyVpcServiceProfile920Ipv6Create["v6_ntp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.sntp_servers.0", accTestPolicyVpcServiceProfile920Ipv6Create["v6_sntp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfile920Ipv6Create["v6_dns"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.log_level", accTestPolicyVpcServiceProfile920Ipv6Create["fwd_log"]),
					resource.TestCheckResourceAttr(testResourceName, "service_subnet_cidrs.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "service_subnet_cidrs.0", accTestPolicyVpcServiceProfile920Ipv6Create["svc_cidr"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtVpcServiceProfile920Ipv6Template(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcServiceProfile920Ipv6Update["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcServiceProfile920Ipv6Update["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.lease_time", accTestPolicyVpcServiceProfile920Ipv6Update["v6_lease"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.preferred_time", accTestPolicyVpcServiceProfile920Ipv6Update["v6_pref"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.ntp_servers.0", accTestPolicyVpcServiceProfile920Ipv6Update["v6_ntp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.sntp_servers.0", accTestPolicyVpcServiceProfile920Ipv6Update["v6_sntp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcpv6_config.0.dhcpv6_server_config.0.dns_client_config.0.dns_server_ips.0", accTestPolicyVpcServiceProfile920Ipv6Update["v6_dns"]),
					resource.TestCheckResourceAttr(testResourceName, "dns_forwarder_config.0.log_level", accTestPolicyVpcServiceProfile920Ipv6Update["fwd_log"]),
					resource.TestCheckResourceAttr(testResourceName, "service_subnet_cidrs.0", accTestPolicyVpcServiceProfile920Ipv6Update["svc_cidr"]),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcServiceProfile920_ipv6Extensions_import(t *testing.T) {
	name := accTestPolicyVpcServiceProfile920Ipv6Create["display_name"]
	testResourceName := "nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfile920Ipv6Template(true),
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

func TestAccResourceNsxtVpcServiceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceProfileMinimalistic(getAccTestResourceName()),
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

func testAccNsxtVpcServiceProfileExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcServiceProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcServiceProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcServiceProfileExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcServiceProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcServiceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_service_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcServiceProfileExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcServiceProfile %s still exists", displayName)
		}
	}
	return nil
}

var testAccNsxtVpcServiceProfileHelper = getAccTestResourceName()

func testAccNsxtProjectProfiles() string {
	return fmt.Sprintf(`
resource "nsxt_policy_mac_discovery_profile" "test" {
  %s
  display_name = "%s"
}

resource "nsxt_policy_ip_discovery_profile" "test" {
  %s
  display_name = "%s"
}

resource "nsxt_policy_spoof_guard_profile" "test" {
  %s
  display_name = "%s"
}

resource "nsxt_policy_segment_security_profile" "test" {
  %s
  display_name = "%s"
}

resource "nsxt_policy_qos_profile" "test" {
  %s
  display_name = "%s"
}`, testAccNsxtProjectContext(), testAccNsxtVpcServiceProfileHelper, testAccNsxtProjectContext(), testAccNsxtVpcServiceProfileHelper, testAccNsxtProjectContext(), testAccNsxtVpcServiceProfileHelper, testAccNsxtProjectContext(), testAccNsxtVpcServiceProfileHelper, testAccNsxtProjectContext(), testAccNsxtVpcServiceProfileHelper)
}

func testAccNsxtVpcServiceProfileTemplate(createFlow bool, withProfiles bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyVpcServiceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyVpcServiceProfileUpdateAttributes
	}
	profileAssignment := ""
	if withProfiles {
		profileAssignment = `
  qos_profile           = nsxt_policy_qos_profile.test.path
  security_profile      = nsxt_policy_segment_security_profile.test.path
  ip_discovery_profile  = nsxt_policy_ip_discovery_profile.test.path
  spoof_guard_profile   = nsxt_policy_spoof_guard_profile.test.path
  mac_discovery_profile = nsxt_policy_mac_discovery_profile.test.path`
	}
	return testAccNsxtProjectProfiles() + fmt.Sprintf(`
resource "nsxt_vpc_service_profile" "test" {
  %s
  %s
  display_name = "%s"
  description  = "%s"

  dhcp_config {
    dhcp_server_config {
      ntp_servers = ["%s"]
      lease_time  = %s

      dns_client_config {
        dns_server_ips = ["%s"]
      }
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"

  depends_on = [nsxt_vpc_service_profile.test]
}`, testAccNsxtProjectContext(), profileAssignment, attrMap["display_name"], attrMap["description"], attrMap["ntp_servers"], attrMap["lease_time"], attrMap["dns_server_ips"], testAccNsxtProjectContext(), attrMap["display_name"])
}

func testAccNsxtVpcServiceProfileDhcpRelay() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"
  dhcp_config {
    dhcp_relay_config {
      server_addresses = ["%s"]
    }
  }

}`, testAccNsxtProjectContext(), accTestPolicyVpcServiceProfileUpdateAttributes["display_name"], accTestPolicyVpcServiceProfileUpdateAttributes["server_addresses"])
}

func testAccNsxtVpcServiceProfileMinimalistic(displayName string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"
  dhcp_config {
  }
}`, testAccNsxtProjectContext(), displayName)
}

func testAccNsxtVpcServiceProfile920Ipv6Template(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyVpcServiceProfile920Ipv6Create
	} else {
		attrMap = accTestPolicyVpcServiceProfile920Ipv6Update
	}
	return fmt.Sprintf(`
resource "nsxt_policy_dns_forwarder_zone" "vsp920_dns_fwd" {
%s
  display_name     = "%s-dnsfwd"
  upstream_servers = ["8.8.8.8"]
}

resource "nsxt_vpc_service_profile" "test" {
  %s
  display_name = "%s"
  description  = "%s"

  dhcp_config {
    dhcp_server_config {
      lease_time  = 86400
      ntp_servers = ["5.5.5.5"]
      dns_client_config {
        dns_server_ips = ["1.1.1.1"]
      }
    }
  }

  dhcpv6_config {
    dhcpv6_server_config {
      lease_time       = %s
      preferred_time   = %s
      ntp_servers      = ["%s"]
      sntp_servers     = ["%s"]
      dns_client_config {
        dns_server_ips = ["%s"]
      }
    }
  }

  dns_forwarder_config {
    cache_size                  = 1024
    default_forwarder_zone_path = nsxt_policy_dns_forwarder_zone.vsp920_dns_fwd.path
    log_level                   = "%s"
  }

  service_subnet_cidrs = ["%s"]
}
`, testAccNsxtProjectContext(), attrMap["display_name"], testAccNsxtProjectContext(), attrMap["display_name"], attrMap["description"],
		attrMap["v6_lease"], attrMap["v6_pref"], attrMap["v6_ntp"], attrMap["v6_sntp"], attrMap["v6_dns"],
		attrMap["fwd_log"], attrMap["svc_cidr"])
}
