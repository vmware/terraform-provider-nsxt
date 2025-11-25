package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceNsxtPolicySegmentPort_basic(t *testing.T) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_segment_port.segmentport1"
	createResourceTag := "profile1"
	tzName := getOverlayTransportZoneName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag string) string {
	return testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag) + `

data "nsxt_policy_segment_port" "segmentport1" {
	display_name = nsxt_policy_segment_port.test.display_name
}
`
}
