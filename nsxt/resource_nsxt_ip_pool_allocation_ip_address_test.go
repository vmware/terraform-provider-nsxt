/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
	"testing"
	"time"
)

const (
	// how long to wait for updated allocation IP address list (to check destroyment)
	waitSeconds = 150
)

func TestAccResourceNsxtIPPoolAllocationIPAddress_basic(t *testing.T) {
	poolName := getIPPoolName()
	if poolName == "" {
		t.Skipf("No NSXT_TEST_IP_POOL set - skipping test")
	}

	poolResourceName := "data.nsxt_ip_pool.acceptance_test"
	resourceName := "nsxt_ip_pool_allocation_ip_address.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXIPPoolAllocationIPAddressCheckDestroy(state, poolResourceName, resourceName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIPPoolAllocationIPAddressCreateTemplate(poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXIPPoolAllocationIPAddressExists(poolResourceName, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_pool_id"),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_id"),
				),
			},
		},
	})
}

func testAccNSXIPPoolAllocationIPAddressExists(poolResourceName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		exists, err := checkAllocationIPAddressExists(state, poolResourceName, resourceName)
		if err != nil {
			return err
		}
		if !*exists {
			return fmt.Errorf("Allocation %s in IP Pool %s is not existing", resourceName, poolResourceName)
		}
		return nil
	}
}

func getIPPoolIDByResourceName(state *terraform.State, poolResourceName string) (string, error) {
	rsPool, ok := state.RootModule().Resources[poolResourceName]
	if !ok {
		return "", fmt.Errorf("IP Pool resource %s not found in resources", getIPPoolName())
	}

	poolID := rsPool.Primary.ID
	if poolID == "" {
		return "", fmt.Errorf("IP Pool resource ID not set in resources ")
	}
	return poolID, nil
}

func checkAllocationIPAddressExists(state *terraform.State, poolResourceName string, resourceName string) (*bool, error) {
	poolID, err := getIPPoolIDByResourceName(state, poolResourceName)
	if err != nil {
		return nil, err
	}

	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("IP Pool allocation_ip_address resource %s not found in resources", resourceName)
	}

	nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
	if nsxClient == nil {
		return nil, resourceNotSupportedError()
	}
	listResult, responseCode, err := nsxClient.PoolManagementApi.ListIpPoolAllocations(nsxClient.Context, poolID)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving allocations for IP Pool ID %s. Error: %v", poolID, err)
	}

	if responseCode.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error while checking if allocations in IP Pool %s exists. HTTP return code was %d", poolID, responseCode.StatusCode)
	}

	exists := false
	for _, allocationIPAddress := range listResult.Results {
		if allocationIPAddress.AllocationId == rs.Primary.ID {
			exists = true
			break
		}
	}

	return &exists, nil
}

func testAccNSXIPPoolAllocationIPAddressCheckDestroy(state *terraform.State, poolResourceName, resourceName string) error {
	// PoolManagementApi.ListIpPoolAllocations() call does not return updated list within two minutes. Need to wait...
	fmt.Printf("testAccNSXIPPoolAllocationIPAddressCheckDestroy: waiting up to %d seconds\n", waitSeconds)

	timeout := time.Now().Add(waitSeconds * time.Second)
	for time.Now().Before(timeout) {
		exists, err := checkAllocationIPAddressExists(state, poolResourceName, resourceName)
		if err != nil {
			return err
		}
		if !*exists {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Timeout on check destroy IP address allocation")
}

func testAccNSXIPPoolAllocationIPAddressCreateTemplate(poolName string) string {
	return fmt.Sprintf(`
data "nsxt_ip_pool" "acceptance_test" {
  display_name = "%s"
}

resource "nsxt_ip_pool_allocation_ip_address" "test" {
  ip_pool_id = "${data.nsxt_ip_pool.acceptance_test.id}"
}`, poolName)
}
