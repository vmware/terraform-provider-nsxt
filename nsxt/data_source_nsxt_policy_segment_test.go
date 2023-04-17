/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicySegment_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_segment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicySegmentDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicySegmentCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicySegmentReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicySegmentCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.Segment{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	uuid, _ := uuid.NewRandom()
	id := uuid.String()

	if testAccIsGlobalManager() {
		gmObj, convErr := convertModelBindingType(obj, model.SegmentBindingType(), gm_model.SegmentBindingType())
		if convErr != nil {
			return convErr
		}

		client := gm_infra.NewSegmentsClient(connector)
		err = client.Patch(id, gmObj.(gm_model.Segment))

	} else {
		client := infra.NewSegmentsClient(connector)
		err = client.Patch(id, obj)
	}
	if err != nil {
		return fmt.Errorf("Error during Segment creation: %v", err)
	}
	return nil
}

func testAccDataSourceNsxtPolicySegmentDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name
	objID, err := testGetObjIDByName(name, "Segment")
	if err != nil {
		return nil
	}
	if testAccIsGlobalManager() {
		client := gm_infra.NewSegmentsClient(connector)
		err = client.Delete(objID)
	} else {
		client := infra.NewSegmentsClient(connector)
		err = client.Delete(objID)
	}
	if err != nil {
		return fmt.Errorf("Error during Segment deletion: %v", err)
	}
	return nil
}

func testAccNsxtPolicySegmentReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_segment" "test" {
  display_name = "%s"
}`, name)
}
