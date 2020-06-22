/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyTransportZone_basic(t *testing.T) {
	transportZoneName := getVlanTransportZoneName()
	testResourceName := "data.nsxt_policy_transport_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXGlobalManagerSitePrecheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTransportZoneReadTemplate(transportZoneName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_default"),
					resource.TestCheckResourceAttr(testResourceName, "transport_type", "VLAN_BACKED"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyTransportZone_withTransportType(t *testing.T) {
	transportZoneName := getVlanTransportZoneName()
	testResourceName := "data.nsxt_policy_transport_zone.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTransportZoneWithTransportTypeTemplate(transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", transportZoneName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_default"),
					resource.TestCheckResourceAttr(testResourceName, "transport_type", "VLAN_BACKED"),
				),
			},
		},
	})
}

func testAccNSXPolicyTransportZoneReadTemplate(transportZoneName string, isVlan bool) string {
	if testAccIsGlobalManager() {
		return testAccNSXGlobalPolicyTransportZoneReadTemplate(isVlan)
	}
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}`, transportZoneName)
}

func testAccNSXPolicyTransportZoneWithTransportTypeTemplate(transportZoneName string) string {
	if testAccIsGlobalManager() {
		return testAccNSXGlobalPolicyTransportZoneWithTransportTypeTemplate(transportZoneName)
	}
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name   = "%s"
  transport_type = "VLAN_BACKED"
}`, transportZoneName)
}

func testAccNSXGlobalPolicyTransportZoneReadTemplate(isVlan bool) string {
	transportType := "OVERLAY_STANDARD"
	if isVlan {
		transportType = "VLAN_BACKED"
	}
	return fmt.Sprintf(`
data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  transport_type = "%s"
  site_path      = data.nsxt_policy_site.test.path
  is_default     = true
}`, getTestSiteName(), transportType)
}

func testAccNSXGlobalPolicyTransportZoneWithTransportTypeTemplate(transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
  site_path = data.nsxt_policy_site.test.path
  transport_type = "VLAN_BACKED"
}`, getTestSiteName(), transportZoneName)
}
