/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var accTestPolicyLdapIdentitySourceCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
}

var accTestPolicyLdapIdentitySourceUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
}

func TestAccResourceNsxtPolicyLdapIdentitySource_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ldap_identity_source.test"
	ldapType := nsxModel.LdapIdentitySource_RESOURCE_TYPE_ACTIVEDIRECTORYIDENTITYSOURCE

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_PASSWORD")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_URL")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_CERT")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_DOMAIN")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_BASE_DN")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLdapIdentitySourceCheckDestroy(state, accTestPolicyLdapIdentitySourceUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLdapIdentitySourceCreate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapUser(), getTestLdapPassword(),
					getTestLdapURL(), getTestLdapCert()),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLdapIdentitySourceExists(accTestPolicyLdapIdentitySourceCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLdapIdentitySourceCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLdapIdentitySourceCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "type", ldapType),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", getTestLdapDomain()),
					resource.TestCheckResourceAttr(testResourceName, "base_dn", getTestLdapBaseDN()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.bind_identity", getTestLdapUser()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.password", getTestLdapPassword()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.url", getTestLdapURL()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.certificates.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyLdapIdentitySourceUpdate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapUser(), getTestLdapPassword(),
					getTestLdapURL(), getTestLdapCert()),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLdapIdentitySourceExists(accTestPolicyLdapIdentitySourceUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLdapIdentitySourceUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLdapIdentitySourceUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "type", ldapType),
					resource.TestCheckResourceAttr(testResourceName, "domain_name", getTestLdapDomain()),
					resource.TestCheckResourceAttr(testResourceName, "base_dn", getTestLdapBaseDN()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.bind_identity", getTestLdapUser()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.password", getTestLdapPassword()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.url", getTestLdapURL()),
					resource.TestCheckResourceAttr(testResourceName, "ldap_server.0.certificates.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLdapIdentitySource_import_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ldap_identity_source.test"
	ldapType := nsxModel.LdapIdentitySource_RESOURCE_TYPE_ACTIVEDIRECTORYIDENTITYSOURCE

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LDAP_USER")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_PASSWORD")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_URL")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_CERT")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_DOMAIN")
			testAccEnvDefined(t, "NSXT_TEST_LDAP_BASE_DN")
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLdapIdentitySourceCheckDestroy(state, accTestPolicyLdapIdentitySourceCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLdapIdentitySourceCreate(
					ldapType, getTestLdapDomain(), getTestLdapBaseDN(), getTestLdapUser(), getTestLdapPassword(),
					getTestLdapURL(), getTestLdapCert()),
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

func testAccNsxtPolicyLdapIdentitySourceCreate(serverType, domainName, baseDn, bindUser, bindPwd, url, cert string) string {
	attrMap := accTestPolicyLdapIdentitySourceCreateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_ldap_identity_source" "test" {
    display_name = "%s"
    description  = "%s"
    type         = "%s"
    domain_name  = "%s"
    base_dn      = "%s"

    ldap_server {
        bind_identity = "%s"
        password      = "%s"
        url           = "%s"
        certificates  = [
            <<-EOT
%s
            EOT
            ,
        ]
    }

    tag {
        scope = "scope1"
        tag = "tag1"
    }
}`, attrMap["display_name"], attrMap["description"], serverType, domainName, baseDn, bindUser, bindPwd, url, cert)
}

func testAccNsxtPolicyLdapIdentitySourceUpdate(serverType, domainName, baseDn, bindUser, bindPwd, url, cert string) string {
	attrMap := accTestPolicyLdapIdentitySourceUpdateAttributes
	return fmt.Sprintf(`
resource "nsxt_policy_ldap_identity_source" "test" {
    display_name = "%s"
    description  = "%s"
    type         = "%s"
    domain_name  = "%s"
    base_dn      = "%s"

    ldap_server {
        bind_identity = "%s"
        password      = "%s"
        url           = "%s"
        certificates  = [
            <<-EOT
%s
            EOT
            ,
        ]
    }
}`, attrMap["display_name"], attrMap["description"], serverType, domainName, baseDn, bindUser, bindPwd, url, cert)
}
