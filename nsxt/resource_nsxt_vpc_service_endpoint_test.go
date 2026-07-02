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

var accTestVpcServiceEndpointCreateAttributes = map[string]string{
	"display_name":             getAccTestResourceName(),
	"description":              "terraform created",
	"service_endpoint_ip":      "192.168.100.10",
	"service_endpoint_ip_type": "WORKLOAD",
}

var accTestVpcServiceEndpointUpdateAttributes = map[string]string{
	"display_name":             getAccTestResourceName(),
	"description":              "terraform updated",
	"service_endpoint_ip":      "192.168.100.10",
	"service_endpoint_ip_type": "WORKLOAD",
}

func TestAccResourceNsxtVpcServiceEndpoint_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_service_endpoint.test"
	testAccOnlyVPC(t)
	testAccNSXVersion(t, "9.2.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceEndpointCheckDestroy(state, accTestVpcServiceEndpointUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceEndpointTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(accTestVpcServiceEndpointCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcServiceEndpointCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcServiceEndpointCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip", accTestVpcServiceEndpointCreateAttributes["service_endpoint_ip"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip_type", accTestVpcServiceEndpointCreateAttributes["service_endpoint_ip_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceEndpointTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(accTestVpcServiceEndpointUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcServiceEndpointUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcServiceEndpointUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip", accTestVpcServiceEndpointUpdateAttributes["service_endpoint_ip"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip_type", accTestVpcServiceEndpointUpdateAttributes["service_endpoint_ip_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceEndpointMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(accTestVpcServiceEndpointCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcServiceEndpoint_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_service_endpoint.test"
	testAccOnlyVPC(t)
	testAccNSXVersion(t, "9.2.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceEndpointCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceEndpointMinimalistic(),
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

func testAccNsxtVpcServiceEndpointExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcServiceEndpoint resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcServiceEndpoint resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcServiceEndpointExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcServiceEndpoint %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcServiceEndpointCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_service_endpoint" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcServiceEndpointExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcServiceEndpoint %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcServiceEndpointTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcServiceEndpointCreateAttributes
	} else {
		attrMap = accTestVpcServiceEndpointUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_service_endpoint" "test" {
  %s
  display_name           = "%s"
  description            = "%s"
  service_endpoint_ip    = "%s"
  service_endpoint_ip_type = "%s"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"],
		attrMap["service_endpoint_ip"], attrMap["service_endpoint_ip_type"])
}

func testAccNsxtVpcServiceEndpointMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_service_endpoint" "test" {
  %s
  display_name        = "%s"
  service_endpoint_ip = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcServiceEndpointUpdateAttributes["display_name"],
		accTestVpcServiceEndpointUpdateAttributes["service_endpoint_ip"])
}
