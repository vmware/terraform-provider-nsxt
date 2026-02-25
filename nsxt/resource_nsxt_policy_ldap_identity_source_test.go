// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/aaa"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var accTestPolicyLdapIdentitySourceCreateAttributes = map[string]string{
	"nsx_id":      getAccTestResourceName(),
	"description": "terraform created",
}

var accTestPolicyLdapIdentitySourceUpdateAttributes = map[string]string{
	"nsx_id":      getAccTestResourceName(),
	"description": "terraform updated",
}

// testAccDeleteAllLdapIdentitySources deletes every LDAP identity source on the manager.
// NSX allows only one source per authentication domain; leftovers from parallel packages or failed
// runs cause "Domain ... is already mapped" (e.g. error 53013) when creating a new source.
// Call from CheckDestroy so the manager is left clean after each test.
func testAccDeleteAllLdapIdentitySources() error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("failed to get policy connector: %w", err)
	}
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	ldapClient := aaa.NewLdapIdentitySourcesClient(sessionContext, connector)
	if ldapClient == nil {
		return fmt.Errorf("LDAP identity sources client is nil")
	}
	converter := bindings.NewTypeConverter()
	var cursor *string
	for {
		listResult, err := ldapClient.List(cursor, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to list LDAP identity sources: %w", err)
		}
		for _, structVal := range listResult.Results {
			if structVal == nil {
				continue
			}
			obj, errs := converter.ConvertToGolang(structVal, nsxModel.LdapIdentitySourceBindingType())
			if errs != nil {
				return fmt.Errorf("failed to convert LDAP identity source from list: %w", errs[0])
			}
			ldapObj := obj.(nsxModel.LdapIdentitySource)
			if ldapObj.Id == nil || *ldapObj.Id == "" {
				continue
			}
			if err := ldapClient.Delete(*ldapObj.Id); err != nil && !isNotFoundError(err) {
				return fmt.Errorf("failed to delete LDAP identity source %s: %w", *ldapObj.Id, err)
			}
		}
		if listResult.Cursor == nil || *listResult.Cursor == "" {
			break
		}
		cursor = listResult.Cursor
	}
	return nil
}

func TestAccResourceNsxtPolicyLdapIdentitySource_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ldap_identity_source.test"
	ldapType := openLdapType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_PASSWORD")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_URL")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_DOMAIN")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_BASE_DN")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccNsxtPolicyLdapIdentitySourceCheckDestroy(state, accTestPolicyLdapIdentitySourceUpdateAttributes["nsx_id"]); err != nil {
				return err
			}
			return testAccDeleteAllLdapIdentitySources()
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLdapIdentitySourceCreate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapAdminUser(), getTestLdapAdminPassword(), getTestLdapURL()),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLdapIdentitySourceExists(accTestPolicyLdapIdentitySourceCreateAttributes["nsx_id"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLdapIdentitySourceCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "type", ldapType),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", getTestLdapDomain()),
					resource.TestCheckResourceAttr(testResourceName, "base_dn", getTestLdapBaseDN()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.bind_identity", getTestLdapAdminUser()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.password", getTestLdapAdminPassword()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.url", getTestLdapURL()),
					//resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.certificates.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "nsx_id", accTestPolicyLdapIdentitySourceCreateAttributes["nsx_id"]),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyLdapIdentitySourceUpdate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapAdminUser(), getTestLdapAdminPassword(), getTestLdapURL()),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLdapIdentitySourceExists(accTestPolicyLdapIdentitySourceUpdateAttributes["nsx_id"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLdapIdentitySourceUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "type", ldapType),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", getTestLdapDomain()),
					resource.TestCheckResourceAttr(testResourceName, "base_dn", getTestLdapBaseDN()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.bind_identity", getTestLdapAdminUser()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.password", getTestLdapAdminPassword()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.url", getTestLdapURL()),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),

					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLdapIdentitySource_import_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ldap_identity_source.test"
	ldapType := activeDirectoryType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_ADMIN_PASSWORD")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_URL")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_DOMAIN")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_BASE_DN")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccNsxtPolicyLdapIdentitySourceCheckDestroy(state, accTestPolicyLdapIdentitySourceCreateAttributes["nsx_id"]); err != nil {
				return err
			}
			return testAccDeleteAllLdapIdentitySources()
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLdapIdentitySourceCreate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapAdminUser(), getTestLdapAdminPassword(), getTestLdapURL()),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ldap_server.0.password"},
			},
		},
	})
}

func testAccNsxtPolicyLdapIdentitySourceExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("LdapIdentitySource resource %s not found in resources", resourceName)
		}

		ldapSourceID := rs.Primary.Attributes["id"]
		if ldapSourceID == "" {
			return fmt.Errorf("LdapIdentitySource resource ID not set in resources")
		}
		exist, err := resourceNsxtPolicyLdapIdentitySourceExists(ldapSourceID, connector, false)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("LdapIdentitySource %s does not exist", displayName)
		}

		return nil
	}
}

func testAccNsxtPolicyLdapIdentitySourceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_ldap_identity_source" {
			continue
		}

		ldapSourceID := rs.Primary.Attributes["id"]
		if ldapSourceID == "" {
			return fmt.Errorf("LdapIdentitySource resource ID not set in resources")
		}
		exist, err := resourceNsxtPolicyLdapIdentitySourceExists(ldapSourceID, connector, false)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("LdapIdentitySource %s still exists", displayName)
		}
		return nil
	}
	return nil
}

func testAccNsxtPolicyLdapIdentitySourceCreate(serverType, domainName, baseDn, bindUser, bindPwd, url string) string {
	attrMap := accTestPolicyLdapIdentitySourceCreateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_ldap_identity_source" "test" {
  nsx_id      = "%s"
  description = "%s"
  type        = "%s"
  domain_name = "%s"
  base_dn     = "%s"

  ldap_server {
    bind_identity = "%s"
    password      = "%s"
    url           = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["nsx_id"], attrMap["description"], serverType, domainName, baseDn, bindUser, bindPwd, url)
}

func testAccNsxtPolicyLdapIdentitySourceUpdate(serverType, domainName, baseDn, bindUser, bindPwd, url string) string {
	attrMap := accTestPolicyLdapIdentitySourceUpdateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_ldap_identity_source" "test" {
  nsx_id      = "%s"
  description = "%s"
  type        = "%s"
  domain_name = "%s"
  base_dn     = "%s"

  ldap_server {
    bind_identity = "%s"
    password      = "%s"
    url           = "%s"
  }
}`, attrMap["nsx_id"], attrMap["description"], serverType, domainName, baseDn, bindUser, bindPwd, url)
}
