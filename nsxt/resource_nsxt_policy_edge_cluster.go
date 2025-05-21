/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var policyEdgeClusterSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"site_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Description:  "Path to the site this Host Transport Node belongs to",
			Optional:     true,
			ForceNew:     true,
			Default:      defaultInfraSitePath,
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType: "string",
			Skip:       true,
		},
	},
	"enforcement_point": {
		Schema: schema.Schema{
			Type:        schema.TypeString,
			Description: "ID of the enforcement point this Host Transport Node belongs to",
			Optional:    true,
			ForceNew:    true,
			Default:     "default",
		},
		Metadata: metadata.Metadata{
			SchemaType: "string",
			Skip:       true,
		},
	},
	"inter_site_forwarding_enabled": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SdkFieldName: "InterSiteForwardingEnabled",
		},
	},
	"edge_cluster_profile_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "EdgeClusterProfile",
		},
	},
	"allocation_rule": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"action_based_on_failure_domain_enabled": {
						Schema: schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						Metadata: metadata.Metadata{
							SchemaType: "bool",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			Skip: true,
		},
	},
	"policy_edge_node": {
		Schema: schema.Schema{
			Type: schema.TypeSet,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"edge_transport_node_path": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validatePolicyPath(),
							Required:     true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "EdgeTransportNodePath",
						},
					},
					"id": {
						Schema: schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "Id",
							OmitIfEmpty:  true,
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "set",
			SdkFieldName: "PolicyEdgeNodes",
			ReflectType:  reflect.TypeOf(model.PolicyEdgeClusterMember{}),
		},
	},
	"password_managed_by_vcf": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "PasswordManagedByVcf",
		},
	},
}

func resourceNsxtPolicyEdgeCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyEdgeClusterCreate,
		Read:   resourceNsxtPolicyEdgeClusterRead,
		Update: resourceNsxtPolicyEdgeClusterUpdate,
		Delete: resourceNsxtPolicyEdgeClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyEdgeClusterImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(policyEdgeClusterSchema),
	}
}

func resourceNsxtPolicyEdgeClusterExists(siteID, epID, id string, connector client.Connector) (bool, error) {
	// Check site existence first
	siteClient := infra.NewSitesClient(connector)
	_, err := siteClient.Get(siteID)
	if err != nil {
		msg := fmt.Sprintf("failed to read site %s", siteID)
		return false, logAPIError(msg, err)
	}

	client := enforcement_points.NewEdgeClustersClient(connector)
	_, err = client.Get(siteID, epID, id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyEdgeClusterCreate(d *schema.ResourceData, m interface{}) error {

	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	}
	sitePath := d.Get("site_path").(string)
	siteID := getResourceIDFromResourcePath(sitePath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining Site ID from site path %s", sitePath)
	}
	epID := d.Get("enforcement_point").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}

	connector := getPolicyConnector(m)
	exists, err := resourceNsxtPolicyEdgeClusterExists(siteID, epID, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	converter := bindings.NewTypeConverter()
	var allocationRules []model.AllocationRule
	allocRules := d.Get("allocation_rule").([]interface{})
	for _, r := range allocRules {
		rule := r.(map[string]interface{})
		var allocRule model.AllocationRule
		if rule["action_based_on_failure_domain_enabled"].(bool) {
			enabled := true
			fdRule := model.AllocationBasedOnFailureDomain{
				Enabled:    &enabled,
				ActionType: model.AllocationRuleAction_ACTION_TYPE_ALLOCATIONBASEDONFAILUREDOMAIN,
			}
			dataValue, errs := converter.ConvertToVapi(fdRule, model.AllocationBasedOnFailureDomainBindingType())
			if errs != nil {
				return handleCreateError("PolicyEdgeCluster", id, errs[0])
			}
			allocRule = model.AllocationRule{
				Action: dataValue.(*data.StructValue),
			}
		}
		allocationRules = append(allocationRules, allocRule)
	}

	obj := model.PolicyEdgeCluster{
		DisplayName:     &displayName,
		Description:     &description,
		AllocationRules: allocationRules,
		Tags:            tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, policyEdgeClusterSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating PolicyEdgeCluster with ID %s", id)

	client := enforcement_points.NewEdgeClustersClient(connector)
	err = client.Patch(siteID, epID, id, obj)
	if err != nil {
		return handleCreateError("PolicyEdgeCluster", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyEdgeClusterRead(d, m)
}

func resourceNsxtPolicyEdgeClusterRead(d *schema.ResourceData, m interface{}) error {

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeClustersClient(connector)
	obj, err := client.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "PolicyEdgeCluster", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	converter := bindings.NewTypeConverter()
	var allocationRules []interface{}
	for _, allocRule := range obj.AllocationRules {
		rule := make(map[string]interface{})
		baseAction, errs := converter.ConvertToGolang(allocRule.Action, model.AllocationRuleActionBindingType())
		if errs != nil {
			return handleReadError(d, "PolicyEdgeCluster", id, errs[0])
		}
		switch baseAction.(model.AllocationRuleAction).ActionType {
		case model.AllocationRuleAction_ACTION_TYPE_ALLOCATIONBASEDONFAILUREDOMAIN:
			action, errs := converter.ConvertToGolang(allocRule.Action, model.AllocationBasedOnFailureDomainBindingType())
			if errs != nil {
				return handleReadError(d, "PolicyEdgeCluster", id, errs[0])
			}
			if *action.(model.AllocationBasedOnFailureDomain).Enabled {
				rule["action_based_on_failure_domain_enabled"] = true
			}

		default:
			return fmt.Errorf("allocation rule action %s is not supported", baseAction.(model.AllocationRuleAction).ActionType)
		}
		allocationRules = append(allocationRules, rule)
	}
	d.Set("allocation_rule", allocationRules)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, policyEdgeClusterSchema, "", nil)
}

func resourceNsxtPolicyEdgeClusterUpdate(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	converter := bindings.NewTypeConverter()
	var allocationRules []model.AllocationRule
	allocRules := d.Get("allocation_rules")
	if allocRules != nil {
		for _, r := range allocRules.([]interface{}) {
			if r != nil {
				rule := r.(map[string]interface{})
				var allocRule model.AllocationRule
				if rule["action_based_on_failure_domain_enabled"].(bool) {
					enabled := true
					fdRule := model.AllocationBasedOnFailureDomain{
						Enabled:    &enabled,
						ActionType: model.AllocationRuleAction_ACTION_TYPE_ALLOCATIONBASEDONFAILUREDOMAIN,
					}
					dataValue, errs := converter.ConvertToVapi(fdRule, model.AllocationBasedOnFailureDomainBindingType())
					if errs != nil {
						return handleCreateError("PolicyEdgeCluster", id, errs[0])
					}
					allocRule = model.AllocationRule{
						Action: dataValue.(*data.StructValue),
					}
				}
				allocationRules = append(allocationRules, allocRule)
			}
		}
	}

	obj := model.PolicyEdgeCluster{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		Revision:        &revision,
		AllocationRules: allocationRules,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, policyEdgeClusterSchema, "", nil); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeClustersClient(connector)
	_, err = client.Update(siteID, epID, id, obj)
	if err != nil {
		return handleUpdateError("PolicyEdgeCluster", id, err)
	}

	return resourceNsxtPolicyEdgeClusterRead(d, m)
}

func resourceNsxtPolicyEdgeClusterDelete(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeClustersClient(connector)
	err = client.Delete(siteID, epID, id)

	if err != nil {
		return handleDeleteError("PolicyEdgeCluster", id, err)
	}

	return nil
}

func resourceNsxtPolicyEdgeClusterImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/edge-clusters/", importID)
	if err != nil {
		return nil, err
	}
	d.Set("enforcement_point", epID)
	sitePath, err := getSitePathFromChildResourcePath(importID)
	if err != nil {
		return rd, err
	}
	d.Set("site_path", sitePath)

	return rd, nil
}
