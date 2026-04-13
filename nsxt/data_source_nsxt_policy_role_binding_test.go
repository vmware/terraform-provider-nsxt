// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicyRoleBinding_basic(t *testing.T) {
	testDataSourceNameByID := "data.nsxt_policy_user_management_role_binding.by_id"
	testDataSourceNameByKey := "data.nsxt_policy_user_management_role_binding.by_name_type"
	identType := model.RoleBinding_IDENTITY_SOURCE_TYPE_LDAP
	userType := model.RoleBinding_TYPE_REMOTE_USER
	ldapUser := getTestLdapUser()
	// Use same display_name as the role binding resource (from shared attributes), not getAccTestResourceName() which returns a new value each call
	expectedDisplayName := accTestPolicyRoleBindingCreateAttributes["display_name"]

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_PASSWORD")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_URL")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_DOMAIN")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_BASE_DN")
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDeleteAllLdapIdentitySources()
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRoleBindingDataSourceTemplate(ldapUser, userType, identType),
				Check: resource.ComposeTestCheckFunc(
					// Lookup by ID
					resource.TestCheckResourceAttrSet(testDataSourceNameByID, "id"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "name", ldapUser),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "type", userType),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "identity_source_type", identType),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "display_name", expectedDisplayName),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.#", "2"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.0.path", "/"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.0.roles.#", "1"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.0.roles.0", "network_engineer"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.1.path", "/orgs/default"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.1.roles.#", "1"),
					resource.TestCheckResourceAttr(testDataSourceNameByID, "roles_for_path.1.roles.0", "org_admin"),
					resource.TestCheckResourceAttrSet(testDataSourceNameByID, "revision"),

					// Lookup by name and type - same binding
					resource.TestCheckResourceAttrSet(testDataSourceNameByKey, "id"),
					resource.TestCheckResourceAttr(testDataSourceNameByKey, "name", ldapUser),
					resource.TestCheckResourceAttr(testDataSourceNameByKey, "type", userType),
					resource.TestCheckResourceAttr(testDataSourceNameByKey, "identity_source_type", identType),
					resource.TestCheckResourceAttr(testDataSourceNameByKey, "display_name", expectedDisplayName),
					resource.TestCheckResourceAttr(testDataSourceNameByKey, "roles_for_path.#", "2"),
					resource.TestCheckResourceAttrSet(testDataSourceNameByKey, "revision"),

					// Both data sources must resolve to the same id
					resource.TestCheckResourceAttrPair(testDataSourceNameByID, "id", testDataSourceNameByKey, "id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRoleBindingDataSourceTemplate(user, userType, identType string) string {
	return testAccNsxtPolicyRoleBindingLdapCreate(user, userType, identType) + `
data "nsxt_policy_user_management_role_binding" "by_id" {
  id = nsxt_policy_user_management_role_binding.test.id
}

data "nsxt_policy_user_management_role_binding" "by_name_type" {
  name                 = nsxt_policy_user_management_role_binding.test.name
  type                 = nsxt_policy_user_management_role_binding.test.type
  identity_source_id   = nsxt_policy_user_management_role_binding.test.identity_source_id
}
`
}
