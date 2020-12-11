/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	gm_segments "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/segments"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyDhcpV6StaticBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDhcpV6StaticBindingCreate,
		Read:   resourceNsxtPolicyDhcpV6StaticBindingRead,
		Update: resourceNsxtPolicyDhcpV6StaticBindingUpdate,
		Delete: resourceNsxtPolicyDhcpV6StaticBindingDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtSegmentResourceImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"segment_path": getPolicyPathSchema(true, true, "segment path"),
			"dns_nameservers": {
				Type:        schema.TypeList,
				Description: "DNS nameservers",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sntp_servers": {
				Type:        schema.TypeList,
				Description: "SNTP server IP addresses",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
			},
			"domain_names": getDomainNamesSchema(),
			"ip_addresses": {
				Type:        schema.TypeList,
				Description: "IP addresses",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
			},
			"lease_time":     getDhcpLeaseTimeSchema(),
			"preferred_time": getDhcpPreferredTimeSchema(),
			"mac_address": {
				Type:         schema.TypeString,
				Description:  "MAC address of the host",
				Required:     true,
				ValidateFunc: validation.IsMACAddress,
			},
		},
	}
}

func policyDhcpV6StaticBindingConvertAndPatch(d *schema.ResourceData, segmentID string, id string, m interface{}) error {

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	ipAddresses := getStringListFromSchemaList(d, "ip_addresses")
	domainNames := getStringListFromSchemaList(d, "domain_names")
	dnsNameservers := getStringListFromSchemaList(d, "dns_nameservers")
	sntpServers := getStringListFromSchemaList(d, "sntp_servers")
	leaseTime := int64(d.Get("lease_time").(int))
	preferredTime := int64(d.Get("preferred_time").(int))
	macAddress := d.Get("mac_address").(string)

	obj := model.DhcpV6StaticBindingConfig{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		LeaseTime:    &leaseTime,
		MacAddress:   &macAddress,
		ResourceType: "DhcpV6StaticBindingConfig",
	}

	if len(ipAddresses) > 0 {
		obj.IpAddresses = ipAddresses
	}

	if len(domainNames) > 0 {
		obj.DomainNames = domainNames
	}

	if len(dnsNameservers) > 0 {
		obj.DnsNameservers = dnsNameservers
	}

	if len(sntpServers) > 0 {
		obj.SntpServers = sntpServers
	}

	if preferredTime > 0 {
		obj.PreferredTime = &preferredTime
	}

	connector := getPolicyConnector(m)

	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	if isPolicyGlobalManager(m) {
		convObj, convErrs := converter.ConvertToVapi(obj, gm_model.DhcpV6StaticBindingConfigBindingType())
		if convErrs != nil {
			return convErrs[0]
		}
		client := gm_segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
		return client.Patch(segmentID, id, convObj.(*data.StructValue))
	}
	convObj, convErrs := converter.ConvertToVapi(obj, model.DhcpV6StaticBindingConfigBindingType())
	if convErrs != nil {
		return convErrs[0]
	}
	client := segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
	return client.Patch(segmentID, id, convObj.(*data.StructValue))
}

func resourceNsxtPolicyDhcpV6StaticBindingCreate(d *schema.ResourceData, m interface{}) error {

	segmentPath := d.Get("segment_path").(string)
	segmentID := getPolicyIDFromPath(segmentPath)
	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDhcpStaticBindingExists(segmentID))
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating DhcpV6StaticBindingConfig with ID %s", id)
	err = policyDhcpV6StaticBindingConvertAndPatch(d, segmentID, id, m)

	if err != nil {
		return handleCreateError("DhcpV6StaticBindingConfig", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDhcpV6StaticBindingRead(d, m)
}

func resourceNsxtPolicyDhcpV6StaticBindingRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}

	segmentPath := d.Get("segment_path").(string)
	segmentID := getPolicyIDFromPath(segmentPath)

	var obj model.DhcpV6StaticBindingConfig
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	var err error
	var dhcpObj *data.StructValue
	if isPolicyGlobalManager(m) {
		client := gm_segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
		dhcpObj, err = client.Get(segmentID, id)

	} else {
		client := segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
		dhcpObj, err = client.Get(segmentID, id)
	}
	if err != nil {
		return handleReadError(d, "DhcpV6StaticBindingConfig", id, err)
	}

	convObj, errs := converter.ConvertToGolang(dhcpObj, model.DhcpV6StaticBindingConfigBindingType())
	if errs != nil {
		return errs[0]
	}
	obj = convObj.(model.DhcpV6StaticBindingConfig)

	if obj.ResourceType != "DhcpV6StaticBindingConfig" {
		return handleReadError(d, "DhcpV6StaticBindingConfig", id, fmt.Errorf("Unexpected ResourceType"))
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("ip_addresses", obj.IpAddresses)
	d.Set("lease_time", obj.LeaseTime)
	d.Set("preferred_time", obj.PreferredTime)
	d.Set("mac_address", obj.MacAddress)
	d.Set("dns_nameservers", obj.DnsNameservers)
	d.Set("domain_names", obj.DomainNames)
	d.Set("sntp_servers", obj.SntpServers)

	return nil
}

func resourceNsxtPolicyDhcpV6StaticBindingUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}
	segmentPath := d.Get("segment_path").(string)
	segmentID := getPolicyIDFromPath(segmentPath)

	log.Printf("[INFO] Updating DhcpV6StaticBindingConfig with ID %s", id)
	err := policyDhcpV6StaticBindingConvertAndPatch(d, segmentID, id, m)
	if err != nil {
		return handleUpdateError("DhcpV6StaticBindingConfig", id, err)
	}

	return resourceNsxtPolicyDhcpV6StaticBindingRead(d, m)
}

func resourceNsxtPolicyDhcpV6StaticBindingDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}
	segmentPath := d.Get("segment_path").(string)
	segmentID := getPolicyIDFromPath(segmentPath)

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
		err = client.Delete(segmentID, id)
	} else {
		client := segments.NewDefaultDhcpStaticBindingConfigsClient(connector)
		err = client.Delete(segmentID, id)
	}

	if err != nil {
		return handleDeleteError("DhcpV6StaticBindingConfig", id, err)
	}

	return nil
}
