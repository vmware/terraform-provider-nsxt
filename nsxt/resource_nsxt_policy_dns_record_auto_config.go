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

var cliProjectDnsAutoRecordConfigsClient = projects.NewDnsAutoRecordConfigsClient

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
			"zone_path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to a locally-owned ProjectDnsZone where auto-created A records are placed.",
			},
			"naming_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "{vm_name}",
				Description: "Template for generating DNS record names from VM attributes. Supports {vm_name}, {ip}, {index}, {vpc_name}. Default: '{vm_name}'.",
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
	c := cliProjectDnsAutoRecordConfigsClient(sessionContext, connector)
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

func policyDnsRecordAutoConfigFromSchema(d *schema.ResourceData) model.ProjectDnsAutoRecordConfig {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	ipBlockPath := d.Get("ip_block_path").(string)
	zonePath := d.Get("zone_path").(string)
	namingPattern := d.Get("naming_pattern").(string)
	ttl := int64(d.Get("ttl").(int))

	return model.ProjectDnsAutoRecordConfig{
		DisplayName:   &displayName,
		Description:   &description,
		Tags:          tags,
		IpBlockPath:   &ipBlockPath,
		ZonePath:      &zonePath,
		NamingPattern: &namingPattern,
		Ttl:           &ttl,
	}
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

	log.Printf("[INFO] Creating ProjectDnsAutoRecordConfig with ID %s", id)
	c := cliProjectDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	err = c.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("ProjectDnsAutoRecordConfig", id, err)
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
		return fmt.Errorf("error obtaining ProjectDnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliProjectDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	obj, err := c.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "ProjectDnsAutoRecordConfig", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("ip_block_path", obj.IpBlockPath)
	d.Set("zone_path", obj.ZonePath)
	if obj.NamingPattern != nil {
		d.Set("naming_pattern", obj.NamingPattern)
	}
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
		return fmt.Errorf("error obtaining ProjectDnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	revision := int64(d.Get("revision").(int))
	obj := policyDnsRecordAutoConfigFromSchema(d)
	obj.Revision = &revision

	c := cliProjectDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	_, err := c.Update(parents[0], parents[1], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("ProjectDnsAutoRecordConfig", id, err)
	}
	return resourceNsxtPolicyDnsRecordAutoConfigRead(d, m)
}

func resourceNsxtPolicyDnsRecordAutoConfigDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsAutoRecordConfig ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliProjectDnsAutoRecordConfigsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS auto record config")
	}
	err := c.Delete(parents[0], parents[1], id)
	if err != nil {
		return handleDeleteError("ProjectDnsAutoRecordConfig", id, err)
	}
	return nil
}
