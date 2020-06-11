/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_domains "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

func resourceNsxtPolicyGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyGatewayPolicyCreate,
		Read:   resourceNsxtPolicyGatewayPolicyRead,
		Update: resourceNsxtPolicyGatewayPolicyUpdate,
		Delete: resourceNsxtPolicyGatewayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},

		Schema: getPolicyGatewayPolicySchema(),
	}
}

func resourceNsxtPolicyGatewayPolicyExistsInDomain(id string, domainName string, connector *client.RestConnector, isGlobalManager bool) bool {
	var err error
	if isGlobalManager {
		client := gm_domains.NewDefaultGatewayPoliciesClient(connector)
		_, err = client.Get(domainName, id)
	} else {
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		_, err = client.Get(domainName, id)
	}

	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Gateway Policy", err)
	return false
}

func resourceNsxtPolicyGatewayPolicyExistsPartial(domainName string) func(id string, connector *client.RestConnector, isGlobalManager bool) bool {
	return func(id string, connector *client.RestConnector, isGlobalManager bool) bool {
		return resourceNsxtPolicyGatewayPolicyExistsInDomain(id, domainName, connector, isGlobalManager)
	}
}

func resourceNsxtPolicyGatewayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyGatewayPolicyExistsPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	category := d.Get("category").(string)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	rules := getPolicyRulesFromSchema(d)

	obj := model.GatewayPolicy{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Category:       &category,
		Comments:       &comments,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		Rules:          rules,
	}

	_, isSet := d.GetOkExists("tcp_strict")
	if isSet {
		tcpStrict := d.Get("tcp_strict").(bool)
		obj.TcpStrict = &tcpStrict
	}

	log.Printf("[INFO] Creating Gateway Policy with ID %s", id)
	if isPolicyGlobalManager(m) {
		client := gm_domains.NewDefaultGatewayPoliciesClient(connector)
		gmObj, err1 := convertModelBindingType(obj, model.GatewayPolicyBindingType(), gm_model.GatewayPolicyBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(d.Get("domain").(string), id, gmObj.(gm_model.GatewayPolicy))
	} else {
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		err = client.Patch(d.Get("domain").(string), id, obj)
	}
	if err != nil {
		return handleCreateError("Gateway Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyGatewayPolicyRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	var obj model.GatewayPolicy
	if isPolicyGlobalManager(m) {
		client := gm_domains.NewDefaultGatewayPoliciesClient(connector)
		gmObj, err := client.Get(d.Get("domain").(string), id)
		if err != nil {
			return handleReadError(d, "Gateway Policy", id, err)
		}
		rawObj, err := convertModelBindingType(gmObj, gm_model.GatewayPolicyBindingType(), model.GatewayPolicyBindingType())
		if err != nil {
			return err
		}
		obj = rawObj.(model.GatewayPolicy)
	} else {
		var err error
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		obj, err = client.Get(d.Get("domain").(string), id)
		if err != nil {
			return handleReadError(d, "Gateway Policy", id, err)
		}
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
	d.Set("sequence_number", obj.SequenceNumber)
	d.Set("stateful", obj.Stateful)
	if obj.TcpStrict != nil {
		// tcp_strict is dependant on stateful and maybe nil
		d.Set("tcp_strict", *obj.TcpStrict)
	}
	d.Set("revision", obj.Revision)
	setPolicyRulesInSchema(d, obj.Rules)

	return nil
}

func resourceNsxtPolicyGatewayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	category := d.Get("category").(string)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	tcpStrict := d.Get("tcp_strict").(bool)
	rules := getPolicyRulesFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.GatewayPolicy{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Category:       &category,
		Comments:       &comments,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		TcpStrict:      &tcpStrict,
		Revision:       &revision,
		Rules:          rules,
	}

	var err error
	if isPolicyGlobalManager(m) {
		rawObj, err1 := convertModelBindingType(obj, model.GatewayPolicyBindingType(), gm_model.GatewayPolicyBindingType())
		if err1 != nil {
			return err1
		}
		gmObj := rawObj.(gm_model.GatewayPolicy)
		client := gm_domains.NewDefaultGatewayPoliciesClient(connector)
		// We need to use PUT, because PATCH will not replace the whole rule list
		_, err = client.Update(d.Get("domain").(string), id, gmObj)
	} else {
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		// We need to use PUT, because PATCH will not replace the whole rule list
		_, err = client.Update(d.Get("domain").(string), id, obj)
	}

	if err != nil {
		return handleUpdateError("Gateway Policy", id, err)
	}

	return resourceNsxtPolicyGatewayPolicyRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_domains.NewDefaultGatewayPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	} else {
		client := domains.NewDefaultGatewayPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	}
	if err != nil {
		return handleDeleteError("Gateway Policy", id, err)
	}

	return nil
}
