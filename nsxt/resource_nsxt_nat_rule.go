/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/go-vmware-nsxt/manager"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"net/http"
	"strings"
)

var natRuleActionValues = []string{
	model.PolicyNatRule_ACTION_SNAT,
	model.PolicyNatRule_ACTION_DNAT,
	model.PolicyNatRule_ACTION_REFLEXIVE,
	model.PolicyNatRule_ACTION_NO_SNAT,
	model.PolicyNatRule_ACTION_NO_DNAT,
	model.PolicyNatRule_ACTION_NAT64,
	"NO_NAT", // NSX < 3.0.0 only
}

func resourceNsxtNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtNatRuleCreate,
		Read:   resourceNsxtNatRuleRead,
		Update: resourceNsxtNatRuleUpdate,
		Delete: resourceNsxtNatRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtNatRuleImport,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"action": {
				Type:         schema.TypeString,
				Description:  "The action for the NAT Rule",
				Required:     true,
				ValidateFunc: validation.StringInSlice(natRuleActionValues, false),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Default:     true,
				Description: "enable/disable the rule",
				Optional:    true,
			},
			"logging": {
				Type:        schema.TypeBool,
				Default:     false,
				Description: "enable/disable the logging of rule",
				Optional:    true,
			},
			"logical_router_id": {
				Type:        schema.TypeString,
				Description: "Logical router id",
				Required:    true,
			},
			"match_destination_network": {
				Type:        schema.TypeString,
				Description: "IP Address | CIDR",
				Optional:    true,
			},
			"match_source_network": {
				Type:        schema.TypeString,
				Description: "IP Address | CIDR",
				Optional:    true,
			},
			"nat_pass": {
				Type:        schema.TypeBool,
				Default:     true,
				Description: "A boolean flag which reflects whether the following firewall stage will be skipped",
				Optional:    true,
			},
			"rule_priority": {
				Type:         schema.TypeInt,
				Description:  "The priority of the rule (ascending). Valid range [0-2147483647]",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"translated_network": {
				Type:        schema.TypeString,
				Description: "IP Address | IP Range | CIDR",
				Optional:    true,
			},
			"translated_ports": {
				Type:        schema.TypeString,
				Description: "port number or port range. DNAT only",
				Optional:    true,
			},
			//TODO(asarfaty): Add match_service field
		},
	}
}

func resourceNsxtNatRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	logicalRouterID := d.Get("logical_router_id").(string)
	if logicalRouterID == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	action := d.Get("action").(string)
	if action == "NO_NAT" && nsxVersionHigherOrEqual("3.0.0") {
		return fmt.Errorf("NO_NAT action is not supported in NSX versions 3.0.0 and greater. Use NO_SNAT and NO_DNAT instead")
	}
	enabled := d.Get("enabled").(bool)
	logging := d.Get("logging").(bool)
	matchDestinationNetwork := d.Get("match_destination_network").(string)
	//match_service := d.Get("match_service").(*NsServiceElement)
	matchSourceNetwork := d.Get("match_source_network").(string)
	natPass := d.Get("nat_pass").(bool)
	rulePriority := int64(d.Get("rule_priority").(int))
	translatedNetwork := d.Get("translated_network").(string)
	translatedPorts := d.Get("translated_ports").(string)
	natRule := manager.NatRule{
		Description:             description,
		DisplayName:             displayName,
		Tags:                    tags,
		Action:                  action,
		Enabled:                 enabled,
		Logging:                 logging,
		LogicalRouterId:         logicalRouterID,
		MatchDestinationNetwork: matchDestinationNetwork,
		//MatchService: match_service,
		MatchSourceNetwork: matchSourceNetwork,
		NatPass:            natPass,
		RulePriority:       rulePriority,
		TranslatedNetwork:  translatedNetwork,
		TranslatedPorts:    translatedPorts,
	}

	natRule, resp, err := nsxClient.LogicalRoutingAndServicesApi.AddNatRule(nsxClient.Context, logicalRouterID, natRule)

	if err != nil {
		return fmt.Errorf("Error during NatRule create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NatRule create: %v", resp.StatusCode)
	}
	d.SetId(natRule.Id)

	return resourceNsxtNatRuleRead(d, m)
}

func resourceNsxtNatRuleRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logicalRouterID := d.Get("logical_router_id").(string)
	if logicalRouterID == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	natRule, resp, err := nsxClient.LogicalRoutingAndServicesApi.GetNatRule(nsxClient.Context, logicalRouterID, id)
	if resp != nil && (resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound) {
		// Due to platform bug, 400 response is returned when NAT rule is not found
		// In this case terraform should not error out
		log.Printf("[DEBUG] NatRule %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NatRule read: %v", err)
	}

	d.Set("revision", natRule.Revision)
	d.Set("description", natRule.Description)
	d.Set("display_name", natRule.DisplayName)
	setTagsInSchema(d, natRule.Tags)
	d.Set("action", natRule.Action)
	d.Set("enabled", natRule.Enabled)
	d.Set("logging", natRule.Logging)
	d.Set("logical_router_id", natRule.LogicalRouterId)
	d.Set("match_destination_network", natRule.MatchDestinationNetwork)
	//d.Set("match_service", natRule.MatchService)
	d.Set("match_source_network", natRule.MatchSourceNetwork)
	d.Set("nat_pass", natRule.NatPass)
	d.Set("rule_priority", natRule.RulePriority)
	d.Set("translated_network", natRule.TranslatedNetwork)
	d.Set("translated_ports", natRule.TranslatedPorts)

	return nil
}

func resourceNsxtNatRuleUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logicalRouterID := d.Get("logical_router_id").(string)
	if logicalRouterID == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	action := d.Get("action").(string)
	if action == "NO_NAT" && nsxVersionHigherOrEqual("3.0.0") {
		return fmt.Errorf("NO_NAT action is not supported in NSX versions 3.0.0 and greater. Use NO_SNAT and NO_DNAT instead")
	}
	enabled := d.Get("enabled").(bool)
	logging := d.Get("logging").(bool)
	matchDestinationNetwork := d.Get("match_destination_network").(string)
	//match_service := d.Get("match_service").(*NsServiceElement)
	matchSourceNetwork := d.Get("match_source_network").(string)
	natPass := d.Get("nat_pass").(bool)
	rulePriority := int64(d.Get("rule_priority").(int))
	translatedNetwork := d.Get("translated_network").(string)
	translatedPorts := d.Get("translated_ports").(string)
	natRule := manager.NatRule{
		Revision:                revision,
		Description:             description,
		DisplayName:             displayName,
		Tags:                    tags,
		Action:                  action,
		Enabled:                 enabled,
		Logging:                 logging,
		LogicalRouterId:         logicalRouterID,
		MatchDestinationNetwork: matchDestinationNetwork,
		//MatchService: match_service,
		MatchSourceNetwork: matchSourceNetwork,
		NatPass:            natPass,
		RulePriority:       rulePriority,
		TranslatedNetwork:  translatedNetwork,
		TranslatedPorts:    translatedPorts,
	}

	natRule, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateNatRule(nsxClient.Context, logicalRouterID, id, natRule)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NatRule update: %v", err)
	}

	return resourceNsxtNatRuleRead(d, m)
}

func resourceNsxtNatRuleDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	logicalRouterID := d.Get("logical_router_id").(string)
	if logicalRouterID == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteNatRule(nsxClient.Context, logicalRouterID, id)
	if err != nil {
		return fmt.Errorf("Error during NatRule delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NatRule %s not found", id)
		d.SetId("")
	}
	return nil
}

func resourceNsxtNatRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Please provide <router-id>/<nat-rule-id> as an input")
	}
	d.SetId(s[1])
	d.Set("logical_router_id", s[0])
	return []*schema.ResourceData{d}, nil
}
