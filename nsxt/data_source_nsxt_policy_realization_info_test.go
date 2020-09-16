/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1DataSource(t *testing.T) {
	resourceDataType := "nsxt_policy_tier1_gateway"
	resourceName := "terraform_test_tier1_1"
	entityType := ""
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(resourceName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier1GatewayCreate(resourceName); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttrSet(testResourceName, "entity_type"),
					resource.TestCheckResourceAttrSet(testResourceName, "realized_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1DataSourceEntity(t *testing.T) {
	resourceDataType := "nsxt_policy_tier1_gateway"
	resourceName := "terraform_test_tier1_2"
	entityType := "RealizedLogicalRouter"
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(resourceName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier1GatewayCreate(resourceName); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "entity_type", entityType),
					resource.TestCheckResourceAttrSet(testResourceName, "realized_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1Resource(t *testing.T) {
	resourceType := "nsxt_policy_tier1_gateway"
	resourceName := "terraform_test_tier1_3"
	entityType := "RealizedLogicalRouter"
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0"); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRealizationInfoReadResourceTemplate(resourceType, resourceName, entityType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttr(testResourceName, "entity_type", entityType),
					resource.TestCheckResourceAttrSet(testResourceName, "realized_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_errorState(t *testing.T) {
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceErrorTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "ERROR"),
					resource.TestCheckResourceAttrSet(testResourceName, "entity_type"),
					resource.TestCheckResourceAttr(testResourceName, "realized_id", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_gmServiceDataSource(t *testing.T) {
	resourceDataType := "nsxt_policy_service"
	resourceName := "DNS"
	entityType := ""
	testResourceName := "data.nsxt_policy_realization_info.realization_info"
	site := getTestSiteName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_SITE_NAME")
			testAccOnlyGlobalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGMRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType, site),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "state", "REALIZED"),
					resource.TestCheckResourceAttrSet(testResourceName, "entity_type"),
					resource.TestCheckResourceAttrSet(testResourceName, "realized_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyGMRealizationInfoReadDataSourceTemplate(resourceDataType string, resourceName string, entityType string, site string) string {
	return fmt.Sprintf(`
data "%s" "policy_resource" {
  display_name = "%s"
}

data "nsxt_policy_site" "test" {
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = data.%s.policy_resource.path
  entity_type = "%s"
  site_path = data.nsxt_policy_site.test.path
}`, resourceDataType, resourceName, site, resourceDataType, entityType)
}

func testAccNsxtPolicyRealizationInfoReadDataSourceErrorTemplate() string {
	return `
resource "nsxt_policy_ip_pool" "test" {
  display_name = "tfippool1"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name = "tfssnet1"
  pool_path    = nsxt_policy_ip_pool.test.path
  cidr         = "12.12.12.0/24"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}

resource "nsxt_policy_ip_address_allocation" "test" {
  display_name  = "tfipallocationerror"
  pool_path     = nsxt_policy_ip_pool.test.path
  allocation_ip = "12.12.12.21"
  depends_on    = [nsxt_policy_ip_pool_static_subnet.test]
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_ip_address_allocation.test.path
}
`
}

func testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType string, resourceName string, entityType string) string {
	return fmt.Sprintf(`
data "%s" "policy_resource" {
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = data.%s.policy_resource.path
  entity_type = "%s"
}`, resourceDataType, resourceName, resourceDataType, entityType)
}

func testAccNsxtPolicyRealizationInfoReadResourceTemplate(resourceType string, resourceName string, entityType string) string {
	return fmt.Sprintf(`
resource "%s" "policy_resource" {
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = %s.policy_resource.path
  entity_type = "%s"
}`, resourceType, resourceName, resourceType, entityType)
}
