// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyParentIntrusionServicePolicy_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_parent_intrusion_service_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Parent Policy"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckNoResourceAttr(testResourceName, "rule"), // Should NOT have rules
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyParentIntrusionServicePolicy_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_parent_intrusion_service_policy.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Parent Policy"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckNoResourceAttr(testResourceName, "rule"), // Should NOT have rules
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyParentIntrusionServicePolicy_byCategory(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_parent_intrusion_service_policy.by_category"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyReadByCategory(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Parent Policy"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckNoResourceAttr(testResourceName, "rule"), // Should NOT have rules
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyParentIntrusionServicePolicy_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_parent_intrusion_service_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyReadByNameMultitenancy(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Parent Policy"),
					resource.TestCheckResourceAttr(testResourceName, "category", "ThreatRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckNoResourceAttr(testResourceName, "rule"), // Should NOT have rules
				),
			},
		},
	})
}

func testAccNsxtPolicyParentIntrusionServicePolicyReadByName(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "parent_ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test Parent Policy"
  category        = "ThreatRules"
  sequence_number = 3
}

data "nsxt_policy_parent_intrusion_service_policy" "by_name" {
  display_name = nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds.display_name
  depends_on   = [nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyParentIntrusionServicePolicyReadByCategory(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "parent_ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test Parent Policy"
  category        = "ThreatRules"
  sequence_number = 3
}

data "nsxt_policy_parent_intrusion_service_policy" "by_category" {
  display_name = nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds.display_name
  category     = "ThreatRules"
  depends_on   = [nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyParentIntrusionServicePolicyReadByID(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "parent_ids_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test Parent Policy"
  category        = "ThreatRules"
  sequence_number = 3
}

data "nsxt_policy_parent_intrusion_service_policy" "by_id" {
  id         = nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds.id
  depends_on = [nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds]
}`, name)
}

func testAccNsxtPolicyParentIntrusionServicePolicyReadByNameMultitenancy(name string) string {
	context := testAccNsxtPolicyMultitenancyContext()
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "parent_ids_policy_for_ds" {
%s
  display_name    = "%s"
  description     = "Acceptance Test Parent Policy"
  category        = "ThreatRules"
  sequence_number = 3
}

data "nsxt_policy_parent_intrusion_service_policy" "by_name" {
%s
  display_name = nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds.display_name
  depends_on   = [nsxt_policy_parent_intrusion_service_policy.parent_ids_policy_for_ds]
}`, context, name, context)
}
