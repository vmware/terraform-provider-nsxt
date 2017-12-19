/* Copyright © 2017 VMware, Inc. All Rights Reserved.
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

func TestNSXLogicalRouterBasic(t *testing.T) {

	name := fmt.Sprintf("test-nsx-logical-router")
	updateName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_logical_router.test"
        failoverMode := "PREEMPTIVE"
        haMode := "ACTIVE_STANDBY"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalRouterCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalRouterCreateTemplate(name, failoverMode, haMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "high_availability_mode", haMode),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "2"),
				),
			},
			{
				Config: testAccNSXLogicalRouterUpdateTemplate(updateName, failoverMode, haMode),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalRouterExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "high_availability_mode", haMode),
					resource.TestCheckResourceAttr(testResourceName, "failover_mode", failoverMode),
					resource.TestCheckResourceAttr(testResourceName, "tags.#", "1"),
				),
			},
		},
	})
}

func testAccNSXLogicalRouterExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical router resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical router resource ID not set in resources")
		}

		resource, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical router ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking verifying logical router existance. HTTP returned %d", resourceID, responseCode)
		}

		if displayName == resource.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX logical router %s not found", displayName)
	}
}

func testAccNSXLogicalRouterCheckDestroy(state *terraform.State, displayName string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_logical_router" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		resource, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving logical router ID %s. Error: %v", resourceID, err)
		}

		if displayName == resource.DisplayName {
			return fmt.Errorf("NSX logical router %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXLogicalRouterCreateTemplate(name string, failoverMode string, haMode string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_router" "test" {
display_name = "%s"
description = "Acceptance Test"
failover_mode = "%s"
high_availability_mode = "%s"
tags = [
    {
	scope = "scope1"
        tag = "tag1"
    }, {
	scope = "scope2"
    	tag = "tag2"
    }
]
}`, name, failoverMode, haMode)
}

func testAccNSXLogicalRouterUpdateTemplate(name string, failoverMode string, haMode string) string {
	return fmt.Sprintf(`
resource "nsxt_logical_router" "test" {
display_name = "%s"
description = "Acceptance Test Update"
failover_mode = "%s"
high_availability_mode = "%s"
tags = [
	{
	scope = "scope3"
    	tag = "tag3"
    },
]
}`, name, failoverMode, haMode)
}
