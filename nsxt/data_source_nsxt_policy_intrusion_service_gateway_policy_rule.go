// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDataSourceDomainNameSchema(),
			"policy_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy path for this rule",
			},
			"sequence_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Sequence number of the rule",
			},
			"source_groups": {
				Type:        schema.TypeList,
				Description: "List of source groups",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"destination_groups": {
				Type:        schema.TypeList,
				Description: "List of destination groups",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"services": {
				Type:        schema.TypeList,
				Description: "List of services",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scope": {
				Type:        schema.TypeList,
				Description: "List of policy objects where the rule is enforced",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"action": {
				Type:        schema.TypeString,
				Description: "Action for this rule",
				Computed:    true,
			},
			"direction": {
				Type:         schema.TypeString,
				Description:  "Traffic direction",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"IN", "OUT", "IN_OUT"}, false),
			},
			"disabled": {
				Type:        schema.TypeBool,
				Description: "Flag to disable the rule",
				Computed:    true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Description:  "IP version",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "IPV4_IPV6"}, false),
			},
			"logged": {
				Type:        schema.TypeBool,
				Description: "Flag to enable logging",
				Computed:    true,
			},
			"log_label": {
				Type:        schema.TypeString,
				Description: "Additional information (string) which will be propagated to the rule syslog",
				Computed:    true,
			},
			"notes": {
				Type:        schema.TypeString,
				Description: "Text for additional notes on changes for the rule",
				Computed:    true,
			},
			"sources_excluded": {
				Type:        schema.TypeBool,
				Description: "Flag to indicate whether sources are negated",
				Computed:    true,
			},
			"destinations_excluded": {
				Type:        schema.TypeBool,
				Description: "Flag to indicate whether destinations are negated",
				Computed:    true,
			},
			"ids_profiles": {
				Type:        schema.TypeList,
				Description: "List of IDS profiles for this rule",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique positive number that is assigned by the system and is useful for debugging",
			},
			"tag":      getTagsSchema(),
			"revision": getRevisionSchema(),
		},
	}
}

func listIntrusionServiceGatewayPolicyRules(d *schema.ResourceData, policyPath string, m interface{}) ([]model.IdsRule, error) {
	connector := getPolicyConnector(m)

	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.IdsRule
	var cursor *string
	total := 0

	for {
		includeMarkForDeleteObjectsParam := false
		rules, err := client.List(domain, policyID, cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, rules.Results...)
		if total == 0 && rules.ResultCount != nil {
			total = int(*rules.ResultCount)
		}

		cursor = rules.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	policyPath := d.Get("policy_path").(string)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	if policyPath == "" {
		return fmt.Errorf("policy_path must be specified")
	}

	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	if objID != "" {
		connector := getPolicyConnector(m)
		client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		obj, err := client.Get(domain, policyID, objID)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("Intrusion Service Gateway Policy Rule with ID %s was not found in policy %s", objID, policyID)
			}
			return fmt.Errorf("Error while reading Intrusion Service Gateway Policy Rule %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		setIntrusionServiceGatewayPolicyRuleDataSourceSchema(d, obj)
		return nil
	}

	if objName == "" {
		return fmt.Errorf("Intrusion Service Gateway Policy Rule id or display_name must be specified")
	}

	rules, err := listIntrusionServiceGatewayPolicyRules(d, policyPath, m)
	if err != nil {
		return fmt.Errorf("Error while reading Intrusion Service Gateway Policy Rules: %v", err)
	}

	var perfectMatch []model.IdsRule
	var prefixMatch []model.IdsRule
	for _, rule := range rules {
		if rule.DisplayName != nil {
			if strings.HasPrefix(*rule.DisplayName, objName) {
				prefixMatch = append(prefixMatch, rule)
			}
			if *rule.DisplayName == objName {
				perfectMatch = append(perfectMatch, rule)
			}
		}
	}

	var obj model.IdsRule
	if len(perfectMatch) > 0 {
		if len(perfectMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Gateway Policy Rules with name '%s'", objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Gateway Policy Rules with name starting with '%s'", objName)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Intrusion Service Gateway Policy Rule with name '%s' was not found in policy %s", objName, policyID)
	}

	d.SetId(*obj.Id)
	setIntrusionServiceGatewayPolicyRuleDataSourceSchema(d, obj)
	return nil
}

func setIntrusionServiceGatewayPolicyRuleDataSourceSchema(d *schema.ResourceData, obj model.IdsRule) {
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("notes", obj.Notes)
	d.Set("logged", obj.Logged)
	d.Set("log_label", obj.Tag)
	d.Set("action", obj.Action)
	d.Set("destinations_excluded", obj.DestinationsExcluded)
	d.Set("sources_excluded", obj.SourcesExcluded)
	d.Set("ip_version", obj.IpProtocol)
	d.Set("direction", obj.Direction)
	d.Set("disabled", obj.Disabled)
	if obj.SequenceNumber != nil {
		d.Set("sequence_number", int(*obj.SequenceNumber))
	}
	d.Set("destination_groups", obj.DestinationGroups)
	d.Set("source_groups", obj.SourceGroups)
	d.Set("services", obj.Services)
	d.Set("scope", obj.Scope)
	d.Set("ids_profiles", obj.IdsProfiles)
	if obj.RuleId != nil {
		d.Set("rule_id", *obj.RuleId)
	}
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("revision", obj.Revision)
}
