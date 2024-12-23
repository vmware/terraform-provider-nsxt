/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceNsxtPolicyGroups_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyGroupsBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func TestAccDataSourceNsxtPolicyGroups_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyGroupsBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicyGroupsBasic(t *testing.T, withContext bool, preCheck func()) {
	domain := "default"
	groupName := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_Groups.test"
	checkResourceName := "data.nsxt_policy_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyGroupDeleteByName(domain, groupName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGroupCreate(domain, groupName); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNSXPolicyGroupsReadTemplate(groupName, withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(checkResourceName, "display_name", groupName),
				),
			},
		},
	})
}

func testAccNSXPolicyGroupsReadTemplate(groupName string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "nsxt_policy_groups" "test" {
%s
}

locals {
  group_from_list = data.nsxt_policy_groups.test.items[index(data.nsxt_policy_groups.test.items.*.display_name, "%s")]
}

data "nsxt_policy_group" "test" {
    id = local.group_from_list.id
}
`, context, groupName)
}
