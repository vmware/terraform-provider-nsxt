/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/ip_pools"
	"testing"
)

var nsxtPolicyStaticSubnetPoolID = "tfpool1"

func TestAccResourceNsxtPolicyIPPoolStaticSubnet_minimal(t *testing.T) {
	name := "staticsubnet1"
	testResourceName := "nsxt_policy_ip_pool_static_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolStaticSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolStaticSubnetCreateMinimalTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolStaticSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_range.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "cidr", "12.12.12.0/24"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "gateway", ""),
					resource.TestCheckResourceAttr(testResourceName, "dns_suffix", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPoolStaticSubnet_basic(t *testing.T) {
	name := "tfpool1"
	updatedName := fmt.Sprintf("%s-updated", name)
	testResourceName := "nsxt_policy_ip_pool_static_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolStaticSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolStaticSubnetCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolStaticSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_range.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "cidr", "12.12.12.0/24"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "gateway", ""),
					resource.TestCheckResourceAttr(testResourceName, "dns_suffix", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolStaticSubnet3AllocationsTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolStaticSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_range.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "cidr", "12.12.12.0/24"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "gateway", "12.12.12.1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_suffix", "tf.test"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolStaticSubnet2AllocationsTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolStaticSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_range.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "dns_nameservers.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "cidr", "12.12.12.0/24"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "gateway", "12.12.12.1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_suffix", "tf.test"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPoolStaticSubnet_import_basic(t *testing.T) {
	name := "tfpool1"
	testResourceName := "nsxt_policy_ip_pool_static_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolStaticSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolStaticSubnetCreateTemplate(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyIPPoolStaticSubnetImporterGetID,
			},
		},
	})
}

func testAccNSXPolicyIPPoolStaticSubnetImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["nsxt_policy_ip_pool_static_subnet.test"]
	if !ok {
		return "", fmt.Errorf("NSX Policy Static Subnet resource %s not found in resources", "nsxt_policy_ip_pool_static_subnet.test")
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy Static Subnet resource ID not set in resources ")
	}
	poolPath := rs.Primary.Attributes["pool_path"]
	if poolPath == "" {
		return "", fmt.Errorf("NSX Policy Static Subnet pool_path not set in resources ")
	}
	poolID := getPolicyIDFromPath(poolPath)
	return fmt.Sprintf("%s/%s", poolID, resourceID), nil
}

func testAccNSXPolicyIPPoolStaticSubnetCheckExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		client := ip_pools.NewDefaultIpSubnetsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX Policy Static Subnet resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX Policy Static Subnet resource ID not set in resources")
		}

		_, err := client.Get(nsxtPolicyStaticSubnetPoolID, resourceID)
		if err != nil {
			return fmt.Errorf("Failed to find Static Subnet %s", resourceID)
		}

		return nil
	}
}

func testAccNSXPolicyIPPoolStaticSubnetCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	client := ip_pools.NewDefaultIpSubnetsClient(connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_ip_pool_static_subnet" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := client.Get(nsxtPolicyStaticSubnetPoolID, resourceID)
		if err == nil {
			return fmt.Errorf("Static Subnet still exists %s", resourceID)
		}

	}
	return nil
}

func testAccNSXPolicyIPPoolStaticSubnetCreateMinimalTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name = "%s"
  pool_path    = "${nsxt_policy_ip_pool.pool1.path}"
  cidr         = "12.12.12.0/24"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}`, nsxtPolicyStaticSubnetPoolID, nsxtPolicyStaticSubnetPoolID, name)
}

func testAccNSXPolicyIPPoolStaticSubnetCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  pool_path    = "${nsxt_policy_ip_pool.pool1.path}"
  cidr         = "12.12.12.0/24"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, nsxtPolicyStaticSubnetPoolID, nsxtPolicyStaticSubnetPoolID, name)
}

func testAccNSXPolicyIPPoolStaticSubnet3AllocationsTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  pool_path       = "${nsxt_policy_ip_pool.pool1.path}"
  cidr            = "12.12.12.0/24"
  dns_nameservers = ["12.12.12.3", "12.12.12.4"]
  dns_suffix      = "tf.test"
  gateway         = "12.12.12.1"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
  allocation_range {
    start = "12.12.12.30"
    end   = "12.12.12.40"
  }
  allocation_range {
    start = "12.12.12.50"
    end   = "12.12.12.60"
  }
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, nsxtPolicyStaticSubnetPoolID, nsxtPolicyStaticSubnetPoolID, name)
}

func testAccNSXPolicyIPPoolStaticSubnet2AllocationsTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name    = "%s"
  description     = "Acceptance Test"
  pool_path       = "${nsxt_policy_ip_pool.pool1.path}"
  cidr            = "12.12.12.0/24"
  dns_nameservers = ["12.12.12.3"]
  dns_suffix      = "tf.test"
  gateway         = "12.12.12.1"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
  allocation_range {
    start = "12.12.12.50"
    end   = "12.12.12.60"
  }
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, nsxtPolicyStaticSubnetPoolID, nsxtPolicyStaticSubnetPoolID, name)
}
