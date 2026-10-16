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

var accTestVpcLBSourceIpPersistenceProfileCreateAttributes = map[string]string{
	"description":                      "terraform created",
	"persistence_shared":               "true",
	"purge":                            "FULL",
	"timeout":                          "100",
	"ha_persistence_mirroring_enabled": "true",
}

var accTestVpcLBSourceIpPersistenceProfileUpdateAttributes = map[string]string{
	"description":                      "terraform updated",
	"persistence_shared":               "false",
	"purge":                            "NO_PURGE",
	"timeout":                          "200",
	"ha_persistence_mirroring_enabled": "false",
}

func TestAccResourceNsxtVpcLBSourceIpPersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_source_ip_persistence_profile.test"
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
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_source_ip_persistence_profile", updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBSourceIpPersistenceProfileTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBSourceIpPersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestVpcLBSourceIpPersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "purge", accTestVpcLBSourceIpPersistenceProfileCreateAttributes["purge"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBSourceIpPersistenceProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestVpcLBSourceIpPersistenceProfileCreateAttributes["ha_persistence_mirroring_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBSourceIpPersistenceProfileTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPersistenceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBSourceIpPersistenceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestVpcLBSourceIpPersistenceProfileUpdateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "purge", accTestVpcLBSourceIpPersistenceProfileUpdateAttributes["purge"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBSourceIpPersistenceProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestVpcLBSourceIpPersistenceProfileUpdateAttributes["ha_persistence_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBSourceIpPersistenceProfileMinimalistic(updateName),
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

func TestAccResourceNsxtVpcLBSourceIpPersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_source_ip_persistence_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPersistenceProfileCheckDestroy(state, "nsxt_vpc_lb_source_ip_persistence_profile", name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBSourceIpPersistenceProfileMinimalistic(name),
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

func testAccNsxtVpcLBSourceIpPersistenceProfileTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBSourceIpPersistenceProfileUpdateAttributes
	if createFlow {
		attrMap = accTestVpcLBSourceIpPersistenceProfileCreateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_source_ip_persistence_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"

  persistence_shared               = %s
  purge                            = "%s"
  timeout                          = %s
  ha_persistence_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["persistence_shared"], attrMap["purge"], attrMap["timeout"], attrMap["ha_persistence_mirroring_enabled"])
}

func testAccNsxtVpcLBSourceIpPersistenceProfileMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_source_ip_persistence_profile" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}
