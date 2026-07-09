// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	enforcement_points "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
)

const testVNAClusterRealizationResourceName = "data.nsxt_policy_virtual_network_appliance_cluster_realization.test"

func TestAccResourceNsxtPolicyVirtualNetworkApplianceCluster_basic(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance_cluster.test"
	displayName := getAccTestResourceName()
	updatedDisplayName := displayName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.1")
			testAccEnvDefined(t, "NSXT_TEST_OVERLAY_TRANSPORT_ZONE")
			testAccNsxtExtraCoverage(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterCreateTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceClusterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "appliance_form_factor", "MEDIUM"),
					resource.TestCheckResourceAttr(testResourceName, "service_type", "VPC_SERVICES"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testVNAClusterRealizationResourceName, "state", "SUCCESS"),
				),
			},
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterUpdateTemplate(updatedDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceClusterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "appliance_form_factor", "LARGE"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_configuration.#", "1"),
					resource.TestCheckResourceAttr(testVNAClusterRealizationResourceName, "state", "SUCCESS"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVirtualNetworkApplianceCluster_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance_cluster.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.1")
			testAccEnvDefined(t, "NSXT_TEST_OVERLAY_TRANSPORT_ZONE")
			testAccNsxtExtraCoverage(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterCreateTemplate(displayName),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
				ImportStateVerifyIgnore: []string{"revision"},
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VirtualNetworkApplianceCluster resource %s not found in resources", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("VirtualNetworkApplianceCluster resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta())
		client := enforcement_points.NewVirtualNetworkApplianceClustersClient(connector)
		_, err := client.Get("default", "default", id)
		if err != nil {
			return fmt.Errorf("error while retrieving VirtualNetworkApplianceCluster ID %s: %v", id, err)
		}

		return nil
	}
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta())
		client := enforcement_points.NewVirtualNetworkApplianceClustersClient(connector)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "nsxt_policy_virtual_network_appliance_cluster" {
				continue
			}

			id := rs.Primary.ID
			_, err := client.Get("default", "default", id)
			if err == nil {
				return fmt.Errorf("VirtualNetworkApplianceCluster %s still exists", id)
			}
			if !isNotFoundError(err) {
				return fmt.Errorf("unexpected error checking VirtualNetworkApplianceCluster %s: %v", id, err)
			}
		}

		return nil
	}
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterCreateTemplate(displayName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s"
  description           = "Acceptance test cluster"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.test.path
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_virtual_network_appliance_cluster_realization" "test" {
  path = nsxt_policy_virtual_network_appliance_cluster.test.path
}
`, getOverlayTransportZoneName(), displayName)
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterUpdateTemplate(displayName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s"
  description           = "Acceptance test cluster - updated"
  appliance_form_factor = "LARGE"
  service_type          = "VPC_SERVICES"

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.test.path
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}

data "nsxt_policy_virtual_network_appliance_cluster_realization" "test" {
  path = nsxt_policy_virtual_network_appliance_cluster.test.path
}
`, getOverlayTransportZoneName(), displayName)
}

// TestAccResourceNsxtPolicyVirtualNetworkApplianceCluster_noAdvancedConfig covers
// the idempotency regression reported in bug 3713480: when advanced_configuration
// is omitted from the manifest, NSX auto-populates server-side defaults for
// high_availability_profile, overlay_transport_zone_path and
// core_allocation_profile. A second plan must detect zero drift.
func TestAccResourceNsxtPolicyVirtualNetworkApplianceCluster_noAdvancedConfig(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance_cluster.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.1")
			testAccNsxtExtraCoverage(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterNoAdvancedConfigTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceClusterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "appliance_form_factor", "MEDIUM"),
					resource.TestCheckResourceAttr(testResourceName, "service_type", "VPC_SERVICES"),
				),
			},
			{
				// Idempotency check: re-planning with the same config (no
				// advanced_configuration) must produce no changes even after
				// NSX has auto-populated the block with server-side defaults.
				Config:             testAccNsxtPolicyVirtualNetworkApplianceClusterNoAdvancedConfigTemplate(displayName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterNoAdvancedConfigTemplate(displayName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s"
  description           = "Acceptance test cluster - no advanced configuration"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"
}
`, displayName)
}

func TestAccResourceNsxtPolicyVirtualNetworkApplianceCluster_coreAllocationProfile(t *testing.T) {
	testResourceName := "nsxt_policy_virtual_network_appliance_cluster.test"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
			testAccEnvDefined(t, "NSXT_TEST_OVERLAY_TRANSPORT_ZONE")
			testAccNsxtExtraCoverage(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyVirtualNetworkApplianceClusterCheckDestroy(testResourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVirtualNetworkApplianceClusterCoreAllocationTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVirtualNetworkApplianceClusterExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "advanced_configuration.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_configuration.0.core_allocation_profile", "L4LBSERVICE"),
					resource.TestCheckResourceAttr(testVNAClusterRealizationResourceName, "state", "SUCCESS"),
				),
			},
		},
	})
}

func testAccNsxtPolicyVirtualNetworkApplianceClusterCoreAllocationTemplate(displayName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "test" {
  display_name          = "%s"
  description           = "Acceptance test cluster - core allocation profile"
  appliance_form_factor = "LARGE"
  service_type          = "VPC_SERVICES"

  advanced_configuration {
    core_allocation_profile     = "L4LBSERVICE"
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.test.path
  }
}

data "nsxt_policy_virtual_network_appliance_cluster_realization" "test" {
  path = nsxt_policy_virtual_network_appliance_cluster.test.path
}
`, getOverlayTransportZoneName(), displayName)
}
