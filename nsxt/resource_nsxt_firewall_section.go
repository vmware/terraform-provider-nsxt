/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceFirewallSection() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallSectionCreate,
		Read:   resourceFirewallSectionRead,
		Update: resourceFirewallSectionUpdate,
		Delete: resourceFirewallSectionDelete,

		Schema: map[string]*schema.Schema{
			"revision":     getRevisionSchema(),
			"system_owned": getSystemOwnedSchema(),
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
			"tags": getTagsSchema(),
			"is_default": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "It is a boolean flag which reflects whether a firewall section is default section or not. Each Layer 3 and Layer 2 section will have at least and at most one default section",
				Computed:    true,
			},
			"rule_count": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Number of rules in this section",
				Computed:    true,
			},
			"section_type": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Type of the rules which a section can contain. Only homogeneous sections are supported",
				Required:     true,
				ValidateFunc: validateSectionType,
			},
			"stateful": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Stateful or Stateless nature of firewall section is enforced on all rules inside the section. Layer3 sections can be stateful or stateless. Layer2 sections can only be stateless",
				Required:    true,
				ForceNew:    true,
			},
			"applied_tos": getResourceReferencesSchema(false, false, []string{"LogicalPort", "LogicalSwitch", "NSGroup"}),
			"rules":       getRulesSchema(),
		},
	}
}

func getRulesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
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
				"action": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Action enforced on the packets which matches the firewall rule",
					Required:     true,
					ValidateFunc: validateRuleAction,
				},
				"applied_tos":  getResourceReferencesSchema(false, false, []string{"LogicalPort", "LogicalSwitch", "NSGroup"}),
				"destinations": getResourceReferencesSchema(false, false, []string{"IPSet", "LogicalPort", "LogicalSwitch", "NSGroup", "MACSet"}),
				"destinations_excluded": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "Negation of the destination",
					Optional:    true,
				},
				"direction": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Rule direction in case of stateless firewall rules. This will only considered if section level parameter is set to stateless. Default to IN_OUT if not specified",
					Optional:     true,
					ValidateFunc: validateRuleDirection,
				},
				"disabled": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "Flag to disable rule. Disabled will only be persisted but never provisioned/realized",
					Optional:    true,
				},
				"ip_protocol": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Type of IP packet that should be matched while enforcing the rule (IPV4, IPV6, IPV4_IPV6)",
					Optional:     true,
					ValidateFunc: validateRuleIPProtocol,
				},
				"logged": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "Flag to enable packet logging. Default is disabled",
					Optional:    true,
				},
				"notes": &schema.Schema{
					Type:        schema.TypeString,
					Description: "User notes specific to the rule",
					Optional:    true,
				},
				"rule_tag": &schema.Schema{
					Type:        schema.TypeString,
					Description: "User level field which will be printed in CLI and packet logs",
					Optional:    true,
				},
				"sources": getResourceReferencesSchema(false, false, []string{"IPSet", "LogicalPort", "LogicalSwitch", "NSGroup", "MACSet"}),
				"sources_excluded": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "Negation of the source",
					Optional:    true,
				},
				"services": getResourceReferencesSchema(false, false, []string{"NSService", "NSServiceGroup"}),
			},
		},
	}
}

func validateRuleIPProtocol(v interface{}, k string) (ws []string, errors []error) {
	legal_values := []string{"IPV4", "IPV6", "IPV4_IPV6"}
	return validateValueInList(v, k, legal_values)
}

func validateRuleAction(v interface{}, k string) (ws []string, errors []error) {
	legal_values := []string{"ALLOW", "DROP", "REJECT"}
	return validateValueInList(v, k, legal_values)
}

func validateRuleDirection(v interface{}, k string) (ws []string, errors []error) {
	legal_values := []string{"IN", "OUT", "IN_OUT"}
	return validateValueInList(v, k, legal_values)
}

func validateSectionType(v interface{}, k string) (ws []string, errors []error) {
	legal_values := []string{"LAYER2", "LAYER3"}
	return validateValueInList(v, k, legal_values)
}

func returnServicesResourceReferences(services []manager.FirewallService) []map[string]interface{} {
	var servicesList []map[string]interface{}
	for _, srv := range services {
		elem := make(map[string]interface{})
		elem["is_valid"] = srv.IsValid
		elem["target_display_name"] = srv.TargetDisplayName
		elem["target_id"] = srv.TargetId
		elem["target_type"] = srv.TargetType
		servicesList = append(servicesList, elem)
	}
	return servicesList
}

func setRulesInSchema(d *schema.ResourceData, rules []manager.FirewallRule) {
	var rulesList []map[string]interface{}
	for _, rule := range rules {
		elem := make(map[string]interface{})
		elem["id"] = rule.Id
		elem["display_name"] = rule.DisplayName
		elem["description"] = rule.Description
		elem["rule_tag"] = rule.RuleTag
		elem["notes"] = rule.Notes
		elem["logged"] = rule.Logged
		elem["action"] = rule.Action
		elem["destinations_excluded"] = rule.DestinationsExcluded
		elem["sources_excluded"] = rule.SourcesExcluded
		elem["ip_protocol"] = rule.IpProtocol
		elem["disabled"] = rule.Disabled
		elem["revision"] = rule.Revision
		elem["direction"] = rule.Direction
		elem["sources"] = returnResourceReferences(rule.Sources)
		elem["destinations"] = returnResourceReferences(rule.Destinations)
		elem["services"] = returnServicesResourceReferences(rule.Services)

		rulesList = append(rulesList, elem)
	}
	d.Set("rules", rulesList)
}

func getServicesResourceReferences(services []interface{}) []manager.FirewallService {
	var servicesList []manager.FirewallService
	for _, srv := range services {
		data := srv.(map[string]interface{})
		elem := manager.FirewallService{
			IsValid:           data["is_valid"].(bool),
			TargetDisplayName: data["target_display_name"].(string),
			TargetId:          data["target_id"].(string),
			TargetType:        data["target_type"].(string),
		}
		servicesList = append(servicesList, elem)
	}
	return servicesList
}

func getRulesFromSchema(d *schema.ResourceData) []manager.FirewallRule {
	rules := d.Get("rules").([]interface{})
	var ruleList []manager.FirewallRule
	for _, rule := range rules {
		data := rule.(map[string]interface{})
		elem := manager.FirewallRule{
			DisplayName:          data["display_name"].(string),
			RuleTag:              data["rule_tag"].(string),
			Notes:                data["notes"].(string),
			Description:          data["description"].(string),
			Action:               data["action"].(string),
			Logged:               data["logged"].(bool),
			Disabled:             data["disabled"].(bool),
			Revision:             int64(data["revision"].(int)),
			SourcesExcluded:      data["sources_excluded"].(bool),
			DestinationsExcluded: data["destinations_excluded"].(bool),
			IpProtocol:           data["ip_protocol"].(string),
			Direction:            data["direction"].(string),
			Sources:              getResourceReferences(data["sources"].([]interface{})),
			Destinations:         getResourceReferences(data["destinations"].([]interface{})),
			Services:             getServicesResourceReferences(data["services"].([]interface{})),
		}

		ruleList = append(ruleList, elem)
	}
	return ruleList
}

func resourceFirewallSectionCreateEmpty(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
	is_default := d.Get("is_default").(bool)
	rule_count := int64(d.Get("rule_count").(int))
	section_type := d.Get("section_type").(string)
	stateful := d.Get("stateful").(bool)

	localVarOptionals := make(map[string]interface{})
	firewall_section := manager.FirewallSection{
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		AppliedTos:  applied_tos,
		IsDefault:   is_default,
		RuleCount:   rule_count,
		SectionType: section_type,
		Stateful:    stateful,
	}
	firewall_section, resp, err := nsxClient.ServicesApi.AddSection(nsxClient.Context, firewall_section, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during FirewallSection create empty: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during FirewallSection create empty: %v", resp.StatusCode)
	}
	d.SetId(firewall_section.Id)

	return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionCreate(d *schema.ResourceData, m interface{}) error {

	rules := getRulesFromSchema(d)
	if len(rules) == 0 {
		return resourceFirewallSectionCreateEmpty(d, m)
	}

	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
	is_default := d.Get("is_default").(bool)
	rule_count := int64(d.Get("rule_count").(int))
	section_type := d.Get("section_type").(string)
	stateful := d.Get("stateful").(bool)

	firewall_section := manager.FirewallSectionRuleList{
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		AppliedTos:  applied_tos,
		IsDefault:   is_default,
		RuleCount:   rule_count,
		SectionType: section_type,
		Stateful:    stateful,
		Rules:       rules,
	}
	localVarOptionals := make(map[string]interface{})
	firewall_section, resp, err := nsxClient.ServicesApi.AddSectionWithRulesCreateWithRules(nsxClient.Context, firewall_section, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during FirewallSection create with rules: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during FirewallSection create with rules: %v", resp.StatusCode)
	}
	d.SetId(firewall_section.Id)

	return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	firewall_section, resp, err := nsxClient.ServicesApi.GetSectionWithRulesListWithRules(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("FirewallSection %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during FirewallSection %s read: %v", id, err)
	}

	d.Set("revision", firewall_section.Revision)
	d.Set("system_owned", firewall_section.SystemOwned)
	d.Set("description", firewall_section.Description)
	d.Set("display_name", firewall_section.DisplayName)
	setTagsInSchema(d, firewall_section.Tags)
	setRulesInSchema(d, firewall_section.Rules)
	d.Set("is_default", firewall_section.IsDefault)
	d.Set("rule_count", firewall_section.RuleCount)
	d.Set("section_type", firewall_section.SectionType)
	d.Set("stateful", firewall_section.Stateful)

	// Getting the applied tos will require another api call
	firewall_section2, resp, err := nsxClient.ServicesApi.GetSection(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("FirewallSection %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during FirewallSection %s read: %v", id, err)
	}
	setResourceReferencesInSchema(d, firewall_section2.AppliedTos, "applied_tos")

	return nil
}

func resourceFirewallSectionUpdateEmpty(d *schema.ResourceData, m interface{}, id string) error {

	nsxClient := m.(*api.APIClient)
	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
	is_default := d.Get("is_default").(bool)
	rule_count := int64(d.Get("rule_count").(int))
	section_type := d.Get("section_type").(string)
	stateful := d.Get("stateful").(bool)
	firewall_section := manager.FirewallSection{
		Revision:    revision,
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		AppliedTos:  applied_tos,
		IsDefault:   is_default,
		RuleCount:   rule_count,
		SectionType: section_type,
		Stateful:    stateful,
	}
	// Update the section ignoring the rules
	firewall_section, resp, err := nsxClient.ServicesApi.UpdateSection(nsxClient.Context, id, firewall_section)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during FirewallSection %s update empty: %v", id, err)
	}

	// Read the section, and delete all current rules from it
	curr_section, resp2, err2 := nsxClient.ServicesApi.GetSectionWithRulesListWithRules(nsxClient.Context, id)
	if resp2.StatusCode == http.StatusNotFound {
		return fmt.Errorf("FirewallSection %s not found during update empty action", id)
	}
	if err2 != nil {
		return fmt.Errorf("Error during FirewallSection %s update empty: cannot read the section: %v", id, err2)
	}
	for _, rule := range curr_section.Rules {
		nsxClient.ServicesApi.DeleteRule(nsxClient.Context, id, rule.Id)
	}
	return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	rules := getRulesFromSchema(d)
	if len(rules) == 0 {
		return resourceFirewallSectionUpdateEmpty(d, m, id)
	}

	nsxClient := m.(*api.APIClient)

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	applied_tos := getResourceReferencesFromSchema(d, "applied_tos")
	is_default := d.Get("is_default").(bool)
	rule_count := int64(d.Get("rule_count").(int))
	section_type := d.Get("section_type").(string)
	stateful := d.Get("stateful").(bool)
	firewall_section := manager.FirewallSectionRuleList{
		Revision:    revision,
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		AppliedTos:  applied_tos,
		IsDefault:   is_default,
		RuleCount:   rule_count,
		SectionType: section_type,
		Stateful:    stateful,
		Rules:       rules,
	}

	firewall_section, resp, err := nsxClient.ServicesApi.UpdateSectionWithRulesUpdateWithRules(nsxClient.Context, id, firewall_section)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during FirewallSection %s update: %v", id, err)
	}

	return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	localVarOptionals["cascade"] = true
	resp, err := nsxClient.ServicesApi.DeleteSection(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during FirewallSection %s delete: %v", id, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("FirewallSection %s not found", id)
		d.SetId("")
	}
	return nil
}
