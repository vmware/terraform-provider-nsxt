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

func TestAccDataSourceNsxtVpcLBMonitor_basic(t *testing.T) {
	name := getAccTestDataSourceName()
	testResourceName := "data.nsxt_vpc_lb_monitor.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtVpcLBMonitorDeleteByName(name)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtVpcLBMonitorCreate(name); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtVpcLBMonitorReadTemplate(name),
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
				Config: testAccNsxtVpcLBMonitorNameOnlyTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				/* fetch test profile by name prefix */
				Config: testAccNsxtVpcLBMonitorNameOnlyTemplate(name[:len(name)-2]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
				),
			},
			{
				/* fetch test profile by ID */
				Config: testAccNsxtVpcLBMonitorIDTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "type", "TCP"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		}),
	})
}

func testAccDataSourceNsxtVpcLBMonitorCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	sessionContext := testAccGetSessionContext()
	parents := getVpcParentsFromContext(sessionContext)
	client := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	displayName := name
	description := name
	obj := model.LBTcpMonitorProfile{
		Description:  &description,
		DisplayName:  &displayName,
		ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE,
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBTcpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error during conversion of LBTcpMonitorProfile: %v", errs[0])
	}

	// the name is used as ID, so that it is known when the config is rendered
	id := name
	err = client.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("MonitorProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtVpcLBMonitorDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	sessionContext := testAccGetSessionContext()
	parents := getVpcParentsFromContext(sessionContext)
	client := cliVpcLbMonitorProfilesClient(sessionContext, connector)

	objList, err := client.List(parents[0], parents[1], parents[2], nil, nil, nil, nil, nil, nil)
	if err != nil {
		return handleListError("MonitorProfile", err)
	}
	force := true
	for _, objInList := range objList.Results {
		result, err := policyLbMonitorConvert(objInList, "ANY")
		if err != nil {
			return fmt.Errorf("Error during MonitorProfile conversion: %v", err)
		}
		if result != nil && *result.DisplayName == name {
			err := client.Delete(parents[0], parents[1], parents[2], *result.Id, &force)
			if err != nil {
				return handleDeleteError("MonitorProfile", *result.Id, err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting MonitorProfile '%s': resource not found", name)
}

func testAccNsxtVpcLBMonitorReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_monitor" "test" {
%s
  type         = "TCP"
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}

func testAccNsxtVpcLBMonitorNameOnlyTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_monitor" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), name)
}

func testAccNsxtVpcLBMonitorIDTemplate(id string) string {
	return fmt.Sprintf(`
data "nsxt_vpc_lb_monitor" "test" {
%s
  id = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), id)
}
