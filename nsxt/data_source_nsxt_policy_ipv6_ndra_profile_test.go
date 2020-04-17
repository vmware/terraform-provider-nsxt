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

func TestAccDataSourceNsxtPolicyIpv6NdraProfile_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_ipv6_ndra_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyIpv6NdraProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyIpv6NdraProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyIpv6NdraProfileReadTemplate(name),
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

func testAccDataSourceNsxtPolicyIpv6NdraProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpv6NdraProfilesClient(connector)

	displayName := name
	description := name
	mode := model.Ipv6NdraProfile_RA_MODE_DISABLED
	config := model.RAConfig{}
	obj := model.Ipv6NdraProfile{
		Description: &description,
		DisplayName: &displayName,
		RaMode:      &mode,
		RaConfig:    &config,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Ipv6NdraProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyIpv6NdraProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultIpv6NdraProfilesClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("Ipv6NdraProfile", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == name {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return fmt.Errorf("Error during Ipv6NdraProfile deletion: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting Ipv6NdraProfile '%s': resource not found", name)
}

func testAccNsxtPolicyIpv6NdraProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_ipv6_ndra_profile" "test" {
  display_name = "%s"
}`, name)
}
