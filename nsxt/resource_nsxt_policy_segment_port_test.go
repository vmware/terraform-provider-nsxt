package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNsxtPolicySegmentPort_basic(t *testing.T) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	updatedSegmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "nsxt_policy_segment_port.test"
	createResourceTag := "profile1"
	updateResourceTag := "profile2"
	tzName := getOverlayTransportZoneName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				// Create
				Config: testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", segmentPortName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", "/infra/ip-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", "/infra/mac-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttrSet(testResourceName, "security_profile.0.spoofguard_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", "/infra/segment-security-profiles/"+profilesPrefix+"create"),
				),
			},
			{
				// Update
				Config: testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, updatedSegmentPortName, updateResourceTag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedSegmentPortName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", "/infra/ip-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", "/infra/mac-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttrSet(testResourceName, "security_profile.0.spoofguard_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", "/infra/segment-security-profiles/"+profilesPrefix+"update"),
				),
			},
		},
	})
}

func testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName string, profilesPrefix string, segmentPortName string, resourceTag string) string {
	return testAccNsxtPolicySegmentBasicTemplate(tzName, segmentName) + segmentProfilesDataSourceTemplates(profilesPrefix+"create", "profile1") + segmentProfilesDataSourceTemplates(profilesPrefix+"update", "profile2") + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "test" {
  display_name = "%s"
  description = "Acceptance tests"
  segment_path = nsxt_policy_segment.test.path
  discovery_profile {
    ip_discovery_profile_path = nsxt_policy_ip_discovery_profile.%s.path
    mac_discovery_profile_path = nsxt_policy_mac_discovery_profile.%s.path
  }
  security_profile {
    spoofguard_profile_path = data.nsxt_policy_spoofguard_profile.%s.path
    security_profile_path = nsxt_policy_segment_security_profile.%s.path
  }
}
`, segmentPortName, resourceTag, resourceTag, resourceTag, resourceTag)
}

func segmentProfilesDataSourceTemplates(name string, resourceTag string) string {
	return fmt.Sprintf(`

resource "nsxt_policy_ip_discovery_profile" "%s" {
  display_name = "%s"
  nsx_id = "%s"
}

resource "nsxt_policy_mac_discovery_profile" "%s" {
  display_name = "%s"
  nsx_id = "%s"
}

resource "nsxt_policy_segment_security_profile" "%s" {
  display_name = "%s"
  nsx_id = "%s"
}

data "nsxt_policy_spoofguard_profile" "%s" {
    display_name = "default-spoofguard-profile"
}
`, resourceTag, name, name, resourceTag, name, name, resourceTag, name, name, resourceTag)
}
