// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	dnssvcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/dns_services"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliDnsRulesClient = dnssvcs.NewRulesClient

var policyDnsRulePathExample = "/orgs/[org]/projects/[project]/dns-services/[dns-service]"

var dnsRuleActionTypeValues = []string{
	model.DnsRule_ACTION_TYPE_UPDATE_MEMBERSHIP,
	model.DnsRule_ACTION_TYPE_FORWARD,
}

func resourceNsxtPolicyDnsRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDnsRuleCreate,
		Read:   resourceNsxtPolicyDnsRuleRead,
		Update: resourceNsxtPolicyDnsRuleUpdate,
		Delete: resourceNsxtPolicyDnsRuleDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent DnsService"),
			"action_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(dnsRuleActionTypeValues, false),
				Description:  "Action to apply to DNS queries matching domain_patterns. UPDATE_MEMBERSHIP or FORWARD.",
			},
			"domain_patterns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Domain name patterns matched via longest-prefix match. Supports wildcards (e.g. '*.example.com'). Required unless action_type is FORWARD and shared_zone_path is set.",
			},
			"upstream_servers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Upstream DNS server IPs for FORWARD rules. Mutually exclusive with shared_zone_path.",
			},
			"shared_zone_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to a DnsZone shared with this project. Only valid for FORWARD rules. Mutually exclusive with upstream_servers.",
			},
		},
	}
}

func resourceNsxtPolicyDnsRuleExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, policyDnsRulePathExample)
	if pathErr != nil {
		return false, pathErr
	}
	c := cliDnsRulesClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type for DNS rule")
	}
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func policyDnsRuleFromSchema(d *schema.ResourceData) model.DnsRule {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	actionType := d.Get("action_type").(string)
	domainPatterns := getStringListFromSchemaList(d, "domain_patterns")
	upstreamServers := getStringListFromSchemaList(d, "upstream_servers")

	obj := model.DnsRule{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		ActionType:      &actionType,
		DomainPatterns:  domainPatterns,
		UpstreamServers: upstreamServers,
	}

	if v, ok := d.GetOk("shared_zone_path"); ok {
		szp := v.(string)
		obj.SharedZonePath = &szp
	}

	return obj
}

func dnsRuleParentsFromPath(parentPath string) (org, project, dnsServiceID string, err error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, policyDnsRulePathExample)
	if pathErr != nil {
		return "", "", "", pathErr
	}
	return parents[0], parents[1], parents[2], nil
}

func resourceNsxtPolicyDnsRuleCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("Policy DNS Rule resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyDnsRuleExists)
	if err != nil {
		return err
	}

	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	obj := policyDnsRuleFromSchema(d)
	if *obj.ActionType == model.DnsRule_ACTION_TYPE_FORWARD {
		hasUpstream := len(obj.UpstreamServers) > 0
		hasSharedZone := obj.SharedZonePath != nil && *obj.SharedZonePath != ""
		if hasUpstream == hasSharedZone {
			return fmt.Errorf("exactly one of upstream_servers or shared_zone_path must be set for FORWARD rules")
		}
	}

	log.Printf("[INFO] Creating DnsRule with ID %s under %s", id, parentPath)
	c := cliDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	err = c.Patch(org, project, dnsServiceID, id, obj)
	if err != nil {
		return handleCreateError("DnsRule", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsRuleRead(d, m)
}

func resourceNsxtPolicyDnsRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	obj, err := c.Get(org, project, dnsServiceID, id)
	if err != nil {
		return handleReadError(d, "DnsRule", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("action_type", obj.ActionType)
	d.Set("domain_patterns", obj.DomainPatterns)
	d.Set("upstream_servers", obj.UpstreamServers)
	d.Set("shared_zone_path", obj.SharedZonePath)
	return nil
}

func resourceNsxtPolicyDnsRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	revision := int64(d.Get("revision").(int))
	obj := policyDnsRuleFromSchema(d)
	obj.Revision = &revision

	if *obj.ActionType == model.DnsRule_ACTION_TYPE_FORWARD {
		hasUpstream := len(obj.UpstreamServers) > 0
		hasSharedZone := obj.SharedZonePath != nil && *obj.SharedZonePath != ""
		if hasUpstream == hasSharedZone {
			return fmt.Errorf("exactly one of upstream_servers or shared_zone_path must be set for FORWARD rules")
		}
	}

	c := cliDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	_, err := c.Update(org, project, dnsServiceID, id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("DnsRule", id, err)
	}
	return resourceNsxtPolicyDnsRuleRead(d, m)
}

func resourceNsxtPolicyDnsRuleDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	err := c.Delete(org, project, dnsServiceID, id)
	if err != nil {
		return handleDeleteError("DnsRule", id, err)
	}
	return nil
}
