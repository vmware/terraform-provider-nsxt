/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
)

func TestAccResourceNsxtPolicyService_icmp(t *testing.T) {
	name := "test-nsx-policy-icmp-type-service-basic"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				// Step 0: Single ICMP entry
				Config: testAccNsxtPolicyIcmpTypeServiceCreateTypeCodeTemplate(name, "3", "1", "ICMPv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", "3"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 1: Update service & entry
				Config: testAccNsxtPolicyIcmpTypeServiceCreateTypeOnlyTemplate(updateName, "34", "ICMPv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", "34"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 2: Update service & entry for ICMPv6
				Config: testAccNsxtPolicyIcmpTypeServiceCreateTypeOnlyTemplate(name, "3", "ICMPv6"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", "3"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv6"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 3: Update service & entry with no type & code
				Config: testAccNsxtPolicyIcmpTypeServiceCreateNoTypeCodeTemplate(name, "ICMPv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 4: Add another ICMP service entry
				Config: testAccNsxtPolicyIcmpTypeServiceCreate2Template(name, "3", "1", "ICMPv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Updated Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.1.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.1.description", "Updated Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.1.icmp_type", "3"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.1.icmp_code", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.1.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", "no-type-code"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Updated Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_icmpNoEntryDisplayName(t *testing.T) {
	name := "test-nsx-policy-icmp-type-service-no-display-name"
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIcmpTypeServiceCreateNoDisplayNameTemplate(name, "3", "1", "ICMPv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", ""),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_type", "3"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.icmp_code", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_l4PortSet(t *testing.T) {
	name := "test-nsx-policy-l4-port-set-type-service-basic"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				// Step 0: Single TCP L4 port set entry
				Config: testAccNsxtPolicyL4PortSetTypeServiceCreateTemplate(name, "TCP", "99-101"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "l4 service"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.destination_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.source_ports.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 1: Single UDP L4 port set entry
				Config: testAccNsxtPolicyL4PortSetTypeServiceCreateTemplate(name, "UDP", "88"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "l4 service"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.protocol", "UDP"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.destination_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.source_ports.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 2: 2 service Entries
				Config: testAccNsxtPolicyL4PortSetTypeServiceUpdateTemplate(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "updated l4 service"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.display_name", "entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.1.display_name", "entry-2"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.1.description", "Entry-2"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.1.protocol", "UDP"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 2: Update back to only one entry
				Config: testAccNsxtPolicyL4PortSetTypeServiceCreateTemplate(name, "UDP", "88"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "l4 service"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.description", "Entry"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.protocol", "UDP"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.destination_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.source_ports.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_mixedServices(t *testing.T) {
	name := "test-nsx-policy-mixed-service"
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				// Step 0: mixed service entries
				Config: testAccNsxtPolicyMixedServiceCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.display_name", "entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.display_name", "entry-2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.description", "Entry-2"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.0.protocol", "ICMPv4"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				// Step 1: Single UDP L4 port set entry
				Config: testAccNsxtPolicyL4PortSetTypeServiceCreateTemplate(name, "UDP", "88"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_igmp(t *testing.T) {
	name := "test-nsx-policy-igmp-type-service"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIgmpTypeServiceCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "IGMP entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyIgmpTypeServiceCreateTemplate(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "IGMP entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.0.display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_etherType(t *testing.T) {
	name := "test-nsx-policy-ether-type-service"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEtherTypeServiceCreateTemplate(name, "1536"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Ether type entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.ether_type", "1536"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyEtherTypeServiceCreateTemplate(updateName, "2001"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Ether type entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.0.ether_type", "2001"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_ipProtocolType(t *testing.T) {
	name := "test-nsx-policy-ip-protocol-type-service"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPProtocolTypeServiceCreateTemplate(name, "6"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "IP Protocol type entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.protocol", "6"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyIPProtocolTypeServiceCreateTemplate(updateName, "17"),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "IP Protocol type entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.0.protocol", "17"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_algType(t *testing.T) {
	name := "test-nsx-policy-alg-service"
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_service.test"
	alg := "SUN_RPC_UDP"
	destPort := "210"
	updatedDestPort := "211"
	sourcePorts := "9000-9001"
	updatedSourcePorts := "9000-9001\", \"500-504"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyAlgServiceCreateTemplate(name, alg, sourcePorts, destPort),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Algorithm entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.algorithm", alg),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.destination_port", destPort),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.source_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyAlgServiceCreateTemplate(updateName, alg, updatedSourcePorts, updatedDestPort),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Algorithm entry"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.description", "Entry-1"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.algorithm", alg),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.source_ports.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "algorithm_entry.0.destination_port", updatedDestPort),
					resource.TestCheckResourceAttr(testResourceName, "l4_port_set_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "icmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "igmp_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ether_type_entry.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "ip_protocol_entry.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyService_importBasic(t *testing.T) {
	name := "test-nsx-policy-service-import"
	testResourceName := "nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyServiceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIcmpTypeServiceCreateNoTypeCodeTemplate(name, "ICMPv4"),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyServiceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy service resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy service resource ID not set in resources")
		}

		var err error
		if testAccIsGlobalManager() {
			nsxClient := gm_infra.NewDefaultServicesClient(connector)
			_, err = nsxClient.Get(resourceID)
		} else {
			nsxClient := infra.NewDefaultServicesClient(connector)
			_, err = nsxClient.Get(resourceID)
		}
		if err != nil {
			return fmt.Errorf("Error while retrieving policy service ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyServiceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	var err error
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		if testAccIsGlobalManager() {
			nsxClient := gm_infra.NewDefaultServicesClient(connector)
			_, err = nsxClient.Get(resourceID)
		} else {
			nsxClient := infra.NewDefaultServicesClient(connector)
			_, err = nsxClient.Get(resourceID)
		}
		if err == nil {
			return fmt.Errorf("Policy service %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIcmpTypeServiceCreateTypeCodeTemplate(name string, icmpType string, icmpCode string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  icmp_entry {
	display_name = "%s"
	description  = "Entry"
	icmp_type    = "%s"
	icmp_code    = "%s"
	protocol     = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, name, name, icmpType, icmpCode, protocol)
}

func testAccNsxtPolicyIcmpTypeServiceCreateTypeOnlyTemplate(name string, icmpType string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  icmp_entry {
	display_name = "%s"
	description  = "Entry"
	icmp_type    = "%s"
	protocol     = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, name, name, icmpType, protocol)
}

func testAccNsxtPolicyIcmpTypeServiceCreateNoTypeCodeTemplate(name string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  icmp_entry {
    display_name = "%s"
    description  = "Entry"
    protocol     = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, name, name, protocol)
}

func testAccNsxtPolicyIcmpTypeServiceCreateNoDisplayNameTemplate(name string, icmpType string, icmpCode string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  icmp_entry {
	description  = "Entry"
	icmp_type    = "%s"
	icmp_code    = "%s"
	protocol     = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, name, icmpType, icmpCode, protocol)
}

func testAccNsxtPolicyIcmpTypeServiceCreate2Template(name string, icmpType string, icmpCode string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  display_name = "%s"
  description  = "Updated Acceptance Test"

  icmp_entry {
    display_name = "%s"
    description  = "Updated Entry"
    icmp_type    = "%s"
    icmp_code    = "%s"
    protocol     = "%s"
  }

  icmp_entry {
    display_name = "no-type-code"
    description  = "Updated Entry"
    protocol     = "ICMPv4"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }

}`, name, name, icmpType, icmpCode, protocol)
}

func testAccNsxtPolicyL4PortSetTypeServiceCreateTemplate(serviceName string, protocol string, port string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description       = "l4 service"
  display_name      = "%s"

  l4_port_set_entry {
    display_name      = "%s"
    description       = "Entry"
    protocol          = "%s"
    destination_ports = [ "%s" ]
    source_ports      = [ "100", "200-300" ]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, serviceName, protocol, port)
}

func testAccNsxtPolicyL4PortSetTypeServiceUpdateTemplate(serviceName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description       = "updated l4 service"
  display_name      = "%s"

  l4_port_set_entry {
    display_name      = "entry-1"
    description       = "Entry-1"
    protocol          = "TCP"
    destination_ports = [ "100" ]
  }

  l4_port_set_entry {
    display_name      = "entry-2"
    description       = "Entry-2"
    protocol          = "UDP"
    destination_ports = [ "101" ]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, serviceName)
}

func testAccNsxtPolicyMixedServiceCreateTemplate(serviceName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description       = "mixed services"
  display_name      = "%s"

  l4_port_set_entry {
    display_name      = "entry-1"
    description       = "Entry-1"
    protocol          = "TCP"
    destination_ports = [ "80" ]
  }

  icmp_entry {
    display_name = "entry-2"
    description  = "Entry-2"
    protocol     = "ICMPv4"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName)
}

func testAccNsxtPolicyIgmpTypeServiceCreateTemplate(serviceName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description  = "IGMP entry"
  display_name = "%s"

  igmp_entry {
    display_name = "%s"
    description  = "Entry-1"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, serviceName)
}

func testAccNsxtPolicyEtherTypeServiceCreateTemplate(serviceName string, etherType string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description  = "Ether type entry"
  display_name = "%s"

  ether_type_entry {
    display_name = "%s"
    description  = "Entry-1"
    ether_type   = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, serviceName, etherType)
}

func testAccNsxtPolicyIPProtocolTypeServiceCreateTemplate(serviceName string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description  = "IP Protocol type entry"
  display_name = "%s"

  ip_protocol_entry {
    display_name = "%s"
    description  = "Entry-1"
    protocol     = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, serviceName, protocol)
}

func testAccNsxtPolicyAlgServiceCreateTemplate(serviceName string, alg string, sourcePorts string, destPort string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_service" "test" {
  description  = "Algorithm entry"
  display_name = "%s"

  algorithm_entry {
    display_name      = "%s"
    description       = "Entry-1"
    algorithm         = "%s"
    source_ports      = ["%s"]
    destination_port  = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, serviceName, serviceName, alg, sourcePorts, destPort)
}
