// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIdpsSignatureVersion_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_signature_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionReadBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "state"),
					resource.TestCheckResourceAttrSet(testResourceName, "status"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSignatureVersion_byDisplayName(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_signature_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionReadByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", "DEFAULT"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "state"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSignatureVersion_attributes(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_signature_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSignatureVersionReadBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "state"),
					resource.TestCheckResourceAttrSet(testResourceName, "status"),
					// These may or may not be set depending on NSX version
					// Just verify the attributes exist in schema
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsSignatureVersionReadBasic() string {
	return `
data "nsxt_policy_idps_signature_version" "test" {
  id = "DEFAULT"
}`
}

func testAccNsxtPolicyIdpsSignatureVersionReadByName() string {
	return `
data "nsxt_policy_idps_signature_version" "test" {
  display_name = "DEFAULT"
}`
}
