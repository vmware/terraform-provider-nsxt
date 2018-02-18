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

func TestAccDataSourceNsxtEdgeCluster_basic(t *testing.T) {
	edgeClusterName := edgeClusterDefaultName
	testResourceName := "data.nsxt_edge_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXEdgeClusterReadTemplate(edgeClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXEdgeClusterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", edgeClusterName),
				),
			},
		},
	})
}

func testAccNSXEdgeClusterExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX edge cluster resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX edge cluster resource ID not set in resources ")
		}

		_, responseCode, err := nsxClient.NetworkTransportApi.ReadEdgeCluster(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving edge cluster ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if edge cluster %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		return nil
	}
}

func testAccNSXEdgeClusterReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_edge_cluster" "test" {
     display_name = "%s"
}`, name)
}
