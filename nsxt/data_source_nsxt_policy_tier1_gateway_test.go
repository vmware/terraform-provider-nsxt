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

func TestAccDataSourceNsxtPolicyTier1Gateway_basic(t *testing.T) {
	routerName := "terraform_test_tier1"
	testResourceName := "data.nsxt_policy_tier1_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(routerName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier1GatewayCreate(routerName); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyTier1ReadTemplate(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", routerName),
					resource.TestCheckResourceAttr(testResourceName, "description", routerName),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyNoTier1Template(),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyTier1GatewayCreate(routerName string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultTier1sClient(connector)

	displayName := routerName
	description := routerName
	obj := model.Tier1{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Tier1", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(routerName string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultTier1sClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("Tier1", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == routerName {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return handleDeleteError("Tier1", *objInList.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting Tier1 '%s': resource not found", routerName)
}

func testAccNsxtPolicyTier1ReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
}`, name)
}

func testAccNsxtPolicyNoTier1Template() string {
	return fmt.Sprintf(` `)
}
