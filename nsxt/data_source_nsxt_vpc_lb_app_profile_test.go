// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestAccDataSourceNsxtVpcLBAppProfile_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_vpc_lb_app_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtVpcLBAppProfileDeleteByName(name)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtVpcLBAppProfileCreate(name); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtVpcLBAppProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", testAccGetSessionContext().ProjectID),
					resource.TestCheckResourceAttr(testResourceName, "context.0.vpc_id", testAccGetSessionContext().VPCID),
				),
			},
			{
				/* fetch test profile by name only */
				Config: testAccNsxtVpcLBAppProfileNameOnlyTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				/* fetch test profile by name prefix */
				Config: testAccNsxtVpcLBAppProfileNameOnlyTemplate(name[:len(name)-2]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
				),
			},
			{
				/* fetch test profile by ID */
				Config: testAccNsxtVpcLBAppProfileIDTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		}),
	})
}

func testAccDataSourceNsxtVpcLBAppProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	sessionContext := testAccGetSessionContext()
	parents := getVpcParentsFromContext(sessionContext)
	client := cliVpcLbAppProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	displayName := name
	description := name
	obj := model.LBFastTcpProfile{
		Description:  &description,
		DisplayName:  &displayName,
		ResourceType: model.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE,
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBFastTcpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error during conversion of LBFastTcpProfile: %v", errs[0])
	}

	// the name is used as ID, so that it is known when the config is rendered
	id := name
	err = client.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("AppProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtVpcLBAppProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	sessionContext := testAccGetSessionContext()
	parents := getVpcParentsFromContext(sessionContext)
	client := cliVpcLbAppProfilesClient(sessionContext, connector)

	objList, err := client.List(parents[0], parents[1], parents[2], nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("AppProfile", err)
	}
	force := true
	for _, objInList := range objList.Results {
		result, err := policyLbAppProfileConvert(objInList, "ANY")
		if err != nil {
			return fmt.Errorf("Error during AppProfile conversion: %v", err)
		}
		if result != nil && *result.DisplayName == name {
			err := client.Delete(parents[0], parents[1], parents[2], *result.Id, &force)
			if err != nil {
				return handleDeleteError("AppProfile", *result.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting AppProfile '%s': resource not found", name)
}

func testAccNsxtVpcLBAppProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_app_profile" "test" {
%s
  type         = "TCP"
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}

func testAccNsxtVpcLBAppProfileNameOnlyTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_app_profile" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}

func testAccNsxtVpcLBAppProfileIDTemplate(id string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_app_profile" "test" {
%s
  id = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), id)
}
