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

var accTestVpcLBServiceCreateAttributes = map[string]string{
	"description":     "terraform created",
	"enabled":         "true",
	"error_log_level": "ERROR",
	"size":            "SMALL",
}

var accTestVpcLBServiceUpdateAttributes = map[string]string{
	"description":     "terraform updated",
	"enabled":         "false",
	"error_log_level": "EMERGENCY",
	"size":            "SMALL",
}

func TestAccResourceNsxtVpcLBService_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_service.test"
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
			return testAccNsxtVpcLBServiceCheckDestroy(state, updateName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBServiceTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBServiceCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestVpcLBServiceCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "error_log_level", accTestVpcLBServiceCreateAttributes["error_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "size", accTestVpcLBServiceCreateAttributes["size"]),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", os.Getenv("NSXT_VPC_PROJECT_ID")),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBServiceTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServiceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBServiceUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestVpcLBServiceUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "error_log_level", accTestVpcLBServiceUpdateAttributes["error_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "size", accTestVpcLBServiceUpdateAttributes["size"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBServiceMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBServiceExists(testResourceName),
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

func TestAccResourceNsxtVpcLBService_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBServiceCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBServiceMinimalistic(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		})),
	})
}

func testAccNsxtVpcLBServiceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPC LBService resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VPC LBService resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBServiceExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VPC LBService %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtVpcLBServiceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_lb_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBServiceExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("VPC LBService %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBServiceTemplate(createFlow bool, name string) string {
	attrMap := accTestVpcLBServiceCreateAttributes
	if !createFlow {
		attrMap = accTestVpcLBServiceUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_service" "test" {
  %s
  display_name    = "%s"
  description     = "%s"
  enabled         = %s
  error_log_level = "%s"
  size            = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["enabled"], attrMap["error_log_level"], attrMap["size"])
}

func testAccNsxtVpcLBServiceMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_service" "test" {
  %s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}
