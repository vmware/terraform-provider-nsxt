// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliDnsAutoRecordConfigsClient = projects.NewDnsAutoRecordConfigsClient

var policyDnsAutoRecordConfigPathExample = "/orgs/[org]/projects/[project]/dns-auto-record-configs/[config]"

func resourceNsxtPolicyDnsRecordAutoConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDnsRecordAutoConfigCreate,
		Read:   resourceNsxtPolicyDnsRecordAutoConfigRead,
		Update: resourceNsxtPolicyDnsRecordAutoConfigUpdate,
		Delete: resourceNsxtPolicyDnsRecordAutoConfigDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "Policy DNS Record Auto Config", getPolicyPathResourceImporter(policyDnsAutoRecordConfigPathExample)),
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"ip_block_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to the IP block from which workload IPs are allocated. Immutable after creation.",
			},
			"a_record_zone_path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to a DnsZone (local or shared) in which auto-created A records are placed.",
			},
			"ptr_record_zone_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Optional policy path to a DnsZone for auto-created PTR (reverse DNS) records. When absent, no PTR record is auto-created.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Time-To-Live in seconds for auto-created A records (30-86400). Default: 300.",
			},
		},
	}
}

func resourceNsxtPolicyDnsRecordAutoConfigExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type for DNS auto record config")
	}
	_, err := c.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func policyDnsRecordAutoConfigFromSchema(d *schema.ResourceData) model.DnsAutoRecordConfig {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	ipBlockPath := d.Get("ip_block_path").(string)
	aRecordZonePath := d.Get("a_record_zone_path").(string)
	ttl := int64(d.Get("ttl").(int))

	obj := model.DnsAutoRecordConfig{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		IpBlockPath:     &ipBlockPath,
		ARecordZonePath: &aRecordZonePath,
		Ttl:             &ttl,
	}

	if v, ok := d.GetOk("ptr_record_zone_path"); ok {
		ptrPath := v.(string)
		obj.PtrRecordZonePath = &ptrPath
	}

	return obj
}

func resourceNsxtPolicyDnsRecordAutoConfigCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("Policy DNS Record Auto Config resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyDnsRecordAutoConfigExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(sessionContext)
	obj := policyDnsRecordAutoConfigFromSchema(d)

	log.Printf("[INFO] Creating DnsAutoRecordConfig with ID %s", id)
	c := cliDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	err = c.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("DnsAutoRecordConfig", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsRecordAutoConfigRead(d, m)
}

func resourceNsxtPolicyDnsRecordAutoConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	obj, err := c.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "DnsAutoRecordConfig", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("ip_block_path", obj.IpBlockPath)
	d.Set("a_record_zone_path", obj.ARecordZonePath)
	d.Set("ptr_record_zone_path", obj.PtrRecordZonePath)
	if obj.Ttl != nil {
		d.Set("ttl", int(*obj.Ttl))
	}
	return nil
}

func resourceNsxtPolicyDnsRecordAutoConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	revision := int64(d.Get("revision").(int))
	obj := policyDnsRecordAutoConfigFromSchema(d)
	obj.Revision = &revision

	c := cliDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	_, err := c.Update(parents[0], parents[1], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("DnsAutoRecordConfig", id, err)
	}
	return resourceNsxtPolicyDnsRecordAutoConfigRead(d, m)
}

func resourceNsxtPolicyDnsRecordAutoConfigDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining DnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	err := c.Delete(parents[0], parents[1], id)
	if err != nil {
		return handleDeleteError("DnsAutoRecordConfig", id, err)
	}
	return nil
}
