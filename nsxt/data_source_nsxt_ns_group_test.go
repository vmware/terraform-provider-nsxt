/* Copyright © 2017 VMware, Inc. All Rights Reserved.
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

func TestAccDataSourceNsxtNsGroup_basic(t *testing.T) {
	groupName := "terraform_test_ns_group"
	testResourceName := "data.nsxt_ns_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtNsGroupDeleteByName(groupName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtNsGroupCreate(groupName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXNsGroupReadTemplate(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", groupName),
					resource.TestCheckResourceAttr(testResourceName, "description", groupName),
				),
			},
			{
				Config: testAccNSXNoNsGroupTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtNsGroupCreate(groupName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	nsGroup := manager.NsGroup{
		Description: groupName,
		DisplayName: groupName,
	}

	nsGroup, responseCode, err := nsxClient.GroupingObjectsApi.CreateNSGroup(nsxClient.Context, nsGroup)
	if err != nil {
		return fmt.Errorf("Error during nsGroup creation: %v", err)
	}

	if responseCode.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during nsGroup creation: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtNsGroupDeleteByName(groupName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name
	objList, _, err := nsxClient.GroupingObjectsApi.ListNSGroups(nsxClient.Context, nil)
	if err != nil {
		return fmt.Errorf("Error while reading NS groups: %v", err)
	}
	// go over the list to find the correct one
	for _, objInList := range objList.Results {
		if objInList.DisplayName == groupName {
			localVarOptionals := make(map[string]interface{})
			responseCode, err := nsxClient.GroupingObjectsApi.DeleteNSGroup(nsxClient.Context, objInList.Id, localVarOptionals)
			if err != nil {
				return fmt.Errorf("Error during nsGroup deletion: %v", err)
			}

			if responseCode.StatusCode != http.StatusOK {
				return fmt.Errorf("Unexpected status returned during nsGroup deletion: %v", responseCode.StatusCode)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting NS group '%s': group not found", groupName)
}

func testAccNSXNsGroupReadTemplate(groupName string) string {
	return fmt.Sprintf(`
data "nsxt_ns_group" "test" {
  display_name = "%s"
}`, groupName)
}

func testAccNSXNoNsGroupTemplate() string {
	return fmt.Sprintf(` `)
}
