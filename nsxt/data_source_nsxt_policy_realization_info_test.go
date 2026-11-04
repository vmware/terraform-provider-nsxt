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

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1DataSource(t *testing.T) {
	testAccDataSourceNsxtPolicyRealizationInfoTier1DataSource(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1DataSource_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyRealizationInfoTier1DataSource(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicyRealizationInfoTier1DataSource(t *testing.T, withContext bool, preCheck func()) {
	resourceDataType := "nsxt_policy_tier1_gateway"
	resourceName := getAccTestDataSourceName()
	entityType := ""
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(resourceName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier1GatewayCreate(resourceName); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType, withContext),
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
	testAccDataSourceNsxtPolicyRealizationInfoTier1DataSourceEntity(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1DataSourceEntity_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyRealizationInfoTier1DataSourceEntity(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicyRealizationInfoTier1DataSourceEntity(t *testing.T, withContext bool, preCheck func()) {
	resourceDataType := "nsxt_policy_tier1_gateway"
	resourceName := getAccTestDataSourceName()
	entityType := "RealizedLogicalRouter"
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyTier1GatewayDeleteByName(resourceName)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyTier1GatewayCreate(resourceName); err != nil {
						t.Error(err)
					}
				},
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType, resourceName, entityType, withContext),
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
	testAccDataSourceNsxtPolicyRealizationInfoTier1Resource(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_tier1Resource_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyRealizationInfoTier1Resource(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccDataSourceNsxtPolicyRealizationInfoTier1Resource(t *testing.T, withContext bool, preCheck func()) {
	resourceType := "nsxt_policy_tier1_gateway"
	resourceName := getAccTestDataSourceName()
	entityType := "RealizedLogicalRouter"
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRealizationInfoReadResourceTemplate(resourceType, resourceName, entityType, withContext),
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
	testAccDataSourceNsxtPolicyRealizationInfoErrorState(t, false, 71234, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccDataSourceNsxtPolicyRealizationInfo_errorState_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyRealizationInfoErrorState(t, true, 71235, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

// The ERROR realization state is provoked by creating two segments with the same overlay_id (VNI):
// NSX accepts both, but the second one fails to realize. The overlay_id differs per test, as the
// tests run in parallel against the same NSX Manager.
func testAccDataSourceNsxtPolicyRealizationInfoErrorState(t *testing.T, withContext bool, overlayID int, preCheck func()) {
	testResourceName := "data.nsxt_policy_realization_info.realization_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRealizationInfoReadDataSourceErrorTemplate(withContext, overlayID),
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
	resourceName := "ICMPv6-ALL"
	entityType := ""
	testResourceName := "data.nsxt_policy_realization_info.realization_info"
	site := getTestSiteName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyGlobalManager(t)
			testAccEnvDefined(t, "NSXT_TEST_SITE_NAME")
			testAccPreCheck(t)
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
  path        = data.%s.policy_resource.path
  entity_type = "%s"
  site_path   = data.nsxt_policy_site.test.path
}`, resourceDataType, resourceName, site, resourceDataType, entityType)
}

func testAccNsxtPolicyRealizationInfoReadDataSourceErrorTemplate(withContext bool, overlayID int) string {
	context := ""
	transportZone := ""
	tzSetting := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	} else {
		transportZone = testAccNSXPolicyTransportZoneReadTemplate(getOverlayTransportZoneName(), false, true)
		tzSetting = "transport_zone_path = data.nsxt_policy_transport_zone.test.path"
	}
	return transportZone + fmt.Sprintf(`
resource "nsxt_policy_segment" "first" {
%[1]s
  display_name = "tfrealizationfirst"
  overlay_id   = %[3]d
  %[2]s
}

// Wait for the first segment to be realized, so that it is the one owning the overlay_id and the
// duplicate is the one that fails to realize
data "nsxt_policy_realization_info" "first" {
%[1]s
  path        = nsxt_policy_segment.first.path
  entity_type = "RealizedLogicalSwitch"
}

resource "nsxt_policy_segment" "duplicate" {
%[1]s
  display_name = "tfrealizationduplicate"
  overlay_id   = %[3]d
  %[2]s
  depends_on   = [data.nsxt_policy_realization_info.first]
}

data "nsxt_policy_realization_info" "realization_info" {
%[1]s
  path        = nsxt_policy_segment.duplicate.path
  entity_type = "RealizedLogicalSwitch"
}
`, context, tzSetting, overlayID)
}

func testAccNsxtPolicyRealizationInfoReadDataSourceTemplate(resourceDataType string, resourceName string, entityType string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "%s" "policy_resource" {
%s
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
%s
  path        = data.%s.policy_resource.path
  entity_type = "%s"
}`, resourceDataType, context, resourceName, context, resourceDataType, entityType)
}

func testAccNsxtPolicyRealizationInfoReadResourceTemplate(resourceType string, resourceName string, entityType string, withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "%s" "policy_resource" {
%s
  display_name = "%s"
}

data "nsxt_policy_realization_info" "realization_info" {
%s
  path        = %s.policy_resource.path
  entity_type = "%s"
}`, resourceType, context, resourceName, context, resourceType, entityType)
}
