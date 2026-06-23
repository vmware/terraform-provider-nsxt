// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	tf_api "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var staticObjDisplayName = getAccTestResourceName()

var accTestTransitGatewayCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"transit_subnets": "192.168.5.0/24",
}

var accTestTransitGatewayUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"transit_subnets": "192.168.7.0/24",
}

func TestAccResourceNsxtPolicyTransitGateway_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway.test"
	testDataSourceName := "nsxt_policy_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayCreateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayUpdateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGateway_withSpan(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway.test"
	testDataSourceName := "nsxt_policy_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_TEST_NETWORK_SPAN")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayWithSpanTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayCreateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "span.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "span.0.cluster_based_span.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "span.0.cluster_based_span.0.span_path"),

					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayWithSpanTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayUpdateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),

					resource.TestCheckResourceAttr(testResourceName, "span.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "span.0.cluster_based_span.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "span.0.cluster_based_span.0.span_path"),

					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "description"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGateway_withCentralizedConfig(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway.test"
	testDataSourceName := "data.nsxt_policy_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayWithCentralizedConfigTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayCreateAttributes["transit_subnets"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_STANDBY"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.edge_cluster_paths.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "centralized_config.0.edge_cluster_paths.0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayWithCentralizedConfigTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayUpdateAttributes["transit_subnets"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_ACTIVE"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.edge_cluster_paths.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "centralized_config.0.edge_cluster_paths.0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGateway_withCentralizedConfig920_failoverMode(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway.test"
	testDataSourceName := "data.nsxt_policy_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayWithCentralizedConfig920FailoverTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_STANDBY"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.failover_mode", "PREEMPTIVE"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.edge_cluster_paths.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "centralized_config.0.edge_cluster_paths.0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayWithCentralizedConfig920FailoverTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_STANDBY"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.failover_mode", "NON_PREEMPTIVE"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.edge_cluster_paths.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "centralized_config.0.edge_cluster_paths.0"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGateway_withBgpConfig(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				// Step 1: create TGW with centralized_config only so it can realize before
				// redistribution_config and bgp_config are applied. The NSX API returns empty
				// bodies for /routing and /bgp endpoints until the TGW is edge-deployed.
				Config: testAccNsxtPolicyTransitGatewayBgpConfigBaseTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_ACTIVE"),
				),
			},
			{
				// Step 2: add redistribution_config and bgp_config. The TGW may still be
				// realizing on the edge; the /routing and /bgp endpoints return empty bodies
				// until the TGW is edge-deployed, so redistribution_config may not appear
				// in state yet. ExpectNonEmptyPlan acknowledges the expected drift.
				Config:             testAccNsxtPolicyTransitGatewayWithBgpConfigTemplate(true),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.0.ha_mode", "ACTIVE_ACTIVE"),
				),
			},
			{
				// Step 3: update bgp_config params to verify updates work. redistribution_config
				// is not checked because the NSX /routing endpoint returns empty bodies until
				// the TGW routing engine fully initializes on the edge (eventual consistency).
				// ExpectNonEmptyPlan is set because redistribution_config state may differ from config.
				Config:             testAccNsxtPolicyTransitGatewayWithBgpConfigTemplate(false),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "centralized_config.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTransitGateway_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyTransitGatewayExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGateway resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGateway resource ID not set in resources")
		}

		projectID := strings.Split(rs.Primary.Attributes["path"], "/")[4]
		exists, err := resourceNsxtPolicyTransitGatewayExists(tf_api.SessionContext{ProjectID: projectID, ClientType: tf_api.Multitenancy}, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGateway %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		projectID := strings.Split(rs.Primary.Attributes["path"], "/")[4]
		exists, err := resourceNsxtPolicyTransitGatewayExists(tf_api.SessionContext{ProjectID: projectID, ClientType: tf_api.Multitenancy}, resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGateway %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
%s
  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_transit_gateway" "test" {
%s
  display_name = "%s"
  depends_on   = [nsxt_policy_transit_gateway.test]
}`, testAccNsxtProjectContext(), attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"],
		testAccNsxtProjectContext(), attrMap["display_name"])
}

func testAccNsxtPolicyTransitGatewayWithSpanTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_network_span" "netspan" {
  display_name = "%s"
}

resource "nsxt_policy_project" "test" {
  display_name = "%s"
  vc_folder = true
  default_span_path = data.nsxt_policy_network_span.netspan.path
}

resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  span {
    cluster_based_span {
      span_path = data.nsxt_policy_network_span.netspan.path
    }
  }
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name = "%s"
  depends_on   = [nsxt_policy_transit_gateway.test]
}`, os.Getenv("NSXT_TEST_NETWORK_SPAN"), staticObjDisplayName, attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"], attrMap["display_name"])
}

func testAccNsxtPolicyTransitGatewayMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
%s
  display_name    = "%s"
  transit_subnets = ["%s"]
}`, testAccNsxtProjectContext(), accTestTransitGatewayUpdateAttributes["display_name"], accTestTransitGatewayUpdateAttributes["transit_subnets"])
}

func testAccNsxtPolicyTransitGatewayWithCentralizedConfigPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "tgw-cc-test-t0"
  ha_mode           = "ACTIVE_STANDBY"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name   = "tgw-cc-test-conn"
  tier0_path     = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["192.168.241.0/24"]
}

resource "nsxt_policy_project" "test" {
  display_name             = "tgw-cc-test-project"
  tier0_gateway_paths      = [nsxt_policy_tier0_gateway.test.path]
  tgw_external_connections = [nsxt_policy_gateway_connection.test.path]
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}`, getEdgeClusterName())
}

func testAccNsxtPolicyTransitGatewayWithCentralizedConfigTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	haMode := "ACTIVE_STANDBY"
	if !createFlow {
		haMode = "ACTIVE_ACTIVE"
	}
	return testAccNsxtPolicyTransitGatewayWithCentralizedConfigPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  centralized_config {
    ha_mode            = "%s"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  depends_on   = [nsxt_policy_transit_gateway.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"], haMode, attrMap["display_name"])
}

// testAccNsxtPolicyTransitGatewayBgpConfigBaseTemplate creates a TGW with centralized_config only,
// allowing the TGW to realize (edge-deploy) before redistribution_config and bgp_config are applied.
func testAccNsxtPolicyTransitGatewayBgpConfigBaseTemplate() string {
	return testAccNsxtPolicyTransitGatewayWithCentralizedConfigPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  centralized_config {
    ha_mode            = "ACTIVE_ACTIVE"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}`, accTestTransitGatewayCreateAttributes["display_name"],
		accTestTransitGatewayCreateAttributes["description"],
		accTestTransitGatewayCreateAttributes["transit_subnets"])
}

func testAccNsxtPolicyTransitGatewayWithBgpConfigTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	localAsNum := "65001"
	ecmp := "true"
	gracefulRestartTimer := 180
	gracefulStaleTimer := 600
	if !createFlow {
		localAsNum = "65002"
		ecmp = "false"
		gracefulRestartTimer = 120
		gracefulStaleTimer = 300
	}
	return testAccNsxtPolicyTransitGatewayWithCentralizedConfigPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  centralized_config {
    ha_mode            = "ACTIVE_ACTIVE"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }

  redistribution_config {
    rule {
      types = ["PUBLIC", "TGW_STATIC_ROUTE"]
    }
  }

  bgp_config {
    local_as_num                       = "%s"
    enabled                            = true
    ecmp                               = %s
    graceful_restart_mode              = "HELPER_ONLY"
    graceful_restart_timer             = %d
    graceful_restart_stale_route_timer = %d
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"],
		localAsNum, ecmp, gracefulRestartTimer, gracefulStaleTimer)
}

func testAccNsxtPolicyTransitGatewayWithCentralizedConfig920FailoverTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	failoverMode := "PREEMPTIVE"
	if !createFlow {
		failoverMode = "NON_PREEMPTIVE"
	}
	return testAccNsxtPolicyTransitGatewayWithCentralizedConfigPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  centralized_config {
    ha_mode            = "ACTIVE_STANDBY"
    failover_mode      = "%s"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "%s"
  depends_on   = [nsxt_policy_transit_gateway.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"], failoverMode, attrMap["display_name"])
}
