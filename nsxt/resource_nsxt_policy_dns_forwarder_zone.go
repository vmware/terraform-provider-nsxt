/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyDNSForwarderZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDNSForwarderZoneCreate,
		Read:   resourceNsxtPolicyDNSForwarderZoneRead,
		Update: resourceNsxtPolicyDNSForwarderZoneUpdate,
		Delete: resourceNsxtPolicyDNSForwarderZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":           getNsxIDSchema(),
			"path":             getPathSchema(),
			"display_name":     getDisplayNameSchema(),
			"description":      getDescriptionSchema(),
			"revision":         getRevisionSchema(),
			"tag":              getTagsSchema(),
			"dns_domain_names": getDomainNamesSchema(),
			"source_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The source IP used by the DNS Forwarder zone",
				ValidateFunc: validation.IsIPv4Address,
			},
			"upstream_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "DNS servers to which the DNS request needs to be forwarded",
				MaxItems:    3,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},
		},
	}
}

func resourceNsxtPolicyDNSForwarderZoneExists(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewDnsForwarderZonesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDnsForwarderZonesClient(connector)
		_, err = client.Get(id)
	}
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func policyDNSForwarderZonePatch(id string, d *schema.ResourceData, connector client.Connector, isGlobalManager bool) error {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	dnsDomainNames := getStringListFromSchemaList(d, "dns_domain_names")
	sourceIP := d.Get("source_ip").(string)
	upstreamServers := getStringListFromSchemaList(d, "upstream_servers")

	obj := model.PolicyDnsForwarderZone{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		UpstreamServers: upstreamServers,
	}

	if len(dnsDomainNames) > 0 {
		obj.DnsDomainNames = dnsDomainNames
	}

	if len(sourceIP) > 0 {
		obj.SourceIp = &sourceIP
	}

	// Create the resource using PATCH
	if isGlobalManager {
		gmObj, convErr := convertModelBindingType(obj, model.PolicyDnsForwarderZoneBindingType(), gm_model.PolicyDnsForwarderZoneBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewDnsForwarderZonesClient(connector)
		return client.Patch(id, gmObj.(gm_model.PolicyDnsForwarderZone))
	}

	client := infra.NewDnsForwarderZonesClient(connector)
	return client.Patch(id, obj)
}

func resourceNsxtPolicyDNSForwarderZoneCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDNSForwarderZoneExists)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Dns Forwarder Zone with ID %s", id)
	err = policyDNSForwarderZonePatch(id, d, connector, isPolicyGlobalManager(m))

	if err != nil {
		return handleCreateError("Dns Forwarder Zone", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDNSForwarderZoneRead(d, m)
}

func resourceNsxtPolicyDNSForwarderZoneRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Dns Forwarder Zone ID")
	}

	var obj model.PolicyDnsForwarderZone
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDnsForwarderZonesClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "Dns Forwarder Zone", id, err)
		}

		lmObj, err := convertModelBindingType(gmObj, gm_model.PolicyDnsForwarderZoneBindingType(), model.PolicyDnsForwarderZoneBindingType())
		if err != nil {
			return err
		}
		obj = lmObj.(model.PolicyDnsForwarderZone)
	} else {
		client := infra.NewDnsForwarderZonesClient(connector)
		var err error
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "Dns Forwarder Zone", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("dns_domain_names", obj.DnsDomainNames)
	d.Set("source_ip", obj.SourceIp)
	d.Set("upstream_servers", obj.UpstreamServers)

	return nil
}

func resourceNsxtPolicyDNSForwarderZoneUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Dns Forwarder Zone ID")
	}

	log.Printf("[INFO] Updating Dns Forwarder Zone with ID %s", id)
	err := policyDNSForwarderZonePatch(id, d, connector, isPolicyGlobalManager(m))
	if err != nil {
		return handleUpdateError("Dns Forwarder Zone", id, err)
	}

	return resourceNsxtPolicyDNSForwarderZoneRead(d, m)
}

func resourceNsxtPolicyDNSForwarderZoneDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Dns Forwarder Zone ID")
	}

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDnsForwarderZonesClient(connector)
		err = client.Delete(id)
	} else {
		client := infra.NewDnsForwarderZonesClient(connector)
		err = client.Delete(id)
	}

	if err != nil {
		return handleDeleteError("Dns Forwarder Zone", id, err)
	}

	return nil
}
