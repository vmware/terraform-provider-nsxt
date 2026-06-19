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

func TestAccDataSourceNsxtPolicyBareMetalServerTags_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_tags.test"

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
				Config: testAccNsxtPolicyBareMetalServerTagsReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSServerID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "tag.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerTags_withTags(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_tags.test"
	testName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyBareMetalServerTagsWithTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSServerID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "environment",
						"tag":   "test",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "test-run",
						"tag":   testName,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "application",
						"tag":   "web-server",
					}),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceTags_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_interface_tags.test"

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
				Config: testAccNsxtPolicyBareMetalServerInterfaceTagsReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSInterfaceID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "tag.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceTags_withTags(t *testing.T) {
	testResourceName := "data.nsxt_policy_baremetal_server_interface_tags.test"
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
				Config: testAccNsxtPolicyBareMetalServerInterfaceTagsWithTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "external_id", getTestBMSInterfaceID()),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "network-zone",
						"tag":   "dmz",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "test-run",
						"tag":   testName,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testResourceName, "tag.*", map[string]string{
						"scope": "interface-type",
						"tag":   "data-plane",
					}),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerTags_nonExistent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyBareMetalServerTagsNonExistentTemplate(),
				ExpectError: regexp.MustCompile("Failed to find Bare Metal Server|could not find bare metal server"),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceTags_nonExistent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyBareMetalServerInterfaceTagsNonExistentTemplate(),
				ExpectError: regexp.MustCompile("Failed to find Bare Metal Server Interface|could not find bare metal server interface"),
			},
		},
	})
}

// Template functions

func testAccNsxtPolicyBareMetalServerTagsReadTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "%s"
}
`, getTestBMSServerID())
}

func testAccNsxtPolicyBareMetalServerTagsWithTagsTemplate(testName string) string {
	return fmt.Sprintf(`
# Create tags for BMS server
resource "nsxt_policy_baremetal_server_tags" "tags_datasource_test" {
  external_id = "%s"

  tag {
    scope = "environment"
    tag   = "test"
  }

  tag {
    scope = "test-run"
    tag   = "%s"
  }

  tag {
    scope = "application"
    tag   = "web-server"
  }
}

# Read the tags back
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "%s"
  depends_on  = [nsxt_policy_baremetal_server_tags.tags_datasource_test]
}
`, getTestBMSServerID(), testName, getTestBMSServerID())
}

func testAccNsxtPolicyBareMetalServerInterfaceTagsReadTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "%s"
}
`, getTestBMSInterfaceID())
}

func testAccNsxtPolicyBareMetalServerInterfaceTagsWithTagsTemplate(testName string) string {
	return fmt.Sprintf(`
# Create tags for BMS interface
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags_datasource_test" {
  external_id = "%s"

  tag {
    scope = "network-zone"
    tag   = "dmz"
  }

  tag {
    scope = "test-run"
    tag   = "%s"
  }

  tag {
    scope = "interface-type"
    tag   = "data-plane"
  }
}

# Read the tags back
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "%s"
  depends_on  = [nsxt_policy_baremetal_server_interface_tags.interface_tags_datasource_test]
}
`, getTestBMSInterfaceID(), testName, getTestBMSInterfaceID())
}

func testAccNsxtPolicyBareMetalServerTagsNonExistentTemplate() string {
	return `
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "99999999-9999-9999-9999-999999999999"
}
`
}

func testAccNsxtPolicyBareMetalServerInterfaceTagsNonExistentTemplate() string {
	return `
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "88888888-8888-8888-8888-888888888888"
}
`
}
