/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestPolicyLBServerSslProfileCreateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform created",
	"cipher_group_label":    "BALANCED",
	"ciphers":               "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	"protocols":             "TLS_V1_2",
	"session_cache_enabled": "true",
}

var accTestPolicyLBServerSslProfileUpdateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform updated",
	"cipher_group_label":    "HIGH_SECURITY",
	"ciphers":               "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"protocols":             "TLS_V1_2",
	"session_cache_enabled": "false",
}

func TestAccResourceNsxtPolicyLBServerSslProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_server_ssl_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBServerSslProfileCheckDestroy(state, accTestPolicyLBServerSslProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBServerSslProfileCustomTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServerSslProfileExists(accTestPolicyLBServerSslProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBServerSslProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBServerSslProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "cipher_group_label", "CUSTOM"),
					resource.TestCheckResourceAttr(testResourceName, "ciphers.0", accTestPolicyLBServerSslProfileCreateAttributes["ciphers"]),
					resource.TestCheckResourceAttr(testResourceName, "protocols.0", accTestPolicyLBServerSslProfileCreateAttributes["protocols"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cache_enabled", accTestPolicyLBServerSslProfileCreateAttributes["session_cache_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "is_secure"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_fips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyLBServerSslProfileCustomTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServerSslProfileExists(accTestPolicyLBServerSslProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBServerSslProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBServerSslProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "cipher_group_label", "CUSTOM"),
					resource.TestCheckResourceAttr(testResourceName, "ciphers.0", accTestPolicyLBServerSslProfileUpdateAttributes["ciphers"]),
					resource.TestCheckResourceAttr(testResourceName, "protocols.0", accTestPolicyLBServerSslProfileUpdateAttributes["protocols"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cache_enabled", accTestPolicyLBServerSslProfileUpdateAttributes["session_cache_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "is_secure"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_fips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyLBServerSslProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServerSslProfileExists(accTestPolicyLBServerSslProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBServerSslProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBServerSslProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "cipher_group_label", accTestPolicyLBServerSslProfileCreateAttributes["cipher_group_label"]),
					resource.TestCheckResourceAttr(testResourceName, "protocols.0", accTestPolicyLBServerSslProfileCreateAttributes["protocols"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cache_enabled", accTestPolicyLBServerSslProfileCreateAttributes["session_cache_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "is_secure"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_fips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBServerSslProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBServerSslProfileExists(accTestPolicyLBServerSslProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLBServerSslProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_server_ssl_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBServerSslProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBServerSslProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBServerSslProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBServerSslProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBServerSslProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBServerSslProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBServerSslProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyLBServerSslProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_server_ssl_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBServerSslProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBServerSslProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBServerSslProfileCustomTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBServerSslProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBServerSslProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_server_ssl_profile" "test" {
  display_name          = "%s"
  description           = "%s"
  cipher_group_label    = "CUSTOM"
  ciphers               = ["%s"]
  protocols             = ["%s"]
  session_cache_enabled = %s
}`, attrMap["display_name"], attrMap["description"], attrMap["ciphers"], attrMap["protocols"], attrMap["session_cache_enabled"])
}

func testAccNsxtPolicyLBServerSslProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBServerSslProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBServerSslProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_server_ssl_profile" "test" {
  display_name          = "%s"
  description           = "%s"
  cipher_group_label    = "%s"
  protocols             = ["%s"]
  session_cache_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["cipher_group_label"], attrMap["protocols"], attrMap["session_cache_enabled"])
}

func testAccNsxtPolicyLBServerSslProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_server_ssl_profile" "test" {
  display_name       = "%s"
  cipher_group_label = "%s"
}`, accTestPolicyLBServerSslProfileUpdateAttributes["display_name"], accTestPolicyLBServerSslProfileUpdateAttributes["cipher_group_label"])
}
