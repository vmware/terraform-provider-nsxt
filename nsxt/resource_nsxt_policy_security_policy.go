/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func resourceNsxtPolicySecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySecurityPolicyCreate,
		Read:   resourceNsxtPolicySecurityPolicyRead,
		Update: resourceNsxtPolicySecurityPolicyUpdate,
		Delete: resourceNsxtPolicySecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},
		Schema: getPolicySecurityPolicySchema(false, true),
	}
}

func getSecurityPolicyInDomain(sessionContext utl.SessionContext, id string, domainName string, connector client.Connector) (model.SecurityPolicy, error) {
	client := domains.NewSecurityPoliciesClient(sessionContext, connector)
	return client.Get(domainName, id)

}

func resourceNsxtPolicySecurityPolicyExistsInDomain(sessionContext utl.SessionContext, id string, domainName string, connector client.Connector) (bool, error) {
	client := domains.NewSecurityPoliciesClient(sessionContext, connector)
	_, err := client.Get(domainName, id)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving Security Policy", err)
}

func resourceNsxtPolicySecurityPolicyExistsPartial(domainName string) func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	return func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
		return resourceNsxtPolicySecurityPolicyExistsInDomain(sessionContext, id, domainName, connector)
	}
}

func policySecurityPolicyBuildAndPatch(d *schema.ResourceData, m interface{}, connector client.Connector, isGlobalManager bool, id string, createFlow bool) error {

	domain := d.Get("domain").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	category := d.Get("category").(string)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	scope := getStringListFromSchemaSet(d, "scope")
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	tcpStrict := d.Get("tcp_strict").(bool)
	revision := int64(d.Get("revision").(int))
	objType := "SecurityPolicy"

	obj := model.SecurityPolicy{
		Id:             &id,
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Category:       &category,
		Comments:       &comments,
		Locked:         &locked,
		Scope:          scope,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		TcpStrict:      &tcpStrict,
		ResourceType:   &objType,
	}
	log.Printf("[INFO] Creating Security Policy with ID %s", id)

	if !createFlow {
		// This is update flow
		obj.Revision = &revision
	}

	err := validatePolicyRuleSequence(d)
	if err != nil {
		return err
	}

	policyChildren, err := getUpdatedRuleChildren(d)
	if err != nil {
		return err
	}
	if len(policyChildren) > 0 {
		obj.Children = policyChildren
	}

	log.Printf("[INFO] Using selective H-API for policy with ID %s", id)
	return securityPolicyInfraPatch(getSessionContext(d, m), obj, domain, m)
}

func resourceNsxtPolicySecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicySecurityPolicyExistsPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

	err = policySecurityPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id, true)

	if err != nil {
		return handleCreateError("Security Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySecurityPolicyRead(d, m)
}

func resourceNsxtPolicySecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	domainName := d.Get("domain").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}
	client := domains.NewSecurityPoliciesClient(getSessionContext(d, m), connector)
	obj, err := client.Get(domainName, id)
	if err != nil {
		return handleReadError(d, "SecurityPolicy", id, err)
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("category", obj.Category)
	d.Set("comments", obj.Comments)
	d.Set("locked", obj.Locked)
	if len(obj.Scope) == 1 && obj.Scope[0] == "ANY" {
		d.Set("scope", nil)
	} else {
		d.Set("scope", obj.Scope)
	}
	d.Set("sequence_number", obj.SequenceNumber)
	d.Set("stateful", obj.Stateful)
	d.Set("tcp_strict", obj.TcpStrict)
	d.Set("revision", obj.Revision)
	return setPolicyRulesInSchema(d, obj.Rules)
}

func resourceNsxtPolicySecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}
	err := policySecurityPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id, false)
	if err != nil {
		return handleUpdateError("Security Policy", id, err)
	}

	return resourceNsxtPolicySecurityPolicyRead(d, m)
}

func resourceNsxtPolicySecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}

	connector := getPolicyConnector(m)

	client := domains.NewSecurityPoliciesClient(getSessionContext(d, m), connector)
	err := client.Delete(d.Get("domain").(string), id)

	if err != nil {
		return handleDeleteError("Security Policy", id, err)
	}

	return nil
}
