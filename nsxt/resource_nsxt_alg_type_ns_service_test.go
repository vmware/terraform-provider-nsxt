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

func TestAccResourceNsxtAlgTypeNsService_basic(t *testing.T) {

	serviceName := fmt.Sprintf("test-nsx-alg-service")
	updateServiceName := fmt.Sprintf("%s-update", serviceName)
	testResourceName := "nsxt_alg_type_ns_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXAlgServiceCheckDestroy(state, serviceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXAlgServiceCreateTemplate(serviceName, "FTP", "9000-9001", "21"),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXAlgServiceExists(serviceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "alg service"),
					resource.TestCheckResourceAttr(testResourceName, "alg", "FTP"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXAlgServiceCreateTemplate(updateServiceName, "ORACLE_TNS", "600", "8081"),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXAlgServiceExists(updateServiceName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateServiceName),
					resource.TestCheckResourceAttr(testResourceName, "description", "alg service"),
					resource.TestCheckResourceAttr(testResourceName, "alg", "ORACLE_TNS"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func testAccNSXAlgServiceExists(display_name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX alg ns service resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX alg ns service resource ID not set in resources ")
		}

		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadAlgTypeNSService(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving alg ns service ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if alg ns service %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if display_name == service.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX alg ns service %s wasn't found", display_name)
	}
}

func testAccNSXAlgServiceCheckDestroy(state *terraform.State, display_name string) error {

	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_alg_type_ns_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		service, responseCode, err := nsxClient.GroupingObjectsApi.ReadAlgTypeNSService(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving alg ns service ID %s. Error: %v", resourceID, err)
		}

		if display_name == service.DisplayName {
			return fmt.Errorf("NSX alg ns service %s still exists", display_name)
		}
	}
	return nil
}

func testAccNSXAlgServiceCreateTemplate(serviceName string, protocol string, source_ports string, dest_ports string) string {
	return fmt.Sprintf(`
resource "nsxt_alg_type_ns_service" "test" {
    description = "alg service"
    display_name = "%s"
    alg = "%s"
    source_ports = [ "%s" ]
    destination_ports = "%s"
    tag {
    	scope = "scope1"
        tag = "tag1"
    }
}`, serviceName, protocol, source_ports, dest_ports)
}
