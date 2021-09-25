/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
)

var accTestPolicyLBVirtualServerCreateAttributes = map[string]string{
	"display_name":                getAccTestResourceName(),
	"description":                 "terraform created",
	"access_log_enabled":          "true",
	"default_pool_member_ports":   "777-778",
	"enabled":                     "true",
	"ip_address":                  "1.1.1.1",
	"max_concurrent_connections":  "2",
	"max_new_connection_rate":     "2",
	"ports":                       "80-81",
	"certificate_chain_depth":     "3",
	"log_significant_event_only":  "true",
	"access_list_control_enabled": "true",
	"access_list_control_action":  "ALLOW",
}

var accTestPolicyLBVirtualServerUpdateAttributes = map[string]string{
	"display_name":                getAccTestResourceName(),
	"description":                 "terraform updated",
	"access_log_enabled":          "false",
	"default_pool_member_ports":   "555",
	"enabled":                     "false",
	"ip_address":                  "1.1.2.1",
	"max_concurrent_connections":  "5",
	"max_new_connection_rate":     "5",
	"ports":                       "443",
	"certificate_chain_depth":     "2",
	"log_significant_event_only":  "false",
	"access_list_control_enabled": "false",
	"access_list_control_action":  "DROP",
}

// TODO: add service_path to the test when backend bug is fixed
func TestAccResourceNsxtPolicyLBVirtualServer_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBVirtualServerCheckDestroy(state, accTestPolicyLBVirtualServerCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBVirtualServerTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBVirtualServerCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "access_log_enabled", accTestPolicyLBVirtualServerCreateAttributes["access_log_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "application_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.0", accTestPolicyLBVirtualServerCreateAttributes["default_pool_member_ports"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyLBVirtualServerCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "max_concurrent_connections", accTestPolicyLBVirtualServerCreateAttributes["max_concurrent_connections"]),
					resource.TestCheckResourceAttr(testResourceName, "max_new_connection_rate", accTestPolicyLBVirtualServerCreateAttributes["max_new_connection_rate"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestPolicyLBVirtualServerCreateAttributes["ports"]),
					resource.TestCheckResourceAttrSet(testResourceName, "sorry_pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "persistence_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBVirtualServerUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "access_log_enabled", accTestPolicyLBVirtualServerUpdateAttributes["access_log_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "application_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyLBVirtualServerUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "max_concurrent_connections", accTestPolicyLBVirtualServerUpdateAttributes["max_concurrent_connections"]),
					resource.TestCheckResourceAttr(testResourceName, "max_new_connection_rate", accTestPolicyLBVirtualServerUpdateAttributes["max_new_connection_rate"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestPolicyLBVirtualServerUpdateAttributes["ports"]),
					resource.TestCheckResourceAttrSet(testResourceName, "sorry_pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "persistence_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "pool_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "sorry_pool_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "persistence_profile_path", ""),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

// TODO: add CRLs and other certificates
func TestAccResourceNsxtPolicyLBVirtualServer_withSSL(t *testing.T) {
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_CERTIFICATE_NAME")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBVirtualServerCheckDestroy(state, accTestPolicyLBVirtualServerCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBVirtualServerSSLTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.0.client_auth", "IGNORE"),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.0.certificate_chain_depth", accTestPolicyLBVirtualServerCreateAttributes["certificate_chain_depth"]),
					resource.TestCheckResourceAttrSet(testResourceName, "client_ssl.0.default_certificate_path"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.server_auth", "IGNORE"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.certificate_chain_depth", accTestPolicyLBVirtualServerCreateAttributes["certificate_chain_depth"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerSSLTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.0.client_auth", "IGNORE"),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.0.certificate_chain_depth", accTestPolicyLBVirtualServerUpdateAttributes["certificate_chain_depth"]),
					resource.TestCheckResourceAttrSet(testResourceName, "client_ssl.0.default_certificate_path"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.server_auth", "IGNORE"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.certificate_chain_depth", accTestPolicyLBVirtualServerUpdateAttributes["certificate_chain_depth"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "client_ssl.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLBVirtualServer_withAccessList(t *testing.T) {
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBVirtualServerCheckDestroy(state, accTestPolicyLBVirtualServerCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBVirtualServerWithAccessList(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestPolicyLBVirtualServerCreateAttributes["ports"]),
					resource.TestCheckResourceAttr(testResourceName, "log_significant_event_only", accTestPolicyLBVirtualServerCreateAttributes["log_significant_event_only"]),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.0.enabled", accTestPolicyLBVirtualServerCreateAttributes["access_list_control_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.0.action", accTestPolicyLBVirtualServerCreateAttributes["access_list_control_action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "access_list_control.0.group_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerWithAccessList(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestPolicyLBVirtualServerUpdateAttributes["ports"]),
					resource.TestCheckResourceAttr(testResourceName, "log_significant_event_only", accTestPolicyLBVirtualServerUpdateAttributes["log_significant_event_only"]),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.0.enabled", accTestPolicyLBVirtualServerUpdateAttributes["access_list_control_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.0.action", accTestPolicyLBVirtualServerUpdateAttributes["access_list_control_action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "access_list_control.0.group_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "access_list_control.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "pool_path", ""),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLBVirtualServer_withRules(t *testing.T) {
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBVirtualServerCheckDestroy(state, accTestPolicyLBVirtualServerCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBVirtualServerWithRules(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", "80"),

					resource.TestCheckResourceAttr(testResourceName, "rule.#", "14"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyLBVirtualServerWithRules(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBVirtualServerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyLBVirtualServerUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", "80"),

					resource.TestCheckResourceAttr(testResourceName, "rule.#", "14"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyLBVirtualServer_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBVirtualServerCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBVirtualServerMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBVirtualServerExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := infra.NewDefaultLbVirtualServersClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBVirtualServer resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBVirtualServer resource ID not set in resources")
		}

		_, err := nsxClient.Get(resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy LBVirtualServer ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyLBVirtualServerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := infra.NewDefaultLbVirtualServersClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_virtual_server" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(resourceID)
		if err == nil {
			return fmt.Errorf("Policy LBVirtualServer %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBVirtualServerTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBVirtualServerCreateAttributes
	} else {
		attrMap = accTestPolicyLBVirtualServerUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_lb_app_profile" "default_tcp"{
  type = "TCP"
}

data "nsxt_policy_lb_persistence_profile" "default" {
  type = "SOURCE_IP"
}

resource "nsxt_policy_lb_pool" "pool" {
  display_name = "terraform-vs-test-pool"
}

resource "nsxt_policy_lb_pool" "sorry_pool" {
  display_name = "terraform-vs-test-sorry-pool"
}

resource "nsxt_policy_lb_virtual_server" "test" {
  display_name               = "%s"
  description                = "%s"
  access_log_enabled         = %s
  application_profile_path   = data.nsxt_policy_lb_app_profile.default_tcp.path
  default_pool_member_ports  = ["%s"]
  enabled                    = %s
  ip_address                 = "%s"
  max_concurrent_connections = %s
  max_new_connection_rate    = %s
  ports                      = ["%s"]
  pool_path                  = nsxt_policy_lb_pool.pool.path
  sorry_pool_path            = nsxt_policy_lb_pool.sorry_pool.path
  persistence_profile_path   = data.nsxt_policy_lb_persistence_profile.default.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_virtual_server.test.path
}`, attrMap["display_name"], attrMap["description"], attrMap["access_log_enabled"], attrMap["default_pool_member_ports"], attrMap["enabled"], attrMap["ip_address"], attrMap["max_concurrent_connections"], attrMap["max_new_connection_rate"], attrMap["ports"])
}

func testAccNsxtPolicyLBVirtualServerMinimalistic() string {
	attrMap := accTestPolicyLBVirtualServerCreateAttributes
	return fmt.Sprintf(`
data "nsxt_policy_lb_app_profile" "default_http"{
    type = "HTTP"
}

resource "nsxt_policy_lb_pool" "pool" {
    display_name = "terraform-vs-test-pool"
}

resource "nsxt_policy_lb_pool" "sorry_pool" {
    display_name = "terraform-vs-test-sorry-pool"
}

resource "nsxt_policy_lb_virtual_server" "test" {
  display_name             = "%s"
  application_profile_path = data.nsxt_policy_lb_app_profile.default_http.path
  ip_address               = "%s"
  ports                    = ["%s"]
}

resource "nsxt_policy_group" "group1" {
  display_name = "terraform-vs-test-group1"
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_virtual_server.test.path
}`, accTestPolicyLBVirtualServerUpdateAttributes["display_name"], attrMap["ip_address"], accTestPolicyLBVirtualServerUpdateAttributes["ports"])
}

func testAccNsxtPolicyLBVirtualServerSSLTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBVirtualServerCreateAttributes
	} else {
		attrMap = accTestPolicyLBVirtualServerUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_certificate" "test" {
    display_name = "%s"
}
data "nsxt_policy_lb_app_profile" "default_tcp"{
    type = "TCP"
}

data "nsxt_policy_lb_client_ssl_profile" "default" {
     display_name = "default-balanced-client-ssl-profile"
}

data "nsxt_policy_lb_server_ssl_profile" "default" {
     display_name = "default-balanced-server-ssl-profile"
}

resource "nsxt_policy_lb_virtual_server" "test" {
  display_name             = "%s"
  application_profile_path = data.nsxt_policy_lb_app_profile.default_tcp.path
  ip_address               = "%s"
  ports                    = ["%s"]
  client_ssl {
      client_auth              = "IGNORE"
      certificate_chain_depth  = %s
      ssl_profile_path         = data.nsxt_policy_lb_client_ssl_profile.default.path
      default_certificate_path = data.nsxt_policy_certificate.test.path
  }

  server_ssl {
      server_auth             = "IGNORE"
      certificate_chain_depth = %s
      ssl_profile_path        = data.nsxt_policy_lb_server_ssl_profile.default.path
  }
}`, getTestCertificateName(false), attrMap["display_name"], attrMap["ip_address"], attrMap["ports"], attrMap["certificate_chain_depth"], attrMap["certificate_chain_depth"])
}

func testAccNsxtPolicyLBVirtualServerWithAccessList(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBVirtualServerCreateAttributes
	} else {
		attrMap = accTestPolicyLBVirtualServerUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_lb_app_profile" "default_tcp"{
  type = "TCP"
}

data "nsxt_policy_lb_persistence_profile" "default" {
  type = "SOURCE_IP"
}

resource "nsxt_policy_lb_pool" "pool" {
  display_name = "terraform-vs-test-pool"
}

resource "nsxt_policy_group" "group1" {
  display_name = "terraform-vs-test-group1"
}

resource "nsxt_policy_lb_virtual_server" "test" {
  display_name               = "%s"
  access_log_enabled         = true
  application_profile_path   = data.nsxt_policy_lb_app_profile.default_tcp.path
  enabled                    = true
  ip_address                 = "%s"
  ports                      = ["%s"]
  pool_path                  = nsxt_policy_lb_pool.pool.path
  persistence_profile_path   = data.nsxt_policy_lb_persistence_profile.default.path
  log_significant_event_only = %s
  access_list_control {
      action     = "%s"
      enabled    = %s
      group_path = nsxt_policy_group.group1.path
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_virtual_server.test.path
}`, attrMap["display_name"], attrMap["ip_address"], attrMap["ports"], attrMap["log_significant_event_only"], attrMap["access_list_control_action"], attrMap["access_list_control_enabled"])
}

func testAccNsxtPolicyLBVirtualServerWithRules(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBVirtualServerCreateAttributes
	} else {
		attrMap = accTestPolicyLBVirtualServerUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_lb_app_profile" "default_http"{
  type = "HTTP"
}

data "nsxt_policy_lb_persistence_profile" "default" {
  type = "SOURCE_IP"
}

data "nsxt_policy_lb_persistence_profile" "generic" {
  type = "GENERIC"
}

resource "nsxt_policy_lb_pool" "pool" {
  display_name = "terraform-vs-test-pool"
}

resource "nsxt_policy_group" "group1" {
  display_name = "terraform-vs-test-group1"
}

resource "nsxt_policy_lb_virtual_server" "test" {
  display_name             = "%s"
  application_profile_path = data.nsxt_policy_lb_app_profile.default_http.path
  ip_address               = "%s"
  ports                    = [ "80" ]
  pool_path                  = nsxt_policy_lb_pool.pool.path

  rule {
	display_name = "connection_drop_action_test"
	connection_drop_action {}
	connection_drop_action {}
  } 
  rule {
	display_name = "http_redirect_action_test"
	http_redirect_action {
      redirect_status = "301"
      redirect_url = "dummy"
	}
	http_redirect_action {
      redirect_status = "302"
      redirect_url = "other_dummy"
	}
	http_request_body_condition {
	  body_value = "xyz"	
	  match_type = "REGEX"
	}
  }
  rule {
	display_name = "http_reject_action_test"
	http_reject_action {
	  reply_status = "404"
	}
	http_reject_action {
	  reply_status = "400"
	}
	http_request_cookie_condition {
	  cookie_name = "is"
	  cookie_value = "brownie"
	  match_type = "REGEX"
	}
  }
  rule {
	display_name = "http_request_header_delete_action_test"
	phase = "HTTP_REQUEST_REWRITE"
	http_request_header_delete_action {
	  header_name = "clever"
	}
	http_request_header_delete_action {
	  header_name = "smart"
	}
	http_request_header_condition {
	  header_name = "my_head"	
	  header_value = "has_hair"
	  match_type = "REGEX"
	}
  }
  rule {
	display_name = "http_request_header_rewrite_action_test"
	phase = "HTTP_REQUEST_REWRITE"
	http_request_header_rewrite_action {
	  header_name = "clever"
	  header_value = "and"
	}
	http_request_header_rewrite_action {
	  header_name = "smart"
	  header_value = "!"
	}
	http_request_uri_arguments_condition {
	  uri_arguments = "foo=bar"
	  match_type = "REGEX"	
	}
	http_request_uri_condition {
	  uri = "xyz"
	  match_type = "REGEX"	
	}
  }
  rule {
	display_name = "http_request_uri_rewrite_action_test"
	phase = "HTTP_REQUEST_REWRITE"
	http_request_uri_rewrite_action {
	  uri = "uri_test"
	}
	http_request_uri_rewrite_action {
      uri = "uri_test2"
	  uri_arguments = "foo=bar"
	}
	http_request_version_condition {
	  version = "HTTP_VERSION_1_0"
	}
  }
  rule {
	display_name = "http_response_header_delete_action_test"
	phase = "HTTP_RESPONSE_REWRITE"
	http_response_header_delete_action {
	  header_name = "clever"
    }
    http_response_header_delete_action {
	  header_name = "smart"
    }
	http_response_header_condition {
	  header_name = "clever"
	  header_value = "smart"
	  match_type = "REGEX"
	}
  }
  rule {
	display_name = "http_response_header_rewrite_action_test"
	phase = "HTTP_RESPONSE_REWRITE"
	http_response_header_rewrite_action {
	  header_name = "clever"
	  header_value = "and"
	}
	http_response_header_rewrite_action {
	  header_name = "smart"
	  header_value = "!"
	}
	ip_header_condition {
	  source_address = "1.1.1.1"
	}
  }
  rule {
	display_name = "select_pool_action_test"
	phase = "HTTP_FORWARDING"
	select_pool_action {
      pool_id = nsxt_policy_lb_pool.pool.path
	} 
	ssl_sni_condition {
	  match_type = "REGEX"
	  sni = "HELO"
	} 
  }
  rule {
	display_name = "ssl_mode_selection_action_test"
	phase = "TRANSPORT"
	ssl_mode_selection_action {
      ssl_mode = "SSL_PASSTHROUGH"
	}
  }
  rule {
	display_name = "variable_assignment_action_test"
	phase = "HTTP_ACCESS"
	variable_assignment_action {
      variable_name = "foo"
	  variable_value = "bar"
	}
	variable_condition {
	  variable_name = "my_var"
	  variable_value = "my_value"
	  match_type = "REGEX"	
	}  
	tcp_header_condition {
		source_port = "80"	
	}  
  }
  rule {
	display_name = "variable_persistence_learn_action_test"
	phase = "HTTP_RESPONSE_REWRITE"
	variable_persistence_learn_action {
      persistence_profile_path = data.nsxt_policy_lb_persistence_profile.generic.path
	  variable_hash_enabled = true
	  variable_name = "my_name"
	}  
	http_ssl_condition {
	  session_reused = "IGNORE"
	  used_protocol = "SSL_V2"
	  used_ssl_cipher = "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"	
	  client_certificate_issuer_dn {
		issuer_dn = "something"
		match_type = "REGEX"
	  }
	  client_certificate_subject_dn {
		subject_dn = "something"  
		match_type = "REGEX"
	  }
	  client_supported_ssl_ciphers = [ "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" ]
	}
  }
  rule {
	display_name = "variable_persistence_on_action_test"
	phase = "HTTP_FORWARDING"
	variable_persistence_on_action {
      persistence_profile_path = data.nsxt_policy_lb_persistence_profile.generic.path
	  variable_hash_enabled = false
	  variable_name = "my_name"
	}  
  }
  rule {
	display_name = "jwt_auth_action_test"
	phase = "HTTP_ACCESS"
	jwt_auth_action {
	  key {
		// one of ...
		public_key_content = "xxx"
		//certificate_path = "/path/to/cert"
		// this one only indicates presence, the value is discarded
		//symmetric_key = "dummy_value"
	  }
	  pass_jwt_to_pool = true
	  realm = "realm"
	  // only one token allowed currently
	  tokens = [ "a" ]
	}
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_lb_virtual_server.test.path
}`, attrMap["display_name"], attrMap["ip_address"])
}
