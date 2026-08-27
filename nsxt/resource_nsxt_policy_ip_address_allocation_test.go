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
var accTestPolicyIPAddressAllocationExhaustedID = getAccTestResourceName()

var accTestPolicyIPAddressAllocationExhausted92PoolName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhausted92SubnetName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhausted92ConsumerName = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhausted92Name = getAccTestResourceName()
var accTestPolicyIPAddressAllocationExhausted92ID = getAccTestResourceName()

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

// TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted verifies that, on NSX
// versions below 9.2.0, when an IP pool has no free addresses left, a create with
// allocation_ip unset passes the initial PATCH but then fails once NSX can't
// asynchronously realize an actual IP, and that the failed allocation is not left
// behind as an orphaned object on NSX Manager. This exercises the cleanup path fixed
// in resourceNsxtPolicyIPAddressAllocationCreate for realization failures.
//
// On NSX 9.2.0+ pool exhaustion is instead rejected synchronously by the PATCH call
// itself, so the create never reaches that cleanup path; see the _poolExhausted92
// sibling test below for that behavior.
func TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.consumer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersionLessThan(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationExhaustedConsumerName); err != nil {
				return err
			}
			return testAccNsxtPolicyIPAddressAllocationCheckNotLeaked(state, accTestPolicyIPAddressAllocationExhaustedID)
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

// TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted92 covers the same pool
// exhaustion scenario as TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted,
// but for NSX 9.2.0+, where the backend validates IP availability synchronously
// inside the initial PATCH call and rejects it outright (error code 520054) rather
// than deferring the failure to async realization. Nothing is created on NSX in this
// case, so this does not exercise the realization-failure cleanup path (covered by
// the pre-9.2 sibling test above); it exists so pool exhaustion still has acceptance
// coverage on 9.2+.
func TestAccResourceNsxtPolicyIPAddressAllocation_poolExhausted92(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.consumer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			if err := testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationExhausted92ConsumerName); err != nil {
				return err
			}
			return testAccNsxtPolicyIPAddressAllocationCheckNotLeaked(state, accTestPolicyIPAddressAllocationExhausted92ID)
		},
		Steps: []resource.TestStep{
			{
				// Consume the single IP available in the pool.
				Config: testAccNsxtPolicyIPAddressAllocationExhausted92Template(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", "13.13.13.10"),
				),
			},
			{
				// The pool now has no free IPs left; NSX rejects the PATCH outright.
				Config:      testAccNsxtPolicyIPAddressAllocationExhausted92Template(true),
				ExpectError: regexp.MustCompile("is exhausted, no free IPs available for allocation"),
			},
		},
	})
}

// testAccNsxtPolicyIPAddressAllocationCheckNotLeaked verifies that the IP address
// allocation identified by exhaustedID does not exist on NSX Manager. It is used to
// confirm that an allocation whose create failed (e.g. due to pool exhaustion) was
// properly cleaned up rather than orphaned.
func testAccNsxtPolicyIPAddressAllocationCheckNotLeaked(state *terraform.State, exhaustedID string) error {
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

	_, err := nsxClient.Get(poolID, exhaustedID)
	if err == nil {
		return fmt.Errorf("Policy IPAddressAllocation with ID %s was left behind on NSX Manager after a failed create caused by pool exhaustion", exhaustedID)
	}
	if !isNotFoundError(err) {
		return err
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
		if !isNotFoundError(err) {
			return err
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

// testAccNsxtPolicyIPAddressAllocationExhaustedTemplateNames builds a pool with a
// single available IP. The "consumer" resource always consumes that IP. When
// includeExhausting is true, a second allocation with a known, fixed nsx_id is added;
// since the pool is already exhausted at that point, its create is expected to fail.
func testAccNsxtPolicyIPAddressAllocationExhaustedTemplateNames(poolName, subnetName, consumerName, exhaustedID, exhaustedName string, includeExhausting bool) string {
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
}`, poolName, subnetName, consumerName)

	if !includeExhausting {
		return config
	}

	return config + fmt.Sprintf(`

resource "nsxt_policy_ip_address_allocation" "exhausted" {
  nsx_id       = "%s"
  display_name = "%s"
  pool_path    = nsxt_policy_ip_pool.exhausted.path
  depends_on   = [nsxt_policy_ip_address_allocation.consumer]
}`, exhaustedID, exhaustedName)
}

func testAccNsxtPolicyIPAddressAllocationExhaustedTemplate(includeExhausting bool) string {
	return testAccNsxtPolicyIPAddressAllocationExhaustedTemplateNames(
		accTestPolicyIPAddressAllocationExhaustedPoolName,
		accTestPolicyIPAddressAllocationExhaustedSubnetName,
		accTestPolicyIPAddressAllocationExhaustedConsumerName,
		accTestPolicyIPAddressAllocationExhaustedID,
		accTestPolicyIPAddressAllocationExhaustedName,
		includeExhausting,
	)
}

func testAccNsxtPolicyIPAddressAllocationExhausted92Template(includeExhausting bool) string {
	return testAccNsxtPolicyIPAddressAllocationExhaustedTemplateNames(
		accTestPolicyIPAddressAllocationExhausted92PoolName,
		accTestPolicyIPAddressAllocationExhausted92SubnetName,
		accTestPolicyIPAddressAllocationExhausted92ConsumerName,
		accTestPolicyIPAddressAllocationExhausted92ID,
		accTestPolicyIPAddressAllocationExhausted92Name,
		includeExhausting,
	)
}
