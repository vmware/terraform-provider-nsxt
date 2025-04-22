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

var accTestPolicyLBCookiePersistenceProfileCreateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"persistence_shared": "true",
	"cookie_name":        "test-create",
	"cookie_mode":        "INSERT",
	"cookie_fallback":    "true",
	"cookie_garble":      "true",
	"cookie_domain":      "chocolate-chip.com",
	"cookie_path":        "test-create",
	"cookie_httponly":    "true",
	"cookie_secure":      "true",
	"max_idle":           "180",
	"max_life":           "1800",
}

var accTestPolicyLBCookiePersistenceProfileUpdateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform updated",
	"persistence_shared": "false",
	"cookie_name":        "test-update",
	"cookie_mode":        "PREFIX",
	"cookie_fallback":    "false",
	"cookie_garble":      "false",
	"cookie_domain":      "orange.com",
	"cookie_path":        "test-update",
	"cookie_httponly":    "false",
	"cookie_secure":      "false",
	"max_idle":           "280",
	"max_life":           "2800",
}

func TestAccResourceNsxtPolicyLBCookiePersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_cookie_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBCookiePersistenceProfileCheckDestroy(state, accTestPolicyLBCookiePersistenceProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBCookiePersistenceProfileSessionCookieTimeTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBCookiePersistenceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBCookiePersistenceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBCookiePersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBCookiePersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_name", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_name"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_mode", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_fallback", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_fallback"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_garble", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_garble"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_domain", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_domain"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_path", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_path"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.0.max_idle", accTestPolicyLBCookiePersistenceProfileCreateAttributes["max_idle"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.0.max_life", accTestPolicyLBCookiePersistenceProfileCreateAttributes["max_life"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_httponly", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_httponly"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_secure", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_secure"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBCookiePersistenceProfilePersistenceCookieTimeTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBCookiePersistenceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBCookiePersistenceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBCookiePersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBCookiePersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_name", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_name"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_mode", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_fallback", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_fallback"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_garble", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_garble"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_domain", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_domain"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_path", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_path"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_cookie_time.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_cookie_time.0.max_idle", accTestPolicyLBCookiePersistenceProfileCreateAttributes["max_idle"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_httponly", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_httponly"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_secure", accTestPolicyLBCookiePersistenceProfileCreateAttributes["cookie_secure"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBCookiePersistenceProfileNoCookieTimeTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBCookiePersistenceProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_name", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_name"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_mode", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_fallback", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_fallback"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_garble", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_garble"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_domain", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_domain"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_path", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_path"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "cookie_httponly", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_httponly"]),
					resource.TestCheckResourceAttr(testResourceName, "cookie_secure", accTestPolicyLBCookiePersistenceProfileUpdateAttributes["cookie_secure"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBCookiePersistenceProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBCookiePersistenceProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyLBCookiePersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_cookie_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBCookiePersistenceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBCookiePersistenceProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBPersistenceProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LB Persistence Profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LB Persistence Profile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBPersistenceProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBCookiePersistenceProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyLBCookiePersistenceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_cookie_persistence_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBPersistenceProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBCookiePersistenceProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBCookiePersistenceProfileNoCookieTimeTemplate(createFlow bool) string {
	return testAccNsxtPolicyLBCookiePersistenceProfileTemplate(createFlow, "")
}

func testAccNsxtPolicyLBCookiePersistenceProfileSessionCookieTimeTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBCookiePersistenceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBCookiePersistenceProfileUpdateAttributes
	}
	cookieTime := fmt.Sprintf(`
  session_cookie_time {
     max_idle = %s
     max_life = %s
  }`, attrMap["max_idle"], attrMap["max_life"])

	return testAccNsxtPolicyLBCookiePersistenceProfileTemplate(createFlow, cookieTime)
}

func testAccNsxtPolicyLBCookiePersistenceProfilePersistenceCookieTimeTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBCookiePersistenceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBCookiePersistenceProfileUpdateAttributes
	}
	cookieTime := fmt.Sprintf(`
  persistence_cookie_time {
     max_idle = %s
  }`, attrMap["max_idle"])

	return testAccNsxtPolicyLBCookiePersistenceProfileTemplate(createFlow, cookieTime)
}

func testAccNsxtPolicyLBCookiePersistenceProfileTemplate(createFlow bool, cookieTime string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBCookiePersistenceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBCookiePersistenceProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_cookie_persistence_profile" "test" {
  display_name = "%s"
  description  = "%s"

  persistence_shared = %s
  cookie_name        = "%s"
  cookie_mode        = "%s"
  cookie_fallback    = %s
  cookie_garble      = %s
  cookie_domain      = "%s"
  cookie_path        = "%s"
  cookie_httponly    = %s
  cookie_secure      = %s
  %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["persistence_shared"], attrMap["cookie_name"], attrMap["cookie_mode"], attrMap["cookie_fallback"], attrMap["cookie_garble"], attrMap["cookie_domain"], attrMap["cookie_path"], attrMap["cookie_httponly"], attrMap["cookie_secure"], cookieTime)
}

func testAccNsxtPolicyLBCookiePersistenceProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_cookie_persistence_profile" "test" {
  display_name = "%s"
}`, accTestPolicyLBCookiePersistenceProfileUpdateAttributes["display_name"])
}
