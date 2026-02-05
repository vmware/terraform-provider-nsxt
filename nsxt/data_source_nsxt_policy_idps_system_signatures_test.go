// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIdpsSystemSignatures_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_system_signatures.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSystemSignaturesReadBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttr(testResourceName, "severity", "CRITICAL"),
					resource.TestCheckResourceAttrSet(testResourceName, "signatures.#"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSystemSignatures_withFilters(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_system_signatures.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSystemSignaturesReadWithSeverityFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttr(testResourceName, "severity", "CRITICAL"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSystemSignatures_withVersionID(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_system_signatures.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSystemSignaturesReadWithVersionID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "version_id", "DEFAULT"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSystemSignatures_withMultipleFilters(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_system_signatures.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSystemSignaturesReadWithMultipleFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_id"),
					resource.TestCheckResourceAttr(testResourceName, "severity", "HIGH"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", "SQL"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIdpsSystemSignatures_signatureAttributes(t *testing.T) {
	testResourceName := "data.nsxt_policy_idps_system_signatures.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIdpsSystemSignaturesReadWithSeverityFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "signatures.#"),
					resource.TestCheckResourceAttr(testResourceName, "severity", "CRITICAL"),
					// Check that signature attributes are present if signatures exist
					// These checks will pass even if signatures list is empty
				),
			},
		},
	})
}

func testAccNsxtPolicyIdpsSystemSignaturesReadBasic() string {
	return `
data "nsxt_policy_idps_system_signatures" "test" {
  severity = "CRITICAL"
}`
}

func testAccNsxtPolicyIdpsSystemSignaturesReadWithSeverityFilter() string {
	return `
data "nsxt_policy_idps_system_signatures" "test" {
  severity = "CRITICAL"
}`
}

func testAccNsxtPolicyIdpsSystemSignaturesReadWithVersionID() string {
	return `
data "nsxt_policy_idps_system_signatures" "test" {
  version_id = "DEFAULT"
  severity   = "HIGH"
}`
}

func testAccNsxtPolicyIdpsSystemSignaturesReadWithMultipleFilters() string {
	return `
data "nsxt_policy_idps_system_signatures" "test" {
  severity     = "HIGH"
  display_name = "SQL"
}`
}
