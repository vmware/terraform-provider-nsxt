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
)

var lbTcpApplicationProfilePathExample = "/infra/lb-app-profiles/[profile]"

func resourceNsxtPolicyLBFastTcpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBTcpApplicationProfileCreate,
		Read:   resourceNsxtPolicyLBTcpApplicationProfileRead,
		Update: resourceNsxtPolicyLBTcpApplicationProfileUpdate,
		Delete: resourceNsxtPolicyLBTcpApplicationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(lbTcpApplicationProfilePathExample),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
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

func resourceNsxtPolicyLBTcpApplicationProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
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

	log.Printf("[INFO] Patching LBTcpProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBFastTcpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error converting LBTcpProfile %s", errs[0])
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbAppProfilesClient(sessionContext, connector)
	return client.Patch(id, dataValue.(*data.StructValue))
}

func resourceNsxtPolicyLBTcpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyLBAppProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtPolicyLBTcpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("LBTcpProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyLBTcpApplicationProfileRead(d, m)
}

func resourceNsxtPolicyLBTcpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBTcpProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbAppProfilesClient(sessionContext, connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "LBTcpProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBFastTcpProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("LBAppProfile with id %s is not of type LBTcpProfile %s", id, errs[0])
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

func resourceNsxtPolicyLBTcpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBTcpProfile ID")
	}

	err := resourceNsxtPolicyLBTcpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("LBTcpProfile", id, err)
	}

	return resourceNsxtPolicyLBTcpApplicationProfileRead(d, m)
}

func resourceNsxtPolicyLBTcpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyLBAppProfileDelete(d, m)
}
