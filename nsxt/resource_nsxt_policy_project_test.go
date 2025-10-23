// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var shortID = getAccTestRandomString(6)

var accTestPolicyProjectCreateAttributes = map[string]string{
	"DisplayName": getAccTestResourceName(),
	"Description": "terraform created",
	"ShortId":     shortID,
}

var accTestPolicyProjectUpdateAttributes = map[string]string{
	"DisplayName": getAccTestResourceName(),
	"Description": "terraform updated",
	"ShortId":     shortID,
}

func getExpectedSiteInfoCount(t *testing.T) string {
	if util.NsxVersion == "" {
		connector, err := testAccGetPolicyConnector()
		if err != nil {
			t.Errorf("Failed to get policy connector")
			return "0"
		}

		err = initNSXVersion(connector)
		if err != nil {
			t.Errorf("Failed to retrieve NSX version")
			return "0"
		}
	}
	if util.NsxVersionHigherOrEqual("4.1.1") {
		return "1"
	}
	return "0"
}

func runChecksNsx410(testResourceName string, attributes map[string]string, expectedValues map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccNsxtPolicyProjectExists(accTestPolicyProjectCreateAttributes["DisplayName"], testResourceName),
		resource.TestCheckResourceAttr(testResourceName, "display_name", attributes["DisplayName"]),
		resource.TestCheckResourceAttr(testResourceName, "description", attributes["Description"]),
		resource.TestCheckResourceAttr(testResourceName, "short_id", attributes["ShortId"]),
		//TODO: add site info validation
		resource.TestCheckResourceAttr(testResourceName, "site_info.#", expectedValues["site_count"]),
		resource.TestCheckResourceAttr(testResourceName, "tier0_gateway_paths.#", expectedValues["t0_count"]),

		resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
		resource.TestCheckResourceAttrSet(testResourceName, "path"),
		resource.TestCheckResourceAttrSet(testResourceName, "revision"),
		resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
	)
}

func runChecksNsx411(testResourceName string, expectedValues map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testResourceName, "external_ipv4_blocks.#", expectedValues["ip_block_count"]),
	)
}

func runChecksNsx420(testResourceName string, expectedValues map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testResourceName, "activate_default_dfw_rules", expectedValues["activate_default_dfw_rules"]),
	)
}

func runChecksNsx900(testResourceName string, expectedValues map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testResourceName, "tgw_external_connections.#", expectedValues["tgw_ext_conn_count"]),
	)
}

func runChecksNsx910(testResourceName string, expectedValues map[string]string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(testResourceName, "non_default_span_paths.#", expectedValues["span_reference_count"]),
		resource.TestCheckResourceAttrSet(testResourceName, "non_default_span_paths.0"),
		resource.TestCheckResourceAttrSet(testResourceName, "default_span_path"),
	)
}

func TestAccResourceNsxtPolicyProject_basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	siteCount := getExpectedSiteInfoCount(t)
	expectedValuesStep1 := map[string]string{
		"t0_count":   "1",
		"site_count": siteCount,
	}
	expectedValuesStep2 := expectedValuesStep1
	expectedValuesStep3 := map[string]string{
		"t0_count":   "0",
		"site_count": siteCount,
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.1.0")
			testAccNSXVersionLessThan(t, "4.1.1")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectTemplate410(true, true),
				Check:  runChecksNsx410(testResourceName, accTestPolicyProjectCreateAttributes, expectedValuesStep1),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate410(false, true),
				Check:  runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep2),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate410(false, false),
				Check:  runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep3),
			},
			{
				Config: testAccNsxtPolicyProjectMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectExists(accTestPolicyProjectCreateAttributes["DisplayName"], testResourceName),
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

func TestAccResourceNsxtPolicyProject_411basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	siteCount := getExpectedSiteInfoCount(t)
	expectedValuesStep1 := map[string]string{
		"t0_count":       "1",
		"ip_block_count": "1",
		"site_count":     siteCount,
	}
	expectedValuesStep2 := expectedValuesStep1
	expectedValuesStep3 := map[string]string{
		"t0_count":       "0",
		"ip_block_count": "0",
		"site_count":     siteCount,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.1.1")
			testAccNSXVersionLessThan(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectTemplate411(true, true, true),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectCreateAttributes, expectedValuesStep1),
					runChecksNsx411(testResourceName, expectedValuesStep1),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate411(false, true, true),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep2),
					runChecksNsx411(testResourceName, expectedValuesStep2),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate411(false, false, false),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep3),
					runChecksNsx411(testResourceName, expectedValuesStep3),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_420basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	siteCount := getExpectedSiteInfoCount(t)
	expectedValuesStep1 := map[string]string{
		"t0_count":                   "1",
		"ip_block_count":             "1",
		"site_count":                 siteCount,
		"activate_default_dfw_rules": "false",
	}
	expectedValuesStep2 := expectedValuesStep1
	expectedValuesStep3 := map[string]string{
		"t0_count":                   "0",
		"ip_block_count":             "0",
		"site_count":                 siteCount,
		"activate_default_dfw_rules": "true",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
			testAccNSXVersionLessThan(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectTemplate420(true, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectCreateAttributes, expectedValuesStep1),
					runChecksNsx411(testResourceName, expectedValuesStep1),
					runChecksNsx420(testResourceName, expectedValuesStep1),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate420(false, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep2),
					runChecksNsx411(testResourceName, expectedValuesStep2),
					runChecksNsx420(testResourceName, expectedValuesStep2),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate420(false, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep3),
					runChecksNsx411(testResourceName, expectedValuesStep3),
					runChecksNsx420(testResourceName, expectedValuesStep3),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_900basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	siteCount := getExpectedSiteInfoCount(t)
	expectedValuesStep1 := map[string]string{
		"t0_count":                   "1",
		"ip_block_count":             "1",
		"tgw_ext_conn_count":         "1",
		"activate_default_dfw_rules": "true",
		"site_count":                 siteCount,
	}
	expectedValuesStep2 := map[string]string{
		"t0_count":                   "1",
		"ip_block_count":             "0",
		"tgw_ext_conn_count":         "1",
		"activate_default_dfw_rules": "false",
		"site_count":                 siteCount,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
			testAccNSXVersionLessThan(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				// Create: Set T0, Ext GW connection, Ext IPv4 Block, activate default DFW
				Config: testAccNsxtPolicyProjectTemplate900(true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectCreateAttributes, expectedValuesStep1),
					runChecksNsx411(testResourceName, expectedValuesStep1),
					runChecksNsx420(testResourceName, expectedValuesStep1),
					runChecksNsx900(testResourceName, expectedValuesStep1),
				),
			},
			{
				// Update: Set T0, Ext GW connection, No Ext IPv4 Block, disable default DFW
				Config: testAccNsxtPolicyProjectTemplate900(false, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep2),
					runChecksNsx411(testResourceName, expectedValuesStep2),
					runChecksNsx420(testResourceName, expectedValuesStep2),
					runChecksNsx900(testResourceName, expectedValuesStep2),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_900defaultSecurityProfile(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectCreateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				// Create: Set T0, Ext GW connection, Ext IPv4 Block, activate default DFW
				Config: testAccNsxtPolicyProjectMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectGetSecurityProfileNSEnabled(testResourceName, false),
				),
			},
			{
				// Update: Set T0, Ext GW connection, No Ext IPv4 Block, disable default DFW
				Config: testAccNsxtPolicyProjectDefaultSecurityPolicy(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectGetSecurityProfileNSEnabled(testResourceName, false),
				),
			},
			{
				// Update: Set T0, Ext GW connection, No Ext IPv4 Block, disable default DFW
				Config: testAccNsxtPolicyProjectDefaultSecurityPolicy(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectGetSecurityProfileNSEnabled(testResourceName, true),
				),
			},
			{
				// Update: Set T0, Ext GW connection, No Ext IPv4 Block, disable default DFW
				Config: testAccNsxtPolicyProjectDefaultSecurityPolicy(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectGetSecurityProfileNSEnabled(testResourceName, false),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_910basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	siteCount := getExpectedSiteInfoCount(t)
	expectedValuesStep1 := map[string]string{
		"t0_count":                   "1",
		"ip_block_count":             "1",
		"tgw_ext_conn_count":         "1",
		"span_reference_count":       "1",
		"activate_default_dfw_rules": "true",
		"site_count":                 siteCount,
	}
	expectedValuesStep2 := map[string]string{
		"t0_count":                   "1",
		"ip_block_count":             "0",
		"tgw_ext_conn_count":         "1",
		"span_reference_count":       "1",
		"activate_default_dfw_rules": "false",
		"site_count":                 siteCount,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				// Create: Set T0, Ext GW connection, Ext IPv4 Block, activate default DFW
				Config: testAccNsxtPolicyProjectTemplate910(true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectCreateAttributes, expectedValuesStep1),
					runChecksNsx411(testResourceName, expectedValuesStep1),
					runChecksNsx420(testResourceName, expectedValuesStep1),
					runChecksNsx900(testResourceName, expectedValuesStep1),
					runChecksNsx910(testResourceName, expectedValuesStep1),
				),
			},
			{
				// Update: Set T0, Ext GW connection, No Ext IPv4 Block, disable default DFW
				Config: testAccNsxtPolicyProjectTemplate910(false, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					runChecksNsx410(testResourceName, accTestPolicyProjectUpdateAttributes, expectedValuesStep2),
					runChecksNsx411(testResourceName, expectedValuesStep2),
					runChecksNsx420(testResourceName, expectedValuesStep2),
					runChecksNsx900(testResourceName, expectedValuesStep2),
					runChecksNsx910(testResourceName, expectedValuesStep2),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_DefaultSpanCheck(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectCreateAttributes["DisplayName"])
		},
		Steps: []resource.TestStep{
			{
				// Create
				Config: testAccNsxtPolicyProjectDefaultSpanCheckTemplate(accTestPolicyProjectCreateAttributes["DisplayName"], true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(
						testAccNsxtCheckSpanPath(testResourceName, false),
						resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyProjectCreateAttributes["DisplayName"]),
						resource.TestCheckResourceAttrSet(testResourceName, "default_span_path"),
					),
				),
			},
			{
				// Update
				Config: testAccNsxtPolicyProjectDefaultSpanCheckTemplate(accTestPolicyProjectCreateAttributes["DisplayName"], false, true),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(
						testAccNsxtCheckSpanPath(testResourceName, true),
						resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyProjectCreateAttributes["DisplayName"]),
						resource.TestCheckResourceAttrSet(testResourceName, "default_span_path"),
					),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyProject_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectMinimalistic(),
			},
			{
				ResourceName: testResourceName,
				ImportState:  true,
			},
		},
	})
}

func testAccNsxtPolicyProjectGetSecurityProfileNSEnabled(resourceName string, expectedVal bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Project resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Project resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		client := projects.NewVpcSecurityProfilesClient(connector)
		obj, err := client.Get(defaultOrgID, resourceID, "default")
		if err != nil {
			return err
		}
		if *obj.NorthSouthFirewall.Enabled != expectedVal {
			return fmt.Errorf("expected NorthSouthFirewall to be %v, instead status is %v", *obj.NorthSouthFirewall.Enabled, expectedVal)
		}
		return nil
	}
}

func testAccNsxtPolicyProjectExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Project resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Project resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyProjectExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy Project %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyProjectCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_project" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyProjectExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy Project %s still exists", displayName)
		}
	}
	return nil
}

// TODO: Visibility should not be set for 410 tests (attribute exists, different enum values since 411)
var templateData = `
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "{{.ClusterName}}"
}

data "nsxt_policy_tier0_gateway" "test" {
  display_name = "{{.T0Name}}"
}

resource "nsxt_policy_ip_block" "test_ip_block" {
  display_name = "test_ip_block"
  cidr         = "10.20.0.0/16"
  {{if .SetIpBlockVisibility}}visibility   = "EXTERNAL"{{end}}
}

{{if .ExternalTGWConnectionPath}}resource "nsxt_policy_gateway_connection" "test_gw_conn" {
  display_name     = "test_gw_conn"
  tier0_path       = {{.T0GwPath}}
  aggregate_routes = ["192.168.240.0/24"]
}{{end}}

{{if .PolicyNetworkSpan}}resource "nsxt_policy_network_span" "netspan" {
  display_name = "test_span"
  exclusive    = true
}{{end}}

resource "nsxt_policy_project" "test" {
  display_name               = "{{.DisplayName}}"
  description                = "{{.Description}}"
  {{if .ShortId}}short_id                   = "{{.ShortId}}"{{end}}
  {{if .T0GwPath}}tier0_gateway_paths        = [{{.T0GwPath}}]{{end}}
  {{if .ExternalIPv4BlockPath}}external_ipv4_blocks       = [{{.ExternalIPv4BlockPath}}]{{end}}
  {{if .ExternalTGWConnectionPath}}tgw_external_connections   = [{{.ExternalTGWConnectionPath}}]{{end}}
  {{if .ActivateDefaultDfwRules}}activate_default_dfw_rules = {{.ActivateDefaultDfwRules}}{{end}}
  {{if .PolicyNetworkSpan}}
  vc_folder = true
  non_default_span_paths = [{{.PolicyNetworkSpan}}]
  {{end}}
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.EC.path]
  }
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`

func getBasicAttrMap(createFlow, includeT0GW bool) map[string]string {
	attrMap := make(map[string]string)
	baseMap := accTestPolicyProjectUpdateAttributes
	if createFlow {
		baseMap = accTestPolicyProjectCreateAttributes
		//do not pass short_id on update
		attrMap["ShortId"] = baseMap["ShortId"]
	}
	attrMap["DisplayName"] = baseMap["DisplayName"]
	attrMap["Description"] = baseMap["Description"]
	if includeT0GW {
		attrMap["T0GwPath"] = "data.nsxt_policy_tier0_gateway.test.path"
	}
	attrMap["T0Name"] = getTier0RouterName()
	attrMap["ClusterName"] = getEdgeClusterName()
	return attrMap
}

func testAccNsxtPolicyProjectTemplate410(createFlow, includeT0GW bool) string {
	attrMap := getBasicAttrMap(createFlow, includeT0GW)
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("testAccNsxtPolicyProjectTemplate410").Parse(templateData)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buffer, attrMap)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func testAccNsxtPolicyProjectTemplate411(createFlow, includeT0GW, includeExternalIPv4Block bool) string {
	attrMap := getBasicAttrMap(createFlow, includeT0GW)
	if includeExternalIPv4Block {
		attrMap["ExternalIPv4BlockPath"] = "nsxt_policy_ip_block.test_ip_block.path"
	}
	attrMap["SetIpBlockVisibility"] = "true"
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("testAccNsxtPolicyProjectTemplate411").Parse(templateData)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buffer, attrMap)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func testAccNsxtPolicyProjectTemplate420(createFlow, includeT0GW, includeExternalIPv4Block, activateDefaultDfwRules bool) string {
	attrMap := getBasicAttrMap(createFlow, includeT0GW)
	if includeExternalIPv4Block {
		attrMap["ExternalIPv4BlockPath"] = "nsxt_policy_ip_block.test_ip_block.path"
	}
	attrMap["SetIpBlockVisibility"] = "true"
	attrMap["ActivateDefaultDfwRules"] = strconv.FormatBool(activateDefaultDfwRules)
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("testAccNsxtPolicyProjectTemplate420").Parse(templateData)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buffer, attrMap)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func testAccNsxtPolicyProjectTemplate900(createFlow, includeT0GW, includeExternalIPv4Block, activateDefaultDfwRules bool) string {
	attrMap := getBasicAttrMap(createFlow, includeT0GW)
	if includeExternalIPv4Block {
		attrMap["ExternalIPv4BlockPath"] = "nsxt_policy_ip_block.test_ip_block.path"
	}
	attrMap["SetIpBlockVisibility"] = "true"
	attrMap["ActivateDefaultDfwRules"] = strconv.FormatBool(activateDefaultDfwRules)
	attrMap["ExternalTGWConnectionPath"] = "nsxt_policy_gateway_connection.test_gw_conn.path"
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("testAccNsxtPolicyProjectTemplate900").Parse(templateData)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buffer, attrMap)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func testAccNsxtPolicyProjectTemplate910(createFlow, includeT0GW, includeExternalIPv4Block, activateDefaultDfwRules bool) string {
	attrMap := getBasicAttrMap(createFlow, includeT0GW)
	if includeExternalIPv4Block {
		attrMap["ExternalIPv4BlockPath"] = "nsxt_policy_ip_block.test_ip_block.path"
	}
	attrMap["SetIpBlockVisibility"] = "true"
	attrMap["ActivateDefaultDfwRules"] = strconv.FormatBool(activateDefaultDfwRules)
	attrMap["ExternalTGWConnectionPath"] = "nsxt_policy_gateway_connection.test_gw_conn.path"
	attrMap["PolicyNetworkSpan"] = "nsxt_policy_network_span.netspan.path"
	buffer := new(bytes.Buffer)
	tmpl, err := template.New("testAccNsxtPolicyProjectTemplate910").Parse(templateData)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buffer, attrMap)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}

func testAccNsxtPolicyProjectMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test" {
  display_name = "%s"

}`, accTestPolicyProjectUpdateAttributes["DisplayName"])
}

func testAccNsxtPolicyProjectDefaultSecurityPolicy(enabled bool) string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test" {
  display_name = "%s"
  default_security_profile {
    north_south_firewall {
      enabled = %s
    }
  }

}`, accTestPolicyProjectCreateAttributes["DisplayName"], strconv.FormatBool(enabled))
}

func testAccNsxtPolicyProjectDefaultSpanCheckTemplate(displayName string, withCustomSpan bool, isUpdate bool) string {
	defaultSpan := "span-default"
	nonDefaultSpan := "span-non-default"
	spanConstruct := constructSpanForProject("nsxt_policy_network_span.def.path", "nsxt_policy_network_span.nondef.path")
	if isUpdate {
		spanConstruct = constructSpanForProject("", "nsxt_policy_network_span.def.path")
	}
	projectTemplate := fmt.Sprintf(`

resource "nsxt_policy_project" "test" {
  display_name = "%s"
  vc_folder = true
  %s
}
`, displayName, spanConstruct)

	return testAccNsxtPolicyNetworkSpan(defaultSpan, nonDefaultSpan) + projectTemplate
}

func constructSpanForProject(customDefaultSpan, nonDefaultSpan string) string {
	if customDefaultSpan == "" {
		return fmt.Sprintf(`
non_default_span_paths = [%s]
`, nonDefaultSpan)
	} else {
		return fmt.Sprintf(`
default_span_path = %s
non_default_span_paths = [%s]
`, customDefaultSpan, nonDefaultSpan)
	}
}

func testAccNsxtPolicyNetworkSpan(defaultSpan, nonDefaultSpan string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_network_span" "def" {
  display_name = "%s"

}

resource "nsxt_policy_network_span" "nondef" {
  display_name = "%s"

}

`, defaultSpan, nonDefaultSpan)
}

func testAccNsxtCheckSpanPath(resourceName string, isSpanNsxDefault bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy project resource %s not found in resources", resourceName)
		}
		resourceID := rs.Primary.ID
		var err error
		var defaultSpan string
		if isSpanNsxDefault {
			defaultSpan, err = getDefaultSpan(connector)
			if err != nil {
				return fmt.Errorf("Failed to get the default span path : %v", err)
			}
		} else {
			defaultSpan = rs.Primary.Attributes["default_span_path"]
		}
		client := infra.NewProjectsClient(connector)
		project, err := client.Get(defaultOrgID, resourceID, nil)
		if err != nil {
			return fmt.Errorf("Error retreiving the project %v", resourceID)
		}
		defaultSpanFromProject := ""
		for _, spanRef := range project.VpcDeploymentScope.SpanReferences {
			if *spanRef.IsDefault {
				defaultSpanFromProject = *spanRef.SpanPath
			}
		}
		if defaultSpan == defaultSpanFromProject {
			return nil
		}
		return fmt.Errorf("Default span setting check failed for project %v", resourceID)
	}
}
