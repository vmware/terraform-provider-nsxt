// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliProjectDnsRecordsClient = projects.NewDnsRecordsClient

var policyDnsRecordPathExample = "/orgs/[org]/projects/[project]/dns-records/[dns-record]"

var dnsRecordTypeValues = []string{
	model.ProjectDnsRecord_RECORD_TYPE_A,
	model.ProjectDnsRecord_RECORD_TYPE_AAAA,
	model.ProjectDnsRecord_RECORD_TYPE_CNAME,
	model.ProjectDnsRecord_RECORD_TYPE_PTR,
	model.ProjectDnsRecord_RECORD_TYPE_NS,
}

func resourceNsxtPolicyDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDnsRecordCreate,
		Read:   resourceNsxtPolicyDnsRecordRead,
		Update: resourceNsxtPolicyDnsRecordUpdate,
		Delete: resourceNsxtPolicyDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "Policy DNS Record", getPolicyPathResourceImporter(policyDnsRecordPathExample)),
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"record_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name or host-octet label (for PTR records).",
			},
			"record_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(dnsRecordTypeValues, false),
				Description:  "DNS record type. Immutable after creation.",
			},
			"record_values": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "DNS record data values consistent with record_type.",
			},
			"zone_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Policy path to the ProjectDnsZone. Immutable after creation.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Time-To-Live in seconds for this DNS record (30-86400). Default: 300.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 address mapped by a PTR record. Only applicable when record_type is PTR.",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System-computed fully qualified domain name (read-only).",
			},
		},
	}
}

func resourceNsxtPolicyDnsRecordExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliProjectDnsRecordsClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type for DNS record")
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

func policyDnsRecordFromSchema(d *schema.ResourceData) model.ProjectDnsRecord {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	recordName := d.Get("record_name").(string)
	recordType := d.Get("record_type").(string)
	recordValues := getStringListFromSchemaList(d, "record_values")
	zonePath := d.Get("zone_path").(string)
	ttl := int64(d.Get("ttl").(int))

	obj := model.ProjectDnsRecord{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		RecordName:   &recordName,
		RecordType:   &recordType,
		RecordValues: recordValues,
		ZonePath:     &zonePath,
		Ttl:          &ttl,
	}

	if ipAddr, ok := d.GetOk("ip_address"); ok {
		ipStr := ipAddr.(string)
		obj.IpAddress = &ipStr
	}

	return obj
}

func resourceNsxtPolicyDnsRecordCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("Policy DNS Record resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyDnsRecordExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(sessionContext)
	obj := policyDnsRecordFromSchema(d)

	log.Printf("[INFO] Creating ProjectDnsRecord with ID %s", id)
	c := cliProjectDnsRecordsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS record")
	}
	err = c.Patch(parents[0], parents[1], id, obj)
	if err != nil {
		return handleCreateError("ProjectDnsRecord", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsRecordRead(d, m)
}

func resourceNsxtPolicyDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRecord ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliProjectDnsRecordsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS record")
	}
	obj, err := c.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "ProjectDnsRecord", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("record_name", obj.RecordName)
	d.Set("record_type", obj.RecordType)
	d.Set("record_values", obj.RecordValues)
	d.Set("zone_path", obj.ZonePath)
	if obj.Ttl != nil {
		d.Set("ttl", int(*obj.Ttl))
	}
	if obj.IpAddress != nil {
		d.Set("ip_address", obj.IpAddress)
	}
	if obj.Fqdn != nil {
		d.Set("fqdn", obj.Fqdn)
	}
	return nil
}

func resourceNsxtPolicyDnsRecordUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRecord ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	revision := int64(d.Get("revision").(int))
	obj := policyDnsRecordFromSchema(d)
	obj.Revision = &revision

	c := cliProjectDnsRecordsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS record")
	}
	_, err := c.Update(parents[0], parents[1], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("ProjectDnsRecord", id, err)
	}
	return resourceNsxtPolicyDnsRecordRead(d, m)
}

func resourceNsxtPolicyDnsRecordDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsRecord ID")
	}

	parents := getVpcParentsFromContext(sessionContext)
	c := cliProjectDnsRecordsClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS record")
	}
	err := c.Delete(parents[0], parents[1], id)
	if err != nil {
		return handleDeleteError("ProjectDnsRecord", id, err)
	}
	return nil
}
