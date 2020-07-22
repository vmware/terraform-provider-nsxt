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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.key", "DOMAIN_NAME"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.885794191.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.608540107", "*-myfiles.sharepoint.com"),
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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.951130521.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.951130521.key", "URL_CATEGORY"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.951130521.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.951130521.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.951130521.value.3317229948", "Abortion"),
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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.key", "DOMAIN_NAME"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.885794191.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.608540107", "*-myfiles.sharepoint.com"),
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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "2"),

					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.key", "APP_ID"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.3133243585.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.value.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.value.2076247700", "SSL"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.value.2328579708", "HTTP"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.3133243585.value.531481488", "SSH"),

					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.key", "DOMAIN_NAME"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.885794191.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.885794191.value.608540107", "*-myfiles.sharepoint.com"),
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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.key", "APP_ID"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.2307024415.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.key", "TLS_VERSION"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.value.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.value.2416475543", "TLS_V12"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.value.2721980181", "TLS_V10"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.sub_attribute.1680957162.value.645341863", "SSL_V3"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2307024415.value.2076247700", "SSL"),
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
					resource.TestCheckResourceAttr(testResourceName, "attribute.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.key", "APP_ID"),
					resource.TestCheckResourceAttrSet(testResourceName, "attribute.2585605962.is_alg_type"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.2040405228.data_type", "STRING"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.2040405228.key", "CIFS_SMB_VERSION"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.2040405228.value.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.2040405228.value.3398512106", "CIFS_SMB_V1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.sub_attribute.2040405228.value.3787226665", "CIFS_SMB_V2"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.value.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "attribute.2585605962.value.1810355442", "CIFS"),
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
	return fmt.Sprintf(`
attribute {
  data_type = "STRING"
  key = "DOMAIN_NAME"
  value = ["*-myfiles.sharepoint.com"]
}`)
}

func testAccNsxtPolicyContextProfileAttributeAppIDTemplate() string {
	return fmt.Sprintf(`
attribute {
  data_type = "STRING"
  key = "APP_ID"
  value = ["SSL", "SSH", "HTTP"]
}`)
}

func testAccNsxtPolicyContextProfileAttributeURLCategoryTemplate() string {
	return fmt.Sprintf(`
attribute {
  data_type = "STRING"
  key = "URL_CATEGORY"
  value = ["Abortion"]
}`)
}

func testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesTemplate() string {
	return fmt.Sprintf(`
attribute {
  data_type = "STRING"
  key = "APP_ID"
  value = ["SSL"]
  sub_attribute {
    data_type = "STRING"
    key = "TLS_VERSION"
    value = ["SSL_V3", "TLS_V10", "TLS_V12"]
  }
}`)
}

func testAccNsxtPolicyContextProfileAttributeAppIDSubAttributesUpdatedTemplate() string {
	return fmt.Sprintf(`
attribute {
  data_type = "STRING"
  key = "APP_ID"
  value = ["CIFS"]
  sub_attribute {
    data_type = "STRING"
    key = "CIFS_SMB_VERSION"
    value = ["CIFS_SMB_V1", "CIFS_SMB_V2"]
  }
}`)
}
