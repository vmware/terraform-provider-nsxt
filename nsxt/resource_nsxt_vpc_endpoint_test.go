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

var accTestVpcEndpointCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
}

var accTestVpcEndpointUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
}

func TestAccResourceNsxtVpcEndpoint_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_endpoint.test"
	testAccOnlyVPC(t)
	testAccNSXVersion(t, "9.2.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcEndpointCheckDestroy(state, accTestVpcEndpointUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcEndpointTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcEndpointExists(accTestVpcEndpointCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcEndpointCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcEndpointCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "vpc_service_endpoint"),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_allocation_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcEndpointTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcEndpointExists(accTestVpcEndpointUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcEndpointUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcEndpointUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "vpc_service_endpoint"),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_allocation_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcEndpointMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcEndpointExists(accTestVpcEndpointCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcEndpoint_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_endpoint.test"
	testAccOnlyVPC(t)
	testAccNSXVersion(t, "9.2.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcEndpointCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcEndpointMinimalistic(),
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

func testAccNsxtVpcEndpointExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcEndpoint resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcEndpoint resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcEndpointExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcEndpoint %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcEndpointCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_endpoint" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcEndpointExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcEndpoint %s still exists", displayName)
		}
	}
	return nil
}

// testAccNsxtVpcEndpointDeps creates a service endpoint and IP allocation
// that the VpcEndpoint will reference. Both are created in the test VPC.
func testAccNsxtVpcEndpointDeps() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_service_endpoint" "dep" {
  %s
  display_name        = "%s-svc"
  service_endpoint_ip = "192.168.100.10"
}

resource "nsxt_vpc_ip_address_allocation" "dep" {
  %s
  display_name    = "%s-alloc"
  allocation_size = 1
}
`, testAccNsxtPolicyMultitenancyContext(), accTestVpcEndpointCreateAttributes["display_name"],
		testAccNsxtPolicyMultitenancyContext(), accTestVpcEndpointCreateAttributes["display_name"])
}

func testAccNsxtVpcEndpointTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcEndpointCreateAttributes
	} else {
		attrMap = accTestVpcEndpointUpdateAttributes
	}
	return testAccNsxtVpcEndpointDeps() + fmt.Sprintf(`
resource "nsxt_vpc_endpoint" "test" {
  %s
  display_name         = "%s"
  description          = "%s"
  vpc_service_endpoint = nsxt_vpc_service_endpoint.dep.path
  ip_allocation_path   = nsxt_vpc_ip_address_allocation.dep.path
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"])
}

func testAccNsxtVpcEndpointMinimalistic() string {
	return testAccNsxtVpcEndpointDeps() + fmt.Sprintf(`
resource "nsxt_vpc_endpoint" "test" {
  %s
  display_name         = "%s"
  vpc_service_endpoint = nsxt_vpc_service_endpoint.dep.path
  ip_allocation_path   = nsxt_vpc_ip_address_allocation.dep.path
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcEndpointUpdateAttributes["display_name"])
}
