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

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXGlobalManagerSitePrecheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTransportZoneReadTemplate(transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", transportZoneName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_default"),
					resource.TestCheckResourceAttr(testResourceName, "transport_type", "VLAN_BACKED"),
				),
			},
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

func testAccNSXPolicyTransportZoneReadTemplate(transportZoneName string) string {
	if testAccIsGlobalManager() {
		return testAccNSXGlobalPolicyTransportZoneReadTemplate(transportZoneName)
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

func testAccNSXGlobalPolicyTransportZoneReadTemplate(transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
  site_path = data.nsxt_policy_site.test.path
}`, getTestSiteName(), transportZoneName)
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
