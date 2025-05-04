package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceNsxtPolicyTier0GatewayInterface_basic(t *testing.T) {
	interfaceName := getAccTestDataSourceName()
	gatewayName := "t0testgw"
	testResourceName := "data.nsxt_policy_tier0_gateway_interface.sample"
	transportZoneName := getOverlayTransportZoneName()
	interfaceDescription := "Acceptance Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, interfaceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceDataSourceTemplate(interfaceName, gatewayName, transportZoneName, interfaceDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", interfaceName),
					resource.TestCheckResourceAttr(testResourceName, "description", interfaceDescription),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplate(interfaceName string, gatewayName string, transportZoneName string, interfaceDescription string) string {
	return CreateT0Gateway(gatewayName) + CreateSegment(transportZoneName) + CreateT0GatewayInterface(interfaceName, interfaceDescription) + fmt.Sprintf(`

data "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

data "nsxt_policy_tier0_gateway_interface" "sample" {
    display_name = "%s"
    t0_gateway_path = data.nsxt_policy_tier0_gateway.test.path
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}
`, gatewayName, interfaceName)
}

func CreateT0Gateway(gatewayName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}
`, gatewayName)
}

func CreateSegment(transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "overlay_transport_zone" {
  display_name   = "%s"
}

resource "nsxt_policy_segment" "segment1" {
  display_name        = "segment-acc-test"
  description         = "Terraform provisioned Segment"
  transport_zone_path = data.nsxt_policy_transport_zone.overlay_transport_zone.path

}
`, transportZoneName)

}

func CreateT0GatewayInterface(interfaceName string, interfaceDescription string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name = "%s"
  description  = "%s"
  type         = "SERVICE"
  mtu          = 1500
  gateway_path = nsxt_policy_tier0_gateway.test.path
  segment_path = nsxt_policy_segment.segment1.path
  subnets      = ["1.1.12.2/24"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [nsxt_policy_tier0_gateway.test]
}
`, interfaceName, interfaceDescription)
}

func CreateRealization() string {
	return `
data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_interface.test.path
  depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

data "nsxt_policy_gateway_interface_realization" "gw_realization" {
  gateway_path = nsxt_policy_tier0_gateway_interface.test.path
  depends_on = [nsxt_policy_tier0_gateway_interface.test]
}
`
}
