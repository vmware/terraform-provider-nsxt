// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestPolicyConstraintCreateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform created",
	"target_owner_type":    "GM",
	"message":              "test-create",
	"target_resource_type": "StaticRoutes",
	"count":                "10",
}

var accTestPolicyConstraintUpdateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform updated",
	"target_owner_type":    "LM",
	"message":              "test-update",
	"target_resource_type": "Infra.TlsCertificate",
	"count":                "100",
}

func TestAccResourceNsxtPolicyConstraint_basic(t *testing.T) {
	testResourceName := "nsxt_policy_constraint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "9.0.0"); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyConstraintCheckDestroy(state, accTestPolicyConstraintUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyConstraintTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConstraintExists(accTestPolicyConstraintCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConstraintCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConstraintCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "message", accTestPolicyConstraintCreateAttributes["message"]),
					resource.TestCheckResourceAttr(testResourceName, "target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "target.0.path_prefix"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.0.target_resource_type", accTestPolicyConstraintCreateAttributes["target_resource_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyConstraintTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConstraintExists(accTestPolicyConstraintUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConstraintUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConstraintUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "message", accTestPolicyConstraintUpdateAttributes["message"]),
					resource.TestCheckResourceAttr(testResourceName, "target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "target.0.path_prefix"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.0.target_resource_type", accTestPolicyConstraintUpdateAttributes["target_resource_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyConstraintMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConstraintExists(accTestPolicyConstraintCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyConstraint_vpc(t *testing.T) {
	testResourceName := "nsxt_policy_constraint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyConstraintCheckDestroy(state, accTestPolicyConstraintUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyConstraintVpcTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConstraintExists(accTestPolicyConstraintCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConstraintCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConstraintCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "message", accTestPolicyConstraintCreateAttributes["message"]),
					resource.TestCheckResourceAttr(testResourceName, "target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "target.0.path_prefix"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.0.target_resource_type", accTestPolicyConstraintCreateAttributes["target_resource_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyConstraintVpcTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConstraintExists(accTestPolicyConstraintUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConstraintUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConstraintUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "message", accTestPolicyConstraintUpdateAttributes["message"]),
					resource.TestCheckResourceAttr(testResourceName, "target.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "target.0.path_prefix"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_count.0.target_resource_type", accTestPolicyConstraintUpdateAttributes["target_resource_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyConstraint_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_constraint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyConstraintCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyConstraintMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyConstraintExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Constraint resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Constraint resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyConstraintExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy Constraint %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyConstraintCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_constraint" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyConstraintExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy Constraint %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyConstraintTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyConstraintCreateAttributes
	} else {
		attrMap = accTestPolicyConstraintUpdateAttributes
	}
	return testAccNsxtPolicyProjectTemplate900(true, true, false, false) + fmt.Sprintf(`
resource "nsxt_policy_constraint" "test" {
  display_name = "%s"
  description  = "%s"
  message      = "%s"

  target {
    path_prefix = nsxt_policy_project.test.path
  }

  instance_count {
    count                = "%s"
    target_resource_type = "%s"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["message"], attrMap["count"], attrMap["target_resource_type"])
}

func testAccNsxtPolicyConstraintVpcTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyConstraintCreateAttributes
	} else {
		attrMap = accTestPolicyConstraintUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_vpc" "test" {
  %s
  id = "%s"
}

resource "nsxt_policy_constraint" "test" {
  %s
  display_name = "%s"
  description  = "%s"
  message      = "%s"

  target {
    path_prefix = "${data.nsxt_vpc.test.path}/"
  }

  instance_count {
    count                = "%s"
    target_resource_type = "%s"
  }
}`, testAccNsxtProjectContext(), os.Getenv("NSXT_VPC_ID"), testAccNsxtProjectContext(), attrMap["display_name"], attrMap["description"], attrMap["message"], attrMap["count"], attrMap["target_resource_type"])
}

func testAccNsxtPolicyConstraintMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_constraint" "test" {
  display_name = "%s"
}`, accTestPolicyConstraintUpdateAttributes["display_name"])
}
