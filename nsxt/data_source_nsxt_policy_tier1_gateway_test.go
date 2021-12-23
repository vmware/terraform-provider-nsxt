/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicyTier1Gateway_basic(t *testing.T) {
	routerName := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_tier1_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
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
		},
	})
}

func testAccDataSourceNsxtPolicyTier1GatewayCreate(routerName string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := routerName
	description := routerName
	obj := model.Tier1{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()
	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.Tier1BindingType(), gm_model.Tier1BindingType())
		if convErr != nil {
			return convErr
		}

		client := gm_infra.NewTier1sClient(connector)
		err = client.Patch(id, gmObj.(gm_model.Tier1))

	} else {
		client := infra.NewTier1sClient(connector)
		err = client.Patch(id, obj)
	}
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

	// Find the object by name
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(routerName, "Tier1")
		if err == nil {
			client := gm_infra.NewTier1sClient(connector)
			err := client.Delete(objID)
			if err != nil {
				return handleDeleteError("Tier1", objID, err)
			}
			return nil
		}
	} else {
		client := infra.NewTier1sClient(connector)

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
	}
	return fmt.Errorf("Error while deleting Tier1 '%s': resource not found", routerName)
}

func testAccNsxtPolicyTier1ReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
}`, name)
}
