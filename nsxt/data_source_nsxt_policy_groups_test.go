// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccDataSourceNsxtPolicyGroups_includeSharedGroups(t *testing.T) {
	projectName := getAccTestResourceName()
	groupName := getAccTestDataSourceName()
	shareName := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_groups.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
			testAccNSXVersion(t, "4.1.1")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyGroupsReadSharedTemplate(projectName, groupName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "items."+groupName, fmt.Sprintf("/infra/domains/default/groups/%s", groupName)),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyGroupsBasic(t *testing.T, withContext bool, preCheck func()) {
	domain := "default"
	groupName := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_groups.test"
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
  // Get id from path
  path_split = split("/", data.nsxt_policy_groups.test.items["%s"])
  group_id   = element(local.path_split, length(local.path_split) - 1)
}

data "nsxt_policy_group" "test" {
%s
  id = local.group_id
}
`, context, groupName, context)
}

func testAccNSXPolicyGroupsReadSharedTemplate(projectName, groupName, shareName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test" {
  display_name = "%s"
}

resource "nsxt_policy_group" "group_default" {
  display_name = "%s"
}

resource "nsxt_policy_share" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  sharing_strategy = "ALL_DESCENDANTS"
  shared_with      = [nsxt_policy_project.test.path]
}

resource "nsxt_policy_shared_resource" "test" {
  display_name = "shared_group_resource"
  share_path   = nsxt_policy_share.test.path
  resource_object {
    resource_path    = nsxt_policy_group.group_default.path
    include_children = false
  }
}

data "nsxt_policy_groups" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  include_shared_groups = true
  depends_on            = [nsxt_policy_shared_resource.test]
}
`, projectName, groupName, shareName)
}
