/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"testing"
)

func TestAccDataSourceNsxtPolicyIpPool_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_ip_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyIPPoolDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyIPPoolCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyIPPoolReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyIPPoolCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpPoolsClient(connector)

	displayName := name
	description := name
	obj := model.IpAddressPool{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IpAddressPool", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyIPPoolDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpPoolsClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("IpPool", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == name {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return handleDeleteError("IpAddressPool", *objInList.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting IpPool '%s': resource not found", name)
}

func testAccNsxtPolicyIPPoolReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_ip_pool" "test" {
  display_name = "%s"
}`, name)
}
