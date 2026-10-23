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

func resourceNsxtVpcLBFastUdpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBUdpApplicationProfileCreate,
		Read:   resourceNsxtVpcLBUdpApplicationProfileRead,
		Update: resourceNsxtVpcLBUdpApplicationProfileUpdate,
		Delete: resourceNsxtVpcLBUdpApplicationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Fast UDP Application Profile", getVpcPathResourceImporter(vpcLbAppProfilePathExample)),
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
				Description:  "how long an idle UDP connection should be kept for this application before cleaning up",
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"flow_mirroring_enabled": {
				Type:        schema.TypeBool,
				Description: "If enabled, all the flows to the bounded virtual server are mirrored to the standby node",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceNsxtVpcLBUdpApplicationProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	FlowMirroringEnabled := d.Get("flow_mirroring_enabled").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	resourceType := model.LBAppProfile_RESOURCE_TYPE_LBFASTUDPPROFILE
	obj := model.LBFastUdpProfile{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		FlowMirroringEnabled: &FlowMirroringEnabled,
		IdleTimeout:          &idleTimeout,
		ResourceType:         resourceType,
	}

	log.Printf("[INFO] Patching VPC LBUdpProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBFastUdpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error converting VPC LBUdpProfile %s", errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBUdpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Fast UDP Application Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBAppProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBUdpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBUdpProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBUdpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBUdpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBUdpProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBUdpProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBFastUdpProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("VPC LBAppProfile with id %s is not of type LBUdpProfile %s", id, errs[0])
	}
	lbUdpProfile := baseObj.(model.LBFastUdpProfile)

	d.Set("display_name", lbUdpProfile.DisplayName)
	d.Set("description", lbUdpProfile.Description)
	setPolicyTagsInSchema(d, lbUdpProfile.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbUdpProfile.Path)
	d.Set("revision", lbUdpProfile.Revision)
	d.Set("flow_mirroring_enabled", lbUdpProfile.FlowMirroringEnabled)
	d.Set("idle_timeout", lbUdpProfile.IdleTimeout)

	return nil
}

func resourceNsxtVpcLBUdpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBUdpProfile ID")
	}

	err := resourceNsxtVpcLBUdpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBUdpProfile", id, err)
	}

	return resourceNsxtVpcLBUdpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBUdpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBAppProfileDelete(d, m)
}
