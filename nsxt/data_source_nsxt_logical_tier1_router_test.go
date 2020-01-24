/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"testing"
)

func TestAccDataSourceNsxtLogicalTier1Router_basic(t *testing.T) {
	routerName := "terraform_test_tier1"
	testResourceName := "data.nsxt_logical_tier1_router.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtTier1RouterDeleteByName(routerName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtTier1RouterCreate(routerName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXTier1RouterReadTemplate(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", routerName),
					resource.TestCheckResourceAttr(testResourceName, "description", routerName),
				),
			},
			{
				Config: testAccNSXNoTier1RouterTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtTier1RouterCreate(routerName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := routerName
	description := routerName
	routerType := "TIER1"
	logicalRouter := manager.LogicalRouter{
		Description: description,
		DisplayName: displayName,
		RouterType:  routerType,
	}

	logicalRouter, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logicalRouter)
	if err != nil {
		return fmt.Errorf("Error during router creation: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during router creation: %v", resp.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtTier1RouterDeleteByName(routerName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name
	objList, _, err := nsxClient.LogicalRoutingAndServicesApi.ListLogicalRouters(nsxClient.Context, nil)
	if err != nil {
		return fmt.Errorf("Error while reading routers: %v", err)
	}
	// go over the list to find the correct one
	for _, objInList := range objList.Results {
		if objInList.DisplayName == routerName {
			localVarOptionals := make(map[string]interface{})
			responseCode, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, objInList.Id, localVarOptionals)
			if err != nil {
				return fmt.Errorf("Error during router deletion: %v", err)
			}

			if responseCode.StatusCode != http.StatusOK {
				return fmt.Errorf("Unexpected status returned during router deletion: %v", responseCode.StatusCode)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting router '%s': service not found", routerName)
}

func testAccNSXTier1RouterReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_logical_tier1_router" "test" {
  display_name = "%s"
}`, name)
}

func testAccNSXNoTier1RouterTemplate() string {
	return fmt.Sprintf(` `)
}
