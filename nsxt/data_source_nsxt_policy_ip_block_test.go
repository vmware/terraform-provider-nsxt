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

func TestAccDataSourceNsxtPolicyIpBlock_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_ip_block.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyIPBlockDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyIPBlockCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyIPBlockReadTemplate(name),
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

func testAccDataSourceNsxtPolicyIPBlockCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpBlocksClient(connector)

	displayName := name
	description := name
	cidr := "4001::/64"
	obj := model.IpAddressBlock{
		Description: &description,
		DisplayName: &displayName,
		Cidr:        &cidr,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IpAddressBlock", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyIPBlockDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpBlocksClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("IpiAddressBlock", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == name {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return handleDeleteError("IpAddressBlock", *objInList.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting IpBlock '%s': resource not found", name)
}

func testAccNsxtPolicyIPBlockReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_ip_block" "test" {
  display_name = "%s"
}`, name)
}
