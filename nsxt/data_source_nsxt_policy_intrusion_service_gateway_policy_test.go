// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_byCategory(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_category"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByCategory(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByName(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_name" {
  display_name = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.display_name
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByCategory(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_category" {
  display_name = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.display_name
  category     = "LocalGatewayRules"
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByID(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_id" {
  id         = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.id
  depends_on = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, name)
}
