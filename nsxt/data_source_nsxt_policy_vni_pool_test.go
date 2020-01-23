/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"testing"
)

func TestAccDataSourceNsxtPolicyVniPoolConfig_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_vni_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyVniPoolConfigDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyVniPoolConfigCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyVniPoolConfigReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "start"),
					resource.TestCheckResourceAttrSet(testResourceName, "end"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyVniPoolConfigCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultVniPoolsClient(connector)

	displayName := name
	description := name
	obj := model.VniPoolConfig{
		Description: &description,
		DisplayName: &displayName,
		Start:       75002,
		End:         95001,
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, obj)
	if err != nil {
		return fmt.Errorf("Error during VniPoolConfig creation: %v", err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyVniPoolConfigDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewDefaultVniPoolsClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("Error while reading VniPoolConfigs: %v", err)
	}
	for _, objInList := range objList.Results {
		if *objInList.DisplayName == name {
			err := client.Delete(*objInList.Id)
			if err != nil {
				return fmt.Errorf("Error during VniPoolConfig deletion: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting VniPoolConfig '%s': resource not found", name)
}

func testAccNsxtPolicyVniPoolConfigReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_vni_pool" "test" {
  display_name = "%s"
}`, name)
}
