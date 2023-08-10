/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyTransportZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransportZoneCreate,
		Read:   resourceNsxtPolicyTransportZoneRead,
		Update: resourceNsxtPolicyTransportZoneUpdate,
		Delete: resourceNsxtPolicyTransportZoneDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyTransportZoneImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"is_default": {
				Type:        schema.TypeBool,
				Description: "Indicates whether the transport zone is default",
				Optional:    true,
				Computed:    true,
			},
			"transport_type": {
				Type:         schema.TypeString,
				Description:  "Type of Transport Zone",
				Required:     true,
				ValidateFunc: validation.StringInSlice(policyTransportZoneTransportTypes, false),
				ForceNew:     true,
			},
			"uplink_teaming_policy_names": {
				Type:        schema.TypeList,
				Description: "Names of the switching uplink teaming policies that are supported by this transport zone.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"site_id": {
				Type:        schema.TypeString,
				Description: "ID of the site this Transport Zone belongs to",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"enforcement_point_id": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Transport Zone belongs to",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
		},
	}
}

func resourceNsxtPolicyTransportZoneExists(siteID, epID, tzID string, connector client.Connector) (bool, error) {
	var err error

	// Check site existence first
	siteClient := infra.NewSitesClient(connector)
	_, err = siteClient.Get(siteID)
	if err != nil {
		msg := fmt.Sprintf("Failed to read site %s", siteID)
		return false, logAPIError(msg, err)
	}

	// Check (ep, tz) existence. In case of ep not found, NSX returns BAD_REQUEST
	tzClient := enforcement_points.NewTransportZonesClient(connector)
	_, err = tzClient.Get(siteID, epID, tzID)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func policyTransportZonePatch(siteID, epID, tzID string, d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	isDefault := d.Get("is_default").(bool)
	transportType := d.Get("transport_type").(string)
	uplinkTeamingNames := getStringListFromSchemaList(d, "uplink_teaming_policy_names")

	if len(uplinkTeamingNames) > 0 && transportType != model.PolicyTransportZone_TZ_TYPE_VLAN_BACKED {
		// uplink_teaming_policy_names only valid for VLAN_BACKED TZ
		return fmt.Errorf("Cannot use uplink_teaming_policy_names with transport_type %s", transportType)
	}

	obj := model.PolicyTransportZone{
		DisplayName:              &displayName,
		Description:              &description,
		Tags:                     tags,
		IsDefault:                &isDefault,
		TzType:                   &transportType,
		UplinkTeamingPolicyNames: uplinkTeamingNames,
	}

	// Create the resource using PATCH
	tzClient := enforcement_points.NewTransportZonesClient(connector)
	_, err := tzClient.Patch(siteID, epID, tzID, obj)
	return err
}

func policyTransportZoneIDTuple(d *schema.ResourceData, m interface{}) (id, siteID, epID string, err error) {
	id = d.Id()
	if id == "" {
		err = fmt.Errorf("Error obtaining PolicyTransportZone ID")
		return
	}
	siteID = d.Get("site_id").(string)
	if siteID == "" {
		siteID = defaultSite
	}
	epID = d.Get("enforcement_point_id").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}
	return
}

func resourceNsxtPolicyTransportZoneCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	}
	siteID := d.Get("site_id").(string)
	if siteID == "" {
		siteID = defaultSite
	}
	epID := d.Get("enforcement_point_id").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}
	exists, err := resourceNsxtPolicyTransportZoneExists(siteID, epID, id, connector)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("resource with ID %s already exists", id)
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating PolicyTransportZone with ID %s under site %s enforcement point %s", id, siteID, epID)
	err = policyTransportZonePatch(siteID, epID, id, d, m)
	if err != nil {
		return handleCreateError("PolicyTransportZone", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransportZoneRead(d, m)
}

func resourceNsxtPolicyTransportZoneRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	tzClient := enforcement_points.NewTransportZonesClient(connector)

	id, siteID, epID, err := policyTransportZoneIDTuple(d, m)
	if err != nil {
		return err
	}

	obj, err := tzClient.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "PolicyTransportZone", id, err)
	}

	d.Set("site_id", siteID)
	d.Set("enforcement_point_id", epID)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("is_default", obj.IsDefault)
	d.Set("transport_type", obj.TzType)
	d.Set("UplinkTeamingPolicyNames", obj.UplinkTeamingPolicyNames)

	return nil
}

func resourceNsxtPolicyTransportZoneUpdate(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyTransportZoneIDTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing PolicyTransportZone with ID %s", id)
	err = policyTransportZonePatch(siteID, epID, id, d, m)
	if err != nil {
		return handleUpdateError("PolicyTransportZone", id, err)
	}

	return resourceNsxtPolicyTransportZoneRead(d, m)
}

func resourceNsxtPolicyTransportZoneDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	tzClient := enforcement_points.NewTransportZonesClient(connector)

	id, siteID, epID, err := policyTransportZoneIDTuple(d, m)
	if err != nil {
		return err
	}

	err = tzClient.Delete(siteID, epID, id)
	if err != nil {
		return handleDeleteError("PolicyTransportZone", id, err)
	}

	return nil
}

func resourceNsxtPolicyTransportZoneImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err == nil {
		epID, err := getParameterFromPolicyPath("/enforcement-points/", "/transport-zones/", importID)
		if err != nil {
			return nil, err
		}
		d.Set("enforcement_point_id", epID)
		siteID, err := getParameterFromPolicyPath("/sites/", "/enforcement-points/", importID)
		if err != nil {
			return nil, err
		}
		d.Set("site_id", siteID)
		return rd, nil
	} else if !errors.Is(err, ErrNotAPolicyPath) {
		return rd, err
	}

	var siteID, epID, id string
	if len(s) == 0 {
		siteID = defaultSite
		epID = getPolicyEnforcementPoint(m)
		id = s[0]
	} else if len(s) == 3 {
		siteID, epID, id = s[0], s[1], s[2]
	} else {
		return nil, fmt.Errorf("Please provide either transport-zone-id or <site-id>/<enforcement-point-id>/<transport-zone-id> as an input")
	}

	connector := getPolicyConnector(m)
	exists, err := resourceNsxtPolicyTransportZoneExists(siteID, epID, id, connector)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("PolicyTransportZone %s/%s/%s not found", siteID, epID, id)
	}
	d.Set("site_id", siteID)
	d.Set("enforcement_point_id", epID)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
