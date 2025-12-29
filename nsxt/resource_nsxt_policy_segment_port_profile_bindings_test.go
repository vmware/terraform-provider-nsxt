package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicySegmentPortProfileBindings_basic(t *testing.T) {
	testAccResourceNsxtPolicySegmentPortProfileBindings_basic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicySegmentPortProfileBindings_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicySegmentPortProfileBindings_basic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicySegmentPortProfileBindings_basic(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "nsxt_policy_segment_port_profile_bindings.test"
	createResourceTag := "profile1"
	updateResourceTag := "profile2"
	tzName := getOverlayTransportZoneName()
	mtPrefix := ""
	if withContext {
		projectID := os.Getenv("NSXT_PROJECT_ID")
		mtPrefix = "/orgs/default/projects/" + projectID
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySegmentCheckDestroy(state, segmentName)
		},
		Steps: []resource.TestStep{
			{
				// Create
				Config: testAccResourceNsxtPolicySegmentPortProfileBindingsTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySegmentPortProfileBindingsExists(testResourceName, withContext),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_port_path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", mtPrefix+"/infra/ip-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", mtPrefix+"/infra/mac-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.spoofguard_profile_path", mtPrefix+"/infra/spoofguard-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", mtPrefix+"/infra/segment-security-profiles/"+profilesPrefix+"create"),
				),
			},
			{
				// Update
				Config: testAccResourceNsxtPolicySegmentPortProfileBindingsTemplate(tzName, segmentName, profilesPrefix, segmentPortName, updateResourceTag, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySegmentPortProfileBindingsExists(testResourceName, withContext),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_port_path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", mtPrefix+"/infra/ip-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", mtPrefix+"/infra/mac-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.spoofguard_profile_path", mtPrefix+"/infra/spoofguard-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", mtPrefix+"/infra/segment-security-profiles/"+profilesPrefix+"update"),
				),
			},
		},
	})
}

func testAccNsxtPolicySegmentPortProfileBindingsExists(resourceName string, withContext bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Segment Port Binding resource %s not found in resources", resourceName)
		}
		segmentPortPath := rs.Primary.Attributes["segment_port_path"]
		segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
		resourceID := getPolicyIDFromPath(segmentPortPath)
		if resourceID == "" {
			return fmt.Errorf("Policy Segment Port Binding resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta())

		if err != nil {
			return fmt.Errorf("Error while parsing policy Segment Port Path %s. Error: %v", segmentPortPath, err)
		}

		_, err = getSegmentPort(segmentPath, resourceID, testAccGetSessionContext(), connector)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy Segment Port Binding ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccResourceNsxtPolicySegmentPortProfileBindingsTemplate(tzName, segmentName string, profilesPrefix string, segmentPortName string, resourceTag string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}

	tfConfigTemp := segmentProfilesDataSourceTemplates(profilesPrefix+"create", "profile1", withContext) + segmentProfilesDataSourceTemplates(profilesPrefix+"update", "profile2", withContext) + fmt.Sprintf(`

resource "nsxt_policy_segment_port" "test" {
  %s
  display_name = "%s"
  description = "Acceptance tests"
  segment_path = nsxt_policy_segment.test.path
  
  lifecycle {
    ignore_changes = [discovery_profile, qos_profile, security_profile]
  }
}

resource "nsxt_policy_segment_port_profile_bindings" "test" {
  %s
  segment_port_path = nsxt_policy_segment_port.test.path
  discovery_profile {
    ip_discovery_profile_path = nsxt_policy_ip_discovery_profile.%s.path
    mac_discovery_profile_path = nsxt_policy_mac_discovery_profile.%s.path
  }
  security_profile {
    spoofguard_profile_path = nsxt_policy_spoofguard_profile.%s.path
    security_profile_path = nsxt_policy_segment_security_profile.%s.path
  }
}
`, context, segmentPortName, context, resourceTag, resourceTag, resourceTag, resourceTag)
	if withContext {
		return testAccNsxtPolicySegmentNoTransportZoneTemplate(segmentName, "12.12.2.1/24", withContext) + tfConfigTemp
	}
	return testAccNsxtPolicySegmentImportTemplate(tzName, segmentName, withContext) + tfConfigTemp
}
