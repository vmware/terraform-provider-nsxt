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

var cliProjectDnsRulesClient = dnssvcs.NewRulesClient

var policyDnsRulePathExample = "/orgs/[org]/projects/[project]/dns-services/[dns-service]"

var dnsRuleActionTypeValues = []string{
	model.ProjectDnsRule_ACTION_TYPE_UPDATE_MEMBERSHIP,
	model.ProjectDnsRule_ACTION_TYPE_FORWARD,
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
			"parent_path":  getPolicyPathSchema(true, true, "Policy path of the parent PolicyDnsService"),
			"action_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(dnsRuleActionTypeValues, false),
				Description:  "Action to apply to DNS queries matching domain_patterns. UPDATE_MEMBERSHIP or FORWARD.",
			},
			"domain_patterns": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Domain name patterns matched via longest-prefix match. Supports wildcards (e.g. '*.example.com').",
			},
			"upstream_servers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Upstream DNS server IPs for FORWARD rules. Required when action_type is FORWARD.",
			},
		},
	}
}

func resourceNsxtPolicyDnsRuleExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, policyDnsRulePathExample)
	if pathErr != nil {
		return false, pathErr
	}
	c := cliProjectDnsRulesClient(sessionContext, connector)
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

func policyDnsRuleFromSchema(d *schema.ResourceData) model.ProjectDnsRule {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	actionType := d.Get("action_type").(string)
	domainPatterns := getStringListFromSchemaList(d, "domain_patterns")
	upstreamServers := getStringListFromSchemaList(d, "upstream_servers")

	return model.ProjectDnsRule{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		ActionType:      &actionType,
		DomainPatterns:  domainPatterns,
		UpstreamServers: upstreamServers,
	}
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
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyDnsRuleExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	obj := policyDnsRuleFromSchema(d)
	if *obj.ActionType == model.ProjectDnsRule_ACTION_TYPE_FORWARD && len(obj.UpstreamServers) == 0 {
		return fmt.Errorf("upstream_servers is required when action_type is FORWARD")
	}

	log.Printf("[INFO] Creating ProjectDnsRule with ID %s under %s", id, parentPath)
	c := cliProjectDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	err = c.Patch(org, project, dnsServiceID, id, obj)
	if err != nil {
		return handleCreateError("ProjectDnsRule", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsRuleRead(d, m)
}

func resourceNsxtPolicyDnsRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliProjectDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	obj, err := c.Get(org, project, dnsServiceID, id)
	if err != nil {
		return handleReadError(d, "ProjectDnsRule", id, err)
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
	return nil
}

func resourceNsxtPolicyDnsRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	revision := int64(d.Get("revision").(int))
	obj := policyDnsRuleFromSchema(d)
	obj.Revision = &revision

	if *obj.ActionType == model.ProjectDnsRule_ACTION_TYPE_FORWARD && len(obj.UpstreamServers) == 0 {
		return fmt.Errorf("upstream_servers is required when action_type is FORWARD")
	}

	c := cliProjectDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	_, err := c.Update(org, project, dnsServiceID, id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("ProjectDnsRule", id, err)
	}
	return resourceNsxtPolicyDnsRuleRead(d, m)
}

func resourceNsxtPolicyDnsRuleDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRule ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsRuleParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliProjectDnsRulesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS rule")
	}
	err := c.Delete(org, project, dnsServiceID, id)
	if err != nil {
		return handleDeleteError("ProjectDnsRule", id, err)
	}
	return nil
}
