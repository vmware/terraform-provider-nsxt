package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyTier0GatewayInterface_basic(t *testing.T) {
	name := getAccTestResourceName()
	mtu := "1500"
	subnet := "1.1.12.2/24"
	testResourceName := "data.nsxt_policy_gateway_interface.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceDataSourceTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyTier1GatewayInterface_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_gateway_interface.test2"
	mtu := "1500"
	subnet := "1.1.12.2/24"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceDataSourceTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplate(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier0InterfaceServiceTemplate(name, subnet, mtu) + fmt.Sprintf(`

data "nsxt_policy_gateway_interface" "test1" {
    display_name = "%s"
    gateway_path = nsxt_policy_tier0_gateway.test.path
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

`, name)
}

func testAccNsxtPolicyTier1InterfaceDataSourceTemplate(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyTier1InterfaceTemplate(name, subnet, mtu, false) + fmt.Sprintf(`

data "nsxt_policy_gateway_interface" "test2" {
    display_name = "%s"
    gateway_path = nsxt_policy_tier1_gateway.test.path
	depends_on = [nsxt_policy_tier1_gateway_interface.test]
}

`, name)
}
