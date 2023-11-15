/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/aaa"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var accTestPolicyRoleBindingCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
}

var accTestPolicyRoleBindingUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
}

func TestAccResourceNsxtPolicyRoleBinding_basic(t *testing.T) {
	testResourceName := "nsxt_policy_user_management_role_binding.test"
	identType := nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_LDAP
	userType := nsxModel.RoleBinding_TYPE_REMOTE_USER

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_USER")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRoleBindingCheckDestroy(state, accTestPolicyRoleUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRoleBindingCreate(getTestLdapUser(), userType, identType),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRoleBindingExists(accTestPolicyRoleBindingCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRoleBindingCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRoleBindingCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "name", getTestLdapUser()),
					resource.TestCheckResourceAttr(testResourceName, "type", userType),
					resource.TestCheckResourceAttr(testResourceName, "identity_source_type", identType),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.roles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.roles.0", "auditor"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.1.path", "/orgs/default"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.1.roles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.1.roles.0", "org_admin"),

					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyRoleBindingUpdate(getTestLdapUser(), userType, identType),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyRoleBindingExists(accTestPolicyRoleBindingUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyRoleBindingUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyRoleBindingUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "name", getTestLdapUser()),
					resource.TestCheckResourceAttr(testResourceName, "type", userType),
					resource.TestCheckResourceAttr(testResourceName, "identity_source_type", identType),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.roles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.roles.0", "auditor"),

					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyRoleBinding_import_basic(t *testing.T) {
	testResourceName := "nsxt_policy_user_management_role_binding.test"
	identType := nsxModel.RoleBinding_IDENTITY_SOURCE_TYPE_LDAP
	userType := nsxModel.RoleBinding_TYPE_REMOTE_USER

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_USER")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyRoleBindingCheckDestroy(state, accTestPolicyRoleUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRoleBindingCreate(getTestLdapUser(), userType, identType),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyRoleBindingExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("RoleBinding resource %s not found in resources", resourceName)
		}

		rbID := rs.Primary.Attributes["id"]
		if rbID == "" {
			return fmt.Errorf("RoleBinding resource ID not set in resources")
		}
		rbClient := aaa.NewRoleBindingsClient(connector)
		_, err := rbClient.Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("RoleBinding %s does not exist", rbID)
			}
		}

		return err
	}
}

func testAccNsxtPolicyRoleBindingCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_user_management_role_binding" {
			continue
		}

		rbID := rs.Primary.Attributes["id"]
		if rbID == "" {
			return fmt.Errorf("RoleBinding resource ID not set in resources")
		}
		rbClient := aaa.NewRoleBindingsClient(connector)
		_, err := rbClient.Get(rbID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("RoleBinding %s still exists", displayName)
	}
	return nil
}

func testAccNsxtPolicyRoleBindingCreate(user, userType, identType string) string {
	attrMap := accTestPolicyRoleBindingCreateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_user_management_role_binding" "test" {
    display_name         = "%s"
    description          = "%s"
    name                 = "%s"
    type                 = "%s"
    identity_source_type = "%s"

    roles_for_path {
        path  = "/"
        roles = ["auditor"]
    }

    roles_for_path {
        path  = "/orgs/default"
        roles = ["org_admin"]
    }
}`, attrMap["display_name"], attrMap["description"], user, userType, identType)
}

func testAccNsxtPolicyRoleBindingUpdate(user, userType, identType string) string {
	attrMap := accTestPolicyRoleBindingUpdateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_user_management_role_binding" "test" {
    display_name         = "%s"
    description          = "%s"
    name                 = "%s"
    type                 = "%s"
    identity_source_type = "%s"

    roles_for_path {
        path  = "/"
        roles = ["auditor"]
    }
}`, attrMap["display_name"], attrMap["description"], user, userType, identType)
}
