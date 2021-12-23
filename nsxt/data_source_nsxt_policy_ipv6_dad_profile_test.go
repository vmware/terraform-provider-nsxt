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

func TestAccDataSourceNsxtPolicyIpv6DadProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_ipv6_dad_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyIpv6DadProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyIpv6DadProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyIpv6DadProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyIpv6DadProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.Ipv6DadProfile{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.Ipv6DadProfileBindingType(), gm_model.Ipv6DadProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewIpv6DadProfilesClient(connector)
		err = client.Patch(id, gmObj.(gm_model.Ipv6DadProfile), nil)
	} else {
		client := infra.NewIpv6DadProfilesClient(connector)
		err = client.Patch(id, obj, nil)
	}

	if err != nil {
		return handleCreateError("Ipv6DadProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyIpv6DadProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	// Find the object by name and delete it
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(name, "Ipv6DadProfile")
		if err == nil {
			client := gm_infra.NewIpv6DadProfilesClient(connector)
			err := client.Delete(objID, nil)
			if err != nil {
				return handleDeleteError("Ipv6DadProfile", objID, err)
			}
			return nil
		}
	} else {
		client := infra.NewIpv6DadProfilesClient(connector)

		// Find the object by name
		objList, err := client.List(nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("Ipv6DadProfile", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == name {
				err := client.Delete(*objInList.Id, nil)
				if err != nil {
					return handleDeleteError("Ipv6DadProfile", *objInList.Id, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("Error while deleting Ipv6DadProfile '%s': resource not found", name)
}

func testAccNsxtPolicyIpv6DadProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_ipv6_dad_profile" "test" {
  display_name = "%s"
}`, name)
}
