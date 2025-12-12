// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	security_policies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/security_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func resourceNsxtPolicySecurityPolicyContainerCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySecurityPolicyContainerClusterCreate,
		Read:   resourceNsxtPolicySecurityPolicyContainerClusterRead,
		Update: resourceNsxtPolicySecurityPolicyContainerClusterUpdate,
		Delete: resourceNsxtPolicySecurityPolicyContainerClusterDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtSecurityPolicyContainerClusterImporter,
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"policy_path":  getPolicyPathSchema(true, true, "Security Policy path"),
			"container_cluster_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the container cluster entity in NSX",
			},
		},
	}
}

func resourceNsxtPolicySecurityPolicyContainerClusterExistsPartial(d *schema.ResourceData, m interface{}, policyPath string) func(string, client.Connector, bool) (bool, error) {
	return func(id string, connector client.Connector, isGlobal bool) (bool, error) {
		sessionContext := getSessionContext(d, m)
		return resourceNsxtPolicySecurityPolicyContainerClusterExists(id, connector, policyPath, sessionContext)
	}
}

func resourceNsxtPolicySecurityPolicyContainerClusterExists(id string, connector client.Connector, policyPath string, sessionContext utl.SessionContext) (bool, error) {
	var err error

	client := security_policies.NewContainerClusterSpanClient(sessionContext, connector)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	_, err = client.Get(domain, policyID, id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicySecurityPolicyContainerClusterCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := security_policies.NewContainerClusterSpanClient(sessionContext, connector)

	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicySecurityPolicyContainerClusterExistsPartial(d, m, policyPath))
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	containerClusterPath := d.Get("container_cluster_path").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.SecurityPolicyContainerCluster{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		ContainerClusterPath: &containerClusterPath,
	}

	err = client.Patch(domain, policyID, id, obj)
	if err != nil {
		return handleCreateError("SecurityPolicyContainerCluster", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySecurityPolicyContainerClusterRead(d, m)
}

func resourceNsxtPolicySecurityPolicyContainerClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SecurityPolicyContainerCluster ID")
	}
	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	sessionContext := getSessionContext(d, m)
	client := security_policies.NewContainerClusterSpanClient(sessionContext, connector)

	obj, err := client.Get(domain, policyID, id)
	if err != nil {
		return handleReadError(d, "SecurityPolicyContainerCluster", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("container_cluster_path", obj.ContainerClusterPath)

	return nil
}

func resourceNsxtPolicySecurityPolicyContainerClusterUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SecurityPolicyContainerCluster ID")
	}

	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	containerClusterPath := d.Get("container_cluster_path").(string)

	revision := int64(d.Get("revision").(int))

	obj := model.SecurityPolicyContainerCluster{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		Revision:             &revision,
		ContainerClusterPath: &containerClusterPath,
	}

	sessionContext := getSessionContext(d, m)
	client := security_policies.NewContainerClusterSpanClient(sessionContext, connector)
	_, err := client.Update(domain, policyID, id, obj)
	if err != nil {
		return handleUpdateError("SecurityPolicyContainerCluster", id, err)
	}

	return resourceNsxtPolicySecurityPolicyContainerClusterRead(d, m)
}

func resourceNsxtPolicySecurityPolicyContainerClusterDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SecurityPolicyContainerCluster ID")
	}
	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	client := security_policies.NewContainerClusterSpanClient(sessionContext, connector)
	err := client.Delete(domain, policyID, id)

	if err != nil {
		return handleDeleteError("SecurityPolicyContainerCluster", id, err)
	}

	return nil
}

func nsxtSecurityPolicyContainerClusterImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}
	ruleIdx := strings.Index(importID, "container-cluster-span")
	if ruleIdx <= 0 {
		return nil, fmt.Errorf("invalid path of SecurityPolicyContainerCluster to import")
	}
	d.Set("policy_path", importID[:ruleIdx-1])
	return []*schema.ResourceData{d}, nil
}
