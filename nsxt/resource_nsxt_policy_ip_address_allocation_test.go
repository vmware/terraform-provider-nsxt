// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	ippools "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
)

var accTestPolicyIPAddressAllocationCreateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform created",
	"allocation_ip": "12.12.12.11",
}

var accTestPolicyIPAddressAllocationUpdateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform updated",
	"allocation_ip": "12.12.12.12",
}

var accTestPolicyIPAddressAllocationPoolName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationSubnetName = getAccTestResourceName()

var accTestPolicyIPAddressAllocationExhaustedPoolName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhaustedSubnetName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhaustedConsumerName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhaustedName = getAccTestResourceName()

func TestAccResourceNsxtPolicyIPAddressAllocation_basic(t *testing.T) {
	testAccResourceNsxtPolicyIPAddressAllocationBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyIPAddressAllocationBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyIPAddressAllocationBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(true, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", accTestPolicyIPAddressAllocationCreateAttributes["allocation_ip"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(withContext),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(false, withContext),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", accTestPolicyIPAddressAllocationUpdateAttributes["allocation_ip"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(withContext),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_anyIPBasic(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ip"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(false),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ip"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(false),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_importBasic(t *testing.T) {
	name := accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(true, false),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyIPAddressAllocationImporterGetID,
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(false),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_importBasic_multitenancy(t *testing.T) {
	name := accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(true, true),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(true),
			},
		},
	})
}

// TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted verifies that when an IP
// pool has no free addresses left, a create with allocation_ip unset passes the
// initial PATCH but then fails once NSX can't asynchronously realize an actual IP,
// and that the failed allocation is not left behind as an orphaned object on NSX
// Manager. This exercises the cleanup path fixed in
// resourceNsxtPolicyIPAddressAllocationCreate for realization failures.
func TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.consumer"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationExhaustedConsumerName); err != nil {
				return err
			}
			return testAccNsxtPolicyIPAddressAllocationCheckNotLeaked(state)
		},
		Steps: []resource.TestStep{
			{
				// Consume the single IP available in the pool.
				Config: testAccNsxtPolicyIPAddressAllocationExhaustedTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", "13.13.13.10"),
				),
			},
			{
				// The pool now has no free IPs left; this allocation must fail, and
				// must not leave an orphaned allocation behind on NSX Manager.
				Config:      testAccNsxtPolicyIPAddressAllocationExhaustedTemplate(true),
				ExpectError: regexp.MustCompile("Failed to get realized IP for path"),
			},
		},
	})
}

// testAccNsxtPolicyIPAddressAllocationCheckNotLeaked verifies that the only
// IPAddressAllocation left on NSX Manager under the pool is the "consumer" one. Since
// the failed "exhausted" allocation gets an auto-generated nsx_id (unknown to the
// test), it can't be looked up directly; instead this lists everything under the
// pool and confirms nothing besides the consumer's ID is present, to confirm that an
// allocation whose create failed (e.g. due to pool exhaustion) was properly cleaned
// up rather than orphaned.
func testAccNsxtPolicyIPAddressAllocationCheckNotLeaked(state *terraform.State) error {
	rs, ok := state.RootModule().Resources["nsxt_policy_ip_address_allocation.consumer"]
	if !ok {
		return fmt.Errorf("Policy IPAddressAllocation resource %s not found in resources", "nsxt_policy_ip_address_allocation.consumer")
	}
	poolPath := rs.Primary.Attributes["pool_path"]
	if poolPath == "" {
		return fmt.Errorf("No pool_path found for IP Address Allocation with ID %s", rs.Primary.ID)
	}
	poolID := getPolicyIDFromPath(poolPath)

	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := ippools.NewIpAllocationsClient(testAccGetSessionContext(), connector)
	if nsxClient == nil {
		return policyResourceNotSupportedError()
	}

	result, err := nsxClient.List(poolID, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, alloc := range result.Results {
		if alloc.Id != nil && *alloc.Id != rs.Primary.ID {
			return fmt.Errorf("Policy IPAddressAllocation %s was left behind on NSX Manager after a failed create caused by pool exhaustion", *alloc.Id)
		}
	}
	return nil
}

func testAccNSXPolicyIPAddressAllocationImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["nsxt_policy_ip_address_allocation.test"]
	if !ok {
		return "", fmt.Errorf("NSX Policy IP Allocation resource %s not found in resources", "nsxt_policy_ip_address_allocation.test")
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy IP Allocation resource ID not set in resources ")
	}
	poolPath := rs.Primary.Attributes["pool_path"]
	if poolPath == "" {
		return "", fmt.Errorf("NSX Policy IP Allocation pool_path not set in resources ")
	}
	poolID := getPolicyIDFromPath(poolPath)
	return fmt.Sprintf("%s/%s", poolID, resourceID), nil
}

func testAccNsxtPolicyIPAddressAllocationExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := ippools.NewIpAllocationsClient(testAccGetSessionContext(), connector)
		if nsxClient == nil {
			return policyResourceNotSupportedError()
		}

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPAddressAllocation resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IPAddressAllocation resource ID not set in resources")
		}

		poolPath := rs.Primary.Attributes["pool_path"]
		if poolPath == "" {
			return fmt.Errorf("No pool_path found for IP Address Allocation with ID %s", resourceID)
		}
		poolID := getPolicyIDFromPath(poolPath)

		_, err := nsxClient.Get(poolID, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy IPAddressAllocation ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyIPAddressAllocationCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := ippools.NewIpAllocationsClient(testAccGetSessionContext(), connector)
	if nsxClient == nil {
		return policyResourceNotSupportedError()
	}
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ip_address_allocation" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		poolPath := rs.Primary.Attributes["pool_path"]
		if poolPath == "" {
			return fmt.Errorf("No pool_path found for IP Address Allocation with ID %s", resourceID)
		}
		poolID := getPolicyIDFromPath(poolPath)

		_, err := nsxClient.Get(poolID, resourceID)
		if err == nil {
			return fmt.Errorf("Policy IPAddressAllocation %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIPAddressAllocationTemplate(createFlow, withContext bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPAddressAllocationCreateAttributes
	} else {
		attrMap = accTestPolicyIPAddressAllocationUpdateAttributes
	}
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(withContext) + fmt.Sprintf(`
resource "nsxt_policy_ip_address_allocation" "test" {
%s
  display_name  = "%s"
  description   = "%s"
  allocation_ip = "%s"
  pool_path     = nsxt_policy_ip_pool.test.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, context, attrMap["display_name"], attrMap["description"], attrMap["allocation_ip"])
}

func testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPAddressAllocationCreateAttributes
	} else {
		attrMap = accTestPolicyIPAddressAllocationUpdateAttributes
	}
	return testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(false) + fmt.Sprintf(`

resource "nsxt_policy_ip_address_allocation" "test" {
  display_name = "%s"
  description  = "%s"
  pool_path    = nsxt_policy_ip_pool.test.path
  depends_on   = [data.nsxt_policy_realization_info.subnet_realization]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"])
}

func testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "test" {
%s
  display_name = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
%s
  display_name = "%s"
  pool_path    = nsxt_policy_ip_pool.test.path
  cidr         = "12.12.12.0/24"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}

data "nsxt_policy_realization_info" "subnet_realization" {
  path = nsxt_policy_ip_pool_static_subnet.test.path
}`, context, accTestPolicyIPAddressAllocationPoolName, context, accTestPolicyIPAddressAllocationSubnetName)
}

// testAccNsxtPolicyIPAddressAllocationExhaustedTemplate builds a pool with a single
// available IP. The "consumer" resource always consumes that IP. When
// includeExhausting is true, a second allocation with an auto-generated nsx_id is
// added; since the pool is already exhausted at that point, its create is expected
// to fail. Leaving nsx_id auto-generated (rather than fixed) matches the common case
// and is an exact repro of the reported customer issue.
func testAccNsxtPolicyIPAddressAllocationExhaustedTemplate(includeExhausting bool) string {
	config := fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "exhausted" {
  display_name = "%s"
}

resource "nsxt_policy_ip_pool_static_subnet" "exhausted" {
  display_name = "%s"
  pool_path    = nsxt_policy_ip_pool.exhausted.path
  cidr         = "13.13.13.0/24"
  allocation_range {
    start = "13.13.13.10"
    end   = "13.13.13.10"
  }
}

data "nsxt_policy_realization_info" "exhausted_subnet_realization" {
  path = nsxt_policy_ip_pool_static_subnet.exhausted.path
}

resource "nsxt_policy_ip_address_allocation" "consumer" {
  display_name = "%s"
  pool_path    = nsxt_policy_ip_pool.exhausted.path
  depends_on   = [data.nsxt_policy_realization_info.exhausted_subnet_realization]
}`, accTestPolicyIPAddressAllocationExhaustedPoolName, accTestPolicyIPAddressAllocationExhaustedSubnetName, accTestPolicyIPAddressAllocationExhaustedConsumerName)

	if !includeExhausting {
		return config
	}

	return config + fmt.Sprintf(`

resource "nsxt_policy_ip_address_allocation" "exhausted" {
  display_name = "%s"
  pool_path    = nsxt_policy_ip_pool.exhausted.path
  depends_on   = [nsxt_policy_ip_address_allocation.consumer]
}`, accTestPolicyIPAddressAllocationExhaustedName)
}
