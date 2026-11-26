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

var accTestPolicyIPSecVpnIkeProfileCreateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform created",
	"dh_groups":             "GROUP2",
	"digest_algorithms":     "SHA2_256",
	"encryption_algorithms": "AES_128",
	"ike_version":           "IKE_FLEX",
	"sa_life_time":          "28100",
}

var accTestPolicyIPSecVpnIkeProfileUpdateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform updated",
	"dh_groups":             "GROUP5",
	"digest_algorithms":     "SHA2_512",
	"encryption_algorithms": "AES_256",
	"ike_version":           "IKE_V2",
	"sa_life_time":          "50000",
}

func TestAccResourceNsxtPolicyIPSecVpnIkeProfile_basic(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnIkeProfileBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyIPSecVpnIkeProfile_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnIkeProfileBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIPSecVpnIkeProfileBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_ipsec_vpn_ike_profile.test"
	testDataSourceName := "data.nsxt_policy_ipsec_vpn_ike_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnIkeProfileCheckDestroy(state, accTestPolicyIPSecVpnIkeProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnIkeProfileTemplate(true, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnIkeProfileExists(accTestPolicyIPSecVpnIkeProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnIkeProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnIkeProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dh_groups.0", accTestPolicyIPSecVpnIkeProfileCreateAttributes["dh_groups"]),
					resource.TestCheckResourceAttr(testResourceName, "digest_algorithms.0", accTestPolicyIPSecVpnIkeProfileCreateAttributes["digest_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "encryption_algorithms.0", accTestPolicyIPSecVpnIkeProfileCreateAttributes["encryption_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_version", accTestPolicyIPSecVpnIkeProfileCreateAttributes["ike_version"]),
					resource.TestCheckResourceAttr(testResourceName, "sa_life_time", accTestPolicyIPSecVpnIkeProfileCreateAttributes["sa_life_time"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnIkeProfileTemplate(false, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnIkeProfileExists(accTestPolicyIPSecVpnIkeProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "dh_groups.0", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["dh_groups"]),
					resource.TestCheckResourceAttr(testResourceName, "digest_algorithms.0", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["digest_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "encryption_algorithms.0", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["encryption_algorithms"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_version", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["ike_version"]),
					resource.TestCheckResourceAttr(testResourceName, "sa_life_time", accTestPolicyIPSecVpnIkeProfileUpdateAttributes["sa_life_time"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnIkeProfileMinimalistic(withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnIkeProfileExists(accTestPolicyIPSecVpnIkeProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyIPSecVpnIkeProfile_importBasic(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnIkeProfileImport(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyIPSecVpnIkeProfile_importBasic_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIPSecVpnIkeProfileImport(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIPSecVpnIkeProfileImport(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_ipsec_vpn_ike_profile.test"

	var importStateIDFunc resource.ImportStateIdFunc
	if withContext {
		importStateIDFunc = testAccResourceNsxtPolicyImportIDRetriever(testResourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnIkeProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnIkeProfileMinimalistic(withContext),
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

func testAccNsxtPolicyIPSecVpnIkeProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		sessionContext := testAccGetSessionContext()

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPSecVpnIkeProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IPSecVpnIkeProfile resource ID not set in resources")
		}
		client := cliIpsecVpnIkeProfilesClient(sessionContext, connector)
		_, err := client.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy IPSecVpnIkeProfile ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyIPSecVpnIkeProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	sessionContext := testAccGetSessionContext()
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ipsec_vpn_ike_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		client := cliIpsecVpnIkeProfilesClient(sessionContext, connector)
		_, err := client.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy IPSecVpnIkeProfile %s still exists", displayName)
		}
		if !isNotFoundError(err) {
			return err
		}
	}
	return nil
}

func testAccNsxtPolicyIPSecVpnIkeProfileTemplate(createFlow bool, withContext bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPSecVpnIkeProfileCreateAttributes
	} else {
		attrMap = accTestPolicyIPSecVpnIkeProfileUpdateAttributes
	}
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_ike_profile" "test" {
%s
  display_name          = "%s"
  description           = "%s"
  dh_groups             = ["%s"]
  digest_algorithms     = ["%s"]
  encryption_algorithms = ["%s"]
  ike_version           = "%s"
  sa_life_time          = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
data "nsxt_policy_ipsec_vpn_ike_profile" "test" {
%s
  display_name = "%s"
  depends_on   = [nsxt_policy_ipsec_vpn_ike_profile.test]
}
`, context, attrMap["display_name"], attrMap["description"], attrMap["dh_groups"], attrMap["digest_algorithms"], attrMap["encryption_algorithms"], attrMap["ike_version"], attrMap["sa_life_time"], context, attrMap["display_name"])
}

func testAccNsxtPolicyIPSecVpnIkeProfileMinimalistic(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_ike_profile" "test" {
%s
  display_name          = "%s"
  dh_groups             = ["%s"]
  encryption_algorithms = ["%s"]
}
data "nsxt_policy_ipsec_vpn_ike_profile" "test" {
%s
  id = nsxt_policy_ipsec_vpn_ike_profile.test.id
}`, context, accTestPolicyIPSecVpnIkeProfileUpdateAttributes["display_name"], accTestPolicyIPSecVpnIkeProfileUpdateAttributes["dh_groups"], "AES_GCM_192", context)
}
