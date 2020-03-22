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

func TestAccDataSourceNsxtPolicyTier0Gateway_basic(t *testing.T) {
	routerName := "terraform_test_tier0"
	testResourceName := "data.nsxt_policy_tier0_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier0DeleteByName(routerName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier0Create(routerName); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyTier0GatewayReadTemplate(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", routerName),
					resource.TestCheckResourceAttr(testResourceName, "description", routerName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoTier0Template(),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyTier0Create(routerName string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultTier0sClient(connector)

	displayName := routerName
	description := routerName
	obj := model.Tier0{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Tier0", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyTier0DeleteByName(routerName string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultTier0sClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("Tier0", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == routerName {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return handleDeleteError("Tier0", *objInList.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting Tier0 '%s': resource not found", routerName)
}

func testAccNsxtPolicyTier0GatewayReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtPolicyNoTier0Template() string {
	return fmt.Sprintf(` `)
}
