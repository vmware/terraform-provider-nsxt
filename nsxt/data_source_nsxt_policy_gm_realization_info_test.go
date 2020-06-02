/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyGMRealizationInfo_serviceDataSource(t *testing.T) {
	resourceDataType := "nsxt_policy_service"
	resourceName := "DNS"
	entityType := ""
	testResourceName := "data.nsxt_policy_gm_realization_info.realization_info"
	site := getTestSiteName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_SITE_NAME")
			testAccOnlyGlobalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGMRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType, site),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttrSet(testResourceName, "entity_type"),
					resource.TestCheckResourceAttrSet(testResourceName, "realized_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
				),
			},
			{
				Config: testAccNsxtPolicyGMNoRealizationInfoTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyGMRealizationInfoReadDataSourceTemplate(resourceDataType string, resourceName string, entityType string, site string) string {
	return fmt.Sprintf(`
data "%s" "policy_resource" {
  display_name = "%s"
}

data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_gm_realization_info" "realization_info" {
  path = data.%s.policy_resource.path
  entity_type = "%s"
  site_path = data.nsxt_policy_site.test.path
}`, resourceDataType, resourceName, site, resourceDataType, entityType)
}

func testAccNsxtPolicyGMNoRealizationInfoTemplate() string {
	return fmt.Sprintf(` `)
}
