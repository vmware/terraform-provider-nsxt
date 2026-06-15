// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyServices_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyServicesBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func TestAccDataSourceNsxtPolicyServices_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyServicesBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func TestAccDataSourceNsxtPolicyServices_builtInOnly(t *testing.T) {
	serviceName := getAccTestDataSourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyServicesBuiltInOnlyTemplate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_services.test_built_in", "id"),
					resource.TestCheckNoResourceAttr("data.nsxt_policy_services.test_built_in", fmt.Sprintf("items.%s", serviceName)),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_services.test_all", fmt.Sprintf("items.%s", serviceName)),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyServices_includeSharedServices(t *testing.T) {
	serviceName := getAccTestDataSourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyServicesIncludeSharedServicesTemplate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_services.test_shared", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_services.test_shared", fmt.Sprintf("items.%s", serviceName)),
					resource.TestCheckNoResourceAttr("data.nsxt_policy_services.test_no_shared", fmt.Sprintf("items.%s", serviceName)),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyServicesBasic(t *testing.T, withContext bool, preCheck func()) {
	serviceName := getAccTestDataSourceName()
	testResourceName := "data.nsxt_policy_services.test"
	checkResourceName := "data.nsxt_policy_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyServicesReadTemplate(serviceName, withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(checkResourceName, "display_name", serviceName),
				),
			},
		},
	})
}

func testAccNSXPolicyServicesReadTemplate(serviceName string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return testAccNsxtPolicyIcmpTypeServiceCreateTypeCodeTemplate(serviceName, "3", "1", "ICMPv4", withContext) + fmt.Sprintf(`
data "nsxt_policy_services" "test" {
  depends_on = [nsxt_policy_service.test]
%s
}

locals {
  // Get id from path
  path_split = split("/", data.nsxt_policy_services.test.items["%s"])
  service_id = element(local.path_split, length(local.path_split) - 1)
}

data "nsxt_policy_service" "test" {
%s
  id = local.service_id
}
`, context, serviceName, context)
}

func testAccNSXPolicyServicesBuiltInOnlyTemplate(serviceName string) string {
	projectContext := testAccNsxtPolicyMultitenancyContext()
	return testAccNsxtPolicyIcmpTypeServiceCreateTypeCodeTemplate(serviceName, "3", "1", "ICMPv4", true) + fmt.Sprintf(`
data "nsxt_policy_services" "test_built_in" {
  depends_on    = [nsxt_policy_service.test]
%s
  built_in_only = true
}

data "nsxt_policy_services" "test_all" {
  depends_on = [nsxt_policy_service.test]
%s
}
`, projectContext, projectContext)
}

func testAccNSXPolicyServicesIncludeSharedServicesTemplate(serviceName string) string {
	projectContext := testAccNsxtPolicyMultitenancyContext()
	projectID := os.Getenv("NSXT_VPC_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("NSXT_PROJECT_ID")
	}
	projectPath := fmt.Sprintf("/orgs/default/projects/%s", projectID)
	shareName := getAccTestResourceName()

	// Create the service in default space (withContext = false) so that it resides in default space and can be shared with the project
	serviceDef := testAccNsxtPolicyIcmpTypeServiceCreateTypeCodeTemplate(serviceName, "3", "1", "ICMPv4", false)

	shareDef := fmt.Sprintf(`
resource "nsxt_policy_share" "test" {
  display_name = "%s"

  sharing_strategy = "ALL_DESCENDANTS"
  shared_with      = ["%s"]
}

resource "nsxt_policy_shared_resource" "test" {
  display_name = "%s"

  share_path = nsxt_policy_share.test.path
  resource_object {
    resource_path    = nsxt_policy_service.test.path
    include_children = true
  }
}`, shareName, projectPath, shareName)

	return serviceDef + shareDef + fmt.Sprintf(`
data "nsxt_policy_services" "test_shared" {
  depends_on = [nsxt_policy_shared_resource.test]
%s
  include_shared_services = true
}

data "nsxt_policy_services" "test_no_shared" {
  depends_on = [nsxt_policy_shared_resource.test]
%s
  include_shared_services = false
}
`, projectContext, projectContext)
}
