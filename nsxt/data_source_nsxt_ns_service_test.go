/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"testing"
)

func TestAccDataSourceNsxtNsService_basic(t *testing.T) {
	serviceName := "terraform_test_ns_service"
	testResourceName := "data.nsxt_ns_service.test"
	var s *terraform.State

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtNsServiceCreate(serviceName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXNsServiceReadTemplate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", serviceName),
					resource.TestCheckResourceAttr(testResourceName, "description", serviceName),
					copyStatePtr(&s),
				),
			},
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtNsServiceDelete(s, testResourceName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXNoNsServiceTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtNsServiceCreate(serviceName string) error {
	nsxClient := testAccGetClient()
	nsService := manager.IgmpTypeNsService{
		NsService: manager.NsService{
			DisplayName:    serviceName,
			Description:    serviceName,
			DefaultService: false,
		},
		NsserviceElement: manager.IgmpTypeNsServiceEntry{
			ResourceType: "IGMPTypeNSService",
		},
	}
	nsService, responseCode, err := nsxClient.GroupingObjectsApi.CreateIgmpTypeNSService(nsxClient.Context, nsService)
	if err != nil {
		return fmt.Errorf("Error during nsService creation: %v", err)
	}

	if responseCode.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during nsService creation: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtNsServiceDelete(state *terraform.State, resourceName string) error {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return fmt.Errorf("NSX nsService data source %s not found", resourceName)
	}

	resourceID := rs.Primary.ID
	if resourceID == "" {
		return fmt.Errorf("NSX nsService data source ID not set")
	}
	nsxClient := testAccGetClient()
	localVarOptionals := make(map[string]interface{})
	responseCode, err := nsxClient.GroupingObjectsApi.DeleteNSService(nsxClient.Context, resourceID, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during nsService deletion: %v", err)
	}

	if responseCode.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during nsService deletion: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccNSXNsServiceReadTemplate(serviceName string) string {
	return fmt.Sprintf(`
data "nsxt_ns_service" "test" {
  display_name = "%s"
}`, serviceName)
}

func testAccNSXNoNsServiceTemplate() string {
	return fmt.Sprintf(` `)
}
