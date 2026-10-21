// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func resourceNsxtVpcLBFastTcpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBTcpApplicationProfileCreate,
		Read:   resourceNsxtVpcLBTcpApplicationProfileRead,
		Update: resourceNsxtVpcLBTcpApplicationProfileUpdate,
		Delete: resourceNsxtVpcLBTcpApplicationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Fast TCP Application Profile", getVpcPathResourceImporter(vpcLbAppProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"idle_timeout": {
				Type:         schema.TypeInt,
				Description:  "how long an idle TCP connection in ESTABLISHED state should be kept for this application before cleaning up",
				Optional:     true,
				Default:      1800,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"close_timeout": {
				Type:         schema.TypeInt,
				Description:  "how long a closing TCP connection (both FINs received or a RST is received) should be kept for this application before cleaning up",
				Optional:     true,
				Default:      8,
				ValidateFunc: validation.IntBetween(1, 15),
			},
			"ha_flow_mirroring_enabled": {
				Type:        schema.TypeBool,
				Description: "If enabled, all the flows to the bounded virtual server are mirrored to the standby node",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceNsxtVpcLBTcpApplicationProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	haFlowMirroringEnabled := d.Get("ha_flow_mirroring_enabled").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	closeTimeout := int64(d.Get("close_timeout").(int))
	resourceType := model.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE
	obj := model.LBFastTcpProfile{
		DisplayName:            &displayName,
		Description:            &description,
		Tags:                   tags,
		HaFlowMirroringEnabled: &haFlowMirroringEnabled,
		IdleTimeout:            &idleTimeout,
		CloseTimeout:           &closeTimeout,
		ResourceType:           resourceType,
	}

	log.Printf("[INFO] Patching VPC LBTcpProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBFastTcpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error converting VPC LBTcpProfile %s", errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBTcpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Fast TCP Application Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBAppProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBTcpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBTcpProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBTcpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBTcpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBTcpProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBTcpProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBFastTcpProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("VPC LBAppProfile with id %s is not of type LBTcpProfile %s", id, errs[0])
	}
	lbTcpProfile := baseObj.(model.LBFastTcpProfile)

	d.Set("display_name", lbTcpProfile.DisplayName)
	d.Set("description", lbTcpProfile.Description)
	setPolicyTagsInSchema(d, lbTcpProfile.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbTcpProfile.Path)
	d.Set("revision", lbTcpProfile.Revision)
	d.Set("ha_flow_mirroring_enabled", lbTcpProfile.HaFlowMirroringEnabled)
	d.Set("idle_timeout", lbTcpProfile.IdleTimeout)
	d.Set("close_timeout", lbTcpProfile.CloseTimeout)

	return nil
}

func resourceNsxtVpcLBTcpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBTcpProfile ID")
	}

	err := resourceNsxtVpcLBTcpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBTcpProfile", id, err)
	}

	return resourceNsxtVpcLBTcpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBTcpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBAppProfileDelete(d, m)
}
