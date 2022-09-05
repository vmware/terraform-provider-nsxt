/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtPolicyLBHTTPAppProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_lb_http_app_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyLBHTTPAppProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyLBHTTPAppProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyLBHTTPAppProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "HTTP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				/* fetch test profile by name only */
				Config: testAccNsxtPolicyLBHTTPAppProfileNameOnlyTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "HTTP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				/* fetch default HTTP profile */
				Config: testAccNsxtPolicyLBHTTPAppProfileTypeOnlyTemplate("HTTP"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "type", "HTTP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyLBHTTPAppProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewLbAppProfilesClient(connector)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	displayName := name
	description := name
	profileType := "LBHttpProfile"
	obj := model.LBHttpProfile{
		Description:  &description,
		DisplayName:  &displayName,
		ResourceType: profileType,
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBHttpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error during conversion of LBHttpProfile: %v", errs[0])
	}

	// Generate a random ID for the resource
	id := newUUID()

	err = client.Patch(id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("LBHttpProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyLBHTTPAppProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	client := infra.NewLbAppProfilesClient(connector)

	// Find the object by name
	objList, err := client.List(nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("LBHttpProfile", err)
	}
	force := true
	for _, objInList := range objList.Results {
		result, err := policyLbAppProfileConvert(objInList, "ANY")
		if err != nil {
			return fmt.Errorf("Error during LBHttpProfile conversion: %v", err)
		}
		if result != nil && *result.DisplayName == name {
			err := client.Delete(*result.Id, &force)
			if err != nil {
				return handleDeleteError("LBHttpProfile", *result.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting LBHttpProfile '%s': resource not found", name)
}

func testAccNsxtPolicyLBHTTPAppProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_lb_http_app_profile" "test" {
  type = "TCP"
  display_name = "%s"
}`, name)
}

func testAccNsxtPolicyLBHTTPAppProfileTypeOnlyTemplate(pType string) string {
	return fmt.Sprintf(`
data "nsxt_policy_lb_http_app_profile" "test" {
  type = "%s"
}`, pType)
}

func testAccNsxtPolicyLBHTTPAppProfileNameOnlyTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_lb_http_app_profile" "test" {
  display_name = "%s"
}`, name)
}
