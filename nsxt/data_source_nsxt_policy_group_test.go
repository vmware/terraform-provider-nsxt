/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_domains "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicyGroup_basic(t *testing.T) {
	name := "terraform_ds_test"
	domain := "default"
	testResourceName := "data.nsxt_policy_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyGroupDeleteByName(domain, name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGroupCreate(domain, name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyGroupReadTemplate(domain, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "domain", domain),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroup_withSite(t *testing.T) {
	name := "terraform_gm_ds_test"
	domain := getTestSiteName()
	testResourceName := "data.nsxt_policy_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyGlobalManager(t)
			testAccEnvDefined(t, "NSXT_TEST_SITE_NAME")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyGroupDeleteByName(domain, name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGroupCreate(domain, name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyGroupReadTemplate(domain, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "domain", domain),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyGroupCreate(domain string, name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.Group{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.GroupBindingType(), gm_model.GroupBindingType())
		if convErr != nil {
			return convErr
		}

		client := gm_domains.NewDefaultGroupsClient(connector)
		err = client.Patch(domain, id, gmObj.(gm_model.Group))
	} else {
		client := domains.NewDefaultGroupsClient(connector)
		err = client.Patch(domain, id, obj)
	}

	if err != nil {
		return handleCreateError("Group", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyGroupDeleteByName(domain string, name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name and delete it
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(name, "Group")
		if err == nil {
			client := gm_domains.NewDefaultGroupsClient(connector)
			err := client.Delete(domain, objID, nil, nil)
			if err != nil {
				return handleDeleteError("Group", objID, err)
			}
			return nil
		}
	} else {
		client := domains.NewDefaultGroupsClient(connector)
		objList, err := client.List(domain, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("Group", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == name {
				err := client.Delete(domain, *objInList.Id, nil, nil)
				if err != nil {
					return handleDeleteError("Group", *objInList.Id, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("Error while deleting Group '%s': resource not found", name)
}

func testAccNsxtPolicyGroupReadTemplate(domain string, name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_group" "test" {
  display_name = "%s"
  domain       = "%s"
}`, name, domain)
}
