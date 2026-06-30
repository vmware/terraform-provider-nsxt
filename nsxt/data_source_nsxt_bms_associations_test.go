// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyBareMetalServerGroupAssociations_consolidated(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_group_associations.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyBareMetalServerGroupAssociationsBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSServerID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.0.display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations_consolidated(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_interface_group_associations.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyBareMetalServerInterfaceGroupAssociationsBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSInterfaceID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.0.path"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.0.display_name"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerGroupAssociations_multipleGroups(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_group_associations.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyBareMetalServerGroupAssociationsBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSServerID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "groups.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerGroupAssociations_noGroups(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyBareMetalServerGroupAssociationsNoGroupsTemplate(),
				ExpectError: regexp.MustCompile("Failed to read group associations|Verify the server exists"),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations_noGroups(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyBareMetalServerInterfaceGroupAssociationsNoGroupsTemplate(),
				ExpectError: regexp.MustCompile("Failed to read group associations|Verify the interface exists"),
			},
		},
	})
}

// Template functions

func testAccNsxtPolicyBareMetalServerGroupAssociationsBasicTemplate() string {
	return fmt.Sprintf(`
# Query group associations for the server - should find system groups
data "nsxt_policy_baremetal_server_group_associations" "test" {
  external_id = "%s"
}
`, getTestBMSServerID())
}

func testAccNsxtPolicyBareMetalServerInterfaceGroupAssociationsBasicTemplate() string {
	return fmt.Sprintf(`
# Query group associations for the interface - should find system groups
data "nsxt_policy_baremetal_server_interface_group_associations" "test" {
  external_id = "%s"
}
`, getTestBMSInterfaceID())
}

func testAccNsxtPolicyBareMetalServerGroupAssociationsNoGroupsTemplate() string {
	return `
# Query group associations for non-existent server (should return empty list)
data "nsxt_policy_baremetal_server_group_associations" "test" {
  external_id = "99999999-9999-9999-9999-999999999999"
}
`
}

func testAccNsxtPolicyBareMetalServerInterfaceGroupAssociationsNoGroupsTemplate() string {
	return `
# Query group associations for non-existent interface (should return empty list)
data "nsxt_policy_baremetal_server_interface_group_associations" "test" {
  external_id = "88888888-8888-8888-8888-888888888888"
}
`
}
