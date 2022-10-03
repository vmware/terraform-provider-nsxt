/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes = map[string]string{
	"display_name":               getAccTestResourceName(),
	"description":                "Terraform-provisioned IPsec Route-Based VPN",
	"enabled":                    "true",
	"vpn_type":                   "RouteBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"ip_addresses":               "169.254.152.25",
	"prefix_length":              "24",
	"peer_address":               "18.18.18.22",
	"peer_id":                    "18.18.18.22",
	"psk":                        "VMware123!",
	"connection_initiation_mode": "INITIATOR",
}

var accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes = map[string]string{
	"display_name":               getAccTestResourceName(),
	"description":                "Terraform-provisioned IPsec Route-Based VPN",
	"enabled":                    "true",
	"vpn_type":                   "RouteBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"ip_addresses":               "169.254.152.26",
	"prefix_length":              "24",
	"peer_address":               "18.18.18.21",
	"peer_id":                    "18.18.18.21",
	"psk":                        "VMware123!",
	"connection_initiation_mode": "INITIATOR",
}

var accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes = map[string]string{
	"display_name":               getAccTestResourceName(),
	"description":                "Terraform-provisioned IPsec Route-Based VPN",
	"enabled":                    "true",
	"vpn_type":                   "PolicyBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"peer_address":               "18.18.18.21",
	"peer_id":                    "18.18.18.21",
	"psk":                        "VMware123!",
	"connection_initiation_mode": "RESPOND_ONLY",
	"sources":                    "192.170.10.0/24",
	"destinations":               "192.171.10.0/24",
	"action":                     "PROTECT",
}

var accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes = map[string]string{
	"display_name":               getAccTestResourceName(),
	"description":                "Terraform-provisioned IPsec Route-Based VPN",
	"enabled":                    "true",
	"vpn_type":                   "PolicyBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"peer_address":               "18.18.18.22",
	"peer_id":                    "18.18.18.22",
	"psk":                        "VMware123!",
	"connection_initiation_mode": "RESPOND_ONLY",
	"sources":                    "192.172.10.0/24",
	"destinations":               "192.173.10.0/24",
	"action":                     "PROTECT",
}

func TestAccResourceNsxtPolicyIPSecVpnSessionRouteBased_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ipsec_vpn_session.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnSessionCheckDestroy(state, accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnSessionRouteBasedTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnSessionExists(accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix_length", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["prefix_length"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes["connection_initiation_mode"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnSessionRouteBasedTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnSessionExists(accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix_length", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["prefix_length"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes["connection_initiation_mode"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPSecVpnSessionPolicyBased_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ipsec_vpn_session.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnSessionCheckDestroy(state, accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnSessionPolicyBasedTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnSessionExists(accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["connection_initiation_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources.0", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations.0", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes["action"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnSessionPolicyBasedTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnSessionExists(accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "service_path"),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["connection_initiation_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources.0", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations.0", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes["action"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIPSecVpnSessionExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPSecVpnSession resource %s not found in resources", displayName)
		}

		resourceID := rs.Primary.Attributes["id"]
		servicePath := rs.Primary.Attributes["service_path"]
		isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
		if err != nil {
			return err
		}
		if !isT0 {
			return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
		}
		if resourceID == "" {
			return fmt.Errorf("Policy IPSecVpnSession resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIPSecVpnSessionExists(gwID, localeServiceID, serviceID, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy IPSecVpnSession %s does not exist", displayName)
		}

		return nil
	}
}

func testAccNsxtPolicyIPSecVpnSessionCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ipsec_vpn_session" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		servicePath := rs.Primary.Attributes["service_path"]
		isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
		if err != nil {
			return err
		}
		if !isT0 {
			return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
		}

		exists, err := resourceNsxtPolicyIPSecVpnSessionExists(gwID, localeServiceID, serviceID, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy IPSecVpnSession %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIPSecVpnSessionPreConditionTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_ike_profile" "test" {
	display_name          = "%s"
	description           = "Ike profile for ipsec vpn session"
	encryption_algorithms = ["AES_128"]
	digest_algorithms     = ["SHA2_256"]
	dh_groups             = ["GROUP14"]
	ike_version           = "IKE_V2"
}

resource "nsxt_policy_ipsec_vpn_tunnel_profile" "test" {
	display_name          = "%s"
	description           = "Terraform provisioned IPSec VPN Ike Profile"
	df_policy             = "COPY"
	encryption_algorithms = ["AES_128"]
	digest_algorithms     = ["SHA2_256"]
	dh_groups             = ["GROUP14"]
	sa_life_time          = 7200
}

resource "nsxt_policy_ipsec_vpn_dpd_profile" "test" {
	display_name       = "%s"
	description        = "Terraform provisioned IPSec VPN DPD Profile"
	dpd_probe_mode     = "ON_DEMAND"
	dpd_probe_interval = 1
	enabled            = true
	retry_count        = 8
  }

resource "nsxt_policy_ipsec_vpn_service" "test" {
	display_name          = "%s"
	locale_service_path   = one(nsxt_policy_tier0_gateway.test.locale_service).path
  }

resource "nsxt_policy_ipsec_vpn_local_endpoint" "test" {
	service_path  = nsxt_policy_ipsec_vpn_service.test.path
	display_name  = "%s"
	local_address = "20.20.0.25"
  }
`, getAccTestResourceName(), getAccTestResourceName(), getAccTestResourceName(), getAccTestResourceName(), getAccTestResourceName())
}

func testAccNsxtPolicyIPSecVpnSessionRouteBasedTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPSecVpnSessionRouteBasedCreateAttributes
	} else {
		attrMap = accTestPolicyIPSecVpnSessionRouteBasedUpdateAttributes
	}
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) +
		testAccNsxtPolicyTier0WithEdgeClusterForVPN("test") +
		testAccNsxtPolicyIPSecVpnSessionPreConditionTemplate() +
		fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_session" "test" {
	display_name               = "%s"
	description                = "%s"
	ike_profile_path           = nsxt_policy_ipsec_vpn_ike_profile.test.path
	tunnel_profile_path        = nsxt_policy_ipsec_vpn_tunnel_profile.test.path
	dpd_profile_path		   = nsxt_policy_ipsec_vpn_dpd_profile.test.path
	local_endpoint_path		   = nsxt_policy_ipsec_vpn_local_endpoint.test.path
	enabled                    = "%s"
	service_path               = nsxt_policy_ipsec_vpn_service.test.path
	vpn_type                   = "%s"
	authentication_mode        = "%s"
	compliance_suite           = "%s"
	ip_addresses               = ["%s"]
	prefix_length              = "%s"
	peer_address               = "%s"
	peer_id                    = "%s"
	psk                        = "%s"
	connection_initiation_mode = "%s"

	tag {
	scope = "scope1"
	tag   = "tag1"
	  }
}`, attrMap["display_name"], attrMap["description"], attrMap["enabled"], attrMap["vpn_type"],
			attrMap["authentication_mode"], attrMap["compliance_suite"], attrMap["ip_addresses"], attrMap["prefix_length"], attrMap["peer_address"], attrMap["peer_id"],
			attrMap["psk"], attrMap["connection_initiation_mode"])
}

func testAccNsxtPolicyIPSecVpnSessionPolicyBasedTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPSecVpnSessionPolicyBasedCreateAttributes
	} else {
		attrMap = accTestPolicyIPSecVpnSessionPolicyBasedUpdateAttributes
	}
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) +
		testAccNsxtPolicyTier0WithEdgeClusterForVPN("test") +
		testAccNsxtPolicyIPSecVpnSessionPreConditionTemplate() +
		fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_session" "test" {
	display_name               = "%s"
	description                = "%s"
	ike_profile_path           = nsxt_policy_ipsec_vpn_ike_profile.test.path
	tunnel_profile_path        = nsxt_policy_ipsec_vpn_tunnel_profile.test.path
	dpd_profile_path		   = nsxt_policy_ipsec_vpn_dpd_profile.test.path
	local_endpoint_path		   = nsxt_policy_ipsec_vpn_local_endpoint.test.path
	enabled                    = "%s"
	service_path               = nsxt_policy_ipsec_vpn_service.test.path
	vpn_type                   = "%s"
	authentication_mode        = "%s"
	compliance_suite           = "%s"
	peer_address               = "%s"
	peer_id                    = "%s"
	psk                        = "%s"
	connection_initiation_mode = "%s"

	rule {
		sources             = ["%s"]
		destinations        = ["%s"]
		action              = "%s"
	  }

	tag {
	scope = "scope1"
	tag   = "tag1"
	  }
}`, attrMap["display_name"], attrMap["description"], attrMap["enabled"], attrMap["vpn_type"],
			attrMap["authentication_mode"], attrMap["compliance_suite"], attrMap["peer_address"], attrMap["peer_id"],
			attrMap["psk"], attrMap["connection_initiation_mode"], attrMap["sources"], attrMap["destinations"], attrMap["action"])
}
