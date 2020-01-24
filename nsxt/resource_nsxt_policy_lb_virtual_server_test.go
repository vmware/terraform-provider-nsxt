/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"testing"
)

var accTestPolicyLBVirtualServerCreateAttributes = map[string]string{
	"display_name":                "terra-test",
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
	"display_name":                "terra-test-updated",
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
		PreCheck:  func() { testAccPreCheck(t) },
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
		PreCheck:  func() { testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_CERTIFICATE_NAME") },
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
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
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

func TestAccResourceNsxtPolicyLBVirtualServer_importBasic(t *testing.T) {
	name := "terra-test-import"
	testResourceName := "nsxt_policy_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
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
