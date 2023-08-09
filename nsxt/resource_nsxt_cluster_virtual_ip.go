/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/cluster"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

var ClusterVirtualIPForceType = []string{
	nsxModel.ClusterVirtualIpProperties_FORCE_TRUE,
	nsxModel.ClusterVirtualIpProperties_FORCE_FALSE,
}

var DefaultIPv4VirtualAddress = "0.0.0.0"

var DefaultIPv6VirtualAddress = "::"

func resourceNsxtClusterVirualIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtClusterVirualIPCreate,
		Read:   resourceNsxtClusterVirualIPRead,
		Update: resourceNsxtClusterVirualIPUpdate,
		Delete: resourceNsxtClusterVirualIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"force": {
				Type:        schema.TypeBool,
				Description: "On enable it ignores duplicate address detection and DNS lookup validation check",
				Optional:    true,
				Default:     true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Description:  "Virtual IPv4 address",
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
				Default:      DefaultIPv4VirtualAddress,
			},
			"ipv6_address": {
				Type:         schema.TypeString,
				Description:  "Virtual IPv6 address",
				Optional:     true,
				ValidateFunc: validation.IsIPv6Address,
				Default:      DefaultIPv6VirtualAddress,
			},
		},
	}
}

func resourceNsxtClusterVirualIPCreate(d *schema.ResourceData, m interface{}) error {
	// Create and update workflow are mostly the same for virtual IP resource
	// except that create workflow sets the ID of this resource
	id := d.Id()
	if id == "" {
		id = newUUID()
	}
	d.SetId(id)
	err := setClusterVirtualIP(d, m)
	if err != nil {
		return err
	}
	return resourceNsxtClusterVirualIPRead(d, m)
}

func resourceNsxtClusterVirualIPRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cluster.NewApiVirtualIpClient(connector)

	obj, err := client.Get()
	if err != nil {
		return err
	}

	// For some reason the Get() function of ApiVirtulIPClient will only return ip address information
	// so skip setting force here
	d.Set("ip_address", obj.IpAddress)
	d.Set("ipv6_address", obj.Ip6Address)

	return nil
}

func setClusterVirtualIP(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cluster.NewApiVirtualIpClient(connector)
	force := d.Get("force").(bool)
	ipAddress := d.Get("ip_address").(string)
	ipv6Address := d.Get("ipv6_address").(string)
	var forceStr string
	if force {
		forceStr = nsxModel.ClusterVirtualIpProperties_FORCE_TRUE
	} else {
		forceStr = nsxModel.ClusterVirtualIpProperties_FORCE_FALSE
	}
	_, err := client.Setvirtualip(&forceStr, &ipv6Address, &ipAddress)
	if err != nil {
		return fmt.Errorf("Failed to set cluster virtual ip: %s", err)
	}
	return nil
}

func resourceNsxtClusterVirualIPUpdate(d *schema.ResourceData, m interface{}) error {
	err := setClusterVirtualIP(d, m)
	if err != nil {
		return err
	}
	return resourceNsxtClusterVirualIPRead(d, m)
}

func resourceNsxtClusterVirualIPDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := cluster.NewApiVirtualIpClient(connector)
	_, err := client.Clearvirtualip()
	if err != nil {
		return fmt.Errorf("Failed to clear cluster virtual IPv4 address: %s", err)
	}
	_, err = client.Clearvirtualip6()
	if err != nil {
		return fmt.Errorf("Failed to clear cluster virtual IPv6 address: %s", err)
	}
	return nil
}
