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

func TestAccDataSourceNsxtPolicyQosProfile_basic(t *testing.T) {
	name := "terraform_ds_test"
	testResourceName := "data.nsxt_policy_qos_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyQosProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyQosProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyQosProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyQosProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.QosProfile{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.QosProfileBindingType(), gm_model.QosProfileBindingType())
		if convErr != nil {
			return convErr
		}

		client := gm_infra.NewDefaultQosProfilesClient(connector)
		err = client.Patch(id, gmObj.(gm_model.QosProfile))
	} else {
		client := infra.NewDefaultQosProfilesClient(connector)
		err = client.Patch(id, obj)
	}

	if err != nil {
		return handleCreateError("QosProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyQosProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name and delete it
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(name, "QosProfile")
		if err == nil {
			client := gm_infra.NewDefaultQosProfilesClient(connector)
			err := client.Delete(objID)
			if err != nil {
				return handleDeleteError("QosProfile", objID, err)
			}
			return nil
		}
	} else {
		client := infra.NewDefaultQosProfilesClient(connector)
		objList, err := client.List(nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("QosProfile", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == name {
				err := client.Delete(*objInList.Id)
				if err != nil {
					return handleDeleteError("QosProfile", *objInList.Id, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("Error while deleting QosProfile '%s': resource not found", name)
}

func testAccNsxtPolicyQosProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_qos_profile" "test" {
  display_name = "%s"
}`, name)
}
