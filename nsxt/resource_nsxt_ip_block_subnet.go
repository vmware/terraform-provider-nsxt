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

func resourceNsxtIpBlockSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIpBlockSubnetCreate,
		Read:   resourceNsxtIpBlockSubnetRead,
		// Update ip block subnet is not supported by the NSX
		Delete: resourceNsxtIpBlockSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				ForceNew:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"tag": getTagsSchemaForceNew(),
			"block_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Block id for which the subnet is created",
				Required:    true,
				ForceNew:    true,
			},
			"size": &schema.Schema{
				// TODO(asarfaty): add validation for power of 2
				Type:        schema.TypeInt,
				Description: "Represents the size or number of ip addresses in the subnet",
				Required:    true,
				ForceNew:    true,
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents network address and the prefix length which will be associated with a layer-2 broadcast domain",
				Computed:    true,
			},
			// TODO(asarfaty) add allocation ranges
		},
	}
}

func resourceNsxtIpBlockSubnetCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	blockId := d.Get("block_id").(string)
	size := int64(d.Get("size").(int))
	tags := getTagsFromSchema(d)
	ipBlockSubnet := manager.IpBlockSubnet{
		DisplayName: displayName,
		Description: description,
		BlockId:     blockId,
		Size:        size,
		Tags:        tags,
	}

	ipBlockSubnet, resp, err := nsxClient.PoolManagementApi.CreateIpBlockSubnet(nsxClient.Context, ipBlockSubnet)

	if err != nil {
		return fmt.Errorf("Error during IpBlockSubnet create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during IpBlockSubnet create: %v", resp.StatusCode)
	}
	d.SetId(ipBlockSubnet.Id)

	return resourceNsxtIpBlockSubnetRead(d, m)
}

func resourceNsxtIpBlockSubnetRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ipBlockSubnet, resp, err := nsxClient.PoolManagementApi.ReadIpBlockSubnet(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpBlockSubnet %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during IpBlockSubnet read: %v", err)
	}

	d.Set("display_name", ipBlockSubnet.DisplayName)
	d.Set("description", ipBlockSubnet.Description)
	d.Set("block_id", ipBlockSubnet.BlockId)
	d.Set("size", ipBlockSubnet.Size)
	setTagsInSchema(d, ipBlockSubnet.Tags)
	//d.Set("allocation_range", ipBlockSubnet.AllocationRanges)
	d.Set("cidr", ipBlockSubnet.Cidr)
	d.Set("revision", ipBlockSubnet.Revision)

	return nil
}

func resourceNsxtIpBlockSubnetDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.PoolManagementApi.DeleteIpBlockSubnet(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during IpBlockSubnet delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpBlockSubnet %s not found", id)
		d.SetId("")
	}
	return nil
}
