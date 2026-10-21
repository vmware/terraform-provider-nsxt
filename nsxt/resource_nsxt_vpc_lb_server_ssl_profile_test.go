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

var accTestVpcLBServerSslProfileCreateAttributes = map[string]string{
	"description":           "terraform created",
	"cipher_group_label":    "CUSTOM",
	"ciphers":               "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"protocols":             "TLS_V1_2",
	"session_cache_enabled": "true",
}

var accTestVpcLBServerSslProfileUpdateAttributes = map[string]string{
	"description":           "terraform updated",
	"cipher_group_label":    "CUSTOM",
	"ciphers":               "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
	"protocols":             "TLS_V1_2",
	"session_cache_enabled": "false",
}

func TestAccResourceNsxtVpcLBServerSslProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_server_ssl_profile.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBServerSslProfileCheckDestroy(state, updateName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBServerSslProfileTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServerSslProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBServerSslProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "cipher_group_label", accTestVpcLBServerSslProfileCreateAttributes["cipher_group_label"]),
					resource.TestCheckResourceAttr(testResourceName, "ciphers.0", accTestVpcLBServerSslProfileCreateAttributes["ciphers"]),
					resource.TestCheckResourceAttr(testResourceName, "protocols.0", accTestVpcLBServerSslProfileCreateAttributes["protocols"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cache_enabled", accTestVpcLBServerSslProfileCreateAttributes["session_cache_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
					resource.TestCheckResourceAttrSet(testResourceName, "is_secure"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_fips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBServerSslProfileTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServerSslProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBServerSslProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "cipher_group_label", accTestVpcLBServerSslProfileUpdateAttributes["cipher_group_label"]),
					resource.TestCheckResourceAttr(testResourceName, "ciphers.0", accTestVpcLBServerSslProfileUpdateAttributes["ciphers"]),
					resource.TestCheckResourceAttr(testResourceName, "protocols.0", accTestVpcLBServerSslProfileUpdateAttributes["protocols"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cache_enabled", accTestVpcLBServerSslProfileUpdateAttributes["session_cache_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBServerSslProfileMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServerSslProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		}),
	})
}

func TestAccResourceNsxtVpcLBServerSslProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_server_ssl_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBServerSslProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBServerSslProfileMinimalistic(name),
			},
			{
				// VPC objects can only be imported by policy path
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		})),
	})
}

func testAccNsxtVpcLBServerSslProfileExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPC LB Server SSL Profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VPC LB Server SSL Profile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBServerSslProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VPC LB Server SSL Profile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBServerSslProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_lb_server_ssl_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBServerSslProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("VPC LB Server SSL Profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBServerSslProfileTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBServerSslProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBServerSslProfileCreateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_server_ssl_profile" "test" {
%s
  display_name          = "%s"
  description           = "%s"
  cipher_group_label    = "%s"
  ciphers               = ["%s"]
  protocols             = ["%s"]
  session_cache_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["cipher_group_label"], attrMap["ciphers"], attrMap["protocols"], attrMap["session_cache_enabled"])
}

func testAccNsxtVpcLBServerSslProfileMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_server_ssl_profile" "test" {
%s
  display_name       = "%s"
  cipher_group_label = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name, accTestVpcLBServerSslProfileUpdateAttributes["cipher_group_label"])
}
