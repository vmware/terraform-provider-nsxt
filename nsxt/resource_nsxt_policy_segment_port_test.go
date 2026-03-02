package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicySegmentPort_basic(t *testing.T) {
	testAccResourceNsxtPolicySegmentPort_basic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicySegmentPort_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicySegmentPort_basic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicySegmentPort_basic(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	updatedSegmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "nsxt_policy_segment_port.test"
	createResourceTag := "profile1"
	updateResourceTag := "profile2"
	tzName := getOverlayTransportZoneName()
	mtPrefix := ""
	if withContext {
		projectID := os.Getenv("NSXT_PROJECT_ID")
		mtPrefix = "/orgs/default/projects/" + projectID
	}
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
				// Create
				Config: testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag, withContext, attachmentID, attachmentType, allocateAddresses, hyperbusMode, childAttachmentID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", segmentPortName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", mtPrefix+"/infra/ip-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", mtPrefix+"/infra/mac-discovery-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.spoofguard_profile_path", mtPrefix+"/infra/spoofguard-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", mtPrefix+"/infra/segment-security-profiles/"+profilesPrefix+"create"),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.id", attachmentID),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.type", attachmentType),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.allocate_addresses", allocateAddresses),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.hyperbus_mode", hyperbusMode),
					resource.TestCheckResourceAttr("nsxt_policy_segment_port.test_child", "attachment.0.context_id", attachmentID),
				),
			},
			{
				// Update
				Config: testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, updatedSegmentPortName, updateResourceTag, withContext, attachmentID, attachmentType, allocateAddresses, hyperbusMode, childAttachmentID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedSegmentPortName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.ip_discovery_profile_path", mtPrefix+"/infra/ip-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "discovery_profile.0.mac_discovery_profile_path", mtPrefix+"/infra/mac-discovery-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.spoofguard_profile_path", mtPrefix+"/infra/spoofguard-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "security_profile.0.security_profile_path", mtPrefix+"/infra/segment-security-profiles/"+profilesPrefix+"update"),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.id", attachmentID),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.type", attachmentType),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.allocate_addresses", allocateAddresses),
					resource.TestCheckResourceAttr(testResourceName, "attachment.0.hyperbus_mode", hyperbusMode),
					resource.TestCheckResourceAttr("nsxt_policy_segment_port.test_child", "attachment.0.context_id", attachmentID),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicySegmentPort_importBasic(t *testing.T) {
	testAccResourceNsxtPolicySegmentPort_importBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicySegmentPort_importBasic_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicySegmentPort_importBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicySegmentPort_importBasic(t *testing.T, withContext bool, preCheck func()) {
	segmentName := getAccTestResourceName()
	segmentPortName := getAccTestResourceName()
	profilesPrefix := getAccTestResourceName()
	testResourceName := "nsxt_policy_segment_port.test"
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
				Config: testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName, profilesPrefix, segmentPortName, createResourceTag, withContext, attachmentID, attachmentType, allocateAddresses, hyperbusMode, childAttachmentID),
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

func testAccResourceNsxtPolicySegmentPortTemplate(tzName, segmentName string, profilesPrefix string, segmentPortName string, resourceTag string, withContext bool, attachmentID string, attachmentType string, allocateAddresses string, hyperbusMode string, childAttachmentID string) string {
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
  discovery_profile {
    ip_discovery_profile_path = nsxt_policy_ip_discovery_profile.%s.path
    mac_discovery_profile_path = nsxt_policy_mac_discovery_profile.%s.path
  }
  security_profile {
    spoofguard_profile_path = nsxt_policy_spoofguard_profile.%s.path
    security_profile_path = nsxt_policy_segment_security_profile.%s.path
  }
  attachment {
    id                 = "%s"
    type               = "%s"
    allocate_addresses = "%s"
    hyperbus_mode      = "%s"
  }
}

resource "nsxt_policy_segment_port" "test_child" {
	%s
    display_name = "port_attach_test-child"
    description  = "Child Port with Attachment"
    segment_path = nsxt_policy_segment.test.path
    
    attachment {
		id                 = "%s"
		context_id         = "%s"
		type               = "CHILD"
		allocate_addresses = "DHCP"
		hyperbus_mode      = "DISABLE"
    }
    depends_on = [ nsxt_policy_segment_port.test ]
}
`, context, segmentPortName, resourceTag, resourceTag, resourceTag, resourceTag, attachmentID, attachmentType, allocateAddresses, hyperbusMode, context, childAttachmentID, attachmentID)
	if withContext {
		return testAccNsxtPolicySegmentNoTransportZoneTemplate(segmentName, "12.12.2.1/24", withContext) + tfConfigTemp
	}
	return testAccNsxtPolicySegmentImportTemplate(tzName, segmentName, withContext) + tfConfigTemp
}

func segmentProfilesDataSourceTemplates(name string, resourceTag string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	profileResource := func(resourceType string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  %s
  display_name = "%s"
  nsx_id = "%s"
}`, resourceType, resourceTag, context, name, name)
	}
	return profileResource("nsxt_policy_ip_discovery_profile") +
		profileResource("nsxt_policy_mac_discovery_profile") +
		profileResource("nsxt_policy_segment_security_profile") +
		profileResource("nsxt_policy_spoofguard_profile")
}
