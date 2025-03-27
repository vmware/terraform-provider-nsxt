// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var entry1UUID = newUUID()

var accTestPolicyL7AccessProfileCreateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform created",
	"default_action_logged": "true",
	"default_action":        "ALLOW",

	"entry1_id":              entry1UUID,
	"entry1_action":          "ALLOW",
	"entry1_disabled":        "true",
	"entry1_logged":          "true",
	"entry1_sequence_number": "100",

	"attr1key1":    "APP_ID",
	"attr1value1":  "SSL",
	"attr1source1": "SYSTEM",

	"attr1key2":   "TLS_CIPHER_SUITE",
	"attr1value2": "TLS_RSA_EXPORT_WITH_RC4_40_MD5",

	"entry2_display_name":    "terraform-created",
	"entry2_sequence_number": "200",
	"entry2_action":          "ALLOW",
	"attr2key1":              "URL_CATEGORY",
	"attr2value1":            "Abused Drugs",
	"attr2source1":           "SYSTEM",
}

var accTestPolicyL7AccessProfileUpdateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform updated",
	"default_action_logged": "false",
	"default_action":        "REJECT",

	"entry1_id":              entry1UUID,
	"entry1_action":          "REJECT",
	"entry1_disabled":        "false",
	"entry1_logged":          "false",
	"entry1_sequence_number": "100",

	"attr1key1":    "APP_ID",
	"attr1value1":  "SSL",
	"attr1source1": "SYSTEM",

	"attr1key2":   "TLS_CIPHER_SUITE",
	"attr1value2": "TLS_RSA_EXPORT_WITH_RC4_40_MD5",

	"entry2_display_name":    "terraform-created",
	"entry2_sequence_number": "200",
	"entry2_action":          "REJECT",
	"attr2key1":              "URL_CATEGORY",
	"attr2value1":            "Auctions",
	"attr2source1":           "SYSTEM",
}

func TestAccResourceNsxtPolicyL7AccessProfile_basic(t *testing.T) {
	testAccResourceNsxtPolicyL7AccessProfileBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyL7AccessProfile_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyL7AccessProfileBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyL7AccessProfileBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_l7_access_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyL7AccessProfileCheckDestroy(state, accTestPolicyL7AccessProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyL7AccessProfileTemplate(true, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyL7AccessProfileExists(accTestPolicyL7AccessProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyL7AccessProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyL7AccessProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "default_action_logged", accTestPolicyL7AccessProfileCreateAttributes["default_action_logged"]),
					resource.TestCheckResourceAttr(testResourceName, "default_action", accTestPolicyL7AccessProfileCreateAttributes["default_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.action", accTestPolicyL7AccessProfileCreateAttributes["entry1_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.disabled", accTestPolicyL7AccessProfileCreateAttributes["entry1_disabled"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.logged", accTestPolicyL7AccessProfileCreateAttributes["entry1_logged"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.sequence_number", accTestPolicyL7AccessProfileCreateAttributes["entry1_sequence_number"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.nsx_id", accTestPolicyL7AccessProfileCreateAttributes["entry1_id"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.key", accTestPolicyL7AccessProfileCreateAttributes["attr1key1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.values.0", accTestPolicyL7AccessProfileCreateAttributes["attr1value1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.0.key", accTestPolicyL7AccessProfileCreateAttributes["attr1key2"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.0.values.0", accTestPolicyL7AccessProfileCreateAttributes["attr1value2"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.0.key", accTestPolicyL7AccessProfileCreateAttributes["attr2key1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.0.values.0", accTestPolicyL7AccessProfileCreateAttributes["attr2value1"]),

					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.display_name", accTestPolicyL7AccessProfileCreateAttributes["entry2_display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.action", accTestPolicyL7AccessProfileCreateAttributes["entry2_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.sequence_number", accTestPolicyL7AccessProfileCreateAttributes["entry2_sequence_number"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyL7AccessProfileTemplate(false, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyL7AccessProfileExists(accTestPolicyL7AccessProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyL7AccessProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyL7AccessProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "default_action_logged", accTestPolicyL7AccessProfileUpdateAttributes["default_action_logged"]),
					resource.TestCheckResourceAttr(testResourceName, "default_action", accTestPolicyL7AccessProfileUpdateAttributes["default_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.action", accTestPolicyL7AccessProfileUpdateAttributes["entry1_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.disabled", accTestPolicyL7AccessProfileUpdateAttributes["entry1_disabled"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.logged", accTestPolicyL7AccessProfileUpdateAttributes["entry1_logged"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.sequence_number", accTestPolicyL7AccessProfileUpdateAttributes["entry1_sequence_number"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.nsx_id", accTestPolicyL7AccessProfileUpdateAttributes["entry1_id"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.key", accTestPolicyL7AccessProfileUpdateAttributes["attr1key1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.values.0", accTestPolicyL7AccessProfileUpdateAttributes["attr1value1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.0.key", accTestPolicyL7AccessProfileUpdateAttributes["attr1key2"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.0.attribute.0.sub_attribute.0.values.0", accTestPolicyL7AccessProfileUpdateAttributes["attr1value2"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.0.key", accTestPolicyL7AccessProfileUpdateAttributes["attr2key1"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.0.values.0", accTestPolicyL7AccessProfileUpdateAttributes["attr2value1"]),

					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.display_name", accTestPolicyL7AccessProfileUpdateAttributes["entry2_display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.action", accTestPolicyL7AccessProfileUpdateAttributes["entry1_action"]),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.attribute.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "l7_access_entry.1.sequence_number", accTestPolicyL7AccessProfileUpdateAttributes["entry2_sequence_number"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyL7AccessProfileMinimalistic(withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyL7AccessProfileExists(accTestPolicyL7AccessProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyL7AccessProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_l7_access_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyL7AccessProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyL7AccessProfileMinimalistic(false),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyL7AccessProfile_importBasic_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_l7_access_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyL7AccessProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyL7AccessProfileMinimalistic(true),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyL7AccessProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy L7AccessProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy L7AccessProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyL7AccessProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy L7AccessProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyL7AccessProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_l7_access_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyL7AccessProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy L7AccessProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyL7AccessProfileTemplate(createFlow bool, withContext bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyL7AccessProfileCreateAttributes
	} else {
		attrMap = accTestPolicyL7AccessProfileUpdateAttributes
	}
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_l7_access_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  default_action_logged = %s
  default_action = "%s"

  l7_access_entry {
    nsx_id = "%s"
    action = "%s"

    attribute {
      key = "%s"
      values = ["%s"]
      attribute_source = "%s"

      sub_attribute {
        key = "%s"
        values = ["%s"]
      }
    }

    disabled = %s
    logged = "%s"
    sequence_number = "%s"
  }

  l7_access_entry {
     display_name = "%s"
     action = "%s"

    attribute {
      key = "%s"
      values = ["%s"]
      attribute_source = "%s"
    }
    sequence_number = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, context, attrMap["display_name"], attrMap["description"], attrMap["default_action_logged"], attrMap["default_action"],
		attrMap["entry1_id"], attrMap["entry1_action"], attrMap["attr1key1"], attrMap["attr1value1"], attrMap["attr1source1"],
		attrMap["attr1key2"], attrMap["attr1value2"], attrMap["entry1_disabled"], attrMap["entry1_logged"], attrMap["entry1_sequence_number"],
		attrMap["entry2_display_name"], attrMap["entry2_action"], attrMap["attr2key1"], attrMap["attr2value1"], attrMap["attr2source1"], attrMap["entry2_sequence_number"])
}

func testAccNsxtPolicyL7AccessProfileMinimalistic(withContext bool) string {
	// Minimalistic profile should have at least one entry
	attrMap := accTestPolicyL7AccessProfileUpdateAttributes
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_l7_access_profile" "test" {
%s
  display_name = "%s"
  default_action = "ALLOW"

  l7_access_entry {
     display_name = "%s"
     action = "%s"

    attribute {
      key = "%s"
      values = ["%s"]
      attribute_source = "%s"
    }
    sequence_number = "%s"
  }

}`, context, accTestPolicyL7AccessProfileUpdateAttributes["display_name"], attrMap["entry2_display_name"], attrMap["entry2_action"], attrMap["attr2key1"], attrMap["attr2value1"], attrMap["attr2source1"], attrMap["entry2_sequence_number"])
}
