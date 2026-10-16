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

var accTestVpcLBCookiePersistenceProfileCreateAttributes = map[string]string{
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

var accTestVpcLBCookiePersistenceProfileUpdateAttributes = map[string]string{
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

func testAccNsxtVpcLBCookiePersistenceProfileChecks(testResourceName, name string, attrs map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
		resource.TestCheckResourceAttr(testResourceName, "display_name", name),
		resource.TestCheckResourceAttr(testResourceName, "description", attrs["description"]),
		resource.TestCheckResourceAttr(testResourceName, "persistence_shared", attrs["persistence_shared"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_name", attrs["cookie_name"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_mode", attrs["cookie_mode"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_fallback", attrs["cookie_fallback"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_garble", attrs["cookie_garble"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_domain", attrs["cookie_domain"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_path", attrs["cookie_path"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_httponly", attrs["cookie_httponly"]),
		resource.TestCheckResourceAttr(testResourceName, "cookie_secure", attrs["cookie_secure"]),
		resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
		resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
		resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
		resource.TestCheckResourceAttrSet(testResourceName, "path"),
		resource.TestCheckResourceAttrSet(testResourceName, "revision"),
		resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
	)
}

func TestAccResourceNsxtVpcLBCookiePersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_cookie_persistence_profile.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_cookie_persistence_profile", updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBCookiePersistenceProfileSessionCookieTimeTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBCookiePersistenceProfileChecks(testResourceName, createName, accTestVpcLBCookiePersistenceProfileCreateAttributes),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.0.max_idle", accTestVpcLBCookiePersistenceProfileCreateAttributes["max_idle"]),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.0.max_life", accTestVpcLBCookiePersistenceProfileCreateAttributes["max_life"]),
				),
			},
			{
				Config: testAccNsxtVpcLBCookiePersistenceProfilePersistenceCookieTimeTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBCookiePersistenceProfileChecks(testResourceName, createName, accTestVpcLBCookiePersistenceProfileCreateAttributes),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_cookie_time.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_cookie_time.0.max_idle", accTestVpcLBCookiePersistenceProfileCreateAttributes["max_idle"]),
				),
			},
			{
				Config: testAccNsxtVpcLBCookiePersistenceProfileTemplate(false, updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBCookiePersistenceProfileChecks(testResourceName, updateName, accTestVpcLBCookiePersistenceProfileUpdateAttributes),
					resource.TestCheckResourceAttr(testResourceName, "session_cookie_time.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_cookie_time.#", "0"),
				),
			},
			{
				Config: testAccNsxtVpcLBCookiePersistenceProfileMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
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

func TestAccResourceNsxtVpcLBCookiePersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_cookie_persistence_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_cookie_persistence_profile", name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBCookiePersistenceProfileMinimalistic(name),
			},
			{
				// VPC objects can only be imported by policy path
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtVpcLBPersistenceProfileExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPC LB Persistence Profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VPC LB Persistence Profile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBPersistenceProfileExistsCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VPC LB Persistence Profile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBPersistenceProfileCheckDestroy(state *terraform.State, resourceType string, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != resourceType {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBPersistenceProfileExistsCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("VPC LB Persistence Profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBCookiePersistenceProfileSessionCookieTimeTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBCookiePersistenceProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBCookiePersistenceProfileCreateAttributes
	}
	cookieTime := fmt.Sprintf(`
  session_cookie_time {
    max_idle = %s
    max_life = %s
  }`, attrMap["max_idle"], attrMap["max_life"])

	return testAccNsxtVpcLBCookiePersistenceProfileTemplate(createFlow, name, cookieTime)
}

func testAccNsxtVpcLBCookiePersistenceProfilePersistenceCookieTimeTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBCookiePersistenceProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBCookiePersistenceProfileCreateAttributes
	}
	cookieTime := fmt.Sprintf(`
  persistence_cookie_time {
    max_idle = %s
  }`, attrMap["max_idle"])

	return testAccNsxtVpcLBCookiePersistenceProfileTemplate(createFlow, name, cookieTime)
}

func testAccNsxtVpcLBCookiePersistenceProfileTemplate(createFlow bool, name string, cookieTime string) string {
	attrMap := accTestVpcLBCookiePersistenceProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBCookiePersistenceProfileCreateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_cookie_persistence_profile" "test" {
%s
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
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["persistence_shared"], attrMap["cookie_name"], attrMap["cookie_mode"], attrMap["cookie_fallback"], attrMap["cookie_garble"], attrMap["cookie_domain"], attrMap["cookie_path"], attrMap["cookie_httponly"], attrMap["cookie_secure"], cookieTime)
}

func testAccNsxtVpcLBCookiePersistenceProfileMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_cookie_persistence_profile" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}
