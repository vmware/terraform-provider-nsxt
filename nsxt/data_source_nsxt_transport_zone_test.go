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

func TestNSXTransportZoneBasic(t *testing.T) {
	transportZoneName := vlanTransportZoneName
	testResourceName := "data.nsxt_transport_zone.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXTransportZoneReadTemplate(transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXTransportZoneExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", transportZoneName),
				),
			},
		},
	})
}

func testAccNSXTransportZoneExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX transport zone resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX transport zone resource ID not set in resources ")
		}

		_, responseCode, err := nsxClient.NetworkTransportApi.GetTransportZone(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving transport zone ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if transport zone %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		return nil
	}
}

func testAccNSXTransportZoneReadTemplate(transportZoneName string) string {
	return fmt.Sprintf(`
data "nsxt_transport_zone" "test" {
     display_name = "%s"
}`, transportZoneName)
}
