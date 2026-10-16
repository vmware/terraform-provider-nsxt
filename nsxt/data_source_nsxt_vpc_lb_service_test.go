// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtVpcLBService_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_vpc_lb_service.test"
	resourceName := "nsxt_vpc_lb_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				/* fetch the service by name */
				Config: testAccNsxtVpcLBServiceDataSourceTemplate(name, `display_name = "`+name+`"`),
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
				/* fetch the service by ID */
				Config: testAccNsxtVpcLBServiceDataSourceTemplate(name, `id = nsxt_vpc_lb_service.test.id`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrPair(testResourceName, "id", resourceName, "id"),
				),
			},
		},
	})
}

func testAccNsxtVpcLBServiceDataSourceTemplate(name string, lookup string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_service" "test" {
%s
  display_name = "%s"
  description  = "%s"
}

data "nsxt_vpc_lb_service" "test" {
%s
  %s

  depends_on = [nsxt_vpc_lb_service.test]
}`, testAccNsxtPolicyMultitenancyContext(), name, name, testAccNsxtPolicyMultitenancyContext(), lookup)
}
