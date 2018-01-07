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

func TestNSXLogicalTier0RouterBasic(t *testing.T) {
	routerName := Tier0RouterDefaultName
	testResourceName := "data.nsxt_logical_tier0_router.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalTier0RouterReadTemplate(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalTier0RouterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", routerName),
				),
			},
		},
	})
}

func testAccNSXLogicalTier0RouterExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical tier0 router resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical tier0 router resource ID not set in resources ")
		}

		_, responseCode, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical tier0 router ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical tier0 router %s exists. HTTP return code was %d", resourceID, responseCode)
		}

		return nil
	}
}

func testAccNSXLogicalTier0RouterReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_logical_tier0_router" "test" {
     display_name = "%s"
}`, name)
}
