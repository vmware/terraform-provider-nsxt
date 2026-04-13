// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/terraform-provider-nsxt/api/aaa"
	"github.com/vmware/terraform-provider-nsxt/api/orgs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Acceptance tests in this file configure LDAP identity sources via the NSX API in PreConfig.
// They must use resource.Test (not resource.ParallelTest), for the same reason as tests that
// declare nsxt_policy_ldap_identity_source: only one LDAP identity source configuration can be
// applied to the manager at a time.

var testAccRoleBindingPathLdapID string
var testAccRoleBindingPathProjectID string

// testAccRoleBindingPathInitTestIDs must run before resource.Test builds TestStep Config strings.
// Config is evaluated when the Steps slice is constructed (before PreConfig), so IDs cannot be set only in PreConfig.
func testAccRoleBindingPathInitTestIDs(t *testing.T) {
	t.Helper()
	testAccRoleBindingPathLdapID = getAccTestResourceName()
	testAccRoleBindingPathProjectID = getAccTestResourceName()
}

func TestAccResourceNsxtPolicyUserManagementRoleBindingPath_basic(t *testing.T) {
	testAccRoleBindingPathInitTestIDs(t)
	testResourceName := "nsxt_policy_user_management_role_binding_path.project_scope"

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
			if err := testAccRoleBindingPathPreConfigDestroy(); err != nil {
				return err
			}
			return testAccDeleteAllLdapIdentitySources()
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccRoleBindingPathPreConfigCreate(); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtPolicyUserManagementRoleBindingPathTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "role_binding_id"),
					resource.TestCheckResourceAttrPair(testResourceName, "path", "nsxt_policy_project.test_project", "path"),
					resource.TestCheckResourceAttr(testResourceName, "roles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles.0", "project_admin"),
				),
			},
			{
				Config: testAccNsxtPolicyUserManagementRoleBindingPathTemplateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(testResourceName, "path", "nsxt_policy_project.test_project", "path"),
					resource.TestCheckResourceAttr(testResourceName, "roles.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyUserManagementRoleBindingPath_import(t *testing.T) {
	testAccRoleBindingPathInitTestIDs(t)
	testResourceName := "nsxt_policy_user_management_role_binding_path.project_scope"

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
			if err := testAccRoleBindingPathPreConfigDestroy(); err != nil {
				return err
			}
			return testAccDeleteAllLdapIdentitySources()
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccRoleBindingPathPreConfigCreate(); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtPolicyUserManagementRoleBindingPathTemplate(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccRoleBindingPathPreConfigCreate creates the LDAP identity source and role binding via API.
// Project id is set so the role binding uses the same path as the project that Terraform will create.
func testAccRoleBindingPathPreConfigCreate() error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("failed to get policy connector: %w", err)
	}
	sessionContext := utl.SessionContext{ClientType: utl.Local}

	// Create LDAP identity source (testAccRoleBindingPathLdapID / testAccRoleBindingPathProjectID set in testAccRoleBindingPathInitTestIDs)
	ldapClient := aaa.NewLdapIdentitySourcesClient(sessionContext, connector)
	if ldapClient == nil {
		return fmt.Errorf("LDAP identity sources client is nil")
	}
	description := "terraform test"
	revision := int64(0)
	domainName := getTestLdapDomain()
	baseDn := getTestLdapBaseDN()
	bindIdentity := getTestLdapAdminUser()
	password := getTestLdapAdminPassword()
	url := getTestLdapURL()
	enabled := true
	useStarttls := false
	ldapServers := []model.IdentitySourceLdapServer{
		{
			BindIdentity: &bindIdentity,
			Password:     &password,
			Url:          &url,
			Enabled:      &enabled,
			UseStarttls:  &useStarttls,
		},
	}
	obj := model.OpenLdapIdentitySource{
		Description:  &description,
		Revision:     &revision,
		DomainName:   &domainName,
		BaseDn:       &baseDn,
		LdapServers:  ldapServers,
		ResourceType: model.LdapIdentitySource_RESOURCE_TYPE_OPENLDAPIDENTITYSOURCE,
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(obj, model.OpenLdapIdentitySourceBindingType())
	if errs != nil {
		return errs[0]
	}
	structValue := dataValue.(*data.StructValue)
	if _, err := ldapClient.Probeidentitysource(structValue); err != nil {
		return fmt.Errorf("LDAP probe failed: %w", err)
	}
	if _, err := ldapClient.Update(testAccRoleBindingPathLdapID, structValue); err != nil {
		return fmt.Errorf("failed to create LDAP identity source: %w", err)
	}

	// Create role binding
	rbClient := aaa.NewRoleBindingsClient(sessionContext, connector)
	if rbClient == nil {
		return fmt.Errorf("role bindings client is nil")
	}
	displayName := "test-rb-path"
	rbDescription := "terraform test"
	username := getTestLdapUser()
	userType := model.RoleBinding_TYPE_REMOTE_USER
	identType := model.RoleBinding_IDENTITY_SOURCE_TYPE_LDAP
	projectPath := "/orgs/default/projects/" + testAccRoleBindingPathProjectID
	boolTrue := true
	rbObj := model.RoleBinding{
		DisplayName:        &displayName,
		Description:        &rbDescription,
		Name:               &username,
		Type_:              &userType,
		IdentitySourceId:   &testAccRoleBindingPathLdapID,
		IdentitySourceType: &identType,
		ReadRolesForPaths:  &boolTrue,
		RolesForPaths: []model.RolesForPath{
			{Path: strPtr("/"), Roles: []model.Role{{Role: strPtr("network_engineer")}}},
			{Path: strPtr("/orgs/default"), Roles: []model.Role{{Role: strPtr("org_admin")}}},
			{Path: &projectPath, Roles: []model.Role{{Role: strPtr("project_admin")}}},
		},
	}
	if _, err := rbClient.Create(rbObj); err != nil {
		return fmt.Errorf("failed to create role binding: %w", err)
	}
	return nil
}

// testAccRoleBindingPathPreConfigDestroy removes the project, LDAP identity source and role binding created in PreConfig.
func testAccRoleBindingPathPreConfigDestroy() error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("failed to get policy connector: %w", err)
	}
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	rbClient := aaa.NewRoleBindingsClient(sessionContext, connector)
	if rbClient == nil {
		return fmt.Errorf("role bindings client is nil")
	}
	username := getTestLdapUser()
	userType := model.RoleBinding_TYPE_REMOTE_USER
	rbList, err := rbClient.List(nil, nil, nil, nil, &username, nil, nil, nil, nil, nil, nil, &userType)
	if err != nil {
		return fmt.Errorf("failed to list role bindings: %w", err)
	}
	for i := range rbList.Results {
		rb := &rbList.Results[i]
		if rb.Id != nil && rb.Name != nil && *rb.Name == username && rb.Type_ != nil && *rb.Type_ == userType {
			if err := rbClient.Delete(*rb.Id, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil); err != nil {
				return fmt.Errorf("failed to delete role binding %s: %w", *rb.Id, err)
			}
			break
		}
	}
	ldapClient := aaa.NewLdapIdentitySourcesClient(sessionContext, connector)
	if ldapClient != nil {
		_ = ldapClient.Delete(testAccRoleBindingPathLdapID)
	}
	projClient := orgs.NewProjectsClient(sessionContext, connector)
	if projClient != nil && testAccRoleBindingPathProjectID != "" {
		_ = projClient.Delete(utl.DefaultOrgID, testAccRoleBindingPathProjectID, nil)
	}
	return nil
}

// testAccNsxtPolicyUserManagementRoleBindingPathTemplate returns config with project resource (Terraform creates project with nsx_id from testAccRoleBindingPathInitTestIDs),
// role binding data source (lookup by name/type), and path resource. LDAP and role binding are created in PreConfig.
func testAccNsxtPolicyUserManagementRoleBindingPathTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test_project" {
  nsx_id       = "%s"
  display_name = "test-rb-path-project"
  description  = "terraform test"
}

data "nsxt_policy_user_management_role_binding" "existing" {
  name                 = "%s"
  type                 = "remote_user"
  identity_source_id   = "%s"
}

resource "nsxt_policy_user_management_role_binding_path" "project_scope" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.existing.id
  path            = nsxt_policy_project.test_project.path
  roles           = ["project_admin"]
}
`, testAccRoleBindingPathProjectID, getTestLdapUser(), testAccRoleBindingPathLdapID)
}

func testAccNsxtPolicyUserManagementRoleBindingPathTemplateUpdate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test_project" {
  nsx_id       = "%s"
  display_name = "test-rb-path-project"
  description  = "terraform test"
}

data "nsxt_policy_user_management_role_binding" "existing" {
  name                 = "%s"
  type                 = "remote_user"
  identity_source_id   = "%s"
}

resource "nsxt_policy_user_management_role_binding_path" "project_scope" {
  role_binding_id = data.nsxt_policy_user_management_role_binding.existing.id
  path            = nsxt_policy_project.test_project.path
  roles           = ["network_engineer", "project_admin"]
}
`, testAccRoleBindingPathProjectID, getTestLdapUser(), testAccRoleBindingPathLdapID)
}
