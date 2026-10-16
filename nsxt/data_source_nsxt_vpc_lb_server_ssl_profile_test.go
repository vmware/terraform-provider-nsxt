// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtVpcLBServerSslProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_vpc_lb_server_ssl_profile.test"
	resourceName := "nsxt_vpc_lb_server_ssl_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				/* fetch the profile by name */
				Config: testAccNsxtVpcLBServerSslProfileDataSourceTemplate(name, `display_name = "`+name+`"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
				),
			},
			{
				/* fetch the profile by name prefix */
				Config: testAccNsxtVpcLBServerSslProfileDataSourceTemplate(name, `display_name = "`+name[:len(name)-2]+`"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrPair(testResourceName, "id", resourceName, "id"),
				),
			},
			{
				/* fetch the profile by ID */
				Config: testAccNsxtVpcLBServerSslProfileDataSourceTemplate(name, `id = nsxt_vpc_lb_server_ssl_profile.test.id`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccNsxtVpcLBServerSslProfileDataSourceTemplate(name string, lookup string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_server_ssl_profile" "test" {
%s
  display_name       = "%s"
  description        = "%s"
  cipher_group_label = "BALANCED"
}

data "nsxt_vpc_lb_server_ssl_profile" "test" {
%s
  %s

  depends_on = [nsxt_vpc_lb_server_ssl_profile.test]
}`, testAccNsxtPolicyMultitenancyContext(), name, name, testAccNsxtPolicyMultitenancyContext(), lookup)
}
