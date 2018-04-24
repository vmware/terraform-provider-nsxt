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

// formatLogicalRouterRollbackError defines the verbose error when
// rollback fails on a logical router create operation.
const formatIpBlockRollbackError = `
WARNING:
There was an error during the creation of IP block router %s:
%s
Additionally, there was an error deleting the IP block during rollback:
%s
The IP block may still exist in the NSX. If it does, please manually delete it
and try again.
`

func resourceNsxtIpBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIpBlockCreate,
		Read:   resourceNsxtIpBlockRead,
		Update: resourceNsxtIpBlockUpdate,
		Delete: resourceNsxtIpBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents network address and the prefix length which will be associated with a layer-2 broadcast domain",
				Required:    true,
				ForceNew:    true,
			},
			"subnet": &schema.Schema{
				Type:        schema.TypeList,
				Description: "List of subnets",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The ID of this subnet",
							Computed:    true,
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Description: "Description of this resource",
							Optional:    true,
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The display name of this resource. Defaults to ID if not set",
							Optional:    true,
							Computed:    true,
						},
						"size": &schema.Schema{
							// TODO(asarfaty): add validation for power of 2
							Type:        	schema.TypeInt,
							Description: 	"Represents the size or number of ip addresses in the subnet",
							Required:    	true,
						},
					},
				},
			},
		},
	}
}

func getIpBlockSubnetsFromSchema(d *schema.ResourceData, ip_block_id string) []manager.IpBlockSubnet {
	subnets := d.Get("subnet").([]interface{})
	var subnetsList []manager.IpBlockSubnet
	for _, subnet := range subnets {
		data := subnet.(map[string]interface{})
		elem := manager.IpBlockSubnet{
			BlockId:     ip_block_id,
			Size:        int64(data["size"].(int)),
			DisplayName: data["display_name"].(string),
			Description: data["description"].(string),
		}
		subnetsList = append(subnetsList, elem)
	}
	return subnetsList
}

func resourceNsxtIpBlockReadSubnets(nsxClient *api.APIClient, ip_block_id string) (manager.IpBlockSubnetListResult, error) {
	// Get all block subnets
	localVarOptionals := make(map[string]interface{})
	localVarOptionals["blockId"] = ip_block_id
	subnets, _, err := nsxClient.PoolManagementApi.ListIpBlockSubnets(nsxClient.Context, localVarOptionals)
	return subnets, err
}

// Delete the IP block together with all the backend subnets
func resourceNsxtIpBlockForceDelete(nsxClient *api.APIClient, ip_block_id string) error {
	// Get all block subnets
	subnets, err1 := resourceNsxtIpBlockReadSubnets(nsxClient, ip_block_id)
	if err1 != nil {
		return err1
	}

	// Delete all block subnets
	for _, objInList := range subnets.Results {
		_, err2 := nsxClient.PoolManagementApi.DeleteIpBlockSubnet(nsxClient.Context, objInList.Id)
		if err2 != nil {
			return err2
		}
	}

	// Delete the block itself
	_, err3 := nsxClient.PoolManagementApi.DeleteIpBlock(nsxClient.Context, ip_block_id)
	if err3 != nil {
		return err3
	}
	return nil
}

func resourceNsxtIpBlockCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	cidr := d.Get("cidr").(string)
	ipBlock := manager.IpBlock{
		Description: description,
		DisplayName: displayName,
		Tags:        tags,
		Cidr:        cidr,
	}
	// Create the IP Block
	ipBlock, resp, err := nsxClient.PoolManagementApi.CreateIpBlock(nsxClient.Context, ipBlock)

	if err != nil {
		return fmt.Errorf("Error during IpBlock create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during IpBlock create: %v", resp.StatusCode)
	}
	d.SetId(ipBlock.Id)

	// Create the subnets
	subnets := getIpBlockSubnetsFromSchema(d, ipBlock.Id)
	if len(subnets) > 0 {
		for _, subnet := range subnets {
			_, resp, err1 := nsxClient.PoolManagementApi.CreateIpBlockSubnet(nsxClient.Context, subnet)
			if err1 != nil || resp.StatusCode != http.StatusCreated {
				// Delete the created ip block subnets & ip block
				err2 := resourceNsxtIpBlockForceDelete(nsxClient, ipBlock.Id)
				if err2 != nil {
					// rollback failed
					return fmt.Errorf(formatIpBlockRollbackError, ipBlock.Id, err1, err2)
				}
				return fmt.Errorf("Error during IpBlockSubnet create: %v", err1)
			}
		}
	}
	return resourceNsxtIpBlockRead(d, m)
}

func setIpBlockSubnetsInSchema(d *schema.ResourceData, subnets []manager.IpBlockSubnet) error {
	subnetsList := make([]map[string]interface{}, 0, len(subnets))
	for _, subnet := range subnets {
		elem := make(map[string]interface{})
		elem["id"] = subnet.Id
		elem["display_name"] = subnet.DisplayName
		elem["description"] = subnet.Description
		elem["size"] = subnet.Size
		subnetsList = append(subnetsList, elem)
	}
	err := d.Set("subnet", subnetsList)
	return err
}

func resourceNsxtIpBlockRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ipBlock, resp, err := nsxClient.PoolManagementApi.ReadIpBlock(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpBlock %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during IpBlock read: %v", err)
	}

	d.Set("revision", ipBlock.Revision)
	d.Set("description", ipBlock.Description)
	d.Set("display_name", ipBlock.DisplayName)
	setTagsInSchema(d, ipBlock.Tags)
	d.Set("cidr", ipBlock.Cidr)

	// Read all the ip block subnets
	subnets, err1 := resourceNsxtIpBlockReadSubnets(nsxClient, id)
	if err1 != nil {
		return fmt.Errorf("Error during IpBlockSubnets read: %v", err1)
	}
	err = setIpBlockSubnetsInSchema(d, subnets.Results)
	if err != nil {
		return fmt.Errorf("Error during IpBlockSubnets set in schema: %v", err)
	}

	return nil
}

func resourceNsxtIpBlockUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	// Check if any relevant field of the IP block changed (Ignore the cidr because it has ForceNew)
	if d.HasChange("description") || d.HasChange("display_name") || d.HasChange("tags") {
		displayName := d.Get("display_name").(string)
		description := d.Get("description").(string)
		cidr := d.Get("cidr").(string)
		tags := getTagsFromSchema(d)
		revision := int64(d.Get("revision").(int))
		ipBlock := manager.IpBlock{
			DisplayName: displayName,
			Description: description,
			Cidr:        cidr,
			Tags:        tags,
			Revision:    revision,
		}
		// Update the IP block
		ipBlock, resp, err := nsxClient.PoolManagementApi.UpdateIpBlock(nsxClient.Context, id, ipBlock)

		if err != nil || resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Error during IpBlock update: %v", err)
		}
	}

	// There is no update action for subnets. Only add/delete subnets
	if d.HasChange("subnet") {
		old_data, new_data := d.GetChange("subnet")
		old_subnets := old_data.([]interface{})
		new_subnets := new_data.([]interface{})
		// Go over the new subnets and check which should be created new and which was modified
		log.Printf("[DEBUG] DEBUG ADIT subnets changed?\n old %#v\n new %#v", old_subnets, new_subnets)
		for _, new_sub_data := range new_subnets {
			log.Printf("[DEBUG] DEBUG ADIT new_sub_data %#v", new_sub_data)
			new_subnet := new_sub_data.(map[string]interface{})
			found := false
			for _, old_sub_data := range old_subnets {
				old_subnet := old_sub_data.(map[string]interface{})
				if old_subnet["id"] == new_subnet["id"]{
					log.Printf("[DEBUG] DEBUG ADIT subnet changed?\n old %#v\n new %#v", old_subnet, new_subnet)
					found = true
					// Check if any of the attributes changed
					if old_subnet["display_name"].(string) != new_subnet["display_name"].(string) ||
						old_subnet["description"].(string) != new_subnet["description"].(string) ||
						old_subnet["size"].(int) != new_subnet["size"].(int) {
						return fmt.Errorf("Error during IpBlock update: Updating subnet %s values is not allowed", old_subnet["id"].(string))
					}
					break
				}
			}
			if !found {
				// Create a new subnet
				sub := manager.IpBlockSubnet{
					BlockId:     id,
					Size:        int64(new_subnet["size"].(int)),
					DisplayName: new_subnet["display_name"].(string),
					Description: new_subnet["description"].(string),
				}
				_, resp, err := nsxClient.PoolManagementApi.CreateIpBlockSubnet(nsxClient.Context, sub)
				if err != nil || resp.StatusCode == http.StatusNotFound {
					return fmt.Errorf("Error during creating a new subnet on IpBlock update: %v", err)
				}
			}
		}
		// Go over the old subnets and check which should be deleted
		for _, old_sub_data := range old_subnets {
			old_subnet := old_sub_data.(map[string]interface{})
			found := false
			for _, new_sub_data := range new_subnets {
				new_subnet := new_sub_data.(map[string]interface{})
				if old_subnet["id"] == new_subnet["id"]{
					found = true
					break
				}
			}
			if !found {
				// Delete this subnet
				_, err := nsxClient.PoolManagementApi.DeleteIpBlockSubnet(nsxClient.Context, old_subnet["id"].(string))
				if err != nil {
					return fmt.Errorf("Error during IpBlockSubnet %s delete: %v", old_subnet["id"].(string), err)
				}
			}
		}
	}

	return resourceNsxtIpBlockRead(d, m)
}

func resourceNsxtIpBlockDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	// Delete the IP block with all the subnets
	err := resourceNsxtIpBlockForceDelete(nsxClient, id)
	if err != nil {
		return fmt.Errorf("Error during IpBlock delete: %v", err)
	}
	return nil
}
