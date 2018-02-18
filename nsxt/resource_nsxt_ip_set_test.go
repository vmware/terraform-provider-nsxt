/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtIpSet_basic(t *testing.T) {

	name := fmt.Sprintf("test-nsx-ip-set")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_ip_set.test"
	single_ip := "1.1.1.1"
	additional_cidr := "2.1.1.0/24"
	additional_range := "3.1.1.1-3.1.1.10"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXIpSetCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIpSetCreateTemplate(name, single_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIpSetExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "1"),
				),
			},
			{
				Config: testAccNSXIpSetUpdateTemplate(updateName, single_ip, additional_cidr, additional_range),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIpSetExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.#", "3"),
				),
			},
		},
	})
}

func testAccNSXIpSetExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("IP Set resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("IP Set resource ID not set in resources ")
		}

		profile, responseCode, err := nsxClient.GroupingObjectsApi.ReadIPSet(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving IP Set ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if IP Set %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == profile.DisplayName {
			return nil
		}
		return fmt.Errorf("IP Set %s wasn't found", display_name)
	}
}

func testAccNSXIpSetCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_port" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		profile, responseCode, err := nsxClient.GroupingObjectsApi.ReadIPSet(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving IP Set ID %s. Error: %v", resourceID, err)
		}

		if display_name == profile.DisplayName {
			return fmt.Errorf("IP Set %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXIpSetCreateTemplate(name string, single_ip string) string {
	return fmt.Sprintf(`
resource "nsxt_ip_set" "test" {
	display_name = "%s"
	description = "Acceptance Test"
	ip_addresses = ["%s"]
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
}`, name, single_ip)
}

func testAccNSXIpSetUpdateTemplate(updatedName string, single_ip string, additional_range string, additional_cidr string) string {
	return fmt.Sprintf(`
resource "nsxt_ip_set" "test" {
	display_name = "%s"
	description = "Acceptance Test Update"
	ip_addresses = ["%s", "%s", "%s"]
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
    tag {
    	scope = "scope2"
        tag = "tag2"
    }
}`, updatedName, single_ip, additional_range, additional_cidr)
}
