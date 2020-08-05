/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

func TestAccResourceNsxtPolicyContextProfile_basic(t *testing.T) {
	name := "terraform-test"
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeURLCategoryTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.608540107", "*-myfiles.sharepoint.com"),
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
					resource.TestCheckResourceAttr(testResourceName, "url_category.1857560543.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "url_category.1857560543.value.3317229948", "Abortion"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyContextProfile_importBasic(t *testing.T) {
	name := "terra-test-import"
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
	name := "terraform-test"
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate() + testAccNsxtPolicyContextProfileAttributeAppIDTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.608540107", "*-myfiles.sharepoint.com"),
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
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.124112814.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.124112814.value.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.124112814.value.2076247700", "SSL"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.124112814.value.2328579708", "HTTP"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.124112814.value.531481488", "SSH"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.608540107", "*-myfiles.sharepoint.com"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyContextProfile_subAttributes(t *testing.T) {
	name := "terraform-test"
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_context_profile.test"
	attributes := testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesTemplate()
	updatedAttributes := testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesUpdatedTemplate()
	attributesNoSub := testAccNsxtPolicyContextProfileAttributeAppIDSslTemplate()
	attributesDomainName := testAccNsxtPolicyContextProfileAttributeDomainNameTemplate()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.2652970086.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.value.2076247700", "SSL"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.cifs_smb_version.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.tls_cipher_suite.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.tls_version.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.tls_version.2416475543", "TLS_V12"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.tls_version.2721980181", "TLS_V10"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.2652970086.sub_attribute.250747496.tls_version.645341863", "SSL_V3"),
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
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.832058275.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.value.1810355442", "CIFS"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.2393957168.cifs_smb_version.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.2393957168.tls_cipher_suite.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.2393957168.tls_version.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.2393957168.cifs_smb_version.3398512106", "CIFS_SMB_V1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.832058275.sub_attribute.2393957168.cifs_smb_version.3787226665", "CIFS_SMB_V2"),
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
					resource.TestCheckResourceAttr(testResourceName, "app_id.4162008338.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.4162008338.value.2076247700", "SSL"),
					resource.TestCheckResourceAttrSet(testResourceName, "app_id.4162008338.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "app_id.4162008338.sub_attribute.#", "0"),
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
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "domain_name.238902231.value.608540107", "*-myfiles.sharepoint.com"),
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
		nsxClient := gm_infra.NewDefaultContextProfilesClient(connector)
		_, err = nsxClient.Get(resourceID)
	} else {
		nsxClient := infra.NewDefaultContextProfilesClient(connector)
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
