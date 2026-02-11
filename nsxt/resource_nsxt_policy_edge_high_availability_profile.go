/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var cliSitesClient = infra.NewSitesClient
var cliEdgeClusterHighAvailabilityProfilesClient = enforcement_points.NewEdgeClusterHighAvailabilityProfilesClient

var policyEdgeHighAvailabilityProfileSchema = map[string]*metadata.ExtendedSchema{
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
	"bfd_probe_interval": {
		Schema: schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      500,
			ValidateFunc: validation.IntBetween(50, 60000),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "BfdProbeInterval",
		},
	},
	"bfd_allowed_hops": {
		Schema: schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      255,
			ValidateFunc: validation.IntBetween(1, 255),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "BfdAllowedHops",
		},
	},
	"bfd_declare_dead_multiple": {
		Schema: schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntBetween(2, 16),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "BfdDeclareDeadMultiple",
		},
	},
	"standby_relocation_config": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"standby_relocation_threshold": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "StandbyRelocationThreshold",
						},
					},
				},
			},
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "StandbyRelocationConfig",
			ReflectType:  reflect.TypeOf(model.StandbyRelocationConfig{}),
		},
	},
}

func resourceNsxtPolicyEdgeHighAvailabilityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyEdgeHighAvailabilityProfileCreate,
		Read:   resourceNsxtPolicyEdgeHighAvailabilityProfileRead,
		Update: resourceNsxtPolicyEdgeHighAvailabilityProfileUpdate,
		Delete: resourceNsxtPolicyEdgeHighAvailabilityProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyEdgeHighAvailabilityProfileImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(policyEdgeHighAvailabilityProfileSchema),
	}
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileExists(siteID, epID, id string, connector client.Connector) (bool, error) {
	// Check site existence first
	siteSessionContext := utl.SessionContext{ClientType: utl.Local}
	siteClient := cliSitesClient(siteSessionContext, connector)
	_, err := siteClient.Get(siteID)
	if err != nil {
		msg := fmt.Sprintf("failed to read site %s", siteID)
		return false, logAPIError(msg, err)
	}

	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliEdgeClusterHighAvailabilityProfilesClient(sessionContext, connector)
	_, err = client.Get(siteID, epID, id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileCreate(d *schema.ResourceData, m interface{}) error {
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
	exists, err := resourceNsxtPolicyEdgeHighAvailabilityProfileExists(siteID, epID, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.PolicyEdgeHighAvailabilityProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, policyEdgeHighAvailabilityProfileSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating PolicyEdgeHighAvailabilityProfile with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliEdgeClusterHighAvailabilityProfilesClient(sessionContext, connector)
	err = client.Patch(siteID, epID, id, obj)
	if err != nil {
		return handleCreateError("PolicyEdgeHighAvailabilityProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileRead(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliEdgeClusterHighAvailabilityProfilesClient(sessionContext, connector)
	obj, err := client.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "PolicyEdgeHighAvailabilityProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, policyEdgeHighAvailabilityProfileSchema, "", nil)
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.PolicyEdgeHighAvailabilityProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, policyEdgeHighAvailabilityProfileSchema, "", nil); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliEdgeClusterHighAvailabilityProfilesClient(sessionContext, connector)
	_, err = client.Update(siteID, epID, id, obj)
	if err != nil {
		return handleUpdateError("PolicyEdgeHighAvailabilityProfile", id, err)
	}

	return resourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileDelete(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := cliEdgeClusterHighAvailabilityProfilesClient(sessionContext, connector)
	err = client.Delete(siteID, epID, id)

	if err != nil {
		return handleDeleteError("PolicyEdgeHighAvailabilityProfile", id, err)
	}

	return nil
}

func resourceNsxtPolicyEdgeHighAvailabilityProfileImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/edge-cluster-high-availability-profiles/", importID)
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
