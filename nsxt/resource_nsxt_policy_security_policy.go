/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
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
		Schema: getPolicySecurityPolicySchema(),
	}
}

func resourceNsxtPolicySecurityPolicyExistsInDomain(id string, domainName string, connector *client.RestConnector) bool {
	client := domains.NewDefaultSecurityPoliciesClient(connector)

	_, err := client.Get(domainName, id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Security Policy", err)
	return false
}

func resourceNsxtPolicySecurityPolicyExistsPartial(domainName string) func(id string, connector *client.RestConnector) bool {
	return func(id string, connector *client.RestConnector) bool {
		return resourceNsxtPolicySecurityPolicyExistsInDomain(id, domainName, connector)
	}
}

func resourceNsxtPolicySecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultSecurityPoliciesClient(connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicySecurityPolicyExistsPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

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
	rules := getPolicyRulesFromSchema(d)

	obj := model.SecurityPolicy{
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
		Rules:          rules,
	}

	log.Printf("[INFO] Creating Security Policy with ID %s", id)
	err = client.Patch(d.Get("domain").(string), id, obj)
	if err != nil {
		return handleCreateError("Security Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySecurityPolicyRead(d, m)
}

func resourceNsxtPolicySecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultSecurityPoliciesClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}

	domainName := d.Get("domain").(string)
	obj, err := client.Get(domainName, id)
	if err != nil {
		return handleReadError(d, "Security Policy", id, err)
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
	setPolicyRulesInSchema(d, obj.Rules)

	return nil
}

func resourceNsxtPolicySecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultSecurityPoliciesClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Security Policy id")
	}

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
	rules := getPolicyRulesFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.SecurityPolicy{
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
		Revision:       &revision,
		Rules:          rules,
	}

	// We need to use PUT, because PATCH will not replace the whole rule list
	_, err := client.Update(d.Get("domain").(string), id, obj)
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
	client := domains.NewDefaultSecurityPoliciesClient(connector)

	err := client.Delete(d.Get("domain").(string), id)
	if err != nil {
		return handleDeleteError("Security Policy", id, err)
	}

	return nil
}
