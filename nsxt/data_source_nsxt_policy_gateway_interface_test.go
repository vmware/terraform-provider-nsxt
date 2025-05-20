package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyTier0GatewayInterface(t *testing.T) {
	name := getAccTestResourceName()
	mtu := "1500"
	subnet := "1.1.12.2/24"
	testResourceName := "data.nsxt_policy_gateway_interface.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithoutGatewayOrService(name, subnet, mtu),
				ExpectError: regexp.MustCompile(`One of gateway_path or service_path should be set`),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithGatewayPath(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithServicePath(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyTier1GatewayInterface(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_gateway_interface.test2"
	mtu := "1500"
	subnet := "1.1.12.2/24"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceDataSourceTemplateWithGatewayPath(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceDataSourceTemplateWithServicePath(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithoutGatewayOrService(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier0InterfaceServiceTemplate(name, subnet, mtu) + fmt.Sprintf(`

data "nsxt_policy_gateway_interface" "test1" {
    display_name = "%s"
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

`, name)
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithGatewayPath(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier0InterfaceServiceTemplate(name, subnet, mtu) + fmt.Sprintf(`

data "nsxt_policy_gateway_interface" "test1" {
    display_name = "%s"
    gateway_path = nsxt_policy_tier0_gateway.test.path
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

`, name)
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplateWithServicePath(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier0InterfaceServiceTemplate(name, subnet, mtu) + fmt.Sprintf(`

data "nsxt_policy_gateway_locale_service" "test" {
  gateway_path = nsxt_policy_tier0_gateway.test.path
}

data "nsxt_policy_gateway_interface" "test1" {
    display_name = "%s"
    service_path = data.nsxt_policy_gateway_locale_service.test.path
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

`, name)
}

func testAccNsxtPolicyTier1InterfaceDataSourceTemplateWithGatewayPath(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier1InterfaceTemplate(name, subnet, mtu, false) + fmt.Sprintf(`

data "nsxt_policy_gateway_interface" "test2" {
    display_name = "%s"
    gateway_path = nsxt_policy_tier1_gateway.test.path
	depends_on = [nsxt_policy_tier1_gateway_interface.test]
}

`, name)
}

func testAccNsxtPolicyTier1InterfaceDataSourceTemplateWithServicePath(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier1InterfaceTemplate(name, subnet, mtu, false) + fmt.Sprintf(`

data "nsxt_policy_gateway_locale_service" "test" {
  gateway_path = nsxt_policy_tier1_gateway.test.path
}

data "nsxt_policy_gateway_interface" "test2" {
    display_name = "%s"
    service_path = data.nsxt_policy_gateway_locale_service.test.path
	depends_on = [nsxt_policy_tier1_gateway_interface.test]
}

`, name)
}
