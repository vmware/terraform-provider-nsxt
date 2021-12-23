/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
)

func TestAccResourceNsxtPolicyContextProfile_basic(t *testing.T) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeURLCategoryTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyContextProfileCheckDestroy(state, testResourceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileTemplate(name, attributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.0", "*-myfiles.sharepoint.com"),
				),
			},
			{
				Config: testAccNsxtPolicyContextProfileTemplate(updatedName, updatedAttributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.0.value.0", "Abortion"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyContextProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyContextProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileTemplate(name, attributes),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyContextProfile_multipleAttributes(t *testing.T) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate() + testAccNsxtPolicyContextProfileAttributeAppIDTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyContextProfileCheckDestroy(state, testResourceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileTemplate(name, attributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.0", "*-myfiles.sharepoint.com"),
				),
			},
			{
				Config: testAccNsxtPolicyContextProfileTemplate(updatedName, updatedAttributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.0.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.0", "HTTP"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.1", "SSH"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.2", "SSL"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.0", "*-myfiles.sharepoint.com"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyContextProfile_subAttributes(t *testing.T) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesUpdatedTemplate()
	attributesNoSub := testAccNsxtPolicyContextProfileAttributeAppIDSslTemplate()
	attributesDomainName := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyContextProfileCheckDestroy(state, testResourceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileTemplate(name, attributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.0.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.0", "SSL"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.cifs_smb_version.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_cipher_suite.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_version.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_version.2", "TLS_V12"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_version.1", "TLS_V10"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_version.0", "SSL_V3"),
				),
			},
			{
				Config: testAccNsxtPolicyContextProfileTemplate(updatedName, updatedAttributes),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.0.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.0", "CIFS"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.cifs_smb_version.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_cipher_suite.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.tls_version.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.cifs_smb_version.0", "CIFS_SMB_V1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.0.cifs_smb_version.1", "CIFS_SMB_V2"),
				),
			},
			{
				Config: testAccNsxtPolicyContextProfileTemplate(updatedName, attributesNoSub),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.value.0", "SSL"),
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.0.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.0.sub_attribute.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyContextProfileTemplate(updatedName, attributesDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyContextProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.0.value.0", "*-myfiles.sharepoint.com"),
				),
			},
		},
	})
}

func testAccNsxtPolicyContextProfileExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy ContextProfile resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy ContextProfile resource ID not set in resources")
		}

		err := nsxtPolicyContextProfileExists(resourceID)

		if err != nil {
			return fmt.Errorf("Error while retrieving policy ContextProfile ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func nsxtPolicyContextProfileExists(resourceID string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	var err error
	if testAccIsGlobalManager() {
		nsxClient := gm_infra.NewContextProfilesClient(connector)
		_, err = nsxClient.Get(resourceID)
	} else {
		nsxClient := infra.NewContextProfilesClient(connector)
		_, err = nsxClient.Get(resourceID)
	}
	return err
}

func testAccNsxtPolicyContextProfileCheckDestroy(state *terraform.State, displayName string) error {
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_context_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		err := nsxtPolicyContextProfileExists(resourceID)
		if err == nil {
			return fmt.Errorf("Policy ContextProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyContextProfileTemplate(name string, attributes string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_context_profile" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  tag {
    scope = "color"
    tag   = "orange"
  }
%s
}`, name, attributes)
}

func testAccNsxtPolicyContextProfileAttributeDomainNameTemplate() string {
	return `
domain_name {
  value     = ["*-myfiles.sharepoint.com"]
}`
}

func testAccNsxtPolicyContextProfileAttributeAppIDTemplate() string {
	return `
app_id {
  value     = ["SSL", "SSH", "HTTP"]
}`
}

func testAccNsxtPolicyContextProfileAttributeURLCategoryTemplate() string {
	return `
url_category {
  value     = ["Abortion"]
}`
}

func testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesTemplate() string {
	return `
app_id {
  value     = ["SSL"]
  sub_attribute {
    tls_version = ["SSL_V3", "TLS_V10", "TLS_V12"]
  }
}`
}

func testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesUpdatedTemplate() string {
	return `
app_id {
  value     = ["CIFS"]
  sub_attribute {
    cifs_smb_version = ["CIFS_SMB_V1", "CIFS_SMB_V2"]
  }
}`
}

func testAccNsxtPolicyContextProfileAttributeAppIDSslTemplate() string {
	return `
app_id {
  value     = ["SSL"]
}`
}
