// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnssvcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/dns_services"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliProjectDnsZonesClient = dnssvcs.NewZonesClient

var policyDnsZonePathExample = "/orgs/[org]/projects/[project]/dns-services/[dns-service]"

func resourceNsxtPolicyDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDnsZoneCreate,
		Read:   resourceNsxtPolicyDnsZoneRead,
		Update: resourceNsxtPolicyDnsZoneUpdate,
		Delete: resourceNsxtPolicyDnsZoneDelete,
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
			"dns_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain name for this zone (e.g. 'example.com'). Immutable after creation.",
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
				Description:  "Optional policy path to a single VPC for split-horizon DNS. When unset, all VPCs can resolve this zone.",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Default Time-To-Live in seconds for DNS records in this zone (30-86400).",
			},
			"soa": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_nameserver": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "FQDN of the primary nameserver (must end with a trailing dot).",
						},
						"responsible_party": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Zone administrator email in DNS format (must end with a trailing dot).",
						},
						"serial_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Zone serial number.",
						},
						"refresh_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Interval in seconds for secondary nameservers to refresh zone data.",
						},
						"retry_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Interval in seconds for secondary nameservers to retry a failed refresh.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Time in seconds after which secondary nameservers stop serving zone data.",
						},
						"negative_cache_ttl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "TTL in seconds for NXDOMAIN responses (0-86400).",
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyDnsZoneExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, policyDnsZonePathExample)
	if pathErr != nil {
		return false, pathErr
	}
	c := cliProjectDnsZonesClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type for DNS zone")
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

func policyDnsZoneFromSchema(d *schema.ResourceData) model.ProjectDnsZone {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	dnsDomain := d.Get("dns_domain_name").(string)
	ttl := int64(d.Get("ttl").(int))

	obj := model.ProjectDnsZone{
		DisplayName:   &displayName,
		Description:   &description,
		Tags:          tags,
		DnsDomainName: &dnsDomain,
		Ttl:           &ttl,
	}

	if scope, ok := d.GetOk("scope"); ok {
		scopeStr := scope.(string)
		obj.Scope = &scopeStr
	}

	soaList := d.Get("soa").([]interface{})
	if len(soaList) > 0 {
		soaMap := soaList[0].(map[string]interface{})
		soa := &model.ProjectDnsZoneSoa{}
		if v, ok := soaMap["primary_nameserver"].(string); ok && v != "" {
			soa.PrimaryNameserver = &v
		}
		if v, ok := soaMap["responsible_party"].(string); ok && v != "" {
			soa.ResponsibleParty = &v
		}
		if v, ok := soaMap["serial_number"].(int); ok && v > 0 {
			sn := int64(v)
			soa.SerialNumber = &sn
		}
		if v, ok := soaMap["refresh_interval"].(int); ok && v > 0 {
			ri := int64(v)
			soa.RefreshInterval = &ri
		}
		if v, ok := soaMap["retry_interval"].(int); ok && v > 0 {
			ri := int64(v)
			soa.RetryInterval = &ri
		}
		if v, ok := soaMap["expire_time"].(int); ok && v > 0 {
			et := int64(v)
			soa.ExpireTime = &et
		}
		if v, ok := soaMap["negative_cache_ttl"].(int); ok && v >= 0 {
			nc := int64(v)
			soa.NegativeCacheTtl = &nc
		}
		if !reflect.DeepEqual(soa, &model.ProjectDnsZoneSoa{}) {
			obj.Soa = soa
		}
	}

	return obj
}

func policyDnsZoneToSchema(d *schema.ResourceData, obj model.ProjectDnsZone) {
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)
	d.Set("nsx_id", obj.Id)
	d.Set("dns_domain_name", obj.DnsDomainName)
	if obj.Ttl != nil {
		d.Set("ttl", int(*obj.Ttl))
	}
	if obj.Scope != nil {
		d.Set("scope", obj.Scope)
	}
	if obj.Soa != nil {
		soaMap := map[string]interface{}{}
		if obj.Soa.PrimaryNameserver != nil {
			soaMap["primary_nameserver"] = *obj.Soa.PrimaryNameserver
		}
		if obj.Soa.ResponsibleParty != nil {
			soaMap["responsible_party"] = *obj.Soa.ResponsibleParty
		}
		if obj.Soa.SerialNumber != nil {
			soaMap["serial_number"] = int(*obj.Soa.SerialNumber)
		}
		if obj.Soa.RefreshInterval != nil {
			soaMap["refresh_interval"] = int(*obj.Soa.RefreshInterval)
		}
		if obj.Soa.RetryInterval != nil {
			soaMap["retry_interval"] = int(*obj.Soa.RetryInterval)
		}
		if obj.Soa.ExpireTime != nil {
			soaMap["expire_time"] = int(*obj.Soa.ExpireTime)
		}
		if obj.Soa.NegativeCacheTtl != nil {
			soaMap["negative_cache_ttl"] = int(*obj.Soa.NegativeCacheTtl)
		}
		d.Set("soa", []interface{}{soaMap})
	}
}

func dnsZoneParentsFromPath(parentPath string) (org, project, dnsServiceID string, err error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, policyDnsZonePathExample)
	if pathErr != nil {
		return "", "", "", pathErr
	}
	return parents[0], parents[1], parents[2], nil
}

func resourceNsxtPolicyDnsZoneCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("Policy DNS Zone resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyDnsZoneExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsZoneParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	obj := policyDnsZoneFromSchema(d)
	log.Printf("[INFO] Creating ProjectDnsZone with ID %s under %s", id, parentPath)

	c := cliProjectDnsZonesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS zone")
	}
	err = c.Patch(org, project, dnsServiceID, id, obj)
	if err != nil {
		return handleCreateError("ProjectDnsZone", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyDnsZoneRead(d, m)
}

func resourceNsxtPolicyDnsZoneRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsZone ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsZoneParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliProjectDnsZonesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS zone")
	}
	obj, err := c.Get(org, project, dnsServiceID, id)
	if err != nil {
		return handleReadError(d, "ProjectDnsZone", id, err)
	}
	policyDnsZoneToSchema(d, obj)
	return nil
}

func resourceNsxtPolicyDnsZoneUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsZone ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsZoneParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	revision := int64(d.Get("revision").(int))
	obj := policyDnsZoneFromSchema(d)
	obj.Revision = &revision

	c := cliProjectDnsZonesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS zone")
	}
	_, err := c.Update(org, project, dnsServiceID, id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("ProjectDnsZone", id, err)
	}
	return resourceNsxtPolicyDnsZoneRead(d, m)
}

func resourceNsxtPolicyDnsZoneDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining ProjectDnsZone ID")
	}

	parentPath := d.Get("parent_path").(string)
	org, project, dnsServiceID, pathErr := dnsZoneParentsFromPath(parentPath)
	if pathErr != nil {
		return pathErr
	}

	c := cliProjectDnsZonesClient(sessionContext, connector)
	if c == nil {
		return fmt.Errorf("unsupported client type for DNS zone")
	}
	err := c.Delete(org, project, dnsServiceID, id)
	if err != nil {
		return handleDeleteError("ProjectDnsZone", id, err)
	}
	return nil
}
