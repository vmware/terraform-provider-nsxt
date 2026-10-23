// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtVpcLBPersistenceProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	cookieName := name + "-cookie"
	genericName := name + "-generic"
	sourceIPName := name + "-srcip"
	testResourceName := "data.nsxt_vpc_lb_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				/* fetch by full name */
				Config: testAccNsxtVpcLBPersistenceProfileDataSourceTemplate(name, fmt.Sprintf(`display_name = "%s"`, cookieName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", cookieName),
					resource.TestCheckResourceAttr(testResourceName, "description", cookieName),
					resource.TestCheckResourceAttr(testResourceName, "type", "COOKIE"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", "nsxt_vpc_lb_cookie_persistence_profile.test", "id"),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
				),
			},
			{
				/* fetch by name prefix and type */
				Config: testAccNsxtVpcLBPersistenceProfileDataSourceTemplate(name, fmt.Sprintf(`display_name = "%s"
  type         = "GENERIC"`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", genericName),
					resource.TestCheckResourceAttr(testResourceName, "type", "GENERIC"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", "nsxt_vpc_lb_generic_persistence_profile.test", "id"),
				),
			},
			{
				/* fetch by name prefix and type, source IP */
				Config: testAccNsxtVpcLBPersistenceProfileDataSourceTemplate(name, fmt.Sprintf(`display_name = "%s"
  type         = "SOURCE_IP"`, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", sourceIPName),
					resource.TestCheckResourceAttr(testResourceName, "type", "SOURCE_IP"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", "nsxt_vpc_lb_source_ip_persistence_profile.test", "id"),
				),
			},
			{
				/* fetch by ID */
				Config: testAccNsxtVpcLBPersistenceProfileDataSourceTemplate(name, `id = nsxt_vpc_lb_cookie_persistence_profile.test.id`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", cookieName),
					resource.TestCheckResourceAttr(testResourceName, "type", "COOKIE"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		}),
	})
}

func testAccNsxtVpcLBPersistenceProfileDataSourceTemplate(name string, lookup string) string {
	context := testAccNsxtPolicyMultitenancyContext()
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_cookie_persistence_profile" "test" {
%[1]s
  display_name = "%[2]s-cookie"
  description  = "%[2]s-cookie"
}

resource "nsxt_vpc_lb_generic_persistence_profile" "test" {
%[1]s
  display_name = "%[2]s-generic"
  description  = "%[2]s-generic"
}

resource "nsxt_vpc_lb_source_ip_persistence_profile" "test" {
%[1]s
  display_name = "%[2]s-srcip"
  description  = "%[2]s-srcip"
}

data "nsxt_vpc_lb_persistence_profile" "test" {
%[1]s
  %[3]s

  depends_on = [
    nsxt_vpc_lb_cookie_persistence_profile.test,
    nsxt_vpc_lb_generic_persistence_profile.test,
    nsxt_vpc_lb_source_ip_persistence_profile.test,
  ]
}`, context, name, lookup)
}
