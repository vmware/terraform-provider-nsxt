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

// testAccNsxtVpcLBPersistenceProfileExists and testAccNsxtVpcLBPersistenceProfileCheckDestroy
// are defined in resource_nsxt_vpc_lb_cookie_persistence_profile_test.go

var accTestVpcLBGenericPersistenceProfileCreateAttributes = map[string]string{
	"description":                      "terraform created",
	"persistence_shared":               "true",
	"timeout":                          "100",
	"ha_persistence_mirroring_enabled": "true",
}

var accTestVpcLBGenericPersistenceProfileUpdateAttributes = map[string]string{
	"description":                      "terraform updated",
	"persistence_shared":               "false",
	"timeout":                          "200",
	"ha_persistence_mirroring_enabled": "false",
}

func TestAccResourceNsxtVpcLBGenericPersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_generic_persistence_profile.test"
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
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_generic_persistence_profile", updateName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBGenericPersistenceProfileTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBGenericPersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestVpcLBGenericPersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBGenericPersistenceProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestVpcLBGenericPersistenceProfileCreateAttributes["ha_persistence_mirroring_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBGenericPersistenceProfileTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBGenericPersistenceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestVpcLBGenericPersistenceProfileUpdateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBGenericPersistenceProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestVpcLBGenericPersistenceProfileUpdateAttributes["ha_persistence_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBGenericPersistenceProfileMinimalistic(updateName),
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
		}),
	})
}

func TestAccResourceNsxtVpcLBGenericPersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_generic_persistence_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_generic_persistence_profile", name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBGenericPersistenceProfileMinimalistic(name),
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

func testAccNsxtVpcLBGenericPersistenceProfileTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBGenericPersistenceProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBGenericPersistenceProfileCreateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_generic_persistence_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"

  persistence_shared               = %s
  timeout                          = %s
  ha_persistence_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["persistence_shared"], attrMap["timeout"], attrMap["ha_persistence_mirroring_enabled"])
}

func testAccNsxtVpcLBGenericPersistenceProfileMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_generic_persistence_profile" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}
