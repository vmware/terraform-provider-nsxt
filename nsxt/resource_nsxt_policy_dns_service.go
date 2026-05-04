// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliPolicyDnsServicesClient = projects.NewDnsServicesClient

var policyDnsServicePathExample = "/orgs/[org]/projects/[project]/dns-services/[dns-service]"

func resourceNsxtPolicyDnsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDnsServiceCreate,
		Read:   resourceNsxtPolicyDnsServiceRead,
		Update: resourceNsxtPolicyDnsServiceUpdate,
		Delete: resourceNsxtPolicyDnsServiceDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "Policy DNS Service", getPolicyPathResourceImporter(policyDnsServicePathExample)),
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"allocated_listener_ips": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
				Description: "Listener IP addresses for this DNS service. VPC workloads send DNS queries to these IPs. Maximum of two entries are allowed; if two are provided, one must be IPv4 and one must be IPv6.",
			},
			"transit_gateway": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to the transit gateway providing north-south connectivity.",
			},
			"forwarder_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of DNS cache entries (100-100000).",
						},
						"upstream_servers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 3,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Upstream DNS server IP addresses for catch-all recursive resolution.",
						},
						"shared_zone_forwarding_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "AUTO",
							ValidateFunc: validation.StringInSlice([]string{"AUTO", "MANUAL"}, false),
							Description:  "Controls DNS forwarding rule management for shared zones. AUTO: system creates rules automatically. MANUAL: user must configure DnsRule resources explicitly.",
						},
					},
				},
				Description: "Forwarder and cache settings. When present, enables recursive resolution.",
			},
		},
	}
}

func resourceNsxtPolicyDnsServiceExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliPolicyDnsServicesClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type for DNS service")
	}
	_, err := c.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func policyDnsServiceFromSchema(d *schema.ResourceData) model.DnsService {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	listenerIPs := getStringListFromSchemaList(d, "allocated_listener_ips")

	obj := model.DnsService{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		AllocatedListenerIps: listenerIPs,
	}

	tgw := d.Get("transit_gateway").(string)
	obj.TransitGateway = &tgw

	fwdConfigs := d.Get("forwarder_config").([]interface{})
	if len(fwdConfigs) > 0 {
		fwdMap := fwdConfigs[0].(map[string]interface{})
		fwdCfg := &model.DnsServiceForwarderConfig{}
		if v, ok := fwdMap["cache_size"].(int); ok && v > 0 {
			cacheSize := int64(v)
			fwdCfg.CacheSize = &cacheSize
		}
		if servers, ok := fwdMap["upstream_servers"].([]interface{}); ok {
			for _, s := range servers {
				fwdCfg.UpstreamServers = append(fwdCfg.UpstreamServers, s.(string))
			}
		}
		if mode, ok := fwdMap["shared_zone_forwarding_mode"].(string); ok && mode != "" {
			fwdCfg.SharedZoneForwardingMode = &mode
		}
		obj.ForwarderConfig = fwdCfg
	}

	return obj
}

func policyDnsServiceToSchema(d *schema.ResourceData, obj model.DnsService) {
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("allocated_listener_ips", obj.AllocatedListenerIps)
	d.Set("transit_gateway", obj.TransitGateway)
	if obj.ForwarderConfig != nil {
		fwdMap := map[string]interface{}{
			"upstream_servers": obj.ForwarderConfig.UpstreamServers,
		}
		if obj.ForwarderConfig.CacheSize != nil {
			fwdMap["cache_size"] = int(*obj.ForwarderConfig.CacheSize)
		}
		if obj.ForwarderConfig.SharedZoneForwardingMode != nil {
			fwdMap["shared_zone_forwarding_mode"] = *obj.ForwarderConfig.SharedZoneForwardingMode
		}
		d.Set("forwarder_config", []interface{}{fwdMap})
	} else {
		d.Set("forwarder_config", nil)
	}
}

func resourceNsxtPolicyDnsServiceCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("Policy DNS Service resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyDnsServiceExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(sessionContext)
	obj := policyDnsServiceFromSchema(d)

	log.Printf("[INFO] Creating DnsService with ID %s", id)
	c := cliPolicyDnsServicesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS service")
	}
	err = c.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("DnsService", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsServiceRead(d, m)
}

func resourceNsxtPolicyDnsServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsService ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliPolicyDnsServicesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS service")
	}
	obj, err := c.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "DnsService", id, err)
	}
	policyDnsServiceToSchema(d, obj)
	return nil
}

func resourceNsxtPolicyDnsServiceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsService ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	revision := int64(d.Get("revision").(int))
	obj := policyDnsServiceFromSchema(d)
	obj.Revision = &revision

	c := cliPolicyDnsServicesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS service")
	}
	_, err := c.Update(parents[0], parents[1], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("DnsService", id, err)
	}
	return resourceNsxtPolicyDnsServiceRead(d, m)
}

func resourceNsxtPolicyDnsServiceDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsService ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliPolicyDnsServicesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS service")
	}
	err := c.Delete(parents[0], parents[1], id)
	if err != nil {
		return handleDeleteError("DnsService", id, err)
	}
	return nil
}
