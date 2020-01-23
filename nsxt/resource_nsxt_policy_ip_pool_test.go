/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

func TestAccResourceNsxtPolicyIPPool_minimal(t *testing.T) {
	name := "tfpool1"
	testResourceName := "nsxt_policy_ip_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolCreateMinimalTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPool_basic(t *testing.T) {
	name := "tfpool1"
	updatedName := fmt.Sprintf("%s-updated", name)
	testResourceName := "nsxt_policy_ip_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPool_import_basic(t *testing.T) {
	name := "tfpool1"
	testResourceName := "nsxt_policy_ip_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolCreateTemplate(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXPolicyIPPoolCheckExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		client := infra.NewDefaultIpPoolsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX Policy IP Pool resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX Policy IP Pool resource ID not set in resources")
		}

		_, err := client.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Failed to find IP Pool %s", resourceID)
		}

		return nil
	}
}

func testAccNSXPolicyIPPoolCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	client := infra.NewDefaultIpPoolsClient(connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_ip_pool" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := client.Get(resourceID)
		if err == nil {
			return fmt.Errorf("IP Pool still exists %s", resourceID)
		}

	}
	return nil
}

func testAccNSXPolicyIPPoolCreateMinimalTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "test" {
  display_name = "%s"
}`, name)
}

func testAccNSXPolicyIPPoolCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, name)
}

func testAccNSXPolicyIPPoolUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, name)
}
