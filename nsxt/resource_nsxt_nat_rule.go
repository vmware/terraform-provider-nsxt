/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

var natRuleActionValues = []string{"SNAT", "DNAT", "NO_NAT", "REFLEXIVE"}

func resourceNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatRuleCreate,
		Read:   resourceNatRuleRead,
		Update: resourceNatRuleUpdate,
		Delete: resourceNatRuleDelete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"tag": getTagsSchema(),
			"action": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "valid actions: SNAT, DNAT, NO_NAT, REFLEXIVE. All rules in a logical router are either stateless or stateful. Mix is not supported. SNAT and DNAT are stateful, can NOT be supported when the logical router is running at active-active HA mode; REFLEXIVE is stateless. NO_NAT has no translated_fields, only match fields",
				Required:     true,
				ValidateFunc: validation.StringInSlice(natRuleActionValues, false),
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     true,
				Description: "enable/disable the rule",
				Optional:    true,
			},
			"logging": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Description: "enable/disable the logging of rule",
				Optional:    true,
			},
			"logical_router_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Logical router id",
				Required:    true,
			},
			"match_destination_network": &schema.Schema{
				Type:        schema.TypeString,
				Description: "IP Address | CIDR | (null implies Any)",
				Optional:    true,
			},
			"match_source_network": &schema.Schema{
				Type:        schema.TypeString,
				Description: "IP Address | CIDR | (null implies Any)",
				Optional:    true,
			},
			"nat_pass": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     true,
				Description: "Default is true. If the nat_pass is set to true, the following firewall stage will be skipped. Please note, if action is NO_NAT, then nat_pass must be set to true or omitted",
				Optional:    true,
			},
			"rule_priority": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Ascending, valid range [0-2147483647]. If multiple rules have the same priority, evaluation sequence is undefined",
				Computed:    true,
			},
			"translated_network": &schema.Schema{
				Type:        schema.TypeString,
				Description: "IP Address | IP Range | CIDR. For DNAT rules only a single ip is supported",
				Optional:    true,
			},
			"translated_ports": &schema.Schema{
				Type:        schema.TypeString,
				Description: "port number or port range. DNAT only",
				Optional:    true,
			},
			//TODO(asarfaty): Add match_service field
		},
	}
}

func resourceNatRuleCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	action := d.Get("action").(string)
	enabled := d.Get("enabled").(bool)
	logging := d.Get("logging").(bool)
	match_destination_network := d.Get("match_destination_network").(string)
	//match_service := d.Get("match_service").(*NsServiceElement)
	match_source_network := d.Get("match_source_network").(string)
	nat_pass := d.Get("nat_pass").(bool)
	rule_priority := int64(d.Get("rule_priority").(int))
	translated_network := d.Get("translated_network").(string)
	translated_ports := d.Get("translated_ports").(string)
	nat_rule := manager.NatRule{
		Description:             description,
		DisplayName:             display_name,
		Tags:                    tags,
		Action:                  action,
		Enabled:                 enabled,
		Logging:                 logging,
		LogicalRouterId:         logical_router_id,
		MatchDestinationNetwork: match_destination_network,
		//MatchService: match_service,
		MatchSourceNetwork: match_source_network,
		NatPass:            nat_pass,
		RulePriority:       rule_priority,
		TranslatedNetwork:  translated_network,
		TranslatedPorts:    translated_ports,
	}

	nat_rule, resp, err := nsxClient.LogicalRoutingAndServicesApi.AddNatRule(nsxClient.Context, logical_router_id, nat_rule)

	if err != nil {
		return fmt.Errorf("Error during NatRule create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NatRule create: %v", resp.StatusCode)
	}
	d.SetId(nat_rule.Id)

	return resourceNatRuleRead(d, m)
}

func resourceNatRuleRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	nat_rule, resp, err := nsxClient.LogicalRoutingAndServicesApi.GetNatRule(nsxClient.Context, logical_router_id, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NatRule %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NatRule read: %v", err)
	}

	d.Set("revision", nat_rule.Revision)
	d.Set("description", nat_rule.Description)
	d.Set("display_name", nat_rule.DisplayName)
	setTagsInSchema(d, nat_rule.Tags)
	d.Set("action", nat_rule.Action)
	d.Set("enabled", nat_rule.Enabled)
	d.Set("logging", nat_rule.Logging)
	d.Set("logical_router_id", nat_rule.LogicalRouterId)
	d.Set("match_destination_network", nat_rule.MatchDestinationNetwork)
	//d.Set("match_service", nat_rule.MatchService)
	d.Set("match_source_network", nat_rule.MatchSourceNetwork)
	d.Set("nat_pass", nat_rule.NatPass)
	d.Set("rule_priority", nat_rule.RulePriority)
	d.Set("translated_network", nat_rule.TranslatedNetwork)
	d.Set("translated_ports", nat_rule.TranslatedPorts)

	return nil
}

func resourceNatRuleUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	action := d.Get("action").(string)
	enabled := d.Get("enabled").(bool)
	logging := d.Get("logging").(bool)
	match_destination_network := d.Get("match_destination_network").(string)
	//match_service := d.Get("match_service").(*NsServiceElement)
	match_source_network := d.Get("match_source_network").(string)
	nat_pass := d.Get("nat_pass").(bool)
	rule_priority := int64(d.Get("rule_priority").(int))
	translated_network := d.Get("translated_network").(string)
	translated_ports := d.Get("translated_ports").(string)
	nat_rule := manager.NatRule{
		Revision:                revision,
		Description:             description,
		DisplayName:             display_name,
		Tags:                    tags,
		Action:                  action,
		Enabled:                 enabled,
		Logging:                 logging,
		LogicalRouterId:         logical_router_id,
		MatchDestinationNetwork: match_destination_network,
		//MatchService: match_service,
		MatchSourceNetwork: match_source_network,
		NatPass:            nat_pass,
		RulePriority:       rule_priority,
		TranslatedNetwork:  translated_network,
		TranslatedPorts:    translated_ports,
	}

	nat_rule, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateNatRule(nsxClient.Context, logical_router_id, id, nat_rule)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NatRule update: %v", err)
	}

	return resourceNatRuleRead(d, m)
}

func resourceNatRuleDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteNatRule(nsxClient.Context, logical_router_id, id)
	if err != nil {
		return fmt.Errorf("Error during NatRule delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NatRule %s not found", id)
		d.SetId("")
	}
	return nil
}
