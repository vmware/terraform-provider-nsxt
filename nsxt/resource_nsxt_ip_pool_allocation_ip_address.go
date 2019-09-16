/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

func resourceNsxtIPPoolAllocationIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIPPoolAllocationIPAddressCreate,
		Read:   resourceNsxtIPPoolAllocationIPAddressRead,
		Update: resourceNsxtIPPoolAllocationIPAddressUpdate,
		Delete: resourceNsxtIPPoolAllocationIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_pool_id": {
				Type:        schema.TypeString,
				Description: "ID of IP pool that allocation belongs to",
				Required:    true,
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Description: "IP Address that is allocated from the pool",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceNsxtIPPoolAllocationIPAddressCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	poolId := d.Get("ip_pool_id").(string)
	allocationId := d.Get("allocation_id").(string)
	allocationIpAddress := manager.AllocationIpAddress{
		AllocationId: allocationId,
	}

	allocationIpAddress, resp, err := nsxClient.PoolManagementApi.AllocateOrReleaseFromIpPool(nsxClient.Context, poolId, allocationIpAddress, "ALLOCATE")

	if err != nil {
		return fmt.Errorf("Error during IpPoolAllocationIpAddress create: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during IpPoolAllocationIpAddress create: %v", resp.StatusCode)
	}
	d.SetId(allocationIpAddress.AllocationId)

	return resourceNsxtIPPoolAllocationIPAddressRead(d, m)
}

func resourceNsxtIPPoolAllocationIPAddressRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	poolId := d.Get("ip_pool_id").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining pool id")
	}

	resultList, resp, err := nsxClient.PoolManagementApi.ListIpPoolAllocations(nsxClient.Context, poolId)
	if err != nil {
		return fmt.Errorf("Error during IPPoolAllocationIPAddress read: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during IpPoolAllocationIpAddress read: %v", resp.StatusCode)
	}

	for _, address := range resultList.Results {
		if address.AllocationId == id {
			d.Set("ip_pool_id", poolId)
			d.Set("allocation_id", address.AllocationId)
			return nil
		}
	}

	log.Printf("[DEBUG] IPPoolAllocationIPAddress list%s not found", id)
	d.SetId("")
	return nil
}

func resourceNsxtIPPoolAllocationIPAddressUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Updating IPPoolAllocationIPAddress is not supported")
}

func resourceNsxtIPPoolAllocationIPAddressDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	poolId := d.Get("ip_pool_id").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining pool id")
	}

	allocationIpAddress := manager.AllocationIpAddress{
		AllocationId: d.Id(),
	}
	_, resp, err := nsxClient.PoolManagementApi.AllocateOrReleaseFromIpPool(nsxClient.Context, poolId, allocationIpAddress, "RELEASE")
	if resp == nil && err != nil {
		return fmt.Errorf("Error during IPPoolAllocationIPAddress delete: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error during IPPoolAllocationIPAddress delete: status=%s", resp.Status)
	}
	return nil
}
