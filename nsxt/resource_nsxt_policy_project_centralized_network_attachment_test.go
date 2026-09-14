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

func TestAccResourceNsxtPolicyProjectCentralizedNetworkAttachment_basic(t *testing.T) {
	testResourceName := "nsxt_policy_project_centralized_network_attachment.test"
	prereqName := getAccTestResourceName()
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCentralizedNetworkAttachmentCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectCNATemplate(prereqName, createName, "terraform created", false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectCentralizedNetworkAttachmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform created"),
					resource.TestCheckResourceAttrPair(testResourceName, "subnet_path", "nsxt_vpc_subnet.test", "path"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_private", "false"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.#", "1"),
					resource.TestCheckResourceAttrPair(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.0", "nsxt_policy_ip_block.test", "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectCNATemplate(prereqName, updateName, "terraform updated", true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectCentralizedNetworkAttachmentExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform updated"),
					resource.TestCheckResourceAttrPair(testResourceName, "subnet_path", "nsxt_vpc_subnet.test", "path"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_private", "true"),
					resource.TestCheckResourceAttr(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.#", "1"),
					resource.TestCheckResourceAttrPair(testResourceName, "advertise_outbound_networks.0.allow_external_blocks.0", "nsxt_policy_ip_block.test", "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProjectCentralizedNetworkAttachment_importBasic(t *testing.T) {
	prereqName := getAccTestResourceName()
	displayName := getAccTestResourceName()
	testResourceName := "nsxt_policy_project_centralized_network_attachment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCentralizedNetworkAttachmentCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectCNATemplate(prereqName, displayName, "terraform created", false),
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

func testAccNsxtPolicyProjectCentralizedNetworkAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy CentralizedNetworkAttachment resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy CentralizedNetworkAttachment resource ID not set in resources")
		}

		path := rs.Primary.Attributes["path"]
		sessionContext := testAccGetSessionContextFromParentPath(testAccProvider.Meta(), path)
		exists, err := resourceNsxtPolicyProjectCNAExists(sessionContext, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy CentralizedNetworkAttachment %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyProjectCentralizedNetworkAttachmentCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_project_centralized_network_attachment" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		path := rs.Primary.Attributes["path"]
		sessionContext := testAccGetSessionContextFromParentPath(testAccProvider.Meta(), path)
		exists, err := resourceNsxtPolicyProjectCNAExists(sessionContext, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy CentralizedNetworkAttachment %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyProjectCNATemplate(prereqName string, displayName string, description string, allowPrivate bool) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_ip_block" "test" {
  display_name = "%s"
  cidr         = "10.120.0.0/22"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_project" "test" {
  display_name         = "%s"
  external_ipv4_blocks = [nsxt_policy_ip_block.test.path]
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

resource "nsxt_vpc" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  private_ips  = ["192.168.240.0/24"]
}

resource "nsxt_vpc_subnet" "test" {
  context {
    project_id = nsxt_policy_project.test.id
    vpc_id     = nsxt_vpc.test.id
  }
  display_name = "%s"
  ip_addresses = ["192.168.240.0/26"]
  access_mode  = "Private"
}

resource "nsxt_policy_project_centralized_network_attachment" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  description  = "%s"
  subnet_path  = nsxt_vpc_subnet.test.path

  advertise_outbound_networks {
    allow_private         = %t
    allow_external_blocks = [nsxt_policy_ip_block.test.path]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getEdgeClusterName(), prereqName, prereqName, prereqName, prereqName, displayName, description, allowPrivate)
}
