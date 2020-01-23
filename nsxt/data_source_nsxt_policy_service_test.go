/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyService_basic(t *testing.T) {
	serviceName := "DNS"
	testResourceName := "data.nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyServiceReadTemplate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", serviceName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoServiceTemplate(),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyService_byId(t *testing.T) {
	serviceID := "DNS"
	testResourceName := "data.nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyServiceReadIDTemplate(serviceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceID),
					resource.TestCheckResourceAttr(testResourceName, "description", serviceID),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoServiceTemplate(),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyService_byPrefix(t *testing.T) {
	serviceName := "Heartbeat"
	servicePrefix := "Heart"
	testResourceName := "data.nsxt_policy_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyServiceReadTemplate(servicePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", serviceName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoServiceTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyServiceReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_service" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtPolicyServiceReadIDTemplate(id string) string {
	return fmt.Sprintf(`
data "nsxt_policy_service" "test" {
  id = "%s"
}`, id)
}

func testAccNsxtPolicyNoServiceTemplate() string {
	return fmt.Sprintf(` `)
}
