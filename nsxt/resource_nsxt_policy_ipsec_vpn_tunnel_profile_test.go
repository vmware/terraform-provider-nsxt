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

var accTestPolicyIPSecVpnTunnelProfileCreateAttributes = map[string]string{
	"display_name":                   getAccTestResourceName(),
	"description":                    "terraform created",
	"df_policy":                      "COPY",
	"dh_groups":                      "GROUP2",
	"digest_algorithms":              "SHA2_256",
	"enable_perfect_forward_secrecy": "true",
	"encryption_algorithms":          "AES_256",
	"sa_life_time":                   "46000",
}

var accTestPolicyIPSecVpnTunnelProfileUpdateAttributes = map[string]string{
	"display_name":                   getAccTestResourceName(),
	"description":                    "terraform updated",
	"df_policy":                      "CLEAR",
	"dh_groups":                      "GROUP5",
	"digest_algorithms":              "SHA1",
	"enable_perfect_forward_secrecy": "false",
	"encryption_algorithms":          "AES_128",
	"sa_life_time":                   "28000",
}

func TestAccResourceNsxtPolicyIPSecVpnTunnelProfile_basic(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnTunnelProfileBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyIPSecVpnTunnelProfile_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnTunnelProfileBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIPSecVpnTunnelProfileBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_ipsec_vpn_tunnel_profile.test"
	testDataSourceName := "data.nsxt_policy_ipsec_vpn_tunnel_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnTunnelProfileCheckDestroy(state, accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnTunnelProfileTemplate(true, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnTunnelProfileExists(accTestPolicyIPSecVpnTunnelProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "df_policy", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["df_policy"]),
					resource.TestCheckResourceAttr(testResourceName, "dh_groups.0", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["dh_groups"]),
					resource.TestCheckResourceAttr(testResourceName, "digest_algorithms.0", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["digest_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "enable_perfect_forward_secrecy", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["enable_perfect_forward_secrecy"]),
					resource.TestCheckResourceAttr(testResourceName, "encryption_algorithms.0", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["encryption_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "sa_life_time", accTestPolicyIPSecVpnTunnelProfileCreateAttributes["sa_life_time"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnTunnelProfileTemplate(false, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnTunnelProfileExists(accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "df_policy", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["df_policy"]),
					resource.TestCheckResourceAttr(testResourceName, "dh_groups.0", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["dh_groups"]),
					resource.TestCheckResourceAttr(testResourceName, "digest_algorithms.0", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["digest_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "enable_perfect_forward_secrecy", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["enable_perfect_forward_secrecy"]),
					resource.TestCheckResourceAttr(testResourceName, "encryption_algorithms.0", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["encryption_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "sa_life_time", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["sa_life_time"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnTunnelProfileMinimalistic(withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnTunnelProfileExists(accTestPolicyIPSecVpnTunnelProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPSecVpnTunnelProfile_importBasic(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnTunnelProfileImport(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyIPSecVpnTunnelProfile_importBasic_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnTunnelProfileImport(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIPSecVpnTunnelProfileImport(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_ipsec_vpn_tunnel_profile.test"

	var importStateIDFunc resource.ImportStateIdFunc
	if withContext {
		importStateIDFunc = testAccResourceNsxtPolicyImportIDRetriever(testResourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnTunnelProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnTunnelProfileMinimalistic(withContext),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateIDFunc,
			},
		},
	})
}

func testAccNsxtPolicyIPSecVpnTunnelProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		sessionContext := testAccGetSessionContext()

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPSecVpnTunnelProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IPSecVpnTunnelProfile resource ID not set in resources")
		}
		client := cliIpsecVpnTunnelProfilesClient(sessionContext, connector)
		_, err := client.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy IPSecVpnTunnelProfile ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyIPSecVpnTunnelProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	sessionContext := testAccGetSessionContext()
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ipsec_vpn_tunnel_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		client := cliIpsecVpnTunnelProfilesClient(sessionContext, connector)
		_, err := client.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy IPSecVpnTunnelProfile %s still exists", displayName)
		}
		if !isNotFoundError(err) {
			return err
		}
	}
	return nil
}

func testAccNsxtPolicyIPSecVpnTunnelProfileTemplate(createFlow bool, withContext bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPSecVpnTunnelProfileCreateAttributes
	} else {
		attrMap = accTestPolicyIPSecVpnTunnelProfileUpdateAttributes
	}
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
%s
  display_name                   = "%s"
  description                    = "%s"
  df_policy                      = "%s"
  dh_groups                      = ["%s"]
  digest_algorithms              = ["%s"]
  enable_perfect_forward_secrecy = %s
  encryption_algorithms          = ["%s"]
  sa_life_time                   = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
data "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
%s
  display_name = "%s"
  depends_on   = [nsxt_policy_ipsec_vpn_tunnel_profile.test]
}`, context, attrMap["display_name"], attrMap["description"], attrMap["df_policy"], attrMap["dh_groups"], attrMap["digest_algorithms"], attrMap["enable_perfect_forward_secrecy"], attrMap["encryption_algorithms"], attrMap["sa_life_time"], context, attrMap["display_name"])
}

func testAccNsxtPolicyIPSecVpnTunnelProfileMinimalistic(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
%s
  display_name          = "%s"
  encryption_algorithms = ["%s"]
  dh_groups             = ["%s"]
}
data "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
%s
  id = nsxt_policy_ipsec_vpn_tunnel_profile.test.id
}`, context, accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["display_name"], "AES_GCM_192", accTestPolicyIPSecVpnTunnelProfileUpdateAttributes["dh_groups"], context)
}
