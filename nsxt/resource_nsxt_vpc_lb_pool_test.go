// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestVpcLBPoolCreateAttributes = map[string]string{
	"description":             "terraform created",
	"algorithm":               "IP_HASH",
	"min_active_members":      "2",
	"tcp_multiplexing_number": "2",
}

var accTestVpcLBPoolUpdateAttributes = map[string]string{
	"description":             "terraform updated",
	"algorithm":               "WEIGHTED_ROUND_ROBIN",
	"min_active_members":      "5",
	"tcp_multiplexing_number": "5",
}

func TestAccResourceNsxtVpcLBPool_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_pool.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPoolCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBPoolTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPoolExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBPoolCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "algorithm", accTestVpcLBPoolCreateAttributes["algorithm"]),
					resource.TestCheckResourceAttr(testResourceName, "member.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.admin_state", "ENABLED"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.display_name", "member1"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.backup_member", "false"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.ip_address", "5.5.5.5"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.max_concurrent_connections", "7"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.port", "77"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.weight", "2"),
					resource.TestCheckResourceAttr(testResourceName, "member.1.admin_state", "DISABLED"),
					resource.TestCheckResourceAttr(testResourceName, "member.1.display_name", "member2"),
					resource.TestCheckResourceAttr(testResourceName, "member.1.ip_address", "5.5.5.3"),
					resource.TestCheckResourceAttr(testResourceName, "min_active_members", accTestVpcLBPoolCreateAttributes["min_active_members"]),
					resource.TestCheckResourceAttr(testResourceName, "tcp_multiplexing_number", accTestVpcLBPoolCreateAttributes["tcp_multiplexing_number"]),
					resource.TestCheckResourceAttr(testResourceName, "active_monitor_paths.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "snat.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "snat.0.type", "AUTOMAP"),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", os.Getenv("NSXT_VPC_PROJECT_ID")),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBPoolTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPoolExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBPoolUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "algorithm", accTestVpcLBPoolUpdateAttributes["algorithm"]),
					resource.TestCheckResourceAttr(testResourceName, "member.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.admin_state", "ENABLED"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.display_name", "member1"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.ip_address", "5.5.5.5"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.max_concurrent_connections", "100"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.port", "88"),
					resource.TestCheckResourceAttr(testResourceName, "member.0.weight", "1"),
					resource.TestCheckResourceAttr(testResourceName, "min_active_members", accTestVpcLBPoolUpdateAttributes["min_active_members"]),
					resource.TestCheckResourceAttr(testResourceName, "tcp_multiplexing_number", accTestVpcLBPoolUpdateAttributes["tcp_multiplexing_number"]),
					resource.TestCheckResourceAttr(testResourceName, "active_monitor_paths.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "snat.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "snat.0.type", "DISABLED"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				// the monitor stays in the config so that it is not deleted before the pool stops referencing it
				Config: testAccNsxtVpcLBPoolMonitorDep(updateName) + testAccNsxtVpcLBPoolMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPoolExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "member.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcLBPool_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.2")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPoolCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBPoolMinimalistic(name),
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

func testAccNsxtVpcLBPoolExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPC LBPool resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VPC LBPool resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBPoolExistsCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VPC LBPool %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtVpcLBPoolCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_lb_pool" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBPoolExistsCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("VPC LBPool %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBPoolMonitorDep(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_icmp_monitor_profile" "dep" {
  %s
  display_name = "%s-monitor"
}
`, testAccNsxtPolicyMultitenancyContext(), name)
}

func testAccNsxtVpcLBPoolTemplate(createFlow bool, name string) string {
	var attrMap map[string]string
	var members, snatType string
	if createFlow {
		attrMap = accTestVpcLBPoolCreateAttributes
		snatType = "AUTOMAP"
		members = `
  member {
    admin_state                = "ENABLED"
    backup_member              = false
    display_name               = "member1"
    ip_address                 = "5.5.5.5"
    max_concurrent_connections = 7
    port                       = "77"
    weight                     = 2
  }

  member {
    admin_state  = "DISABLED"
    display_name = "member2"
    ip_address   = "5.5.5.3"
  }
`
	} else {
		attrMap = accTestVpcLBPoolUpdateAttributes
		snatType = "DISABLED"
		members = `
  member {
    admin_state                = "ENABLED"
    backup_member              = false
    display_name               = "member1"
    ip_address                 = "5.5.5.5"
    max_concurrent_connections = 100
    port                       = "88"
    weight                     = 1
  }
`
	}
	return testAccNsxtVpcLBPoolMonitorDep(name) + fmt.Sprintf(`
resource "nsxt_vpc_lb_pool" "test" {
  %s
  display_name            = "%s"
  description             = "%s"
  algorithm               = "%s"
  min_active_members      = %s
  tcp_multiplexing_number = %s
  active_monitor_paths    = [nsxt_vpc_lb_icmp_monitor_profile.dep.path]
  %s

  snat {
    type = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["algorithm"], attrMap["min_active_members"], attrMap["tcp_multiplexing_number"], members, snatType)
}

func testAccNsxtVpcLBPoolMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_pool" "test" {
  %s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}
