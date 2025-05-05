package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyGatewayInterface_basic(t *testing.T) {
	t0InterfaceName := getAccTestDataSourceName()
	t1InterfaceName := getAccTestDataSourceName()
	t0GatewayName := "t0testgw"
	t1GatewayName := "t1testgw"
	t0TestResourceName := "data.nsxt_policy_gateway_interface.test1"
	t1TestResourceName := "data.nsxt_policy_gateway_interface.test2"
	transportZoneName := getOverlayTransportZoneName()
	interfaceDescription := "Acceptance Test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: func(state *terraform.State) error {
		// 	return testAccNsxtPolicyTier0InterfaceCheckDestroy(state, t0InterfaceName)
		// },
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0InterfaceDataSourceTemplate(t0InterfaceName, t1InterfaceName, t0GatewayName, t1GatewayName, transportZoneName, interfaceDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(t0TestResourceName, "display_name", t0InterfaceName),
					resource.TestCheckResourceAttr(t0TestResourceName, "description", interfaceDescription),
					resource.TestCheckResourceAttrSet(t0TestResourceName, "path"),
					resource.TestCheckResourceAttr(t1TestResourceName, "display_name", t1InterfaceName),
					resource.TestCheckResourceAttr(t1TestResourceName, "description", interfaceDescription),
					resource.TestCheckResourceAttrSet(t1TestResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTier0InterfaceDataSourceTemplate(t0InterfaceName string, t1InterfaceName string, t0GatewayName string, t1GatewayName string, transportZoneName string, interfaceDescription string) string {
	return CreateT0Gateway(t0GatewayName) + CreateSegment(transportZoneName) + CreateT0GatewayInterface(t0InterfaceName, interfaceDescription) + CreateT1Gateway(t1GatewayName) + CreateT1GatewayInterface(t1InterfaceName, interfaceDescription) + fmt.Sprintf(`

data "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

data "nsxt_policy_gateway_interface" "test1" {
    display_name = "%s"
    gateway_path = data.nsxt_policy_tier0_gateway.test.path
	depends_on = [nsxt_policy_tier0_gateway_interface.test]
}

data "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  depends_on = [nsxt_policy_tier1_gateway_interface.test]
}

data "nsxt_policy_gateway_interface" "test2" {
    display_name = "%s"
    gateway_path = data.nsxt_policy_tier1_gateway.test.path
	depends_on = [nsxt_policy_tier1_gateway_interface.test]
}
`, t0GatewayName, t0InterfaceName, t1GatewayName, t1InterfaceName)
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

func CreateRealizationT0() string {
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

func CreateT1Gateway(gatewayName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}


`, gatewayName)
}

func CreateT1GatewayInterface(interfaceName string, interfaceDescription string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_tier1_gateway_interface" "test" {
  display_name = "%s"
  description  = "%s"
  mtu          = 1500
  gateway_path = nsxt_policy_tier1_gateway.test.path
  segment_path = nsxt_policy_segment.segment1.path
  subnets      = ["1.1.12.2/24"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [nsxt_policy_tier1_gateway.test]
}
`, interfaceName, interfaceDescription)
}

func CreateRealizationT1() string {
	return `
data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier1_gateway_interface.test.path
  depends_on = [nsxt_policy_tier1_gateway_interface.test]
}

data "nsxt_policy_gateway_interface_realization" "gw_realization" {
  gateway_path = nsxt_policy_tier1_gateway_interface.test.path
  depends_on = [nsxt_policy_tier1_gateway_interface.test]
}
`
}
