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

var accTestPolicyIpBlockQuotaCreateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform created",
	"single_ip_cidrs":       "2",
	"ip_block_visibility":   "PRIVATE",
	"ip_block_address_type": "IPV4",
	"mask":                  "/28",
	"total_count":           "2",
}

var accTestPolicyIpBlockQuotaUpdateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform updated",
	"single_ip_cidrs":       "5",
	"ip_block_visibility":   "PRIVATE",
	"ip_block_address_type": "IPV6",
	"mask":                  "/104",
	"total_count":           "5",
}

func TestAccResourceNsxtPolicyIpBlockQuota_basic(t *testing.T) {
	testAccResourceNsxtPolicyIpBlockQuota_basic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "9.0.0")
	})
}

func TestAccResourceNsxtPolicyIpBlockQuota_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIpBlockQuota_basic(t, true, func() {
		testAccPreCheck(t)
		testAccNSXVersion(t, "9.0.0")
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIpBlockQuota_basic(t *testing.T, multitenancy bool, preCheck func()) {
	testResourceName := "nsxt_policy_ip_block_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIpBlockQuotaCheckDestroy(state, accTestPolicyIpBlockQuotaUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIpBlockQuotaTemplate(true, multitenancy),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIpBlockQuotaExists(accTestPolicyIpBlockQuotaCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIpBlockQuotaCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIpBlockQuotaCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.single_ip_cidrs", accTestPolicyIpBlockQuotaCreateAttributes["single_ip_cidrs"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "quota.0.ip_block_paths.0"),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.ip_block_visibility", accTestPolicyIpBlockQuotaCreateAttributes["ip_block_visibility"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.ip_block_address_type", accTestPolicyIpBlockQuotaCreateAttributes["ip_block_address_type"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.0.mask", accTestPolicyIpBlockQuotaCreateAttributes["mask"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.0.total_count", accTestPolicyIpBlockQuotaCreateAttributes["total_count"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIpBlockQuotaTemplate(false, multitenancy),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIpBlockQuotaExists(accTestPolicyIpBlockQuotaUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIpBlockQuotaUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIpBlockQuotaUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.single_ip_cidrs", accTestPolicyIpBlockQuotaUpdateAttributes["single_ip_cidrs"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "quota.0.ip_block_paths.0"),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.ip_block_visibility", accTestPolicyIpBlockQuotaUpdateAttributes["ip_block_visibility"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.ip_block_address_type", accTestPolicyIpBlockQuotaUpdateAttributes["ip_block_address_type"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.0.mask", accTestPolicyIpBlockQuotaUpdateAttributes["mask"]),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.0.total_count", accTestPolicyIpBlockQuotaUpdateAttributes["total_count"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIpBlockQuotaMinimalistic(multitenancy),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIpBlockQuotaExists(accTestPolicyIpBlockQuotaCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "quota.0.other_cidrs.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIpBlockQuota_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_ip_block_quota.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIpBlockQuotaCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIpBlockQuotaMinimalistic(false),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyIpBlockQuotaExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IpBlockQuota resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IpBlockQuota resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIpBlockQuotaExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy IpBlockQuota %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyIpBlockQuotaCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ip_block_quota" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyIpBlockQuotaExists(testAccGetSessionProjectContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy IpBlockQuota %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIpBlockQuotaTemplate(createFlow bool, withContext bool) string {
	var attrMap map[string]string
	var ipblock string
	if createFlow {
		attrMap = accTestPolicyIpBlockQuotaCreateAttributes
		ipblock = "test-ipv4"
	} else {
		attrMap = accTestPolicyIpBlockQuotaUpdateAttributes
		ipblock = "test-ipv6"
	}
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ip_block" "test-ipv4" {
  %s
  visibility   = "%s"
  display_name = "%s-ipv4"
  cidr         = "20.20.20.0/24"
}

resource "nsxt_policy_ip_block" "test-ipv6" {
  %s
  visibility   = "%s"
  display_name = "%s-ipv6"
  cidr         = "fc7e:0011:da43::0/48"
}

resource "nsxt_policy_ip_block_quota" "test" {
  %s
  display_name = "%s"
  description  = "%s"

  quota {
    ip_block_paths        = [nsxt_policy_ip_block.%s.path]
    ip_block_visibility   = "%s"
    ip_block_address_type = "%s"

    single_ip_cidrs = %s
    other_cidrs {
      mask        = "%s"
      total_count = %s
    }
  }

}`, context, attrMap["ip_block_visibility"], attrMap["display_name"], context, attrMap["ip_block_visibility"], attrMap["display_name"], context, attrMap["display_name"], attrMap["description"], ipblock, attrMap["ip_block_visibility"], attrMap["ip_block_address_type"], attrMap["single_ip_cidrs"], attrMap["mask"], attrMap["total_count"])
}

func testAccNsxtPolicyIpBlockQuotaMinimalistic(withContext bool) string {
	attrMap := accTestPolicyIpBlockQuotaCreateAttributes
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ip_block_quota" "test" {
  %s
  display_name = "%s"
  quota {
    ip_block_visibility   = "%s"
    ip_block_address_type = "%s"
    single_ip_cidrs       = %s

    other_cidrs {
    }
  }
}`, context, attrMap["display_name"], attrMap["ip_block_visibility"], attrMap["ip_block_address_type"], attrMap["single_ip_cidrs"])
}
