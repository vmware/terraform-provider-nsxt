// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var accTestVpcSubnetCreateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"access_mode":        "Isolated",
	"enabled":            "true",
	"ipv4_subnet_size":   "32",
	"ip_addresses":       "192.168.240.0/26",
	"reserved_ip_ranges": "11.11.50.7",
}

var accTestVpcSubnetUpdateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform updated",
	"access_mode":        "Isolated",
	"enabled":            "true",
	"ipv4_subnet_size":   "32",
	"ip_addresses":       "192.168.240.0/26",
	"reserved_ip_ranges": "11.11.50.4-11.11.50.6",
}

var accTestVpcSubnetReservedIPRangeAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"access_mode":        "Isolated",
	"enabled":            "true",
	"ipv4_subnet_size":   "32",
	"ip_addresses":       "10.110.1.0/26",
	"reserved_ip_ranges": "10.110.1.3-10.110.1.15",
}

func TestAccResourceNsxtVpcSubnet_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersionLessThan(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetCreateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", "DHCP_DEACTIVATED"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetUpdateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", "DHCP_DEACTIVATED"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetMinimalisticIsolated(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcSubnet920_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"
	vlanExtResourcesname := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetCreateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", "DHCP_DEACTIVATED"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestVpcSubnetUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestVpcSubnetUpdateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", "DHCP_DEACTIVATED"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetMinimalisticIsolated(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetVLANExtensionTemplate(vlanExtResourcesname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "vlan_connection"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", vlanExtResourcesname+"-vpc-subnet"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcSubnet_subnetSize(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersionLessThan(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetSizeTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestVpcSubnetCreateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetSizeTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestVpcSubnetCreateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetMinimalisticPublic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcSubnet920_subnetSize(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetSizeTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestVpcSubnetCreateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetSizeTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestVpcSubnetCreateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetMinimalisticPublic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcSubnet_reservedIPRange(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersionLessThan(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetReservedIPRangeTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetReservedIPRangeAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetReservedIPRangeAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.reserved_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.reserved_ip_ranges.0", accTestVpcSubnetReservedIPRangeAttributes["reserved_ip_ranges"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcSubnet920_reservedIPRange(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetReservedIPRangeTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetReservedIPRangeAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcSubnetReservedIPRangeAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.reserved_ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.reserved_ip_ranges.0", accTestVpcSubnetReservedIPRangeAttributes["reserved_ip_ranges"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_server_additional_config.0.options.0.other.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

var accTestVpcSubnet920SubnetDhcpv6Create = map[string]string{
	"display_name": getAccTestResourceName(),
}

var accTestVpcSubnet920SubnetDhcpv6Update = map[string]string{
	"display_name": getAccTestResourceName(),
}

var accTestVpcSubnetDhcpv4ConfigCreate = map[string]string{
	"display_name": getAccTestResourceName(),
}

var accTestVpcSubnetDhcpv4ConfigUpdate = map[string]string{
	"display_name": getAccTestResourceName(),
}

func TestAccResourceNsxtVpcSubnet_subnetDhcpv4Config(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetDhcpv4ConfigUpdate["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpv4Template(true, "DHCP_DEACTIVATED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpv4ConfigCreate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpv4ConfigCreate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpv4Template(false, "DHCP_DEACTIVATED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpv4ConfigUpdate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpv4ConfigUpdate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetSubnetDhcpv4Template(false, "DHCP_DEACTIVATED"),
				PlanOnly: true,
			},
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpv4Template(false, "DHCP_SERVER"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpv4ConfigUpdate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpv4ConfigUpdate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_SERVER),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.0", "192.168.240.11/26"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetSubnetDhcpv4Template(false, "DHCP_SERVER"),
				PlanOnly: true,
			},
		},
	})
}

var accTestVpcSubnetDhcpDualStackConfigCreate = map[string]string{
	"display_name": getAccTestResourceName(),
}

var accTestVpcSubnetDhcpDualStackConfigUpdate = map[string]string{
	"display_name": getAccTestResourceName(),
}

func TestAccResourceNsxtVpcSubnet920_subnetDhcpDualStackConfig(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetDhcpDualStackConfigUpdate["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(true, "DHCP_DEACTIVATED", "DHCP_DEACTIVATED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpDualStackConfigCreate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpDualStackConfigCreate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(false, "DHCP_DEACTIVATED", "DHCP_DEACTIVATED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpDualStackConfigUpdate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpDualStackConfigUpdate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(false, "DHCP_DEACTIVATED", "DHCP_DEACTIVATED"),
				PlanOnly: true,
			},
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(false, "DHCP_SERVER", "DHCP_DEACTIVATED"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpDualStackConfigUpdate["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpDualStackConfigUpdate["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_SERVER),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.0", "192.168.240.11/26"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(false, "DHCP_SERVER", "DHCP_DEACTIVATED"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceNsxtVpcSubnet920_subnetDhcpv6Config(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnet920SubnetDhcpv6Update["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpv6Template(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnet920SubnetDhcpv6Create["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnet920SubnetDhcpv6Create["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.dns_server_preference", model.SubnetDhcpv6Config_DNS_SERVER_PREFERENCE_PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetSubnetDhcpv6Template(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnet920SubnetDhcpv6Update["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnet920SubnetDhcpv6Update["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.dns_server_preference", model.SubnetDhcpv6Config_DNS_SERVER_PREFERENCE_DNS_FORWARDER_PREFERRED_OVER_PROFILE_DNS_SERVERS),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetSubnetDhcpv6Template(false),
				PlanOnly: true,
			},
		},
	})
}

var accTestVpcSubnet920DnsPrefCreateAttrs = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"ipv4_subnet_size": "32",
}

var accTestVpcSubnet920DnsPrefUpdateAttrs = map[string]string{
	"display_name": getAccTestResourceName(),
}

func TestAccResourceNsxtVpcSubnet920_dnsServerPreference(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnet920DnsPrefUpdateAttrs["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetDnsServerPreferenceTemplate(model.SubnetDhcpConfig_DNS_SERVER_PREFERENCE_DNS_FORWARDER_PREFERRED_OVER_PROFILE_DNS_SERVERS),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnet920DnsPrefCreateAttrs["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_SERVER),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_server_preference", model.SubnetDhcpConfig_DNS_SERVER_PREFERENCE_DNS_FORWARDER_PREFERRED_OVER_PROFILE_DNS_SERVERS),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtVpcSubnetDnsServerPreferenceTemplate(model.SubnetDhcpConfig_DNS_SERVER_PREFERENCE_PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnet920DnsPrefCreateAttrs["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_server_preference", model.SubnetDhcpConfig_DNS_SERVER_PREFERENCE_PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER),
				),
			},
			{
				Config: testAccNsxtVpcSubnetDnsServerPreferenceMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnet920DnsPrefUpdateAttrs["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcSubnet_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetMinimalisticIsolated(),
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

func testAccNsxtVpcSubnetExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VpcSubnet resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VpcSubnet resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcSubnetExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VpcSubnet %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcSubnetCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_subnet" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcSubnetExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("VpcSubnet %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcSubnetTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetCreateAttributes
	} else {
		attrMap = accTestVpcSubnetUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  description  = "%s"

  ip_addresses = ["%s"]
  access_mode  = "%s"
  dhcp_config {
    mode = "DHCP_DEACTIVATED"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["ip_addresses"], attrMap["access_mode"])
}

func testAccNsxtVpcSubnetVLANExtensionTemplate(vlanExtResourcesname string) string {
	template := `
resource "nsxt_policy_project" "vpc_project" {
  nsx_id       = "${NAME}-proj"
  display_name = "${NAME}-proj"
  description  = "Vlan extension Project"
  tgw_external_connections = [nsxt_policy_distributed_vlan_connection.test.path]
}

resource "nsxt_vpc" "test" {
  context {
    project_id = nsxt_policy_project.vpc_project.id
  }

  display_name = "${NAME}-vpc"
  description  = "Terraform provisioned VPC"

}

resource "nsxt_policy_distributed_vlan_connection" "test" {
  display_name      = "${NAME}-dvc"
  vlan_id           = 888
  subnet_extension_connection = "ENABLED_L2"
}

resource "nsxt_vpc_subnet" "test" {
  context {
    project_id = nsxt_policy_project.vpc_project.id
    vpc_id = nsxt_vpc.test.id
  }
  display_name = "${NAME}-vpc-subnet"

  ip_addresses = ["192.168.240.0/26"]
  access_mode  = "L2_Only"
  vlan_connection = nsxt_policy_distributed_vlan_connection.test.path
}`
	return strings.ReplaceAll(template, "${NAME}", vlanExtResourcesname)
}

func testAccNsxtVpcSubnetSizeTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetCreateAttributes
	} else {
		attrMap = accTestVpcSubnetUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  description  = "%s"

  ipv4_subnet_size = %s
  access_mode      = "Public"

  dhcp_config {
    mode = "DHCP_SERVER"

    dhcp_server_additional_config {
      options {
        option121 {
          static_route {
            network  = "2.1.1.0/24"
            next_hop = "2.3.1.3"
          }
        }
        other {
          code   = "119"
          values = ["abc", "def"]
        }
      }
    }
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["ipv4_subnet_size"])
}

func testAccNsxtVpcSubnetMinimalisticIsolated() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  ip_addresses = ["%s"]
  access_mode  = "Isolated"
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcSubnetUpdateAttributes["display_name"], accTestVpcSubnetUpdateAttributes["ip_addresses"])
}

func testAccNsxtVpcSubnetMinimalisticPublic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  access_mode  = "Public"
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcSubnetUpdateAttributes["display_name"])
}

func testAccNsxtVpcSubnetDnsServerPreferenceTemplate(preference string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  description  = "%s"

  ipv4_subnet_size = %s
  access_mode      = "Public"

  dhcp_config {
    mode                  = "%s"
    dns_server_preference = "%s"
  }
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcSubnet920DnsPrefCreateAttrs["display_name"],
		accTestVpcSubnet920DnsPrefCreateAttrs["description"], accTestVpcSubnet920DnsPrefCreateAttrs["ipv4_subnet_size"],
		model.SubnetDhcpConfig_MODE_SERVER, preference)
}

func testAccNsxtVpcSubnetDnsServerPreferenceMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  access_mode  = "Public"
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcSubnet920DnsPrefUpdateAttrs["display_name"])
}

func testAccNsxtVpcSubnetReservedIPRangeTemplate() string {
	attrMap := accTestVpcSubnetReservedIPRangeAttributes
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  description  = "%s"
  ip_addresses = ["%s"]

  access_mode = "Public"

  dhcp_config {
    mode = "DHCP_SERVER"

    dhcp_server_additional_config {
      options {
        option121 {
          static_route {
            network  = "2.1.1.0/24"
            next_hop = "2.3.1.3"
          }
        }
        other {
          code   = "119"
          values = ["abc", "def"]
        }
      }
      reserved_ip_ranges = ["%s"]
    }
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["ip_addresses"], attrMap["reserved_ip_ranges"])
}

func testAccNsxtVpcSubnetSubnetDhcpv6Template(createFlow bool) string {
	var attrMap map[string]string
	var mode, dnsPref string
	if createFlow {
		attrMap = accTestVpcSubnet920SubnetDhcpv6Create
		mode = model.SubnetDhcpv6Config_MODE_DEACTIVATED
		dnsPref = model.SubnetDhcpv6Config_DNS_SERVER_PREFERENCE_PROFILE_DNS_SERVERS_PREFERRED_OVER_DNS_FORWARDER
	} else {
		attrMap = accTestVpcSubnet920SubnetDhcpv6Update
		mode = model.SubnetDhcpv6Config_MODE_DEACTIVATED
		dnsPref = model.SubnetDhcpv6Config_DNS_SERVER_PREFERENCE_DNS_FORWARDER_PREFERRED_OVER_PROFILE_DNS_SERVERS
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  ip_addresses = ["2001:db8:f100::/64"]
  ip_address_type = "IPV6"
  access_mode  = "Isolated"

  subnet_dhcpv6_config {
    mode                  = "%s"
    dns_server_preference = "%s"
  }

  advanced_config {
    connectivity_state    = "DISCONNECTED"
    gateway_addresses     = ["2001:db8:f100::10/64"]
    dhcp_server_addresses = ["2001:db8:f100::11/64"]
  }
}
`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], mode, dnsPref)
}

func testAccNsxtVpcSubnetSubnetDhcpv4Template(createFlow bool, dhcpMode string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetDhcpv4ConfigCreate
	} else {
		attrMap = accTestVpcSubnetDhcpv4ConfigUpdate
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  ip_addresses = ["192.168.240.0/26"]
  access_mode  = "Isolated"

  dhcp_config {
    mode = "%s"
  }

  advanced_config {
    connectivity_state    = "DISCONNECTED"
    gateway_addresses     = ["192.168.240.10/26"]
    dhcp_server_addresses = ["192.168.240.11/26"]
  }
}
`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], dhcpMode)
}

func testAccNsxtVpcSubnetSubnetDhcpDualStackTemplate(createFlow bool, v4Mode string, v6Mode string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcSubnetDhcpDualStackConfigCreate
	} else {
		attrMap = accTestVpcSubnetDhcpDualStackConfigUpdate
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name = "%s"
  ip_addresses = ["192.168.240.0/26", "2001:db8:f100::/64"]
  ip_address_type = "IPV4_IPV6"
  access_mode  = "Isolated"

  dhcp_config {
    mode = "%s"
  }

  subnet_dhcpv6_config {
    mode = "%s"
  }

  advanced_config {
    connectivity_state    = "DISCONNECTED"
    gateway_addresses     = ["192.168.240.10/26", "2001:db8:f100::10/64"]
    dhcp_server_addresses = ["192.168.240.11/26", "2001:db8:f100::11/64"]
  }
}
`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], v4Mode, v6Mode)
}

var accTestVpcSubnetDhcpv4DeactivatedServerAddress = map[string]string{
	"display_name": getAccTestResourceName(),
}

var accTestVpcSubnetDhcpv6DeactivatedServerAddress = map[string]string{
	"display_name": getAccTestResourceName(),
}

// TestAccResourceNsxtVpcSubnet920_createWithDeactivatedIPv6ServerAddress covers the
// regression from bug 3719347 update #11: creating a dual-stack subnet where DHCPv4 is
// deactivated but an IPv4 entry is still present in advanced_config.dhcp_server_addresses
// alongside an active IPv6 entry used to fail on create with an NSX ip-cidr-block format
// validation error on an empty dhcp_server_addresses entry. The plan-time DiffSuppressFunc
// masked the IPv4 slot as unchanged (since DHCPv4 is deactivated), but
// resourceNsxtVpcSubnetCreate then sent that suppressed slot to NSX as an empty string instead
// of omitting it, and NSX rejected the empty string.
func TestAccResourceNsxtVpcSubnet920_createWithDeactivatedIPv6ServerAddress(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
			t.Skip("DHCPv6 disabled on VPC service profile")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetDhcpv4DeactivatedServerAddress["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetDhcpv4DeactivatedServerAddressTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpv4DeactivatedServerAddress["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpv4DeactivatedServerAddress["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_SERVER),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.0", "2001:db8:f200::11/64"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetDhcpv4DeactivatedServerAddressTemplate(),
				PlanOnly: true,
			},
		},
	})
}

// TestAccResourceNsxtVpcSubnet920_createWithDeactivatedIPv4ServerAddress covers
// the same scenario as TestAccResourceNsxtVpcSubnet920_createWithDeactivatedIPv6ServerAddress
// with the only difference that DHCPv6 is deactivated, and DHCPv4 is active
func TestAccResourceNsxtVpcSubnet920_createWithDeactivatedIPv4ServerAddress(t *testing.T) {
	testResourceName := "nsxt_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcSubnetCheckDestroy(state, accTestVpcSubnetDhcpv6DeactivatedServerAddress["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcSubnetDhcpv6DeactivatedServerAddressTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcSubnetExists(accTestVpcSubnetDhcpv6DeactivatedServerAddress["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcSubnetDhcpv6DeactivatedServerAddress["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.mode", model.SubnetDhcpConfig_MODE_SERVER),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnet_dhcpv6_config.0.mode", model.SubnetDhcpv6Config_MODE_DEACTIVATED),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.dhcp_server_addresses.0", "192.168.249.66/26"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config:   testAccNsxtVpcSubnetDhcpv6DeactivatedServerAddressTemplate(),
				PlanOnly: true,
			},
		},
	})
}

func testAccNsxtVpcSubnetDhcpv4DeactivatedServerAddressTemplate() string {
	return testAccNsxtVpcSubnetDeactivatedFamilyServerAddressTemplate(
		accTestVpcSubnetDhcpv4DeactivatedServerAddress, model.SubnetDhcpConfig_MODE_DEACTIVATED, model.SubnetDhcpv6Config_MODE_SERVER)
}

func testAccNsxtVpcSubnetDhcpv6DeactivatedServerAddressTemplate() string {
	return testAccNsxtVpcSubnetDeactivatedFamilyServerAddressTemplate(
		accTestVpcSubnetDhcpv6DeactivatedServerAddress, model.SubnetDhcpConfig_MODE_SERVER, model.SubnetDhcpv6Config_MODE_DEACTIVATED)
}

func testAccNsxtVpcSubnetDeactivatedFamilyServerAddressTemplate(attrMap map[string]string, v4Mode string, v6Mode string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "test" {
%s
  display_name    = "%s"
  ip_addresses    = ["192.168.249.64/26", "2001:db8:f200::/64"]
  ip_address_type = "IPV4_IPV6"
  access_mode     = "Isolated"

  dhcp_config {
    mode = "%s"
  }

  subnet_dhcpv6_config {
    mode = "%s"
  }

  advanced_config {
    connectivity_state    = "DISCONNECTED"
    gateway_addresses     = ["192.168.249.65/26", "2001:db8:f200::10/64"]
    dhcp_server_addresses = ["192.168.249.66/26", "2001:db8:f200::11/64"]
  }
}
`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], v4Mode, v6Mode)
}
