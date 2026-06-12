// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceNsxtPolicyGroupBareMetalServerMembers_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerMembersReadTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerMembers_withEnforcementPoint(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerMembersWithEnforcementPointTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "enforcement_point_path", "/infra/sites/default/enforcement-points/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerMembers_emptyGroup(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_members.test"
	testName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupBareMetalServerMembersEmptyGroupTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttr(testResourceName, "items.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_interface_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerInterfaceMembersReadTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers_withEnforcementPoint(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_interface_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerInterfaceMembersWithEnforcementPointTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "enforcement_point_path", "/infra/sites/default/enforcement-points/default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers_emptyGroup(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_interface_members.test"
	testName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupBareMetalServerInterfaceMembersEmptyGroupTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttr(testResourceName, "items.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerMembers_multipleMembers(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerMembersMultipleTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
					// Verify first server member attributes
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.external_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.source_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.last_sync_time"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers_multipleMembers(t *testing.T) {
	testResourceName := "data.nsxt_policy_group_baremetal_server_interface_members.test"
	testName := getAccTestResourceName()

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
				Config: testAccNsxtPolicyGroupBareMetalServerInterfaceMembersMultipleTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttrSet(testResourceName, "group_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.#"),
					// Verify first interface member attributes
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.external_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.bms_external_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.source_id"),
					// state attribute should exist (even if empty)
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[testResourceName]
						if rs == nil {
							return fmt.Errorf("resource not found: %s", testResourceName)
						}
						// Check that the state attribute key exists in the terraform state
						if _, ok := rs.Primary.Attributes["items.0.state"]; !ok {
							return fmt.Errorf("attribute 'items.0.state' not found in resource %s", testResourceName)
						}
						return nil
					},
					resource.TestCheckResourceAttrSet(testResourceName, "items.0.last_sync_time"),
				),
			},
		},
	})
}

// Template functions

func testAccNsxtPolicyGroupBareMetalServerMembersReadTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with a server member
resource "nsxt_policy_group" "test" {
  display_name = "%s-bms-group"
  description  = "BMS group for group members test"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Test the group members data source
data "nsxt_policy_group_baremetal_server_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName, getTestBMSServerID())
}

func testAccNsxtPolicyGroupBareMetalServerMembersWithEnforcementPointTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with a server member
resource "nsxt_policy_group" "test" {
  display_name = "%s-bms-group"
  description  = "BMS group for group members test with enforcement point"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Test the group members data source with enforcement point
data "nsxt_policy_group_baremetal_server_members" "test" {
  domain                 = "default"
  group_id              = nsxt_policy_group.test.id
  enforcement_point_path = "/infra/sites/default/enforcement-points/default"
}
`, testName, getTestBMSServerID())
}

func testAccNsxtPolicyGroupBareMetalServerMembersEmptyGroupTemplate(testName string) string {
	return fmt.Sprintf(`
# Create an empty BMS group (no members)
resource "nsxt_policy_group" "test" {
  display_name = "%s-empty-bms-group"
  description  = "Empty BMS group for testing empty results"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["99999999-9999-9999-9999-999999999999"]
    }
  }
}

# Test the group members data source on empty group
data "nsxt_policy_group_baremetal_server_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName)
}

func testAccNsxtPolicyGroupBareMetalServerMembersMultipleTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with multiple server members
resource "nsxt_policy_group" "test" {
  display_name = "%s-multi-bms-group"
  description  = "BMS group with multiple members for testing"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Test the group members data source with multiple members
data "nsxt_policy_group_baremetal_server_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName, getTestBMSServerID())
}

func testAccNsxtPolicyGroupBareMetalServerInterfaceMembersReadTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with an interface member
resource "nsxt_policy_group" "test" {
  display_name = "%s-bmsi-group"
  description  = "BMSI group for group members test"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Test the group interface members data source
data "nsxt_policy_group_baremetal_server_interface_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName, getTestBMSInterfaceID())
}

func testAccNsxtPolicyGroupBareMetalServerInterfaceMembersWithEnforcementPointTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with an interface member
resource "nsxt_policy_group" "test" {
  display_name = "%s-bmsi-group"
  description  = "BMSI group for group members test with enforcement point"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Test the group interface members data source with enforcement point
data "nsxt_policy_group_baremetal_server_interface_members" "test" {
  domain                 = "default"
  group_id              = nsxt_policy_group.test.id
  enforcement_point_path = "/infra/sites/default/enforcement-points/default"
}
`, testName, getTestBMSInterfaceID())
}

func testAccNsxtPolicyGroupBareMetalServerInterfaceMembersEmptyGroupTemplate(testName string) string {
	return fmt.Sprintf(`
# Create an empty BMSI group (no members)
resource "nsxt_policy_group" "test" {
  display_name = "%s-empty-bmsi-group"
  description  = "Empty BMSI group for testing empty results"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["88888888-8888-8888-8888-888888888888"]
    }
  }
}

# Test the group interface members data source on empty group
data "nsxt_policy_group_baremetal_server_interface_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName)
}

func testAccNsxtPolicyGroupBareMetalServerInterfaceMembersMultipleTemplate(testName string) string {
	return fmt.Sprintf(`
# Create a BMS group with multiple interface members
resource "nsxt_policy_group" "test" {
  display_name = "%s-multi-bmsi-group"
  description  = "BMSI group with multiple members for testing"
  group_type   = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Test the group interface members data source with multiple members
data "nsxt_policy_group_baremetal_server_interface_members" "test" {
  domain   = "default"
  group_id = nsxt_policy_group.test.id
}
`, testName, getTestBMSInterfaceID())
}
