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
	"description":              "terraform created",
	"service_endpoint_ip_type": "WORKLOAD",
}

var accTestVpcServiceEndpointUpdateAttributes = map[string]string{
	"description":              "terraform updated",
	"service_endpoint_ip_type": "WORKLOAD",
}

func TestAccResourceNsxtVpcServiceEndpoint_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_service_endpoint.test"
	testAccOnlyVPC(t)
	testAccNSXVersion(t, "9.2.0")
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()
	subnetName := getAccTestResourceName()
	cidr := "192.168.100.0/24"
	ip := "192.168.100.10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceEndpointCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceEndpointTemplate(true, createName, subnetName, cidr, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(createName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcServiceEndpointCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip", ip),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip_type", accTestVpcServiceEndpointCreateAttributes["service_endpoint_ip_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceEndpointTemplate(false, updateName, subnetName, cidr, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcServiceEndpointUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip", ip),
					resource.TestCheckResourceAttr(testResourceName, "service_endpoint_ip_type", accTestVpcServiceEndpointUpdateAttributes["service_endpoint_ip_type"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcServiceEndpointMinimalistic(updateName, subnetName, cidr, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcServiceEndpointExists(updateName, testResourceName),
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
	subnetName := getAccTestResourceName()
	cidr := "192.168.101.0/24"
	ip := "192.168.101.10"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcServiceEndpointCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcServiceEndpointMinimalistic(name, subnetName, cidr, ip),
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

// testAccNsxtVpcServiceEndpointSubnetDeps creates a subnet in the test VPC so
// that the service_endpoint_ip used by the dependent nsxt_vpc_service_endpoint
// falls within a subnet that actually exists in the VPC.
func testAccNsxtVpcServiceEndpointSubnetDeps(subnetName, cidr string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_subnet" "dep" {
  %s
  display_name = "%s-subnet"
  ip_addresses = ["%s"]
  access_mode  = "Isolated"
}
`, testAccNsxtPolicyMultitenancyContext(), subnetName, cidr)
}

func testAccNsxtVpcServiceEndpointTemplate(createFlow bool, displayName, subnetName, cidr, ip string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcServiceEndpointCreateAttributes
	} else {
		attrMap = accTestVpcServiceEndpointUpdateAttributes
	}
	return testAccNsxtVpcServiceEndpointSubnetDeps(subnetName, cidr) + fmt.Sprintf(`
resource "nsxt_vpc_service_endpoint" "test" {
  %s
  display_name             = "%s"
  description              = "%s"
  service_endpoint_ip      = "%s"
  service_endpoint_ip_type = "%s"
  depends_on               = [nsxt_vpc_subnet.dep]
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), displayName, attrMap["description"],
		ip, attrMap["service_endpoint_ip_type"])
}

func testAccNsxtVpcServiceEndpointMinimalistic(displayName, subnetName, cidr, ip string) string {
	return testAccNsxtVpcServiceEndpointSubnetDeps(subnetName, cidr) + fmt.Sprintf(`
resource "nsxt_vpc_service_endpoint" "test" {
  %s
  display_name        = "%s"
  service_endpoint_ip = "%s"
  depends_on          = [nsxt_vpc_subnet.dep]
}`, testAccNsxtPolicyMultitenancyContext(), displayName, ip)
}
