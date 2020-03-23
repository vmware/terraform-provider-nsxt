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

func TestAccDataSourceNsxtPolicyLBClientSslProfile_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_lb_client_ssl_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyLBClientSslProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyLBClientSslProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyLBClientSslProfileReadTemplate(name),
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

func testAccDataSourceNsxtPolicyLBClientSslProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultLbClientSslProfilesClient(connector)

	displayName := name
	description := name
	obj := model.LBClientSslProfile{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("LBClientSslProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyLBClientSslProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultLbClientSslProfilesClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("LBClientSslProfile", err)
	}
	force := true
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == name {
			err := client.Delete(*objInList.Id, &force)
			if err != nil {
				return handleDeleteError("LBClientSslProfile", *objInList.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting LBClientSslProfile '%s': resource not found", name)
}

func testAccNsxtPolicyLBClientSslProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_lb_client_ssl_profile" "test" {
  display_name = "%s"
}`, name)
}
