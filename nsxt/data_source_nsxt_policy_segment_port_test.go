package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceNsxtPolicySegmentPort_basic(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPort_basic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicySegmentPort_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicySegmentPort_basic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicySegmentPort_basic(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_segment_port.segmentport1"
	createResourceTag := "profile1"
	tzName := getOverlayTransportZoneName()
	// Attachment field values
	attachmentID := getAccTestResourceName()
	attachmentType := "PARENT"
	allocateAddresses := "DHCP"
	hyperbusMode := "DISABLE"
	childAttachmentID := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag, withContext, attachmentID, attachmentType, allocateAddresses, hyperbusMode, childAttachmentID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag string, withContext bool, attachmentID string, attachmentType string, allocateAddresses string, hyperbusMode string, childAttachmentID string) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag, withContext, attachmentID, attachmentType, allocateAddresses, hyperbusMode, childAttachmentID) + fmt.Sprintf(`

data "nsxt_policy_segment_port" "segmentport1" {
	%s
	display_name = nsxt_policy_segment_port.test.display_name
}
`, context)
}
