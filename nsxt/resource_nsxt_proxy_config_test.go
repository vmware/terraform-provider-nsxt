// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/proxy"
)

func TestAccResourceNsxtProxyConfig_basic(t *testing.T) {
	testAccResourceNsxtPolicyProxyConfigBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func testAccResourceNsxtPolicyProxyConfigBasic(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_proxy_config.test"
	server1 := "proxy1.example.com"
	server2 := "proxy2.example.com"
	port1 := "8080"
	port2 := "3128"
	username1 := "user1"
	username2 := "user2"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtProxyConfigCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				// Test with disabled proxy first (no connectivity validation)
				Config: testAccNsxtProxyConfigDisabledWithDetails(name, server1, port1, username1, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtProxyConfigExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scheme", "HTTP"),
					resource.TestCheckResourceAttr(testResourceName, "host", server1),
					resource.TestCheckResourceAttr(testResourceName, "port", port1),
					resource.TestCheckResourceAttr(testResourceName, "username", username1),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				// Test update with different values
				Config: testAccNsxtProxyConfigDisabledWithDetails(updatedName, server2, port2, username2, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtProxyConfigExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "scheme", "HTTP"),
					resource.TestCheckResourceAttr(testResourceName, "host", server2),
					resource.TestCheckResourceAttr(testResourceName, "port", port2),
					resource.TestCheckResourceAttr(testResourceName, "username", username2),
				),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "TelemetryConfigIdentifier",
				// Password is sensitive and not returned by API
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccResourceNsxtProxyConfig_disabled(t *testing.T) {
	testAccResourceNsxtPolicyProxyConfigDisabled(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func testAccResourceNsxtPolicyProxyConfigDisabled(t *testing.T, withContext bool, preCheck func()) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_proxy_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtProxyConfigCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtProxyConfigDisabled(name, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtProxyConfigExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtProxyConfigExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Proxy Config resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Proxy Config resource ID not set in resources")
		}

		// Proxy config is a singleton - just verify we can read it
		client := proxy.NewConfigClient(connector)
		_, err := client.Get()
		if err != nil {
			return fmt.Errorf("Policy Proxy Config %s does not exist: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtProxyConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_proxy_config" {
			continue
		}

		// Proxy config is a singleton and cannot be deleted, only disabled
		// Verify it's disabled
		client := proxy.NewConfigClient(connector)
		proxyConfig, err := client.Get()
		if err != nil {
			// If we can't read it, consider it destroyed
			return nil
		}

		if proxyConfig.Enabled != nil && *proxyConfig.Enabled {
			return fmt.Errorf("Policy Proxy Config %s is still enabled", displayName)
		}
	}
	return nil
}

func testAccNsxtProxyConfigDisabled(name string, withContext bool) string {
	context := ""
	// Note: display_name and description are not included as NSX doesn't support them for proxy config singleton
	return fmt.Sprintf(`
resource "nsxt_proxy_config" "test" {
%s
  enabled      = false
  host         = "proxy.example.com"
}`, context)
}

func testAccNsxtProxyConfigDisabledWithDetails(name string, server string, port string, username string, withContext bool) string {
	context := ""
	// Note: display_name and description are not included as NSX doesn't support them for proxy config singleton
	// Note: enabled is false to avoid NSX connectivity validation with fake proxy servers
	return fmt.Sprintf(`
resource "nsxt_proxy_config" "test" {
%s
  enabled              = false
  scheme               = "HTTP"
  host                 = "%s"
  port                 = %s
  username             = "%s"
  password             = "password123"
}`, context, server, port, username)
}
